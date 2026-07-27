<template>
  <a-layout class="app-layout preview-layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :width="224"
      :collapsed-width="76"
      :trigger="null"
      class="app-sider"
      theme="light"
    >
      <!-- Logo 区域 -->
      <div class="logo">
        <span class="logo-icon preview-brand-mark" aria-hidden="true">职</span>
        <span v-show="!collapsed" class="logo-text">职途</span>
      </div>

      <!-- 导航菜单 -->
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        :style="{ borderRight: 0 }"
      >
        <a-menu-item key="resumes" @click="navigateTo('/app-style-preview/resumes')">
          <template #icon>
            <FileTextOutlined />
          </template>
          <span>简历实验室</span>
        </a-menu-item>

        <a-menu-item key="interviews" @click="navigateTo('/app-style-preview/interviews')">
          <template #icon>
            <CommentOutlined />
          </template>
          <span>面试训练场</span>
        </a-menu-item>

        <a-menu-item key="deliveries" @click="navigateTo('/app-style-preview/deliveries')">
          <template #icon>
            <SendOutlined />
          </template>
          <span>投递看板</span>
        </a-menu-item>
      </a-menu>

      <!-- 自定义折叠按钮（仿菜单项样式） -->
      <div class="sider-footer">
        <button class="collapse-btn" @click="collapsed = !collapsed" :title="collapsed ? '展开侧边栏' : '折叠侧边栏'">
          <MenuFoldOutlined v-if="!collapsed" />
          <MenuUnfoldOutlined v-else />
          <span v-show="!collapsed" class="collapse-text">折叠菜单</span>
        </button>
        <p class="sider-copyright" title="职途 © 2026 - 让求职更简单">
          <span v-if="collapsed">职途</span>
          <span v-else>职途 © 2026 - 让求职更简单</span>
        </p>
      </div>
    </a-layout-sider>

    <!-- 右侧布局 -->
    <a-layout class="main-layout">
      <!-- 顶部导航栏 -->
      <a-layout-header class="app-header">
        <div class="header-left">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/app-style-preview/resumes">职途</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>
              {{ currentRouteTitle }}
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <!-- 用户信息下拉菜单 -->
          <a-dropdown overlay-class-name="user-dropdown-overlay">
            <div class="user-info">
              <a-avatar :size="32" class="avatar">
                {{ userInitial }}
              </a-avatar>
              <span class="user-name">{{ authStore.user?.nickname || '用户' }}</span>
              <DownOutlined class="arrow-icon" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile" @click="showProfileModal = true">
                  <UserOutlined />
                  <span class="ml-2">个人资料</span>
                </a-menu-item>
                <a-menu-item key="password" @click="navigateTo('/app-style-preview/settings/password')">
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
      <a-layout-content class="app-content" :class="{ 'editor-content': currentRouteName === 'PreviewResumeEditor', 'kanban-content': currentRouteName === 'PreviewDeliveryKanban' }">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>

    </a-layout>

    <!-- 个人资料弹窗 -->
    <UserProfileModal v-model:open="showProfileModal" />
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserProfileModal from '@/components/UserProfileModal.vue'
import {
  UserOutlined,
  FileTextOutlined,
  CommentOutlined,
  SendOutlined,
  DownOutlined,
  LogoutOutlined,
  LockOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 侧边栏折叠状态
const collapsed = ref(false)
const autoCollapsedForViewport = ref(false)

// 个人资料弹窗
const showProfileModal = ref(false)

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
    PreviewResumeList: '简历实验室',
    PreviewResumeTemplateSelect: '选择简历模板',
    PreviewResumeEditor: '简历实验室',
    PreviewInterviewList: '面试训练场',
    PreviewInterviewSceneSelect: '选择训练场景',
    PreviewInterviewRoom: '面试训练场',
    PreviewDeliveryKanban: '投递看板',
    PreviewChangePassword: '修改密码',
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
    if (path === '/app-style-preview' || path === '/app-style-preview/') {
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
  --primary: #1757d2;
  --brand-50: #edf2ff;
  --brand-500: #1757d2;
  --brand-600: #0d47bb;
  --background-50: #ffffff;
  --background-100: #f4f5f2;
  --background-200: #eceeeb;
  --background-300: #d8dce1;
  --foreground: #151a23;
  --muted-foreground: #68717e;
  --border: #d8dce1;
  --card: #ffffff;
  --radius: 0.5rem;
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --shadow-sm: none;
  --shadow-md: none;
  --shadow-lg: none;
  min-height: 100vh;
  background: #f4f5f2;
}

.main-layout {
  min-width: 0;
  flex: 1 1 0;
  background: #f4f5f2;
}

/* ===== 侧边栏：Pinguo 设计稿对齐 ===== */
.app-sider {
  background: #f0f1ee !important;
  border-right: 1px solid #d8dce1;
  box-shadow: none !important;
  z-index: 10;
  position: sticky !important;
  top: 0;
  align-self: flex-start;
  height: 100vh;
}

.app-sider :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
  padding: 0 14px 10px;
}

.logo {
  height: 72px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 8px;
  border-bottom: 1px solid #d7dbdf;
  justify-content: flex-start;
}

.logo-icon {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  background: #151a23;
  padding: 0;
  box-shadow: none;
}

.preview-brand-mark {
  display: grid;
  place-items: center;
  color: #fff;
  font-size: 16px;
  font-weight: 800;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--sidebar-foreground);
  margin-left: 0;
}

/* ===== 导航菜单：覆盖 Ant Design 默认样式 ===== */
.app-sider :deep(.ant-menu) {
  background: transparent;
  border-right: 0;
  padding-top: 14px;
  gap: 4px;
}

.app-sider :deep(.ant-menu-item) {
  min-height: 50px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--sidebar-foreground);
  margin: 0 0 4px 0;
  transition: background-color 0.18s ease, color 0.18s ease;
}

.app-sider :deep(.ant-menu-item:hover) {
  background: var(--sidebar-accent);
  color: var(--sidebar-foreground);
}

.app-sider :deep(.ant-menu-item-selected) {
  background: rgba(255, 255, 255, 0.66) !important;
  color: var(--primary) !important;
  box-shadow: none;
}

.app-sider :deep(.ant-menu-item-selected::after) {
  display: none;
}

.app-sider :deep(.ant-menu-item .anticon) {
  font-size: 18px;
}

/* ===== 自定义折叠按钮（仿菜单项样式） ===== */
.sider-footer {
  margin-top: auto;
  padding: 14px 4px 6px;
  border-top: 1px solid #d7dbdf;
}
.collapse-btn {
  width: 100%;
  height: 40px;
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 0 11px;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: var(--sidebar-foreground);
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.18s ease, color 0.18s ease;
}
.collapse-btn:hover {
  background: var(--sidebar-accent);
  color: var(--sidebar-foreground);
}
.collapse-btn :deep(.anticon) {
  font-size: 18px;
  flex-shrink: 0;
}
.collapse-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
/* 折叠状态下居中显示图标 */
.app-sider.ant-layout-sider-collapsed .collapse-btn {
  justify-content: center;
  padding: 0;
}
.sider-copyright {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  margin: 8px 0 0;
  color: var(--muted-foreground);
  font-size: 10px;
  line-height: 1.45;
  text-align: center;
  white-space: nowrap;
}
.app-sider.ant-layout-sider-collapsed .sider-copyright {
  font-size: 9px;
}

/* ===== 顶部栏 ===== */
.app-header {
  background: rgba(255, 255, 255, 0.92);
  padding: 0 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #d8dce1;
  box-shadow: none;
  height: 64px;
  z-index: 9;
  position: sticky;
  top: 0;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left :deep(.ant-breadcrumb) {
  font-size: 14px;
}

.header-left :deep(.ant-breadcrumb a) {
  color: var(--muted-foreground);
}

.header-left :deep(.ant-breadcrumb-separator) {
  color: var(--muted-foreground);
}

.header-left :deep(.ant-breadcrumb > span:last-child) {
  color: var(--foreground);
  font-weight: 500;
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
  padding: 6px 10px;
  border-radius: 8px;
  transition: background-color 0.18s ease;
  gap: 10px;
}

.user-info:hover {
  background-color: var(--background-200);
}

:global(.user-dropdown-overlay .ant-dropdown-menu) {
  box-shadow: var(--shadow-lg);
}

.avatar {
  border-radius: 8px !important;
  background-color: var(--brand-50) !important;
  color: var(--primary) !important;
  font-size: 13px;
  font-weight: 600;
}

.user-name {
  margin-left: 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--foreground);
}

.arrow-icon {
  margin-left: 0;
  font-size: 12px;
  color: var(--muted-foreground);
}

/* ===== 内容区 ===== */
.app-content {
  box-sizing: border-box;
  min-width: 0;
  width: auto;
  margin: 26px 30px 40px;
  padding: 0;
  background: transparent;
  border-radius: 0;
  border: 0;
  min-height: calc(100vh - 64px - 48px);
  box-shadow: none;
}

/* ===== 复制版主题桥接：只覆盖现有业务组件的视觉，不改变结构或内容 ===== */
.preview-layout :deep(.page-header) {
  margin-bottom: 30px;
  padding-bottom: 22px;
  border-bottom: 1px solid var(--border);
}

.preview-layout :deep(.page-title) {
  color: var(--foreground);
  font-size: 32px;
  font-weight: 720;
  letter-spacing: -0.045em;
  line-height: 1.12;
}

.preview-layout :deep(.page-subtitle) {
  margin-top: 10px;
  color: var(--muted-foreground);
  font-size: 13px;
}

.preview-layout :deep(.btn-create),
.preview-layout :deep(.new-btn),
.preview-layout :deep(.btn-primary-capsule),
.preview-layout :deep(.export-button),
.preview-layout :deep(.action-btn) {
  border-radius: 6px !important;
  background: var(--primary) !important;
  border-color: var(--primary) !important;
  box-shadow: none !important;
}

.preview-layout :deep(.btn-secondary-capsule),
.preview-layout :deep(.ant-btn),
.preview-layout :deep(.toolbar-select),
.preview-layout :deep(.zoom-control) {
  border-radius: 6px !important;
  box-shadow: none !important;
}

.preview-layout :deep(.resume-card),
.preview-layout :deep(.table-card),
.preview-layout :deep(.stat-card),
.preview-layout :deep(.platform-card),
.preview-layout :deep(.platform-hero-card),
.preview-layout :deep(.kb-table-card),
.preview-layout :deep(.kb-mobile-card),
.preview-layout :deep(.editor-panel),
.preview-layout :deep(.dim-card),
.preview-layout :deep(.strategy-card),
.preview-layout :deep(.ant-card),
.preview-layout :deep(.ant-table-wrapper) {
  border-radius: 8px !important;
  border-color: var(--border) !important;
  box-shadow: none !important;
}

.preview-layout :deep(.resume-card:hover),
.preview-layout :deep(.stat-card:hover),
.preview-layout :deep(.kb-mobile-card:hover) {
  transform: none !important;
  box-shadow: none !important;
  border-color: #aeb5bf !important;
}

.preview-layout :deep(.ant-input),
.preview-layout :deep(.ant-input-affix-wrapper),
.preview-layout :deep(.ant-input-password),
.preview-layout :deep(.ant-picker),
.preview-layout :deep(.ant-select-selector) {
  border-radius: 6px !important;
}

.preview-layout :deep(.ant-tag),
.preview-layout :deep(.status-tag),
.preview-layout :deep(.scene-tag),
.preview-layout :deep(.mode-tag) {
  border-radius: 999px !important;
}

.preview-layout :deep(.qc-input-wrap),
.preview-layout :deep(.search-wrap),
.preview-layout :deep(.filter-select),
.preview-layout :deep(.view-toggle-group),
.preview-layout :deep(.empty-state) {
  border-radius: 8px !important;
  box-shadow: none !important;
}

.preview-layout :deep(.stat-grid) {
  gap: 0 !important;
  overflow: hidden;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--card);
}

.preview-layout :deep(.stat-card) {
  border: 0 !important;
  border-right: 1px solid var(--border) !important;
  border-radius: 0 !important;
  background: transparent !important;
}

.preview-layout :deep(.stat-card:last-child) {
  border-right: 0 !important;
}

.preview-layout :deep(.stat-icon-wrap) {
  background: var(--brand-50) !important;
  color: var(--primary) !important;
}

.preview-layout :deep(.stat-spark polyline) {
  stroke: var(--primary) !important;
}

.preview-layout :deep(.view-toggle-btn) {
  border-radius: 6px !important;
}

@media (max-width: 1100px) {
  .app-content {
    margin: 16px;
    padding: 16px;
  }
}

.app-content.editor-content {
  height: calc(100vh - 64px);
  min-height: 0;
  margin: 0;
  padding: 0;
  border-radius: 0;
  border: none;
  overflow: hidden;
  box-shadow: none;
}

.app-content.kanban-content {
  width: 100%;
  margin: 0;
  padding: 0;
  border-radius: 0;
  border: none;
  box-shadow: none;
  background: var(--background-100);
  max-width: none;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s cubic-bezier(0.32, 0.72, 0, 1);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* 下拉菜单项样式 */
.ml-2 {
  margin-left: 8px;
}

/* 响应式：窄屏隐藏用户名 */
@media (max-width: 768px) {
  .user-name {
    display: none;
  }
}
</style>
