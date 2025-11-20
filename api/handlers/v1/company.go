package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
)

// @Router		/company [POST]
// @Summary		Create company
// @Tags        Company
// @Description	Here company can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.CompanyCreateReq true "post info"
// @Success		200 	{object}  models.CompanyApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CompanyCreate(c *gin.Context) {
	body := &models.CompanyCreateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	apiKey, err := etc.GenerateAPIKey()
	if HandleInternalWithMessage(c, h.log, err, "etc.GenerateAPIKey()") {
		return
	}

	apiSecret, err := etc.GenerateAPISecret()
	if HandleInternalWithMessage(c, h.log, err, "etc.GenerateAPISecret()") {
		return
	}

	password, err := etc.HashPassword(apiSecret)
	if HandleInternalWithMessage(c, h.log, err, "etc.HashPassword(apiSecret)") {
		return
	}

	user, err := h.storage.Postgres().UserCreate(context.Background(), &models.UserCreateReq{
		ID:       uuid.New().String(),
		Username: apiKey,
		Password: password,
		Role:     "project",
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "UserCreate: h.storage.Postgres().UserCreate()") {
		return
	}

	body.OwnerID = user.ID
	body.Status = "active"
	body.ID = uuid.NewString()
	company, err := h.storage.Postgres().CompanyCreate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyCreate: h.storage.Postgres().CompanyCreate()") {
		return
	}

	storage, err := h.storage.Postgres().StorageCreate(context.Background(), &models.StorageCreateReq{
		ID:         uuid.NewString(),
		Type:       "unknown",
		DomainName: "storage.example.com",
		AccessKey:  "Access key",
		SecretKey:  "Secret key",
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyCreate: h.storage.Postgres().StorageCreate()") {
		return
	}

	project, err := h.storage.Postgres().ProjectCreate(context.Background(), &models.ProjectCreateReq{
		ID:        uuid.NewString(),
		Title:     body.Title,
		AccessKey: apiKey,
		SecretKey: apiSecret,
		CompanyID: company.ID,
		OwnerID:   user.ID,
		Status:    "active",
		StorageID: storage.ID,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyCreate: h.storage.Postgres().ProjectCreate()") {
		return
	}

	company.Owner = user

	_, err = h.storage.Postgres().WebhookCreate(context.Background(), &models.WebhookCreateReq{
		ID:          uuid.NewString(),
		ProjectID:   project.ID,
		Title:       "PUT Update Pipeline Status",
		WebhookType: "update_status",
		URL:         "http://example.com/status",
		Active:      false,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyCreate: h.storage.Postgres().WebhookCreate()") {
		return
	}

	c.JSON(http.StatusOK, &models.CompanyApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         company,
	})
}

// @Router		/company/{id} [GET]
// @Summary		Get company by key
// @Tags        Company
// @Description	Here company can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.CompanyApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CompanyGet(c *gin.Context) {
	id := c.Param("id")

	res, err := h.storage.Postgres().CompanyGet(context.Background(), &models.CompanyGetReq{
		ID: id,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyGet: h.storage.Postgres().CompanyGet()") {
		return
	}

	c.JSON(http.StatusOK, models.CompanyApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/company/list [GET]
// @Summary		Get companies list
// @Tags        Company
// @Description	Here all companies can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.CompaniesFindReq true "filters"
// @Success		200 	{object}  models.CompanyApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CompanyFind(c *gin.Context) {
	var (
		dbReq = &models.CompaniesFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "CompanyFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "CompanyFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().CompanyFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyFind: h.storage.Postgres().CompanyFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.CompanyApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/company [PUT]
// @Summary		Update company
// @Tags        Company
// @Description	Here company can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.CompanyUpdateReq true "post info"
// @Success		200 	{object}  models.CompanyApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CompanyUpdate(c *gin.Context) {
	body := &models.CompanyUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "CompanyUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().CompanyUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "CompanyUpdate: h.storage.Postgres().CompanyUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.CompanyApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}
