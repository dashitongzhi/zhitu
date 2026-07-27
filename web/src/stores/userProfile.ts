import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// 个人基础信息（个人资料弹窗维护，同步到简历实验室）
// 持久化到 localStorage，避免重复填写
export interface UserProfileBasic {
  name: string
  phone: string
  email: string
  city: string
  github: string
  gender: 'male' | 'female' | ''
  birth_date: string
  self_introduction: string
  avatar: string
}

const STORAGE_KEY = 'zhitu-user-profile'

const createDefault = (): UserProfileBasic => ({
  name: '',
  phone: '',
  email: '',
  city: '',
  github: '',
  gender: '',
  birth_date: '',
  self_introduction: '',
  avatar: '',
})

// 从 localStorage 读取初始值（避免 SSR / 首次访问空白）
const loadInitial = (): UserProfileBasic => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return createDefault()
    const parsed = JSON.parse(raw)
    return { ...createDefault(), ...parsed }
  } catch {
    return createDefault()
  }
}

export const useUserProfileStore = defineStore('userProfile', () => {
  const basic = ref<UserProfileBasic>(loadInitial())

  const hasFilled = computed(() =>
    !!(basic.value.name || basic.value.phone || basic.value.email)
  )

  // 更新基础信息（同时写入 localStorage）
  const updateBasic = (data: Partial<UserProfileBasic>) => {
    basic.value = { ...basic.value, ...data }
    persist()
  }

  // 单字段更新（适合表单 v-model 同步）
  const setField = <K extends keyof UserProfileBasic>(key: K, value: UserProfileBasic[K]) => {
    basic.value[key] = value
    persist()
  }

  const persist = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(basic.value))
    } catch (e) {
      console.warn('无法写入 localStorage:', e)
    }
  }

  return { basic, hasFilled, updateBasic, setField }
})
