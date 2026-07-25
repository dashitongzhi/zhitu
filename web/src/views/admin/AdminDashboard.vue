<template>
  <div class="admin-dashboard">
    <a-spin :spinning="loading" tip="加载中...">
      <div class="stats-cards">
        <a-row :gutter="16">
          <a-col :xs="12" :sm="12" :md="6">
            <a-card class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <p class="stat-label">总用户数</p>
                  <p class="stat-value">{{ stats?.total_users || 0 }}</p>
                </div>
                <div class="stat-icon user-icon">
                  <UserOutlined />
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="12" :sm="12" :md="6">
            <a-card class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <p class="stat-label">活跃用户</p>
                  <p class="stat-value">{{ stats?.active_users || 0 }}</p>
                </div>
                <div class="stat-icon active-icon">
                  <UserAddOutlined />
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="12" :sm="12" :md="6">
            <a-card class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <p class="stat-label">总投递数</p>
                  <p class="stat-value">{{ stats?.total_deliveries || 0 }}</p>
                </div>
                <div class="stat-icon delivery-icon">
                  <SendOutlined />
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="12" :sm="12" :md="6">
            <a-card class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <p class="stat-label">已获Offer</p>
                  <p class="stat-value">{{ stats?.offer_count || 0 }}</p>
                </div>
                <div class="stat-icon offer-icon">
                  <TrophyOutlined />
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>

      <a-row :gutter="16" class="content-row">
        <a-col :xs="24" :md="12">
          <a-card title="投递漏斗" class="funnel-card">
            <div class="funnel-container">
              <div class="funnel-item">
                <div class="funnel-label">
                  <span>投递申请</span>
                  <span class="funnel-count">{{ funnel?.applied || 0 }}</span>
                </div>
                <a-progress
                  :percent="100"
                  :show-info="false"
                  stroke-color="#1890ff"
                  :stroke-width="24"
                />
              </div>
              <div class="funnel-item">
                <div class="funnel-label">
                  <span>笔试通过</span>
                  <span class="funnel-count">
                    {{ funnel?.written_test_pass || 0 }}
                    <span class="funnel-rate">({{ funnel?.written_test_rate || 0 }}%)</span>
                  </span>
                </div>
                <a-progress
                  :percent="funnel?.written_test_rate || 0"
                  :show-info="false"
                  stroke-color="#52c41a"
                  :stroke-width="24"
                />
              </div>
              <div class="funnel-item">
                <div class="funnel-label">
                  <span>一面通过</span>
                  <span class="funnel-count">
                    {{ funnel?.first_pass || 0 }}
                    <span class="funnel-rate">({{ funnel?.first_rate || 0 }}%)</span>
                  </span>
                </div>
                <a-progress
                  :percent="funnel?.first_rate || 0"
                  :show-info="false"
                  stroke-color="#faad14"
                  :stroke-width="24"
                />
              </div>
              <div class="funnel-item">
                <div class="funnel-label">
                  <span>二面通过</span>
                  <span class="funnel-count">
                    {{ funnel?.second_pass || 0 }}
                    <span class="funnel-rate">({{ funnel?.second_rate || 0 }}%)</span>
                  </span>
                </div>
                <a-progress
                  :percent="funnel?.second_rate || 0"
                  :show-info="false"
                  stroke-color="#722ed1"
                  :stroke-width="24"
                />
              </div>
              <div class="funnel-item">
                <div class="funnel-label">
                  <span>获得Offer</span>
                  <span class="funnel-count">
                    {{ funnel?.offer_count || 0 }}
                    <span class="funnel-rate">({{ funnel?.offer_rate || 0 }}%)</span>
                  </span>
                </div>
                <a-progress
                  :percent="funnel?.offer_rate || 0"
                  :show-info="false"
                  stroke-color="#eb2f96"
                  :stroke-width="24"
                />
              </div>
            </div>
          </a-card>
        </a-col>

        <a-col :xs="24" :md="12">
          <a-card title="最近注册用户" class="recent-users-card">
            <a-table
              :columns="userColumns"
              :data-source="recentUsers"
              :pagination="false"
              size="small"
              row-key="id"
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
              </template>
            </a-table>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  UserOutlined,
  UserAddOutlined,
  SendOutlined,
  TrophyOutlined,
} from '@ant-design/icons-vue'
import * as adminApi from '@/api/admin'
import type { AdminDashboardStats, AdminDeliveryFunnel, AdminUser } from '@/types/models'
import type { TableColumnsType } from 'ant-design-vue'

const loading = ref(false)
const stats = ref<AdminDashboardStats | null>(null)
const funnel = ref<AdminDeliveryFunnel | null>(null)
const recentUsers = ref<AdminUser[]>([])

const userColumns: TableColumnsType = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 60,
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
    width: 80,
  },
  {
    title: '注册时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 160,
  },
]

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

const fetchData = async () => {
  loading.value = true
  try {
    const [statsRes, funnelRes, usersRes] = await Promise.all([
      adminApi.getStats(),
      adminApi.getFunnel(),
      adminApi.getUsers({ page: 1, page_size: 10 }),
    ])
    stats.value = statsRes.data.data
    funnel.value = funnelRes.data.data
    recentUsers.value = usersRes.data.data.list || []
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-dashboard {
  width: 100%;
}

.stats-cards {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-info {
  flex: 1;
}

.stat-label {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.user-icon {
  background: #e6f7ff;
  color: #1890ff;
}

.active-icon {
  background: #f6ffed;
  color: #52c41a;
}

.delivery-icon {
  background: #fff7e6;
  color: #fa8c16;
}

.offer-icon {
  background: #fff0f6;
  color: #eb2f96;
}

.content-row {
  margin-bottom: 0;
}

.funnel-card,
.recent-users-card {
  height: 100%;
}

.funnel-container {
  padding: 16px 0;
}

.funnel-item {
  margin-bottom: 24px;
}

.funnel-item:last-child {
  margin-bottom: 0;
}

.funnel-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #333;
}

.funnel-count {
  font-weight: 600;
  color: #1a1a2e;
}

.funnel-rate {
  color: #999;
  font-size: 12px;
  font-weight: normal;
  margin-left: 4px;
}
</style>
