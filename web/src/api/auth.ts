import { get, post } from '@/utils/request'
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  User,
  ChangePasswordRequest,
} from '@/types/models'

// 用户登录
export const login = (data: LoginRequest) => {
  return post<LoginResponse>('/api/auth/login', data)
}

// 用户注册
export const register = (data: RegisterRequest) => {
  return post<LoginResponse>('/api/auth/register', data)
}

// 获取当前用户信息
export const getCurrentUser = () => {
  return get<ApiResponse<User>>('/api/auth/me')
}

// 修改密码
export const changePassword = (data: ChangePasswordRequest) => {
  return post<ApiResponse>('/api/auth/change-password', data)
}

// 管理员登录
export const adminLogin = (data: LoginRequest) => {
  return post<LoginResponse>('/api/auth/admin/login', data)
}

// 注：后端无 /api/auth/logout 接口，退出登录仅在前端清除 token
