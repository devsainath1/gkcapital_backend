package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanInquiryController struct {
	loanInquiryService *services.LoanInquiryService
}

func NewLoanInquiryController(loanInquiryService *services.LoanInquiryService) *LoanInquiryController {
	return &LoanInquiryController{loanInquiryService: loanInquiryService}
}

// Submit
// @Summary Submit a loan inquiry (Apply Now)
// @Description Submit the loan application form
// @Tags Loan Inquiry
// @Accept json
// @Produce json
// @Param request body dto.LoanInquiryRequest true "Loan inquiry application data"
// @Success 201 {object} utils.APIResponse{data=models.LoanInquiry}
// @Failure 400 {object} utils.APIResponse
// @Router /api/loan-inquiry [post]
func (ctrl *LoanInquiryController) Submit(c *gin.Context) {
	var req dto.LoanInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inquiry, err := ctrl.loanInquiryService.Submit(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Loan application submitted successfully", inquiry)
}

// GetAll
// @Summary Get loan inquiries
// @Description Get list of loan applications with search and status filters
// @Tags Admin Inquiries
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (new/contacted/approved/rejected)"
// @Param search query string false "Search query (name, email, phone, city)"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/loan-inquiries [get]
func (ctrl *LoanInquiryController) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")
	search := c.Query("search")

	inquiries, total, err := ctrl.loanInquiryService.GetAll(page, pageSize, status, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, inquiries, total, page, pageSize)
}

// UpdateStatus
// @Summary Update loan inquiry status
// @Description Update status of an application by ID
// @Tags Admin Inquiries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Inquiry ID"
// @Param request body dto.UpdateInquiryStatusRequest true "New status request"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/loan-inquiries/{id} [patch]
func (ctrl *LoanInquiryController) UpdateStatus(c *gin.Context) {
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

	if err := ctrl.loanInquiryService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Loan inquiry status updated successfully", nil)
}
