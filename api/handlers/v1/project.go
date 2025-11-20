package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
)

// @Router		/project [POST]
// @Summary		Create project
// @Tags        Project
// @Description	Here project can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.ProjectCreateReq true "post info"
// @Success		200 	{object}  models.ProjectApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectCreate(c *gin.Context) {
	claim, err := GetClaims(*h, c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "GetClaims(*h, c)") {
		return
	}

	body := &models.ProjectCreateReq{}
	err = c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	company, err := h.storage.Postgres().CompanyGet(context.Background(), &models.CompanyGetReq{
		OwnerId: claim.Sub,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectCreate: h.storage.Postgres().CompanyGet()") {
		return
	}

	if company.OwnerID != claim.Sub {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("not_allowed_to_create_project"), "ProjectCreate: company.OwnerID != claim.Sub ")
		return
	}

	accessKey, err := etc.GenerateAPIKey()
	if HandleInternalWithMessage(c, h.log, err, "ProjectCreate: etc.GenerateAPIKey()") {
		return
	}

	secretKey, err := etc.GenerateAPISecret()
	if HandleInternalWithMessage(c, h.log, err, "ProjectCreate: etc.GenerateAPISecret()") {
		return
	}
	password, err := etc.HashPassword(secretKey)
	if HandleInternalWithMessage(c, h.log, err, "HashPassword: etc.HashPassword(secretKey)") {
		return
	}

	user, err := h.storage.Postgres().UserCreate(context.Background(), &models.UserCreateReq{
		ID:       uuid.New().String(),
		Username: accessKey,
		Password: password,
		Role:     "project",
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectCreate: h.storage.Postgres().UserCreate()") {
		return
	}

	storage, err := h.storage.Postgres().StorageCreate(context.Background(), &models.StorageCreateReq{
		ID:         uuid.NewString(),
		Type:       "unknown",
		DomainName: "storage.example.com",
		AccessKey:  "Access key",
		SecretKey:  "Secret key",
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectCreate: h.storage.Postgres().StorageCreate()") {
		return
	}

	project, err := h.storage.Postgres().ProjectCreate(context.Background(), &models.ProjectCreateReq{
		ID:        uuid.NewString(),
		Title:     body.Title,
		AccessKey: accessKey,
		SecretKey: secretKey,
		CompanyID: company.ID,
		OwnerID:   user.ID,
		Status:    "active",
		StorageID: storage.ID,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectCreate: h.storage.Postgres().ProjectCreate()") {
		return
	}

	_, err = h.storage.Postgres().WebhookCreate(context.Background(), &models.WebhookCreateReq{
		ID:          uuid.NewString(),
		ProjectID:   project.ID,
		Title:       "PUT Update Pipeline Status",
		WebhookType: "update_status",
		URL:         "http://example.com/status",
		Active:      false,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectCreate: h.storage.Postgres().WebhookCreate()") {
		return
	}

	project.Owner = user
	project.Company = company
	project.Storage = storage

	c.JSON(http.StatusOK, &models.ProjectApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         project,
	})
}

// @Router		/project/{id} [GET]
// @Summary		Get project by key
// @Tags        Project
// @Description	Here project can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.ProjectApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectGet(c *gin.Context) {
	var res *models.ProjectResponse
	claims, err := GetClaims(*h, c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "error while getting claims: ") {
		return
	}

	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		res, err = h.storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{OwnerId: claims.Sub})
		if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().ProjectGet()") {
			return
		}
	} else {
		res, err = h.storage.Postgres().ProjectGet(context.Background(), &models.ProjectGetReq{
			ID: id,
		})
		if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().ProjectGet()") {
			return
		}
	}

	company, err := h.storage.Postgres().CompanyGet(context.Background(), &models.CompanyGetReq{ID: res.CompanyID})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().StorageGet()") {
		return
	}
	res.Company = company

	if res.OwnerID != claims.Sub && company.OwnerID != claims.Sub {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("no-access"), "no access")
		return
	}

	storage, err := h.storage.Postgres().StorageGet(context.Background(), &models.StorageGetReq{ID: res.StorageID})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().StorageGet()") {
		return
	}
	res.Storage = storage

	user, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{ID: res.OwnerID})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().UserGet()") {
		return
	}
	res.Owner = user

	relatedProjects := []models.RelatedProjects{
		{
			Value:    res.ID,
			Label: res.Title,
		},
	}

	if company.OwnerID == claims.Sub {
		projects, err := h.storage.Postgres().ProjectFind(context.Background(), &models.ProjectsFindReq{
			CompanyId: res.CompanyID,
			Page:      1,
			Limit:     1000,
		})
		if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().ProjectFind()") {
			return
		}

		for _, e := range projects.Projects {
			if e.ID != res.ID {
				relatedProjects = append(relatedProjects, models.RelatedProjects{Value: e.ID, Label: e.Title})
			}
		}

	}
	res.RelatedProjects = relatedProjects

	webhooks, err := h.storage.Postgres().WebhookFind(context.Background(), &models.WebhookFindReq{
		Page:      1,
		Limit:     1,
		ProjectId: res.ID,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGet: h.storage.Postgres().WebhookFind()") {
		return
	}

	if len(webhooks.Webhooks) > 0 {
		res.Webhook = webhooks.Webhooks[0]
	}

	c.JSON(http.StatusOK, models.ProjectApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/project/pid [GET]
// @Summary		Get project by key
// @Tags        Project
// @Description	Here project can be got by project ID -> 1000001.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       project_id query  int true "project_id"
// @Success		200 	{object}  models.ProjectApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectGetProjectID(c *gin.Context) {
	id := c.Query("project_id")
	pID, err := strconv.Atoi(id)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGetID: strconv.Atoi(id)") {
		return
	}
	res, err := h.storage.Postgres().ProjectGetID(context.Background(), pID)

	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectGetID: h.storage.Postgres().ProjectGetID()") {
		return
	}
	c.JSON(http.StatusOK, models.ProjectApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/project/list [GET]
// @Summary		Get projects list
// @Tags        Project
// @Description	Here all projects can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.ProjectsFindReq true "filters"
// @Success		200 	{object}  models.ProjectApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectFind(c *gin.Context) {
	var (
		dbReq = &models.ProjectsFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "ProjectFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "ProjectFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().ProjectFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectFind: h.storage.Postgres().ProjectFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.ProjectApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/project [PUT]
// @Summary		Update project
// @Tags        Project
// @Description	Here project can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.ProjectUpdateReq true "post info"
// @Success		200 	{object}  models.ProjectApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectUpdate(c *gin.Context) {
	body := &models.ProjectUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "ProjectUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().ProjectUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectUpdate: h.storage.Postgres().ProjectUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.ProjectApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/project/name [PUT]
// @Summary		Update project
// @Tags        Project
// @Description	Here project name can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.ProjectNameUpdateReq true "post info"
// @Success		200 	{object}  models.ProjectApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectNmUpdate(c *gin.Context) {
	body := &models.ProjectNameUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "ProjectUpdate: c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().ProjectUpdateName(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectUpdate: h.storage.Postgres().ProjectUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.ProjectApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/project/{id} [DELETE]
// @Summary		Delete project
// @Tags        Project
// @Description	Here project can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.DefaultResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) ProjectDelete(c *gin.Context) {
	id := c.Param("id")

	err := h.storage.Postgres().ProjectDelete(context.Background(), &models.ProjectDeleteReq{ID: id})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "ProjectDelete: h.storage.Postgres().ProjectDelete()") {
		return
	}

	c.JSON(http.StatusOK, models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Successfully deleted",
	})
}
