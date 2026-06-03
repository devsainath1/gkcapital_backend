package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	contactService *services.ContactService
}

func NewContactController(contactService *services.ContactService) *ContactController {
	return &ContactController{contactService: contactService}
}

// Submit
// @Summary Submit a contact inquiry
// @Description Submit the public contact form
// @Tags Contact Inquiry
// @Accept json
// @Produce json
// @Param request body dto.ContactRequest true "Contact inquiry submission data"
// @Success 201 {object} utils.APIResponse{data=models.ContactInquiry}
// @Failure 400 {object} utils.APIResponse
// @Router /api/contact [post]
func (ctrl *ContactController) Submit(c *gin.Context) {
	var req dto.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inquiry, err := ctrl.contactService.Submit(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Inquiry submitted successfully", inquiry)
}

// GetAll
// @Summary Get contact inquiries
// @Description Get list of contact inquiries with search and status filters
// @Tags Admin Inquiries
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (new/read/responded)"
// @Param search query string false "Search query (name, email, phone)"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/contact-inquiries [get]
func (ctrl *ContactController) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")
	search := c.Query("search")

	inquiries, total, err := ctrl.contactService.GetAll(page, pageSize, status, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, inquiries, total, page, pageSize)
}

// UpdateStatus
// @Summary Update contact inquiry status
// @Description Update status of an inquiry by ID
// @Tags Admin Inquiries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Inquiry ID"
// @Param request body dto.UpdateInquiryStatusRequest true "New status request"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/contact-inquiries/{id} [patch]
func (ctrl *ContactController) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid inquiry ID")
		return
	}

	var req dto.UpdateInquiryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.contactService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inquiry status updated successfully", nil)
}
