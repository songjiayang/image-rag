<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <div class="logo">
          <el-icon><Picture /></el-icon>
          <span>Image RAG Service</span>
        </div>
        <el-menu
          :default-active="activeIndex"
          class="header-menu"
          mode="horizontal"
          router
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>Dashboard</span>
          </el-menu-item>
          <el-menu-item index="/records">
            <el-icon><Files /></el-icon>
            <span>Records</span>
          </el-menu-item>
          <el-menu-item index="/search">
            <el-icon><Search /></el-icon>
            <span>Search</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>Settings</span>
          </el-menu-item>
        </el-menu>
      </div>
    </el-header>

    <el-main class="app-main">
      <router-view />
    </el-main>

    <el-footer class="app-footer">
      <div class="footer-content">
        <span>&copy; 2025 Image RAG Service. Built with Vue.js & Element Plus.</span>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeIndex = computed(() => route.path)
</script>

<style scoped>
.app-container {
  min-height: 100vh;
}

.app-header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
}

.logo .el-icon {
  margin-right: 8px;
  font-size: 24px;
}

.header-menu {
  border-bottom: none;
}

.app-main {
  background-color: #f5f7fa;
  padding: 20px;
  min-height: calc(100vh - 120px);
}

.app-footer {
  background-color: #fff;
  border-top: 1px solid #e4e7ed;
  text-align: center;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-content {
  color: #909399;
  font-size: 14px;
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    height: auto;
    padding: 10px;
  }

  .header-menu {
    margin-top: 10px;
  }

  .app-main {
    padding: 10px;
  }
}
</style>