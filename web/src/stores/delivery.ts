import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Delivery,
  DeliveryRound,
  DeliveryFeedback,
  DeliveryStats,
  DeliveryFunnel,
  CreateDeliveryRequest,
  ChangeStatusRequest,
  CreateRoundRequest,
  CreateFeedbackRequest,
  DeliveryListParams,
} from '@/types/models'
import * as deliveryApi from '@/api/delivery'
import { message } from 'ant-design-vue'

export const useDeliveryStore = defineStore('delivery', () => {
  // 状态
  const deliveries = ref<Delivery[]>([])
  const currentDelivery = ref<Delivery | null>(null)
  const rounds = ref<DeliveryRound[]>([])
  const feedbacks = ref<DeliveryFeedback[]>([])
  const stats = ref<DeliveryStats | null>(null)
  const funnel = ref<DeliveryFunnel | null>(null)
  const loading = ref(false)

  // 获取投递列表（后端无分页，返回数组）
  const fetchDeliveries = async (params?: DeliveryListParams) => {
    loading.value = true
    try {
      const response = await deliveryApi.list(params)
      deliveries.value = response.data.data || []
      return true
    } catch (error) {
      console.error('获取投递列表失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取单个投递详情（含轮次与反馈）
  const fetchDelivery = async (id: number) => {
    loading.value = true
    try {
      const response = await deliveryApi.getDelivery(id)
      const detail = response.data.data
      if (detail) {
        currentDelivery.value = detail.delivery
        rounds.value = detail.rounds || []
        feedbacks.value = detail.feedbacks || []
      }
      return true
    } catch (error) {
      console.error('获取投递详情失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 创建投递
  const createDelivery = async (data: CreateDeliveryRequest) => {
    try {
      const response = await deliveryApi.create(data)
      const d = response.data.data
      if (d) deliveries.value.unshift(d)
      message.success('投递创建成功')
      return d || null
    } catch (error) {
      console.error('创建投递失败:', error)
      return null
    }
  }

  // 更新投递
  const updateDelivery = async (id: number, data: Partial<Delivery>) => {
    try {
      await deliveryApi.update(id, data)
      // 更新成功后刷新列表中对应项
      const index = deliveries.value.findIndex((d) => d.id === id)
      if (index !== -1) {
        await fetchDelivery(id)
        if (currentDelivery.value) {
          deliveries.value[index] = currentDelivery.value
        }
      }
      message.success('投递更新成功')
      return true
    } catch (error) {
      console.error('更新投递失败:', error)
      return false
    }
  }

  // 删除投递
  const deleteDelivery = async (id: number) => {
    try {
      await deliveryApi.deleteDelivery(id)
      deliveries.value = deliveries.value.filter((d) => d.id !== id)
      if (currentDelivery.value?.id === id) {
        clearCurrentDelivery()
      }
      message.success('投递删除成功')
      return true
    } catch (error) {
      console.error('删除投递失败:', error)
      return false
    }
  }

  // 变更投递状态
  const changeStatus = async (id: number, data: ChangeStatusRequest) => {
    try {
      const response = await deliveryApi.changeStatus(id, data)
      const d = response.data.data
      const index = deliveries.value.findIndex((item) => item.id === id)
      if (index !== -1 && d) {
        deliveries.value[index] = d
      }
      if (currentDelivery.value?.id === id && d) {
        currentDelivery.value = d
      }
      message.success('状态变更成功')
      return true
    } catch (error) {
      console.error('变更状态失败:', error)
      return false
    }
  }

  // 获取轮次列表
  const fetchRounds = async (deliveryId: number) => {
    try {
      const response = await deliveryApi.listRounds(deliveryId)
      rounds.value = response.data.data || []
      return true
    } catch (error) {
      console.error('获取轮次列表失败:', error)
      return false
    }
  }

  // 创建轮次
  const createRound = async (deliveryId: number, data: CreateRoundRequest) => {
    try {
      const response = await deliveryApi.createRound(deliveryId, data)
      const r = response.data.data
      if (r) rounds.value.push(r)
      message.success('轮次创建成功')
      return r || null
    } catch (error) {
      console.error('创建轮次失败:', error)
      return null
    }
  }

  // 更新轮次
  const updateRound = async (
    deliveryId: number,
    roundId: number,
    data: Partial<DeliveryRound>
  ) => {
    try {
      await deliveryApi.updateRound(deliveryId, roundId, data)
      // 更新成功后刷新轮次列表
      await fetchRounds(deliveryId)
      message.success('轮次更新成功')
      return true
    } catch (error) {
      console.error('更新轮次失败:', error)
      return false
    }
  }

  // 删除轮次
  const deleteRound = async (deliveryId: number, roundId: number) => {
    try {
      await deliveryApi.deleteRound(deliveryId, roundId)
      rounds.value = rounds.value.filter((r) => r.id !== roundId)
      message.success('轮次删除成功')
      return true
    } catch (error) {
      console.error('删除轮次失败:', error)
      return false
    }
  }

  // 获取反馈列表
  const fetchFeedbacks = async (deliveryId: number) => {
    try {
      const response = await deliveryApi.listFeedbacks(deliveryId)
      feedbacks.value = response.data.data || []
      return true
    } catch (error) {
      console.error('获取反馈列表失败:', error)
      return false
    }
  }

  // 创建反馈
  const createFeedback = async (
    deliveryId: number,
    data: CreateFeedbackRequest
  ) => {
    try {
      const response = await deliveryApi.createFeedback(deliveryId, data)
      const f = response.data.data
      if (f) feedbacks.value.unshift(f)
      message.success('反馈创建成功')
      return f || null
    } catch (error) {
      console.error('创建反馈失败:', error)
      return null
    }
  }

  // 删除反馈
  const deleteFeedback = async (deliveryId: number, feedbackId: number) => {
    try {
      await deliveryApi.deleteFeedback(deliveryId, feedbackId)
      feedbacks.value = feedbacks.value.filter((f) => f.id !== feedbackId)
      message.success('反馈删除成功')
      return true
    } catch (error) {
      console.error('删除反馈失败:', error)
      return false
    }
  }

  // 获取统计数据
  const fetchStats = async () => {
    try {
      const response = await deliveryApi.getStats()
      stats.value = response.data.data
      return true
    } catch (error) {
      console.error('获取统计数据失败:', error)
      return false
    }
  }

  // 获取漏斗数据
  const fetchFunnel = async () => {
    try {
      const response = await deliveryApi.getFunnel()
      funnel.value = response.data.data
      return true
    } catch (error) {
      console.error('获取漏斗数据失败:', error)
      return false
    }
  }

  // 清空当前投递详情
  const clearCurrentDelivery = () => {
    currentDelivery.value = null
    rounds.value = []
    feedbacks.value = []
  }

  return {
    // 状态
    deliveries,
    currentDelivery,
    rounds,
    feedbacks,
    stats,
    funnel,
    loading,
    // 操作
    fetchDeliveries,
    fetchDelivery,
    createDelivery,
    updateDelivery,
    deleteDelivery,
    changeStatus,
    fetchRounds,
    createRound,
    updateRound,
    deleteRound,
    fetchFeedbacks,
    createFeedback,
    deleteFeedback,
    fetchStats,
    fetchFunnel,
    clearCurrentDelivery,
  }
})