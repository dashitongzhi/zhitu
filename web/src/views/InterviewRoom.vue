<template>
  <div class="interview-room-page">
    <a-spin :spinning="interviewStore.loading" tip="加载面试...">
      <!-- 顶部信息栏 -->
      <div class="room-header">
        <div class="header-left">
          <a-button @click="backToList">
            <ArrowLeftOutlined /> 返回
          </a-button>
          <div class="interview-meta" v-if="interviewStore.currentInterview">
            <a-tag :color="sceneColor(interviewStore.currentInterview.scene)">
              {{ sceneLabel(interviewStore.currentInterview.scene) }}
            </a-tag>
            <span class="meta-text">
              {{ interviewStore.currentInterview.target_company || '未指定' }}
              · {{ interviewStore.currentInterview.target_position }}
            </span>
            <a-tag>{{ difficultyLabel(interviewStore.currentInterview.difficulty) }}</a-tag>
            <a-tag :color="modeColor(interviewStore.currentInterview.mode)">
              {{ modeLabel(interviewStore.currentInterview.mode) }}
            </a-tag>
            <a-divider type="vertical" />
            <span class="meta-progress">
              进度：{{ interviewStore.currentInterview.current_question_no }}
              / {{ interviewStore.currentInterview.total_questions }}
            </span>
            <a-badge
              :status="statusBadge(interviewStore.currentInterview.status)"
              :text="statusLabel(interviewStore.currentInterview.status)"
            />
          </div>
        </div>
        <div class="header-right">
          <a-button
            v-if="interviewStore.currentInterview?.status === 'ongoing'"
            danger
            @click="handleEndInterview"
          >
            <PoweroffOutlined /> 结束面试
          </a-button>
        </div>
      </div>

      <section v-if="isTeachingScene" class="classroom-stage">
        <div class="stage-meta">
          <span>模拟教室</span>
          <strong>{{ currentTeachingPhase }}</strong>
          <small>第 {{ interviewStore.currentInterview?.current_question_no || 1 }} 环节</small>
        </div>
        <div class="examiner-bench">
          <div class="examiner-person side"><i>考</i><span>学科考官</span></div>
          <div class="examiner-person lead"><i>主</i><span>主考官</span></div>
          <div class="examiner-person side"><i>记</i><span>记录员</span></div>
        </div>
        <div class="teaching-board">
          <span>当前考题</span>
          <p>{{ latestQuestion }}</p>
          <div class="chalk-line"></div>
        </div>
        <div class="candidate-position">
          <span>考生位置</span>
          <strong>请面向考官开始作答</strong>
        </div>
        <div class="stage-timer">
          <ClockCircleOutlined />
          <span>本环节建议</span>
          <strong>{{ teachingTimeHint }}</strong>
        </div>
      </section>

      <!-- 主体：左消息流 + 右信息/复盘 -->
      <a-row :gutter="16" class="room-body">
        <!-- 左侧消息流 -->
        <a-col :span="16">
          <a-card class="chat-card" :bordered="false" :body-style="{ padding: '0' }">
            <template #title>
              <span><MessageOutlined /> 面试对话</span>
            </template>

            <!-- 消息列表 -->
            <div class="chat-messages" ref="messagesContainer">
              <a-empty
                v-if="interviewStore.messages.length === 0"
                description="暂无消息"
                class="chat-empty"
              />
              <div
                v-for="msg in displayMessages"
                :key="msg.id"
                :class="['msg-row', msg.role]"
              >
                <div class="msg-avatar">
                  <a-avatar :size="32" :style="avatarStyle(msg.role)">
                    {{ avatarText(msg.role) }}
                  </a-avatar>
                </div>
                <div class="msg-bubble-wrap">
                  <div class="msg-meta">
                    <span class="msg-role">{{ roleLabel(msg.role) }}</span>
                    <span v-if="msg.question_no > 0" class="msg-qno">
                      第 {{ msg.question_no }} 题
                    </span>
                    <span v-if="msg.question_type" class="msg-qtype">
                      · {{ msg.question_type }}
                    </span>
                    <span v-if="msg.duration_sec > 0" class="msg-duration">
                      · {{ msg.duration_sec }}s
                    </span>
                    <!-- TTS 播放按钮：仅 assistant 消息 + mode 为 voice/hybrid -->
                    <a-button
                      v-if="canPlayTts && msg.role === 'assistant'"
                      type="link"
                      size="small"
                      :loading="ttsLoadingId === msg.id"
                      @click="handlePlayTts(msg.id)"
                    >
                      <SoundOutlined /> 播放
                    </a-button>
                  </div>
                  <div class="msg-bubble">
                    {{ msg.content }}
                  </div>
                  <div class="msg-time">{{ formatTime(msg.created_at) }}</div>
                </div>
              </div>

              <!-- 流式 AI 回复占位 -->
              <div v-if="interviewStore.streamingText" class="msg-row assistant">
                <div class="msg-avatar">
                  <a-avatar :size="32" style="background: #1890ff">AI</a-avatar>
                </div>
                <div class="msg-bubble-wrap">
                  <div class="msg-meta">
                    <span class="msg-role">AI 助手</span>
                    <a-spin :spinning="interviewStore.sending" size="small" />
                  </div>
                  <div class="msg-bubble streaming-bubble">
                    {{ interviewStore.streamingText }}
                    <span class="cursor-blink">|</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 输入区 -->
            <div class="chat-input-area" v-if="isOngoing">
              <a-textarea
                v-model:value="inputText"
                :rows="3"
                placeholder="输入你的回答..."
                :disabled="interviewStore.sending"
                @keydown.enter.exact.prevent="handleSendText"
              />
              <div class="input-toolbar">
                <div class="toolbar-left">
                  <!-- 语音上传（mode 为 voice/hybrid 时可用） -->
                  <a-upload
                    v-if="canSendVoice"
                    :before-upload="handleSendVoice"
                    :show-upload-list="false"
                    accept="audio/*,mp3,wav,m4a,webm,ogg"
                  >
                    <a-button :disabled="interviewStore.sending">
                      <AudioOutlined /> 语音回答
                    </a-button>
                  </a-upload>
                  <!-- 发送简历：把简历绑定到本次面试，AI 后续提问会结合简历内容 -->
                  <a-button
                    :type="resumeAttached ? 'default' : 'dashed'"
                    :disabled="interviewStore.sending || !isOngoing"
                    @click="openResumeModal"
                  >
                    <FileTextOutlined />
                    <span v-if="resumeAttached">简历已发送</span>
                    <span v-else>发送简历</span>
                  </a-button>
                  <span class="input-hint">Enter 发送 · Shift+Enter 换行</span>
                </div>
                <div class="toolbar-right">
                  <a-button
                    type="primary"
                    :loading="interviewStore.sending"
                    :disabled="!inputText.trim()"
                    @click="handleSendText"
                  >
                    <SendOutlined /> 发送
                  </a-button>
                </div>
              </div>
            </div>
            <div v-else class="chat-input-disabled">
              <a-alert
                :message="`面试已${interviewStore.currentInterview?.status === 'completed' ? '完成' : '取消'}，请在右侧查看复盘报告`"
                type="info"
                show-icon
              />
            </div>
          </a-card>
        </a-col>

        <!-- 右侧信息/复盘 -->
        <a-col :span="8">
          <a-card class="side-card" :bordered="false">
            <a-tabs v-model:activeKey="sideTab">
              <!-- 面试信息 -->
              <a-tab-pane key="info" tab="面试信息">
                <div v-if="interviewStore.currentInterview" class="info-block">
                  <a-descriptions :column="1" size="small" bordered>
                    <a-descriptions-item label="场景">
                      {{ sceneLabel(interviewStore.currentInterview.scene) }}
                    </a-descriptions-item>
                    <a-descriptions-item label="目标公司">
                      {{ interviewStore.currentInterview.target_company || '未指定' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="目标职位">
                      {{ interviewStore.currentInterview.target_position }}
                    </a-descriptions-item>
                    <a-descriptions-item label="难度">
                      {{ difficultyLabel(interviewStore.currentInterview.difficulty) }}
                    </a-descriptions-item>
                    <a-descriptions-item label="模式">
                      {{ modeLabel(interviewStore.currentInterview.mode) }}
                    </a-descriptions-item>
                    <a-descriptions-item label="进度">
                      {{ interviewStore.currentInterview.current_question_no }}
                      / {{ interviewStore.currentInterview.total_questions }}
                    </a-descriptions-item>
                    <a-descriptions-item label="状态">
                      <a-badge
                        :status="statusBadge(interviewStore.currentInterview.status)"
                        :text="statusLabel(interviewStore.currentInterview.status)"
                      />
                    </a-descriptions-item>
                    <a-descriptions-item label="开始时间">
                      {{ formatTime(interviewStore.currentInterview.started_at) }}
                    </a-descriptions-item>
                    <a-descriptions-item label="结束时间">
                      {{ formatTime(interviewStore.currentInterview.ended_at) }}
                    </a-descriptions-item>
                  </a-descriptions>

                  <div v-if="interviewStore.currentInterview.target_jd" class="jd-block">
                    <div class="block-title">目标 JD</div>
                    <div class="jd-content">
                      {{ interviewStore.currentInterview.target_jd }}
                    </div>
                  </div>

                  <div class="resume-block-info">
                    <div class="block-title">已发送简历</div>
                    <a-alert
                      v-if="resumeAttached"
                      type="success"
                      show-icon
                      :message="`简历：${interviewStore.currentInterview.resume_name || '未命名简历'}`"
                      description="AI 将在后续提问中结合该简历内容进行追问，评分与复盘也会参考简历经历。"
                    />
                    <a-empty
                      v-else
                      description="尚未发送简历，可在输入区点击「发送简历」绑定"
                      :image="simpleImage"
                    />
                  </div>
                </div>
                <a-empty v-else description="暂无信息" />
              </a-tab-pane>

              <!-- 复盘报告 -->
              <a-tab-pane key="report" tab="复盘报告" :disabled="!isCompleted">
                <a-spin :spinning="reportLoading" tip="加载报告...">
                  <div v-if="interviewStore.report" class="report-block">
                    <a-statistic
                      title="整体评分"
                      :value="interviewStore.report.overall_score"
                      suffix="/ 100"
                      :value-style="scoreStyle(interviewStore.report.overall_score)"
                      class="overall-score"
                    />
                    <div class="block-title">综合评语</div>
                    <a-typography-paragraph>
                      {{ interviewStore.report.summary }}
                    </a-typography-paragraph>

                    <div class="block-title">亮点</div>
                    <a-list
                      v-if="parseJsonArray(interviewStore.report.highlights).length"
                      size="small"
                      :data-source="parseJsonArray(interviewStore.report.highlights)"
                    >
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <CheckCircleOutlined style="color: #52c41a" />
                          <span class="list-item-text">{{ item }}</span>
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-empty v-else description="无" :image="simpleImage" />

                    <div class="block-title">改进建议</div>
                    <a-list
                      v-if="parseJsonArray(interviewStore.report.improvements).length"
                      size="small"
                      :data-source="parseJsonArray(interviewStore.report.improvements)"
                    >
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <WarningOutlined style="color: #faad14" />
                          <span class="list-item-text">{{ item }}</span>
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-empty v-else description="无" :image="simpleImage" />

                    <div class="block-title">推荐资源</div>
                    <a-list
                      v-if="parseJsonArray(interviewStore.report.recommendations).length"
                      size="small"
                      :data-source="parseJsonArray(interviewStore.report.recommendations)"
                    >
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <BulbOutlined style="color: #1890ff" />
                          <span class="list-item-text">{{ item }}</span>
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-empty v-else description="无" :image="simpleImage" />
                  </div>
                  <a-empty
                    v-else-if="!reportLoading"
                    description="暂无报告，请先结束面试"
                    :image="simpleImage"
                  />
                </a-spin>
              </a-tab-pane>

              <!-- 评分明细 -->
              <a-tab-pane key="scores" tab="评分明细" :disabled="!isCompleted">
                <a-spin :spinning="scoresLoading" tip="加载评分...">
                  <a-table
                    v-if="interviewStore.scores.length > 0"
                    :columns="scoreColumns"
                    :data-source="interviewStore.scores"
                    :pagination="false"
                    row-key="id"
                    size="small"
                  >
                    <template #bodyCell="{ column, record }">
                      <template v-if="column.key === 'dimension'">
                        {{ dimensionLabel(record.dimension) }}
                      </template>
                      <template v-else-if="column.key === 'score'">
                        <a-tag :color="scoreTagColor(record.score)">{{ record.score }}</a-tag>
                      </template>
                    </template>
                  </a-table>
                  <a-empty
                    v-else-if="!scoresLoading"
                    description="暂无评分明细"
                    :image="simpleImage"
                  />
                </a-spin>
              </a-tab-pane>
            </a-tabs>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>

    <!-- 隐藏的 audio 元素用于 TTS 播放 -->
    <audio ref="audioPlayer" class="hidden-audio" />

    <!-- 发送简历弹窗：选择简历 + 版本 -->
    <a-modal
      v-model:open="resumeModalVisible"
      title="发送简历给 AI 面试官"
      ok-text="发送简历"
      cancel-text="取消"
      :ok-button-props="{ loading: attaching, disabled: !selectedResumeId }"
      :cancel-button-props="{ disabled: attaching }"
      :mask-closable="false"
      width="520px"
      @ok="confirmAttachResume"
    >
      <a-spin :spinning="loadingResumes" tip="加载简历列表...">
        <a-empty
          v-if="!loadingResumes && resumeList.length === 0"
          description="还没有简历，请先到「简历实验室」创建简历"
        >
          <a-button type="primary" @click="goToResumeLab">前往简历实验室</a-button>
        </a-empty>
        <div v-else class="resume-modal-body">
          <p class="modal-tip">
            选择一份简历发送给 AI 面试官，AI 将在后续提问中结合简历的项目、工作经历、技能进行深挖追问，评分与复盘也会参考简历。
          </p>
          <a-form layout="vertical">
            <a-form-item label="选择简历">
              <a-select
                v-model:value="selectedResumeId"
                placeholder="请选择简历"
                :disabled="attaching"
                @change="onResumeChange"
              >
                <a-select-option
                  v-for="r in resumeList"
                  :key="r.id"
                  :value="r.id"
                >
                  {{ r.name }}
                  <span class="resume-meta">
                    （{{ r.target_position || '未指定岗位' }}）
                  </span>
                </a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="选择版本">
              <a-select
                v-model:value="selectedVersionId"
                placeholder="默认使用当前版本"
                :disabled="attaching || !selectedResumeId || loadingVersions"
                :loading="loadingVersions"
              >
                <a-select-option :value="0">当前版本</a-select-option>
                <a-select-option
                  v-for="v in versionList"
                  :key="v.id"
                  :value="v.id"
                >
                  {{ v.version_label }}
                  <span class="resume-meta">· {{ formatVersionTime(v.created_at) }}</span>
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-form>
          <a-alert
            v-if="resumeAttached"
            type="info"
            show-icon
            message="本次面试已绑定过简历，再次发送将覆盖之前绑定的简历。"
          />
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, message, Empty, type TableColumnsType } from 'ant-design-vue'
import {
  ArrowLeftOutlined,
  PoweroffOutlined,
  MessageOutlined,
  SoundOutlined,
  AudioOutlined,
  SendOutlined,
  FileTextOutlined,
  CheckCircleOutlined,
  WarningOutlined,
  BulbOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import { useInterviewStore } from '@/stores/interview'
import { listResumes, listVersions } from '@/api/resume'
import type {
  InterviewScene,
  InterviewMode,
  InterviewDifficulty,
  InterviewStatus,
  InterviewDimension,
  Resume,
  ResumeVersion,
} from '@/types/models'

const route = useRoute()
const router = useRouter()
const interviewStore = useInterviewStore()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const interviewId = computed(() => Number(route.params.id))

// 输入与界面状态
const inputText = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const audioPlayer = ref<HTMLAudioElement | null>(null)
const sideTab = ref<'info' | 'report' | 'scores'>('info')

// 加载状态
const reportLoading = ref(false)
const scoresLoading = ref(false)
const ttsLoadingId = ref<number | null>(null)

// 发送简历相关状态
const resumeModalVisible = ref(false)
const loadingResumes = ref(false)
const loadingVersions = ref(false)
const attaching = ref(false)
const resumeList = ref<Resume[]>([])
const versionList = ref<ResumeVersion[]>([])
const selectedResumeId = ref<number | null>(null)
const selectedVersionId = ref<number>(0)
const resumeAttached = computed(
  () => !!interviewStore.currentInterview?.resume_snapshot
)

// 计算属性
const isOngoing = computed(
  () => interviewStore.currentInterview?.status === 'ongoing'
)
const isCompleted = computed(
  () => interviewStore.currentInterview?.status === 'completed'
)
const canSendVoice = computed(() => {
  const mode = interviewStore.currentInterview?.mode
  return mode === 'voice' || mode === 'hybrid'
})
const canPlayTts = computed(() => {
  const mode = interviewStore.currentInterview?.mode
  return mode === 'voice' || mode === 'hybrid'
})
const isTeachingScene = computed(() => interviewStore.currentInterview?.scene === 'teaching')
const latestQuestion = computed(() => {
  const assistantMessages = interviewStore.messages.filter((item) => item.role === 'assistant')
  return assistantMessages[assistantMessages.length - 1]?.content || '考官正在准备第一道结构化问题…'
})
const currentTeachingPhase = computed(() => {
  const no = interviewStore.currentInterview?.current_question_no || 1
  if (no <= 2) return '结构化问答'
  if (no <= 4) return '模拟试讲'
  return '考官答辩'
})
const teachingTimeHint = computed(() => currentTeachingPhase.value === '模拟试讲' ? '8 分钟' : '2 分钟')

// 显示的消息列表（隐藏流式时的临时占位）
const displayMessages = computed(() => interviewStore.messages)

// 自动滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

watch(
  () => interviewStore.messages.length,
  () => scrollToBottom()
)
watch(
  () => interviewStore.streamingText,
  () => scrollToBottom()
)

// 场景
const sceneLabel = (s: InterviewScene | string): string => {
  const map: Record<string, string> = {
    teaching: '模拟教室',
    corporate: '企业会议室',
    group: '群面讨论室',
    defense: '项目答辩室',
    client: '客户会议室',
    pressure: '压力面试室',
    public: '结构化面试厅',
    medical: '医疗面试室',
    media: '媒体演播室',
    remote: '远程面试间',
    system: '系统设计室',
    aviation: '航空面试厅',
    tech: '技术面',
    behavior: '行为面',
    hr: 'HR 面',
  }
  return map[s] || s
}
const sceneColor = (s: InterviewScene | string): string => {
  const map: Record<string, string> = {
    teaching: 'green',
    corporate: 'blue',
    group: 'orange',
    defense: 'geekblue',
    client: 'cyan',
    pressure: 'red',
    public: 'gold',
    medical: 'green',
    media: 'magenta',
    remote: 'purple',
    system: 'volcano',
    aviation: 'lime',
    tech: 'blue',
    behavior: 'cyan',
    hr: 'purple',
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

// 消息角色
const roleLabel = (r: string): string => (r === 'assistant' ? 'AI 面试官' : '我')
const avatarText = (r: string): string => (r === 'assistant' ? 'AI' : '我')
const avatarStyle = (r: string): Record<string, string> => {
  if (r === 'assistant') return { background: '#1890ff', color: '#fff' }
  return { background: '#52c41a', color: '#fff' }
}

// 维度
const dimensionLabel = (d: InterviewDimension | string): string => {
  const map: Record<string, string> = {
    professional: '专业能力',
    expression: '表达能力',
    logic: '逻辑思维',
    adaptability: '应变能力',
    pace: '节奏掌控',
  }
  return map[d] || d
}

// 评分相关样式
const scoreStyle = (score: number): Record<string, string> => {
  if (score >= 80) return { color: '#52c41a' }
  if (score >= 60) return { color: '#faad14' }
  return { color: '#f5222d' }
}
const scoreTagColor = (score: number): string => {
  if (score >= 80) return 'green'
  if (score >= 60) return 'orange'
  return 'red'
}

// 评分表列
const scoreColumns: TableColumnsType = [
  { title: '维度', key: 'dimension', width: 120 },
  { title: '分数', key: 'score', width: 80, align: 'center' },
  { title: '评语', dataIndex: 'comment', key: 'comment' },
]

// JSON 数组解析（容错）
const parseJsonArray = (str: string | null | undefined): string[] => {
  if (!str) return []
  try {
    const parsed = JSON.parse(str)
    if (Array.isArray(parsed)) return parsed.map((x) => String(x))
    if (typeof parsed === 'string') return [parsed]
    return [JSON.stringify(parsed)]
  } catch {
    // 不是 JSON，按行分割
    return str.split('\n').filter((s) => s.trim())
  }
}

// 时间格式化
const formatTime = (dateStr: string | null | undefined): string => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 返回列表
const backToList = () => {
  interviewStore.clearCurrent()
  router.push('/app/interviews')
}

// ============ 发送简历相关 ============

// 前往简历实验室（无简历时）
const goToResumeLab = () => {
  resumeModalVisible.value = false
  router.push('/app/resumes')
}

// 版本时间格式化（用于下拉展示）
const formatVersionTime = (dateStr: string): string => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 拉取用户简历列表
const loadResumes = async () => {
  loadingResumes.value = true
  try {
    const resp = await listResumes()
    resumeList.value = resp.data.data || []
    // 默认选中第一份简历
    if (resumeList.value.length > 0 && selectedResumeId.value === null) {
      selectedResumeId.value = resumeList.value[0].id
      await loadVersions(resumeList.value[0].id)
    }
  } catch (error) {
    console.error('加载简历列表失败:', error)
    message.error('加载简历列表失败')
  } finally {
    loadingResumes.value = false
  }
}

// 选择简历时加载版本列表
const onResumeChange = async (resumeId: number) => {
  selectedVersionId.value = 0
  versionList.value = []
  if (!resumeId) return
  await loadVersions(resumeId)
}

const loadVersions = async (resumeId: number) => {
  loadingVersions.value = true
  try {
    const resp = await listVersions(resumeId)
    versionList.value = resp.data.data || []
  } catch (error) {
    console.error('加载简历版本失败:', error)
  } finally {
    loadingVersions.value = false
  }
}

// 打开发送简历弹窗
const openResumeModal = async () => {
  resumeModalVisible.value = true
  if (resumeList.value.length === 0) {
    await loadResumes()
  }
}

// 确认发送简历
const confirmAttachResume = async () => {
  if (!selectedResumeId.value) {
    message.warning('请先选择一份简历')
    return
  }
  if (!isOngoing.value) {
    message.warning('面试已结束，无法发送简历')
    return
  }
  attaching.value = true
  try {
    const attached = await interviewStore.attachResume(interviewId.value, {
      resume_id: selectedResumeId.value,
      version_id: selectedVersionId.value || undefined,
    })
    if (!attached) return
    resumeModalVisible.value = false
    // 切到面试信息 Tab，让用户看到已发送简历的展示
    sideTab.value = 'info'
  } finally {
    attaching.value = false
  }
}

// 发送文字回答
const handleSendText = async () => {
  const content = inputText.value.trim()
  if (!content) return
  if (!isOngoing.value) {
    message.warning('面试已结束，无法继续作答')
    return
  }
  inputText.value = ''
  await interviewStore.sendMessage(interviewId.value, content)
}

// 发送语音回答
const handleSendVoice = async (file: File): Promise<boolean> => {
  if (!isOngoing.value) {
    message.warning('面试已结束，无法继续作答')
    return false
  }
  if (file.size > 25 * 1024 * 1024) {
    message.error('音频文件不能超过 25MB')
    return false
  }
  await interviewStore.sendVoice(interviewId.value, file)
  return false // 阻止 a-upload 自动上传
}

// TTS 播放
const handlePlayTts = async (msgId: number) => {
  if (ttsLoadingId.value !== null) return
  ttsLoadingId.value = msgId
  try {
    const url = await interviewStore.playTts(interviewId.value, msgId)
    if (url && audioPlayer.value) {
      // 如果有旧的 URL，先释放
      const oldSrc = audioPlayer.value.src
      if (oldSrc && oldSrc.startsWith('blob:')) {
        URL.revokeObjectURL(oldSrc)
      }
      audioPlayer.value.src = url
      audioPlayer.value.onended = () => {
        if (url.startsWith('blob:')) URL.revokeObjectURL(url)
      }
      await audioPlayer.value.play()
    }
  } finally {
    ttsLoadingId.value = null
  }
}

// 结束面试
const handleEndInterview = () => {
  Modal.confirm({
    title: '确认结束面试',
    content: '结束后将生成复盘报告，无法继续作答。是否继续？',
    okText: '结束面试',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      const report = await interviewStore.endInterview(interviewId.value)
      if (report) {
        sideTab.value = 'report'
        // 同时加载评分明细
        scoresLoading.value = true
        await interviewStore.fetchScores(interviewId.value)
        scoresLoading.value = false
      }
    },
  })
}

// 初始化加载
onMounted(async () => {
  if (!interviewId.value || isNaN(interviewId.value)) {
    message.error('无效的面试 ID')
    router.push('/app/interviews')
    return
  }
  await interviewStore.fetchOne(interviewId.value)
  // 已完成的面试：自动加载报告 + 评分
  if (interviewStore.currentInterview?.status === 'completed') {
    sideTab.value = 'report'
    reportLoading.value = true
    await interviewStore.fetchReport(interviewId.value)
    reportLoading.value = false
    scoresLoading.value = true
    await interviewStore.fetchScores(interviewId.value)
    scoresLoading.value = false
  }
  scrollToBottom()
})
</script>

<style scoped>
.interview-room-page {
  width: 100%;
}

.room-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.interview-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 13px;
}

.meta-text {
  color: #333;
  font-weight: 500;
}

.meta-progress {
  color: #666;
}

.room-body {
  margin-bottom: 16px;
}

.classroom-stage {
  position: relative;
  min-height: 360px;
  margin-bottom: 18px;
  overflow: hidden;
  border: 1px solid #bfcac3;
  background:
    linear-gradient(115deg, rgba(255,255,255,.72), transparent 32%),
    linear-gradient(#dfe7e1 0 69%, #b7875e 69% 72%, #caa77e 72%);
  color: #173e34;
  box-shadow: 0 16px 38px rgba(25, 55, 46, .1);
}

.stage-meta {
  position: absolute;
  z-index: 3;
  top: 20px;
  left: 22px;
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.stage-meta span,
.stage-meta small {
  font-size: 11px;
  color: #667b73;
}

.stage-meta strong {
  color: #b45c36;
  font-size: 16px;
}

.examiner-bench {
  position: absolute;
  top: 32px;
  left: 50%;
  z-index: 2;
  display: flex;
  gap: 54px;
  transform: translateX(-50%);
}

.examiner-person {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 7px;
  color: #587068;
  font-size: 11px;
}

.examiner-person i {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: #285347;
  color: #fff;
  font-style: normal;
}

.examiner-person.lead i {
  background: #b45c36;
}

.teaching-board {
  position: absolute;
  top: 126px;
  left: 15%;
  right: 15%;
  min-height: 142px;
  padding: 22px 28px;
  border: 10px solid #896547;
  background: #173e34;
  color: #eef2e8;
  box-shadow: 0 12px 24px rgba(20, 46, 38, .18);
}

.teaching-board > span {
  color: #9eb4a8;
  font-size: 10px;
  letter-spacing: .14em;
}

.teaching-board p {
  max-width: 900px;
  margin: 12px 0;
  font-family: "Songti SC", "STSong", serif;
  font-size: 18px;
  line-height: 1.65;
}

.chalk-line {
  width: 42%;
  height: 1px;
  background: #8da89a;
}

.candidate-position {
  position: absolute;
  bottom: 18px;
  left: 50%;
  min-width: 240px;
  padding: 12px 20px;
  transform: translateX(-50%);
  background: #704c32;
  color: #fff5e8;
  text-align: center;
  box-shadow: 0 10px 0 #543823;
}

.candidate-position span {
  display: block;
  color: #d4b99f;
  font-size: 10px;
}

.candidate-position strong {
  font-size: 12px;
}

.stage-timer {
  position: absolute;
  right: 20px;
  bottom: 20px;
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 2px 8px;
  padding: 10px 14px;
  background: rgba(249, 249, 244, .92);
}

.stage-timer .anticon {
  grid-row: 1 / 3;
  align-self: center;
}

.stage-timer span {
  font-size: 10px;
  color: #718079;
}

.stage-timer strong {
  font-size: 15px;
}

@media (max-width: 760px) {
  .classroom-stage { min-height: 410px; }
  .examiner-bench { gap: 20px; }
  .teaching-board { left: 5%; right: 5%; top: 136px; }
  .stage-timer { right: 10px; bottom: 72px; }
}

.chat-card {
  border-radius: 8px;
  height: calc(100vh - 220px);
  display: flex;
  flex-direction: column;
}

.chat-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: #f9f9f9;
}

.chat-empty {
  padding: 80px 0;
}

.msg-row {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.msg-row.user {
  flex-direction: row-reverse;
}

.msg-avatar {
  flex-shrink: 0;
}

.msg-bubble-wrap {
  max-width: 75%;
  display: flex;
  flex-direction: column;
}

.msg-row.user .msg-bubble-wrap {
  align-items: flex-end;
}

.msg-meta {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.msg-role {
  font-weight: 500;
  color: #666;
}

.msg-qno,
.msg-qtype,
.msg-duration {
  color: #999;
}

.msg-bubble {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
  white-space: pre-wrap;
  background: #fff;
  border: 1px solid #e8e8e8;
}

.msg-row.assistant .msg-bubble {
  background: #fff;
  border-color: #d9d9d9;
}

.msg-row.user .msg-bubble {
  background: #1890ff;
  color: #fff;
  border-color: #1890ff;
}

.streaming-bubble {
  background: #e6f7ff !important;
  border-color: #91d5ff !important;
  color: #333 !important;
}

.cursor-blink {
  animation: blink 1s infinite;
  font-weight: bold;
  color: #1890ff;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.msg-time {
  font-size: 11px;
  color: #bbb;
  margin-top: 4px;
}

.chat-input-area {
  padding: 12px 16px;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}

.input-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.input-hint {
  color: #999;
  font-size: 12px;
}

.chat-input-disabled {
  padding: 16px;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}

.side-card {
  border-radius: 8px;
  height: calc(100vh - 220px);
  overflow-y: auto;
}

.side-card :deep(.ant-card-body) {
  padding: 12px 16px;
}

.info-block {
  font-size: 13px;
}

.jd-block {
  margin-top: 16px;
}

.block-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 16px 0 8px 0;
}

.jd-content {
  background: #fafafa;
  padding: 12px;
  border-radius: 6px;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  max-height: 200px;
  overflow-y: auto;
}

.report-block {
  font-size: 13px;
}

.overall-score {
  text-align: center;
  margin-bottom: 16px;
  padding: 16px;
  background: #fafafa;
  border-radius: 6px;
}

.list-item-text {
  margin-left: 8px;
}

.hidden-audio {
  display: none;
}

.resume-block-info {
  margin-top: 16px;
}

.resume-modal-body {
  padding-top: 4px;
}

.modal-tip {
  margin: 0 0 16px 0;
  color: #595959;
  font-size: 13px;
  line-height: 1.7;
}

.resume-meta {
  color: #8c8c8c;
  font-size: 12px;
}
</style>
