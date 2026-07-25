import { get, post, put, del } from '@/utils/request'
import { streamSSE } from '@/utils/sse'
import type {
  ApiResponse,
  Resume,
  CreateResumeRequest,
  UpdateResumeRequest,
  ResumeVersion,
  CreateVersionRequest,
  GenerateInput,
  PolishInput,
  ScoreResult,
  JDMatchResult,
  SSECallbacks,
} from '@/types/models'

// ==================== 简历主表 CRUD ====================

// 获取简历列表（无分页，返回数组）
export const listResumes = () => {
  return get<ApiResponse<Resume[]>>('/api/v1/resumes')
}

// 兼容旧名称
export const getResumes = listResumes

// 获取单个简历
export const getResume = (id: number) => {
  return get<ApiResponse<Resume>>(`/api/v1/resumes/${id}`)
}

// 创建简历
export const createResume = (data: CreateResumeRequest) => {
  return post<ApiResponse<Resume>>('/api/v1/resumes', data)
}

// 更新简历
export const updateResume = (id: number, data: UpdateResumeRequest) => {
  return put<ApiResponse>(`/api/v1/resumes/${id}`, data)
}

// 删除简历
export const deleteResume = (id: number) => {
  return del<ApiResponse>(`/api/v1/resumes/${id}`)
}

// ==================== 版本管理 ====================

// 列出版本
export const listVersions = (resumeId: number) => {
  return get<ApiResponse<ResumeVersion[]>>(`/api/v1/resumes/${resumeId}/versions`)
}

// 创建版本
export const createVersion = (resumeId: number, data: CreateVersionRequest) => {
  return post<ApiResponse<ResumeVersion>>(`/api/v1/resumes/${resumeId}/versions`, data)
}

// 获取版本
export const getVersion = (resumeId: number, versionId: number) => {
  return get<ApiResponse<ResumeVersion>>(
    `/api/v1/resumes/${resumeId}/versions/${versionId}`
  )
}

// 回滚到指定版本
export const rollbackVersion = (resumeId: number, versionId: number) => {
  return post<ApiResponse<ResumeVersion>>(
    `/api/v1/resumes/${resumeId}/rollback/${versionId}`
  )
}

// ==================== AI 操作 ====================

// AI 生成简历（SSE 流式）
export const aiGenerate = (
  resumeId: number,
  input: GenerateInput,
  callbacks: SSECallbacks,
  signal?: AbortSignal
) => {
  return streamSSE(
    `/api/v1/resumes/${resumeId}/ai/generate`,
    callbacks,
    { body: input, signal }
  )
}

// AI 润色
export const aiPolish = (resumeId: number, input: PolishInput) => {
  return post<ApiResponse<ResumeVersion>>(
    `/api/v1/resumes/${resumeId}/ai/polish`,
    input
  )
}

// AI 评分
export const aiScore = (resumeId: number, jd?: string) => {
  return post<ApiResponse<ScoreResult>>(
    `/api/v1/resumes/${resumeId}/ai/score`,
    { jd: jd || '' }
  )
}

// AI JD 匹配
export const aiJdMatch = (resumeId: number, jd: string) => {
  return post<ApiResponse<JDMatchResult>>(
    `/api/v1/resumes/${resumeId}/ai/jd-match`,
    { jd }
  )
}

// ==================== 同步档案 ====================

// 将当前简历版本内容反向同步回用户档案
export const syncProfile = (resumeId: number) => {
  return post<ApiResponse>(`/api/v1/resumes/${resumeId}/sync-profile`)
}
