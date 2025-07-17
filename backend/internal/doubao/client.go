package doubao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image-rag-backend/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Client struct {
	apiKey string
	model  string
	url    string
	client *http.Client
}

type EmbeddingRequest struct {
	Model string `json:"model"`
	Input struct {
		Images []ImageData `json:"images"`
	} `json:"input"`
}

type ImageData struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
		Object    string    `json:"object"`
	} `json:"data"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func NewClient(cfg *config.DoubaoConfig) *Client {
	return &Client{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		url:    cfg.URL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GenerateEmbedding(imagePath string) ([]float32, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("doubao api key is required")
	}

	// Read and encode image
	imageData, err := encodeImageToBase64(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	// Get image format
	format := getImageFormat(imagePath)

	// Prepare request
	req := EmbeddingRequest{
		Model: c.model,
	}
	req.Input.Images = []ImageData{
		{
			Data:   imageData,
			Format: format,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("doubao api error: %s", string(body))
	}

	// Parse response
	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	return response.Data[0].Embedding, nil
}

func (c *Client) GenerateEmbeddingFromFile(file multipart.File, filename string) ([]float32, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("doubao api key is required")
	}

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Encode to base64
	imageData := base64.StdEncoding.EncodeToString(fileContent)

	// Get image format
	format := getImageFormat(filename)

	// Prepare request
	req := EmbeddingRequest{
		Model: c.model,
	}
	req.Input.Images = []ImageData{
		{
			Data:   imageData,
			Format: format,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("doubao api error: %s", string(body))
	}

	// Parse response
	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	return response.Data[0].Embedding, nil
}

// encodeImageToBase64 reads an image file and returns base64 encoded data
func encodeImageToBase64(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

// getImageFormat determines the image format from file extension
func getImageFormat(filename string) string {
	ext := getFileExtension(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".webp":
		return "webp"
	default:
		return "jpeg" // default fallback
	}
}

// getFileExtension returns the file extension
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}
