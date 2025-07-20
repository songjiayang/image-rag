package milvus

import (
	"context"
	"fmt"
	"image-rag-backend/internal/config"
	"strconv"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Client struct {
	client client.Client
	cfg    *config.MilvusConfig
	ctx    context.Context
}

type VectorData struct {
	VectorID string
	Vector   []float32
}

type SearchResult struct {
	VectorID string
	Distance float32
}

func NewClient(cfg *config.MilvusConfig) (*Client, error) {
	ctx := context.Background()

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cli, err := client.NewClient(timeout, client.Config{
		Address: cfg.Host + ":" + cfg.Port,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to milvus at %s:%s %w", cfg.Host, cfg.Port, err)
	}

	return &Client{
		client: cli,
		cfg:    cfg,
		ctx:    ctx,
	}, nil
}

// CreateCollection creates the image embeddings collection
func (c *Client) CreateCollection() error {
	// Check if collection already exists
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	exists, err := c.client.HasCollection(ctx, "image_embeddings")
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if exists {
		// Collection already exists, ensure it's loaded
		return c.LoadCollection()
	}

	// Define schema
	schema := &entity.Schema{
		CollectionName: "image_embeddings",
		Description:    "Image embeddings for similarity search",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "image_id",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: "100",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					entity.TypeParamDim: "1024",
				},
			},
		},
	}

	// Create collection
	if err := c.client.CreateCollection(ctx, schema, 2); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// Create index using proper Milvus index constructor
	index, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	if err := c.client.CreateIndex(ctx, "image_embeddings", "embedding", index, false); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Load collection
	return c.LoadCollection()
}

// InsertVector inserts a vector into the collection
func (c *Client) InsertVector(imageID string, vector []float32) (int64, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	// Prepare data
	ids := []string{imageID}
	vectors := [][]float32{vector}

	// Insert data
	columnID := entity.NewColumnVarChar("image_id", ids)
	columnVector := entity.NewColumnFloatVector("embedding", 1024, vectors)

	_, err := c.client.Insert(c.ctx, "image_embeddings", "", columnID, columnVector)
	if err != nil {
		return 0, fmt.Errorf("failed to insert vector: %w", err)
	}

	// Flush to ensure data is persisted
	if err := c.client.Flush(ctx, "image_embeddings", false); err != nil {
		return 0, fmt.Errorf("failed to flush collection: %w", err)
	}

	// Since we're not using auto-generated IDs, return 1 to indicate success
	// The actual ID is the image_id we provided
	return 1, nil
}

// SearchSimilar searches for similar vectors
func (c *Client) SearchSimilar(vector []float32, topK int) ([]SearchResult, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	if topK <= 0 {
		topK = 10
	}

	// Search parameters using index params
	searchParams, err := entity.NewIndexIvfFlatSearchParam(10)
	if err != nil {
		return nil, fmt.Errorf("failed to create search parameters: %w", err)
	}

	// Perform search
	results, err := c.client.Search(
		ctx,
		"image_embeddings",
		[]string{},
		"",
		[]string{"image_id"},
		[]entity.Vector{entity.FloatVector(vector)},
		"embedding",
		entity.L2,
		topK,
		searchParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// Parse results
	var searchResults []SearchResult
	for _, result := range results {
		if len(result.Fields) > 0 {
			for _, field := range result.Fields {
				if col, ok := field.(*entity.ColumnVarChar); ok {
					data := col.Data()
					for i, imageID := range data {
						var score float32
						if len(result.Scores) > i {
							score = result.Scores[i]
						}
						searchResults = append(searchResults, SearchResult{
							VectorID: imageID,
							Distance: score,
						})
					}
					break
				}
			}
		}
	}

	return searchResults, nil
}

// DeleteVector deletes a vector by image ID
func (c *Client) DeleteVector(imageID string) error {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	// This is a simplified delete operation
	// In production, you might want to use a more sophisticated approach
	expr := fmt.Sprintf(`image_id == "%s"`, imageID)

	return c.client.Delete(ctx, "image_embeddings", "", expr)
}

// GetVectorCount returns the total number of vectors
func (c *Client) GetVectorCount() (int64, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	stats, err := c.client.GetCollectionStatistics(ctx, "image_embeddings")
	if err != nil {
		return 0, fmt.Errorf("failed to get collection statistics: %w", err)
	}

	count, ok := stats["row_count"]
	if !ok {
		return 0, fmt.Errorf("row count not found in statistics")
	}

	return strconv.ParseInt(count, 10, 64)
}

// Close closes the Milvus client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// Ping checks if the Milvus server is available
func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	// Use GetCollectionStatistics as a simple health check
	_, err := c.client.GetCollectionStatistics(ctx, "image_embeddings")
	return err
}

// DeleteByExpr deletes vectors by expression
func (c *Client) DeleteByExpr(expr string) error {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.client.Delete(ctx, "image_embeddings", "", expr)
}

// LoadCollection loads the collection into memory
func (c *Client) LoadCollection() error {
	ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer cancel()

	// Check if collection is already loaded
	loaded, err := c.client.GetLoadState(ctx, "image_embeddings", []string{})
	if err != nil {
		return fmt.Errorf("failed to check load state: %w", err)
	}

	if loaded == entity.LoadStateLoaded {
		return nil // Already loaded
	}

	// Load collection with retry
	var lastErr error
	for i := 0; i < 3; i++ {
		if err := c.client.LoadCollection(ctx, "image_embeddings", false); err != nil {
			lastErr = fmt.Errorf("failed to load collection (attempt %d): %w", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return lastErr
}

// ReleaseCollection releases the collection from memory
func (c *Client) ReleaseCollection() error {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.client.ReleaseCollection(ctx, "image_embeddings")
}

// HasCollection checks if collection exists
func (c *Client) HasCollection(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.client.HasCollection(ctx, name)
}

// DropCollection drops the collection
func (c *Client) DropCollection(name string) error {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.client.DropCollection(ctx, name)
}
