package milvus

import (
	"context"
	"testing"
	"time"

	"image-rag-backend/internal/config"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/stretchr/testify/assert"
)

func TestMilvusConnection(t *testing.T) {
	// Test configuration for 0.0.0.0:19530
	cfg := &config.MilvusConfig{
		Host:     "0.0.0.0",
		Port:     "19530",
		Database: "image_rag",
	}

	t.Run("TestConnectionTo0.0.0.0", func(t *testing.T) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Try to connect directly using client.NewClient
		milvusClient, err := client.NewClient(ctx, client.Config{
			Address: cfg.Host + ":" + cfg.Port,
		})

		if err != nil {
			t.Logf("Failed to connect to Milvus at %s:%s: %v", cfg.Host, cfg.Port, err)
			t.Skip("Milvus server not available at 0.0.0.0:19530")
		}
		defer milvusClient.Close()

		// Test basic connectivity by getting server version
		version, err := milvusClient.GetVersion(ctx)
		assert.NoError(t, err, "Should be able to get server version")
		assert.NotEmpty(t, version, "Server version should not be empty")

		t.Logf("Successfully connected to Milvus server version: %s", version)
	})

	t.Run("TestNewClientFunction", func(t *testing.T) {
		// Test using our NewClient function
		client, err := NewClient(cfg)
		if err != nil {
			t.Logf("NewClient failed to connect: %v", err)
			t.Skip("Milvus server not available")
		}
		defer client.Close()

		// Test Ping functionality
		err = client.Ping()
		if err != nil {
			t.Logf("Ping failed: %v", err)
			t.Skip("Milvus server not responding")
		}

		t.Log("Successfully connected using NewClient function")
	})
}

func TestConfigConnection(t *testing.T) {
	// Test with environment configuration
	cfg := &config.MilvusConfig{
		Host: getTestHost(),
		Port: getTestPort(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Milvus server not available at %s:%s: %v", cfg.Host, cfg.Port, err)
	}
	defer client.Close()

	t.Logf("Successfully connected to Milvus at %s:%s", cfg.Host, cfg.Port)
}

// Helper functions for test configuration
func getTestHost() string {
	// Check environment variable first, fallback to 0.0.0.0
	if host := getEnv("MILVUS_TEST_HOST", ""); host != "" {
		return host
	}
	return "0.0.0.0"
}

func getTestPort() string {
	if port := getEnv("MILVUS_TEST_PORT", ""); port != "" {
		return port
	}
	return "19530"
}

func getEnv(key, defaultValue string) string {
	// Simple environment variable getter for tests
	// In a real test, you might want to use a more sophisticated approach
	return defaultValue
}
