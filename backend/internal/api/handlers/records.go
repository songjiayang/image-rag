package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"image-rag-backend/internal/logger"
	"image-rag-backend/internal/models"
	"image-rag-backend/internal/services"
)

// @Summary Create a new record with images
// @Description Create a new record with associated images and generate vectors for similarity search
// @Tags Records
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Record name"
// @Param description formData string false "Record description"
// @Param images formData []file true "Image files to upload"
// @Success 201 {object} models.RecordResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /records [post]

type RecordHandler struct {
	recordService *services.RecordService
	vectorService *services.VectorService
	logger        *logger.Logger
}

func NewRecordHandler(recordService *services.RecordService, vectorService *services.VectorService, logger *logger.Logger) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
		vectorService: vectorService,
		logger:        logger,
	}
}

// CreateRecord creates a new record with images
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// Create record first
	record, err := h.recordService.CreateRecord(name, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process uploaded images
	form, _ := c.MultipartForm()
	files := form.File["images"]

	var uploadedImages []*models.Image

	for _, file := range files {
		// Validate file
		if err := services.ValidateImageFile(file.Filename); err != nil {
			// Skip invalid files but don't fail the entire operation
			continue
		}

		// Generate unique filename
		filename := services.GenerateUniqueFilename(file.Filename)
		filePath := filepath.Join("uploads", filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			// Log error but continue with other files
			fmt.Printf("Failed to save file %s: %v\n", filename, err)
			continue
		}

		// Generate vector
		vectorID, err := h.vectorService.GenerateVector(filePath)
		if err != nil {
			// Clean up file if vector generation fails
			_ = services.NewRecordService().DeleteImageByPath(filePath)
			fmt.Printf("Failed to generate vector for %s: %v\n", filename, err)
			continue
		}

		// Add image to record
		image, err := h.recordService.AddImageToRecord(record.ID, filename, vectorID)
		if err != nil {
			// Clean up file and vector if adding to record fails
			_ = services.NewRecordService().DeleteImageByPath(filePath)
			_ = h.vectorService.DeleteVector(vectorID)
			fmt.Printf("Failed to add image to record: %v\n", err)
			continue
		}

		uploadedImages = append(uploadedImages, image)
	}

	// Reload record with images
	record, _ = h.recordService.GetRecord(record.ID)

	c.JSON(http.StatusCreated, record)
}

// GetRecords lists all records with pagination
func (h *RecordHandler) GetRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	records, total, err := h.recordService.GetRecords(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRecord retrieves a single record by ID
func (h *RecordHandler) GetRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record ID"})
		return
	}

	record, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateRecord updates a record
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record ID"})
		return
	}

	var req models.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.recordService.UpdateRecord(uint(id), req.Name, req.Description)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord deletes a record and its associated images
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record ID"})
		return
	}

	// Get record to delete associated images
	record, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete vectors from vector service
	for _, image := range record.Images {
		if err := h.vectorService.DeleteVector(image.VectorID); err != nil {
			// Log error but continue with deletion
			fmt.Printf("Failed to delete vector %s: %v\n", image.VectorID, err)
		}
	}

	// Delete record (will cascade to images due to foreign key constraint)
	if err := h.recordService.DeleteRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "record deleted successfully"})
}

// AddImageToRecord adds an image to an existing record
func (h *RecordHandler) AddImageToRecord(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record ID"})
		return
	}

	// Ensure record exists
	_, err = h.recordService.GetRecord(uint(recordID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	defer file.Close()

	// Validate file
	if err := services.ValidateImageFile(header.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique filename
	filename := services.GenerateUniqueFilename(header.Filename)
	filePath := filepath.Join("uploads", filename)

	// Save file
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Generate vector
	vectorID, err := h.vectorService.GenerateVector(filePath)
	if err != nil {
		h.logger.Error("generarte image vector with error: %v", err)
		// Clean up file
		_ = services.NewRecordService().DeleteImageByPath(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate vector"})
		return
	}

	// Add image to record
	image, err := h.recordService.AddImageToRecord(uint(recordID), filename, vectorID)
	if err != nil {
		// Clean up file and vector
		_ = services.NewRecordService().DeleteImageByPath(filePath)
		_ = h.vectorService.DeleteVector(vectorID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add image to record"})
		return
	}

	c.JSON(http.StatusCreated, image)
}

// DeleteImage deletes an image from a record
func (h *RecordHandler) DeleteImage(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// Get image to delete vector
	image, err := h.recordService.GetImage(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	// Delete vector from vector service
	if err := h.vectorService.DeleteVector(image.VectorID); err != nil {
		fmt.Printf("Failed to delete vector %s: %v\n", image.VectorID, err)
	}

	// Delete image
	if err := h.recordService.DeleteImage(uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted successfully"})
}

// GetImagePreview serves an image file for preview
func (h *RecordHandler) GetImagePreview(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// Get image metadata
	image, err := h.recordService.GetImage(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(image.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "image file not found"})
		return
	}

	// Serve the image file
	c.File(image.Path)
}
