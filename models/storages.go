package models

type StorageCreateReq struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	DomainName string `json:"domain_name"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Region     string `json:"region"`
}

type StorageUpdateReq struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	DomainName string `json:"domain_name"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Region     string `json:"region"`
}

type StorageGetReq struct {
	ID string `json:"id"`
}

type StorageFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
}

type StorageDeleteReq struct {
	ID string `json:"id"`
}

type StorageFindResponse struct {
	Storages []*StorageResponse `json:"storages"`
	Count    int                `json:"count"`
}

type StorageResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	DomainName string `json:"domain_name"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Region     string `json:"region"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}

type StorageApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *StorageResponse
}

type StorageApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *StorageFindResponse
}
