package postgres

import (
	"context"

	"gitlab.com/transcodeuz/transcode-rest/models"
)

type PostgresI interface {
	// common
	UpdateSingleField(ctx context.Context, req *models.UpdateSingleFieldReq) error
	CheckIfExists(ctx context.Context, req *models.CheckIfExistsReq) (*models.CheckIfExistsRes, error)

	// User
	UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error)
	UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error)
	UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error)
	UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error)
	UserDelete(ctx context.Context, req *models.UserDeleteReq) error

	// Company
	CompanyCreate(ctx context.Context, req *models.CompanyCreateReq) (*models.CompanyResponse, error)
	CompanyGet(ctx context.Context, req *models.CompanyGetReq) (*models.CompanyResponse, error)
	CompanyFind(ctx context.Context, req *models.CompaniesFindReq) (*models.CompaniesFindResponse, error)
	CompanyUpdate(ctx context.Context, req *models.CompanyUpdateReq) (*models.CompanyResponse, error)
	CompanyDelete(ctx context.Context, req *models.CompanyDeleteReq) error

	// Pipeline
	PipelineCreate(ctx context.Context, req *models.PipelineCreateReq) (*models.PipelineResponse, error)
	PipelineGet(ctx context.Context, req *models.PipelineGetReq) (*models.PipelineResponse, error)
	PipelinesFind(ctx context.Context, req *models.PipelineFindReq) (*models.PipelinesFindResponse, error)
	PipelineUpdate(ctx context.Context, req *models.PipelineUpdateReq) (*models.PipelineResponse, error)
	PipelineDelete(ctx context.Context, req *models.PipelineDeleteReq) error
	PipelineDashboarStatistics(ctx context.Context, req *models.DashboardStatisticsRequest) (*models.DashboardStatisticsResponse, error)
	
	// Project
	ProjectCreate(ctx context.Context, req *models.ProjectCreateReq) (*models.ProjectResponse, error)
	ProjectGet(ctx context.Context, req *models.ProjectGetReq) (*models.ProjectResponse, error)
	ProjectGetID(ctx context.Context, ID int) (*models.ProjectResponse, error)
	ProjectFind(ctx context.Context, req *models.ProjectsFindReq) (*models.ProjectsFindResponse, error)
	ProjectUpdate(ctx context.Context, req *models.ProjectUpdateReq) (*models.ProjectResponse, error)
	ProjectUpdateName(ctx context.Context, req *models.ProjectNameUpdateReq) (*models.ProjectResponse, error)
	ProjectDelete(ctx context.Context, req *models.ProjectDeleteReq) error

	// Storage
	StorageCreate(ctx context.Context, req *models.StorageCreateReq) (*models.StorageResponse, error)
	StorageGet(ctx context.Context, req *models.StorageGetReq) (*models.StorageResponse, error)
	StorageFind(ctx context.Context, req *models.StorageFindReq) (*models.StorageFindResponse, error)
	StorageUpdate(ctx context.Context, req *models.StorageUpdateReq) (*models.StorageResponse, error)
	StorageDelete(ctx context.Context, req *models.StorageDeleteReq) error

	// Webhook
	WebhookCreate(ctx context.Context, req *models.WebhookCreateReq) (*models.WebhookResponse, error)
	WebhookGet(ctx context.Context, req *models.WebhookGetReq) (*models.WebhookResponse, error)
	WebhookFind(ctx context.Context, req *models.WebhookFindReq) (*models.WebhookFindResponse, error)
	WebhookUpdate(ctx context.Context, req *models.WebhookUpdateReq) (*models.WebhookResponse, error)
	WebhookDelete(ctx context.Context, req *models.WebhookDeleteReq) error

	// Don't delete this line, it is used to modify the file automatically
}
