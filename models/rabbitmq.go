package models

type PipelineRabbitMq struct {
	Id           string       `json:"id"`
	InputURI     string       `json:"input_uri"`
	OutputKey    string       `json:"output_key"`
	OutputPath   string       `json:"output_path"`
	CdnUrl       string       `json:"cdn_url"`
	CdnAccessKey string       `json:"cdn_access_key"`
	CdnSecretKey string       `json:"cdn_secret_key"`
	CdnRegion    string       `json:"cdn_region"`
	CdnBucket    string       `json:"cdn_bucket"`
	CdnType      string       `json:"cdn_type"`
	Resolutions  []string     `json:"resolutions"`
	AudioTracks  []AudioTrack `json:"audio_tracks"`
	Subtitle     []Subtitle   `json:"subtitle"`
	Language     string       `json:"language"`
	LanguageCode string       `json:"language_code"`
}

type UpdatePipelineStatus struct {
	Id                  string       `json:"id"`
	Stage               string       `json:"stage"`
	Status              string       `json:"status"`
	PreparationDuration int          `json:"preparation_duration"` // milliseconds
	TranscodeDuration   int          `json:"transcode_duration"`   // milliseconds
	UploadDuration      int          `json:"upload_duration"`      // milliseconds
	VideoDuration       float64      `json:"video_duration"`
	Resolutions         []Resolution `json:"resolutions"`
	FailDescription     string       `json:"fail_description"`
	ErrorCode           string       `json:"error_code"`
}

type Resolution struct {
	Resolution string `json:"resolution"`
	Measure    string `json:"measure"`
	BitRate    string `json:"bitrate"`
}
