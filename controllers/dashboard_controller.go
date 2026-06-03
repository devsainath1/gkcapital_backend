package controllers

import (
	"net/http"

	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService *services.DashboardService
}

func NewDashboardController(dashboardService *services.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

// GetStats
// @Summary Get Admin Dashboard Statistics
// @Description Fetch aggregated metrics for services, testimonials, and inquiries
// @Tags Admin Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=services.DashboardStats}
// @Router /api/admin/dashboard [get]
func (ctrl *DashboardController) GetStats(c *gin.Context) {
	stats, err := ctrl.dashboardService.GetStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Dashboard statistics fetched successfully", stats)
}
