# Image RAG Service API Documentation

## Overview
RESTful API for image-based Retrieval-Augmented Generation service with vector similarity search using Doubao AI and Milvus.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently uses API key authentication (configured via environment variable).

## Content Types
- Request: `application/json`, `multipart/form-data`
- Response: `application/json`

## Error Handling
All errors follow a consistent format:
```json
{
  "error": "Error message description"
}
```

## Rate Limiting
- Default: 100 requests per minute per IP
- Configurable via environment variables

## Endpoints

### Health Check
```
GET /api/v1/health
```

### Records Management

#### Create Record with Images
```
POST /api/v1/records
Content-Type: multipart/form-data

Parameters:
- name (string, required): Record name
- description (string, optional): Record description
- images (files, required): Image files to upload

Response: 201 Created
{
  "id": 1,
  "name": "Sample Record",
  "description": "Description of the record",
  "images": [
    {
      "id": 1,
      "filename": "image1.jpg",
      "path": "./uploads/image1.jpg"
    }
  ],
  "created_at": "2025-07-17T10:00:00Z",
  "updated_at": "2025-07-17T10:00:00Z"
}
```

#### List Records
```
GET /api/v1/records?page=1&limit=10

Parameters:
- page (int, optional): Page number (default: 1)
- limit (int, optional): Items per page (default: 10, max: 100)

Response: 200 OK
{
  "data": [...],
  "total": 25,
  "page": 1,
  "limit": 10
}
```

#### Get Record
```
GET /api/v1/records/{id}

Response: 200 OK
{
  "id": 1,
  "name": "Sample Record",
  "description": "Description of the record",
  "images": [...],
  "created_at": "...",
  "updated_at": "..."
}
```

#### Update Record
```
PUT /api/v1/records/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "description": "Updated description"
}

Response: 200 OK
{
  "id": 1,
  "name": "Updated Name",
  "description": "Updated description",
  ...
}
```

#### Delete Record
```
DELETE /api/v1/records/{id}

Response: 200 OK
{
  "message": "record deleted successfully"
}
```

#### Add Image to Record
```
POST /api/v1/records/{id}/images
Content-Type: multipart/form-data

Parameters:
- image (file, required): Image file to add

Response: 201 Created
{
  "id": 1,
  "filename": "new_image.jpg",
  "path": "./uploads/new_image.jpg",
  "vector_id": "vec_123..."
}
```

#### Delete Image
```
DELETE /api/v1/images/{image_id}

Response: 200 OK
{
  "message": "image deleted successfully"
}
```

### Search

#### Search Similar Images
```
POST /api/v1/search
Content-Type: multipart/form-data

Parameters:
- image (file, required): Image file to search for
- top_k (query, optional): Number of results (default: 10, max: 100)

Response: 200 OK
{
  "results": [
    {
      "record_id": 1,
      "record_name": "Similar Item",
      "description": "Description",
      "image_id": 1,
      "filename": "image1.jpg",
      "distance": 0.1234
    }
  ],
  "count": 5,
  "message": "Search completed successfully"
}
```

#### Find Similar Images to Existing
```
GET /api/v1/search/similar/{image_id}?top_k=10

Response: 200 OK
{
  "results": [...],
  "count": 3
}
```

#### Advanced Search
```
POST /api/v1/search/advanced
Content-Type: multipart/form-data

Parameters:
- image (file, required): Image file to search for
- q (query, optional): Text search in descriptions
- record_name (query, optional): Filter by record name
- min_distance (query, optional): Minimum similarity threshold
- max_distance (query, optional): Maximum similarity threshold
- top_k (query, optional): Number of results (default: 10)

Response: 200 OK
{
  "results": [...],
  "count": 2,
  "query": "cat",
  "filters": {
    "record_name": "pets",
    "min_distance": 0.1,
    "max_distance": 0.5
  }
}
```

### File Serving

#### Serve Uploaded Images
```
GET /uploads/{filename}

Returns: Image file (JPEG, PNG, WebP)
```

## Image Formats
Supported formats:
- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)

## Error Codes
- 400: Bad Request - Invalid parameters or missing required fields
- 404: Not Found - Resource not found
- 422: Unprocessable Entity - Validation errors
- 429: Too Many Requests - Rate limit exceeded
- 500: Internal Server Error - Server-side errors

## Example Usage

### Create Record with Images
```bash
curl -X POST http://localhost:8080/api/v1/records \
  -F "name=My Cat Photos" \
  -F "description=Collection of cat pictures" \
  -F "images=@cat1.jpg" \
  -F "images=@cat2.jpg"
```

### Search Similar Images
```bash
curl -X POST http://localhost:8080/api/v1/search?top_k=5 \
  -F "image=@query_image.jpg"
```

### List Records
```bash
curl http://localhost:8080/api/v1/records?page=1&limit=20
```

## Environment Variables
```bash
# Database
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=image_rag
MYSQL_PASSWORD=image_rag_password
MYSQL_DATABASE=image_rag

# Milvus
MILVUS_HOST=localhost
MILVUS_PORT=19530

# Doubao API
DOUBAO_API_KEY=your_api_key
DOUBAO_MODEL=doubao-embedding-vision-250615

# Server
SERVER_PORT=8080
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE_MB=10
```

## Swagger Documentation
Interactive API documentation available at: http://localhost:8080/api/v1/swagger/index.html