package main

import (
	"time"

	"gitlab.com/transcodeuz/transcode-rest/api"
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/db"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/pkg/rabbitmq"
	"gitlab.com/transcodeuz/transcode-rest/storage"
)

func main() {
	// if os.Getenv("SECRET_KEY") != "ADFaAdfadf2ds2w11f3r5tq" {
	// 	fmt.Println("You are lying!!!")
	// 	return
	// }
	config := config.Load()
	logger := logger.New(config.LogLevel)

	posgtres, err := db.New(config)
	if err != nil {
		logger.Error("error while connecting to postgresql")
		return
	}

	storage := storage.New(posgtres, logger, config)

	rabbit, err := rabbitmq.New(&config, *logger, storage)
	if err != nil {
		logger.Error("error while connecting to rabbti mq.")
		return
	}

	go rabbit.StartListening()
	go func(r *rabbitmq.RabbitMQ) {
		logger.Info("Cron job started: ")
		for {
			r.Retry()
			time.Sleep(time.Minute)
		}
	}(rabbit)

	engine := api.New(logger, config, storage, rabbit)

	err = engine.Run(":" + config.HTTPPort)
	if err != nil {
		logger.Error("error while running rest server")
		return
	}
}
