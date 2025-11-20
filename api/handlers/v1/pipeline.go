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

// @Router		/pipeline [POST]
// @Summary		Create pipeline
// @Tags        Pipeline
// @Description	Here pipeline can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.PipelineCreateReq true "post info"
// @Success		200 	{object}  models.PipelineApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineCreate(c *gin.Context) {
	body := &models.PipelineCreateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "invalid request body") {
		return
	}

	msg, valid := body.Validate()
	if !valid {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("not valid request"), msg)
		return
	}

	body.SizeKB, err = etc.GetFileSizeFromUrl(body.InputURL)
	if err != nil {
		h.log.Error("error while getting file size by its url.", err)
		HandleBadRequestErrWithMessage(c, h.log, err, "send valid input url")
		return
	}

	body.ID = uuid.NewString()
	res, err := h.storage.Postgres().PipelineCreate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().PipelineCreate()") {
		return
	}

	project, err := h.storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{
		ID: res.ProjectID,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineCreate: h.storage.Postgres().ProjectGet()") {
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

// @Router		/pipeline/{id} [GET]
// @Summary		Get pipeline by key
// @Tags        Pipeline
// @Description	Here pipeline can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.PipelineApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineGet(c *gin.Context) {
	id := c.Param("id")

	res, err := h.storage.Postgres().PipelineGet(context.Background(), &models.PipelineGetReq{
		ID: id,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineGet: h.storage.Postgres().PipelineGet()") {
		return
	}

	c.JSON(http.StatusOK, models.PipelineApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/pipeline/list [GET]
// @Summary		Get pipelines list
// @Tags        Pipeline
// @Description	Here all pipelines can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.PipelineFindReq true "filters"
// @Success		200 	{object}  models.PipelineApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineFind(c *gin.Context) {
	var (
		dbReq = &models.PipelineFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "PipelineFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "PipelineFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)
	dbReq.ProjectId = c.Query("project_id")

	res, err := h.storage.Postgres().PipelinesFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineFind: h.storage.Postgres().PipelineFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.PipelineApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/pipeline [PUT]
// @Summary		Update pipeline
// @Tags        Pipeline
// @Description	Here pipeline can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.PipelineUpdateReq true "post info"
// @Success		200 	{object}  models.PipelineApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineUpdate(c *gin.Context) {
	body := &models.PipelineUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "PipelineUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().PipelineUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineUpdate: h.storage.Postgres().PipelineUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.PipelineApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/pipeline/{id} [DELETE]
// @Summary		Delete pipeline
// @Tags        Pipeline
// @Description	Here pipeline can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.DefaultResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineDelete(c *gin.Context) {
	id := c.Param("id")

	err := h.storage.Postgres().PipelineDelete(context.Background(), &models.PipelineDeleteReq{ID: id})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineDelete: h.storage.Postgres().PipelineDelete()") {
		return
	}

	c.JSON(http.StatusOK, models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Successfully deleted",
	})
}
