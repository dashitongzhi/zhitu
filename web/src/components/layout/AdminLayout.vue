<template>
  <a-layout class="admin-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :width="240"
      :collapsed-width="80"
      class="admin-sider"
      theme="dark"
    >
      <div class="logo">
        <img src="/favicon.svg" alt="职途AI" class="logo-icon" />
        <span v-show="!collapsed" class="logo-text">职途AI管理</span>
      </div>

      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        theme="dark"
        :style="{ borderRight: 0 }"
      >
        <a-menu-item key="dashboard" @click="navigateTo('/admin')">
          <template #icon>
            <DashboardOutlined />
          </template>
          <span>仪表盘</span>
        </a-menu-item>

        <a-menu-item key="users" @click="navigateTo('/admin/users')">
          <template #icon>
            <TeamOutlined />
          </template>
          <span>用户管理</span>
        </a-menu-item>

        <a-menu-item key="deliveries" @click="navigateTo('/admin/deliveries')">
          <template #icon>
            <SendOutlined />
          </template>
          <span>投递管理</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="admin-header">
        <div class="header-left">
          <span class="header-title">管理后台</span>
        </div>

        <div class="header-right">
          <span class="admin-email">{{ adminStore.adminUser?.email || '管理员' }}</span>
          <a-button type="text" @click="handleLogout">
            <LogoutOutlined />
            退出
          </a-button>
        </div>
      </a-layout-header>

      <a-layout-content class="admin-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>

      <a-layout-footer class="admin-footer">
        职途AI 管理后台 © 2024
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminAuthStore } from '@/stores/admin'
import {
  DashboardOutlined,
  TeamOutlined,
  SendOutlined,
  LogoutOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminAuthStore()

const collapsed = ref(false)
const selectedKeys = ref<string[]>(['dashboard'])

watch(
  () => route.path,
  (path) => {
    if (path === '/admin') {
      selectedKeys.value = ['dashboard']
    } else if (path.startsWith('/admin/users')) {
      selectedKeys.value = ['users']
    } else if (path.startsWith('/admin/deliveries')) {
      selectedKeys.value = ['deliveries']
    }
  },
  { immediate: true }
)

const navigateTo = (path: string) => {
  router.push(path)
}

const handleLogout = () => {
  adminStore.adminLogout()
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
}

.admin-sider {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-left: 12px;
}

.admin-header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  z-index: 9;
  height: 64px;
  line-height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.admin-email {
  color: #666;
  font-size: 14px;
}

.admin-content {
  margin: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  min-height: calc(100vh - 64px - 70px - 48px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.admin-footer {
  text-align: center;
  color: #999;
  font-size: 14px;
  padding: 24px;
  background: #fff;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
