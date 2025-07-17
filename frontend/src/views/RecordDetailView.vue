<template>
  <div class="record-detail-view">
    <el-row :gutter="20" class="header-row">
      <el-col :span="24">
        <el-button @click="goBack" :icon="ArrowLeft" text>Back to Records</el-button>
        <h1>{{ record?.name || 'Loading...' }}</h1>
        <p class="subtitle">Detailed view of the image record</p>
      </el-col>
    </el-row>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <div v-else-if="!record" class="error-container">
      <el-empty description="Record not found">
        <el-button type="primary" @click="goBack">Go Back</el-button>
      </el-empty>
    </div>

    <div v-else>
      <el-row :gutter="20">
        <el-col :xs="24" :lg="16">
          <!-- Image Gallery -->
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>Images</span>
                <div class="card-actions">
                  <el-button 
                    type="primary" 
                    @click="showAddImageDialog = true"
                    :icon="Plus"
                    size="small"
                  >
                    Add Image
                  </el-button>
                </div>
              </div>
            </template>

            <div v-if="record.images.length === 0" class="empty-images">
              <el-empty description="No images in this record">
                <el-button type="primary" @click="showAddImageDialog = true">
                  Add First Image
                </el-button>
              </el-empty>
            </div>

            <div v-else class="image-gallery">
              <el-carousel v-if="record.images.length > 1" height="400px" arrow="always">
                <el-carousel-item v-for="image in record.images" :key="image.id">
                  <img 
                    :src="`/api/v1/images/${image.id}/preview`" 
                    :alt="record.name"
                    class="gallery-image"
                    @error="handleImageError"
                  />
                </el-carousel-item>
              </el-carousel>

              <div v-else class="single-image">
                <img 
                  :src="`/api/v1/images/${record.images[0].id}/preview`" 
                  :alt="record.name"
                  class="gallery-image"
                  @error="handleImageError"
                />
              </div>

              <div class="image-thumbnails" v-if="record.images.length > 1">
                <div 
                  v-for="image in record.images" 
                  :key="image.id"
                  class="thumbnail"
                  :class="{ active: selectedImage?.id === image.id }"
                  @click="selectImage(image)"
                >
                  <img 
                    :src="`/api/v1/images/${image.id}/preview`" 
                    :alt="record.name"
                    @error="handleImageError"
                  />
                </div>
              </div>
            </div>
          </el-card>

          <!-- Similar Images -->
          <el-card shadow="hover" class="similar-images-card" v-if="selectedImage">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Search /></el-icon>
                  Similar Images
                </span>
                <div class="card-actions">
                  <el-button 
                    @click="findSimilar(selectedImage.id)"
                    :loading="findingSimilar"
                    size="small"
                  >
                    Find Similar
                  </el-button>
                </div>
              </div>
            </template>

            <div v-if="similarImages.length === 0 && !findingSimilar" class="empty-similar">
              <el-empty description="No similar images found">
                <el-button @click="findSimilar(selectedImage.id)">
                  Search for Similar
                </el-button>
              </el-empty>
            </div>

            <div v-else-if="findingSimilar" class="loading-similar">
              <el-skeleton :rows="3" animated />
            </div>

            <div v-else class="similar-images">
              <el-row :gutter="15">
                <el-col
                  v-for="(similar, index) in similarImages"
                  :key="`${similar.record_id}-${similar.image_id}`"
                  :xs="24"
                  :sm="12"
                  :md="8"
                >
                  <el-card class="similar-card" shadow="hover">
                    <div class="similar-rank">#{{ index + 1 }}</div>
                    <img 
                      :src="`/api/v1/images/${similar.image_id}/preview`" 
                      :alt="similar.record_name"
                      class="similar-image"
                      @error="handleImageError"
                    />
                    <div class="similar-info">
                      <h5>{{ similar.record_name }}</h5>
                      <p class="similarity">{{ ((1 - similar.distance) * 100).toFixed(1) }}% match</p>
                    </div>
                    <div class="similar-actions">
                      <el-button 
                        type="primary" 
                        size="small" 
                        @click="viewRecord(similar.record_id)"
                      >
                        View
                      </el-button>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :lg="8">
          <!-- Record Details -->
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>Record Details</span>
                <div class="card-actions">
                  <el-dropdown @command="handleRecordAction">
                    <el-button text :icon="More" circle />
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="edit">
                          <el-icon><Edit /></el-icon>
                          Edit
                        </el-dropdown-item>
                        <el-dropdown-item command="delete" divided style="color: #f56c6c">
                          <el-icon><Delete /></el-icon>
                          Delete
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </template>

            <div class="record-details">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="Name">
                  {{ record.name }}
                </el-descriptions-item>
                <el-descriptions-item label="Description">
                  {{ record.description || 'No description provided' }}
                </el-descriptions-item>
                <el-descriptions-item label="Images">
                  {{ record.images.length }}
                </el-descriptions-item>
                <el-descriptions-item label="Created">
                  {{ formatDate(record.created_at) }}
                </el-descriptions-item>
                <el-descriptions-item label="Updated">
                  {{ formatDate(record.updated_at) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-card>

          <!-- Image Info -->
          <el-card shadow="hover" v-if="selectedImage">
            <template #header>
              <div class="card-header">
                <span>Image Details</span>
              </div>
            </template>

            <div class="image-details">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="Filename">
                  {{ selectedImage.filename }}
                </el-descriptions-item>
                <el-descriptions-item label="Vector ID">
                  <el-tag size="small">{{ selectedImage.vector_id }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Created">
                  {{ formatDate(selectedImage.created_at) }}
                </el-descriptions-item>
              </el-descriptions>

              <div class="image-actions" style="margin-top: 16px">
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="deleteImage(selectedImage.id)"
                  :icon="Delete"
                >
                  Delete Image
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- Add Image Dialog -->
    <el-dialog
      v-model="showAddImageDialog"
      title="Add New Image"
      width="500px"
    >
      <el-upload
        ref="uploadRef"
        drag
        :auto-upload="false"
        :multiple="false"
        :on-change="handleImageSelect"
        accept="image/*"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">
          Drop image here or <em>click to upload</em>
        </div>
      </el-upload>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddImageDialog = false">Cancel</el-button>
          <el-button type="primary" @click="uploadImage" :loading="uploading">
            Upload
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { recordService, searchService } from '@/services/api'
import type { Record, Image } from '@/types'
import { ArrowLeft, Plus, More, Edit, Delete, Search, Upload } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// Data
const record = ref<Record | null>(null)
const selectedImage = ref<Image | null>(null)
const loading = ref(false)
const uploading = ref(false)
const findingSimilar = ref(false)
const similarImages = ref<SearchResult[]>([])

// Dialogs
const showAddImageDialog = ref(false)
const uploadRef = ref()
const selectedFile = ref<File | null>(null)

// Methods
const loadRecord = async () => {
  const recordId = Number(route.params.id)
  if (!recordId) {
    router.push('/records')
    return
  }

  loading.value = true
  try {
    record.value = await recordService.getRecord(recordId)
    if (record.value.images.length > 0) {
      selectedImage.value = record.value.images[0]
    }
  } catch (error) {
    console.error('Failed to load record:', error)
    ElMessage.error('Failed to load record')
  } finally {
    loading.value = false
  }
}

const selectImage = (image: Image) => {
  selectedImage.value = image
  clearSimilarResults()
}

const handleImageSelect = (file: any) => {
  selectedFile.value = file.raw
}

const uploadImage = async () => {
  if (!selectedFile.value || !record.value) return

  uploading.value = true
  try {
    await recordService.addImageToRecord(record.value.id, selectedFile.value)
    ElMessage.success('Image uploaded successfully')
    showAddImageDialog.value = false
    selectedFile.value = null
    loadRecord()
  } catch (error) {
    console.error('Failed to upload image:', error)
    ElMessage.error('Failed to upload image')
  } finally {
    uploading.value = false
  }
}

const deleteImage = async (imageId: number) => {
  if (!record.value) return

  try {
    await ElMessageBox.confirm(
      'Are you sure you want to delete this image? This action cannot be undone.',
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )

    await recordService.deleteImage(imageId)
    ElMessage.success('Image deleted successfully')
    loadRecord()
    clearSimilarResults()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete image:', error)
      ElMessage.error('Failed to delete image')
    }
  }
}

const findSimilar = async (imageId: number) => {
  if (findingSimilar.value) return

  findingSimilar.value = true
  try {
    const results = await searchService.findSimilar(imageId, 12)
    similarImages.value = results
    if (results.length === 0) {
      ElMessage.info('No similar images found')
    }
  } catch (error) {
    console.error('Failed to find similar images:', error)
    ElMessage.error('Failed to find similar images')
  } finally {
    findingSimilar.value = false
  }
}

const clearSimilarResults = () => {
  similarImages.value = []
}

const handleRecordAction = async (command: string) => {
  if (!record.value) return

  switch (command) {
    case 'edit':
      ElMessage.info('Edit functionality coming soon')
      break
    case 'delete':
      await deleteRecord()
      break
  }
}

const deleteRecord = async () => {
  if (!record.value) return

  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete "${record.value.name}"? This action cannot be undone.`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )

    await recordService.deleteRecord(record.value.id)
    ElMessage.success('Record deleted successfully')
    router.push('/records')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete record:', error)
      ElMessage.error('Failed to delete record')
    }
  }
}

const goBack = () => {
  router.push('/records')
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/placeholder-image.jpg'
}

// Initialize
onMounted(() => {
  loadRecord()
})
</script>

<style scoped>
.record-detail-view {
  max-width: 1400px;
  margin: 0 auto;
}

.header-row {
  margin-bottom: 30px;
}

.header-row h1 {
  font-size: 32px;
  color: #303133;
  margin: 10px 0;
}

.subtitle {
  font-size: 16px;
  color: #909399;
  margin: 0;
}

.loading-container,
.error-container {
  text-align: center;
  padding: 60px 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: bold;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.image-gallery {
  margin-bottom: 20px;
}

.gallery-image {
  width: 100%;
  height: 400px;
  object-fit: contain;
  border-radius: 4px;
}

.single-image {
  text-align: center;
}

.single-image .gallery-image {
  max-width: 100%;
  height: auto;
}

.image-thumbnails {
  display: flex;
  gap: 10px;
  margin-top: 15px;
  overflow-x: auto;
  padding: 10px 0;
}

.thumbnail {
  flex-shrink: 0;
  width: 80px;
  height: 80px;
  border: 2px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  overflow: hidden;
}

.thumbnail.active {
  border-color: #409eff;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.empty-images {
  text-align: center;
  padding: 40px 0;
}

.similar-images-card {
  margin-top: 20px;
}

.empty-similar {
  text-align: center;
  padding: 40px 0;
}

.loading-similar {
  padding: 20px;
}

.similar-images .el-row {
  margin: 0 -7.5px;
}

.similar-images .el-col {
  padding: 0 7.5px;
}

.similar-card {
  margin-bottom: 15px;
  text-align: center;
}

.similar-rank {
  position: absolute;
  top: 8px;
  left: 8px;
  background-color: #409eff;
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.similar-image {
  width: 100%;
  height: 150px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 8px;
}

.similar-info h5 {
  margin: 0 0 4px 0;
  font-size: 14px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.similarity {
  color: #67c23a;
  font-size: 12px;
  margin: 0;
}

.similar-actions {
  margin-top: 8px;
}

.record-details,
.image-details {
  padding: 10px 0;
}

@media (max-width: 768px) {
  .header-row h1 {
    font-size: 24px;
  }
  
  .gallery-image {
    height: 300px;
  }
  
  .similar-card {
    margin-bottom: 10px;
  }
}
</style>