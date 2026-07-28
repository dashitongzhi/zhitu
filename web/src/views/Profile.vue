<template src="./templates/Profile.html"></template>

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

<style scoped src="./styles/profile.css"></style>
