# Image RAG Service Code Examples

## Backend API - Go + Gin Example

### Record Creation with Multiple Images
```go
// internal/models/record.go
package models

import (
    "time"
)

type Record struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Images      []Image   `json:"images" gorm:"foreignKey:RecordID"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    RecordID uint   `json:"record_id" gorm:"not null"`
    Filename string `json:"filename" gorm:"not null"`
    Path     string `json:"path" gorm:"not null"`
    VectorID string `json:"vector_id" gorm:"not null"` // Milvus vector ID
}

// internal/api/handlers/records.go
package api

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type RecordHandler struct {
    recordService *services.RecordService
    vectorService *services.VectorService
}

func (h *RecordHandler) CreateRecord(c *gin.Context) {
    var req struct {
        Name        string                `form:"name" binding:"required"`
        Description string                `form:"description"`
        Images      []*multipart.FileHeader `form:"images" binding:"required"`
    }
    
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Create record
    record := &models.Record{
        Name:        req.Name,
        Description: req.Description,
    }
    
    if err := h.recordService.Create(record); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Process images and generate vectors
    for _, file := range req.Images {
        filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
        filepath := filepath.Join("uploads", filename)
        
        if err := c.SaveUploadedFile(file, filepath); err != nil {
            continue
        }
        
        // Generate vector using Python service
        vectorID, err := h.vectorService.GenerateVector(filepath)
        if err != nil {
            continue
        }
        
        // Save image record
        image := &models.Image{
            RecordID: record.ID,
            Filename: filename,
            Path:     filepath,
            VectorID: vectorID,
        }
        
        h.recordService.AddImage(image)
    }
    
    c.JSON(http.StatusCreated, record)
}
```

## Frontend - Vue.js Example

### Image Upload Component
```vue
<!-- components/ImageUpload.vue -->
<template>
  <div class="image-upload">
    <el-form ref="formRef" :model="form" label-width="120px">
      <el-form-item label="Record Name" prop="name" required>
        <el-input v-model="form.name" placeholder="Enter record name" />
      </el-form-item>
      
      <el-form-item label="Description" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="Enter description"
        />
      </el-form-item>
      
      <el-form-item label="Images" prop="images" required>
        <el-upload
          ref="uploadRef"
          action="#"
          list-type="picture-card"
          :auto-upload="false"
          :multiple="true"
          :on-change="handleImageChange"
          :on-remove="handleImageRemove"
        >
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" @click="submitForm" :loading="loading">
          Create Record
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { createRecord } from '@/api/records'

const formRef = ref()
const uploadRef = ref()
const loading = ref(false)

const form = reactive({
  name: '',
  description: '',
  images: []
})

const handleImageChange = (file, fileList) => {
  form.images = fileList
}

const handleImageRemove = (file, fileList) => {
  form.images = fileList
}

const submitForm = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      const formData = new FormData()
      formData.append('name', form.name)
      formData.append('description', form.description)
      
      form.images.forEach(file => {
        formData.append('images', file.raw)
      })
      
      await createRecord(formData)
      ElMessage.success('Record created successfully')
      
      // Reset form
      formRef.value.resetFields()
      uploadRef.value.clearFiles()
    } catch (error) {
      ElMessage.error('Failed to create record')
    } finally {
      loading.value = false
    }
  })
}
</script>
```

### Image Search Component
```vue
<!-- components/ImageSearch.vue -->
<template>
  <div class="image-search">
    <el-upload
      class="upload-demo"
      drag
      action="#"
      :auto-upload="false"
      :on-change="handleQueryImage"
      :show-file-list="false"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        Drop image here or <em>click to upload</em>
      </div>
    </el-upload>
    
    <div v-if="queryImage" class="query-preview">
      <img :src="queryImage" alt="Query image" />
    </div>
    
    <div v-if="results.length > 0" class="search-results">
      <h3>Search Results</h3>
      <el-row :gutter="20">
        <el-col
          v-for="result in results"
          :key="result.id"
          :span="6"
        >
          <el-card :body-style="{ padding: '0px' }">
            <img :src="getImageUrl(result.filename)" class="image" />
            <div style="padding: 14px">
              <h4>{{ result.name }}</h4>
              <p>{{ result.description }}</p>
              <p class="similarity">Similarity: {{ result.distance.toFixed(2) }}</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'
import { searchImages } from '@/api/search'

const queryImage = ref('')
const results = ref([])
const loading = ref(false)

const handleQueryImage = async (file) => {
  queryImage.value = URL.createObjectURL(file.raw)
  
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('image', file.raw)
    
    const response = await searchImages(formData)
    results.value = response.data
  } catch (error) {
    console.error('Search failed:', error)
  } finally {
    loading.value = false
  }
}

const getImageUrl = (filename) => {
  return `${import.meta.env.VITE_API_URL}/uploads/${filename}`
}
</script>
```

## Environment Configuration

### .env.example
```bash
# Backend (Go)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=image_rag

# Vector Service (Python)
DOUBAO_API_KEY=your_doubao_api_key
MILVUS_HOST=localhost
MILVUS_PORT=19530

# Frontend (Vue.js)
VITE_API_URL=http://localhost:8080
```

## API Endpoints

### Backend API Routes
```go
// internal/api/routes.go
package api

func SetupRoutes(router *gin.Engine) {
    api := router.Group("/api/v1")
    
    // Records
    api.POST("/records", recordHandler.CreateRecord)
    api.GET("/records", recordHandler.GetRecords)
    api.GET("/records/:id", recordHandler.GetRecord)
    api.PUT("/records/:id", recordHandler.UpdateRecord)
    api.DELETE("/records/:id", recordHandler.DeleteRecord)
    
    // Search
    api.POST("/search", searchHandler.SearchImages)
    api.GET("/search/similar/:id", searchHandler.FindSimilar)
    
    // Serve uploaded images
    router.Static("/uploads", "./uploads")
}
```