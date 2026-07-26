<template>
  <a-layout class="app-layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :width="240"
      :collapsed-width="80"
      class="app-sider"
      theme="light"
    >
      <!-- Logo 区域 -->
      <div class="logo">
        <img src="/favicon.svg" alt="职途AI" class="logo-icon" />
        <span v-show="!collapsed" class="logo-text">职途AI</span>
      </div>

      <!-- 导航菜单 -->
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        :style="{ borderRight: 0 }"
      >
        <a-menu-item key="resumes" @click="navigateTo('/app/resumes')">
          <template #icon>
            <FileTextOutlined />
          </template>
          <span>简历实验室</span>
        </a-menu-item>

        <a-menu-item key="interviews" @click="navigateTo('/app/interviews')">
          <template #icon>
            <CommentOutlined />
          </template>
          <span>面试训练场</span>
        </a-menu-item>

        <a-menu-item key="deliveries" @click="navigateTo('/app/deliveries')">
          <template #icon>
            <SendOutlined />
          </template>
          <span>投递看板</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <!-- 右侧布局 -->
    <a-layout>
      <!-- 顶部导航栏 -->
      <a-layout-header class="app-header">
        <div class="header-left">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/app/resumes">职途AI</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>
              {{ currentRouteTitle }}
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <!-- 用户信息下拉菜单 -->
          <a-dropdown>
            <div class="user-info">
              <a-avatar :size="32" class="avatar">
                {{ userInitial }}
              </a-avatar>
              <span class="user-name">{{ authStore.user?.nickname || '用户' }}</span>
              <DownOutlined class="arrow-icon" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile" @click="navigateTo('/app/profile')">
                  <UserOutlined />
                  <span class="ml-2">个人资料</span>
                </a-menu-item>
                <a-menu-item key="password" @click="navigateTo('/app/settings/password')">
                  <LockOutlined />
                  <span class="ml-2">修改密码</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span class="ml-2">退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 内容区域 -->
      <a-layout-content class="app-content" :class="{ 'editor-content': currentRouteName === 'ResumeEditor' }">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>

      <!-- 底部 -->
      <a-layout-footer v-if="currentRouteName !== 'ResumeEditor'" class="app-footer">
        职途AI © 2024 - 让求职更简单
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  UserOutlined,
  FileTextOutlined,
  CommentOutlined,
  SendOutlined,
  DownOutlined,
  LogoutOutlined,
  LockOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 侧边栏折叠状态
const collapsed = ref(false)
const autoCollapsedForViewport = ref(false)

const syncSidebarToViewport = () => {
  const shouldCollapse = window.innerWidth <= 1100
  if (shouldCollapse === autoCollapsedForViewport.value) return
  autoCollapsedForViewport.value = shouldCollapse
  collapsed.value = shouldCollapse
}

onMounted(() => {
  syncSidebarToViewport()
  window.addEventListener('resize', syncSidebarToViewport)
})

onBeforeUnmount(() => window.removeEventListener('resize', syncSidebarToViewport))

// 菜单选中状态
const selectedKeys = ref<string[]>(['resumes'])
const openKeys = ref<string[]>([])

// 当前路由名称
const currentRouteName = computed(() => route.name as string)

// 当前路由标题
const currentRouteTitle = computed(() => {
  const titleMap: Record<string, string> = {
    ResumeList: '简历实验室',
    ResumeTemplateSelect: '选择简历模板',
    ResumeEditor: '简历实验室',
    InterviewList: '面试训练场',
    InterviewRoom: '面试训练场',
    DeliveryKanban: '投递看板',
    ChangePassword: '修改密码',
  }
  return titleMap[currentRouteName.value] || ''
})

// 用户名首字母
const userInitial = computed(() => {
  const nickname = authStore.user?.nickname || '用户'
  return nickname.charAt(0).toUpperCase()
})

// 监听路由变化，更新菜单选中状态
watch(
  () => route.path,
  (path) => {
    if (path === '/app' || path === '/app/') {
      selectedKeys.value = ['resumes']
    } else {
      // 形如 /app/profile -> 取第三段
      const seg = path.split('/')[2]
      if (seg) selectedKeys.value = [seg]
    }
  },
  { immediate: true }
)

// 导航到指定路径
const navigateTo = (path: string) => {
  router.push(path)
}

// 退出登录
const handleLogout = () => {
  authStore.logout()
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
}

.app-sider {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #1890ff;
  margin-left: 12px;
}

.app-header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  z-index: 9;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.avatar {
  background-color: #1890ff;
}

.user-name {
  margin-left: 8px;
  font-size: 14px;
}

.arrow-icon {
  margin-left: 4px;
  font-size: 12px;
}

.app-content {
  margin: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  min-height: calc(100vh - 64px - 70px - 48px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.app-content.editor-content {
  height: calc(100vh - 64px);
  min-height: 0;
  margin: 0;
  padding: 0;
  border-radius: 0;
  overflow: hidden;
  box-shadow: none;
}

.app-footer {
  text-align: center;
  color: #999;
  font-size: 14px;
  padding: 24px;
  background: #fff;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 下拉菜单项样式 */
.ml-2 {
  margin-left: 8px;
}
</style>
