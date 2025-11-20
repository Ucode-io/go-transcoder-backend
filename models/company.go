package models

type CompanyResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	OwnerID   string        `json:"owner_id"`
	Status    string        `json:"status"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Owner     *UserResponse `json:"owner"`
}

type CompanyCreateReq struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	OwnerID string `json:"owner_id"`
	Status  string `json:"status"`
}

type CompanyUpdateReq struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type CompanyGetReq struct {
	ID    string `json:"id"`
	OwnerId string `json:"owner_id"`
}

type CompaniesFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
}

type CompanyDeleteReq struct {
	ID string `json:"id"`
}

type CompaniesFindResponse struct {
	Companies []*CompanyResponse `json:"companies"`
	Count     int                `json:"count"`
}

type CompanyApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *CompanyResponse
}

type CompanyApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *CompaniesFindResponse
}

type Company struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	OwnerID   string `json:"owner_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
