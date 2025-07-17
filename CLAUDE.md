### 🔄 Project Awareness & Context
- **Always read `PLANNING.md`** at the start of a new conversation to understand the project's architecture, goals, style, and constraints.
- **Check `TASK.md`** before starting a new task. If the task isn't listed, add it with a brief description and today's date.
- **Use consistent naming conventions, file structure, and architecture patterns** as described in `PLANNING.md`.
- **Project Goal**: Image RAG service for vector-based image search returning image info (id, name, description)

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
- **Mark completed tasks in `TASK.md`** immediately after finishing them.
- Add new sub-tasks or TODOs discovered during development to `TASK.md` under a “Discovered During Work” section.

### 📎 Style & Conventions
- **Backend**: Go with Gin framework, follow Go conventions and use `gofmt` for formatting
- **Frontend**: Vue.js 3 with Composition API, TypeScript preferred
- **Vector Service**: Integrated into Go backend for direct Doubao API calls
- **Database**: 
  - MySQL with GORM for Go backend
  - Milvus for vector storage with Python client
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

### 🧠 AI Behavior Rules
- **Never assume missing context. Ask questions if uncertain.**
- **Never hallucinate libraries or functions** – only use known, verified packages for Go, Python, and Node.js
- **Always confirm file paths and module names** exist before referencing them in code or tests
- **Never delete or overwrite existing code** unless explicitly instructed to or if part of a task from `TASK.md`
- **Service Architecture**:
  - **Backend API**: Go + Gin (port 8080)
  - **Frontend**: Vue.js (port 3000)
  - **Databases**: MySQL (3306) + Milvus (19530)
- **Image Processing**: Use Doubao doubao-embedding-vision-250615 for vectorization
- **File Structure**: Follow microservices pattern with clear separation of concerns