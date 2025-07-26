package tool

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"image-rag-backend/internal/utils"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type imageRetrieveReq struct {
	ImagePath string `json:"image_path" description:"输入图片文件路径"`
}

type imageRetrieveResp struct {
	Record struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"record"`
}

func NewImageRagTool() *protocol.Tool {
	tool, err := protocol.NewTool(
		"image-rag-retrive",
		"根据输入图片文件路径，进行向量搜索，返回图片的具体信息，方便用于后续的图片内容解析和问答需求。",
		imageRetrieveReq{},
	)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	return tool
}

func ImageRagHandler(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(imageRetrieveReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}
	log.Printf("输入参数为 %#v", req)

	data, err := os.ReadFile(req.ImagePath)
	if err != nil {
		log.Printf("read image with error: %v", err)
		return nil, err
	}

	inputB64Image := base64.StdEncoding.EncodeToString(data)
	apiResponse, err := doImageRetrive(inputB64Image)
	if err != nil {
		return nil, errors.New("查询图片信息失败")
	}

	// Return the record name and description to the MCP client
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("图片名 \"%s\", 图片描述: \"%s\"", apiResponse.Record.Name, apiResponse.Record.Description),
			},
		},
	}, nil

}

func doImageRetrive(inputB64Image string) (*imageRetrieveResp, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Prepare request payload for the existing API endpoint
	requestPayload := map[string]interface{}{
		"image_base64": inputB64Image,
		"format":       "jpeg",
		"top_k":        1,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("failed to marshal request payload: %v", err)
		return nil, fmt.Errorf("failed to prepare request: %v", err)
	}

	// Make HTTP request to the existing API endpoint
	endpoint := utils.GetEnv("IMAGE_RAG_ENDPOINT", "http://localhost:8080")
	apiURL := endpoint + "/api/v1/search/record-by-image"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to call API endpoint: %v", err)
		return nil, fmt.Errorf("failed to search image: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response JSON
	var apiResponse imageRetrieveResp
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("failed to parse API response: %v", err)
		log.Printf("response body: %s", string(body))
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &apiResponse, nil
}
