package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"image-rag-backend/internal/logger"
	"image-rag-backend/internal/models"
	"image-rag-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	recordService *services.RecordService
	vectorService *services.VectorService
	logger        *logger.Logger
}

type SearchResult struct {
	RecordID    uint    `json:"record_id"`
	RecordName  string  `json:"record_name"`
	Description string  `json:"description"`
	ImageID     uint    `json:"image_id"`
	Filename    string  `json:"filename"`
	Distance    float64 `json:"distance"`
}

func NewSearchHandler(recordService *services.RecordService, vectorService *services.VectorService, logger *logger.Logger) *SearchHandler {
	return &SearchHandler{
		recordService: recordService,
		vectorService: vectorService,
		logger:        logger,
	}
}

// SearchImages searches for similar images based on query image
// @Summary Search similar images
// @Description Upload an image to search for similar images in the database using vector similarity
// @Tags Search
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to search for"
// @Param top_k query int false "Number of results to return (default: 10, max: 100)"
// @Success 200 {object} SearchResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [post]
func (h *SearchHandler) SearchImages(c *gin.Context) {
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

	// Get top_k parameter
	topK, _ := strconv.Atoi(c.DefaultQuery("top_k", "10"))
	if topK < 1 || topK > 100 {
		topK = 10
	}

	// Save temporary file
	filename := services.GenerateUniqueFilename(header.Filename)
	tempPath := filepath.Join("uploads", "temp", filename)

	// Ensure temp directory exists
	services.EnsureDirectoryExists(filepath.Join("uploads", "temp"))

	// Save file
	if err := c.SaveUploadedFile(header, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Clean up temp file after processing
	defer func() {
		_ = services.NewRecordService().DeleteImageByPath(tempPath)
	}()

	// Search for similar images
	results, err := h.vectorService.SearchSimilar(tempPath, topK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get record information for each result
	var searchResults []SearchResult
	for _, result := range results {
		// Find image by vector ID
		image, err := h.findImageByVectorID(result.ImageID)
		if err != nil {
			continue // Skip if image not found
		}

		// Get record information
		record, err := h.recordService.GetRecord(image.RecordID)
		if err != nil {
			continue // Skip if record not found
		}

		searchResults = append(searchResults, SearchResult{
			RecordID:    record.ID,
			RecordName:  record.Name,
			Description: record.Description,
			ImageID:     image.ID,
			Filename:    image.Filename,
			Distance:    float64(result.Distance),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": searchResults,
		"count":   len(searchResults),
	})
}

// GetImageByVectorID retrieves image information by vector ID
// @Summary Get image by vector ID
// @Description Retrieve image information using its vector ID from Milvus
// @Tags Search
// @Accept json
// @Produce json
// @Param vector_id path string true "Vector ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search/by-vector/{vector_id} [get]
func (h *SearchHandler) GetImageByVectorID(c *gin.Context) {
	vectorID := c.Param("vector_id")
	if vectorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vector ID is required"})
		return
	}

	// Find image by vector ID
	image, err := h.findImageByVectorID(vectorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get record information
	record, err := h.recordService.GetRecord(image.RecordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image":  image,
		"record": record,
	})
}

// FindSimilar finds similar images to an existing image
func (h *SearchHandler) FindSimilar(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// Get image information
	image, err := h.recordService.GetImage(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	// Get top_k parameter
	topK, _ := strconv.Atoi(c.DefaultQuery("top_k", "10"))
	if topK < 1 || topK > 100 {
		topK = 10
	}

	// Get image path
	imagePath := filepath.Join("uploads", image.Filename)
	if _, err := services.NewRecordService().FileExists(imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image file not found"})
		return
	}

	// Search for similar images
	results, err := h.vectorService.SearchSimilar(imagePath, topK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get record information for each result
	var searchResults []SearchResult
	for _, result := range results {
		// Skip the same image
		if result.ImageID == image.VectorID {
			continue
		}

		// Find image by vector ID
		similarImage, err := h.findImageByVectorID(result.ImageID)
		if err != nil {
			continue // Skip if image not found
		}

		// Get record information
		record, err := h.recordService.GetRecord(similarImage.RecordID)
		if err != nil {
			continue // Skip if record not found
		}

		searchResults = append(searchResults, SearchResult{
			RecordID:    record.ID,
			RecordName:  record.Name,
			Description: record.Description,
			ImageID:     similarImage.ID,
			Filename:    similarImage.Filename,
			Distance:    float64(result.Distance),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": searchResults,
		"count":   len(searchResults),
	})
}

// AdvancedSearch performs advanced search with filters
func (h *SearchHandler) AdvancedSearch(c *gin.Context) {
	// Get search parameters
	query := c.Query("q")
	recordName := c.Query("record_name")
	minDistance, _ := strconv.ParseFloat(c.DefaultQuery("min_distance", "0"), 32)
	maxDistance, _ := strconv.ParseFloat(c.DefaultQuery("max_distance", "1"), 32)
	topK, _ := strconv.Atoi(c.DefaultQuery("top_k", "10"))

	if topK < 1 || topK > 100 {
		topK = 10
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

	// Save temporary file
	filename := services.GenerateUniqueFilename(header.Filename)
	tempPath := filepath.Join("uploads", "temp", filename)

	// Ensure temp directory exists
	services.EnsureDirectoryExists(filepath.Join("uploads", "temp"))

	// Save file
	if err := c.SaveUploadedFile(header, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Clean up temp file after processing
	defer func() {
		_ = services.NewRecordService().DeleteImageByPath(tempPath)
	}()

	// Search for similar images
	results, err := h.vectorService.SearchSimilar(tempPath, topK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get record information and apply filters
	var searchResults []SearchResult
	for _, result := range results {
		// Check distance range
		if result.Distance < float32(minDistance) || result.Distance > float32(maxDistance) {
			continue
		}

		// Find image by vector ID
		image, err := h.findImageByVectorID(result.ImageID)
		if err != nil {
			continue
		}

		// Get record information
		record, err := h.recordService.GetRecord(image.RecordID)
		if err != nil {
			continue
		}

		// Apply text filters
		if query != "" && !strings.Contains(strings.ToLower(record.Description), strings.ToLower(query)) {
			continue
		}

		if recordName != "" && !strings.Contains(strings.ToLower(record.Name), strings.ToLower(recordName)) {
			continue
		}

		searchResults = append(searchResults, SearchResult{
			RecordID:    record.ID,
			RecordName:  record.Name,
			Description: record.Description,
			ImageID:     image.ID,
			Filename:    image.Filename,
			Distance:    float64(result.Distance),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": searchResults,
		"count":   len(searchResults),
		"query":   query,
		"filters": gin.H{
			"record_name":  recordName,
			"min_distance": minDistance,
			"max_distance": maxDistance,
		},
	})
}

// Helper function to find image by vector ID
func (h *SearchHandler) findImageByVectorID(vectorID string) (*models.Image, error) {
	// Query database for image with matching vector ID
	var image models.Image
	if err := h.recordService.GetDB().Where("vector_id = ?", vectorID).First(&image).Error; err != nil {
		return nil, fmt.Errorf("image not found for vector ID: %s", vectorID)
	}
	return &image, nil
}

// Base64SearchRequest represents the request structure for base64 image search
// @Base64SearchRequest represents the request structure for base64 image search
type Base64SearchRequest struct {
	Base64Data string `json:"image_base64" binding:"required"`
	Format     string `json:"format" binding:"omitempty"`
	TopK       int    `json:"top_k" binding:"omitempty,min=1,max=100"`
}

// SearchByBase64 searches for similar images using base64 image data
// @Summary Search similar images using base64
// @Description Search for similar images using base64 encoded image data
// @Tags Search
// @Accept json
// @Produce json
// @Param search body Base64SearchRequest true "Base64 image data and search parameters"
// @Success 200 {object} SearchResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search/base64 [post]
func (h *SearchHandler) SearchByBase64(c *gin.Context) {
	var req Base64SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default top_k if not provided
	if req.TopK == 0 {
		req.TopK = 10
	}

	// Validate base64 data
	if req.Base64Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_base64 is required"})
		return
	}

	// Search for similar images using base64 data
	results, err := h.vectorService.SearchSimilarFromBase64(req.Base64Data, req.Format, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get record information for each result
	var searchResults []SearchResult
	for _, result := range results {
		// Find image by vector ID
		image, err := h.findImageByVectorID(result.ImageID)
		if err != nil {
			continue // Skip if image not found
		}

		// Get record information
		record, err := h.recordService.GetRecord(image.RecordID)
		if err != nil {
			continue // Skip if record not found
		}

		searchResults = append(searchResults, SearchResult{
			RecordID:    record.ID,
			RecordName:  record.Name,
			Description: record.Description,
			ImageID:     image.ID,
			Filename:    image.Filename,
			Distance:    float64(result.Distance),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": searchResults,
		"count":   len(searchResults),
		"format":  req.Format,
	})
}

// GetRecordDetailsByImage searches for a single image by base64 and returns the most similar record details
// @Summary Get record details by base64 image
// @Description Upload a base64 image and return the most similar record's name and description
// @Tags Search
// @Accept json
// @Produce json
// @Param image body Base64SearchRequest true "Base64 image data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search/record-by-image [post]
func (h *SearchHandler) GetRecordDetailsByImage(c *gin.Context) {
	var req Base64SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate base64 data
	if req.Base64Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_base64 is required"})
		return
	}

	// Set default to get only the best match
	if req.TopK == 0 {
		req.TopK = 1
	}

	// Search for similar images using base64 data
	results, err := h.vectorService.SearchSimilarFromBase64(req.Base64Data, req.Format, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no matching records found"})
		return
	}

	// Get the best match (first result)
	bestMatch := results[0]

	// Find image by vector ID
	image, err := h.findImageByVectorID(bestMatch.ImageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found for vector ID"})
		return
	}

	// Get record information
	record, err := h.recordService.GetRecord(image.RecordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"record": gin.H{
			"id":          record.ID,
			"name":        record.Name,
			"description": record.Description,
			"created_at":  record.CreatedAt,
			"updated_at":  record.UpdatedAt,
		},
		"image": gin.H{
			"id":       image.ID,
			"filename": image.Filename,
			"distance": bestMatch.Distance,
		},
	})
}
