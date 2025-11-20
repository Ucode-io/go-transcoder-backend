package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

// @Router		/storage/list [GET]
// @Summary		Get storages list
// @Tags        Storage
// @Description	Here all storages can be got.
// @Accept      json
// @Produce		json
// @Param       filters query models.StorageFindReq true "filters"
// @Success		200 	{object}  models.StorageApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) StorageFind(c *gin.Context) {
	var (
		dbReq = &models.StorageFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "StorageFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "StorageFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().StorageFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "StorageFind: h.storage.Postgres().StorageFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.StorageApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/storage [PUT]
// @Summary		Update storage
// @Tags        Storage
// @Description	Here storage can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.StorageUpdateReq true "post info"
// @Success		200 	{object}  models.StorageApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) StorageUpdate(c *gin.Context) {
	body := &models.StorageUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "StorageUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().StorageUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "StorageUpdate: h.storage.Postgres().StorageUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.StorageApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/storage/{id} [DELETE]
// @Summary		Delete storage
// @Tags        Storage
// @Description	Here storage can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.DefaultResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) StorageDelete(c *gin.Context) {
	id := c.Param("id")

	err := h.storage.Postgres().StorageDelete(context.Background(), &models.StorageDeleteReq{ID: id})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "StorageDelete: h.storage.Postgres().StorageDelete()") {
		return
	}

	c.JSON(http.StatusOK, models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Successfully deleted",
	})
}
