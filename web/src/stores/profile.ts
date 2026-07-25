import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  FullProfile,
  UpdateProfileRequest,
  Education,
  Work,
  Project,
  Skill,
  Honor,
  Practice,
} from '@/types/models'
import * as profileApi from '@/api/profile'
import { message } from 'ant-design-vue'

export const useProfileStore = defineStore('profile', () => {
  const profile = ref<FullProfile | null>(null)
  const completion = ref(0)
  const loading = ref(false)

  // 获取完整档案
  const fetchProfile = async () => {
    loading.value = true
    try {
      const response = await profileApi.getProfile()
      profile.value = response.data.data
      if (profile.value) {
        completion.value = profile.value.completion_pct || 0
      }
      return true
    } catch (error) {
      console.error('获取档案失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 更新基础信息
  const updateProfile = async (data: UpdateProfileRequest) => {
    try {
      const response = await profileApi.updateProfile(data)
      profile.value = response.data.data
      message.success('档案保存成功')
      await fetchCompletion()
      return true
    } catch (error) {
      console.error('更新档案失败:', error)
      return false
    }
  }

  // 获取完成度
  const fetchCompletion = async () => {
    try {
      const response = await profileApi.getCompletion()
      completion.value = response.data.data?.completion_pct || 0
      return true
    } catch (error) {
      console.error('获取完成度失败:', error)
      return false
    }
  }

  // 上传简历解析
  const parseResume = async (file: File) => {
    try {
      await profileApi.parseResume(file)
      message.success('简历解析成功，已合并到档案')
      await fetchProfile()
      return true
    } catch (error) {
      console.error('简历解析失败:', error)
      return false
    }
  }

  // ==================== 子资源操作（操作后刷新档案）====================

  const refreshAfterSubChange = async () => {
    await fetchProfile()
    await fetchCompletion()
  }

  // 教育背景
  const createEducation = async (data: Omit<Education, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createEducation(data)
      message.success('教育背景已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updateEducation = async (id: number, data: Partial<Education>) => {
    try {
      await profileApi.updateEducation(id, data)
      message.success('教育背景已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deleteEducation = async (id: number) => {
    try {
      await profileApi.deleteEducation(id)
      message.success('教育背景已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  // 工作经历
  const createWork = async (data: Omit<Work, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createWork(data)
      message.success('工作经历已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updateWork = async (id: number, data: Partial<Work>) => {
    try {
      await profileApi.updateWork(id, data)
      message.success('工作经历已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deleteWork = async (id: number) => {
    try {
      await profileApi.deleteWork(id)
      message.success('工作经历已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  // 项目经历
  const createProject = async (data: Omit<Project, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createProject(data)
      message.success('项目经历已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updateProject = async (id: number, data: Partial<Project>) => {
    try {
      await profileApi.updateProject(id, data)
      message.success('项目经历已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deleteProject = async (id: number) => {
    try {
      await profileApi.deleteProject(id)
      message.success('项目经历已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  // 技能
  const createSkill = async (data: Omit<Skill, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createSkill(data)
      message.success('技能已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updateSkill = async (id: number, data: Partial<Skill>) => {
    try {
      await profileApi.updateSkill(id, data)
      message.success('技能已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deleteSkill = async (id: number) => {
    try {
      await profileApi.deleteSkill(id)
      message.success('技能已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  // 荣誉奖项
  const createHonor = async (data: Omit<Honor, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createHonor(data)
      message.success('荣誉奖项已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updateHonor = async (id: number, data: Partial<Honor>) => {
    try {
      await profileApi.updateHonor(id, data)
      message.success('荣誉奖项已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deleteHonor = async (id: number) => {
    try {
      await profileApi.deleteHonor(id)
      message.success('荣誉奖项已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  // 校内外实践
  const createPractice = async (data: Omit<Practice, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    try {
      await profileApi.createPractice(data)
      message.success('实践经历已添加')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const updatePractice = async (id: number, data: Partial<Practice>) => {
    try {
      await profileApi.updatePractice(id, data)
      message.success('实践经历已更新')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }
  const deletePractice = async (id: number) => {
    try {
      await profileApi.deletePractice(id)
      message.success('实践经历已删除')
      await refreshAfterSubChange()
      return true
    } catch {
      return false
    }
  }

  return {
    profile,
    completion,
    loading,
    fetchProfile,
    updateProfile,
    fetchCompletion,
    parseResume,
    // 教育背景
    createEducation,
    updateEducation,
    deleteEducation,
    // 工作经历
    createWork,
    updateWork,
    deleteWork,
    // 项目经历
    createProject,
    updateProject,
    deleteProject,
    // 技能
    createSkill,
    updateSkill,
    deleteSkill,
    // 荣誉奖项
    createHonor,
    updateHonor,
    deleteHonor,
    // 校内外实践
    createPractice,
    updatePractice,
    deletePractice,
  }
})
