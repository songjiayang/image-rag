<template>
  <div class="dashboard">
    <el-row :gutter="20" class="dashboard-header">
      <el-col :span="24">
        <h1>Image RAG Service Dashboard</h1>
        <p class="subtitle">Manage and search your image collections with AI-powered vector search</p>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#409eff"><Document /></el-icon>
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalRecords }}</div>
              <div class="stat-label">Total Records</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#67c23a"><Picture /></el-icon>
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalImages }}</div>
              <div class="stat-label">Total Images</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#e6a23c"><DocumentAdd /></el-icon>
            <div class="stat-info">
              <div class="stat-number">{{ stats.todayRecords }}</div>
              <div class="stat-label">Today's Records</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#f56c6c"><Picture /></el-icon>
            <div class="stat-info">
              <div class="stat-number">{{ stats.todayImages }}</div>
              <div class="stat-label">Today's Images</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="action-row">
      <el-col :span="12">
        <el-card class="action-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span><el-icon><Plus /></el-icon> Quick Add Record</span>
            </div>
          </template>
          <div class="action-content">
            <p>Add new images to your collection with metadata</p>
            <el-button type="primary" @click="goToRecords" size="large">
              <el-icon><Plus /></el-icon>
              Add New Record
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="action-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span><el-icon><Search /></el-icon> Quick Search</span>
            </div>
          </template>
          <div class="action-content">
            <p>Search similar images using AI-powered vector search</p>
            <el-button type="success" @click="goToSearch" size="large">
              <el-icon><Search /></el-icon>
              Search Images
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="recent-row">
      <el-col :span="24">
        <el-card class="recent-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span><el-icon><Clock /></el-icon> Recent Records</span>
              <el-button text @click="goToRecords" type="primary">
                View All
              </el-button>
            </div>
          </template>
          
          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="3" animated />
          </div>
          
          <div v-else-if="recentRecords.length === 0" class="empty-container">
            <el-empty description="No records yet. Add your first record!" />
          </div>
          
          <el-row v-else :gutter="20">
            <el-col 
              v-for="record in recentRecords" 
              :key="record.id" 
              :xs="24" 
              :sm="12" 
              :md="8" 
              :lg="6"
            >
              <el-card class="record-card" shadow="hover" @click="goToRecord(record.id)">
                <div class="record-image">
                  <img 
                    v-if="record.images.length > 0"
                    :src="`/api/v1/images/${record.images[0].id}/preview`" 
                    :alt="record.name"
                  />
                  <div v-else class="no-image">
                    <el-icon><Picture /></el-icon>
                  </div>
                </div>
                <div class="record-info">
                  <h4>{{ record.name }}</h4>
                  <p class="record-description">{{ record.description || 'No description' }}</p>
                  <p class="record-date">{{ formatDate(record.created_at) }}</p>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { recordService, statsService } from '@/services/api'
import type { Record } from '@/types'
import { formatDate } from '@/utils/date-utils'
import { 
  Document, 
  Picture, 
  DocumentAdd, 
  Plus, 
  Search, 
  Clock 
} from '@element-plus/icons-vue'

const router = useRouter()

const stats = ref({
  totalRecords: 0,
  totalImages: 0,
  todayRecords: 0,
  todayImages: 0
})

const recentRecords = ref<Record[]>([])
const loading = ref(false)

const loadDashboardData = async () => {
  loading.value = true
  try {
    const [statsResponse, recentResponse] = await Promise.all([
      statsService.getDashboardStats(),
      recordService.getRecords(1, 4)
    ])
    
    stats.value.totalRecords = statsResponse.data.total_records
    stats.value.totalImages = statsResponse.data.total_images
    stats.value.todayRecords = statsResponse.data.today_records
    stats.value.todayImages = statsResponse.data.today_images
    
    recentRecords.value = recentResponse.data
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  } finally {
    loading.value = false
  }
}

const goToRecords = () => {
  router.push('/records')
}

const goToSearch = () => {
  router.push('/search')
}

const goToRecord = (id: number) => {
  router.push(`/records/${id}`)
}


onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 30px;
  text-align: center;
}

.dashboard-header h1 {
  font-size: 32px;
  color: #303133;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 16px;
  color: #909399;
  margin: 0;
}

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  height: 100%;
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 20px 0;
}

.stat-icon {
  font-size: 48px;
  margin-right: 16px;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.action-row {
  margin-bottom: 30px;
}

.action-card {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: bold;
}

.action-content {
  text-align: center;
  padding: 20px 0;
}

.action-content p {
  color: #909399;
  margin-bottom: 20px;
}

.recent-row {
  margin-bottom: 20px;
}

.recent-card {
  height: 100%;
}

.loading-container {
  padding: 40px;
  text-align: center;
}

.empty-container {
  padding: 40px;
  text-align: center;
}

.record-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.2s;
}

.record-card:hover {
  transform: translateY(-2px);
}

.record-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 12px;
}

.record-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  color: #909399;
  font-size: 48px;
}

.record-info h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
}

.record-description {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.record-date {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

@media (max-width: 768px) {
  .dashboard-header h1 {
    font-size: 24px;
  }
  
  .subtitle {
    font-size: 14px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-icon {
    font-size: 36px;
  }
}
</style>