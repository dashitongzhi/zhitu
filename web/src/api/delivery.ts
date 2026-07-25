import { get, post, put, del, patch } from '@/utils/request'
import type {
  ApiResponse,
  Delivery,
  DeliveryDetail,
  DeliveryListParams,
  CreateDeliveryRequest,
  ChangeStatusRequest,
  DeliveryRound,
  CreateRoundRequest,
  DeliveryFeedback,
  CreateFeedbackRequest,
  DeliveryStats,
  DeliveryFunnel,
} from '@/types/models'

// ==================== 投递主表 ====================

// 获取投递列表（无分页，返回数组）
export const list = (params?: DeliveryListParams) => {
  return get<ApiResponse<Delivery[]>>('/api/v1/deliveries', { params })
}

// 兼容旧名称
export const listDeliveries = list

// 获取投递详情（含轮次与反馈）
export const getDelivery = (id: number) => {
  return get<ApiResponse<DeliveryDetail>>(`/api/v1/deliveries/${id}`)
}

// 创建投递
export const create = (data: CreateDeliveryRequest) => {
  return post<ApiResponse<Delivery>>('/api/v1/deliveries', data)
}

// 兼容旧名称
export const createDelivery = create

// 更新投递
export const update = (id: number, data: Partial<Delivery>) => {
  return put<ApiResponse>(`/api/v1/deliveries/${id}`, data)
}

// 删除投递
export const deleteDelivery = (id: number) => {
  return del<ApiResponse>(`/api/v1/deliveries/${id}`)
}

// 变更投递状态
export const changeStatus = (id: number, data: ChangeStatusRequest) => {
  return patch<ApiResponse<Delivery>>(`/api/v1/deliveries/${id}/status`, data)
}

// ==================== 面试轮次 ====================

export const listRounds = (deliveryId: number) => {
  return get<ApiResponse<DeliveryRound[]>>(
    `/api/v1/deliveries/${deliveryId}/rounds`
  )
}

export const createRound = (deliveryId: number, data: CreateRoundRequest) => {
  return post<ApiResponse<DeliveryRound>>(
    `/api/v1/deliveries/${deliveryId}/rounds`,
    data
  )
}

export const updateRound = (
  deliveryId: number,
  roundId: number,
  data: Partial<DeliveryRound>
) => {
  return put<ApiResponse>(
    `/api/v1/deliveries/${deliveryId}/rounds/${roundId}`,
    data
  )
}

export const deleteRound = (deliveryId: number, roundId: number) => {
  return del<ApiResponse>(
    `/api/v1/deliveries/${deliveryId}/rounds/${roundId}`
  )
}

// ==================== HR 反馈 ====================

export const listFeedbacks = (deliveryId: number) => {
  return get<ApiResponse<DeliveryFeedback[]>>(
    `/api/v1/deliveries/${deliveryId}/feedbacks`
  )
}

export const createFeedback = (
  deliveryId: number,
  data: CreateFeedbackRequest
) => {
  return post<ApiResponse<DeliveryFeedback>>(
    `/api/v1/deliveries/${deliveryId}/feedbacks`,
    data
  )
}

export const deleteFeedback = (deliveryId: number, feedbackId: number) => {
  return del<ApiResponse>(
    `/api/v1/deliveries/${deliveryId}/feedbacks/${feedbackId}`
  )
}

// ==================== 统计与漏斗 ====================

export const getStats = () => {
  return get<ApiResponse<DeliveryStats>>('/api/v1/deliveries/stats')
}

export const getFunnel = () => {
  return get<ApiResponse<DeliveryFunnel>>('/api/v1/deliveries/funnel')
}
