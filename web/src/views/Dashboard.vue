<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <h1 class="welcome-title">欢迎回来，{{ authStore.user?.nickname || '用户' }}！</h1>
      <p class="welcome-subtitle">今天是 {{ currentDate }}，祝您求职顺利！</p>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="[16, 16]" class="stats-row">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic
            title="投递总数"
            :value="stats?.total || 0"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <SendOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic
            title="进行中"
            :value="stats?.in_progress || 0"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <SyncOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic
            title="已获Offer"
            :value="stats?.offer_count || 0"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <TrophyOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic
            title="已拒绝"
            :value="stats?.rejected || 0"
            :value-style="{ color: '#999' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 投递漏斗 -->
    <a-card title="投递漏斗" class="funnel-card">
      <div class="funnel-container">
        <div v-if="funnel" class="funnel-chart">
          <div class="funnel-item" :style="{ width: '100%' }">
            <div class="funnel-bar" style="background: #1890ff">
              <span class="funnel-label">已投递</span>
              <span class="funnel-value">{{ funnel.applied }}</span>
            </div>
          </div>
          <div class="funnel-item" :style="{ width: `${(funnel.written_test_pass / (funnel.applied || 1)) * 100}%` }">
            <div class="funnel-bar" style="background: #52c41a">
              <span class="funnel-label">笔试通过</span>
              <span class="funnel-value">{{ funnel.written_test_pass }}</span>
            </div>
          </div>
          <div class="funnel-item" :style="{ width: `${(funnel.first_pass / (funnel.applied || 1)) * 100}%` }">
            <div class="funnel-bar" style="background: #faad14">
              <span class="funnel-label">一面通过</span>
              <span class="funnel-value">{{ funnel.first_pass }}</span>
            </div>
          </div>
          <div class="funnel-item" :style="{ width: `${(funnel.second_pass / (funnel.applied || 1)) * 100}%` }">
            <div class="funnel-bar" style="background: #722ed1">
              <span class="funnel-label">二面通过</span>
              <span class="funnel-value">{{ funnel.second_pass }}</span>
            </div>
          </div>
          <div class="funnel-item" :style="{ width: `${(funnel.offer_count / (funnel.applied || 1)) * 100}%` }">
            <div class="funnel-bar" style="background: '#f5222d'">
              <span class="funnel-label">已获Offer</span>
              <span class="funnel-value">{{ funnel.offer_count }}</span>
            </div>
          </div>
        </div>
        <a-empty v-else description="暂无数据" />
      </div>
    </a-card>

    <!-- 快捷入口 -->
    <a-row :gutter="[16, 16]" class="quick-actions">
      <a-col :xs="24" :sm="8">
        <a-card hoverable class="action-card" @click="router.push('/app/profile')">
          <div class="action-content">
            <UserOutlined class="action-icon" />
            <h3>完善档案</h3>
            <p>填写个人基本信息</p>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="8">
        <a-card hoverable class="action-card" @click="router.push('/app/resumes')">
          <div class="action-content">
            <FileTextOutlined class="action-icon" />
            <h3>创建简历</h3>
            <p>使用AI生成简历</p>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="8">
        <a-card hoverable class="action-card" @click="router.push('/app/deliveries')">
          <div class="action-content">
            <SendOutlined class="action-icon" />
            <h3>开始投递</h3>
            <p>记录投递进度</p>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDeliveryStore } from '@/stores/delivery'
import {
  SendOutlined,
  SyncOutlined,
  TrophyOutlined,
  ClockCircleOutlined,
  UserOutlined,
  FileTextOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'

const router = useRouter()
const authStore = useAuthStore()
const deliveryStore = useDeliveryStore()

// 当前日期
const currentDate = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))

// 统计数据
const stats = computed(() => deliveryStore.stats)

// 漏斗数据
const funnel = computed(() => deliveryStore.funnel)

// 页面加载时获取数据
onMounted(async () => {
  await deliveryStore.fetchStats()
  await deliveryStore.fetchFunnel()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.welcome-section {
  margin-bottom: 24px;
}

.welcome-title {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.welcome-subtitle {
  font-size: 14px;
  color: #8c8c8c;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
  transition: all 0.3s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.funnel-card {
  margin-bottom: 24px;
}

.funnel-container {
  padding: 16px 0;
}

.funnel-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.funnel-item {
  margin: 0 auto;
  transition: width 0.3s;
}

.funnel-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-radius: 8px;
  color: #fff;
  font-weight: 500;
}

.funnel-label {
  font-size: 14px;
}

.funnel-value {
  font-size: 18px;
  font-weight: 600;
}

.quick-actions {
  margin-top: 24px;
}

.action-card {
  text-align: center;
  border-radius: 8px;
  cursor: pointer;
}

.action-content {
  padding: 24px;
}

.action-icon {
  font-size: 48px;
  color: #1890ff;
  margin-bottom: 16px;
}

.action-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.action-content p {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0;
}
</style>