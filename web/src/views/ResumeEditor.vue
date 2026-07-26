<template>
  <div class="resume-lab" :class="{ 'is-demo': demoMode }">
    <header class="lab-toolbar">
      <div class="toolbar-leading">
        <button class="icon-button" type="button" aria-label="返回简历列表" @click="backToList">
          <ArrowLeftOutlined />
        </button>
        <div class="brand-mark"><span>职</span></div>
        <div class="resume-title-wrap">
          <input v-model="editableName" class="resume-name" @blur="handleNameSave" />
          <span class="save-state"><CheckCircleOutlined /> {{ demoMode ? '本地预览模式' : '内容实时预览' }}</span>
        </div>
      </div>

      <div class="toolbar-controls">
        <label class="toolbar-select">
          <LayoutOutlined />
          <select v-model="templateStyle">
            <option value="classic">经典单栏</option>
            <option value="modern">现代强调</option>
          </select>
        </label>
        <label class="toolbar-select">
          <FontSizeOutlined />
          <select v-model="fontFamily">
            <option value="serif">宋体正文</option>
            <option value="sans">黑体正文</option>
          </select>
        </label>
        <label class="toolbar-select compact">
          <ColumnHeightOutlined />
          <select v-model="density">
            <option value="compact">紧凑</option>
            <option value="comfortable">舒展</option>
          </select>
        </label>
        <div class="zoom-control">
          <button type="button" @click="zoom = Math.max(60, zoom - 5)">−</button>
          <span>{{ zoom }}%</span>
          <button type="button" @click="zoom = Math.min(110, zoom + 5)">＋</button>
        </div>
        <a-button @click="versionDrawerOpen = true"><HistoryOutlined /> 版本</a-button>
        <a-button @click="showSaveVersionModal = true"><SaveOutlined /> 保存版本</a-button>
        <a-button class="export-button" type="primary" @click="exportResume">
          <DownloadOutlined /> 导出 PDF
        </a-button>
      </div>
    </header>

    <main class="lab-workspace">
      <aside class="editor-panel">
        <div class="panel-tabs">
          <button class="active" type="button">结构化编辑</button>
          <button type="button" @click="showRawJson = !showRawJson">{{ showRawJson ? '返回表单' : '查看 JSON' }}</button>
        </div>

        <div v-if="showRawJson" class="raw-json-panel">
          <div class="panel-hint">当前结构化数据会同步生成版本 JSON</div>
          <pre>{{ serializedContent }}</pre>
        </div>

        <div v-else class="editor-scroll">
          <div class="panel-hint">填写内容后，右侧 A4 简历会即时更新</div>

          <section class="form-section" :class="{ collapsed: collapsed.personal }">
            <button class="section-heading" type="button" @click="toggleSection('personal')">
              <span><UserOutlined /> 基本信息</span><DownOutlined />
            </button>
            <div class="section-body grid-two">
              <label class="field full"><span>姓名</span><input v-model="resumeContent.personal.name" /></label>
              <label class="field"><span>电话</span><input v-model="resumeContent.personal.phone" /></label>
              <label class="field"><span>邮箱</span><input v-model="resumeContent.personal.email" /></label>
              <label class="field"><span>所在城市</span><input v-model="resumeContent.personal.city" /></label>
              <label class="field"><span>GitHub / 主页</span><input v-model="resumeContent.personal.github" /></label>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.intention }">
            <button class="section-heading" type="button" @click="toggleSection('intention')">
              <span><AimOutlined /> 求职意向</span><DownOutlined />
            </button>
            <div class="section-body grid-two">
              <label class="field full"><span>目标岗位</span><input v-model="resumeContent.intention.position" /></label>
              <label class="field"><span>期望城市</span><input v-model="resumeContent.intention.city" /></label>
              <label class="field"><span>期望薪资</span><input v-model="resumeContent.intention.salary" /></label>
              <label class="field"><span>到岗时间</span><input v-model="resumeContent.intention.arrival" /></label>
              <label class="field"><span>目标行业</span><input v-model="resumeContent.intention.industry" /></label>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.education }">
            <button class="section-heading" type="button" @click="toggleSection('education')">
              <span><ReadOutlined /> 教育经历</span><DownOutlined />
            </button>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.education" :key="`edu-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.education.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>学校</span><input v-model="item.school" /></label>
                  <label class="field"><span>专业</span><input v-model="item.major" /></label>
                  <label class="field"><span>学历</span><input v-model="item.degree" /></label>
                  <label class="field"><span>开始</span><input v-model="item.start" /></label>
                  <label class="field"><span>结束</span><input v-model="item.end" /></label>
                  <label class="field full"><span>主修课程 / GPA</span><input v-model="item.courses" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addEducation"><PlusOutlined /> 添加教育经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.work }">
            <button class="section-heading" type="button" @click="toggleSection('work')">
              <span><BankOutlined /> 工作经历</span><DownOutlined />
            </button>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.work" :key="`work-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.work.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>公司</span><input v-model="item.company" /></label>
                  <label class="field"><span>职位</span><input v-model="item.position" /></label>
                  <label class="field"><span>时间</span><input v-model="item.start" placeholder="2023.07 - 至今" /></label>
                  <label class="field full"><span>工作成果（每行一条）</span><textarea v-model="item.description" rows="5" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addWork"><PlusOutlined /> 添加工作经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.project }">
            <button class="section-heading" type="button" @click="toggleSection('project')">
              <span><ProjectOutlined /> 项目经历</span><DownOutlined />
            </button>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.project" :key="`project-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.project.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>项目名称</span><input v-model="item.name" /></label>
                  <label class="field"><span>项目角色</span><input v-model="item.role" /></label>
                  <label class="field"><span>时间</span><input v-model="item.start" placeholder="2024.01 - 2024.06" /></label>
                  <label class="field full"><span>项目成果（每行一条）</span><textarea v-model="item.description" rows="5" /></label>
                  <label class="field full"><span>技术栈（逗号分隔）</span><input :value="item.tech_stack.join(', ')" @input="updateTechStack(item, $event)" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addProject"><PlusOutlined /> 添加项目经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.skills }">
            <button class="section-heading" type="button" @click="toggleSection('skills')">
              <span><ToolOutlined /> 专业技能</span><DownOutlined />
            </button>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.skills" :key="`skill-${index}`" class="skill-row">
                <input v-model="item.category" placeholder="分类" />
                <input v-model="item.name" placeholder="技能内容" />
                <button type="button" @click="resumeContent.skills.splice(index, 1)"><CloseOutlined /></button>
              </article>
              <button class="add-item" type="button" @click="addSkill"><PlusOutlined /> 添加技能分类</button>
            </div>
          </section>
        </div>
      </aside>

      <section class="preview-stage">
        <div class="stage-meta">
          <span><FileTextOutlined /> A4 · 实时预览</span>
          <span>{{ contentStats.chars }} 字符 · {{ contentStats.sections }} 个内容模块</span>
        </div>
        <div class="paper-viewport">
          <article
            id="resume-paper"
            class="resume-paper"
            :class="[`template-${templateStyle}`, `font-${fontFamily}`, `density-${density}`]"
            :style="{ transform: `scale(${zoom / 100})` }"
          >
            <header class="paper-header">
              <h1>{{ resumeContent.personal.name || '你的姓名' }}</h1>
              <p class="paper-role">{{ resumeContent.intention.position || '目标职位' }}</p>
              <div class="contact-line">
                <span v-if="resumeContent.personal.phone">{{ resumeContent.personal.phone }}</span>
                <span v-if="resumeContent.personal.email">{{ resumeContent.personal.email }}</span>
                <span v-if="resumeContent.personal.city">{{ resumeContent.personal.city }}</span>
                <span v-if="resumeContent.personal.github">{{ resumeContent.personal.github }}</span>
              </div>
            </header>

            <section v-if="isVisible('intention')" class="paper-section intention-section">
              <h2>求职意向</h2>
              <p>{{ intentionSummary }}</p>
            </section>

            <section v-if="isVisible('education') && resumeContent.education.length" class="paper-section">
              <h2>教育经历</h2>
              <article v-for="(item, index) in resumeContent.education" :key="`preview-edu-${index}`" class="resume-entry">
                <div class="entry-heading">
                  <strong>{{ item.school || '学校名称' }}</strong>
                  <span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span>
                </div>
                <div class="entry-subheading"><span>{{ item.major }}</span><span>{{ item.degree }}</span></div>
                <p v-if="item.courses" class="entry-note">{{ item.courses }}</p>
              </article>
            </section>

            <section v-if="isVisible('work') && resumeContent.work.length" class="paper-section">
              <h2>工作经历</h2>
              <article v-for="(item, index) in resumeContent.work" :key="`preview-work-${index}`" class="resume-entry">
                <div class="entry-heading"><strong>{{ item.company || '公司名称' }}</strong><span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span></div>
                <div class="entry-subheading"><span>{{ item.position }}</span></div>
                <ul><li v-for="(line, lineIndex) in toLines(item.description)" :key="lineIndex">{{ line }}</li></ul>
              </article>
            </section>

            <section v-if="isVisible('project') && resumeContent.project.length" class="paper-section">
              <h2>项目经历</h2>
              <article v-for="(item, index) in resumeContent.project" :key="`preview-project-${index}`" class="resume-entry">
                <div class="entry-heading"><strong>{{ item.name || '项目名称' }}</strong><span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span></div>
                <div class="entry-subheading"><span>{{ item.role }}</span><span class="tech-stack">{{ item.tech_stack.join(' · ') }}</span></div>
                <ul><li v-for="(line, lineIndex) in toLines(item.description)" :key="lineIndex">{{ line }}</li></ul>
              </article>
            </section>

            <section v-if="isVisible('skills') && resumeContent.skills.length" class="paper-section skill-section">
              <h2>专业技能</h2>
              <p v-for="(item, index) in resumeContent.skills" :key="`preview-skill-${index}`"><strong>{{ item.category }}：</strong>{{ item.name }}</p>
            </section>
          </article>
        </div>
      </section>

      <aside class="ai-panel">
        <div class="ai-heading">
          <div><span class="ai-kicker">AI 诊断</span><h2>针对 {{ resumeContent.intention.position || '目标岗位' }}</h2></div>
          <RobotOutlined />
        </div>

        <section class="match-card">
          <span>JD 匹配度</span>
          <div class="score-line"><strong>{{ diagnosticScore }}%</strong><em>{{ diagnosticScore >= 80 ? '较好' : '待完善' }}</em></div>
          <div class="score-track"><i :style="{ width: `${diagnosticScore}%` }" /></div>
          <small>基于完整度、成果量化和关键词覆盖估算</small>
        </section>

        <section class="dimension-card">
          <h3>六维评估</h3>
          <div v-for="item in dimensions" :key="item.label" class="dimension-row">
            <span>{{ item.label }}</span><div><i :style="{ width: `${item.value}%` }" /></div><b>{{ item.value }}</b>
          </div>
        </section>

        <section class="suggestion-card">
          <h3>优化建议 <span>{{ suggestions.length }}</span></h3>
          <article v-for="(item, index) in suggestions" :key="item.title">
            <b>{{ index + 1 }}</b>
            <div><strong>{{ item.title }}</strong><p>{{ item.detail }}</p></div>
          </article>
        </section>

        <label class="jd-field">
          <span>目标 JD</span>
          <textarea v-model="targetJd" rows="5" placeholder="粘贴岗位描述，诊断会随内容更新" />
        </label>
        <button class="optimize-button" type="button" @click="handleOptimize"><RobotOutlined /> 根据 JD 一键优化</button>
        <p class="privacy-note"><SafetyCertificateOutlined /> 简历内容仅用于当前编辑与诊断</p>
      </aside>
    </main>

    <a-drawer v-model:open="versionDrawerOpen" title="版本历史" width="380">
      <a-empty v-if="!resumeStore.versions.length" description="暂无已保存版本" />
      <div v-else class="version-list">
        <button v-for="version in resumeStore.versions" :key="version.id" type="button" :class="{ active: resumeStore.currentVersion?.id === version.id }" @click="selectVersion(version)">
          <span><strong>{{ version.version_label }}</strong><small>{{ formatVersionDate(version.created_at) }}</small></span>
          <em>{{ version.change_note || '无版本备注' }}</em>
        </button>
      </div>
    </a-drawer>

    <a-modal v-model:open="showSaveVersionModal" title="保存当前版本" :confirm-loading="savingVersion" @ok="handleSaveVersion">
      <a-form layout="vertical">
        <a-form-item label="版本备注" required><a-input v-model:value="newVersionNote" placeholder="例如：针对后端工程师 JD 完成优化" /></a-form-item>
        <a-alert message="当前结构化内容会保存为一个新的简历版本" type="info" show-icon />
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  AimOutlined, ArrowLeftOutlined, BankOutlined, CheckCircleOutlined, CloseOutlined,
  ColumnHeightOutlined, DeleteOutlined, DownloadOutlined, DownOutlined, FileTextOutlined,
  FontSizeOutlined, HistoryOutlined, LayoutOutlined, PlusOutlined, ProjectOutlined,
  ReadOutlined, RobotOutlined, SafetyCertificateOutlined, SaveOutlined, ToolOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { useResumeStore } from '@/stores/resume'
import type { ResumeVersion } from '@/types/models'

interface ResumeEducation { school: string; major: string; degree: string; start: string; end: string; courses: string; gpa: string }
interface ResumeWork { company: string; position: string; start: string; end: string; description: string; leave_reason: string }
interface ResumeProject { name: string; role: string; start: string; end: string; description: string; tech_stack: string[]; url: string }
interface ResumeSkill { category: string; name: string; proficiency: string }
interface ResumeContent {
  personal: { name: string; gender: string; age: string; phone: string; email: string; github: string; avatar: string; city: string }
  intention: { position: string; city: string; salary: string; arrival: string; industry: string }
  education: ResumeEducation[]; work: ResumeWork[]; project: ResumeProject[]; skills: ResumeSkill[]
  honor: unknown[]; custom: unknown[]; module_order: string[]; module_visibility: Record<string, boolean>
}

const route = useRoute()
const router = useRouter()
const resumeStore = useResumeStore()
const demoMode = computed(() => route.name === 'ResumeLabPreview' || route.query.demo === '1')
const resumeId = computed(() => Number(route.params.id))
const editableName = ref('后端工程师－张明')
const templateStyle = ref('classic')
const fontFamily = ref('sans')
const density = ref('comfortable')
const zoom = ref(82)
const showRawJson = ref(false)
const versionDrawerOpen = ref(false)
const showSaveVersionModal = ref(false)
const savingVersion = ref(false)
const newVersionNote = ref('')
const targetJd = ref('负责高并发业务系统设计与研发；熟悉 Go、Java、MySQL、Redis、消息队列与微服务架构；具备性能优化和工程化经验。')
const collapsed = reactive<Record<string, boolean>>({ personal: false, intention: false, education: false, work: false, project: false, skills: false })
const resumeContent = reactive<ResumeContent>(createSampleContent())

const serializedContent = computed(() => JSON.stringify(resumeContent, null, 2))
const contentStats = computed(() => ({
  chars: serializedContent.value.length,
  sections: [resumeContent.education, resumeContent.work, resumeContent.project, resumeContent.skills].filter((items) => items.length).length + 2,
}))
const intentionSummary = computed(() => [resumeContent.intention.position, resumeContent.intention.city, resumeContent.intention.salary, resumeContent.intention.arrival].filter(Boolean).join(' ｜ '))
const diagnosticScore = computed(() => {
  let score = 42
  if (resumeContent.personal.phone && resumeContent.personal.email) score += 8
  score += Math.min(12, resumeContent.education.length * 6)
  score += Math.min(16, resumeContent.work.length * 8)
  score += Math.min(14, resumeContent.project.length * 7)
  if (resumeContent.skills.length >= 3) score += 8
  return Math.min(score, 96)
})
const dimensions = computed(() => [
  { label: '技能匹配', value: Math.min(96, 58 + resumeContent.skills.length * 8) },
  { label: '项目经验', value: Math.min(94, 52 + resumeContent.project.length * 14) },
  { label: '工作经历', value: Math.min(92, 54 + resumeContent.work.length * 16) },
  { label: '教育背景', value: Math.min(90, 62 + resumeContent.education.length * 12) },
  { label: '求职意向', value: resumeContent.intention.position && resumeContent.intention.city ? 88 : 48 },
  { label: '关键词覆盖', value: targetJd.value ? 76 : 42 },
])
const suggestions = computed(() => {
  const items = []
  const allDescriptions = [...resumeContent.work, ...resumeContent.project].map((item) => item.description).join(' ')
  if (!/[0-9]+[%万倍ms]/.test(allDescriptions)) items.push({ title: '补充量化成果', detail: '在项目与工作经历中加入性能、规模、效率或业务结果指标。' })
  if (resumeContent.project.length < 2) items.push({ title: '项目证据偏少', detail: '建议补充一个能体现架构设计或复杂问题解决能力的项目。' })
  if (!/Kafka|RocketMQ|消息队列/i.test(allDescriptions + resumeContent.skills.map((item) => item.name).join(' '))) items.push({ title: '关键词覆盖不足', detail: '目标 JD 强调消息队列，可在真实使用过的经历中补充对应技术。' })
  items.push({ title: '强化岗位定位', detail: `将最相关的${resumeContent.intention.position || '目标岗位'}经历放在前半页，降低招聘者理解成本。` })
  return items.slice(0, 3)
})

watch(() => resumeStore.currentVersion, (version) => {
  if (!version || demoMode.value) return
  try { Object.assign(resumeContent, normalizeContent(JSON.parse(version.content))) } catch { message.warning('当前版本内容不是可识别的结构化简历 JSON') }
})

function createSampleContent(): ResumeContent {
  return normalizeContent({
    personal: { name: '张明', phone: '138-0000-0000', email: 'zhangming@email.com', city: '上海市浦东新区', github: 'github.com/zhangming' },
    intention: { position: '后端工程师', city: '上海', salary: '25K–35K', arrival: '1个月内', industry: '互联网 / AI' },
    education: [{ school: '上海交通大学', major: '计算机科学与技术', degree: '本科', start: '2016.09', end: '2020.06', courses: 'GPA 3.7/4.0 · 数据结构、操作系统、数据库系统', gpa: '3.7/4.0' }],
    work: [
      { company: '上海云舟科技有限公司', position: '后端工程师', start: '2021.07', end: '至今', description: '负责公司核心业务系统设计与开发，支撑日均 1000 万+ 请求量\n拆分单体服务为 12 个微服务，系统可用性提升至 99.95%\n优化数据库慢查询与缓存策略，接口平均响应时间从 120ms 降至 45ms\n推动 CI/CD 落地，发布效率提升 70%', leave_reason: '' },
      { company: '上海智汇互联有限公司', position: '后端开发工程师', start: '2020.07', end: '2021.06', description: '参与电商平台后端开发，负责订单、支付与库存核心模块\n基于 Spring Cloud Alibaba 完成服务治理与配置管理', leave_reason: '' },
    ],
    project: [
      { name: '云舟微服务平台重构项目', role: '后端负责人', start: '2022.03', end: '2022.11', description: '设计并实现基于 Spring Cloud 的微服务框架，集成 Gateway、Sentinel、SkyWalking\n项目上线后系统稳定性提升 80%，运维成本降低 30%', tech_stack: ['Spring Cloud', 'MySQL', 'Redis'], url: '' },
      { name: '实时数据分析平台', role: '核心开发', start: '2021.09', end: '2022.02', description: '使用 Kafka + Flink 实现实时流处理，数据延迟从分钟级降至毫秒级\n基于 ClickHouse 建设查询链路，查询性能提升 5 倍', tech_stack: ['Kafka', 'Flink', 'ClickHouse'], url: '' },
    ],
    skills: [
      { category: '编程语言', name: 'Java（熟练）、Go（熟练）、Python（熟悉）', proficiency: '' },
      { category: '框架', name: 'Spring Boot、Spring Cloud Alibaba、Gin、MyBatis', proficiency: '' },
      { category: '数据库', name: 'MySQL、Redis、MongoDB、ClickHouse', proficiency: '' },
      { category: '工程工具', name: 'Docker、Kubernetes、Prometheus、Grafana', proficiency: '' },
    ],
  })
}

function normalizeContent(value: any): ResumeContent {
  return {
    personal: { name: '', gender: '', age: '', phone: '', email: '', github: '', avatar: '', city: '', ...(value.personal || {}) },
    intention: { position: '', city: '', salary: '', arrival: '', industry: '', ...(value.intention || {}) },
    education: Array.isArray(value.education) ? value.education : [],
    work: Array.isArray(value.work) ? value.work : [],
    project: Array.isArray(value.project) ? value.project.map((item) => ({ ...item, tech_stack: Array.isArray(item.tech_stack) ? item.tech_stack : [] })) : [],
    skills: Array.isArray(value.skills) ? value.skills : [],
    honor: Array.isArray(value.honor) ? value.honor : [], custom: Array.isArray(value.custom) ? value.custom : [],
    module_order: value.module_order || ['personal', 'intention', 'education', 'work', 'project', 'skills', 'honor'],
    module_visibility: value.module_visibility || { personal: true, intention: true, education: true, work: true, project: true, skills: true, honor: true },
  }
}

const toggleSection = (key: string) => { collapsed[key] = !collapsed[key] }
const isVisible = (key: string) => resumeContent.module_visibility[key] !== false
const toLines = (text: string) => text.split(/\n+/).map((line) => line.replace(/^[•·\-]\s*/, '').trim()).filter(Boolean)
const addEducation = () => resumeContent.education.push({ school: '', major: '', degree: '', start: '', end: '', courses: '', gpa: '' })
const addWork = () => resumeContent.work.push({ company: '', position: '', start: '', end: '', description: '', leave_reason: '' })
const addProject = () => resumeContent.project.push({ name: '', role: '', start: '', end: '', description: '', tech_stack: [], url: '' })
const addSkill = () => resumeContent.skills.push({ category: '', name: '', proficiency: '' })
const updateTechStack = (item: ResumeProject, event: Event) => { item.tech_stack = (event.target as HTMLInputElement).value.split(/[,，]/).map((value) => value.trim()).filter(Boolean) }
const selectVersion = (version: ResumeVersion) => { resumeStore.setCurrentVersion(version); versionDrawerOpen.value = false }

const handleNameSave = async () => {
  if (demoMode.value || !resumeStore.currentResume || !editableName.value.trim()) return
  await resumeStore.update(resumeId.value, { name: editableName.value.trim() })
}
const handleSaveVersion = async () => {
  if (!newVersionNote.value.trim()) return message.warning('请输入版本备注')
  if (demoMode.value) { showSaveVersionModal.value = false; newVersionNote.value = ''; return message.success('本地预览：版本已模拟保存') }
  savingVersion.value = true
  const version = await resumeStore.createVersion(resumeId.value, { content: JSON.stringify(resumeContent), change_note: newVersionNote.value.trim() })
  savingVersion.value = false
  if (version) { showSaveVersionModal.value = false; newVersionNote.value = '' }
}
const handleOptimize = () => message.success(demoMode.value ? '已完成本地诊断演示；接入后端后可生成真实优化版本' : '请使用现有 AI 操作生成优化版本')
const exportResume = () => window.print()
const backToList = () => demoMode.value ? router.push('/') : router.push('/app/resumes')
const formatVersionDate = (value: string) => value ? new Date(value).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : ''

onMounted(async () => {
  if (demoMode.value) return
  if (!resumeId.value || Number.isNaN(resumeId.value)) return router.push('/app/resumes')
  await resumeStore.fetchOne(resumeId.value)
  editableName.value = resumeStore.currentResume?.name || '未命名简历'
  targetJd.value = resumeStore.currentResume?.target_jd || targetJd.value
  await resumeStore.fetchVersions(resumeId.value)
  if (resumeStore.versions.length) resumeStore.setCurrentVersion(resumeStore.versions[0])
})
</script>

<style scoped>
.resume-lab { --ink: #17211d; --muted: #68736d; --line: #dfe4e0; --green: #087443; --green-soft: #e8f3ed; height: 100vh; min-height: 760px; overflow: hidden; background: #eef0ed; color: var(--ink); }
.lab-toolbar { height: 64px; display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 0 24px; background: #fff; border-bottom: 1px solid var(--line); position: relative; z-index: 5; }
.toolbar-leading, .toolbar-controls { display: flex; align-items: center; gap: 10px; min-width: 0; }
.toolbar-controls { overflow-x: auto; scrollbar-width: none; }
.icon-button { width: 36px; height: 36px; border: 0; background: transparent; border-radius: 7px; cursor: pointer; font-size: 16px; }
.icon-button:hover { background: #f1f3f1; }
.brand-mark { width: 34px; height: 34px; border-radius: 50%; display: grid; place-items: center; color: #fff; background: var(--green); font-weight: 800; }
.resume-title-wrap { display: flex; align-items: center; gap: 12px; min-width: 0; }
.resume-name { width: 190px; border: 0; border-bottom: 1px solid transparent; outline: 0; color: var(--ink); font-size: 16px; font-weight: 750; background: transparent; }
.resume-name:focus { border-color: var(--green); }
.save-state { color: var(--muted); font-size: 12px; white-space: nowrap; }
.save-state :deep(.anticon) { color: var(--green); margin-right: 4px; }
.toolbar-select { height: 34px; display: flex; align-items: center; gap: 6px; padding: 0 9px; border-left: 1px solid var(--line); color: #33403a; }
.toolbar-select select { border: 0; outline: 0; background: transparent; font-size: 13px; cursor: pointer; }
.toolbar-select.compact { padding-right: 4px; }
.zoom-control { height: 34px; display: flex; align-items: center; border: 1px solid var(--line); border-radius: 6px; overflow: hidden; background: #fafbfa; }
.zoom-control button { width: 30px; height: 100%; border: 0; background: transparent; cursor: pointer; }
.zoom-control span { width: 52px; text-align: center; font-size: 12px; }
.export-button { background: var(--green); box-shadow: none; }
.lab-workspace { height: calc(100vh - 64px); display: grid; grid-template-columns: 340px minmax(640px, 1fr) 320px; }
.editor-panel, .ai-panel { background: #fff; min-width: 0; overflow: hidden; }
.editor-panel { border-right: 1px solid var(--line); }
.ai-panel { border-left: 1px solid var(--line); padding: 20px; overflow-y: auto; }
.panel-tabs { height: 51px; display: grid; grid-template-columns: 1fr 1fr; border-bottom: 1px solid var(--line); }
.panel-tabs button { border: 0; background: #fff; color: var(--muted); font-weight: 650; cursor: pointer; position: relative; }
.panel-tabs button.active { color: var(--green); }
.panel-tabs button.active::after { content: ''; height: 2px; background: var(--green); position: absolute; left: 22px; right: 22px; bottom: 0; }
.panel-hint { padding: 12px 16px; color: var(--muted); font-size: 12px; background: #fafbfa; border-bottom: 1px solid #edf0ed; }
.editor-scroll, .raw-json-panel { height: calc(100% - 51px); overflow-y: auto; }
.raw-json-panel pre { margin: 0; padding: 18px; color: #34423b; font: 12px/1.7 ui-monospace, SFMono-Regular, Menlo, monospace; white-space: pre-wrap; }
.form-section { border-bottom: 1px solid var(--line); }
.section-heading { width: 100%; height: 48px; padding: 0 16px; display: flex; align-items: center; justify-content: space-between; border: 0; background: #fff; color: var(--ink); cursor: pointer; font-weight: 720; }
.section-heading span { display: flex; align-items: center; gap: 9px; }
.section-heading > :deep(.anticon) { transition: transform .2s; color: #89918d; }
.form-section.collapsed .section-heading > :deep(.anticon) { transform: rotate(-90deg); }
.form-section.collapsed .section-body { display: none; }
.section-body { padding: 0 16px 16px; }
.grid-two { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.field { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.field.full { grid-column: 1 / -1; }
.field span { color: var(--muted); font-size: 11px; font-weight: 650; }
.field input, .field textarea, .skill-row input, .jd-field textarea { width: 100%; border: 1px solid #d9dedb; border-radius: 6px; outline: none; padding: 8px 9px; color: var(--ink); background: #fff; font: inherit; font-size: 12px; transition: border-color .2s, box-shadow .2s; }
.field textarea { resize: vertical; line-height: 1.55; }
.field input:focus, .field textarea:focus, .skill-row input:focus, .jd-field textarea:focus { border-color: var(--green); box-shadow: 0 0 0 3px rgba(8,116,67,.08); }
.repeat-list { display: flex; flex-direction: column; gap: 10px; }
.repeat-item { position: relative; padding: 12px; border: 1px solid #e1e5e2; border-radius: 7px; background: #fbfcfb; }
.remove-item { position: absolute; top: 5px; right: 5px; width: 26px; height: 26px; border: 0; border-radius: 5px; background: transparent; color: #929a96; cursor: pointer; z-index: 2; }
.remove-item:hover { color: #b42318; background: #fff0ee; }
.add-item { min-height: 34px; border: 1px dashed #bdc8c1; border-radius: 6px; background: #fff; color: var(--green); cursor: pointer; font-weight: 650; }
.skill-row { display: grid; grid-template-columns: 90px 1fr 26px; gap: 6px; }
.skill-row button { border: 0; background: transparent; color: #999; cursor: pointer; }
.preview-stage { min-width: 0; overflow: hidden; background: #e9ebe8; display: flex; flex-direction: column; }
.stage-meta { height: 42px; padding: 0 18px; display: flex; align-items: center; justify-content: space-between; color: #69736e; font-size: 11px; border-bottom: 1px solid #d8dcd9; }
.paper-viewport { flex: 1; overflow: auto; padding: 26px 32px 80px; display: flex; justify-content: center; align-items: flex-start; }
.resume-paper { width: 794px; min-height: 1123px; transform-origin: top center; padding: 56px 60px; background: #fff; box-shadow: 0 8px 28px rgba(31,40,35,.13); color: #161d19; transition: transform .18s ease; }
.resume-paper.font-serif { font-family: 'Songti SC', 'Noto Serif SC', serif; }
.resume-paper.font-sans { font-family: 'PingFang SC', 'Noto Sans SC', sans-serif; }
.paper-header { padding-bottom: 16px; border-bottom: 2px solid #26312b; }
.paper-header h1 { margin: 0; font-size: 31px; line-height: 1.15; letter-spacing: .08em; font-weight: 800; }
.paper-role { margin: 8px 0 10px; color: var(--green); font-size: 14px; font-weight: 750; }
.contact-line { display: flex; flex-wrap: wrap; gap: 5px 18px; color: #4d5651; font-size: 10.5px; }
.contact-line span:not(:last-child)::after { content: '·'; margin-left: 18px; color: #9ca39f; }
.paper-section { margin-top: 20px; }
.paper-section h2 { margin: 0 0 9px; padding-bottom: 5px; border-bottom: 1px solid #39433e; font-size: 15px; line-height: 1; letter-spacing: .08em; }
.paper-section p { margin: 0; font-size: 10.8px; line-height: 1.65; }
.resume-entry + .resume-entry { margin-top: 13px; }
.entry-heading, .entry-subheading { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; }
.entry-heading strong { font-size: 11.5px; }
.entry-heading span, .entry-subheading { color: #545e58; font-size: 10px; }
.entry-subheading { margin-top: 3px; justify-content: flex-start; }
.entry-subheading .tech-stack { margin-left: auto; color: #68736d; }
.entry-note { color: #505a54; }
.resume-entry ul { margin: 6px 0 0; padding-left: 17px; }
.resume-entry li { margin: 2px 0; font-size: 10.4px; line-height: 1.55; }
.skill-section p + p { margin-top: 3px; }
.template-modern .paper-header { padding: 18px 20px; border: 0; background: #edf5f0; }
.template-modern .paper-section h2 { color: var(--green); border-color: #8fb9a3; }
.density-compact { padding-top: 46px; padding-bottom: 46px; }
.density-compact .paper-section { margin-top: 15px; }
.density-compact .resume-entry + .resume-entry { margin-top: 9px; }
.ai-heading { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.ai-heading h2 { margin: 3px 0 0; font-size: 15px; }
.ai-heading > :deep(.anticon) { color: var(--green); font-size: 20px; }
.ai-kicker { color: var(--green); font-size: 12px; font-weight: 800; letter-spacing: .08em; }
.match-card, .dimension-card, .suggestion-card { padding: 15px; border: 1px solid var(--line); border-radius: 8px; margin-bottom: 12px; }
.match-card > span { color: var(--muted); font-size: 12px; }
.score-line { display: flex; align-items: center; gap: 10px; margin: 5px 0 9px; }
.score-line strong { color: var(--green); font-size: 36px; line-height: 1; }
.score-line em { padding: 3px 8px; border-radius: 4px; background: var(--green-soft); color: var(--green); font-style: normal; font-size: 11px; }
.score-track { height: 6px; border-radius: 10px; background: #e7eae8; overflow: hidden; }
.score-track i, .dimension-row div i { display: block; height: 100%; background: var(--green); border-radius: inherit; }
.match-card small { display: block; margin-top: 9px; color: #8a938e; font-size: 10px; }
.dimension-card h3, .suggestion-card h3 { margin: 0 0 12px; font-size: 13px; }
.dimension-row { display: grid; grid-template-columns: 62px 1fr 26px; align-items: center; gap: 8px; margin: 8px 0; color: #5c6660; font-size: 10px; }
.dimension-row div { height: 4px; border-radius: 8px; background: #edf0ee; overflow: hidden; }
.dimension-row b { color: #65706a; text-align: right; }
.suggestion-card h3 span { color: var(--green); }
.suggestion-card article { display: flex; gap: 9px; padding: 10px 0; border-top: 1px solid #edf0ed; }
.suggestion-card article > b { flex: 0 0 22px; height: 22px; display: grid; place-items: center; border-radius: 5px; color: var(--green); background: var(--green-soft); font-size: 11px; }
.suggestion-card strong { font-size: 11.5px; }
.suggestion-card p { margin: 4px 0 0; color: var(--muted); font-size: 10px; line-height: 1.5; }
.jd-field { display: flex; flex-direction: column; gap: 7px; margin-top: 14px; }
.jd-field span { font-size: 12px; font-weight: 700; }
.jd-field textarea { resize: vertical; font-size: 11px; line-height: 1.5; }
.optimize-button { width: 100%; min-height: 42px; margin-top: 10px; border: 0; border-radius: 7px; background: var(--green); color: #fff; cursor: pointer; font-weight: 750; }
.privacy-note { margin: 10px 0 0; color: #89918d; font-size: 9.5px; text-align: center; }
.version-list { display: flex; flex-direction: column; gap: 8px; }
.version-list button { display: flex; flex-direction: column; gap: 5px; padding: 12px; border: 1px solid var(--line); border-radius: 7px; background: #fff; text-align: left; cursor: pointer; }
.version-list button.active { border-color: var(--green); background: var(--green-soft); }
.version-list button span { display: flex; justify-content: space-between; }
.version-list small, .version-list em { color: var(--muted); font-size: 11px; font-style: normal; }
@media (max-width: 1180px) { .lab-workspace { grid-template-columns: 320px minmax(600px, 1fr); } .ai-panel { display: none; } .toolbar-select { display: none; } }
@media (max-width: 820px) { .resume-lab { height: auto; min-height: 100vh; overflow: visible; } .lab-toolbar { height: auto; min-height: 64px; padding: 10px 14px; flex-wrap: wrap; } .toolbar-controls { width: 100%; padding-bottom: 2px; } .save-state { display: none; } .lab-workspace { height: auto; display: flex; flex-direction: column; } .editor-panel { height: 58vh; border: 0; border-bottom: 1px solid var(--line); } .preview-stage { height: 72vh; } .resume-paper { flex: 0 0 794px; } }
@media print { .lab-toolbar, .editor-panel, .ai-panel, .stage-meta { display: none !important; } .resume-lab, .lab-workspace, .preview-stage, .paper-viewport { display: block; height: auto; overflow: visible; background: #fff; padding: 0; } .resume-paper { transform: none !important; box-shadow: none; margin: 0; width: 210mm; min-height: 297mm; } }
</style>
