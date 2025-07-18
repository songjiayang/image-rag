<template>
  <div class="records-view">
    <el-row :gutter="20" class="header-row">
      <el-col :span="24">
        <h1>Image Records</h1>
        <p class="subtitle">Manage your image collections and metadata</p>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="actions-row">
      <el-col :span="24">
        <el-button type="primary" @click="showCreateDialog = true" size="large">
          <el-icon><Plus /></el-icon>
          Add New Record
        </el-button>
        <el-button @click="loadRecords" :loading="loading">
          <el-icon><Refresh /></el-icon>
          Refresh
        </el-button>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>All Records</span>
              <div class="card-actions">
                <el-input
                  v-model="searchQuery"
                  placeholder="Search records..."
                  clearable
                  style="width: 200px; margin-right: 10px"
                  @clear="loadRecords"
                >
                  <template ><el-icon><Search /></el-icon></template>
                </el-input>
              </div>
            </div>
          </template>

          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="5" animated />
          </div>

          <div v-else-if="records.length === 0" class="empty-container">
            <el-empty description="No records found">
              <template #default>
                <el-button type="primary" @click="showCreateDialog = true">
                  Add Your First Record
                </el-button>
              </template>
            </el-empty>
          </div>

          <div v-else class="records-container">
            <el-row :gutter="20">
              <el-col
                v-for="record in records"
                :key="record.id"
                :xs="24"
                :sm="12"
                :md="8"
                :lg="6"
              >
                <el-card class="record-card" shadow="hover">
                  <template #header>
                    <div class="record-header">
                      <h4>{{ record.name }}</h4>
                      <el-dropdown @command="handleRecordAction($event, record)">
                        <el-button text circle>
                          <el-icon><More /></el-icon>
                        </el-button>
                        <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item command="view">
                              <el-icon><View /></el-icon> View Details
                            </el-dropdown-item>
                            <el-dropdown-item command="edit">
                              <el-icon><Edit /></el-icon> Edit
                            </el-dropdown-item>
                            <el-dropdown-item command="delete" divided style="color: #f56c6c">
                              <el-icon><Delete /></el-icon> Delete
                            </el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                    </div>
                  </template>

                  <div class="record-content">
                    <div class="record-images">
                      <div v-if="record.images.length === 0" class="no-images">
                        <el-icon><Picture /></el-icon>
                        <span>No images</span>
                      </div>
                      
                      <div v-else class="image-gallery">
                        <div class="main-image">
                          <img
                            v-if="record.images[0]"
                            :src="`/api/v1/images/${record.images[0].id}/preview`"
                            :alt="record.name"
                            @error="handleImageError"
                          />
                        </div>
                        <div v-if="record.images.length > 1" class="image-count">
                          +{{ record.images.length - 1 }} more
                        </div>
                      </div>
                    </div>

                    <div class="record-info">
                      <p class="description">{{ record.description || 'No description provided' }}</p>
                      <p class="date">Created: {{ formatDate(record.created_at) }}</p>
                      <p class="images-count">{{ record.images.length }} image{{ record.images.length !== 1 ? 's' : '' }}</p>
                    </div>
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <div class="pagination-container">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="totalRecords"
                :page-sizes="[12, 24, 48, 96]"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Create Record Dialog -->
    <el-dialog
      v-model="showCreateDialog"
      title="Create New Record"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="createForm.name" placeholder="Enter record name" />
        </el-form-item>

        <el-form-item label="Description" prop="description">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="Enter record description (optional)"
          />
        </el-form-item>

        <el-form-item label="Images" prop="images">
          <el-upload
            ref="uploadRef"
            class="upload-demo"
            drag
            multiple
            :auto-upload="false"
            :file-list="fileList"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            accept="image/*"
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">
              Drop image files here or <em>click to upload</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                jpg/png files with a size less than 10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">Cancel</el-button>
          <el-button type="primary" @click="createRecord" :loading="saving"
            >Create Record</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { recordService } from '@/services/api'
import type { Record, CreateRecordRequest } from '@/types'

const router = useRouter()

// Data
const records = ref<Record[]>([])
const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const totalRecords = ref(0)

// Dialog
const showCreateDialog = ref(false)
const createFormRef = ref()
const uploadRef = ref()

// Form
const createForm = ref<CreateRecordRequest & { images: File[] }>({
  name: '',
  description: '',
  images: []
})

const fileList = ref([])

const createRules = {
  name: [
    { required: true, message: 'Please enter record name', trigger: 'blur' },
    { min: 2, max: 100, message: 'Length should be 2 to 100', trigger: 'blur' }
  ],
  images: [
    { required: true, message: 'Please upload at least one image', trigger: 'change' }
  ]
}

// Methods
const loadRecords = async () => {
  loading.value = true
  try {
    const response = await recordService.getRecords(currentPage.value, pageSize.value)
    records.value = response.data
    totalRecords.value = response.total
  } catch (error) {
    console.error('Failed to load records:', error)
    ElMessage.error('Failed to load records')
  } finally {
    loading.value = false
  }
}

const handleFileChange = (_file: any, fileList: any[]) => {
  createForm.value.images = fileList.map(f => f.raw)
}

const handleFileRemove = (_file: any, fileList: any[]) => {
  createForm.value.images = fileList.map(f => f.raw)
}

const createRecord = async () => {
  if (!createFormRef.value) return
  
  try {
    await createFormRef.value.validate()
    
    if (createForm.value.images.length === 0) {
      ElMessage.warning('Please upload at least one image')
      return
    }

    saving.value = true
    const newRecord = await recordService.createRecord(
      {
        name: createForm.value.name,
        description: createForm.value.description
      },
      createForm.value.images
    )

    ElMessage.success('Record created successfully')
    showCreateDialog.value = false
    resetCreateForm()
    loadRecords()
    
    // Navigate to the new record
    router.push(`/records/${newRecord.id}`)
  } catch (error: any) {
    console.error('Failed to create record:', error)
    ElMessage.error(error.response?.data?.message || 'Failed to create record')
  } finally {
    saving.value = false
  }
}

const resetCreateForm = () => {
  createForm.value = {
    name: '',
    description: '',
    images: []
  }
  fileList.value = []
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

const handleRecordAction = (action: string, record: Record) => {
  switch (action) {
    case 'view':
      router.push(`/records/${record.id}`)
      break
    case 'edit':
      // TODO: Implement edit functionality
      ElMessage.info('Edit functionality coming soon')
      break
    case 'delete':
      handleDeleteRecord(record)
      break
  }
}

const handleDeleteRecord = async (record: Record) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete "${record.name}"? This action cannot be undone.`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )

    await recordService.deleteRecord(record.id)
    ElMessage.success('Record deleted successfully')
    loadRecords()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete record:', error)
      ElMessage.error('Failed to delete record')
    }
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  loadRecords()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadRecords()
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/placeholder-image.jpg'
}

watch(searchQuery, (newQuery) => {
  // TODO: Implement search functionality
  if (newQuery) {
    // Implement search
  } else {
    loadRecords()
  }
})

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.records-view {
  max-width: 1400px;
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

.actions-row {
  margin-bottom: 30px;
  text-align: right;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: bold;
}

.card-actions {
  display: flex;
  align-items: center;
}

.loading-container {
  padding: 40px;
  text-align: center;
}

.empty-container {
  padding: 40px;
  text-align: center;
}

.records-container {
  padding: 20px 0;
}

.record-card {
  margin-bottom: 20px;
  height: 100%;
}

.record-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.record-header h4 {
  margin: 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.record-content {
  padding: 10px 0;
}

.record-images {
  margin-bottom: 12px;
}

.no-images {
  height: 150px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  color: #909399;
  border-radius: 4px;
}

.no-images .el-icon {
  font-size: 48px;
  margin-bottom: 8px;
}

.image-gallery {
  position: relative;
}

.main-image {
  width: 100%;
  height: 150px;
  overflow: hidden;
  border-radius: 4px;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-count {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.record-info .description {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.record-info .date,
.record-info .images-count {
  font-size: 12px;
  color: #909399;
  margin: 4px 0;
}

.pagination-container {
  margin-top: 30px;
  text-align: center;
}

.upload-demo {
  width: 100%;
}

:deep(.el-upload-dragger) {
  padding: 40px;
}

@media (max-width: 768px) {
  .header-row h1 {
    font-size: 24px;
  }
  
  .subtitle {
    font-size: 14px;
  }
  
  .actions-row {
    text-align: center;
  }
}
</style>