import { get, post, put, del } from '@/utils/request'
import type {
  ApiResponse,
  FullProfile,
  UpdateProfileRequest,
  ParseResumeResult,
  Education,
  Work,
  Project,
  Skill,
  Honor,
  Practice,
} from '@/types/models'

// ==================== 档案主表 ====================

// 获取完整档案（含子资源）
export const getProfile = () => {
  return get<ApiResponse<FullProfile>>('/api/v1/profile')
}

// 更新基础信息（仅白名单字段生效）
export const updateProfile = (data: UpdateProfileRequest) => {
  return put<ApiResponse<FullProfile>>('/api/v1/profile', data)
}

// 获取完成度
export const getCompletion = () => {
  return get<ApiResponse<{ completion_pct: number }>>('/api/v1/profile/completion')
}

// 上传简历解析（multipart/form-data）
export const parseResume = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return post<ApiResponse<ParseResumeResult>>(
    '/api/v1/profile/parse-resume',
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
    }
  )
}

// ==================== 教育背景 ====================

export const listEducations = () => {
  return get<ApiResponse<Education[]>>('/api/v1/profile/educations')
}

export const createEducation = (data: Omit<Education, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Education>>('/api/v1/profile/educations', data)
}

export const updateEducation = (id: number, data: Partial<Education>) => {
  return put<ApiResponse>(`/api/v1/profile/educations/${id}`, data)
}

export const deleteEducation = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/educations/${id}`)
}

// ==================== 工作经历 ====================

export const listWorks = () => {
  return get<ApiResponse<Work[]>>('/api/v1/profile/works')
}

export const createWork = (data: Omit<Work, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Work>>('/api/v1/profile/works', data)
}

export const updateWork = (id: number, data: Partial<Work>) => {
  return put<ApiResponse>(`/api/v1/profile/works/${id}`, data)
}

export const deleteWork = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/works/${id}`)
}

// ==================== 项目经历 ====================

export const listProjects = () => {
  return get<ApiResponse<Project[]>>('/api/v1/profile/projects')
}

export const createProject = (data: Omit<Project, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Project>>('/api/v1/profile/projects', data)
}

export const updateProject = (id: number, data: Partial<Project>) => {
  return put<ApiResponse>(`/api/v1/profile/projects/${id}`, data)
}

export const deleteProject = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/projects/${id}`)
}

// ==================== 技能 ====================

export const listSkills = () => {
  return get<ApiResponse<Skill[]>>('/api/v1/profile/skills')
}

export const createSkill = (data: Omit<Skill, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Skill>>('/api/v1/profile/skills', data)
}

export const updateSkill = (id: number, data: Partial<Skill>) => {
  return put<ApiResponse>(`/api/v1/profile/skills/${id}`, data)
}

export const deleteSkill = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/skills/${id}`)
}

// ==================== 荣誉奖项 ====================

export const listHonors = () => {
  return get<ApiResponse<Honor[]>>('/api/v1/profile/honors')
}

export const createHonor = (data: Omit<Honor, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Honor>>('/api/v1/profile/honors', data)
}

export const updateHonor = (id: number, data: Partial<Honor>) => {
  return put<ApiResponse>(`/api/v1/profile/honors/${id}`, data)
}

export const deleteHonor = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/honors/${id}`)
}

// ==================== 校内外实践 ====================

export const listPractices = () => {
  return get<ApiResponse<Practice[]>>('/api/v1/profile/practices')
}

export const createPractice = (data: Omit<Practice, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
  return post<ApiResponse<Practice>>('/api/v1/profile/practices', data)
}

export const updatePractice = (id: number, data: Partial<Practice>) => {
  return put<ApiResponse>(`/api/v1/profile/practices/${id}`, data)
}

export const deletePractice = (id: number) => {
  return del<ApiResponse>(`/api/v1/profile/practices/${id}`)
}
