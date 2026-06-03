package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type TestimonialController struct {
	testimonialService *services.TestimonialService
}

func NewTestimonialController(testimonialService *services.TestimonialService) *TestimonialController {
	return &TestimonialController{testimonialService: testimonialService}
}

// GetAll
// @Summary Get active testimonials
// @Description Get list of all active testimonials for public view
// @Tags Public Testimonials
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.Testimonial}
// @Router /api/testimonials [get]
func (ctrl *TestimonialController) GetAll(c *gin.Context) {
	testimonials, err := ctrl.testimonialService.GetAllActive()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Testimonials fetched successfully", testimonials)
}

// AdminGetAll
// @Summary Get all testimonials for administration
// @Description Get paginated list of all testimonials including inactive ones
// @Tags Admin Testimonials
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/testimonials [get]
func (ctrl *TestimonialController) AdminGetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	testimonials, total, err := ctrl.testimonialService.GetAllAdmin(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, testimonials, total, page, pageSize)
}

// Create
// @Summary Create a new testimonial
// @Description Create a testimonial for administration
// @Tags Admin Testimonials
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTestimonialRequest true "Create Testimonial Request"
// @Success 201 {object} utils.APIResponse{data=models.Testimonial}
// @Router /api/admin/testimonials [post]
func (ctrl *TestimonialController) Create(c *gin.Context) {
	var req dto.CreateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	testimonial, err := ctrl.testimonialService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Testimonial created successfully", testimonial)
}

// Update
// @Summary Update an existing testimonial
// @Description Update details of a testimonial by ID
// @Tags Admin Testimonials
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Testimonial ID"
// @Param request body dto.UpdateTestimonialRequest true "Update Testimonial Request"
// @Success 200 {object} utils.APIResponse{data=models.Testimonial}
// @Router /api/admin/testimonials/{id} [put]
func (ctrl *TestimonialController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid testimonial ID")
		return
	}

	var req dto.UpdateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	testimonial, err := ctrl.testimonialService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Testimonial updated successfully", testimonial)
}

// Delete
// @Summary Delete a testimonial
// @Description Delete a testimonial by ID
// @Tags Admin Testimonials
// @Security BearerAuth
// @Param id path int true "Testimonial ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/testimonials/{id} [delete]
func (ctrl *TestimonialController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid testimonial ID")
		return
	}

	if err := ctrl.testimonialService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Testimonial deleted successfully", nil)
}
