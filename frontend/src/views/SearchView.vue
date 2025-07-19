<template>
  <div class="search-view">
    <el-row :gutter="20" class="header-row">
      <el-col :span="24">
        <h1>AI-Powered Image Search</h1>
        <p class="subtitle">Find similar images using vector similarity search</p>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Search Images</span>
            </div>
          </template>

          <div class="search-section">
            <el-tabs v-model="activeTab" class="search-tabs">
              <el-tab-pane label="Upload Image" name="upload">
                <div class="upload-section">
                  <el-upload
                    ref="uploadRef"
                    class="upload-area"
                    drag
                    :auto-upload="false"
                    :show-file-list="false"
                    :on-change="handleImageUpload"
                    accept="image/*"
                  >
                    <div v-if="!searchImage">
                      <el-icon class="upload-icon"><Upload /></el-icon>
                      <div class="upload-text">
                        Drop image here or <em>click to upload</em>
                      </div>
                      <div class="upload-tip">
                        Supports JPG, PNG, WebP formats
                      </div>
                    </div>
                    <div v-else class="preview-area">
                      <img :src="searchImage" alt="Search image" class="preview-image" />
                      <div class="preview-actions">
                        <el-button type="danger" @click="clearImage" size="small">
                          <el-icon><Delete /></el-icon>
                          Clear
                        </el-button>
                      </div>
                    </div>
                  </el-upload>
                </div>
              </el-tab-pane>

              <el-tab-pane label="From URL" name="url">
                <div class="url-section">
                  <el-input
                    v-model="imageUrl"
                    placeholder="Enter image URL"
                    clearable
                    @change="handleUrlChange"
                  >
                    <template #prefix>
                      <el-icon><Link /></el-icon>
                    </template>
                  </el-input>
                  <div v-if="searchImage" class="url-preview">
                    <img :src="searchImage" alt="URL image" class="preview-image" />
                    <el-button type="danger" @click="clearImage" size="small">
                      <el-icon><Delete /></el-icon>
                      Clear
                    </el-button>
                  </div>
                </div>
              </el-tab-pane>

              <el-tab-pane label="From Records" name="records">
                <div class="record-section">
                  <p>Select an image from your existing records</p>
                  <el-select
                    v-model="selectedRecordId"
                    placeholder="Select a record"
                    filterable
                    clearable
                    @change="handleRecordSelect"
                  >
                    <el-option
                      v-for="record in records"
                      :key="record.id"
                      :label="record.name"
                      :value="record.id"
                    />
                  </el-select>
                  <div v-if="selectedImage" class="record-preview">
                    <img :src="selectedImage" alt="Record image" class="preview-image" />
                    <el-button type="danger" @click="clearImage" size="small">
                      <el-icon><Delete /></el-icon>
                      Clear
                    </el-button>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>

            <div class="search-options">
              <el-form :model="searchOptions" label-width="120px">
                <el-form-item label="Results to show">
                  <el-slider
                    v-model="searchOptions.topK"
                    :min="5"
                    :max="50"
                    :step="5"
                    show-stops
                    show-input
                  />
                </el-form-item>

                <el-form-item label="Advanced Options">
                  <el-checkbox v-model="searchOptions.useAdvanced">Enable advanced search</el-checkbox>
                </el-form-item>

                <div v-if="searchOptions.useAdvanced" class="advanced-options">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-form-item label="Text Query">
                        <el-input
                          v-model="searchOptions.textQuery"
                          placeholder="Optional text search"
                          clearable
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="12">
                      <el-form-item label="Record Name">
                        <el-input
                          v-model="searchOptions.recordName"
                          placeholder="Filter by record name"
                          clearable
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </div>
              </el-form>
            </div>

            <div class="search-actions">
              <el-button
                type="primary"
                size="large"
                @click="performSearch"
                :disabled="!canSearch"
                :loading="searching"
              >
                <el-icon><Search /></el-icon>
                Search Images
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Search Results -->
    <el-row v-if="searchResults.length > 0" :gutter="20" class="results-row">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="results-header">
              <span>
                <el-icon><List /></el-icon>
                Search Results ({{ searchResults.length }} found)
              </span>
              <div class="results-actions">
                <el-button @click="clearResults" size="small">
                  <el-icon><Close /></el-icon>
                  Clear Results
                </el-button>
              </div>
            </div>
          </template>

          <div class="results-container">
            <el-row :gutter="20">
              <el-col
                v-for="(result, index) in searchResults"
                :key="`${result.record_id}-${result.image_id}`"
                :xs="24"
                :sm="12"
                :md="8"
                :lg="6"
              >
                <el-card class="result-card" shadow="hover">
                  <div class="result-rank">
                    <span class="rank-number">#{{ index + 1 }}</span>
                    <span class="similarity">{{ (1 - result.distance).toFixed(3) }}</span>
                  </div>
                  
                  <div class="result-image">
                    <img
                      :src="`/api/v1/images/${result.image_id}/preview`"
                      :alt="result.record_name"
                    />
                  </div>

                  <div class="result-info">
                    <h4>{{ result.record_name }}</h4>
                    <p class="description">{{ result.description || 'No description' }}</p>
                    <p class="distance">Distance: {{ result.distance.toFixed(4) }}</p>
                  </div>

                  <div class="result-actions">
                    <el-button
                      type="primary"
                      size="small"
                      @click="viewRecord(result.record_id)"
                    >
                      <el-icon><View /></el-icon>
                      View Details
                    </el-button>
                    <el-button
                      size="small"
                      @click="findSimilar(result.image_id)"
                    >
                      <el-icon><Search /></el-icon>
                      Find Similar
                    </el-button>
                  </div>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { searchService, recordService } from '@/services/api'
import type { SearchResult, Record, SearchResponse } from '@/types'

const router = useRouter()

// Search state
const activeTab = ref<'upload' | 'url' | 'records'>('upload')
const searchImage = ref<string | null>(null)
const imageUrl = ref('')
const selectedRecordId = ref<number | null>(null)
const selectedImage = ref<string | null>(null)
const searching = ref(false)

// Records for selection
const records = ref<Record[]>([])

// Search options
const searchOptions = ref({
  topK: 10,
  useAdvanced: false,
  textQuery: '',
  recordName: ''
})

// Results
const searchResults = ref<SearchResult[]>([])

// Computed
const canSearch = computed(() => {
  return searchImage.value !== null || selectedImage.value !== null
})

// Methods
const loadRecords = async () => {
  try {
    const response = await recordService.getRecords(1, 100)
    records.value = response.data
  } catch (error) {
    console.error('Failed to load records:', error)
    ElMessage.error('Failed to load records for selection')
  }
}

const handleImageUpload = (file: any) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    searchImage.value = e.target?.result as string
  }
  reader.readAsDataURL(file.raw)
  activeTab.value = 'upload'
}

const handleUrlChange = () => {
  if (imageUrl.value) {
    searchImage.value = imageUrl.value
    activeTab.value = 'url'
  }
}

const handleRecordSelect = async (recordId: number) => {
  const record = records.value.find(r => r.id === recordId)
  if (record && record.images.length > 0) {
    selectedImage.value = `/api/v1/images/${record.images[0].id}/preview`
    searchImage.value = selectedImage.value
    activeTab.value = 'records'
  }
}

const clearImage = () => {
  searchImage.value = null
  imageUrl.value = ''
  selectedRecordId.value = null
  selectedImage.value = null
}

const performSearch = async () => {
  if (!searchImage.value) return

  searching.value = true
  try {
    let response: SearchResponse
    
    if (searchOptions.value.useAdvanced) {
      // Convert data URL or URL to File for advanced search
      const file = await convertToFile(searchImage.value)
      response = await searchService.advancedSearch(file, {
        q: searchOptions.value.textQuery,
        record_name: searchOptions.value.recordName,
        top_k: searchOptions.value.topK
      })
    } else {
      const file = await convertToFile(searchImage.value)
      response = await searchService.searchImages(file, searchOptions.value.topK)
    }

    searchResults.value = response.results
    
    if (response.results.length === 0) {
      ElMessage.info('No similar images found')
    }
  } catch (error) {
    console.error('Search failed:', error)
    ElMessage.error('Search failed. Please try again.')
  } finally {
    searching.value = false
  }
}

const convertToFile = async (imageSrc: string): Promise<File> => {
  if (imageSrc.startsWith('data:')) {
    // Data URL
    const response = await fetch(imageSrc)
    const blob = await response.blob()
    return new File([blob], 'search-image.jpg', { type: blob.type })
  } else {
    // URL
    const response = await fetch(imageSrc)
    const blob = await response.blob()
    return new File([blob], 'search-image.jpg', { type: blob.type })
  }
}

const clearResults = () => {
  searchResults.value = []
}

const viewRecord = (recordId: number) => {
  router.push(`/records/${recordId}`)
}

const findSimilar = async (imageId: number) => {
  if (searching.value) return
  
  searching.value = true
  try {
    const response = await searchService.findSimilar(imageId, searchOptions.value.topK)
    searchResults.value = response.results
  } catch (error) {
    console.error('Find similar failed:', error)
    ElMessage.error('Failed to find similar images')
  } finally {
    searching.value = false
  }
}

// Initialize
onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.search-view {
  max-width: 1200px;
  margin: 0 auto;
}

.header-row {
  margin-bottom: 30px;
  text-align: center;
}

.header-row h1 {
  font-size: 32px;
  color: #303133;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 16px;
  color: #909399;
  margin: 0;
}

.search-section {
  padding: 20px 0;
}

.search-tabs {
  margin-bottom: 30px;
}

.upload-section {
  text-align: center;
}

.upload-area {
  width: 100%;
}

.upload-area :deep(.el-upload-dragger) {
  padding: 60px 20px;
  border: 2px dashed #dcdfe6;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 20px;
}

.upload-text {
  color: #606266;
  font-size: 16px;
  margin-bottom: 8px;
}

.upload-tip {
  color: #909399;
  font-size: 14px;
}

.preview-area,
.url-preview,
.record-preview {
  text-align: center;
}

.preview-image {
  max-width: 300px;
  max-height: 300px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.preview-actions {
  margin-top: 16px;
}

.url-section,
.record-section {
  text-align: center;
}

.url-section .el-input,
.record-section .el-select {
  max-width: 400px;
  margin-bottom: 20px;
}

.search-options {
  margin: 30px 0;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.advanced-options {
  margin-top: 20px;
  padding: 20px;
  background-color: white;
  border-radius: 4px;
}

.search-actions {
  text-align: center;
  margin-top: 30px;
}

.results-row {
  margin-top: 30px;
}

.results-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: bold;
}

.results-actions {
  display: flex;
  gap: 8px;
}

.results-container {
  padding: 20px 0;
}

.result-card {
  margin-bottom: 20px;
  height: 100%;
}

.result-rank {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.rank-number {
  background-color: #409eff;
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.similarity {
  color: #67c23a;
  font-weight: bold;
  font-size: 14px;
}

.result-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 12px;
}

.result-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.result-info h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-info .description {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.result-info .distance {
  font-size: 12px;
  color: #909399;
  margin: 0;
}

.result-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

@media (max-width: 768px) {
  .header-row h1 {
    font-size: 24px;
  }
  
  .subtitle {
    font-size: 14px;
  }
  
  .preview-image {
    max-width: 200px;
    max-height: 200px;
  }
  
  .result-actions {
    flex-direction: column;
  }
}
</style>