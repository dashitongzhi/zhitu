export type ResumeTemplateId =
  | 'classic'
  | 'modern'
  | 'executive'
  | 'compact'
  | 'sidebar'
  | 'editorial'
  | 'minimal'
  | 'academic'
  | 'creative'
  | 'graduate'

export interface ResumeTemplate {
  id: ResumeTemplateId
  name: string
  description: string
  audience: string
  tone: string
  accent: string
  preview: 'single' | 'band' | 'formal' | 'dense' | 'split' | 'magazine' | 'air' | 'paper' | 'block' | 'fresh'
}

export const resumeTemplates: ResumeTemplate[] = [
  { id: 'classic', name: '经典单栏', description: '清晰稳妥的标准结构，兼顾阅读效率与 ATS 识别。', audience: '通用岗位', tone: '稳健', accent: '#087443', preview: 'single' },
  { id: 'modern', name: '现代强调', description: '用浅色标题区强化个人定位，适合互联网与产品岗位。', audience: '产品 / 运营', tone: '现代', accent: '#16705a', preview: 'band' },
  { id: 'executive', name: '管理层', description: '克制的深蓝与衬线标题，突出领导力和业务成果。', audience: '管理 / 咨询', tone: '专业', accent: '#25354a', preview: 'formal' },
  { id: 'compact', name: '技术密集', description: '高信息密度与紧凑节奏，让项目和技术栈占据主角。', audience: '研发 / 数据', tone: '高效', accent: '#263f36', preview: 'dense' },
  { id: 'sidebar', name: '双栏侧边', description: '侧栏承载联系方式和技能，主栏集中展示核心经历。', audience: '设计 / 市场', tone: '有序', accent: '#214c3d', preview: 'split' },
  { id: 'editorial', name: '编辑部', description: '杂志式大标题与暖色标记，强调表达与个人观点。', audience: '内容 / 品牌', tone: '编辑感', accent: '#b5523b', preview: 'magazine' },
  { id: 'minimal', name: '极简留白', description: '减少装饰与分隔，以留白和文字层级建立高级感。', audience: '全行业', tone: '克制', accent: '#252b28', preview: 'air' },
  { id: 'academic', name: '学术研究', description: '论文式排版和明确的时间线，适合研究与教育经历。', audience: '科研 / 教育', tone: '严谨', accent: '#403b35', preview: 'paper' },
  { id: 'creative', name: '创意品牌', description: '鲜明色块与大胆姓名区，适合作品导向的求职者。', audience: '创意 / 设计', tone: '鲜明', accent: '#126b51', preview: 'block' },
  { id: 'graduate', name: '校园新锐', description: '轻快蓝绿层级，把教育、实践与校园项目放在前面。', audience: '应届 / 实习', tone: '清新', accent: '#277c87', preview: 'fresh' },
]

export const getResumeTemplate = (id: string | null | undefined) =>
  resumeTemplates.find((template) => template.id === id) || resumeTemplates[0]

export const createTemplateContent = (templateId: ResumeTemplateId) => JSON.stringify({
  template_style: templateId,
  personal: {
    name: '你的姓名',
    gender: '',
    age: '',
    phone: '',
    email: '',
    github: '',
    avatar: '',
    city: '',
  },
  intention: {
    position: '目标岗位',
    city: '',
    salary: '',
    arrival: '',
    industry: '',
  },
  education: [],
  work: [],
  project: [],
  skills: [],
  honor: [],
  custom: [],
  module_order: ['personal', 'intention', 'education', 'work', 'project', 'skills', 'honor'],
  module_visibility: {
    personal: true,
    intention: true,
    education: true,
    work: true,
    project: true,
    skills: true,
    honor: true,
  },
})
