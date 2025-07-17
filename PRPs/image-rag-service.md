name: "Image RAG Service - Vector-based Image Search with Management Interface"
description: |
  Complete implementation of an Image RAG (Retrieval Augmented Generation) service for vector-based image search, including web-based management interface and automated vectorization pipeline.

## Purpose
Build a production-ready Image RAG service that enables users to upload images and retrieve similar images based on content similarity, with a comprehensive management interface for database operations.

## Goal
Create a complete Image RAG service with:
1. **Backend API**: Go + Gin REST API for record management and image search
2. **Vector Service**: Python service for image vectorization using Doubao API
3. **Frontend**: Vue.js SPA for record management and visual search interface
4. **Database**: MySQL for metadata, Milvus for vector storage
5. **Integration**: Seamless pipeline from image upload to vector search

## Why
- **Business Value**: Enable content-based image retrieval for applications like e-commerce, digital asset management, and visual search
- **Technical Innovation**: Leverage cutting-edge vision models for semantic image understanding
- **User Experience**: Provide intuitive web interface for non-technical users to manage image databases
- **Scalability**: Microservices architecture allows independent scaling of components

## What
### User-Visible Features:
- **Record Management**: Create, read, update, delete records with multiple images
- **Visual Search**: Upload query image to find similar images from database
- **Admin Interface**: Web-based management of image database and search results
- **Batch Operations**: Upload multiple images at once

### Technical Requirements:
- **API Endpoints**: Complete CRUD for records, image upload, similarity search
- **Vector Pipeline**: Automatic image vectorization using Doubao doubao-embedding-vision-250615
- **Database Design**: MySQL for metadata, Milvus for high-dimensional vectors
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Performance**: Async processing for large images, caching for frequent searches

### Success Criteria
- [ ] All API endpoints return correct responses with proper HTTP status codes
- [ ] Images are successfully vectorized and stored in Milvus
- [ ] Search returns relevant similar images with similarity scores
- [ ] Frontend interface allows complete record management workflow
- [ ] Error cases handled gracefully with informative messages
- [ ] All tests pass (unit, integration, and end-to-end)
- [ ] Performance: <2s response time for search with 1000+ images

## All Needed Context

### Documentation & References
```yaml
# MUST READ - Include these in your context window
- url: https://www.volcengine.com/docs/82379/1523520
  why: Doubao API documentation for doubao-embedding-vision-250615 model usage, authentication, and request/response formats
  
- url: https://milvus.io/docs
  why: Milvus documentation for vector collection setup, indexing strategies (IVF_FLAT, HNSW), and search optimization
  
- url: https://gin-gonic.com/docs/
  why: Gin framework documentation for REST API design, middleware, and file upload handling
  
- url: https://gorm.io/docs/
  why: GORM documentation for MySQL integration, associations (one-to-many), and migrations
  
- url: https://vuejs.org/guide/
  why: Vue.js 3 Composition API documentation for component design and reactivity
  
- url: https://element-plus.org/en-US/
  why: Element Plus UI library documentation for form components, file upload, and layout

# Code Examples from Project
- file: examples/image-rag-example.md
  why: Complete implementation examples for Go backend, Python vector service, and Vue.js frontend
  
- file: CLAUDE.md
  why: Project-wide rules including file size limits, testing requirements, and naming conventions
```

### Current Codebase Tree
```
image-rag/
├── PRPs/
│   ├── templates/
│   │   └── prp_base.md
│   └── EXAMPLE_multi_agent_prp.md
├── examples/
│   └── image-rag-example.md
├── CLAUDE.md
├── README.md
├── INITIAL.md
├── INITIAL_EXAMPLE.md
└── LICENSE
```

### Desired Codebase Tree
```
image-rag/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── records.go
│   │   │   │   └── search.go
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go
│   │   │   │   └── rate_limit.go
│   │   │   └── routes.go
│   │   ├── models/
│   │   │   └── record.go
│   │   ├── services/
│   │   │   ├── record_service.go
│   │   │   └── vector_service.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── database/
│   │       └── connection.go
│   ├── uploads/
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
├── vector-service/
│   ├── src/
│   │   ├── services/
│   │   │   └── vector_service.py
│   │   ├── api/
│   │   │   └── app.py
│   │   └── models/
│   │       └── schemas.py
│   ├── tests/
│   │   └── test_vector_service.py
│   ├── requirements.txt
│   ├── .env.example
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ImageUpload.vue
│   │   │   ├── ImageSearch.vue
│   │   │   └── RecordList.vue
│   │   ├── views/
│   │   │   ├── Records.vue
│   │   │   └── Search.vue
│   │   ├── api/
│   │   │   └── records.js
│   │   ├── router/
│   │   │   └── index.js
│   │   └── main.js
│   ├── public/
│   ├── package.json
│   ├── vite.config.js
│   └── .env.example
├── docker-compose.yml
├── .env.example
└── README.md
```

### Known Gotchas & Library Quirks
```python
# CRITICAL: Doubao API requires base64 encoded images with specific format
# Supported formats: JPEG, PNG, WebP (max 10MB per image)
# Base64 string must be properly formatted without data URI prefix

# CRITICAL: Milvus vector dimensions must match model output
# doubao-embedding-vision-250615 outputs 1024-dimensional vectors
# Collection schema must specify dim=1024

# CRITICAL: Go file upload requires proper multipart form handling
# Use gin.Context.FormFile() and c.SaveUploadedFile() for file handling
# Set appropriate max memory limits: router.MaxMultipartMemory = 64 << 20

# CRITICAL: CORS configuration required for frontend-backend communication
# Allow origins, methods, and headers for file uploads
# Use github.com/gin-contrib/cors middleware

# CRITICAL: Async processing for image vectorization
# Use goroutines for non-blocking vector generation
# Implement proper error handling and cleanup for failed operations

# CRITICAL: Database consistency between MySQL and Milvus
# Use transactions where possible
# Implement cleanup for orphaned vectors when records are deleted

# CRITICAL: Image file validation
# Check file extensions (.jpg, .jpeg, .png, .webp)
# Validate file size (max 10MB)
# Verify image format using Go's image.DecodeConfig()
```

## Implementation Blueprint

### Data Models and Structure

#### MySQL Models (GORM)
```go
// internal/models/record.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Record struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null;size:255"`
    Description string    `json:"description" gorm:"type:text"`
    Images      []Image   `json:"images" gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    RecordID uint   `json:"record_id" gorm:"not null;index"`
    Filename string `json:"filename" gorm:"not null;size:255"`
    Path     string `json:"path" gorm:"not null;size:500"`
    VectorID string `json:"vector_id" gorm:"not null;size:100;index"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Milvus Collection Schema
```python
# vector_service/models/collection.py
from pymilvus import CollectionSchema, FieldSchema, DataType

def create_image_collection_schema():
    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
        FieldSchema(name="image_id", dtype=DataType.VARCHAR, max_length=100),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1024)
    ]
    
    schema = CollectionSchema(
        fields=fields,
        description="Image embeddings for similarity search"
    )
    return schema
```

### Task List (in order)

#### Task 1: Database Setup and Configuration
```yaml
CREATE backend/internal/config/config.go:
  - Load environment variables
  - Database connection strings
  - API configuration
  - File upload settings

CREATE backend/internal/database/connection.go:
  - MySQL connection setup
  - Database migrations
  - Connection pooling

CREATE vector-service/.env.example:
  - Doubao API key configuration
  - Milvus connection settings
  - Python service configuration
```

#### Task 2: Backend API Structure
```yaml
CREATE backend/cmd/server/main.go:
  - Initialize Gin router
  - Setup middleware (CORS, rate limiting)
  - Initialize database connections
  - Start HTTP server

CREATE backend/internal/models/record.go:
  - Record and Image structs with GORM tags
  - Database relationships
  - JSON serialization tags

CREATE backend/internal/api/middleware/cors.go:
  - CORS configuration for frontend
  - Allow file uploads
  - Security headers

CREATE backend/internal/api/middleware/rate_limit.go:
  - Rate limiting middleware
  - Configurable limits per endpoint
  - IP-based limiting
```

#### Task 3: Service Layer
```yaml
CREATE backend/internal/services/record_service.go:
  - CRUD operations for records
  - Image file management
  - Database transactions

CREATE backend/internal/services/vector_service.go:
  - HTTP client for vector service
  - Retry logic for API calls
  - Error handling

CREATE vector-service/src/services/vector_service.py:
  - Doubao API integration
  - Image preprocessing
  - Vector storage in Milvus
```

#### Task 4: API Handlers
```yaml
CREATE backend/internal/api/handlers/records.go:
  - POST /records - Create record with images
  - GET /records - List all records
  - GET /records/:id - Get single record
  - PUT /records/:id - Update record
  - DELETE /records/:id - Delete record

CREATE backend/internal/api/handlers/search.go:
  - POST /search - Search similar images
  - GET /search/similar/:id - Find similar to existing image

CREATE backend/internal/api/routes.go:
  - Route definitions
  - Handler registration
  - Static file serving for uploads
```

#### Task 5: Frontend Components
```yaml
CREATE frontend/src/components/ImageUpload.vue:
  - Multi-file upload with preview
  - Form validation
  - Progress indicators

CREATE frontend/src/components/ImageSearch.vue:
  - Drag-and-drop image upload
  - Search results display
  - Similarity scores

CREATE frontend/src/components/RecordList.vue:
  - Paginated record listing
  - Image gallery view
  - CRUD operations

CREATE frontend/src/views/Records.vue:
  - Main records management page
  - Integrate upload and list components

CREATE frontend/src/views/Search.vue:
  - Dedicated search interface
  - Results visualization
```

#### Task 6: Testing and Validation
```yaml
CREATE backend/tests/record_service_test.go:
  - Unit tests for record operations
  - Mock database interactions
  - Error scenario testing

CREATE vector-service/tests/test_vector_service.py:
  - Mock Doubao API responses
  - Milvus integration tests
  - Image processing tests

CREATE docker-compose.yml:
  - MySQL service
  - Milvus service
  - Backend service
  - Frontend service
  - Python vector service
```

### Integration Points
```yaml
DATABASE:
  - MySQL schema: records and images tables
  - Indexes: record_id on images table, vector_id for lookups
  - Foreign key constraints with CASCADE delete

VECTOR DATABASE:
  - Milvus collection: image_embeddings
  - Index type: IVF_FLAT or HNSW for performance
  - Vector dimension: 1024 (matches Doubao output)

CONFIGURATION:
  - Environment variables for all services
  - CORS origins for frontend-backend communication
  - File upload limits and storage paths

SECURITY:
  - API key management for Doubao
  - File type validation
  - Rate limiting per IP
  - CORS configuration
```

## Validation Loop

### Level 1: Backend Syntax & Style
```bash
cd backend
# Install dependencies
go mod init image-rag-backend
go mod tidy

# Format code
go fmt ./...

# Run vet
go vet ./...

# Build check
go build -o server cmd/server/main.go
```

### Level 2: Backend Unit Tests
```bash
cd backend
# Run tests
go test ./... -v

# Coverage check
go test ./... -cover
```

### Level 3: Python Vector Service
```bash
cd vector-service
# Install dependencies
pip install -r requirements.txt

# Run tests
python -m pytest tests/ -v

# Type checking
mypy src/

# Code style
flake8 src/
```

### Level 4: Frontend Build
```bash
cd frontend
# Install dependencies
npm install

# Type checking
npm run type-check

# Build
npm run build

# Linting
npm run lint
```

### Level 5: Integration Testing
```bash
# Start all services
docker-compose up -d

# Wait for services to be ready
sleep 30

# Test backend health
curl http://localhost:8080/health

# Test record creation
curl -X POST http://localhost:8080/api/v1/records \
  -F "name=test-record" \
  -F "description=test description" \
  -F "images=@test-image.jpg"

# Test search
curl -X POST http://localhost:8080/api/v1/search \
  -F "image=@query-image.jpg"

# Test frontend
open http://localhost:3000
```

### Level 6: End-to-End Testing
```bash
# Create test record with images
node scripts/create-test-record.js

# Perform visual search
node scripts/test-search.js

# Verify results
node scripts/verify-results.js
```

## Final Validation Checklist
- [ ] All services start successfully with docker-compose
- [ ] MySQL database initialized with correct schema
- [ ] Milvus collection created with proper index
- [ ] Backend API responds to all CRUD operations
- [ ] Image upload works with file size/type validation
- [ ] Vector generation succeeds for uploaded images
- [ ] Similarity search returns relevant results
- [ ] Frontend displays records and search results correctly
- [ ] Error handling works for all failure scenarios
- [ ] All unit tests pass with >80% coverage
- [ ] Performance benchmarks meet requirements
- [ ] Security measures implemented (CORS, rate limiting, input validation)
- [ ] Documentation complete (API docs, setup instructions)

## Anti-Patterns to Avoid
- ❌ Don't store images as BLOBs in MySQL - use filesystem with path references
- ❌ Don't process images synchronously - use goroutines for vectorization
- ❌ Don't skip input validation - validate file types, sizes, and content
- ❌ Don't hardcode API keys - use environment variables
- ❌ Don't ignore database migrations - use GORM auto-migration safely
- ❌ Don't create monolithic handlers - separate concerns into services
- ❌ Don't skip error logging - implement comprehensive logging with levels
- ❌ Don't forget cleanup - remove temporary files and orphaned vectors
- ❌ Don't skip rate limiting - implement per-IP and global limits
- ❌ Don't ignore CORS - configure properly for frontend communication

---

## Confidence Score: 9/10
This PRP provides comprehensive context for implementing the Image RAG service. The architecture is well-defined, examples are provided, and validation steps are clear. The modular design allows for iterative development and testing.