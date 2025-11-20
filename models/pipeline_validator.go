package models

import (
	"strconv"

	"gitlab.com/transcodeuz/transcode-rest/config"
	"golang.org/x/text/language"
)

func (t *PipelineCreateReq) Validate() (string, bool) {
	switch {
	case t.InputURL == "":
		return "input_url is required", false
	case t.OutputKey == "":
		return "output_key is required", false
	case t.OutputPath == "":
		return "outputpath is required", false
	case t.ProjectID == "":
		return "project_id is required", false
	}

	return "", true
}

func (t *CreatePipelineIntegrationReq) Validate() (string, bool) {
	switch {
	case t.InputURL == "":
		return "input_url is required", false
	case t.OutputKey == "":
		return "output_key is required", false
	case t.BucketName == "":
		return "bucket_name is required", false
	case t.ProjectID == "":
		return "project_id is required", false
	case len(t.ProjectID) != 10:
		return "project_id should be 10 in length", false
	case string(t.ProjectID[:3]) != "TR-":
		return "project_id should be started with TR-", false
	case t.Language == "":
		return "language is required field", false
	case !isValidLanguageCode(t.LanguageCode):
		return t.LanguageCode + " language code is not valid", false
	}

	if _, err := strconv.Atoi(string(t.ProjectID[3:])); err != nil {
		return "project_id should be end with 7 digit number", false
	}

	// Validate resolutions
	for _, e := range t.Resolutions {
		if !config.Resolutions[e] {
			return "invalid resolution " + e, false
		}
	}

	for _, e := range t.AudioTracks {
		status, success := e.Validate()
		if !success {
			return status, success
		}
	}

	for _, e := range t.Subtitle {
		status, success := e.Validate()
		if !success {
			return status, success
		}
	}

	return "", true
}

func isValidLanguageCode(code string) bool {
	_, err := language.Parse(code)
	return err == nil
}

func (t *AudioTrackRequest) Validate() (string, bool) {
	switch {
	case t.InputURL == "":
		return "input_url is required field", false
	case t.Language == "":
		return "language is required field", false
	case !isValidLanguageCode(t.LanguageCode):
		return t.LanguageCode + " language code is not valid", false
	}

	return "", true
}

func (t *SubtitleRequest) Validate() (string, bool) {
	switch {
	case t.InputURL == "":
		return "subtitle.input_url is required field", false
	case t.Language == "":
		return "subtitle.language is required field", false
	case !isValidLanguageCode(t.LanguageCode):
		return t.LanguageCode + " language code is not valid in subtitle", false
	}

	return "", true
}
