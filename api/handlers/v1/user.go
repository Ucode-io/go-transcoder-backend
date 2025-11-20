package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	token "gitlab.com/transcodeuz/transcode-rest/api/tokens"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
)

// @Router		/user/login [POST]
// @Summary		Login user
// @Tags        User
// @Description	Here user can be logged in.
// @Accept      json
// @Produce		json
// @Param       post     body       models.UserLoginRequest true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserLogin(c *gin.Context) {
	body := &models.UserLoginRequest{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	user, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		Username: body.UserName,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "Error while getting user by user_name, ") {
		return
	}

	if !etc.CheckPasswordHash(body.Password, user.Password) {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("password_incorrect"), "password_incorrect")
		return
	}

	h.jwthandler = token.JWTHandler{
		Sub:       user.ID,
		Role:      user.Role,
		SigninKey: h.cfg.SignInKey,
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}

	user.AccessToken, user.RefreshToken, err = h.jwthandler.GenerateAuthJWT()
	if HandleInternalWithMessage(c, h.log, err, "error_while_creating_jwt_token") {
		return
	}

	_, err = h.storage.Postgres().UserUpdate(context.Background(), &models.UserUpdateReq{
		ID:           user.ID,
		RefreshToken: user.RefreshToken,
	})
	if HandleInternalWithMessage(c, h.log, err, "error_while_updating_refresh_token_of_user") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         user,
	})
}

// @Router		/user [POST]
// @Summary		Create user
// @Tags        User
// @Description	Here user can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserCreateReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserCreate(c *gin.Context) {
	body := &models.UserCreateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	if body.Role != "superadmin" {
		HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("only_superadmin_can_be_created"), "only_superadmin_can_be_created")
		return
	}
	body.Password, err = etc.HashPassword(body.Password)
	if HandleInternalWithMessage(c, h.log, err, "password_should_be_valid") {
		return
	}

	body.ID = uuid.New().String()
	res, err := h.storage.Postgres().UserCreate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "UserCreate: h.storage.Postgres().UserCreate()") {
		return
	}

	h.jwthandler = token.JWTHandler{
		Sub:       body.ID,
		Role:      body.Role,
		SigninKey: h.cfg.SignInKey,
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}

	res.AccessToken, res.RefreshToken, err = h.jwthandler.GenerateAuthJWT()
	if HandleInternalWithMessage(c, h.log, err, "error_while_creating_jwt_token") {
		return
	}

	_, err = h.storage.Postgres().UserUpdate(context.Background(), &models.UserUpdateReq{
		ID:           body.ID,
		RefreshToken: res.RefreshToken,
	})
	if HandleInternalWithMessage(c, h.log, err, "error_while_updating_refresh_token_of_user") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/user/{id} [GET]
// @Summary		Get user by key
// @Tags        User
// @Description	Here user can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "uuid"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserGet(c *gin.Context) {
	id := c.Param("id")
	claim, err := GetClaims(*h, c)

	if HandleBadRequestErrWithMessage(c, h.log, err, "UserGet:GetClaims()") {
		return
	}

	// super admin can see others profiles
	if id != "" && claim.Role == "superadmin" {
		claim.Sub = id
	}

	res, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		ID: claim.Sub,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "UserGet:h.storage.Postgres().UserGet()") {
		return
	}
	res.Password = ""
	c.JSON(http.StatusOK, models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/user/list [GET]
// @Summary		Get users list
// @Tags        User
// @Description	Here all users can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.UserFindReq true "filters"
// @Success		200 	{object}  models.UserApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserFind(c *gin.Context) {
	var (
		dbReq = &models.UserFindReq{}
		err   error
	)
	dbReq.Page, err = ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "UserFind: helper.ParsePageQueryParam(c)") {
		return
	}
	dbReq.Limit, err = ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "UserFind: helper.ParseLimitQueryParam(c)") {
		return
	}

	dbReq.Search = c.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(c.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().UserFind(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "UserFind: h.storage.Postgres().UserFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}
