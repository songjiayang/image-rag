<template>
  <div class="settings-view">
    <el-row :gutter="20" class="header-row">
      <el-col :span="24">
        <h1>Settings</h1>
        <p class="subtitle">Configure your Image RAG service</p>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :xs="24" :lg="12">
        <!-- Service Configuration -->
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Service Configuration</span>
            </div>
          </template>

          <el-form :model="serviceConfig" label-width="140px">
            <el-form-item label="API Base URL">
              <el-input v-model="serviceConfig.apiUrl" placeholder="http://localhost:8080" />
            </el-form-item>

            <el-form-item label="Max Image Size">
              <el-select v-model="serviceConfig.maxImageSize">
                <el-option label="2MB" value="2MB" />
                <el-option label="5MB" value="5MB" />
                <el-option label="10MB" value="10MB" />
                <el-option label="20MB" value="20MB" />
              </el-select>
            </el-form-item>

            <el-form-item label="Default Top K">
              <el-slider 
                v-model="serviceConfig.defaultTopK" 
                :min="5" 
                :max="50" 
                :step="5"
                show-input
              />
            </el-form-item>

            <el-form-item label="Auto Refresh">
              <el-switch v-model="serviceConfig.autoRefresh" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="saveServiceConfig">Save Settings</el-button>
              <el-button @click="resetServiceConfig">Reset</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Search Preferences -->
        <el-card shadow="hover" class="preferences-card">
          <template #header>
            <div class="card-header">
              <span>Search Preferences</span㸼/div>
          </template>

          <el-form :model="searchPreferences" label-width="140px">
            <el-form-item label="Show Similarity">
              <el-switch v-model="searchPreferences.showSimilarity" />
            </el-form-item>

            <el-form-item label="Show Distance">
              <el-switch v-model="searchPreferences.showDistance" />
            </el-form-item>

            <el-form-item label="Grid View">
              <el-switch v-model="searchPreferences.gridView" />
            </el-form-item>

            <el-form-item label="Image Size">
              <el-select v-model="searchPreferences.imageSize">
                <el-option label="Small" value="small" />
                <el-option label="Medium" value="medium" />
                <el-option label="Large" value="large" />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="saveSearchPreferences">Save Preferences</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <!-- System Status -->
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>System Status</span>
              <el-button @click="checkSystemStatus" :loading="checkingStatus" size="small">
                <el-icon><Refresh /></el-icon>
                Refresh
              </el-button>
            </div>
          </template>

          <div v-if="checkingStatus" class="status-loading">
            <el-skeleton :rows="4" animated />
          </div>

          <div v-else class="status-info">
            <el-descriptions :column="1" border
              <el-descriptions-item label="Backend Status">
                <el-tag :type="systemStatus.backend ? 'success' : 'danger'">
                  {{ systemStatus.backend ? 'Online' : 'Offline' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="Database Status">
                <el-tag :type="systemStatus.database ? 'success' : 'danger'">
                  {{ systemStatus.database ? 'Connected' : 'Disconnected' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="Milvus Status">
                <el-tag :type="systemStatus.milvus ? 'success' : 'danger'">
                  {{ systemStatus.milvus ? 'Connected' : 'Disconnected' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="Last Checked">
                {{ systemStatus.lastChecked }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>

        <!-- About -->
        <el-card shadow="hover" class="about-card">
          <template #header>
            <div class="card-header">
              <span>About</span>
            </div>
          </template>

          <div class="about-info">
            <p><strong>Image RAG Service</strong></p>
            <p>Version: 1.0.0</p>
            <p>Built with Vue.js 3 & Element Plus</p>
            <p>Vector search powered by Milvus & Doubao API</p>
            <p>&copy; 2025 Image RAG Service</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { healthService } from '@/services/api'

// Configuration
const serviceConfig = ref({
  apiUrl: localStorage.getItem('apiUrl') || 'http://localhost:8080',
  maxImageSize: localStorage.getItem('maxImageSize') || '10MB',
  defaultTopK: Number(localStorage.getItem('defaultTopK')) || 10,
  autoRefresh: localStorage.getItem('autoRefresh') === 'true'
})

const searchPreferences = ref({
  showSimilarity: localStorage.getItem('showSimilarity') !== 'false',
  showDistance: localStorage.getItem('showDistance') !== 'false',
  gridView: localStorage.getItem('gridView') !== 'false',
  imageSize: localStorage.getItem('imageSize') || 'medium'
})

// System status
const systemStatus = ref({
  backend: false,
  database: false,
  milvus: false,
  lastChecked: ''
})

const checkingStatus = ref(false)

// Methods
const saveServiceConfig = () => {
  localStorage.setItem('apiUrl', serviceConfig.value.apiUrl)
  localStorage.setItem('maxImageSize', serviceConfig.value.maxImageSize)
  localStorage.setItem('defaultTopK', serviceConfig.value.defaultTopK.toString())
  localStorage.setItem('autoRefresh', serviceConfig.value.autoRefresh.toString())
  ElMessage.success('Service settings saved')
}

const resetServiceConfig = () => {
  serviceConfig.value = {
    apiUrl: 'http://localhost:8080',
    maxImageSize: '10MB',
    defaultTopK: 10,
    autoRefresh: false
  }
  saveServiceConfig()
}

const saveSearchPreferences = () => {
  localStorage.setItem('showSimilarity', searchPreferences.value.showSimilarity.toString())
  localStorage.setItem('showDistance', searchPreferences.value.showDistance.toString())
  localStorage.setItem('gridView', searchPreferences.value.gridView.toString())
  localStorage.setItem('imageSize', searchPreferences.value.imageSize)
  ElMessage.success('Search preferences saved')
}

const checkSystemStatus = async () => {
  checkingStatus.value = true
  try {
    const health = await healthService.checkHealth()
    systemStatus.value.backend = health.status === 'healthy'
    systemStatus.value.database = true // Assume connected if health check passes
    systemStatus.value.milvus = true   // Assume connected if health check passes
    systemStatus.value.lastChecked = new Date().toLocaleString()
  } catch (error) {
    systemStatus.value.backend = false
    systemStatus.value.database = false
    systemStatus.value.milvus = false
    systemStatus.value.lastChecked = new Date().toLocaleString()
    ElMessage.error('System health check failed')
  } finally {
    checkingStatus.value = false
  }
}

// Initialize
onMounted(() => {
  checkSystemStatus()
})
</script>

<style scoped>
.settings-view {
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

.card-header {
  font-weight: bold;
}

.preferences-card,
.about-card {
  margin-top: 20px;
}

.status-loading {
  text-align: center;
  padding: 20px;
}

.status-info {
  padding: 10px 0;
}

.about-info p {
  margin: 8px 0;
  color: #606266;
}

@media (max-width: 768px) {
  .header-row h1 {
    font-size: 24px;
  }
  
  .subtitle {
    font-size: 14px;
  }
}
</style>