<template>
  <div class="delivery-kanban">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <a-radio-group v-model:value="viewMode" button-style="solid">
          <a-radio-button value="table">
            <TableOutlined /> 表格视图
          </a-radio-button>
          <a-radio-button value="kanban">
            <AppstoreOutlined /> 看板视图
          </a-radio-button>
        </a-radio-group>
      </div>

      <div class="toolbar-right">
        <a-space>
          <a-select
            v-model:value="filters.status"
            placeholder="状态筛选"
            style="width: 140px"
            allowClear
            @change="handleFilter"
          >
            <a-select-option v-for="s in statusList" :key="s.value" :value="s.value">
              {{ s.label }}
            </a-select-option>
          </a-select>

          <a-select
            v-model:value="filters.channel"
            placeholder="渠道筛选"
            style="width: 140px"
            allowClear
            @change="handleFilter"
          >
            <a-select-option v-for="c in channelOptions" :key="c.value" :value="c.value">
              {{ c.label }}
            </a-select-option>
          </a-select>

          <a-button type="primary" @click="showCreateModal">
            <PlusOutlined /> 新增投递
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 表格视图 -->
    <div v-if="viewMode === 'table'" class="table-view">
      <a-table
        :columns="columns"
        :data-source="deliveryStore.deliveries"
        :loading="deliveryStore.loading"
        :pagination="false"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'company'">
            <div class="company-cell">
              <div class="company-name">{{ record.company }}</div>
              <div class="position-name">{{ record.position }}</div>
            </div>
          </template>

          <template v-else-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>

          <template v-else-if="column.key === 'channel'">
            {{ getChannelLabel(record.channel) }}
          </template>

          <template v-else-if="column.key === 'priority'">
            <a-tag :color="getPriorityColor(record.priority)">{{ getPriorityText(record.priority) }}</a-tag>
          </template>

          <template v-else-if="column.key === 'apply_date'">
            {{ formatDate(record.apply_date) }}
          </template>

          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="showDetailDrawer(record)">
                详情
              </a-button>
              <a-dropdown>
                <a-button type="link" size="small">
                  更多 <DownOutlined />
                </a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item
                      v-for="s in getAvailableTransitions(record.status)"
                      :key="s.value"
                      @click="handleStatusChange(record, s.value)"
                    >
                      转为{{ s.label }}
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="handleDelete(record.id)">
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>

    <!-- 看板视图 -->
    <div v-else class="kanban-view">
      <a-row :gutter="12">
        <a-col :span="4" v-for="status in statusList" :key="status.value">
          <div class="kanban-column">
            <div class="kanban-column-header" :style="{ background: status.color }">
              {{ status.label }} ({{ getDeliveriesByStatus(status.value).length }})
            </div>
            <div class="kanban-column-body">
              <div
                v-for="delivery in getDeliveriesByStatus(status.value)"
                :key="delivery.id"
                class="kanban-card"
                @click="showDetailDrawer(delivery)"
              >
                <div class="card-header">
                  <span class="company">{{ delivery.company }}</span>
                  <span class="position">{{ delivery.position }}</span>
                </div>
                <div class="card-footer">
                  <span class="date">{{ formatDate(delivery.apply_date) }}</span>
                  <a-tag v-if="delivery.channel" size="small">{{ getChannelLabel(delivery.channel) }}</a-tag>
                </div>
              </div>
            </div>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 新增投递模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="新增投递"
      width="640px"
      @ok="handleCreate"
      @cancel="resetCreateForm"
    >
      <a-form
        :model="createForm"
        :rules="createRules"
        layout="vertical"
        ref="createFormRef"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公司" name="company">
              <a-input v-model:value="createForm.company" placeholder="请输入公司名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="职位" name="position">
              <a-input v-model:value="createForm.position" placeholder="请输入职位名称" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="投递渠道" name="channel">
              <a-select v-model:value="createForm.channel" placeholder="请选择投递渠道">
                <a-select-option v-for="c in channelOptions" :key="c.value" :value="c.value">
                  {{ c.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="投递日期" name="apply_date">
              <a-date-picker
                v-model:value="createForm.apply_date"
                placeholder="选择日期"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="优先级" name="priority">
              <a-select v-model:value="createForm.priority" placeholder="请选择优先级">
                <a-select-option value="high">高</a-select-option>
                <a-select-option value="medium">中</a-select-option>
                <a-select-option value="low">低</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="简历版本ID" name="resume_version_id">
              <a-input-number
                v-model:value="createForm.resume_version_id"
                placeholder="可选"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="JD 内容" name="jd_text">
          <a-textarea
            v-model:value="createForm.jd_text"
            placeholder="职位 JD 描述（可选）"
            :rows="3"
          />
        </a-form-item>

        <a-form-item label="备注" name="remark">
          <a-textarea
            v-model:value="createForm.remark"
            placeholder="备注信息（可选）"
            :rows="2"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情抽屉 -->
    <a-drawer
      v-model:open="detailDrawerVisible"
      title="投递详情"
      width="720"
      @close="handleDrawerClose"
    >
      <template v-if="deliveryStore.currentDelivery">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="公司">{{ deliveryStore.currentDelivery.company }}</a-descriptions-item>
          <a-descriptions-item label="职位">{{ deliveryStore.currentDelivery.position }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(deliveryStore.currentDelivery.status)">
              {{ getStatusText(deliveryStore.currentDelivery.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="投递渠道">{{ getChannelLabel(deliveryStore.currentDelivery.channel) }}</a-descriptions-item>
          <a-descriptions-item label="投递日期">{{ formatDate(deliveryStore.currentDelivery.apply_date) }}</a-descriptions-item>
          <a-descriptions-item label="优先级">
            <a-tag :color="getPriorityColor(deliveryStore.currentDelivery.priority)">{{ getPriorityText(deliveryStore.currentDelivery.priority) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="JD 内容" :span="2">{{ deliveryStore.currentDelivery.jd_text || '-' }}</a-descriptions-item>
          <a-descriptions-item label="HR 联系方式" :span="2">{{ formatJsonField(deliveryStore.currentDelivery.hr_contact) }}</a-descriptions-item>
          <a-descriptions-item label="下一步" :span="2">{{ formatJsonField(deliveryStore.currentDelivery.next_step) }}</a-descriptions-item>
          <a-descriptions-item label="备注" :span="2">{{ deliveryStore.currentDelivery.remark || '-' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 面试轮次 -->
        <a-divider orientation="left">
          面试轮次
          <a-button type="link" size="small" @click="showRoundModal()">
            <PlusOutlined /> 新增轮次
          </a-button>
        </a-divider>
        <a-timeline v-if="deliveryStore.rounds.length > 0">
          <a-timeline-item
            v-for="round in deliveryStore.rounds"
            :key="round.id"
            :color="getRoundResultColor(round.result)"
          >
            <div class="round-item">
              <div class="round-header">
                <span class="round-type">{{ getRoundTypeText(round.round_type) }}</span>
                <a-tag :color="getRoundResultColor(round.result)">
                  {{ getRoundResultText(round.result) }}
                </a-tag>
                <a-dropdown>
                  <a-button type="link" size="small">操作 <DownOutlined /></a-button>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item @click="showRoundModal(round)">编辑</a-menu-item>
                      <a-menu-item @click="handleDeleteRound(round.id)">删除</a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </div>
              <div class="round-info">
                <div>时间: {{ round.interview_time ? formatDateTime(round.interview_time) : '-' }}</div>
                <div>形式: {{ getRoundFormatText(round.format) }}</div>
                <div v-if="round.interviewer_name">面试官: {{ round.interviewer_name }} {{ round.interviewer_title }}</div>
                <div v-if="round.question_summary">问题摘要: {{ round.question_summary }}</div>
                <div v-if="round.feedback">反馈: {{ round.feedback }}</div>
              </div>
            </div>
          </a-timeline-item>
        </a-timeline>
        <a-empty v-else description="暂无面试轮次" />

        <!-- HR 反馈 -->
        <a-divider orientation="left">
          HR 反馈
          <a-button type="link" size="small" @click="showFeedbackModal()">
            <PlusOutlined /> 新增反馈
          </a-button>
        </a-divider>
        <a-list
          v-if="deliveryStore.feedbacks.length > 0"
          :data-source="deliveryStore.feedbacks"
          size="small"
          bordered
        >
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta>
                <template #title>
                  {{ formatDateTime(item.contact_time) }} · {{ item.method }}
                </template>
                <template #description>
                  <div>{{ item.summary }}</div>
                  <div v-if="item.next_action" class="feedback-next">下一步: {{ item.next_action }}</div>
                </template>
              </a-list-item-meta>
              <template #actions>
                <a-button type="link" size="small" danger @click="handleDeleteFeedback(item.id)">删除</a-button>
              </template>
            </a-list-item>
          </template>
        </a-list>
        <a-empty v-else description="暂无 HR 反馈" />
      </template>
    </a-drawer>

    <!-- 轮次表单弹窗 -->
    <a-modal
      v-model:open="roundModalVisible"
      :title="editingRound ? '编辑轮次' : '新增轮次'"
      width="560px"
      @ok="handleSaveRound"
    >
      <a-form :model="roundForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="轮次类型" name="round_type">
              <a-select v-model:value="roundForm.round_type" placeholder="请选择">
                <a-select-option value="written_test">笔试</a-select-option>
                <a-select-option value="first_tech">一面</a-select-option>
                <a-select-option value="second_tech">二面</a-select-option>
                <a-select-option value="third_tech">三面</a-select-option>
                <a-select-option value="cross">交叉面</a-select-option>
                <a-select-option value="hr">HR 面</a-select-option>
                <a-select-option value="additional">加面</a-select-option>
                <a-select-option value="final">终面</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="面试形式" name="format">
              <a-select v-model:value="roundForm.format" placeholder="请选择">
                <a-select-option value="onsite">现场</a-select-option>
                <a-select-option value="video">视频</a-select-option>
                <a-select-option value="phone">电话</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="面试时间" name="interview_time">
              <a-date-picker
                v-model:value="roundForm.interview_time"
                show-time
                placeholder="选择时间"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结果" name="result">
              <a-select v-model:value="roundForm.result" placeholder="请选择">
                <a-select-option value="pass">通过</a-select-option>
                <a-select-option value="pending">待定</a-select-option>
                <a-select-option value="rejected">未通过</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="面试官姓名" name="interviewer_name">
              <a-input v-model:value="roundForm.interviewer_name" placeholder="可选" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="面试官头衔" name="interviewer_title">
              <a-input v-model:value="roundForm.interviewer_title" placeholder="可选" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="问题摘要" name="question_summary">
          <a-textarea v-model:value="roundForm.question_summary" :rows="2" placeholder="可选" />
        </a-form-item>
        <a-form-item label="反馈" name="feedback">
          <a-textarea v-model:value="roundForm.feedback" :rows="2" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 反馈表单弹窗 -->
    <a-modal
      v-model:open="feedbackModalVisible"
      title="新增反馈"
      width="480px"
      @ok="handleSaveFeedback"
    >
      <a-form :model="feedbackForm" layout="vertical">
        <a-form-item label="联系时间" name="contact_time" required>
          <a-date-picker
            v-model:value="feedbackForm.contact_time"
            show-time
            placeholder="选择时间"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="联系方式" name="method" required>
          <a-select v-model:value="feedbackForm.method" placeholder="请选择">
            <a-select-option value="wechat">微信</a-select-option>
            <a-select-option value="phone">电话</a-select-option>
            <a-select-option value="email">邮件</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="反馈摘要" name="summary" required>
          <a-textarea v-model:value="feedbackForm.summary" :rows="3" placeholder="HR 反馈内容" />
        </a-form-item>
        <a-form-item label="下一步行动" name="next_action">
          <a-textarea v-model:value="feedbackForm.next_action" :rows="2" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useDeliveryStore } from '@/stores/delivery'
import { message, Modal } from 'ant-design-vue'
import {
  TableOutlined,
  AppstoreOutlined,
  PlusOutlined,
  DownOutlined,
} from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'
import type {
  Delivery,
  DeliveryStatus,
  DeliveryRound,
  CreateDeliveryRequest,
  CreateRoundRequest,
  CreateFeedbackRequest,
} from '@/types/models'
import dayjs, { Dayjs } from 'dayjs'

const deliveryStore = useDeliveryStore()

// 视图模式
const viewMode = ref<'table' | 'kanban'>('table')

// 状态枚举（6 列，对齐后端）
const statusList = [
  { label: '待处理', value: 'pending' as DeliveryStatus, color: '#d9d9d9' },
  { label: '笔试中', value: 'written_test' as DeliveryStatus, color: '#1890ff' },
  { label: '面试中', value: 'interview' as DeliveryStatus, color: '#52c41a' },
  { label: '待Offer', value: 'waiting_offer' as DeliveryStatus, color: '#faad14' },
  { label: '已获Offer', value: 'offer' as DeliveryStatus, color: '#722ed1' },
  { label: '已拒绝', value: 'rejected' as DeliveryStatus, color: '#f5222d' },
]

// 渠道枚举（英文 value，中文 label）
const channelOptions = [
  { value: 'boss', label: 'BOSS直聘' },
  { value: 'official', label: '官网' },
  { value: 'referral', label: '内推' },
  { value: 'campus', label: '校园招聘' },
  { value: 'headhunt', label: '猎头' },
  { value: 'other', label: '其他' },
]

// 筛选条件
const filters = reactive({
  status: '' as string,
  channel: '' as string,
})

// 表格列定义
const columns = [
  { title: '公司/职位', key: 'company', dataIndex: 'company' },
  { title: '状态', key: 'status', dataIndex: 'status', width: 100 },
  { title: '投递渠道', key: 'channel', dataIndex: 'channel', width: 110 },
  { title: '优先级', key: 'priority', dataIndex: 'priority', width: 90 },
  { title: '投递日期', key: 'apply_date', dataIndex: 'apply_date', width: 120 },
  { title: '操作', key: 'action', width: 160 },
]

// 状态合法流转
const transitionMap: Record<DeliveryStatus, DeliveryStatus[]> = {
  pending: ['written_test', 'interview', 'rejected'],
  written_test: ['interview', 'waiting_offer', 'rejected'],
  interview: ['waiting_offer', 'offer', 'rejected'],
  waiting_offer: ['offer', 'rejected'],
  offer: [],
  rejected: [],
}

const getAvailableTransitions = (status: DeliveryStatus) => {
  return transitionMap[status]?.map((v) => ({
    value: v,
    label: getStatusText(v),
  })) || []
}

// 新增投递模态框
const createModalVisible = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = reactive<CreateDeliveryRequest & { apply_date?: Dayjs }>({
  company: '',
  position: '',
  channel: 'boss',
  apply_date: undefined,
  priority: 'medium',
  jd_text: '',
  remark: '',
  resume_version_id: undefined,
})

const createRules = {
  company: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  position: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  channel: [{ required: true, message: '请选择投递渠道', trigger: 'change' }],
  apply_date: [{ required: true, message: '请选择投递日期', trigger: 'change' }],
}

// 详情抽屉
const detailDrawerVisible = ref(false)

// 轮次表单
const roundModalVisible = ref(false)
const editingRound = ref<DeliveryRound | null>(null)
const roundForm = reactive<{
  round_type: string
  interview_time: Dayjs | null
  format: string
  interviewer_name: string
  interviewer_title: string
  question_summary: string
  feedback: string
  result: string
}>({
  round_type: 'first_tech',
  interview_time: null,
  format: 'video',
  interviewer_name: '',
  interviewer_title: '',
  question_summary: '',
  feedback: '',
  result: 'pending',
})

// 反馈表单
const feedbackModalVisible = ref(false)
const feedbackForm = reactive<{
  contact_time: Dayjs | null
  method: string
  summary: string
  next_action: string
}>({
  contact_time: null,
  method: 'wechat',
  summary: '',
  next_action: '',
})

// 页面加载时获取数据
onMounted(async () => {
  await deliveryStore.fetchDeliveries()
})

// 处理筛选
const handleFilter = async () => {
  await deliveryStore.fetchDeliveries({
    status: filters.status || undefined,
    channel: filters.channel || undefined,
  })
}

// 显示新增模态框
const showCreateModal = () => {
  createModalVisible.value = true
}

// 处理新增
const handleCreate = async () => {
  try {
    await createFormRef.value?.validateFields()
    const data: CreateDeliveryRequest = {
      company: createForm.company,
      position: createForm.position,
      channel: createForm.channel as CreateDeliveryRequest['channel'],
      apply_date: createForm.apply_date?.format('YYYY-MM-DD') || '',
      priority: createForm.priority as CreateDeliveryRequest['priority'],
      jd_text: createForm.jd_text,
      remark: createForm.remark,
      resume_version_id: createForm.resume_version_id,
    }
    await deliveryStore.createDelivery(data)
    createModalVisible.value = false
    resetCreateForm()
  } catch (error) {
    console.error('验证失败:', error)
  }
}

// 重置新增表单
const resetCreateForm = () => {
  createFormRef.value?.resetFields()
  Object.assign(createForm, {
    company: '',
    position: '',
    channel: 'boss',
    apply_date: undefined,
    priority: 'medium',
    jd_text: '',
    remark: '',
    resume_version_id: undefined,
  })
}

// 显示详情抽屉
const showDetailDrawer = async (delivery: Delivery) => {
  detailDrawerVisible.value = true
  await deliveryStore.fetchDelivery(delivery.id)
}

// 关闭抽屉
const handleDrawerClose = () => {
  deliveryStore.clearCurrentDelivery()
}

// 处理状态变更
const handleStatusChange = async (delivery: Delivery, status: DeliveryStatus) => {
  await deliveryStore.changeStatus(delivery.id, { status })
}

// 处理删除
const handleDelete = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这条投递记录吗？',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      await deliveryStore.deleteDelivery(id)
    },
  })
}

// ==================== 轮次操作 ====================

const showRoundModal = (round?: DeliveryRound) => {
  if (round) {
    editingRound.value = round
    Object.assign(roundForm, {
      round_type: round.round_type,
      interview_time: round.interview_time ? dayjs(round.interview_time) : null,
      format: round.format,
      interviewer_name: round.interviewer_name,
      interviewer_title: round.interviewer_title,
      question_summary: round.question_summary,
      feedback: round.feedback,
      result: round.result,
    })
  } else {
    editingRound.value = null
    Object.assign(roundForm, {
      round_type: 'first_tech',
      interview_time: null,
      format: 'video',
      interviewer_name: '',
      interviewer_title: '',
      question_summary: '',
      feedback: '',
      result: 'pending',
    })
  }
  roundModalVisible.value = true
}

const handleSaveRound = async () => {
  if (!deliveryStore.currentDelivery) return
  const deliveryId = deliveryStore.currentDelivery.id
  const data: CreateRoundRequest = {
    round_type: roundForm.round_type,
    interview_time: roundForm.interview_time?.format('YYYY-MM-DD HH:mm:ss') || '',
    format: roundForm.format,
    interviewer_name: roundForm.interviewer_name,
    interviewer_title: roundForm.interviewer_title,
    question_summary: roundForm.question_summary,
    feedback: roundForm.feedback,
    result: roundForm.result,
  }
  if (editingRound.value) {
    await deliveryStore.updateRound(deliveryId, editingRound.value.id, data)
  } else {
    await deliveryStore.createRound(deliveryId, data)
  }
  roundModalVisible.value = false
}

const handleDeleteRound = (roundId: number) => {
  if (!deliveryStore.currentDelivery) return
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该面试轮次吗？',
    onOk: async () => {
      await deliveryStore.deleteRound(deliveryStore.currentDelivery!.id, roundId)
    },
  })
}

// ==================== 反馈操作 ====================

const showFeedbackModal = () => {
  Object.assign(feedbackForm, {
    contact_time: dayjs(),
    method: 'wechat',
    summary: '',
    next_action: '',
  })
  feedbackModalVisible.value = true
}

const handleSaveFeedback = async () => {
  if (!deliveryStore.currentDelivery) return
  if (!feedbackForm.contact_time || !feedbackForm.summary) {
    message.warning('请填写联系时间和反馈摘要')
    return
  }
  const data: CreateFeedbackRequest = {
    contact_time: feedbackForm.contact_time.format('YYYY-MM-DD HH:mm'),
    method: feedbackForm.method,
    summary: feedbackForm.summary,
    next_action: feedbackForm.next_action,
  }
  await deliveryStore.createFeedback(deliveryStore.currentDelivery.id, data)
  feedbackModalVisible.value = false
}

const handleDeleteFeedback = (feedbackId: number) => {
  if (!deliveryStore.currentDelivery) return
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该反馈吗？',
    onOk: async () => {
      await deliveryStore.deleteFeedback(deliveryStore.currentDelivery!.id, feedbackId)
    },
  })
}

// ==================== 辅助函数 ====================

const getDeliveriesByStatus = (status: string) => {
  return deliveryStore.deliveries.filter((d) => d.status === status)
}

const getStatusColor = (status: string): string => {
  const found = statusList.find((s) => s.value === status)
  return found?.color || 'default'
}

const getStatusText = (status: string): string => {
  const found = statusList.find((s) => s.value === status)
  return found?.label || status
}

const getChannelLabel = (channel: string): string => {
  return channelOptions.find((c) => c.value === channel)?.label || channel || '-'
}

const getPriorityColor = (priority: string): string => {
  const map: Record<string, string> = { high: 'red', medium: 'orange', low: 'default' }
  return map[priority] || 'default'
}

const getPriorityText = (priority: string): string => {
  const map: Record<string, string> = { high: '高', medium: '中', low: '低' }
  return map[priority] || priority || '-'
}

const getRoundResultColor = (result?: string): string => {
  if (!result) return 'default'
  return result === 'pass' ? 'green' : result === 'rejected' ? 'red' : 'blue'
}

const getRoundResultText = (result?: string): string => {
  if (!result) return '待定'
  const map: Record<string, string> = { pass: '通过', pending: '待定', rejected: '未通过' }
  return map[result] || result
}

const getRoundTypeText = (type?: string): string => {
  const map: Record<string, string> = {
    written_test: '笔试',
    first_tech: '一面',
    second_tech: '二面',
    third_tech: '三面',
    cross: '交叉面',
    hr: 'HR 面',
    additional: '加面',
    final: '终面',
  }
  return map[type] || type || '未知'
}

const getRoundFormatText = (format?: string): string => {
  const map: Record<string, string> = { onsite: '现场', video: '视频', phone: '电话' }
  return map[format] || format || '-'
}

const formatJsonField = (field: string): string => {
  if (!field) return '-'
  try {
    return JSON.stringify(JSON.parse(field))
  } catch {
    return field
  }
}

const formatDate = (date: string): string => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string): string => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped>
.delivery-kanban {
  padding: 0;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.toolbar-left {
  display: flex;
  align-items: center;
}

.toolbar-right {
  display: flex;
  align-items: center;
}

.company-cell {
  line-height: 1.5;
}

.company-name {
  font-weight: 500;
  color: #262626;
}

.position-name {
  font-size: 12px;
  color: #8c8c8c;
}

/* 看板视图样式 */
.kanban-column {
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
}

.kanban-column-header {
  padding: 12px 16px;
  color: #fff;
  font-weight: 500;
  font-size: 14px;
}

.kanban-column-body {
  padding: 8px;
  max-height: calc(100vh - 280px);
  overflow-y: auto;
}

.kanban-card {
  background: #fff;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #e8e8e8;
}

.kanban-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-color: #1890ff;
}

.card-header {
  margin-bottom: 8px;
}

.card-header .company {
  font-weight: 500;
  color: #262626;
  display: block;
}

.card-header .position {
  font-size: 12px;
  color: #8c8c8c;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-footer .date {
  font-size: 12px;
  color: #999;
}

/* 详情抽屉样式 */
.round-item {
  padding: 8px 0;
}

.round-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  gap: 8px;
}

.round-type {
  font-weight: 500;
}

.round-info {
  font-size: 13px;
  color: #8c8c8c;
  line-height: 1.8;
}

.feedback-next {
  margin-top: 4px;
  color: #1890ff;
}
</style>
