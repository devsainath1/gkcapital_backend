package controllers

import (
	"net/http"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type WebsiteSettingController struct {
	settingService *services.WebsiteSettingService
}

func NewWebsiteSettingController(settingService *services.WebsiteSettingService) *WebsiteSettingController {
	return &WebsiteSettingController{settingService: settingService}
}

// GetAll
// @Summary Get all website settings
// @Description Get a list of all settings (admin view)
// @Tags Admin Settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]dto.WebsiteSettingResponse}
// @Router /api/admin/settings [get]
func (ctrl *WebsiteSettingController) GetAll(c *gin.Context) {
	settings, err := ctrl.settingService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Settings retrieved successfully", settings)
}

// BulkUpsert
// @Summary Update or insert multiple settings
// @Description Update settings in bulk
// @Tags Admin Settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateWebsiteSettingsRequest true "Bulk Upsert Settings Request"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/settings [put]
func (ctrl *WebsiteSettingController) BulkUpsert(c *gin.Context) {
	var req dto.UpdateWebsiteSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.settingService.BulkUpsert(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settings updated successfully", nil)
}

// GetPublic
// @Summary Get public settings as a map
// @Description Get all website settings as a key-value map (public view)
// @Tags Public Settings
// @Produce json
// @Success 200 {object} utils.APIResponse{data=map[string]string}
// @Router /api/settings [get]
func (ctrl *WebsiteSettingController) GetPublic(c *gin.Context) {
	settingsMap, err := ctrl.settingService.GetAllMap()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Public settings retrieved successfully", settingsMap)
}
