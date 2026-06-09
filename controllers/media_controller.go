package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gk-capital-backend/dto"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"github.com/gin-gonic/gin"
)

type MediaController struct {
	mediaService *services.MediaService
}

func NewMediaController(mediaService *services.MediaService) *MediaController {
	return &MediaController{mediaService: mediaService}
}

// Serve serves a media asset by ID as raw binary with proper headers
// @Summary Serve media asset by ID
// @Description Serve image binary data with correct Content-Type and caching headers
// @Tags Public Media
// @Produce octet-stream
// @Param id path int true "Media Asset ID"
// @Success 200
// @Failure 404 {object} utils.APIResponse
// @Router /api/media/{id} [get]
func (ctrl *MediaController) Serve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	asset, err := ctrl.mediaService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// Set caching headers — cache for 30 days
	c.Header("Cache-Control", "public, max-age=2592000, immutable")
	c.Header("ETag", fmt.Sprintf(`"%d-%d"`, asset.ID, asset.UpdatedAt.Unix()))
	c.Data(http.StatusOK, asset.MimeType, asset.Data)
}

// ServeByName serves a media asset by name
// @Summary Serve media asset by name
// @Description Serve image binary data by unique name (e.g., "logo")
// @Tags Public Media
// @Produce octet-stream
// @Param name path string true "Media Asset Name"
// @Success 200
// @Failure 404 {object} utils.APIResponse
// @Router /api/media/name/{name} [get]
func (ctrl *MediaController) ServeByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	asset, err := ctrl.mediaService.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	etag := fmt.Sprintf(`"%d-%d"`, asset.ID, asset.UpdatedAt.Unix())
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", etag)

	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Data(http.StatusOK, asset.MimeType, asset.Data)
}

// Upload handles multipart file upload and stores in DB
// @Summary Upload a media asset
// @Description Upload an image file to store as binary in the database
// @Tags Admin Media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "Image file"
// @Param name formData string true "Unique asset name"
// @Success 201 {object} utils.APIResponse{data=dto.MediaAssetResponse}
// @Failure 400 {object} utils.APIResponse
// @Router /api/admin/media/upload [post]
func (ctrl *MediaController) Upload(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Name is required")
		return
	}

	// Sanitize name
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	asset, err := ctrl.mediaService.Upload(name, mimeType, data)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.MediaAssetResponse{
		ID:       asset.ID,
		Name:     asset.Name,
		MimeType: asset.MimeType,
		Size:     asset.Size,
		URL:      fmt.Sprintf("/api/media/%d", asset.ID),
	}

	utils.SuccessResponse(c, http.StatusCreated, "Media uploaded successfully", resp)
}

// AdminGetAll lists all media assets (without binary data)
// @Summary List all media assets
// @Description Get list of all media assets for administration
// @Tags Admin Media
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/media [get]
func (ctrl *MediaController) AdminGetAll(c *gin.Context) {
	assets, err := ctrl.mediaService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var resp []dto.MediaAssetResponse
	for _, a := range assets {
		resp = append(resp, dto.MediaAssetResponse{
			ID:       a.ID,
			Name:     a.Name,
			MimeType: a.MimeType,
			Size:     a.Size,
			URL:      fmt.Sprintf("/api/media/%d", a.ID),
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Media assets fetched successfully", resp)
}

// Delete removes a media asset
// @Summary Delete a media asset
// @Description Delete a media asset by ID
// @Tags Admin Media
// @Security BearerAuth
// @Param id path int true "Media Asset ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/admin/media/{id} [delete]
func (ctrl *MediaController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid media ID")
		return
	}

	if err := ctrl.mediaService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Media deleted successfully", nil)
}
