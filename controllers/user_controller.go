package controllers

import (
	"net/http"
	"strconv"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetAll
// @Summary Get all users
// @Description Get paginated list of users for administration
// @Tags Admin Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search query"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/admin/users [get]
func (ctrl *UserController) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	search := c.Query("search")

	resp, err := ctrl.userService.FindAll(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, resp.Users, resp.Total, page, pageSize)
}

// Create
// @Summary Create a new user
// @Description Create a MANAGER or SUPER_ADMIN user
// @Tags Admin Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create User Request"
// @Success 201 {object} utils.APIResponse{data=dto.UserResponse}
// @Router /api/admin/users [post]
func (ctrl *UserController) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	user, err := ctrl.userService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", user)
}

// Update
// @Summary Update an existing user
// @Description Update user details by ID
// @Tags Admin Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update User Request"
// @Success 200 {object} utils.APIResponse{data=dto.UserResponse}
// @Router /api/admin/users/{id} [put]
func (ctrl *UserController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID")
		return
	}

	currentUserIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	currentUserID := currentUserIDVal.(uint)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	user, err := ctrl.userService.Update(uint(id), req, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user)
}

// Delete
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Admin Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/users/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID")
		return
	}

	currentUserIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	currentUserID := currentUserIDVal.(uint)

	if err := ctrl.userService.Delete(uint(id), currentUserID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}
