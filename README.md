# Image RAG Service

A production-ready image retrieval-augmented generation (RAG) service that enables vector-based image search using Doubao API for image embeddings and Milvus for vector storage.

## 🚀 Features

- **Image Upload & Processing**: Upload images with automatic vectorization
- **Vector Search**: Semantic search across uploaded images
- **RESTful API**: Complete Go + Gin backend with OpenAPI documentation
- **Web Interface**: Vue.js frontend for easy management and search
- **Production Deployment**: Docker containers with CI/CD pipeline
- **Health Monitoring**: Comprehensive health checks for all services

## 🏗️ Architecture

### Tech Stack
- **Backend**: Go with Gin framework
- **Frontend**: Vue.js 3 with TypeScript and Composition API
- **Vector Database**: Milvus for image embeddings
- **Metadata Storage**: MySQL for image metadata
- **Image Storage**: Local filesystem
- **Vectorization**: Doubao doubao-embedding-vision-250615 API

### Service Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Vue.js        │────│   Go API        │────│   Doubao API    │
│   Frontend      │    │   Backend       │    │   (Embeddings)  │
│   Port: 3000    │    │   Port: 8080    │    │   External      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
              ┌─────────────────┬─────────────────┐
              │                 │                 │
    ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
    │   MySQL         │ │   Milvus        │ │   File System   │
    │   Port: 3306    │ │   Port: 19530   │ │   Images        │
    │   Metadata      │ │   Vectors       │ │   Storage       │
    └─────────────────┘ └─────────────────┘ └─────────────────┘
```

## 🛠️ Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

### Local Development

1. **Clone and setup**:
```bash
git clone <repository-url>
cd image-rag
cp .env.example .env
```

2. **Configure environment**:
Edit `.env` with your API keys and database configurations.

3. **Start services**:
```bash
docker-compose up -d
```

4. **Access services**:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/index.html

### Manual Development Setup

1. **Backend**:
```bash
cd backend
go mod tidy
go run cmd/main.go
```

2. **Frontend**:
```bash
cd frontend
npm install
npm run dev
```

## 📋 API Endpoints

### Core Endpoints
- `POST /api/images/upload` - Upload images
- `POST /api/search` - Search similar images
- `GET /api/images` - List all images
- `GET /api/images/:id` - Get image details
- `DELETE /api/images/:id` - Delete image

### Health & Monitoring
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /live` - Liveness check

## 🐳 Production Deployment

### Using GitHub Actions (Recommended)
1. Push to `main` branch → Automatic staging deployment
2. Create release → Automatic production deployment

### Manual Deployment
```bash
# Build and push images
docker build -t image-rag-backend ./backend
docker build -t image-rag-frontend ./frontend

# Deploy to Kubernetes
kubectl apply -f k8s/production/
```

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=image_rag

# Milvus
MILVUS_HOST=localhost
MILVUS_PORT=19530

# Doubao API
DOUBAO_API_KEY=your-api-key
DOUBAO_ENDPOINT=https://ark.cn-beijing.volces.com/api/v3

# Service URLs
BACKEND_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
```

## 📊 Monitoring

### Health Checks
All services include comprehensive health monitoring:
- Database connectivity
- Milvus connection
- API endpoint availability
- Resource usage

### Logs
```bash
# View service logs
docker-compose logs -f [service-name]

# Backend logs
docker-compose logs -f backend

# Frontend logs
docker-compose logs -f frontend
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 📁 Project Structure

```
image-rag/
├── backend/
│   ├── internal/
│   │   ├── api/           # API handlers and routes
│   │   ├── models/        # Database models
│   │   ├── services/      # Business logic
│   │   └── config/        # Configuration
│   ├── cmd/              # Application entry points
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/    # Vue components
│   │   ├── views/         # Page views
│   │   ├── services/      # API services
│   │   └── types/         # TypeScript types
│   └── Dockerfile
├── k8s/                  # Kubernetes manifests
├── docker-compose.yml    # Local development
└── .github/workflows/    # CI/CD pipelines
```

## 🔍 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Check MySQL is running: `docker-compose ps`
   - Verify credentials in `.env`

2. **Milvus connection failed**
   - Ensure Milvus container is healthy
   - Check port 19530 is accessible

3. **Doubao API errors**
   - Verify API key is set correctly
   - Check network connectivity to Volces API

### Debug Mode
```bash
# Enable debug logging
docker-compose -f docker-compose.yml -f docker-compose.debug.yml up
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit pull request

## 📄 License

MIT License - see LICENSE file for details.