<template>
  <div class="profile-page">
    <a-spin :spinning="profileStore.loading">
      <!-- 顶部完成度 -->
      <a-card class="completion-card" :bordered="false">
        <div class="completion-header">
          <span class="completion-title">档案完成度</span>
          <span class="completion-value">{{ profileStore.completion }}%</span>
        </div>
        <a-progress
          :percent="profileStore.completion"
          :stroke-color="completionColor"
          size="large"
        />
      </a-card>

      <!-- 简历解析上传 -->
      <a-card title="简历解析" class="section-card" :bordered="false">
        <a-upload
          :before-upload="handleParseResume"
          :show-upload-list="false"
          accept=".pdf,.doc,.docx"
        >
          <a-button>
            <UploadOutlined /> 上传简历自动解析
          </a-button>
        </a-upload>
        <span class="upload-tip">支持 PDF / DOC / DOCX，解析后自动合并到档案</span>
      </a-card>

      <!-- 基础信息 -->
      <a-card title="基础信息" class="section-card" :bordered="false">
        <a-form
          :model="basicForm"
          layout="vertical"
          ref="basicFormRef"
        >
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="真实姓名" name="real_name">
                <a-input v-model:value="basicForm.real_name" placeholder="请输入" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="性别" name="gender">
                <a-radio-group v-model:value="basicForm.gender">
                  <a-radio value="male">男</a-radio>
                  <a-radio value="female">女</a-radio>
                  <a-radio value="other">其他</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="出生日期" name="birth_date">
                <a-date-picker
                  v-model:value="basicForm.birth_date"
                  placeholder="选择日期"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="手机号" name="phone">
                <a-input v-model:value="basicForm.phone" placeholder="请输入" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="目标职位" name="target_position">
                <a-input v-model:value="basicForm.target_position" placeholder="请输入" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="目标城市" name="target_city">
                <a-input v-model:value="basicForm.target_city" placeholder="请输入" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="期望薪资" name="expected_salary">
                <a-input v-model:value="basicForm.expected_salary" placeholder="例如 20k-30k" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="求职状态" name="job_status">
                <a-select v-model:value="basicForm.job_status" placeholder="请选择" allow-clear>
                  <a-select-option value="fresh">应届</a-select-option>
                  <a-select-option value="graduated">已毕业</a-select-option>
                  <a-select-option value="employed">在职</a-select-option>
                  <a-select-option value="resigned">离职</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>

          <a-form-item label="自我介绍" name="self_introduction">
            <a-textarea
              v-model:value="basicForm.self_introduction"
              :rows="4"
              placeholder="请输入自我介绍"
            />
          </a-form-item>

          <div class="form-actions">
            <a-button type="primary" :loading="profileStore.loading" @click="handleSaveBasic">
              保存基础信息
            </a-button>
          </div>
        </a-form>
      </a-card>

      <!-- 6 类子资源 Tabs -->
      <a-card :bordered="false" class="section-card">
        <a-tabs v-model:activeKey="activeTab">
          <!-- 教育背景 -->
          <a-tab-pane key="educations" tab="教育背景">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showEducationModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="educationColumns"
              :data-source="profileStore.profile?.educations || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showEducationModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeleteEducation(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <!-- 工作经历 -->
          <a-tab-pane key="works" tab="工作经历">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showWorkModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="workColumns"
              :data-source="profileStore.profile?.works || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showWorkModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeleteWork(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <!-- 项目经历 -->
          <a-tab-pane key="projects" tab="项目经历">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showProjectModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="projectColumns"
              :data-source="profileStore.profile?.projects || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showProjectModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeleteProject(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <!-- 技能 -->
          <a-tab-pane key="skills" tab="技能">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showSkillModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="skillColumns"
              :data-source="profileStore.profile?.skills || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showSkillModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeleteSkill(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <!-- 荣誉奖项 -->
          <a-tab-pane key="honors" tab="荣誉奖项">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showHonorModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="honorColumns"
              :data-source="profileStore.profile?.honors || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showHonorModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeleteHonor(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <!-- 校内外实践 -->
          <a-tab-pane key="practices" tab="校内外实践">
            <div class="tab-toolbar">
              <a-button type="primary" size="small" @click="showPracticeModal()">
                <PlusOutlined /> 新增
              </a-button>
            </div>
            <a-table
              :columns="practiceColumns"
              :data-source="profileStore.profile?.practices || []"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'action'">
                  <a-button type="link" size="small" @click="showPracticeModal(record)">编辑</a-button>
                  <a-button type="link" size="small" danger @click="handleDeletePractice(record.id)">删除</a-button>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </a-spin>

    <!-- 教育背景弹窗 -->
    <a-modal
      v-model:open="educationModalVisible"
      :title="editingEducation ? '编辑教育背景' : '新增教育背景'"
      width="640px"
      @ok="handleSaveEducation"
    >
      <a-form :model="educationForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="学校" required>
              <a-input v-model:value="educationForm.school" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="专业" required>
              <a-input v-model:value="educationForm.major" placeholder="请输入" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="学历">
              <a-input v-model:value="educationForm.degree" placeholder="如 本科" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="开始日期">
              <a-date-picker v-model:value="educationForm.start_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="结束日期">
              <a-date-picker v-model:value="educationForm.end_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="GPA">
              <a-input v-model:value="educationForm.gpa" placeholder="如 3.8/4.0" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="主修课程">
          <a-textarea v-model:value="educationForm.courses" :rows="2" placeholder="可选" />
        </a-form-item>
        <a-form-item label="交换经历">
          <a-textarea v-model:value="educationForm.exchange" :rows="2" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 工作经历弹窗 -->
    <a-modal
      v-model:open="workModalVisible"
      :title="editingWork ? '编辑工作经历' : '新增工作经历'"
      width="640px"
      @ok="handleSaveWork"
    >
      <a-form :model="workForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公司" required>
              <a-input v-model:value="workForm.company" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="职位" required>
              <a-input v-model:value="workForm.position" placeholder="请输入" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开始日期">
              <a-date-picker v-model:value="workForm.start_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束日期">
              <a-date-picker v-model:value="workForm.end_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="工作描述">
          <a-textarea v-model:value="workForm.description" :rows="3" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="离职原因">
          <a-input v-model:value="workForm.leave_reason" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 项目经历弹窗 -->
    <a-modal
      v-model:open="projectModalVisible"
      :title="editingProject ? '编辑项目经历' : '新增项目经历'"
      width="640px"
      @ok="handleSaveProject"
    >
      <a-form :model="projectForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="项目名称" required>
              <a-input v-model:value="projectForm.name" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="担任角色">
              <a-input v-model:value="projectForm.role" placeholder="如 前端负责人" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开始日期">
              <a-date-picker v-model:value="projectForm.start_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束日期">
              <a-date-picker v-model:value="projectForm.end_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="项目描述">
          <a-textarea v-model:value="projectForm.description" :rows="3" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="技术栈">
          <a-input v-model:value="projectForm.tech_stack" placeholder="如 Go,React,MySQL" />
        </a-form-item>
        <a-form-item label="项目链接">
          <a-input v-model:value="projectForm.url" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 技能弹窗 -->
    <a-modal
      v-model:open="skillModalVisible"
      :title="editingSkill ? '编辑技能' : '新增技能'"
      width="480px"
      @ok="handleSaveSkill"
    >
      <a-form :model="skillForm" layout="vertical">
        <a-form-item label="分类" required>
          <a-input v-model:value="skillForm.category" placeholder="如 后端开发" />
        </a-form-item>
        <a-form-item label="技能名称" required>
          <a-input v-model:value="skillForm.name" placeholder="如 Go" />
        </a-form-item>
        <a-form-item label="熟练度">
          <a-select v-model:value="skillForm.proficiency" placeholder="请选择">
            <a-select-option value="beginner">入门</a-select-option>
            <a-select-option value="familiar">熟悉</a-select-option>
            <a-select-option value="proficient">熟练</a-select-option>
            <a-select-option value="expert">精通</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 荣誉奖项弹窗 -->
    <a-modal
      v-model:open="honorModalVisible"
      :title="editingHonor ? '编辑荣誉奖项' : '新增荣誉奖项'"
      width="480px"
      @ok="handleSaveHonor"
    >
      <a-form :model="honorForm" layout="vertical">
        <a-form-item label="奖项名称" required>
          <a-input v-model:value="honorForm.name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="颁发机构">
          <a-input v-model:value="honorForm.issuer" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="获奖日期">
          <a-date-picker v-model:value="honorForm.award_date" placeholder="选择" style="width: 100%" />
        </a-form-item>
        <a-form-item label="级别">
          <a-input v-model:value="honorForm.level" placeholder="如 国家级、省级" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 校内外实践弹窗 -->
    <a-modal
      v-model:open="practiceModalVisible"
      :title="editingPractice ? '编辑实践经历' : '新增实践经历'"
      width="640px"
      @ok="handleSavePractice"
    >
      <a-form :model="practiceForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="实践标题" required>
              <a-input v-model:value="practiceForm.title" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="组织">
              <a-input v-model:value="practiceForm.organization" placeholder="请输入" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开始日期">
              <a-date-picker v-model:value="practiceForm.start_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束日期">
              <a-date-picker v-model:value="practiceForm.end_date" placeholder="选择" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="实践描述">
          <a-textarea v-model:value="practiceForm.description" :rows="3" placeholder="请输入" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useProfileStore } from '@/stores/profile'
import { message, Modal } from 'ant-design-vue'
import { UploadOutlined, PlusOutlined } from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'
import type {
  Education,
  Work,
  Project,
  Skill,
  Honor,
  Practice,
} from '@/types/models'
import dayjs, { Dayjs } from 'dayjs'

const profileStore = useProfileStore()
const basicFormRef = ref<FormInstance>()

// 当前激活的 Tab
const activeTab = ref('educations')

// 完成度颜色
const completionColor = computed(() => {
  const v = profileStore.completion
  if (v >= 80) return '#52c41a'
  if (v >= 50) return '#faad14'
  return '#1890ff'
})

// 基础信息表单
const basicForm = reactive<{
  real_name: string
  gender: string
  birth_date: Dayjs | null
  phone: string
  target_position: string
  target_city: string
  expected_salary: string
  job_status: string
  self_introduction: string
}>({
  real_name: '',
  gender: '',
  birth_date: null,
  phone: '',
  target_position: '',
  target_city: '',
  expected_salary: '',
  job_status: '',
  self_introduction: '',
})

// 同步 profile 到 basicForm
const syncBasicForm = () => {
  const p = profileStore.profile
  if (!p) return
  basicForm.real_name = p.real_name || ''
  basicForm.gender = p.gender || ''
  basicForm.birth_date = p.birth_date ? dayjs(p.birth_date) : null
  basicForm.phone = p.phone || ''
  basicForm.target_position = p.target_position || ''
  basicForm.target_city = p.target_city || ''
  basicForm.expected_salary = p.expected_salary || ''
  basicForm.job_status = p.job_status || ''
  basicForm.self_introduction = p.self_introduction || ''
}

// 页面加载
onMounted(async () => {
  await profileStore.fetchProfile()
  syncBasicForm()
})

// 简历解析上传
const handleParseResume = async (file: File) => {
  await profileStore.parseResume(file)
  syncBasicForm()
  return false // 阻止自动上传
}

// 保存基础信息
const handleSaveBasic = async () => {
  const data = {
    real_name: basicForm.real_name,
    gender: basicForm.gender as 'male' | 'female' | 'other' | undefined,
    birth_date: basicForm.birth_date?.format('YYYY-MM-DD'),
    phone: basicForm.phone,
    target_position: basicForm.target_position,
    target_city: basicForm.target_city,
    expected_salary: basicForm.expected_salary,
    job_status: basicForm.job_status as 'fresh' | 'graduated' | 'employed' | 'resigned' | undefined,
    self_introduction: basicForm.self_introduction,
  }
  const ok = await profileStore.updateProfile(data)
  if (ok) syncBasicForm()
}

// ==================== 表格列定义 ====================

const educationColumns = [
  { title: '学校', dataIndex: 'school', key: 'school' },
  { title: '专业', dataIndex: 'major', key: 'major' },
  { title: '学历', dataIndex: 'degree', key: 'degree', width: 80 },
  { title: '开始', dataIndex: 'start_date', key: 'start_date', width: 110 },
  { title: '结束', dataIndex: 'end_date', key: 'end_date', width: 110 },
  { title: 'GPA', dataIndex: 'gpa', key: 'gpa', width: 80 },
  { title: '操作', key: 'action', width: 130 },
]

const workColumns = [
  { title: '公司', dataIndex: 'company', key: 'company' },
  { title: '职位', dataIndex: 'position', key: 'position' },
  { title: '开始', dataIndex: 'start_date', key: 'start_date', width: 110 },
  { title: '结束', dataIndex: 'end_date', key: 'end_date', width: 110 },
  { title: '操作', key: 'action', width: 130 },
]

const projectColumns = [
  { title: '项目名称', dataIndex: 'name', key: 'name' },
  { title: '角色', dataIndex: 'role', key: 'role', width: 120 },
  { title: '开始', dataIndex: 'start_date', key: 'start_date', width: 110 },
  { title: '结束', dataIndex: 'end_date', key: 'end_date', width: 110 },
  { title: '技术栈', dataIndex: 'tech_stack', key: 'tech_stack' },
  { title: '操作', key: 'action', width: 130 },
]

const skillColumns = [
  { title: '分类', dataIndex: 'category', key: 'category' },
  { title: '技能', dataIndex: 'name', key: 'name' },
  { title: '熟练度', dataIndex: 'proficiency', key: 'proficiency', width: 100 },
  { title: '操作', key: 'action', width: 130 },
]

const honorColumns = [
  { title: '奖项名称', dataIndex: 'name', key: 'name' },
  { title: '颁发机构', dataIndex: 'issuer', key: 'issuer' },
  { title: '获奖日期', dataIndex: 'award_date', key: 'award_date', width: 120 },
  { title: '级别', dataIndex: 'level', key: 'level', width: 100 },
  { title: '操作', key: 'action', width: 130 },
]

const practiceColumns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '组织', dataIndex: 'organization', key: 'organization' },
  { title: '开始', dataIndex: 'start_date', key: 'start_date', width: 110 },
  { title: '结束', dataIndex: 'end_date', key: 'end_date', width: 110 },
  { title: '操作', key: 'action', width: 130 },
]

// ==================== 教育背景 ====================

const educationModalVisible = ref(false)
const editingEducation = ref<Education | null>(null)
const educationForm = reactive({
  school: '',
  major: '',
  degree: '',
  start_date: null as Dayjs | null,
  end_date: null as Dayjs | null,
  gpa: '',
  courses: '',
  exchange: '',
})

const showEducationModal = (item?: Education) => {
  if (item) {
    editingEducation.value = item
    Object.assign(educationForm, {
      school: item.school,
      major: item.major,
      degree: item.degree,
      start_date: item.start_date ? dayjs(item.start_date) : null,
      end_date: item.end_date ? dayjs(item.end_date) : null,
      gpa: item.gpa,
      courses: item.courses,
      exchange: item.exchange,
    })
  } else {
    editingEducation.value = null
    Object.assign(educationForm, {
      school: '', major: '', degree: '', start_date: null, end_date: null,
      gpa: '', courses: '', exchange: '',
    })
  }
  educationModalVisible.value = true
}

const handleSaveEducation = async () => {
  if (!educationForm.school || !educationForm.major) {
    message.warning('请填写学校和专业')
    return
  }
  const data = {
    school: educationForm.school,
    major: educationForm.major,
    degree: educationForm.degree,
    start_date: educationForm.start_date?.format('YYYY-MM-DD') || '',
    end_date: educationForm.end_date?.format('YYYY-MM-DD') || '',
    gpa: educationForm.gpa,
    courses: educationForm.courses,
    exchange: educationForm.exchange,
  }
  if (editingEducation.value) {
    await profileStore.updateEducation(editingEducation.value.id!, data)
  } else {
    await profileStore.createEducation(data)
  }
  educationModalVisible.value = false
}

const handleDeleteEducation = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该教育背景吗？',
    onOk: async () => {
      await profileStore.deleteEducation(id)
    },
  })
}

// ==================== 工作经历 ====================

const workModalVisible = ref(false)
const editingWork = ref<Work | null>(null)
const workForm = reactive({
  company: '',
  position: '',
  start_date: null as Dayjs | null,
  end_date: null as Dayjs | null,
  description: '',
  leave_reason: '',
})

const showWorkModal = (item?: Work) => {
  if (item) {
    editingWork.value = item
    Object.assign(workForm, {
      company: item.company,
      position: item.position,
      start_date: item.start_date ? dayjs(item.start_date) : null,
      end_date: item.end_date ? dayjs(item.end_date) : null,
      description: item.description,
      leave_reason: item.leave_reason,
    })
  } else {
    editingWork.value = null
    Object.assign(workForm, {
      company: '', position: '', start_date: null, end_date: null,
      description: '', leave_reason: '',
    })
  }
  workModalVisible.value = true
}

const handleSaveWork = async () => {
  if (!workForm.company || !workForm.position) {
    message.warning('请填写公司和职位')
    return
  }
  const data = {
    company: workForm.company,
    position: workForm.position,
    start_date: workForm.start_date?.format('YYYY-MM-DD') || '',
    end_date: workForm.end_date?.format('YYYY-MM-DD') || '',
    description: workForm.description,
    leave_reason: workForm.leave_reason,
  }
  if (editingWork.value) {
    await profileStore.updateWork(editingWork.value.id!, data)
  } else {
    await profileStore.createWork(data)
  }
  workModalVisible.value = false
}

const handleDeleteWork = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该工作经历吗？',
    onOk: async () => {
      await profileStore.deleteWork(id)
    },
  })
}

// ==================== 项目经历 ====================

const projectModalVisible = ref(false)
const editingProject = ref<Project | null>(null)
const projectForm = reactive({
  name: '',
  role: '',
  start_date: null as Dayjs | null,
  end_date: null as Dayjs | null,
  description: '',
  tech_stack: '',
  url: '',
})

const showProjectModal = (item?: Project) => {
  if (item) {
    editingProject.value = item
    Object.assign(projectForm, {
      name: item.name,
      role: item.role,
      start_date: item.start_date ? dayjs(item.start_date) : null,
      end_date: item.end_date ? dayjs(item.end_date) : null,
      description: item.description,
      tech_stack: item.tech_stack,
      url: item.url,
    })
  } else {
    editingProject.value = null
    Object.assign(projectForm, {
      name: '', role: '', start_date: null, end_date: null,
      description: '', tech_stack: '', url: '',
    })
  }
  projectModalVisible.value = true
}

const handleSaveProject = async () => {
  if (!projectForm.name) {
    message.warning('请填写项目名称')
    return
  }
  const data = {
    name: projectForm.name,
    role: projectForm.role,
    start_date: projectForm.start_date?.format('YYYY-MM-DD') || '',
    end_date: projectForm.end_date?.format('YYYY-MM-DD') || '',
    description: projectForm.description,
    tech_stack: projectForm.tech_stack,
    url: projectForm.url,
  }
  if (editingProject.value) {
    await profileStore.updateProject(editingProject.value.id!, data)
  } else {
    await profileStore.createProject(data)
  }
  projectModalVisible.value = false
}

const handleDeleteProject = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该项目经历吗？',
    onOk: async () => {
      await profileStore.deleteProject(id)
    },
  })
}

// ==================== 技能 ====================

const skillModalVisible = ref(false)
const editingSkill = ref<Skill | null>(null)
const skillForm = reactive({
  category: '',
  name: '',
  proficiency: 'familiar',
})

const showSkillModal = (item?: Skill) => {
  if (item) {
    editingSkill.value = item
    Object.assign(skillForm, {
      category: item.category,
      name: item.name,
      proficiency: item.proficiency,
    })
  } else {
    editingSkill.value = null
    Object.assign(skillForm, { category: '', name: '', proficiency: 'familiar' })
  }
  skillModalVisible.value = true
}

const handleSaveSkill = async () => {
  if (!skillForm.category || !skillForm.name) {
    message.warning('请填写分类和技能名称')
    return
  }
  const data = {
    category: skillForm.category,
    name: skillForm.name,
    proficiency: skillForm.proficiency,
  }
  if (editingSkill.value) {
    await profileStore.updateSkill(editingSkill.value.id!, data)
  } else {
    await profileStore.createSkill(data)
  }
  skillModalVisible.value = false
}

const handleDeleteSkill = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该技能吗？',
    onOk: async () => {
      await profileStore.deleteSkill(id)
    },
  })
}

// ==================== 荣誉奖项 ====================

const honorModalVisible = ref(false)
const editingHonor = ref<Honor | null>(null)
const honorForm = reactive({
  name: '',
  issuer: '',
  award_date: null as Dayjs | null,
  level: '',
})

const showHonorModal = (item?: Honor) => {
  if (item) {
    editingHonor.value = item
    Object.assign(honorForm, {
      name: item.name,
      issuer: item.issuer,
      award_date: item.award_date ? dayjs(item.award_date) : null,
      level: item.level,
    })
  } else {
    editingHonor.value = null
    Object.assign(honorForm, { name: '', issuer: '', award_date: null, level: '' })
  }
  honorModalVisible.value = true
}

const handleSaveHonor = async () => {
  if (!honorForm.name) {
    message.warning('请填写奖项名称')
    return
  }
  const data = {
    name: honorForm.name,
    issuer: honorForm.issuer,
    award_date: honorForm.award_date?.format('YYYY-MM-DD') || '',
    level: honorForm.level,
  }
  if (editingHonor.value) {
    await profileStore.updateHonor(editingHonor.value.id!, data)
  } else {
    await profileStore.createHonor(data)
  }
  honorModalVisible.value = false
}

const handleDeleteHonor = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该荣誉奖项吗？',
    onOk: async () => {
      await profileStore.deleteHonor(id)
    },
  })
}

// ==================== 校内外实践 ====================

const practiceModalVisible = ref(false)
const editingPractice = ref<Practice | null>(null)
const practiceForm = reactive({
  title: '',
  organization: '',
  start_date: null as Dayjs | null,
  end_date: null as Dayjs | null,
  description: '',
})

const showPracticeModal = (item?: Practice) => {
  if (item) {
    editingPractice.value = item
    Object.assign(practiceForm, {
      title: item.title,
      organization: item.organization,
      start_date: item.start_date ? dayjs(item.start_date) : null,
      end_date: item.end_date ? dayjs(item.end_date) : null,
      description: item.description,
    })
  } else {
    editingPractice.value = null
    Object.assign(practiceForm, {
      title: '', organization: '', start_date: null, end_date: null, description: '',
    })
  }
  practiceModalVisible.value = true
}

const handleSavePractice = async () => {
  if (!practiceForm.title) {
    message.warning('请填写实践标题')
    return
  }
  const data = {
    title: practiceForm.title,
    organization: practiceForm.organization,
    start_date: practiceForm.start_date?.format('YYYY-MM-DD') || '',
    end_date: practiceForm.end_date?.format('YYYY-MM-DD') || '',
    description: practiceForm.description,
  }
  if (editingPractice.value) {
    await profileStore.updatePractice(editingPractice.value.id!, data)
  } else {
    await profileStore.createPractice(data)
  }
  practiceModalVisible.value = false
}

const handleDeletePractice = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该实践经历吗？',
    onOk: async () => {
      await profileStore.deletePractice(id)
    },
  })
}
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.completion-card {
  border-radius: 8px;
}

.completion-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.completion-title {
  font-size: 16px;
  font-weight: 500;
  color: #262626;
}

.completion-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
}

.section-card {
  border-radius: 8px;
}

.upload-tip {
  margin-left: 12px;
  color: #999;
  font-size: 13px;
}

.form-actions {
  margin-top: 8px;
}

.tab-toolbar {
  margin-bottom: 12px;
}
</style>
