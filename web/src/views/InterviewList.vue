<template>
  <div class="interview-list-page">
    <!-- 顶部工具栏 -->
    <div class="page-header">
      <div class="header-info">
        <h2 class="page-title">模拟面试</h2>
        <p class="page-desc">AI 自动发问，支持文字/语音/混合模式，结束生成复盘报告</p>
      </div>
      <a-button type="primary" @click="showCreateModal = true">
        <PlusOutlined /> 开始新面试
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

    <!-- 创建面试弹窗 -->
    <a-modal
      v-model:open="showCreateModal"
      title="开始新面试"
      :confirm-loading="creating"
      @ok="handleCreate"
      @cancel="resetCreateForm"
      width="640px"
    >
      <a-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        layout="vertical"
      >
        <a-form-item label="面试场景" name="scene" required>
          <a-select v-model:value="createForm.scene" placeholder="请选择">
            <a-select-option value="tech">技术面</a-select-option>
            <a-select-option value="behavior">行为面</a-select-option>
            <a-select-option value="pressure">压力面</a-select-option>
            <a-select-option value="hr">HR 面</a-select-option>
            <a-select-option value="group">群面</a-select-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="目标公司" name="target_company">
              <a-input v-model:value="createForm.target_company" placeholder="如：字节跳动" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="目标职位" name="target_position" required>
              <a-input v-model:value="createForm.target_position" placeholder="如：后端开发" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="目标 JD（可选）" name="target_jd">
          <a-textarea
            v-model:value="createForm.target_jd"
            :rows="3"
            placeholder="粘贴岗位描述，AI 将据此出题"
          />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="难度" name="difficulty">
              <a-select v-model:value="createForm.difficulty">
                <a-select-option value="junior">初级</a-select-option>
                <a-select-option value="mid">中级</a-select-option>
                <a-select-option value="senior">高级</a-select-option>
                <a-select-option value="mixed">混合自适应</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="题数" name="total_questions">
              <a-input-number
                v-model:value="createForm.total_questions"
                :min="1"
                :max="20"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="模式" name="mode">
              <a-select v-model:value="createForm.mode">
                <a-select-option value="text">纯文字</a-select-option>
                <a-select-option value="voice">语音</a-select-option>
                <a-select-option value="hybrid">混合</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-alert
          message="创建后将自动生成第一道题，可在面试房间内开始作答"
          type="info"
          show-icon
          banner
        />
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, type FormInstance, type TableColumnsType } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useInterviewStore } from '@/stores/interview'
import type {
  InterviewScene,
  InterviewMode,
  InterviewDifficulty,
  InterviewStatus,
  CreateInterviewRequest,
} from '@/types/models'

const router = useRouter()
const interviewStore = useInterviewStore()

const showCreateModal = ref(false)
const creating = ref(false)
const createFormRef = ref<FormInstance>()

const createForm = reactive<CreateInterviewRequest>({
  scene: 'tech',
  target_company: '',
  target_position: '',
  target_jd: '',
  difficulty: 'mid',
  total_questions: 5,
  mode: 'text',
})

const createRules = {
  scene: [{ required: true, message: '请选择面试场景', trigger: 'change' }],
  target_position: [{ required: true, message: '请输入目标职位', trigger: 'blur' }],
}

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

// 创建面试
const handleCreate = async () => {
  try {
    await createFormRef.value?.validate()
  } catch {
    return
  }
  creating.value = true
  const payload: CreateInterviewRequest = {
    scene: createForm.scene,
    target_company: createForm.target_company?.trim() || undefined,
    target_position: createForm.target_position.trim(),
    target_jd: createForm.target_jd?.trim() || undefined,
    difficulty: createForm.difficulty,
    total_questions: createForm.total_questions,
    mode: createForm.mode,
  }
  const created = await interviewStore.create(payload)
  creating.value = false
  if (created) {
    showCreateModal.value = false
    resetCreateForm()
    router.push(`/app/interviews/${created.id}`)
  } else {
    message.error('创建失败，请重试')
  }
}

const resetCreateForm = () => {
  createForm.scene = 'tech'
  createForm.target_company = ''
  createForm.target_position = ''
  createForm.target_jd = ''
  createForm.difficulty = 'mid'
  createForm.total_questions = 5
  createForm.mode = 'text'
  createFormRef.value?.resetFields()
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
