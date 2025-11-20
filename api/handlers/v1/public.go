package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
)

// @Router		/public/pipeline [POST]
// @Summary		Create pipeline
// @Tags        Public Pipeline
// @Description	Here pipeline can be created.
// @Accept      json
// @Produce		json
// @Param       AccessKey header   string  true "Access key of a project"
// @Param       SecretKey header   string  true "Secret key of a project"
// @Param       post      body     models.CreatePipelineIntegrationReq true "post info"
// @Success		200 	  {object} models.PipelineApiResponse
// @Failure     default   {object} models.DefaultResponse
func (h *handlerV1) PublicPipelineCreate(c *gin.Context) {
	var (
		err       error
		body      = &models.CreatePipelineIntegrationReq{}
		accessKey = c.Request.Header.Get("AccessKey")
		secretKey = c.Request.Header.Get("SecretKey")
	)

	// validate if these access and secret key are exist and match.
	user, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		Username: accessKey,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().UserGet()") {
		return
	}

	if !etc.CheckPasswordHash(secretKey, user.Password) {
		c.AbortWithStatusJSON(http.StatusForbidden, models.DefaultResponse{
			ErrorCode:    ErrorCodeNotAllowed,
			ErrorMessage: "access or secret key is incorrect",
		})
		return
	}

	err = c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "invalid request body") {
		return
	}

	msg, valid := body.Validate()
	if !valid {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("not valid request"), msg)
		return
	}

	sizeKb, err := etc.GetFileSizeFromUrl(body.InputURL)
	if err != nil {
		HandleBadRequestErrWithMessage(c, h.log, err, "send valid input url")
		return
	}

	projectId, err := strconv.Atoi(string(body.ProjectID[3:]))
	if HandleBadRequestErrWithMessage(c, h.log, err, "project id is not valid") {
		return
	}

	project, err := h.storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{
		ProjectId: projectId,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().ProjectGet()") {
		return
	}

	company, err := h.storage.Postgres().CompanyGet(context.Background(), &models.CompanyGetReq{ID: project.CompanyID})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().CompanyGet()") {
		return
	}

	// check if the user has access to the project.
	if company.OwnerID != user.ID && project.OwnerID != user.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, models.DefaultResponse{
			ErrorCode:    ErrorCodeNotAllowed,
			ErrorMessage: "there is not access to this project",
		})
		return
	}

	// prepare audo tracks here
	audioTracks := []models.AudioTrack{}
	for _, e := range body.AudioTracks {
		sizeKb, err = etc.GetFileSizeFromUrl(e.InputURL)
		if err != nil {
			HandleBadRequestErrWithMessage(c, h.log, err, "send valid input url")
			return
		}
		audioTracks = append(audioTracks, models.AudioTrack{
			Id:           uuid.NewString(),
			SizeKb:       float32(sizeKb),
			InputURL:     e.InputURL,
			LanguageCode: e.LanguageCode,
			Language:     e.Language,
		})
	}

	// prepare subtitles
	subtitles := []models.Subtitle{}
	for _, e := range body.Subtitle {
		sizeKb, err = etc.GetFileSizeFromUrl(e.InputURL)
		if err != nil {
			HandleBadRequestErrWithMessage(c, h.log, err, "send valid input url")
			return
		}
		subtitles = append(subtitles, models.Subtitle{
			Id:           uuid.NewString(),
			SizeKb:       float32(sizeKb),
			InputURL:     e.InputURL,
			LanguageCode: e.LanguageCode,
			Language:     e.Language,
		})
	}

	res, err := h.storage.Postgres().PipelineCreate(context.Background(), &models.PipelineCreateReq{
		ID:                uuid.NewString(),
		ProjectID:         project.ID,
		Stage:             "initial",
		InputURL:          body.InputURL,
		OutputKey:         body.OutputKey,
		OutputPath:        body.BucketName,
		MaxResolution:     "1080p",
		SizeKB:            sizeKb,
		ResolutionsString: body.Resolutions,
		AudioTracks:       audioTracks,
		Subtitle:          subtitles,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().PipelineCreate()") {
		return
	}

	storage, err := h.storage.Postgres().StorageGet(context.Background(), &models.StorageGetReq{
		ID: project.StorageID,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().StorageGet()") {
		return
	}

	err = h.rabbit.PublishPipeline(&models.PipelineRabbitMq{
		Id:           res.ID,
		InputURI:     res.InputURL,
		OutputKey:    res.OutputKey,
		CdnUrl:       storage.DomainName,
		CdnAccessKey: storage.AccessKey,
		CdnSecretKey: storage.SecretKey,
		CdnRegion:    storage.Region,
		CdnBucket:    res.OutputPath,
		CdnType:      storage.Type,
		Resolutions:  body.Resolutions,
		AudioTracks:  audioTracks,
		Language:     body.Language,
		LanguageCode: body.LanguageCode,
		Subtitle:     subtitles,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			ErrorCode:    ErrorCodeInternal,
			ErrorMessage: "couldn't publish to message broker",
		})
		return
	}

	c.JSON(http.StatusOK, &models.PipelineApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}
