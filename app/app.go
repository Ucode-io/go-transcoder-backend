package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/transcodeuz/transcode-rest/api"
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/db"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/storage"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func LoggerInit(config config.Config) *logger.Logger {
	fmt.Println("logInit")
	log := logger.New(config.LogLevel)
	return log
}

func ServerStart(r *gin.Engine, cfg config.Config) {
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", err)
		panic(err)
	}
}

type App struct {
	engine *fx.App
}

// Start starts app with context spesified
func (a *App) Start(ctx context.Context) {
	a.engine.Start(ctx)
}

// Run starts the application, blocks on the signals channel, and then gracefully shuts the application down
func (a *App) Run() {
	a.engine.Run()
}

// New returns fx app
func New() App {

	engine := fx.New(
		fx.Provide(
			config.Load,
			LoggerInit,
			db.New,
			storage.New,
			api.New,
		),

		fx.Invoke(
			ServerStart,
		),

		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)

	return App{engine: engine}
}
