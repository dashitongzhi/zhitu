<template>
  <div class="admin-deliveries">
    <div class="search-bar">
      <a-form layout="inline" :model="searchForm">
        <a-form-item label="状态">
          <a-select
            v-model:value="searchForm.status"
            placeholder="全部状态"
            style="width: 140px"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option value="pending">待处理</a-select-option>
            <a-select-option value="written_test">笔试中</a-select-option>
            <a-select-option value="interview">面试中</a-select-option>
            <a-select-option value="waiting_offer">待Offer</a-select-option>
            <a-select-option value="offer">已获Offer</a-select-option>
            <a-select-option value="rejected">已拒绝</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="公司">
          <a-input
            v-model:value="searchForm.company"
            placeholder="公司名称"
            allow-clear
            style="width: 160px"
            @change="handleSearch"
          />
        </a-form-item>
        <a-form-item label="用户邮箱">
          <a-input
            v-model:value="searchForm.user_email"
            placeholder="用户邮箱"
            allow-clear
            style="width: 180px"
            @change="handleSearch"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            <SearchOutlined />
            搜索
          </a-button>
          <a-button style="margin-left: 8px" @click="handleReset">
            <ReloadOutlined />
            重置
          </a-button>
        </a-form-item>
      </a-form>
    </div>

    <a-table
      :columns="columns"
      :data-source="deliveryList"
      :loading="loading"
      :pagination="paginationConfig"
      row-key="id"
      :expandable="{ expandedRowRender: expandedRowRender }"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'apply_date'">
          {{ formatDate(record.apply_date) }}
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import * as adminApi from '@/api/admin'
import type { Delivery, DeliveryStatus, AdminDeliveryListParams } from '@/types/models'
import type { TableColumnsType, TablePaginationConfig } from 'ant-design-vue'

const loading = ref(false)
const deliveryList = ref<Delivery[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<AdminDeliveryListParams>({
  status: '',
  company: '',
  user_email: '',
  page: 1,
  page_size: 10,
})

const columns: TableColumnsType = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80,
  },
  {
    title: '用户邮箱',
    dataIndex: 'user_id',
    key: 'user_id',
    width: 180,
    customRender: ({ record }) => {
      return h('span', record.user_email || `用户#${record.user_id}`)
    },
  },
  {
    title: '公司',
    dataIndex: 'company',
    key: 'company',
  },
  {
    title: '职位',
    dataIndex: 'position',
    key: 'position',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
  },
  {
    title: '投递日期',
    dataIndex: 'apply_date',
    key: 'apply_date',
    width: 180,
  },
]

const paginationConfig = reactive<TablePaginationConfig>({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total) => `共 ${total} 条记录`,
})

const statusMap: Record<DeliveryStatus, { text: string; color: string }> = {
  pending: { text: '待处理', color: 'default' },
  written_test: { text: '笔试中', color: 'blue' },
  interview: { text: '面试中', color: 'orange' },
  waiting_offer: { text: '待Offer', color: 'gold' },
  offer: { text: '已获Offer', color: 'green' },
  rejected: { text: '已拒绝', color: 'red' },
}

const getStatusText = (status: DeliveryStatus): string => {
  return statusMap[status]?.text || status
}

const getStatusColor = (status: DeliveryStatus): string => {
  return statusMap[status]?.color || 'default'
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const expandedRowRender = (record: Delivery) => {
  const nextStep = record.next_step ? safeParseJson(record.next_step) : null
  const hrContact = record.hr_contact ? safeParseJson(record.hr_contact) : null
  return h('div', { class: 'expanded-detail' }, [
    h('a-descriptions', {
      column: 2,
      bordered: true,
      size: 'small',
    }, {
      default: () => [
        h('a-descriptions-item', { label: '投递ID' }, () => record.id),
        h('a-descriptions-item', { label: '用户ID' }, () => record.user_id),
        h('a-descriptions-item', { label: '公司' }, () => record.company),
        h('a-descriptions-item', { label: '职位' }, () => record.position),
        h('a-descriptions-item', { label: '状态' }, () => getStatusText(record.status)),
        h('a-descriptions-item', { label: '投递渠道' }, () => record.channel || '-'),
        h('a-descriptions-item', { label: '优先级' }, () => record.priority || '-'),
        h('a-descriptions-item', { label: '投递时间' }, () => formatDate(record.apply_date)),
        h('a-descriptions-item', { label: '更新时间' }, () => formatDate(record.updated_at)),
        h('a-descriptions-item', { label: 'HR联系方式' }, () => hrContact ? JSON.stringify(hrContact) : '-'),
        record.jd_text
          ? h('a-descriptions-item', { label: 'JD内容', span: 2 }, () => record.jd_text)
          : null,
        nextStep
          ? h('a-descriptions-item', { label: '下一步', span: 2 }, () => JSON.stringify(nextStep))
          : null,
        record.remark
          ? h('a-descriptions-item', { label: '备注', span: 2 }, () => record.remark)
          : null,
      ]
    }),
  ])
}

// 安全解析 JSON 字符串
const safeParseJson = (str: string): unknown => {
  try {
    return JSON.parse(str)
  } catch {
    return str
  }
}

const fetchDeliveries = async () => {
  loading.value = true
  try {
    const response = await adminApi.getDeliveries({
      status: searchForm.status || undefined,
      company: searchForm.company || undefined,
      user_email: searchForm.user_email || undefined,
      page: currentPage.value,
      page_size: pageSize.value,
    })
    deliveryList.value = response.data.data.list || []
    total.value = response.data.data.total
    paginationConfig.total = response.data.data.total
  } catch (error) {
    console.error('获取投递列表失败:', error)
    message.error('获取投递列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  paginationConfig.current = 1
  fetchDeliveries()
}

const handleReset = () => {
  searchForm.status = ''
  searchForm.company = ''
  searchForm.user_email = ''
  currentPage.value = 1
  paginationConfig.current = 1
  fetchDeliveries()
}

const handleTableChange = (pagination: TablePaginationConfig) => {
  currentPage.value = pagination.current || 1
  pageSize.value = pagination.pageSize || 10
  fetchDeliveries()
}

onMounted(() => {
  fetchDeliveries()
})
</script>

<style scoped>
.admin-deliveries {
  width: 100%;
}

.search-bar {
  margin-bottom: 16px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.expanded-detail {
  padding: 8px 0;
}

.rounds-section {
  margin-top: 16px;
}
</style>
