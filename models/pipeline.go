package models

import "time"

type PipelineCreateReq struct {
	ID                    string       `json:"id"`
	ProjectID             string       `json:"project_id"`
	Stage                 string       `json:"stage"`
	FailDescription       string       `json:"fail_description"`
	InputURL              string       `json:"input_url"`
	OutputKey             string       `json:"output_key"`
	OutputPath            string       `json:"output_path"`
	SizeKB                float64      `json:"size_kb"`
	TranscodeDurationSecs float64      `json:"transcode_duration_secs"`
	UploadDurationSecs    float64      `json:"upload_duration_secs"`
	MaxResolution         string       `json:"max_resolution"`
	ResolutionsString     []string     `json:"resolutions_string"`
	AudioTracks           []AudioTrack `json:"audio_tracks"`
	Subtitle              []Subtitle   `json:"subtitle"`
	LanguageCode          string       `json:"language_code"`
	Language              string       `json:"language"`
}

type PipelineUpdateReq struct {
	ID                         string       `json:"id"`
	Stage                      string       `json:"stage"`
	StageStatus                string       `json:"stage_status"`
	FailDescription            string       `json:"fail_description"`
	InputURL                   string       `json:"input_url"`
	OutputKey                  string       `json:"output_key"`
	OutputPath                 string       `json:"output_path"`
	SizeKB                     float64      `json:"size_kb"`
	VideoDuration              float32      `json:"video_duration"`
	PreparationDurationSeconds float64      `json:"preparation_duration_seconds"`
	TranscodeDurationSecs      float64      `json:"transcode_duration_secs"`
	UploadDurationSecs         float64      `json:"upload_duration_secs"`
	MaxResolution              string       `json:"max_resolution"`
	Resolutions                []Resolution `json:"resolutions"`
	ResolutionsString          []string     `json:"resolutions_string"`
	WebhookStatus              bool         `json:"webhook_status"`
	WebhookRetryCount          int          `json:"webhook_retry_count"`
	WebhookLastRetry           time.Time    `json:"webhook_last_retry"`
}

type PipelineGetReq struct {
	ID string `json:"id"`
}

type PipelineFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
	ProjectId        string `json:"project_id"`
	WebhookStatus    int    `json:"webhook_status"` // -1 0 1
}

type PipelineDeleteReq struct {
	ID string `json:"id"`
}

type PipelinesFindResponse struct {
	Pipelines []*PipelineResponse `json:"pipelines"`
	Count     int                 `json:"count"`
}

type PipelineResponse struct {
	ID                         string       `json:"id"`
	ProjectID                  string       `json:"project_id"`
	Stage                      string       `json:"stage"`
	StageStatus                string       `json:"stage_status"`
	FailDescription            string       `json:"fail_description"`
	InputURL                   string       `json:"input_url"`
	OutputKey                  string       `json:"output_key"`
	OutputPath                 string       `json:"output_path"`
	SizeKB                     float64      `json:"size_kb"`
	TranscodeDurationSecs      float64      `json:"transcode_duration_secs"`
	UploadDurationSecs         float64      `json:"upload_duration_secs"`
	MaxResolution              string       `json:"max_resolution"`
	Resolutions                []Resolution `json:"resolutions"`
	ResolutionsString          []string     `json:"resolutions_string"`
	VideoDuration              float32      `json:"video_duration"`
	PreparationDurationSeconds float64      `json:"preparation_duration_seconds"`
	WebhookStatus              bool         `json:"webhook_status"`
	WebhookRetryCount          int          `json:"webhook_retry_count"`
	WebhookLastRetry           time.Time    `json:"webhook_last_retry"`
	AudioTracks                []AudioTrack `json:"audio_tracks"`
	Subtitle                   []Subtitle   `json:"subtitle"`
	LanguageCode               string       `json:"language_code"`
	Language                   string       `json:"language"`
	CreatedAt                  string       `json:"created_at"`
	UpdatedAt                  string       `json:"updated_at"`
}

type PipelineApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *PipelineResponse
}

type PipelineApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *PipelinesFindResponse
}

type CreatePipelineIntegrationReq struct {
	ProjectID    string              `json:"project_id"`
	InputURL     string              `json:"input_url"`
	OutputKey    string              `json:"output_key"`
	BucketName   string              `json:"bucket_name"`
	Resolutions  []string            `json:"resolutions"`
	LanguageCode string              `json:"language_code"`
	Language     string              `json:"language"`
	AudioTracks  []AudioTrackRequest `json:"audio_tracks"`
	Subtitle     []SubtitleRequest   `json:"subtitle"`
}

type AudioTrackRequest struct {
	InputURL     string `json:"input_url"`
	LanguageCode string `json:"language_code"`
	Language     string `json:"language"`
}

type AudioTrack struct {
	Id           string  `json:"id"`
	PipelineId   string  `json:"pipeline_id"`
	SizeKb       float32 `json:"size_kb"`
	InputURL     string  `json:"input_url"`
	LanguageCode string  `json:"lang_code"`
	Language     string  `json:"language"`
	CreatedAt    string  `json:"created_at"`
}

type SubtitleRequest struct {
	InputURL     string `json:"input_url"`
	LanguageCode string `json:"language_code"`
	Language     string `json:"language"`
}

type Subtitle struct {
	Id           string  `json:"id"`
	PipelineId   string  `json:"pipeline_id"`
	SizeKb       float32 `json:"size_kb"`
	InputURL     string  `json:"input_url"`
	LanguageCode string  `json:"language_code"`
	Language     string  `json:"language"`
	CreatedAt    string  `json:"created_at"`
}

type DashboardStatisticsRequest struct {
	LastNdays uint32 `json:"last_n_days"`
	ProjectId string `json:"project_id"`
	CompanyId string `json:"company_id"`
}

type DashboardStatisticsApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *DashboardStatisticsResponse
}

type DashboardStatisticsResponse struct {
	SuccessPercent uint32           `json:"success_percent"`
	TotalCount     uint32           `json:"total_count"`
	TotalSize      float64          `json:"total_size"`
	PipelinePerDay []PipelinePerDay `json:"pipelines_by_day"`
}

type PipelinePerDay struct {
	Date  string  `json:"date"`
	Size  float64 `json:"size"`
	Count uint32  `json:"count"`
}
