package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type HomepageController struct {
	homepageService *services.HomepageService
}

func NewHomepageController(homepageService *services.HomepageService) *HomepageController {
	return &HomepageController{homepageService: homepageService}
}

// GetAll
// @Summary Get active homepage sections
// @Description Get list of all active homepage sections
// @Tags Public Content
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.HomepageSection}
// @Router /api/homepage [get]
func (ctrl *HomepageController) GetAll(c *gin.Context) {
	sections, err := ctrl.homepageService.GetAllActive()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Homepage sections fetched successfully", sections)
}

// AdminGetAll
// @Summary Get all homepage sections for administration
// @Description Get paginated list of all homepage sections
// @Tags Admin Homepage Sections
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/homepage [get]
func (ctrl *HomepageController) AdminGetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	sections, total, err := ctrl.homepageService.GetAllAdmin(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, sections, total, page, pageSize)
}

// Create
// @Summary Create a new homepage section
// @Description Create a homepage section
// @Tags Admin Homepage Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateHomepageSectionRequest true "Create Homepage Section Request"
// @Success 201 {object} utils.APIResponse{data=models.HomepageSection}
// @Router /api/admin/homepage [post]
func (ctrl *HomepageController) Create(c *gin.Context) {
	var req dto.CreateHomepageSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	section, err := ctrl.homepageService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Homepage section created successfully", section)
}

// Update
// @Summary Update an existing homepage section
// @Description Update details of a homepage section by ID
// @Tags Admin Homepage Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Homepage Section ID"
// @Param request body dto.UpdateHomepageSectionRequest true "Update Homepage Section Request"
// @Success 200 {object} utils.APIResponse{data=models.HomepageSection}
// @Router /api/admin/homepage/{id} [put]
func (ctrl *HomepageController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid section ID")
		return
	}

	var req dto.UpdateHomepageSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	section, err := ctrl.homepageService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Homepage section updated successfully", section)
}

// Delete
// @Summary Delete a homepage section
// @Description Delete a homepage section by ID
// @Tags Admin Homepage Sections
// @Security BearerAuth
// @Param id path int true "Homepage Section ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/homepage/{id} [delete]
func (ctrl *HomepageController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid section ID")
		return
	}

	if err := ctrl.homepageService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Homepage section deleted successfully", nil)
}
