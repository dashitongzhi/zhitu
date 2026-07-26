<template>
  <div class="resume-list-page">
    <!-- 顶部工具栏 -->
    <div class="page-header">
      <div class="header-info">
        <h2 class="page-title">简历管理</h2>
        <p class="page-desc">支持多份简历、版本管理、AI 生成与评分</p>
      </div>
      <a-button type="primary" @click="router.push('/app/resumes/new')">
        <PlusOutlined /> 创建简历
      </a-button>
    </div>

    <!-- 简历卡片网格 -->
    <a-spin :spinning="resumeStore.loading" tip="加载中...">
      <a-empty
        v-if="!resumeStore.loading && resumeStore.resumes.length === 0"
        description="暂无简历，点击右上角创建"
        class="empty-state"
      />

      <a-row v-else :gutter="[16, 16]">
        <a-col
          v-for="resume in resumeStore.resumes"
          :key="resume.id"
          :xs="24"
          :sm="12"
          :md="12"
          :lg="8"
          :xl="6"
        >
          <a-card class="resume-card" hoverable @click="enterEditor(resume.id)">
            <div class="card-header">
              <div class="card-title">
                <FileTextOutlined class="card-icon" />
                <span class="title-text" :title="resume.name">{{ resume.name }}</span>
              </div>
              <a-dropdown :trigger="['click']" @click.stop>
                <MoreOutlined class="more-btn" @click.stop />
                <template #overlay>
                  <a-menu @click="(e) => handleMenuClick(e.key, resume)">
                    <a-menu-item key="edit">进入编辑</a-menu-item>
                    <a-menu-item key="delete" danger>删除简历</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>

            <div class="card-body">
              <div class="info-row">
                <span class="info-label">目标公司</span>
                <span class="info-value">{{ resume.target_company || '未指定' }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">目标职位</span>
                <span class="info-value">{{ resume.target_position || '未指定' }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">生成场景</span>
                <a-tag :color="sceneColor(resume.scene)">
                  {{ sceneLabel(resume.scene) }}
                </a-tag>
              </div>
            </div>

            <div class="card-footer">
              <span class="footer-time">
                <ClockCircleOutlined />
                {{ formatDate(resume.updated_at) }}
              </span>
              <a-button type="link" size="small" @click.stop="enterEditor(resume.id)">
                编辑 <ArrowRightOutlined />
              </a-button>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>

  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  FileTextOutlined,
  MoreOutlined,
  ClockCircleOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons-vue'
import { useResumeStore } from '@/stores/resume'
import type { Resume, ResumeScene } from '@/types/models'

const router = useRouter()
const resumeStore = useResumeStore()

// 场景标签
const sceneLabel = (scene: ResumeScene | string): string => {
  const map: Record<string, string> = {
    manual: '手动编辑',
    jd: '基于 JD',
    scenario: '场景化',
  }
  return map[scene] || scene
}

const sceneColor = (scene: ResumeScene | string): string => {
  const map: Record<string, string> = {
    manual: 'default',
    jd: 'blue',
    scenario: 'purple',
  }
  return map[scene] || 'default'
}

// 日期格式化
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

// 进入编辑器
const enterEditor = (id: number) => {
  router.push(`/app/resumes/${id}`)
}

// 菜单点击
const handleMenuClick = (key: string, resume: Resume) => {
  if (key === 'edit') {
    enterEditor(resume.id)
  } else if (key === 'delete') {
    Modal.confirm({
      title: '确认删除',
      content: `确定删除简历「${resume.name}」吗？所有版本将一并删除，且不可恢复。`,
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        await resumeStore.remove(resume.id)
      },
    })
  }
}

onMounted(() => {
  resumeStore.fetchList()
})
</script>

<style scoped>
.resume-list-page {
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

.resume-card {
  border-radius: 8px;
  transition: all 0.3s;
  cursor: pointer;
}

.resume-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.card-icon {
  color: #1890ff;
  font-size: 18px;
}

.title-text {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a2e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.more-btn {
  font-size: 16px;
  color: #999;
  cursor: pointer;
  padding: 4px;
}

.more-btn:hover {
  color: #1890ff;
}

.card-body {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.info-label {
  color: #999;
}

.info-value {
  color: #333;
  max-width: 60%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px dashed #f0f0f0;
}

.footer-time {
  color: #999;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
