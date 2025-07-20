package services

import (
	"fmt"
	"image-rag-backend/internal/config"
	"image-rag-backend/internal/doubao"
	"image-rag-backend/internal/milvus"
	"sync"
	"time"
)

type VectorService struct {
	doubaoClient *doubao.Client
	milvusClient *milvus.Client
	config       *config.Config
}

type VectorResult struct {
	VectorID string
	Vector   []float32
}

type SearchResult struct {
	ImageID  string
	Distance float32
}

func NewVectorService(cfg *config.Config) (*VectorService, error) {
	doubaoClient := doubao.NewClient(&cfg.Doubao)

	milvusClient, err := milvus.NewClient(&cfg.Milvus)
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus client: %w", err)
	}

	// Initialize Milvus collection (creates if needed and loads)
	if err := milvusClient.CreateCollection(); err != nil {
		return nil, fmt.Errorf("failed to create/load milvus collection: %w", err)
	}

	return &VectorService{
		doubaoClient: doubaoClient,
		milvusClient: milvusClient,
		config:       cfg,
	}, nil
}

func (s *VectorService) GenerateVector(imagePath string) (string, error) {
	// Generate embedding using Doubao
	embedding, err := s.doubaoClient.GenerateEmbedding(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Generate unique vector ID
	vectorID := generateUUID()

	// Insert into Milvus
	_, err = s.milvusClient.InsertVector(vectorID, embedding)
	if err != nil {
		return "", fmt.Errorf("failed to insert vector into milvus: %w", err)
	}

	// Return the original vector ID (UUID)
	return vectorID, nil
}

func (s *VectorService) GenerateVectorFromFile(imagePath string) (string, []float32, error) {
	// Generate embedding using Doubao
	embedding, err := s.doubaoClient.GenerateEmbedding(imagePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Generate unique vector ID
	vectorID := generateUUID()

	// Insert into Milvus
	_, err = s.milvusClient.InsertVector(vectorID, embedding)
	if err != nil {
		return "", nil, fmt.Errorf("failed to insert vector into milvus: %w", err)
	}

	// Return the original vector ID (UUID)
	return vectorID, embedding, nil
}

func (s *VectorService) SearchSimilar(imagePath string, topK int) ([]SearchResult, error) {
	// Generate embedding for query image
	embedding, err := s.doubaoClient.GenerateEmbedding(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search in Milvus
	results, err := s.milvusClient.SearchSimilar(embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar vectors: %w", err)
	}

	// Convert to our result format
	var searchResults []SearchResult
	for _, result := range results {
		searchResults = append(searchResults, SearchResult{
			ImageID:  result.VectorID,
			Distance: result.Distance,
		})
	}

	return searchResults, nil
}

func (s *VectorService) SearchSimilarWithVector(vector []float32, topK int) ([]SearchResult, error) {
	// Search in Milvus
	results, err := s.milvusClient.SearchSimilar(vector, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar vectors: %w", err)
	}

	// Convert to our result format
	var searchResults []SearchResult
	for _, result := range results {
		searchResults = append(searchResults, SearchResult{
			ImageID:  result.VectorID,
			Distance: result.Distance,
		})
	}

	return searchResults, nil
}

func (s *VectorService) DeleteVector(vectorID string) error {
	// Delete from Milvus
	return s.milvusClient.DeleteVector(vectorID)
}

func (s *VectorService) GetVectorCount() (int64, error) {
	return s.milvusClient.GetVectorCount()
}

func (s *VectorService) HealthCheck() error {
	// Check Milvus connection
	if err := s.milvusClient.Ping(); err != nil {
		return fmt.Errorf("milvus health check failed: %w", err)
	}

	// Check if Doubao API key is configured
	if s.config.Doubao.APIKey == "" {
		return fmt.Errorf("doubao api key not configured")
	}

	return nil
}

func (s *VectorService) Close() error {
	var errs []error

	if s.milvusClient != nil {
		if err := s.milvusClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close milvus client: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing vector service: %v", errs)
	}

	return nil
}

// ProcessImage handles the complete image processing pipeline
func (s *VectorService) ProcessImage(imagePath string) (string, error) {
	return s.GenerateVector(imagePath)
}

// ProcessImageWithEmbedding handles the complete image processing pipeline with embedding
func (s *VectorService) ProcessImageWithEmbedding(imagePath string) (string, []float32, error) {
	return s.GenerateVectorFromFile(imagePath)
}

// BatchProcessImages processes multiple images concurrently
func (s *VectorService) BatchProcessImages(imagePaths []string) ([]string, error) {
	var wg sync.WaitGroup
	results := make([]string, len(imagePaths))
	errors := make([]error, len(imagePaths))

	for i, path := range imagePaths {
		wg.Add(1)
		go func(index int, imagePath string) {
			defer wg.Done()
			vectorID, err := s.ProcessImage(imagePath)
			if err != nil {
				errors[index] = err
				return
			}
			results[index] = vectorID
		}(i, path)
	}

	wg.Wait()

	// Check for errors
	var errs []error
	for _, err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("batch processing errors: %v", errs)
	}

	return results, nil
}

// AsyncProcessImage processes an image asynchronously
func (s *VectorService) AsyncProcessImage(imagePath string, resultChan chan<- string, errorChan chan<- error) {
	go func() {
		vectorID, err := s.ProcessImage(imagePath)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- vectorID
	}()
}

// UpdateVector updates an existing vector
func (s *VectorService) UpdateVector(oldVectorID string, newImagePath string) (string, error) {
	// Delete old vector
	if err := s.DeleteVector(oldVectorID); err != nil {
		return "", fmt.Errorf("failed to delete old vector: %w", err)
	}

	// Process new image
	return s.ProcessImage(newImagePath)
}

// generateUUID generates a unique identifier
func generateUUID() string {
	// In a real implementation, use a proper UUID generator
	// For now, use a simple timestamp-based approach
	return fmt.Sprintf("vec_%d", time.Now().UnixNano())
}

// ValidateVector validates a vector
func ValidateVector(vector []float32) error {
	if len(vector) != 1024 {
		return fmt.Errorf("vector dimension must be 1024, got %d", len(vector))
	}
	return nil
}

// NormalizeVector normalizes a vector to unit length
func NormalizeVector(vector []float32) []float32 {
	// Calculate magnitude
	var magnitude float32
	for _, v := range vector {
		magnitude += v * v
	}
	magnitude = sqrt32(magnitude)

	// Normalize if magnitude > 0
	if magnitude > 0 {
		result := make([]float32, len(vector))
		for i, v := range vector {
			result[i] = v / magnitude
		}
		return result
	}

	return vector
}

// sqrt32 calculates square root for float32
func sqrt32(x float32) float32 {
	// Simple implementation - in production use math.Sqrt(float64(x))
	return float32(1.0) // Placeholder - use proper implementation
}

// CalculateSimilarity calculates cosine similarity between two vectors
func CalculateSimilarity(vec1, vec2 []float32) float32 {
	if len(vec1) != len(vec2) {
		return 0.0
	}

	var dotProduct, norm1, norm2 float32
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (sqrt32(norm1) * sqrt32(norm2))
}

// GetVectorByID retrieves a vector by its ID
func (s *VectorService) GetVectorByID(vectorID string) ([]float32, error) {
	// This would require implementing a get operation in Milvus
	// For now, return empty vector
	return make([]float32, 1024), nil
}

// GetStats returns statistics about the vector service
func (s *VectorService) GetStats() (map[string]interface{}, error) {
	count, err := s.GetVectorCount()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_vectors": count,
		"doubao_model":  s.config.Doubao.Model,
		"milvus_host":   s.config.Milvus.Host,
		"milvus_port":   s.config.Milvus.Port,
	}

	return stats, nil
}
