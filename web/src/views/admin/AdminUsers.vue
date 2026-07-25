<template>
  <div class="admin-users">
    <div class="search-bar">
      <a-form layout="inline" :model="searchForm">
        <a-form-item label="关键词">
          <a-input
            v-model:value="searchForm.keyword"
            placeholder="邮箱/昵称"
            allow-clear
            style="width: 200px"
            @change="handleSearch"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-select
            v-model:value="searchForm.status"
            placeholder="全部"
            style="width: 120px"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="disabled">禁用</a-select-option>
          </a-select>
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
      :data-source="userList"
      :loading="loading"
      :pagination="paginationConfig"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'created_at'">
          {{ formatDate(record.created_at) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space size="small">
            <a-button type="link" size="small" @click="handleViewDetail(record)">
              查看详情
            </a-button>
            <a-button
              type="link"
              size="small"
              :danger="record.status === 'active'"
              @click="handleToggleStatus(record)"
            >
              {{ record.status === 'active' ? '禁用' : '启用' }}
            </a-button>
            <a-button type="link" size="small" @click="handleResetPassword(record)">
              重置密码
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-drawer
      v-model:open="detailDrawerVisible"
      title="用户详情"
      :width="480"
      :mask-closable="false"
    >
      <div v-if="userDetail" class="user-detail">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="用户ID">
            {{ userDetail.id }}
          </a-descriptions-item>
          <a-descriptions-item label="邮箱">
            {{ userDetail.email }}
          </a-descriptions-item>
          <a-descriptions-item label="昵称">
            {{ userDetail.nickname }}
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="userDetail.status === 'active' ? 'green' : 'red'">
              {{ userDetail.status === 'active' ? '启用' : '禁用' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="注册时间">
            {{ formatDate(userDetail.created_at) }}
          </a-descriptions-item>
          <a-descriptions-item label="更新时间">
            {{ formatDate(userDetail.updated_at) }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>数据统计</a-divider>

        <a-row :gutter="16">
          <a-col :span="8">
            <a-card size="small" class="stat-mini-card">
              <p class="stat-mini-label">简历数</p>
              <p class="stat-mini-value">{{ userDetail.resume_count }}</p>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card size="small" class="stat-mini-card">
              <p class="stat-mini-label">投递数</p>
              <p class="stat-mini-value">{{ userDetail.delivery_count }}</p>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card size="small" class="stat-mini-card">
              <p class="stat-mini-label">面试数</p>
              <p class="stat-mini-value">{{ userDetail.interview_count }}</p>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </a-drawer>

    <a-modal
      v-model:open="passwordModalVisible"
      title="重置密码"
      :mask-closable="false"
      @ok="handleConfirmResetPassword"
      @cancel="handleCancelResetPassword"
    >
      <a-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef">
        <a-form-item label="新密码" name="newPassword">
          <a-input-password v-model:value="passwordForm.newPassword" placeholder="请输入新密码" />
        </a-form-item>
        <a-form-item label="确认密码" name="confirmPassword">
          <a-input-password v-model:value="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import * as adminApi from '@/api/admin'
import type { AdminUser, AdminUserDetail, AdminUserListParams } from '@/types/models'
import type { TableColumnsType, TablePaginationConfig } from 'ant-design-vue'
import type { Rule } from 'ant-design-vue/es/form'

const loading = ref(false)
const userList = ref<AdminUser[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<AdminUserListParams>({
  keyword: '',
  status: '',
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
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: '昵称',
    dataIndex: 'nickname',
    key: 'nickname',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
  },
  {
    title: '注册时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 180,
  },
  {
    title: '操作',
    key: 'action',
    width: 220,
    fixed: 'right',
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

const detailDrawerVisible = ref(false)
const userDetail = ref<AdminUserDetail | null>(null)

const passwordModalVisible = ref(false)
const passwordFormRef = ref()
const currentUserId = ref<number | null>(null)
const passwordForm = reactive({
  newPassword: '',
  confirmPassword: '',
})

const validateConfirmPassword = (_rule: Rule, value: string) => {
  if (!value) {
    return Promise.reject('请再次输入新密码')
  }
  if (value !== passwordForm.newPassword) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const passwordRules: Record<string, Rule[]> = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能少于8位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' },
  ],
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

const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await adminApi.getUsers({
      keyword: searchForm.keyword || undefined,
      status: searchForm.status || undefined,
      page: currentPage.value,
      page_size: pageSize.value,
    })
    userList.value = response.data.data.list || []
    total.value = response.data.data.total
    paginationConfig.total = response.data.data.total
  } catch (error) {
    console.error('获取用户列表失败:', error)
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  paginationConfig.current = 1
  fetchUsers()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  currentPage.value = 1
  paginationConfig.current = 1
  fetchUsers()
}

const handleTableChange = (pagination: TablePaginationConfig) => {
  currentPage.value = pagination.current || 1
  pageSize.value = pagination.pageSize || 10
  fetchUsers()
}

const handleViewDetail = async (record: AdminUser) => {
  detailDrawerVisible.value = true
  userDetail.value = null
  try {
    const response = await adminApi.getUserDetail(record.id)
    userDetail.value = response.data.data
  } catch (error) {
    console.error('获取用户详情失败:', error)
    message.error('获取用户详情失败')
  }
}

const handleToggleStatus = (record: AdminUser) => {
  const newStatus = record.status === 'active' ? 'disabled' : 'active'
  const actionText = newStatus === 'active' ? '启用' : '禁用'

  Modal.confirm({
    title: `${actionText}用户`,
    content: `确定要${actionText}用户「${record.nickname}」吗？`,
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await adminApi.toggleUserStatus(record.id, newStatus)
        message.success(`用户${actionText}成功`)
        fetchUsers()
      } catch (error) {
        console.error('操作用户状态失败:', error)
        message.error('操作失败')
      }
    },
  })
}

const handleResetPassword = (record: AdminUser) => {
  currentUserId.value = record.id
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  passwordModalVisible.value = true
}

const handleConfirmResetPassword = async () => {
  if (!passwordFormRef.value || !currentUserId.value) return

  try {
    await passwordFormRef.value.validate()
  } catch {
    return
  }

  try {
    await adminApi.resetUserPassword(currentUserId.value, passwordForm.newPassword)
    message.success('密码重置成功')
    passwordModalVisible.value = false
  } catch (error) {
    console.error('重置密码失败:', error)
    message.error('重置密码失败')
  }
}

const handleCancelResetPassword = () => {
  passwordModalVisible.value = false
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  currentUserId.value = null
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.admin-users {
  width: 100%;
}

.search-bar {
  margin-bottom: 16px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.user-detail {
  padding: 8px 0;
}

.stat-mini-card {
  text-align: center;
}

.stat-mini-label {
  color: #999;
  font-size: 12px;
  margin-bottom: 4px;
}

.stat-mini-value {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
}
</style>
