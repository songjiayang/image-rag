## FEATURE:

Image RAG (Retrieval Augmented Generation) service for vector-based image search with the following capabilities:

1. **Core Image Search**: Users can upload an image and retrieve similar images with their metadata (id, name, description)
2. **RAG Service Management Interface**: Web-based admin interface for:
   - Database creation and management
   - Adding records with name, description, and multiple images
   - One-to-many relationship between records and images
3. **Image Vectorization Pipeline**: 
   - Automatic image vectorization upon record creation
   - Vector storage in Milvus vector database
   - Integration with Doubao doubao-embedding-vision-250615 API

## EXAMPLES:

- **Complete implementation example**: Available in `examples/image-rag-example.md`
- **Backend API**: Go + Gin REST API with endpoints for record CRUD operations and image search
- **Vector Integration**: Go service directly integrates with Doubao API for image vectorization and Milvus client for vector storage
- **Frontend**: Vue.js SPA with components for image upload, record management, and visual search
- **Database Models**: Complete schema examples for MySQL records/images tables and Milvus collections

## DOCUMENTATION:

- **Doubao API Documentation**: https://www.volcengine.com/docs/82379/1523520
  - doubao-embedding-vision-250615 model specifications
  - API authentication and request/response formats
  - Rate limiting and usage guidelines
- **Milvus Documentation**: https://milvus.io/docs
  - Vector collection setup and indexing
  - Search parameters and similarity metrics
  - Python client usage
- **Go Gin Framework**: https://gin-gonic.com/docs/
- **Vue.js 3 Documentation**: https://vuejs.org/guide/
- **GORM Documentation**: https://gorm.io/docs/
- **Element Plus**: https://element-plus.org/en-US/
- **Image Processing Libraries**:
  - Go standard library for image handling (image, image/jpeg, image/png)
  - Base64 encoding for Doubao API requests
  - Milvus Go client for vector operations

## OTHER CONSIDERATIONS:

### Security & Performance:
- **Image file validation**: Ensure uploaded files are valid images (jpeg, png, webp)
- **File size limits**: Implement upload size restrictions (e.g., 10MB per image)
- **Rate limiting**: Add API rate limiting to prevent abuse
- **CORS configuration**: Configure properly for frontend-backend communication
- **API key security**: Store Doubao API key securely (environment variables)

### Scalability & Architecture:
- **Asynchronous processing**: Use goroutines for image vectorization to avoid blocking
- **Batch processing**: Support batch image uploads for better performance
- **Caching**: Implement Redis caching for frequent searches
- **File storage**: Consider cloud storage (AWS S3) for production
- **Vector indexing**: Use appropriate Milvus index type (IVF_FLAT or HNSW)

### Error Handling:
- **Graceful degradation**: Handle Doubao API failures with fallback
- **User feedback**: Clear error messages for upload failures
- **Retry logic**: Implement retry mechanism for API calls
- **Logging**: Comprehensive logging for debugging

### Data Management:
- **Data consistency**: Ensure MySQL and Milvus stay synchronized
- **Cleanup**: Implement orphaned vector cleanup
- **Migration**: Database migration scripts for schema changes
- **Backup**: Regular backup strategy for both MySQL and Milvus

### Common Gotchas:
- **Image encoding**: Ensure proper base64 encoding for Doubao API
- **Vector dimensions**: Verify vector dimensions match Milvus collection schema
- **Memory usage**: Monitor Go service memory usage with large images
- **File cleanup**: Remove temporary files after vectorization
- **Cross-origin issues**: Configure proper CORS headers for file uploads
- **Timeout handling**: Set appropriate timeouts for Doubao API calls
