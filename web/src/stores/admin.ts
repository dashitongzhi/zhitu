import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest } from '@/types/models'
import * as authApi from '@/api/auth'
import { message } from 'ant-design-vue'
import router from '@/router'

// 管理员认证 store，独立持久化到 zhitu-admin-auth 键名
// 与普通用户 auth store 完全隔离，token 不混用
export const useAdminAuthStore = defineStore('adminAuth', () => {
  // 状态
  const adminToken = ref<string>('')
  const adminUser = ref<User | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!adminToken.value)

  // 管理员登录
  const adminLogin = async (credentials: LoginRequest) => {
    try {
      const response = await authApi.adminLogin(credentials)
      // 后端统一响应：{ code, message, data: { token, token_type, expires_in } }
      const authToken = response.data?.data?.token
      if (!authToken) {
        message.error('登录失败：未返回 token')
        return false
      }

      adminToken.value = authToken
      // 后端 admin login 不返回 user 详情，用邮箱占位
      adminUser.value = {
        id: 0,
        email: credentials.email,
        nickname: '管理员',
        avatar: '',
        created_at: '',
        updated_at: '',
      }

      message.success('管理员登录成功')
      router.push('/admin')

      return true
    } catch (error) {
      console.error('管理员登录失败:', error)
      return false
    }
  }

  // 退出登录
  const adminLogout = () => {
    adminToken.value = ''
    adminUser.value = null
    router.push('/admin/login')
    message.success('已退出登录')
  }

  return {
    // 状态
    adminToken,
    adminUser,
    isAuthenticated,
    // 操作
    adminLogin,
    adminLogout,
  }
}, {
  // 持久化配置：独立键名，与用户 token 隔离
  persist: {
    key: 'zhitu-admin-auth',
    paths: ['adminToken', 'adminUser'],
  },
})
