package models

type ProjectResponse struct {
	ID              string            `json:"id"`
	ProjectID       string            `json:"project_id"`
	Title           string            `json:"title"`
	AccessKey       string            `json:"access_key"`
	SecretKey       string            `json:"secret_key"`
	CompanyID       string            `json:"company_id"`
	OwnerID         string            `json:"owner_id"`
	Status          string            `json:"status"`
	StorageID       string            `json:"storage_id"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	DeletedAt       string            `json:"deleted_at"`
	Owner           *UserResponse     `json:"owner"`
	Storage         *StorageResponse  `json:"storage"`
	Company         *CompanyResponse  `json:"company"`
	RelatedProjects []RelatedProjects `json:"related_projects"`
	Webhook         *WebhookResponse  `json:"webhook"`
}

type RelatedProjects struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ProjectCreateReq struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	CompanyID string `json:"company_id"`
	OwnerID   string `json:"owner_id"`
	Status    string `json:"status"`
	StorageID string `json:"storage_id"`
}

type ProjectUpdateReq struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Status    string `json:"status"`
}

type ProjectNameUpdateReq struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ProjectGetReq struct {
	ID        string `json:"id"`
	OwnerId   string `json:"owner_id"`
	ProjectId int    `json:"project_id"`
}

type ProjectsFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
	CompanyId        string `json:"company_id"`
}

type ProjectDeleteReq struct {
	ID string `json:"id"`
}

type ProjectsFindResponse struct {
	Projects []*ProjectResponse `json:"projects"`
	Count    int                `json:"count"`
}

type ProjectApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *ProjectResponse
}

type ProjectApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *ProjectsFindResponse
}
