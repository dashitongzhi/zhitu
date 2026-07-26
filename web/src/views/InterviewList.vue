<template>
  <div class="interview-list-page">
    <!-- 顶部工具栏 -->
    <div class="page-header">
      <div class="header-info">
        <h2 class="page-title">面试训练场</h2>
        <p class="page-desc">进入真实场景反复演练，让考官、流程和规则共同塑造训练</p>
      </div>
      <a-button type="primary" @click="router.push('/app/interviews/new')">
        <PlusOutlined /> 选择训练场景
      </a-button>
    </div>

    <!-- 面试会话表格 -->
    <a-spin :spinning="interviewStore.loading" tip="加载中...">
      <a-empty
        v-if="!interviewStore.loading && interviewStore.interviews.length === 0"
        description="暂无面试记录，点击右上角开始新面试"
        class="empty-state"
      />
      <a-table
        v-else
        :columns="columns"
        :data-source="interviewStore.interviews"
        :pagination="false"
        row-key="id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'scene'">
            <a-tag :color="sceneColor(record.scene)">{{ sceneLabel(record.scene) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'target'">
            <div class="target-cell">
              <div>{{ record.target_company || '未指定' }}</div>
              <div class="sub-text">{{ record.target_position }}</div>
            </div>
          </template>
          <template v-else-if="column.key === 'difficulty'">
            <a-tag>{{ difficultyLabel(record.difficulty) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'mode'">
            <a-tag :color="modeColor(record.mode)">{{ modeLabel(record.mode) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'progress'">
            <span>{{ record.current_question_no }} / {{ record.total_questions }}</span>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-badge :status="statusBadge(record.status)" :text="statusLabel(record.status)" />
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ formatDate(record.created_at) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <a-button type="link" size="small" @click="enterRoom(record.id)">
              {{ record.status === 'completed' ? '查看复盘' : '继续面试' }}
            </a-button>
          </template>
        </template>
      </a-table>
    </a-spin>

  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { type TableColumnsType } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useInterviewStore } from '@/stores/interview'
import type {
  InterviewScene,
  InterviewMode,
  InterviewDifficulty,
  InterviewStatus,
} from '@/types/models'

const router = useRouter()
const interviewStore = useInterviewStore()

const columns: TableColumnsType = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '场景', dataIndex: 'scene', key: 'scene', width: 90 },
  { title: '目标', key: 'target', width: 200 },
  { title: '难度', dataIndex: 'difficulty', key: 'difficulty', width: 90 },
  { title: '模式', dataIndex: 'mode', key: 'mode', width: 90 },
  { title: '进度', key: 'progress', width: 100, align: 'center' },
  { title: '状态', dataIndex: 'status', key: 'status', width: 110 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 110, fixed: 'right' },
]

// 场景
const sceneLabel = (s: InterviewScene | string): string => {
  const map: Record<string, string> = {
    teaching: '模拟教室',
    tech: '技术面',
    behavior: '行为面',
    pressure: '压力面',
    hr: 'HR 面',
    group: '群面',
  }
  return map[s] || s
}
const sceneColor = (s: InterviewScene | string): string => {
  const map: Record<string, string> = {
    teaching: 'green',
    tech: 'blue',
    behavior: 'cyan',
    pressure: 'red',
    hr: 'purple',
    group: 'orange',
  }
  return map[s] || 'default'
}

// 难度
const difficultyLabel = (d: InterviewDifficulty | string): string => {
  const map: Record<string, string> = {
    junior: '初级',
    mid: '中级',
    senior: '高级',
    mixed: '混合',
  }
  return map[d] || d
}

// 模式
const modeLabel = (m: InterviewMode | string): string => {
  const map: Record<string, string> = {
    text: '文字',
    voice: '语音',
    hybrid: '混合',
  }
  return map[m] || m
}
const modeColor = (m: InterviewMode | string): string => {
  const map: Record<string, string> = {
    text: 'default',
    voice: 'green',
    hybrid: 'geekblue',
  }
  return map[m] || 'default'
}

// 状态
const statusLabel = (s: InterviewStatus | string): string => {
  const map: Record<string, string> = {
    ongoing: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return map[s] || s
}
const statusBadge = (s: InterviewStatus | string): 'success' | 'processing' | 'default' | 'error' => {
  const map: Record<string, 'success' | 'processing' | 'default' | 'error'> = {
    ongoing: 'processing',
    completed: 'success',
    cancelled: 'default',
  }
  return map[s] || 'default'
}

// 日期
const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 进入面试房间
const enterRoom = (id: number) => {
  router.push(`/app/interviews/${id}`)
}

onMounted(() => {
  interviewStore.fetchList()
})
</script>

<style scoped>
.interview-list-page {
  width: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.page-desc {
  color: #999;
  font-size: 13px;
  margin: 0;
}

.empty-state {
  padding: 80px 0;
}

.target-cell {
  line-height: 1.4;
}

.sub-text {
  color: #999;
  font-size: 12px;
}
</style>
