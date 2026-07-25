<template>
  <div class="resume-editor-page">
    <a-spin :spinning="resumeStore.loading" tip="加载简历...">
      <!-- 顶部工具栏 -->
      <div class="editor-toolbar">
        <div class="toolbar-left">
          <a-button @click="backToList">
            <ArrowLeftOutlined /> 返回
          </a-button>
          <a-input
            v-model:value="editableName"
            class="name-input"
            placeholder="简历名称"
            @blur="handleNameSave"
          />
        </div>
        <div class="toolbar-right">
          <a-button @click="showSaveVersionModal = true">
            <SaveOutlined /> 保存版本
          </a-button>
          <a-dropdown>
            <a-button type="primary">
              <RobotOutlined /> AI 操作 <DownOutlined />
            </a-button>
            <template #overlay>
              <a-menu @click="(e) => handleAiMenuClick(e.key)">
                <a-menu-item key="generate">AI 生成（流式）</a-menu-item>
                <a-menu-item key="polish">AI 润色</a-menu-item>
                <a-menu-item key="score">AI 评分</a-menu-item>
                <a-menu-item key="jdMatch">JD 匹配度分析</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          <a-button @click="handleSyncProfile">
            <SyncOutlined /> 同步档案
          </a-button>
        </div>
      </div>

      <!-- 主体布局 -->
      <a-row :gutter="16" class="editor-body">
        <!-- 左侧版本列表 -->
        <a-col :span="6">
          <a-card title="版本历史" class="version-card" :bordered="false">
            <template #extra>
              <a-tooltip title="刷新版本">
                <a-button type="text" size="small" @click="refreshVersions">
                  <ReloadOutlined />
                </a-button>
              </a-tooltip>
            </template>
            <a-empty
              v-if="resumeStore.versions.length === 0"
              description="暂无版本"
              :image="simpleImage"
            />
            <a-list v-else :data-source="resumeStore.versions" size="small">
              <template #renderItem="{ item }">
                <a-list-item
                  :class="[
                    'version-item',
                    { active: resumeStore.currentVersion?.id === item.id },
                  ]"
                  @click="selectVersion(item)"
                >
                  <div class="version-item-content">
                    <div class="version-header">
                      <a-tag color="blue">{{ item.version_label }}</a-tag>
                      <span class="version-time">{{ formatVersionDate(item.created_at) }}</span>
                    </div>
                    <div class="version-note" :title="item.change_note">
                      {{ item.change_note || '无备注' }}
                    </div>
                    <div class="version-actions">
                      <a-button
                        type="link"
                        size="small"
                        @click.stop="handleRollback(item.id)"
                      >
                        回滚
                      </a-button>
                    </div>
                  </div>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>

        <!-- 右侧内容编辑区 -->
        <a-col :span="18">
          <a-card class="content-card" :bordered="false">
            <template #title>
              <div class="content-header">
                <span>
                  {{ resumeStore.currentVersion
                    ? `当前版本：${resumeStore.currentVersion.version_label}`
                    : '请选择版本' }}
                </span>
                <span v-if="resumeStore.currentVersion" class="content-meta">
                  {{ contentStats.chars }} 字符 · {{ contentStats.lines }} 行
                </span>
              </div>
            </template>

            <a-empty
              v-if="!resumeStore.currentVersion"
              description="请从左侧选择一个版本"
              :image="simpleImage"
            />
            <a-textarea
              v-else
              v-model:value="editableContent"
              :rows="22"
              placeholder="简历内容（JSON 字符串）"
              class="content-textarea"
            />
          </a-card>
        </a-col>
      </a-row>
    </a-spin>

    <!-- 保存版本弹窗 -->
    <a-modal
      v-model:open="showSaveVersionModal"
      title="保存新版本"
      :confirm-loading="savingVersion"
      @ok="handleSaveVersion"
    >
      <a-form layout="vertical">
        <a-form-item label="版本备注" required>
          <a-input
            v-model:value="newVersionNote"
            placeholder="如：润色了工作经历"
            :maxlength="100"
            show-count
          />
        </a-form-item>
        <a-alert
          message="将基于当前编辑框内容创建新版本"
          type="info"
          show-icon
          banner
        />
      </a-form>
    </a-modal>

    <!-- AI 生成弹窗（SSE 流式） -->
    <a-modal
      v-model:open="showGenerateModal"
      title="AI 生成简历（流式）"
      width="720px"
      :mask-closable="false"
      :footer="null"
      @cancel="handleGenerateCancel"
    >
      <a-form layout="vertical">
        <a-form-item label="目标 JD（可选）">
          <a-textarea
            v-model:value="generateForm.target_jd"
            placeholder="粘贴岗位描述，AI 将据此生成更精准的简历"
            :rows="3"
          />
        </a-form-item>
        <a-form-item label="生成场景">
          <a-select v-model:value="generateForm.scene" allow-clear placeholder="留空使用默认">
            <a-select-option value="manual">手动</a-select-option>
            <a-select-option value="jd">基于 JD</a-select-option>
            <a-select-option value="scenario">场景化</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模块提示（可选）">
          <a-textarea
            v-model:value="generateForm.module_hints"
            placeholder="如：重点突出项目经历、强调技术深度"
            :rows="2"
          />
        </a-form-item>
      </a-form>

      <div v-if="resumeStore.streamingText" class="streaming-box">
        <div class="streaming-header">
          <a-spin :spinning="resumeStore.aiLoading" size="small" />
          <span>AI 正在生成...</span>
        </div>
        <a-typography-paragraph class="streaming-text">
          <pre>{{ resumeStore.streamingText }}</pre>
        </a-typography-paragraph>
      </div>

      <div class="modal-footer">
        <a-button @click="handleGenerateCancel">取消</a-button>
        <a-button
          type="primary"
          :loading="resumeStore.aiLoading"
          @click="handleGenerate"
        >
          开始生成
        </a-button>
      </div>
    </a-modal>

    <!-- AI 润色弹窗 -->
    <a-modal
      v-model:open="showPolishModal"
      title="AI 润色"
      :confirm-loading="resumeStore.aiLoading"
      @ok="handlePolish"
    >
      <a-form layout="vertical">
        <a-form-item label="润色模块" required>
          <a-select v-model:value="polishForm.module">
            <a-select-option value="work">工作经历</a-select-option>
            <a-select-option value="project">项目经历</a-select-option>
            <a-select-option value="all">全部</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="参考 JD（可选）">
          <a-textarea v-model:value="polishForm.jd" :rows="3" placeholder="粘贴 JD" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- AI 评分弹窗 -->
    <a-modal
      v-model:open="showScoreModal"
      title="AI 评分"
      :confirm-loading="resumeStore.aiLoading"
      @ok="handleScore"
    >
      <a-form layout="vertical">
        <a-form-item label="参考 JD（可选）">
          <a-textarea v-model:value="scoreJd" :rows="3" placeholder="粘贴 JD" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 评分结果弹窗 -->
    <a-modal
      v-model:open="showScoreResult"
      title="评分结果"
      :footer="null"
      width="640px"
    >
      <pre class="result-pre">{{ scoreResultText }}</pre>
    </a-modal>

    <!-- JD 匹配弹窗 -->
    <a-modal
      v-model:open="showJdMatchModal"
      title="JD 匹配度分析"
      :confirm-loading="resumeStore.aiLoading"
      @ok="handleJdMatch"
    >
      <a-form layout="vertical">
        <a-form-item label="目标 JD" required>
          <a-textarea v-model:value="jdMatchInput" :rows="5" placeholder="粘贴 JD" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- JD 匹配结果弹窗 -->
    <a-modal
      v-model:open="showJdMatchResult"
      title="匹配度分析结果"
      :footer="null"
      width="640px"
    >
      <pre class="result-pre">{{ jdMatchResultText }}</pre>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal, Empty } from 'ant-design-vue'
import {
  ArrowLeftOutlined,
  SaveOutlined,
  RobotOutlined,
  DownOutlined,
  SyncOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'
import { useResumeStore } from '@/stores/resume'
import type {
  ResumeVersion,
  GenerateInput,
  PolishInput,
  ScoreResult,
  JDMatchResult,
} from '@/types/models'

const route = useRoute()
const router = useRouter()
const resumeStore = useResumeStore()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

// 解析路由参数 id
const resumeId = computed(() => Number(route.params.id))

// 可编辑的简历名称
const editableName = ref('')
// 可编辑的内容（与 currentVersion.content 双向同步）
const editableContent = ref('')

// 当 currentVersion 变化时同步到 editableContent
watch(
  () => resumeStore.currentVersion,
  (v) => {
    editableContent.value = v?.content || ''
  }
)

// 内容统计
const contentStats = computed(() => {
  const text = editableContent.value || ''
  return {
    chars: text.length,
    lines: text ? text.split('\n').length : 0,
  }
})

// 保存版本弹窗
const showSaveVersionModal = ref(false)
const savingVersion = ref(false)
const newVersionNote = ref('')

// AI 生成弹窗
const showGenerateModal = ref(false)
const generateForm = reactive<GenerateInput>({
  target_jd: '',
  scene: '',
  module_hints: '',
})
let generateAbort: AbortController | null = null

// AI 润色弹窗
const showPolishModal = ref(false)
const polishForm = reactive<PolishInput>({
  module: 'work',
  jd: '',
})

// AI 评分
const showScoreModal = ref(false)
const scoreJd = ref('')
const showScoreResult = ref(false)
const scoreResultText = ref('')

// JD 匹配
const showJdMatchModal = ref(false)
const jdMatchInput = ref('')
const showJdMatchResult = ref(false)
const jdMatchResultText = ref('')

// 返回列表
const backToList = () => {
  resumeStore.clearCurrent()
  router.push('/app/resumes')
}

// 刷新版本
const refreshVersions = async () => {
  await resumeStore.fetchVersions(resumeId.value)
}

// 选择版本
const selectVersion = (v: ResumeVersion) => {
  resumeStore.setCurrentVersion(v)
}

// 名称保存
const handleNameSave = async () => {
  if (!resumeStore.currentResume) return
  const newName = editableName.value.trim()
  if (!newName) {
    message.warning('简历名称不能为空')
    editableName.value = resumeStore.currentResume.name
    return
  }
  if (newName === resumeStore.currentResume.name) return
  await resumeStore.update(resumeId.value, { name: newName })
}

// 保存版本
const handleSaveVersion = async () => {
  if (!newVersionNote.value.trim()) {
    message.warning('请输入版本备注')
    return
  }
  savingVersion.value = true
  const v = await resumeStore.createVersion(resumeId.value, {
    content: editableContent.value,
    change_note: newVersionNote.value.trim(),
  })
  savingVersion.value = false
  if (v) {
    showSaveVersionModal.value = false
    newVersionNote.value = ''
  }
}

// 回滚
const handleRollback = (versionId: number) => {
  Modal.confirm({
    title: '确认回滚',
    content: '将基于该版本创建新版本，不会删除中间版本。是否继续？',
    onOk: async () => {
      await resumeStore.rollbackVersion(resumeId.value, versionId)
    },
  })
}

// 同步档案
const handleSyncProfile = async () => {
  await resumeStore.syncProfile(resumeId.value)
}

// AI 菜单点击
const handleAiMenuClick = (key: string) => {
  if (key === 'generate') {
    resumeStore.clearStreaming()
    showGenerateModal.value = true
  } else if (key === 'polish') {
    showPolishModal.value = true
  } else if (key === 'score') {
    showScoreModal.value = true
  } else if (key === 'jdMatch') {
    showJdMatchModal.value = true
  }
}

// AI 生成
const handleGenerate = async () => {
  generateAbort = new AbortController()
  const input: GenerateInput = {
    target_jd: generateForm.target_jd?.trim() || undefined,
    scene: generateForm.scene || undefined,
    module_hints: generateForm.module_hints?.trim() || undefined,
  }
  const ok = await resumeStore.aiGenerate(resumeId.value, input, undefined, generateAbort.signal)
  if (ok) {
    showGenerateModal.value = false
    generateForm.target_jd = ''
    generateForm.scene = ''
    generateForm.module_hints = ''
  }
}

const handleGenerateCancel = () => {
  generateAbort?.abort()
  showGenerateModal.value = false
  resumeStore.clearStreaming()
}

// AI 润色
const handlePolish = async () => {
  const input: PolishInput = {
    module: polishForm.module,
    jd: polishForm.jd?.trim() || undefined,
  }
  const v = await resumeStore.aiPolish(resumeId.value, input)
  if (v) {
    showPolishModal.value = false
    polishForm.jd = ''
  }
}

// AI 评分
const handleScore = async () => {
  const result = await resumeStore.aiScore(resumeId.value, scoreJd.value.trim() || undefined)
  if (result) {
    scoreResultText.value = formatResult(result)
    showScoreResult.value = true
    showScoreModal.value = false
  }
}

// JD 匹配
const handleJdMatch = async () => {
  if (!jdMatchInput.value.trim()) {
    message.warning('请输入目标 JD')
    return
  }
  const result = await resumeStore.aiJdMatch(resumeId.value, jdMatchInput.value.trim())
  if (result) {
    jdMatchResultText.value = formatResult(result)
    showJdMatchResult.value = true
    showJdMatchModal.value = false
  }
}

// 格式化结果对象为可读文本
const formatResult = (data: ScoreResult | JDMatchResult | null): string => {
  if (!data) return '（空结果）'
  try {
    return JSON.stringify(data, null, 2)
  } catch {
    return String(data)
  }
}

// 日期格式化
const formatVersionDate = (dateStr: string): string => {
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

// 初始化加载
onMounted(async () => {
  if (!resumeId.value || isNaN(resumeId.value)) {
    message.error('无效的简历 ID')
    router.push('/app/resumes')
    return
  }
  await resumeStore.fetchOne(resumeId.value)
  if (resumeStore.currentResume) {
    editableName.value = resumeStore.currentResume.name
  }
  await resumeStore.fetchVersions(resumeId.value)
  // 若没有 currentVersion，取首个
  if (!resumeStore.currentVersion && resumeStore.versions.length > 0) {
    resumeStore.setCurrentVersion(resumeStore.versions[0])
  }
})
</script>

<style scoped>
.resume-editor-page {
  width: 100%;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-input {
  width: 280px;
  font-size: 15px;
  font-weight: 600;
}

.editor-body {
  margin-bottom: 16px;
}

.version-card {
  border-radius: 8px;
}

.version-item {
  cursor: pointer;
  padding: 10px 12px !important;
  border-radius: 6px;
  transition: all 0.2s;
  margin-bottom: 6px !important;
  border: 1px solid transparent !important;
}

.version-item:hover {
  background: #f5f5f5;
}

.version-item.active {
  background: #e6f7ff;
  border-color: #91d5ff !important;
}

.version-item-content {
  width: 100%;
}

.version-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.version-time {
  color: #999;
  font-size: 12px;
}

.version-note {
  color: #333;
  font-size: 13px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.version-actions {
  display: flex;
  justify-content: flex-end;
}

.content-card {
  border-radius: 8px;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.content-meta {
  color: #999;
  font-size: 12px;
  font-weight: normal;
}

.content-textarea {
  font-family: 'SF Mono', Monaco, Consolas, 'Liberation Mono', monospace;
  font-size: 13px;
  resize: vertical;
}

.streaming-box {
  margin-top: 12px;
  padding: 12px;
  background: #fafafa;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  max-height: 360px;
  overflow-y: auto;
}

.streaming-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #1890ff;
  font-size: 13px;
}

.streaming-text pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.result-pre {
  background: #fafafa;
  padding: 16px;
  border-radius: 6px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 480px;
  overflow-y: auto;
  margin: 0;
}
</style>
