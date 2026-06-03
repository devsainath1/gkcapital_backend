package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type FAQController struct {
	faqService *services.FAQService
}

func NewFAQController(faqService *services.FAQService) *FAQController {
	return &FAQController{faqService: faqService}
}

// GetAll
// @Summary Get active FAQs
// @Description Get list of all active FAQs for public view
// @Tags Public FAQs
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.FAQ}
// @Router /api/faqs [get]
func (ctrl *FAQController) GetAll(c *gin.Context) {
	faqs, err := ctrl.faqService.GetAllActive()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "FAQs fetched successfully", faqs)
}

// AdminGetAll
// @Summary Get all FAQs for administration
// @Description Get paginated list of all FAQs including inactive ones
// @Tags Admin FAQs
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/faqs [get]
func (ctrl *FAQController) AdminGetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	faqs, total, err := ctrl.faqService.GetAllAdmin(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, faqs, total, page, pageSize)
}

// Create
// @Summary Create a new FAQ
// @Description Create a FAQ for administration
// @Tags Admin FAQs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateFAQRequest true "Create FAQ Request"
// @Success 201 {object} utils.APIResponse{data=models.FAQ}
// @Router /api/admin/faqs [post]
func (ctrl *FAQController) Create(c *gin.Context) {
	var req dto.CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	faq, err := ctrl.faqService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "FAQ created successfully", faq)
}

// Update
// @Summary Update an existing FAQ
// @Description Update details of a FAQ by ID
// @Tags Admin FAQs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "FAQ ID"
// @Param request body dto.UpdateFAQRequest true "Update FAQ Request"
// @Success 200 {object} utils.APIResponse{data=models.FAQ}
// @Router /api/admin/faqs/{id} [put]
func (ctrl *FAQController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid FAQ ID")
		return
	}

	var req dto.UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	faq, err := ctrl.faqService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "FAQ updated successfully", faq)
}

// Delete
// @Summary Delete a FAQ
// @Description Delete a FAQ by ID
// @Tags Admin FAQs
// @Security BearerAuth
// @Param id path int true "FAQ ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/faqs/{id} [delete]
func (ctrl *FAQController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid FAQ ID")
		return
	}

	if err := ctrl.faqService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "FAQ deleted successfully", nil)
}
