package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type SEOController struct {
	seoService *services.SEOService
}

func NewSEOController(seoService *services.SEOService) *SEOController {
	return &SEOController{seoService: seoService}
}

// GetBySlug
// @Summary Get SEO configurations by page slug
// @Description Fetch meta tags, canonical URL and open graph fields for a public page
// @Tags Public Content
// @Produce json
// @Param slug query string true "Page slug (e.g. home, about, services, contact)"
// @Success 200 {object} utils.APIResponse{data=models.SEOPage}
// @Failure 404 {object} utils.APIResponse
// @Router /api/seo [get]
func (ctrl *SEOController) GetBySlug(c *gin.Context) {
	slug := c.Query("slug")
	if slug == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Page slug is required")
		return
	}

	page, err := ctrl.seoService.GetBySlug(slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "SEO page configuration not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "SEO configuration fetched successfully", page)
}

// GetAll
// @Summary Get all SEO page configurations
// @Description Get list of all SEO entries for management
// @Tags Admin SEO Pages
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.SEOPage}
// @Router /api/admin/seo [get]
func (ctrl *SEOController) GetAll(c *gin.Context) {
	pages, err := ctrl.seoService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "SEO configurations fetched successfully", pages)
}

// Create
// @Summary Create a new SEO configuration
// @Description Setup SEO meta tags for a page slug
// @Tags Admin SEO Pages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateSEOPageRequest true "Create SEO configuration request"
// @Success 201 {object} utils.APIResponse{data=models.SEOPage}
// @Router /api/admin/seo [post]
func (ctrl *SEOController) Create(c *gin.Context) {
	var req dto.CreateSEOPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctrl.seoService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "SEO configuration created successfully", page)
}

// Update
// @Summary Update an existing SEO configuration
// @Description Update SEO parameters by ID
// @Tags Admin SEO Pages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "SEO Configuration ID"
// @Param request body dto.UpdateSEOPageRequest true "Update SEO configuration request"
// @Success 200 {object} utils.APIResponse{data=models.SEOPage}
// @Router /api/admin/seo/{id} [put]
func (ctrl *SEOController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	var req dto.UpdateSEOPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctrl.seoService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "SEO configuration updated successfully", page)
}

// Delete
// @Summary Delete an SEO configuration
// @Description Delete an SEO configuration by ID
// @Tags Admin SEO Pages
// @Security BearerAuth
// @Param id path int true "SEO Configuration ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/seo/{id} [delete]
func (ctrl *SEOController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	if err := ctrl.seoService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "SEO configuration deleted successfully", nil)
}
