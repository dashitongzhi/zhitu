import { get, post, patch } from '@/utils/request'
import type {
  ApiResponse,
  AdminDashboardStats,
  AdminUserListResponse,
  AdminUserListParams,
  AdminUserDetail,
  AdminDeliveryListResponse,
  AdminDeliveryListParams,
  AdminDeliveryFunnel,
  UserStatus,
} from '@/types/models'

// 获取仪表盘统计
export const getStats = () => {
  return get<ApiResponse<AdminDashboardStats>>('/api/admin/stats')
}

// 获取用户列表
export const getUsers = (params: AdminUserListParams) => {
  return get<ApiResponse<AdminUserListResponse>>('/api/admin/users', { params })
}

// 获取用户详情
export const getUserDetail = (id: number) => {
  return get<ApiResponse<AdminUserDetail>>(`/api/admin/users/${id}`)
}

// 切换用户状态
export const toggleUserStatus = (id: number, status: UserStatus) => {
  return patch<ApiResponse>(`/api/admin/users/${id}/status`, { status })
}

// 重置用户密码（后端要求至少 6 位，建议前端校验 8 位）
export const resetUserPassword = (id: number, newPassword: string) => {
  return post<ApiResponse>(`/api/admin/users/${id}/reset-password`, {
    new_password: newPassword,
  })
}

// 获取投递列表
export const getDeliveries = (params: AdminDeliveryListParams) => {
  return get<ApiResponse<AdminDeliveryListResponse>>('/api/admin/deliveries', {
    params,
  })
}

// 获取投递漏斗
export const getFunnel = () => {
  return get<ApiResponse<AdminDeliveryFunnel>>('/api/admin/deliveries/funnel')
}
