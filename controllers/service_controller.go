package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	serviceService *services.ServiceService
}

func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{serviceService: serviceService}
}

// GetAll
// @Summary Get active services
// @Description Get list of all active services for public view
// @Tags Public Services
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.Service}
// @Failure 500 {object} utils.APIResponse
// @Router /api/services [get]
func (ctrl *ServiceController) GetAll(c *gin.Context) {
	services, err := ctrl.serviceService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Services fetched successfully", services)
}

// GetByID
// @Summary Get service by ID
// @Description Get details of a single service by ID
// @Tags Public Services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} utils.APIResponse{data=models.Service}
// @Failure 404 {object} utils.APIResponse
// @Router /api/services/{id} [get]
func (ctrl *ServiceController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID")
		return
	}

	service, err := ctrl.serviceService.GetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Service not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service fetched successfully", service)
}

// AdminGetAll
// @Summary Get all services for administration
// @Description Get paginated list of all services including inactive ones
// @Tags Admin Services
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/admin/services [get]
func (ctrl *ServiceController) AdminGetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	services, total, err := ctrl.serviceService.GetAllAdmin(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, services, total, page, pageSize)
}

// Create
// @Summary Create a new service
// @Description Create a service for administration
// @Tags Admin Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateServiceRequest true "Create Service Request"
// @Success 201 {object} utils.APIResponse{data=models.Service}
// @Failure 400 {object} utils.APIResponse
// @Router /api/admin/services [post]
func (ctrl *ServiceController) Create(c *gin.Context) {
	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	service, err := ctrl.serviceService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Service created successfully", service)
}

// Update
// @Summary Update an existing service
// @Description Update details of a service by ID
// @Tags Admin Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param request body dto.UpdateServiceRequest true "Update Service Request"
// @Success 200 {object} utils.APIResponse{data=models.Service}
// @Failure 400 {object} utils.APIResponse
// @Router /api/admin/services/{id} [put]
func (ctrl *ServiceController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID")
		return
	}

	var req dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	service, err := ctrl.serviceService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service updated successfully", service)
}

// Delete
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags Admin Services
// @Security BearerAuth
// @Param id path int true "Service ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/services/{id} [delete]
func (ctrl *ServiceController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID")
		return
	}

	if err := ctrl.serviceService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service deleted successfully", nil)
}
