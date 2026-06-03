package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type AboutController struct {
	aboutService *services.AboutService
}

func NewAboutController(aboutService *services.AboutService) *AboutController {
	return &AboutController{aboutService: aboutService}
}

// GetAll
// @Summary Get active about sections
// @Description Get list of all active about sections
// @Tags Public Content
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.AboutSection}
// @Router /api/about [get]
func (ctrl *AboutController) GetAll(c *gin.Context) {
	sections, err := ctrl.aboutService.GetAllActive()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "About sections fetched successfully", sections)
}

// AdminGetAll
// @Summary Get all about sections for administration
// @Description Get paginated list of all about sections
// @Tags Admin About Sections
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/about [get]
func (ctrl *AboutController) AdminGetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	sections, total, err := ctrl.aboutService.GetAllAdmin(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, sections, total, page, pageSize)
}

// Create
// @Summary Create a new about section
// @Description Create an about section
// @Tags Admin About Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAboutSectionRequest true "Create About Section Request"
// @Success 201 {object} utils.APIResponse{data=models.AboutSection}
// @Router /api/admin/about [post]
func (ctrl *AboutController) Create(c *gin.Context) {
	var req dto.CreateAboutSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	section, err := ctrl.aboutService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "About section created successfully", section)
}

// Update
// @Summary Update an existing about section
// @Description Update details of an about section by ID
// @Tags Admin About Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "About Section ID"
// @Param request body dto.UpdateAboutSectionRequest true "Update About Section Request"
// @Success 200 {object} utils.APIResponse{data=models.AboutSection}
// @Router /api/admin/about/{id} [put]
func (ctrl *AboutController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid section ID")
		return
	}

	var req dto.UpdateAboutSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	section, err := ctrl.aboutService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "About section updated successfully", section)
}

// Delete
// @Summary Delete an about section
// @Description Delete an about section by ID
// @Tags Admin About Sections
// @Security BearerAuth
// @Param id path int true "About Section ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/about/{id} [delete]
func (ctrl *AboutController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid section ID")
		return
	}

	if err := ctrl.aboutService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "About section deleted successfully", nil)
}
