package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/storage"
)

// RabbitMQ - structure that contains rabbit queue and channel
type RabbitMQ struct {
	Queues  map[string]amqp.Queue
	Channel *amqp.Channel
	Logger  logger.Logger
	Cfg     config.Config
	Storage storage.StorageI
}

// New - returns new RabbitMQ queue and channel
func New(cfg *config.Config, log logger.Logger, storage storage.StorageI) (*RabbitMQ, error) {
	log.Info(fmt.Sprintf("Dialing to rabbitmq host with host:%s user:%s", cfg.RabbitMqHost, cfg.RabbitMqUser))

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			cfg.RabbitMqUser,
			cfg.RabbitMqPassword,
			cfg.RabbitMqHost,
			cfg.RabbitMqPort,
		),
	)

	if err != nil {
		log.Error("Error while connecting to rabbitmq err:" + err.Error())
		return &RabbitMQ{}, err
	}

	log.Info("RabbitMQ connection is created...")

	channel, err := conn.Channel()
	if err != nil {
		log.Error("Error while connecting to channel err: " + err.Error())
		return &RabbitMQ{}, err
	}

	log.Info("RabbitMQ channel is created...")

	listen, err := channel.QueueDeclare(
		cfg.ListenQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Error while declaring queue err:" + err.Error())
		return &RabbitMQ{}, err
	}

	write, err := channel.QueueDeclare(
		cfg.WriteQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Error while declaring queue err:" + err.Error())
		return &RabbitMQ{}, err
	}

	return &RabbitMQ{
		Queues: map[string]amqp.Queue{
			cfg.ListenQueue: listen,
			cfg.WriteQueue:  write,
		},
		Channel: channel,
		Logger:  log,
		Cfg:     *cfg,
		Storage: storage,
	}, nil
}

func (r *RabbitMQ) PublishPipeline(req *models.PipelineRabbitMq) error {
	jsonByte, err := json.MarshalIndent(req, "", "    ")
	if err != nil {
		r.Logger.Error("Error while publishing new pipeline status")
		return err
	}

	err = r.Channel.Publish(
		"",
		r.Queues[r.Cfg.WriteQueue].Name,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonByte,
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), `"channel/connection is not open"`) {
			panic(`panic reason: "channel/connection is not open"`)
		}
		r.Logger.Error("Error while publishing the message err:" + err.Error())
		return err
	}

	return nil
}

func (r *RabbitMQ) StartListening() {
	for {
		msgs, err := r.Channel.Consume(
			r.Queues[r.Cfg.ListenQueue].Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			r.Logger.Error("Error while consuming messages err:" + err.Error())
			err = r.Reconnect()
			if err != nil {
				panic("couldn't reconnect to rabbitmq")
			}
			time.Sleep(time.Second * 5)
			continue
		}

		for data := range msgs {
			r.Logger.Info("Message is recieved")
			err = r.UpdatePipelineStatus(data.Body)
			if err == nil {
				err = data.Ack(true)
				if err != nil {
					r.Logger.Error("error while acknoledging message")
					break
				}
			} else {
				r.Logger.Error("error while udpatin pipelines status err:" + err.Error())
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func (r *RabbitMQ) Reconnect() error {
	r.Logger.Info("reconnecting to rabbitmq")

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			r.Cfg.RabbitMqUser,
			r.Cfg.RabbitMqPassword,
			r.Cfg.RabbitMqHost,
			r.Cfg.RabbitMqPort,
		),
	)
	if err != nil {
		return err
	}

	r.Channel, err = conn.Channel()
	if err != nil {
		r.Logger.Error("Error while connecting to channel err: " + err.Error())
		return err
	}

	r.Logger.Info("RabbitMQ channel is recreated...")

	listen, err := r.Channel.QueueDeclare(
		r.Cfg.ListenQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		r.Logger.Error("Error while declaring queue err:" + err.Error())
		return err
	}

	write, err := r.Channel.QueueDeclare(
		r.Cfg.WriteQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		r.Logger.Error("Error while declaring queue err:" + err.Error())
		return err
	}

	r.Queues = map[string]amqp.Queue{
		r.Cfg.ListenQueue: listen,
		r.Cfg.WriteQueue:  write,
	}

	return nil
}

func (r *RabbitMQ) UpdatePipelineStatus(body []byte) error {
	req := models.UpdatePipelineStatus{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		r.Logger.Error("error while unmarshaling request")
		return err
	}
	resolutions := []string{}
	for _, e := range req.Resolutions {
		resolutions = append(resolutions, e.Resolution)
	}

	// updating pipeline status
	res, err := r.Storage.Postgres().PipelineUpdate(context.Background(), &models.PipelineUpdateReq{
		ID:                         req.Id,
		Stage:                      req.Stage,
		StageStatus:                req.Status,
		FailDescription:            req.FailDescription,
		PreparationDurationSeconds: float64(req.PreparationDuration) / 1000,
		TranscodeDurationSecs:      float64(req.TranscodeDuration) / 1000,
		UploadDurationSecs:         float64(req.UploadDuration) / 1000,
		VideoDuration:              float32(req.VideoDuration),
		ResolutionsString:          resolutions,
	})
	if err != nil {
		return err
	}

	if res.StageStatus == "fail" || (res.StageStatus == "success" && res.Stage == "upload") {
		webhooks, err := r.Storage.Postgres().WebhookFind(context.Background(), &models.WebhookFindReq{
			Page:      1,
			Limit:     1,
			ProjectId: res.ProjectID,
		})
		if err != nil || len(webhooks.Webhooks) == 0 {
			return err
		}
		webhook := webhooks.Webhooks[0]
		if webhook.Active {
			webhookStatus := false

			project, err := r.Storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{ID: webhook.ProjectID})
			if err != nil {
				return err
			}
			res.Resolutions = req.Resolutions

			bodybyte, err := etc.DoRequest(webhook.URL, "POST", res, project.AccessKey, project.SecretKey)
			if err == nil {
				body := models.PipelineWebhookResponse{}
				err = json.Unmarshal(bodybyte, &body)
				if err == nil {
					webhookStatus = body.Success
				}

				_, err = r.Storage.Postgres().PipelineUpdate(context.Background(), &models.PipelineUpdateReq{
					ID:                req.Id,
					WebhookStatus:     webhookStatus,
					WebhookRetryCount: res.WebhookRetryCount + 1,
					WebhookLastRetry:  time.Now(),
				})
				if err != nil {
					return err
				}
			} else {
				return err
			}

		}

	}
	return nil
}

// this funcition retries webhooks in every 1 minute for max 10 times per pipeline.
func (h *RabbitMQ) Retry() {
	var (
		page  = 1
		limit = 50
	)

	for {
		pipelines, err := h.Storage.Postgres().PipelinesFind(context.Background(), &models.PipelineFindReq{
			Page:             page,
			Limit:            limit,
			OrderByCreatedAt: 1,
			WebhookStatus:    -1,
		})
		if err != nil {
			h.Logger.Error("Error while gettign pipelines")
			continue
		}
		page++

		for _, pipeline := range pipelines.Pipelines {
			webhooks, err := h.Storage.Postgres().WebhookFind(context.Background(), &models.WebhookFindReq{
				ProjectId:        pipeline.ProjectID,
				OrderByCreatedAt: 1,
				Page:             1,
				Limit:            1,
			})
			fmt.Println(pipeline.ProjectID)
			if err != nil || len(webhooks.Webhooks) == 0 {
				h.Logger.Error("Error while getting webhook", err)
				continue
			}
			webhook := webhooks.Webhooks[0]

			if !webhook.Active {
				_, err = h.Storage.Postgres().PipelineUpdate(context.TODO(), &models.PipelineUpdateReq{
					ID:                pipeline.ID,
					WebhookRetryCount: -1,
					WebhookLastRetry:  time.Now(),
				})
				if err != nil {
					h.Logger.Error("Error while updating status of pipeline")
					continue
				}
			}

			project, err := h.Storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{ID: webhook.ProjectID})
			if err != nil {
				h.Logger.Error("Error while getting project")
				continue
			}

			bodybyte, err := etc.DoRequest(webhook.URL, "POST", pipeline, project.AccessKey, project.SecretKey)
			if err == nil {
				body := models.PipelineWebhookResponse{}
				json.Unmarshal(bodybyte, &body)
				webhookStatus := body.Success

				_, err = h.Storage.Postgres().PipelineUpdate(context.Background(), &models.PipelineUpdateReq{
					ID:                pipeline.ID,
					WebhookStatus:     webhookStatus,
					WebhookRetryCount: pipeline.WebhookRetryCount + 1,
					WebhookLastRetry:  time.Now(),
				})
				if err != nil {
					h.Logger.Error("Error while getting project")
				}
			} else {
				h.Logger.Error("Error while getting project")
			}

		}

		if len(pipelines.Pipelines) < limit {
			break
		}
	}
}
