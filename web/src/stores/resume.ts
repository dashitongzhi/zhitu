import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Resume,
  ResumeVersion,
  CreateResumeRequest,
  UpdateResumeRequest,
  CreateVersionRequest,
  GenerateInput,
  PolishInput,
  ScoreResult,
  JDMatchResult,
} from '@/types/models'
import * as resumeApi from '@/api/resume'
import { message } from 'ant-design-vue'

export const useResumeStore = defineStore('resume', () => {
  // 状态
  const resumes = ref<Resume[]>([])
  const currentResume = ref<Resume | null>(null)
  const versions = ref<ResumeVersion[]>([])
  const currentVersion = ref<ResumeVersion | null>(null)
  const loading = ref(false)
  const aiLoading = ref(false)
  const streamingText = ref('')

  // 获取简历列表（无分页）
  const fetchList = async () => {
    loading.value = true
    try {
      const response = await resumeApi.listResumes()
      resumes.value = response.data.data || []
      return true
    } catch (error) {
      console.error('获取简历列表失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取单个简历
  const fetchOne = async (id: number) => {
    loading.value = true
    try {
      const response = await resumeApi.getResume(id)
      currentResume.value = response.data.data
      return true
    } catch (error) {
      console.error('获取简历失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 创建简历
  const create = async (data: CreateResumeRequest) => {
    try {
      const response = await resumeApi.createResume(data)
      const r = response.data.data
      if (r) resumes.value.unshift(r)
      message.success('简历创建成功')
      return r || null
    } catch (error) {
      console.error('创建简历失败:', error)
      return null
    }
  }

  // 更新简历
  const update = async (id: number, data: UpdateResumeRequest) => {
    try {
      await resumeApi.updateResume(id, data)
      await fetchOne(id)
      message.success('简历更新成功')
      return true
    } catch (error) {
      console.error('更新简历失败:', error)
      return false
    }
  }

  // 删除简历
  const remove = async (id: number) => {
    try {
      await resumeApi.deleteResume(id)
      resumes.value = resumes.value.filter((r) => r.id !== id)
      if (currentResume.value?.id === id) {
        clearCurrent()
      }
      message.success('简历删除成功')
      return true
    } catch (error) {
      console.error('删除简历失败:', error)
      return false
    }
  }

  // 获取版本列表
  const fetchVersions = async (resumeId: number) => {
    try {
      const response = await resumeApi.listVersions(resumeId)
      versions.value = response.data.data || []
      // 默认选中最新版本（列表第一个）
      if (versions.value.length > 0 && !currentVersion.value) {
        currentVersion.value = versions.value[0]
      }
      return true
    } catch (error) {
      console.error('获取版本列表失败:', error)
      return false
    }
  }

  // 创建版本
  const createVersion = async (
    resumeId: number,
    data: CreateVersionRequest
  ) => {
    try {
      const response = await resumeApi.createVersion(resumeId, data)
      const v = response.data.data
      if (v) {
        versions.value.unshift(v)
        currentVersion.value = v
      }
      message.success('版本保存成功')
      return v || null
    } catch (error) {
      console.error('创建版本失败:', error)
      return null
    }
  }

  // 回滚到指定版本
  const rollbackVersion = async (resumeId: number, versionId: number) => {
    try {
      await resumeApi.rollbackVersion(resumeId, versionId)
      message.success('回滚成功')
      await fetchVersions(resumeId)
      return true
    } catch (error) {
      console.error('回滚版本失败:', error)
      return false
    }
  }

  // 切换当前版本
  const setCurrentVersion = (version: ResumeVersion) => {
    currentVersion.value = version
  }

  // AI 生成简历（SSE 流式）
  // onDelta: 每次收到增量文本的回调
  const aiGenerate = async (
    resumeId: number,
    input: GenerateInput,
    onDelta?: (delta: string) => void,
    signal?: AbortSignal
  ) => {
    aiLoading.value = true
    streamingText.value = ''
    try {
      await resumeApi.aiGenerate(
        resumeId,
        input,
        {
          onDelta: (delta) => {
            streamingText.value += delta
            onDelta?.(delta)
          },
          onStatus: (msg) => {
            message.info(msg)
          },
          onDone: () => {
            message.success('AI 生成完成')
          },
          onError: (errMsg) => {
            message.error(errMsg || 'AI 生成失败')
          },
        },
        signal
      )
      // 生成完成后刷新版本列表
      await fetchVersions(resumeId)
      return true
    } catch (error) {
      console.error('AI 生成失败:', error)
      return false
    } finally {
      aiLoading.value = false
    }
  }

  // AI 润色
  const aiPolish = async (resumeId: number, input: PolishInput) => {
    aiLoading.value = true
    try {
      const response = await resumeApi.aiPolish(resumeId, input)
      const v = response.data.data
      if (v) {
        versions.value.unshift(v)
        currentVersion.value = v
      }
      message.success('AI 润色完成，已生成新版本')
      return v || null
    } catch (error) {
      console.error('AI 润色失败:', error)
      return null
    } finally {
      aiLoading.value = false
    }
  }

  // AI 评分
  const aiScore = async (resumeId: number, jd?: string) => {
    aiLoading.value = true
    try {
      const response = await resumeApi.aiScore(resumeId, jd)
      return response.data.data as ScoreResult | null
    } catch (error) {
      console.error('AI 评分失败:', error)
      return null
    } finally {
      aiLoading.value = false
    }
  }

  // AI JD 匹配
  const aiJdMatch = async (resumeId: number, jd: string) => {
    aiLoading.value = true
    try {
      const response = await resumeApi.aiJdMatch(resumeId, jd)
      return response.data.data as JDMatchResult | null
    } catch (error) {
      console.error('AI JD 匹配失败:', error)
      return null
    } finally {
      aiLoading.value = false
    }
  }

  // 同步档案
  const syncProfile = async (resumeId: number) => {
    try {
      await resumeApi.syncProfile(resumeId)
      message.success('已同步到档案')
      return true
    } catch (error) {
      console.error('同步档案失败:', error)
      return false
    }
  }

  // 清空流式文本
  const clearStreaming = () => {
    streamingText.value = ''
  }

  // 清空当前简历
  const clearCurrent = () => {
    currentResume.value = null
    versions.value = []
    currentVersion.value = null
    streamingText.value = ''
  }

  return {
    // 状态
    resumes,
    currentResume,
    versions,
    currentVersion,
    loading,
    aiLoading,
    streamingText,
    // 操作
    fetchList,
    fetchOne,
    create,
    update,
    remove,
    fetchVersions,
    createVersion,
    rollbackVersion,
    setCurrentVersion,
    aiGenerate,
    aiPolish,
    aiScore,
    aiJdMatch,
    syncProfile,
    clearStreaming,
    clearCurrent,
  }
})
