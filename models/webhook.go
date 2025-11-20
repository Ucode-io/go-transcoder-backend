package models

type WebhookCreateReq struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	WebhookType string `json:"webhook_type"`
	URL         string `json:"url"`
	Active      bool   `json:"active"`
}

type WebhookUpdateReq struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	WebhookType string `json:"webhook_type"`
	URL         string `json:"url"`
	Active      bool   `json:"active"`
}

type WebhookGetReq struct {
	ID string `json:"id"`
}

type WebhookFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
	ProjectId        string `json:"project_id"`
}

type WebhookDeleteReq struct {
	ID string `json:"id"`
}

type WebhookFindResponse struct {
	Webhooks []*WebhookResponse `json:"webhooks"`
	Count    int                `json:"count"`
}

type WebhookResponse struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	WebhookType string `json:"webhook_type"`
	URL         string `json:"url"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type WebhookApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *WebhookResponse
}

type WebhookApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *WebhookFindResponse
}

type PipelineWebhookResponse struct {
	Success bool `json:"success"`
}

type PublicWebhookRequest struct {
	Stage                      string   `json:"stage"`
	StageStatus                string   `json:"stage_status"`
	FailDescription            string   `json:"fail_description"`
	InputURL                   string   `json:"input_url"`
	OutputKey                  string   `json:"output_key"`
	OutputPath                 string   `json:"output_path"`
	SizeKB                     float64  `json:"size_kb"`
	TranscodeDurationSecs      float64  `json:"transcode_duration_secs"`
	UploadDurationSecs         float64  `json:"upload_duration_secs"`
	Resolutions                []string `json:"resolutions"`
	VideoDuration              float32  `json:"video_duration"`
	PreparationDurationSeconds float64  `json:"preparation_duration_seconds"`
}