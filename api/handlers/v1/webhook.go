package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

// @Router		/webhook [POST]
// @Summary		Create webhook
// @Tags        Webhook
// @Description	Here webhook can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.WebhookCreateReq true "post info"
// @Success		200 	{object}  models.WebhookApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) WebhookCreate(c *gin.Context) {
	body := &models.WebhookCreateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().WebhookCreate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "WebhookCreate: h.storage.Postgres().WebhookCreate()") {
		return
	}

	c.JSON(http.StatusOK, &models.WebhookApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/webhook/{id} [GET]
// @Summary		Get webhook by key
// @Tags        Webhook
// @Description	Here webhook can be got.
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.WebhookApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) WebhookGet(c *gin.Context) {
	id := c.Param("id")

	res, err := h.storage.Postgres().WebhookGet(context.Background(), &models.WebhookGetReq{
		ID: id,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "WebhookGet: h.storage.Postgres().WebhookGet()") {
		return
	}

	c.JSON(http.StatusOK, models.WebhookApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/webhook/list [GET]
// @Summary		Get webhooks list
// @Tags        Webhook
// @Description	Here all webhooks can be got.
// @Accept      json
// @Produce		json
// @Param       filters query models.WebhookFindReq true "filters"
// @Success		200 	{object}  models.WebhookApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) WebhookFind(c *gin.Context) {
	var (
		dbReq = &models.WebhookFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "WebhookFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "WebhookFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().WebhookFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "WebhookFind: h.storage.Postgres().WebhookFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.WebhookApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/webhook [PUT]
// @Summary		Update webhook
// @Tags        Webhook
// @Description	Here webhook can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.WebhookUpdateReq true "post info"
// @Success		200 	{object}  models.WebhookApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) WebhookUpdate(c *gin.Context) {
	body := &models.WebhookUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "WebhookUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().WebhookUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "WebhookUpdate: h.storage.Postgres().WebhookUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.WebhookApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/webhook/{id} [DELETE]
// @Summary		Delete webhook
// @Tags        Webhook
// @Description	Here webhook can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.DefaultResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) WebhookDelete(c *gin.Context) {
	id := c.Param("id")

	err := h.storage.Postgres().WebhookDelete(context.Background(), &models.WebhookDeleteReq{ID: id})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "WebhookDelete: h.storage.Postgres().WebhookDelete()") {
		return
	}

	c.JSON(http.StatusOK, models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Successfully deleted",
	})
}
