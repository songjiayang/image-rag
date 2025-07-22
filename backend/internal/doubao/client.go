package doubao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"image-rag-backend/internal/config"
)

type Client struct {
	apiKey string
	model  string
	url    string
	client *http.Client
}

type EmbeddingRequest struct {
	Model      string      `json:"model"`
	Input      []ImageData `json:"input"`
	Dimensions int         `json:"dimensions"`
}

type ImageData struct {
	Type     string   `json:"type"`
	ImageUrl ImageUrl `json:"image_url"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

type EmbeddingResponse struct {
	Created int64 `json:"created"`
	Data    struct {
		Embedding []float32 `json:"embedding"`
		Object    string    `json:"object"`
	} `json:"data"`
	ID     string `json:"id"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens        int `json:"prompt_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			ImageTokens int `json:"image_tokens"`
			TextTokens  int `json:"text_tokens"`
		} `json:"prompt_tokens_details"`
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

	return c.generateEmbeddingFromData(imageData, format)
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

	return c.generateEmbeddingFromData(imageData, format)
}

// GenerateEmbeddingFromBase64 processes base64 image data and returns embedding
func (c *Client) GenerateEmbeddingFromBase64(base64Data string, format string) ([]float32, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("doubao api key is required")
	}

	// Clean base64 data (remove data URL prefix if present)
	cleanBase64 := cleanBase64Data(base64Data)

	// Validate format
	if format == "" {
		format = "jpeg" // default fallback
	}

	return c.generateEmbeddingFromData(cleanBase64, format)
}

// generateEmbeddingFromData is a private method that handles the common HTTP request logic
// for generating embeddings from base64 encoded image data
func (c *Client) generateEmbeddingFromData(base64Data string, format string) ([]float32, error) {
	// Prepare request
	req := EmbeddingRequest{
		Model: c.model,
		Input: []ImageData{
			{
				Type: "image_url",
				ImageUrl: ImageUrl{
					Url: fmt.Sprintf("data:image/%s;base64,%s", format, base64Data),
				},
			},
		},
		Dimensions: 1024,
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

	if response.Data.Object != "embedding" {
		return nil, fmt.Errorf("unexpected response object type: %s", response.Data.Object)
	}

	if len(response.Data.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding in response")
	}

	return response.Data.Embedding, nil
}

// cleanBase64Data removes data URL prefix from base64 string if present
func cleanBase64Data(base64Data string) string {
	// Remove data URL prefix like "data:image/jpeg;base64,"
	if idx := strings.Index(base64Data, "base64,"); idx != -1 {
		return base64Data[idx+7:]
	}
	return base64Data
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
