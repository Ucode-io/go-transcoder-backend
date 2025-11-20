package storage

import (
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/db"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/storage/postgres"
)

type StorageI interface {
	Postgres() postgres.PostgresI
}

type StoragePg struct {
	postgres postgres.PostgresI
}

// NewStoragePg
func New(db *db.Postgres, log *logger.Logger, cfg config.Config) StorageI {
	return &StoragePg{
		postgres: postgres.New(db, log, cfg),
	}
}

func (s *StoragePg) Postgres() postgres.PostgresI {
	return s.postgres
}
