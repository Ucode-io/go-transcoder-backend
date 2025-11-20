package v1

import (
	t "gitlab.com/transcodeuz/transcode-rest/api/tokens"
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/pkg/rabbitmq"
	"gitlab.com/transcodeuz/transcode-rest/storage"
	"gitlab.com/transcodeuz/transcode-rest/storage/redisrepo"
)

type handlerV1 struct {
	log        *logger.Logger
	cfg        config.Config
	storage    storage.StorageI
	jwthandler t.JWTHandler
	redis      redisrepo.InMemoryStorageI
	rabbit     *rabbitmq.RabbitMQ
}

type HandlerV1Config struct {
	Logger     *logger.Logger
	Cfg        config.Config
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redisrepo.InMemoryStorageI
	Rabbit     *rabbitmq.RabbitMQ
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:        c.Logger,
		cfg:        c.Cfg,
		storage:    c.Postgres,
		jwthandler: c.JWTHandler,
		redis:      c.Redis,
		rabbit:     c.Rabbit,
	}
}
