package api

import (
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/golanguzb70/middleware/gin/basicauth"
	_ "gitlab.com/transcodeuz/transcode-rest/api/docs" // docs
	v1 "gitlab.com/transcodeuz/transcode-rest/api/handlers/v1"
	"gitlab.com/transcodeuz/transcode-rest/api/middleware"
	t "gitlab.com/transcodeuz/transcode-rest/api/tokens"
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/pkg/rabbitmq"
	"gitlab.com/transcodeuz/transcode-rest/storage"
	"gitlab.com/transcodeuz/transcode-rest/storage/redisrepo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf       config.Config
	Logger     *logger.Logger
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redisrepo.InMemoryStorageI
}

// New ...
// @title           Monolithic project API Endpoints
// @version         1.0
// @description     Here QA can test and frontend or mobile developers can get information of API endpoints.

// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(log *logger.Logger, cfg config.Config, strg storage.StorageI, rabbitMq *rabbitmq.RabbitMQ) *gin.Engine {
	// basicAuth := middleware.BasicAuth()
	// basicAuth.Users = append(basicAuth.Users, basicauth.User{
	// 	UserName: "admin",
	// 	Password: "admin",
	// })

	// basicAuth.RestrictedUrls = append(basicAuth.RestrictedUrls, "/v1/swagger/*")

	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, cfg.CSVFilePath)
	if err != nil {
		log.Error("casbin enforcer error", err)
	}
	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin error load policy", err)
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	jwtHandler := t.JWTHandler{
		SigninKey: cfg.SignInKey,
		Log:       log,
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	h := v1.New(&v1.HandlerV1Config{
		Logger:     log,
		Cfg:        cfg,
		Postgres:   strg,
		JWTHandler: jwtHandler,
		Rabbit:     rabbitMq,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.NewAuth(casbinEnforcer, jwtHandler, cfg))
	api := router.Group("/v1")
	api.Any("/ping", h.Ping)

	user := api.Group("/user")
	user.POST("/login", h.UserLogin)
	user.POST("", h.UserCreate)
	user.GET("/:id", h.UserGet)
	user.GET("/list", h.UserFind)

	company := api.Group("/company")
	company.POST("", h.CompanyCreate)
	company.GET("/:id", h.CompanyGet)
	company.GET("/list", h.CompanyFind)
	company.PUT("", h.CompanyUpdate)

	storage := api.Group("/storage")
	storage.GET("/list", h.StorageFind)
	storage.PUT("", h.StorageUpdate)
	storage.DELETE(":id", h.StorageDelete)

	project := api.Group("/project")
	project.POST("", h.ProjectCreate)
	project.GET("/:id", h.ProjectGet)
	project.GET("/pid", h.ProjectGetProjectID)
	project.GET("/list", h.ProjectFind)
	project.PUT("", h.ProjectUpdate)
	project.PUT("/name", h.ProjectNmUpdate)
	project.DELETE(":id", h.ProjectDelete)

	pipeline := api.Group("/pipeline")
	pipeline.POST("", h.PipelineCreate)
	pipeline.GET("/:id", h.PipelineGet)
	pipeline.GET("/list", h.PipelineFind)
	pipeline.PUT("", h.PipelineUpdate)
	pipeline.DELETE(":id", h.PipelineDelete)

	webhook := api.Group("/webhook")
	webhook.POST("", h.WebhookCreate)
	webhook.GET("/:id", h.WebhookGet)
	webhook.GET("/list", h.WebhookFind)
	webhook.PUT("", h.WebhookUpdate)
	webhook.DELETE(":id", h.WebhookDelete)

	public := api.Group("/public")
	public.POST("/pipeline", h.PublicPipelineCreate)

	dashboard := api.Group("/dashboard")
	dashboard.GET("/statistics", h.PipelineDashboarStatistics)
	dashboard.GET("/wait", h.Wait)

	media := api.Group("/media")
	api.Static("/media", "./media")
	media.POST("/photo", h.UploadMedia)

	// Don't delete this line, it is used to modify the file automatically
	url := ginSwagger.URL("swagger/doc.json")
	api.Use().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
