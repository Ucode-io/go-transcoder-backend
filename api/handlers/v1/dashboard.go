package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

// @Router		/dashboard/statistics [GET]
// @Summary		Get PipelineDashboarStatistics
// @Tags        Dashboard
// @Description	Here all dashboard statistics can be taken.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.DashboardStatisticsRequest true "filters"
// @Success		200 	{object}  models.PipelineApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) PipelineDashboarStatistics(c *gin.Context) {
	var (
		dbReq = &models.DashboardStatisticsRequest{}
		err   error
	)

	claims, err := GetClaims(*h, c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "We couldn't authorize you!") {
		return
	}
	
	lastNdays, err := strconv.Atoi(c.Query("last_n_days"))
	if err != nil {
		fmt.Println("here is error")
		lastNdays = 7
	}
	dbReq.LastNdays = uint32(lastNdays)
	dbReq.ProjectId = c.Query("project_id")
	dbReq.CompanyId = c.Query("company_id")

	if claims.Sub != "superadmin" {
		if dbReq.ProjectId == "" {
			HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("project_id is required for not super users"), "project_id is required for not super users")
			return
		}
	}

	res, err := h.storage.Postgres().PipelineDashboarStatistics(context.Background(), dbReq)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "PipelineDashboarStatistics: h.storage.Postgres().PipelineDashboarStatistics()") {
		return
	}

	c.JSON(http.StatusOK, &models.DashboardStatisticsApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}
