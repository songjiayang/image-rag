### 🔄 Project Awareness & Context
- **Always read `./PRPS/image-rag-service.md`** at the start of a new conversation to understand the project's architecture, goals, style, and constraints.
- **Check `TASK.md`** before starting a new task. If the task isn't listed, add it with a brief description and today's date.
- **Use consistent naming conventions, file structure, and architecture patterns** as described in `./PRPS/image-rag-service.md`.
- **Project Goal**: Image RAG service for vector-based image search returning image info (id, name, description)

### 📊 Project Status Summary
**Current Status**: ✅ **Core Features Complete** - Image RAG service is fully functional with all major components integrated and tested.

**Key Achievements**:
- ✅ Direct Doubao API integration in Go backend (no Python dependency)
- ✅ Complete Vue.js frontend with TypeScript
- ✅ MySQL + Milvus database integration
- ✅ Docker containerization with compose setup
- ✅ REST API with OpenAPI documentation
- ✅ Production-ready deployment configuration

**Services Architecture**:
- **Backend**: Go + Gin REST API (port 8080)
- **Frontend**: Vue.js 3 SPA (port 3000)
- **Database**: MySQL 8.0 + Milvus 2.3
- **Storage**: Local filesystem for images
- **Vectorization**: Doubao doubao-embedding-vision-250615 API

### 🧱 Code Structure & Modularity
- **Never create a file longer than 500 lines of code.** If a file approaches this limit, refactor by splitting it into modules or helper files.
- **Organize code into clearly separated modules**, grouped by feature or responsibility.
  - **Frontend**: Vue.js SPA for RAG service management
  - **Backend**: Go + Gin REST API with separate packages for handlers, models, services, middleware
  - **Vector Service**: Integrated into Go backend for direct Doubao API calls
- **Architecture**: 
  - Database: MySQL for metadata storage (records with name, description)
  - Vector DB: Milvus for image embeddings
  - Image Storage: Local filesystem for uploaded images
  - Vectorization: Doubao doubao-embedding-vision-250615 API
- **Use clear, consistent imports** (prefer relative imports within packages).
- **Use environment variables** for configuration (database URLs, API keys, etc.)

### ✅ Task Completion
- **Always update `TASK.md` at the start of new requirements** before beginning work.
- **Mark completed tasks in `TASK.md`** immediately after finishing them.
- Add new sub-tasks or TODOs discovered during development to `TASK.md` under a "Discovered During Work" section.

### 📎 Style & Conventions
- **Backend**: Go with Gin framework, follow Go conventions and use `gofmt` for formatting
- **Frontend**: Vue.js 3 with Composition API, TypeScript preferred
- **Vector Service**: Integrated into Go backend for direct Doubao API calls
- **Database**: 
  - MySQL with GORM for Go backend
  - Milvus for vector storage
- **API Documentation**: OpenAPI/Swagger for Go API endpoints
- **Code Style**: Follow Go conventions with gofmt for formatting
- **Configuration**: Use `.env` files for all external service configurations

### 📚 Documentation & Explainability
- **Update `README.md`** when new features are added, dependencies change, or setup steps are modified.
- **Comment non-obvious code** and ensure everything is understandable to a mid-level developer.
- When writing complex logic, **add an inline `# Reason:` comment** explaining the why, not just the what.
- **External Services**:
  - **Doubao API**: https://www.volcengine.com/docs/82379/1523520 for doubao-embedding-vision-250615
  - **Milvus**: Vector database for image embeddings
  - **MySQL**: Relational database for metadata storage

### 💡 Key Memories
- to memorize