<template>
  <div class="not-found-view">
    <el-result
      icon="warning"
      title="404"
      sub-title="Sorry, the page you visited does not exist."
    >
      <template #extra>
        <el-button type="primary" @click="goHome">Back Home</el-button>
        <el-button @click="goBack">Go Back</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.not-found-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 20px;
}

:deep(.el-result__icon svg) {
  color: #f56c6c;
}

:deep(.el-result__title) {
  color: #303133;
}

:deep(.el-result__subtitle) {
  color: #909399;
}
</style>