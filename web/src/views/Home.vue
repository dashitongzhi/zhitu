<template>
  <div class="home-page" :style="{ '--page-bg': `url(${pageBackgroundImage})` }">
    <!-- 整页底层背景图：极淡写实摄影，fixed 固定视差 -->
    <div class="page-bg-layer" aria-hidden="true"></div>

    <!-- 顶部导航栏：固定定位 + 毛玻璃背景 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="brand">
          <span class="brand-logo">
            <ThunderboltOutlined />
          </span>
          <span class="brand-name">职途</span>
        </div>
        <nav class="nav-links">
          <a href="#features" @click.prevent="scrollTo('features')">功能特性</a>
          <a href="#workflow" @click.prevent="scrollTo('workflow')">使用流程</a>
          <a href="#faq" @click.prevent="scrollTo('faq')">常见问题</a>
        </nav>
        <div class="nav-actions">
          <button class="btn-ghost" type="button" @click="openAuthModal('login')">
            登录
          </button>
          <button class="btn-primary" type="button" @click="openAuthModal('register')">
            免费注册
          </button>
        </div>
      </div>
    </header>

    <!-- Hero 区：大标题 + 副标题 + 统计数据 + 双 CTA -->
    <section class="hero">
      <!-- 写实摄影背景层 -->
      <div class="hero-bg-layer" :style="{ backgroundImage: `url(${heroBackgroundImage})` }" aria-hidden="true"></div>
      <div class="hero-bg-overlay" aria-hidden="true"></div>
      <!-- 装饰性光晕背景 -->
      <div class="hero-blob hero-blob-1"></div>
      <div class="hero-blob hero-blob-2"></div>

      <div class="hero-container">
        <div class="hero-left">
          <!-- 顶部标签 -->
          <div class="hero-badge">
            <span class="hero-badge-dot"></span>
            <span>AI 驱动的智能求职助手</span>
          </div>

          <!-- 主标题：响应式缩放，48px / 700 -->
          <h1 class="hero-title">
            让每一次求职都<br />
            <span class="hero-title-accent">精准而高效</span>
          </h1>

          <!-- 副标题：18px / muted -->
          <p class="hero-desc">
            职途 融合智能简历生成、模拟面试与投递看板，为求职者提供端到端的求职支持，让每一次机会都被认真对待。
          </p>

          <!-- 数据统计 -->
          <div class="hero-stats">
            <div class="hero-stat">
              <div class="hero-stat-value">4+</div>
              <div class="hero-stat-label">核心模块</div>
            </div>
            <div class="hero-stat">
              <div class="hero-stat-value">AI</div>
              <div class="hero-stat-label">智能生成</div>
            </div>
            <div class="hero-stat">
              <div class="hero-stat-value">100%</div>
              <div class="hero-stat-label">本地化部署</div>
            </div>
          </div>

          <!-- 双 CTA 按钮：主按钮「立即开始」+ 次按钮「了解更多」 -->
          <div class="hero-cta">
            <button class="btn-primary btn-lg" type="button" @click="handlePrimaryCta">
              立即开始
              <ArrowRightOutlined />
            </button>
            <button class="btn-secondary btn-lg" type="button" @click="scrollTo('features')">
              了解更多
            </button>
          </div>
        </div>

        <!-- Hero 右侧：产品预览浮卡组合 -->
        <div class="hero-right">
          <!-- 主预览卡片：模拟简历编辑器 -->
          <div class="preview-card preview-main">
            <div class="preview-header">
              <span class="preview-dot" style="background:#ff5f57"></span>
              <span class="preview-dot" style="background:#febc2e"></span>
              <span class="preview-dot" style="background:#28c840"></span>
              <span class="preview-title">简历实验室</span>
            </div>
            <div class="preview-cover">
              <img :src="heroPreviewImage" alt="简历实验室工作场景" loading="lazy" />
              <div class="preview-cover-overlay">
                <div class="preview-cover-row">
                  <span class="preview-cover-label">JD 匹配度</span>
                  <span class="preview-cover-value">92%</span>
                </div>
                <div class="preview-cover-row">
                  <span class="preview-cover-label">AI 优化</span>
                  <span class="preview-cover-tags">
                    <span class="preview-tag">措辞润色</span>
                    <span class="preview-tag">经历强化</span>
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 浮卡 1：面试训练 -->
          <div class="preview-card preview-float preview-float-1">
            <div class="preview-float-icon" style="background:var(--brand-50);color:var(--primary)">
              <MessageOutlined />
            </div>
            <div>
              <div class="preview-float-title">面试训练中</div>
              <div class="preview-float-sub">第 3 轮 · 深度追问</div>
            </div>
          </div>

          <!-- 浮卡 2：投递看板统计 -->
          <div class="preview-card preview-float preview-float-2">
            <div class="preview-float-icon" style="background:rgba(52,199,89,0.14);color:#34c759">
              <SendOutlined />
            </div>
            <div>
              <div class="preview-float-title">本周投递</div>
              <div class="preview-float-sub">12 次 · 3 次面试</div>
            </div>
          </div>

          <!-- 浮卡 3：Offer 数 -->
          <div class="preview-card preview-float preview-float-3">
            <div class="preview-float-icon" style="background:rgba(175,82,222,0.14);color:#af52de">
              <CheckCircleFilled />
            </div>
            <div>
              <div class="preview-float-title">已获 Offer</div>
              <div class="preview-float-sub">2 个 · 待抉择</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 功能入口区：3 个卡片（简历实验室、面试训练场、投递看板） -->
    <section id="features" class="section">
      <div class="container">
        <div class="section-head">
          <span class="section-eyebrow">功能特性</span>
          <h2 class="section-title">一站式求职解决方案</h2>
          <p class="section-sub">
            从简历到面试，从投递到追踪，职途 覆盖求职全流程。
          </p>
        </div>

        <div class="feature-grid">
          <article
            v-for="feature in features"
            :key="feature.title"
            class="feature-card"
            @click="handleFeatureClick(feature)"
          >
            <!-- 卡片图标 -->
            <div class="feature-icon">
              <component :is="feature.icon" />
            </div>
            <!-- 卡片正文 -->
            <div class="feature-body">
              <h3 class="feature-title">{{ feature.title }}</h3>
              <p class="feature-desc">{{ feature.desc }}</p>
              <!-- 特性标签 -->
              <div class="feature-tags">
                <span v-for="tag in feature.tags" :key="tag" class="feature-tag">{{ tag }}</span>
              </div>
              <!-- 亮点列表 -->
              <ul class="feature-highlights">
                <li v-for="h in feature.highlights" :key="h">
                  <CheckCircleFilled class="feature-check" />
                  <span>{{ h }}</span>
                </li>
              </ul>
            </div>
            <!-- 卡片箭头：hover 时右移 -->
            <div class="feature-arrow">
              <ArrowRightOutlined />
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- 核心能力：6 项技术亮点（2x3 网格） -->
    <section class="section section-tint">
      <div class="container">
        <div class="section-head">
          <span class="section-eyebrow">核心能力</span>
          <h2 class="section-title">技术驱动的求职引擎</h2>
          <p class="section-sub">
            从 AI 匹配到数据看板，从安全隔离到极速响应，每一项能力都为求职效率服务。
          </p>
        </div>

        <div class="capability-grid">
          <article
            v-for="cap in capabilities"
            :key="cap.title"
            class="capability-card"
          >
            <div class="capability-icon">
              <component :is="cap.icon" />
            </div>
            <div class="capability-body">
              <h3 class="capability-title">{{ cap.title }}</h3>
              <p class="capability-desc">{{ cap.desc }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- 数据成果：渐变背景 + 4 个大数字 -->
    <section class="section section-grad">
      <div class="container">
        <div class="section-head section-head-light">
          <span class="section-eyebrow section-eyebrow-light">数据成果</span>
          <h2 class="section-title section-title-light">用数字说话</h2>
          <p class="section-sub section-sub-light">
            一组关键指标，看清职途 的能力边界。
          </p>
        </div>

        <div class="achievement-grid">
          <div v-for="a in achievements" :key="a.label" class="achievement-card">
            <div class="achievement-value">{{ a.value }}</div>
            <div class="achievement-label">{{ a.label }}</div>
            <div class="achievement-sub">{{ a.sub }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 用户证言：3 张卡片墙 -->
    <section class="section section-photo">
      <div class="section-bg-layer" :style="{ backgroundImage: `url(${testimonialBackgroundImage})` }" aria-hidden="true"></div>
      <div class="container">
        <div class="section-head">
          <span class="section-eyebrow">用户证言</span>
          <h2 class="section-title">求职者的真实反馈</h2>
          <p class="section-sub">
            来自不同岗位、不同经验的求职者，分享他们与职途 的故事。
          </p>
        </div>

        <div class="testimonial-grid">
          <article
            v-for="t in testimonials"
            :key="t.name"
            class="testimonial-card"
          >
            <div class="testimonial-stars">
              <StarFilled v-for="n in 5" :key="n" />
            </div>
            <p class="testimonial-quote">"{{ t.quote }}"</p>
            <div class="testimonial-author">
              <img :src="t.photo" :alt="t.name" class="testimonial-avatar" loading="lazy" />
              <div>
                <div class="testimonial-name">{{ t.name }}</div>
                <div class="testimonial-role">{{ t.role }}</div>
              </div>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- 使用流程：4 步时间轴 -->
    <section id="workflow" class="section section-alt">
      <div class="section-bg-layer" :style="{ backgroundImage: `url(${workflowBackgroundImage})` }" aria-hidden="true"></div>
      <div class="container">
        <div class="section-head">
          <span class="section-eyebrow">使用流程</span>
          <h2 class="section-title">四步开启你的求职之旅</h2>
          <p class="section-sub">
            从注册到入职，职途 陪你走完求职全流程。
          </p>
        </div>

        <div class="workflow-grid">
          <article
            v-for="(step, idx) in workflows"
            :key="step.title"
            class="workflow-card"
          >
            <!-- 步骤序号 -->
            <div class="workflow-index">
              <span class="workflow-index-num">{{ idx + 1 }}</span>
              <span
                v-if="idx < workflows.length - 1"
                class="workflow-index-line"
              ></span>
            </div>
            <!-- 步骤内容 -->
            <div class="workflow-content">
              <div class="workflow-icon">
                <component :is="step.icon" />
              </div>
              <h3 class="workflow-title">{{ step.title }}</h3>
              <p class="workflow-desc">{{ step.desc }}</p>
              <ul class="workflow-points">
                <li v-for="point in step.points" :key="point">{{ point }}</li>
              </ul>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- 常见问题：折叠面板 -->
    <section id="faq" class="section">
      <div class="container container-narrow">
        <div class="section-head">
          <span class="section-eyebrow">常见问题</span>
          <h2 class="section-title">还有疑问？看这里</h2>
          <p class="section-sub">
            收录用户最常问的问题，帮你快速上手职途。
          </p>
        </div>

        <a-collapse
          :bordered="false"
          :accordion="true"
          class="faq-collapse"
          :default-active-key="['1']"
        >
          <a-collapse-panel
            v-for="faq in faqs"
            :key="faq.key"
            :header="faq.q"
            class="faq-panel"
          >
            <p class="faq-answer">{{ faq.a }}</p>
          </a-collapse-panel>
        </a-collapse>

        <!-- 底部辅助 CTA -->
        <div class="faq-cta">
          <span class="faq-cta-text">没有找到你要的答案？</span>
          <button
            class="btn-ghost"
            type="button"
            @click="openAuthModal('register')"
          >
            立即体验，遇到问题随时反馈
          </button>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-grid">
          <div class="footer-col footer-col-brand">
            <div class="footer-brand">
              <span class="footer-logo">
                <ThunderboltOutlined />
              </span>
              <span>职途</span>
            </div>
            <p class="footer-tagline">让求职更简单 · 智能求职助手平台</p>
            <p class="footer-copy">© 2026 职途. All rights reserved.</p>
          </div>
          <div class="footer-col">
            <h4 class="footer-title">产品</h4>
            <a @click.prevent="handleFeatureClick({ path: '/app/resumes' })">简历实验室</a>
            <a @click.prevent="handleFeatureClick({ path: '/app/interviews' })">面试训练场</a>
            <a @click.prevent="handleFeatureClick({ path: '/app/applications' })">投递看板</a>
          </div>
          <div class="footer-col">
            <h4 class="footer-title">资源</h4>
            <a @click.prevent="scrollTo('features')">功能特性</a>
            <a @click.prevent="scrollTo('workflow')">使用流程</a>
            <a @click.prevent="scrollTo('faq')">常见问题</a>
          </div>
          <div class="footer-col">
            <h4 class="footer-title">关于</h4>
            <a @click.prevent="openAuthModal('register')">免费注册</a>
            <a @click.prevent="openAuthModal('login')">登录</a>
            <a @click.prevent="handlePrimaryCta">立即开始</a>
          </div>
        </div>
      </div>
    </footer>

    <!-- 登录 / 注册弹窗：使用 a-modal + a-form -->
    <a-modal
      v-model:open="authModalVisible"
      :footer="null"
      :width="420"
      class="auth-modal"
      :destroy-on-close="false"
    >
      <div class="auth-modal-body">
        <!-- 弹窗标题 -->
        <h3 class="auth-modal-title">
          {{ activeTab === 'login' ? '欢迎回来' : '创建账号' }}
        </h3>
        <p class="auth-modal-sub">
          {{ activeTab === 'login' ? '登录以继续你的求职之旅' : '注册即可免费使用全部功能' }}
        </p>

        <!-- 登录 / 注册 切换 tabs -->
        <div class="auth-tabs">
          <button
            class="auth-tab"
            :class="{ 'auth-tab-active': activeTab === 'login' }"
            type="button"
            @click="activeTab = 'login'"
          >
            登录
          </button>
          <button
            class="auth-tab"
            :class="{ 'auth-tab-active': activeTab === 'register' }"
            type="button"
            @click="activeTab = 'register'"
          >
            注册
          </button>
        </div>

        <!-- 登录表单 -->
        <a-form
          v-if="activeTab === 'login'"
          :model="loginForm"
          :rules="loginRules"
          layout="vertical"
          class="auth-form"
          @finish="handleLogin"
        >
          <a-form-item label="邮箱" name="email">
            <a-input
              v-model:value="loginForm.email"
              placeholder="请输入邮箱"
              size="large"
            >
              <template #prefix>
                <MailOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="密码" name="password">
            <a-input-password
              v-model:value="loginForm.password"
              placeholder="请输入密码"
              size="large"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <button
              type="submit"
              class="btn-primary btn-lg auth-submit"
              :class="{ 'is-loading': loading }"
              :disabled="loading"
            >
              {{ loading ? '登录中…' : '登录' }}
            </button>
          </a-form-item>

          <div class="auth-switch">
            还没有账号？
            <a @click="activeTab = 'register'">立即注册</a>
          </div>
        </a-form>

        <!-- 注册表单 -->
        <a-form
          v-else
          :model="registerForm"
          :rules="registerRules"
          layout="vertical"
          class="auth-form"
          @finish="handleRegister"
        >
          <a-form-item label="邮箱" name="email">
            <a-input
              v-model:value="registerForm.email"
              placeholder="请输入邮箱"
              size="large"
            >
              <template #prefix>
                <MailOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="昵称" name="nickname">
            <a-input
              v-model:value="registerForm.nickname"
              placeholder="请输入昵称"
              size="large"
            >
              <template #prefix>
                <UserOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="密码" name="password">
            <a-input-password
              v-model:value="registerForm.password"
              placeholder="请输入密码（不少于 6 位）"
              size="large"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item label="确认密码" name="confirmPassword">
            <a-input-password
              v-model:value="registerForm.confirmPassword"
              placeholder="请再次输入密码"
              size="large"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <button
              type="submit"
              class="btn-primary btn-lg auth-submit"
              :class="{ 'is-loading': loading }"
              :disabled="loading"
            >
              {{ loading ? '注册中…' : '注册' }}
            </button>
          </a-form-item>

          <div class="auth-switch">
            已有账号？
            <a @click="activeTab = 'login'">立即登录</a>
          </div>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  MailOutlined,
  LockOutlined,
  UserOutlined,
  ArrowRightOutlined,
  FileTextOutlined,
  MessageOutlined,
  SendOutlined,
  ThunderboltOutlined,
  UserAddOutlined,
  EditOutlined,
  VideoCameraOutlined,
  FundProjectionScreenOutlined,
  SafetyCertificateOutlined,
  RocketOutlined,
  BulbOutlined,
  LineChartOutlined,
  CloudOutlined,
  ClockCircleOutlined,
  CheckCircleFilled,
  StarFilled,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const authStore = useAuthStore()

// 弹窗状态
const authModalVisible = ref(false)
const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)

// 登录表单
const loginForm = reactive({
  email: '',
  password: '',
})

const loginRules: Record<string, Rule[]> = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
}

// 注册表单
const registerForm = reactive({
  email: '',
  nickname: '',
  password: '',
  confirmPassword: '',
})

// 确认密码校验
const validateConfirmPassword = async (_rule: Rule, value: string) => {
  if (value !== registerForm.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const registerRules: Record<string, Rule[]> = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度为 2-20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

// 打开认证弹窗
const openAuthModal = (tab: 'login' | 'register') => {
  activeTab.value = tab
  authModalVisible.value = true
}

// 处理登录
const handleLogin = async () => {
  loading.value = true
  try {
    const success = await authStore.login(loginForm)
    if (success) {
      authModalVisible.value = false
      router.push('/app')
    }
  } finally {
    loading.value = false
  }
}

// 处理注册
const handleRegister = async () => {
  loading.value = true
  try {
    const success = await authStore.register({
      email: registerForm.email,
      nickname: registerForm.nickname,
      password: registerForm.password,
    })
    if (success) {
      // 注册成功后切换到登录 tab，并预填邮箱
      loginForm.email = registerForm.email
      loginForm.password = ''
      activeTab.value = 'login'
      message.info('请使用注册的账号登录')
    }
  } finally {
    loading.value = false
  }
}

// 平滑滚动到指定锚点
const scrollTo = (id: string) => {
  const el = document.getElementById(id)
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

// 主 CTA 按钮：已登录跳转 /app，未登录打开注册弹窗
const handlePrimaryCta = () => {
  if (authStore.isAuthenticated) {
    router.push('/app')
  } else {
    openAuthModal('register')
  }
}

// 功能入口卡片数据：3 个核心模块
const features = [
  {
    title: '简历实验室',
    desc: '基于目标岗位 JD 智能生成简历，AI 润色表达，多份简历灵活切换。',
    icon: markRaw(FileTextOutlined),
    path: '/app/resumes',
    tags: ['JD 匹配', 'AI 润色', '多份管理'],
    highlights: ['智能生成结构化简历', '一键优化措辞表达', 'A4 实时预览'],
  },
  {
    title: '面试训练场',
    desc: '多轮深度追问，语音与文字双模式，面试记录随时回看复盘。',
    icon: markRaw(MessageOutlined),
    path: '/app/interviews',
    tags: ['多轮追问', '语音文字', '复盘回放'],
    highlights: ['覆盖主流岗位题库', 'AI 实时深度追问', '完整记录随时回看'],
  },
  {
    title: '投递看板',
    desc: '看板式管理投递状态，记录面试轮次，洞察投递漏斗转化。',
    icon: markRaw(SendOutlined),
    path: '/app/applications',
    tags: ['看板视图', '漏斗分析', '面试追踪'],
    highlights: ['6 列状态灵活拖拽', '面试轮次完整记录', '投递漏斗转化洞察'],
  },
]

// 核心能力：6 项技术亮点
const capabilities = [
  {
    title: 'JD 智能匹配',
    desc: '基于目标岗位 JD 自动生成结构化简历内容，匹配度提升 80%。',
    icon: markRaw(BulbOutlined),
  },
  {
    title: '多轮深度追问',
    desc: 'AI 模拟真实面试官，根据你的回答进行多轮深度追问。',
    icon: markRaw(MessageOutlined),
  },
  {
    title: '实时数据看板',
    desc: '投递漏斗、转化率、行业分布一目了然，数据驱动求职决策。',
    icon: markRaw(LineChartOutlined),
  },
  {
    title: '本地化部署',
    desc: '所有数据存储在你自己的服务器，隐私不外泄，安全可控。',
    icon: markRaw(CloudOutlined),
  },
  {
    title: 'JWT 双密钥隔离',
    desc: '用户与管理员 token 独立签发，权限边界清晰，安全保障。',
    icon: markRaw(SafetyCertificateOutlined),
  },
  {
    title: '极速响应',
    desc: '流式 AI 输出，秒级响应，简历生成与面试对答流畅无卡顿。',
    icon: markRaw(RocketOutlined),
  },
]

// 数据成果
const achievements = [
  { value: '4+', label: '核心功能模块', sub: '简历 / 面试 / 投递 / 数据' },
  { value: '6', label: '看板状态列', sub: '从投递到 Offer 全链路' },
  { value: '100%', label: '本地化部署', sub: '数据自托管，隐私可控' },
  { value: '∞', label: 'AI 对话轮次', sub: '多轮深度追问不限制' },
]

// 用户证言
const testimonials = [
  {
    quote: '从粘贴 JD 到生成简历只用了 30 秒，AI 优化的措辞比我手写的专业太多了。',
    name: '张同学',
    role: '前端工程师 · 3 年经验',
    avatar: 'Z',
    color: '#007aff',
    photo:
      'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20headshot%20portrait%20of%20young%20Chinese%20woman%2C%20mid%20twenties%2C%20software%20engineer%2C%20soft%20natural%20light%2C%20blurred%20office%20background%2C%20subtle%20smile%2C%20smart%20casual%20blouse%2C%204k&image_size=square_hd',
  },
  {
    quote: '模拟面试的追问很真实，第一次面试被问懵的题，第二次在这里练熟了，正式面试拿下 Offer。',
    name: '李同学',
    role: '后端开发 · 应届生',
    avatar: 'L',
    color: '#34c759',
    photo:
      'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20headshot%20portrait%20of%20young%20Chinese%20man%2C%20early%20twenties%2C%20graduate%20developer%2C%20soft%20studio%20light%2C%20neutral%20background%2C%20friendly%20smile%2C%20collared%20shirt%2C%204k&image_size=square_hd',
  },
  {
    quote: '投递看板让我看清自己的漏斗，原来拒了 20 家只有 3 家面试，调整策略后面试率翻倍。',
    name: '王同学',
    role: '产品经理 · 转行求职',
    avatar: 'W',
    color: '#af52de',
    photo:
      'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20headshot%20portrait%20of%20Chinese%20man%20around%20thirty%2C%20product%20manager%2C%20soft%20natural%20light%2C%20blurred%20office%20background%2C%20confident%20expression%2C%20business%20casual%2C%204k&image_size=square_hd',
  },
]

// Hero 右侧主预览图：写实摄影风的简历桌面场景
const heroPreviewImage =
  'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20editorial%20workspace%20photography%2C%20modern%20wooden%20desk%20with%20open%20laptop%20showing%20resume%20document%2C%20notebook%2C%20pen%2C%20soft%20natural%20window%20light%20from%20left%2C%20shallow%20depth%20of%20field%2C%20warm%20neutral%20color%20palette%2C%20professional%204k%20photograph&image_size=landscape_4_3'

// 整页底层背景：极淡的办公室写实摄影，做整体氛围底色
const pageBackgroundImage =
  'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20minimalist%20modern%20office%20interior%2C%20soft%20natural%20daylight%2C%20light%20wooden%20floor%2C%20white%20walls%2C%20subtle%20plant%20silhouette%2C%20extremely%20subtle%20and%20calm%2C%20muted%20neutral%20beige%20tones%2C%20professional%20corporate%20atmosphere%2C%204k%20photograph&image_size=landscape_16_9'

// Hero 区背景：城市写字楼天际线写实摄影
const heroBackgroundImage =
  'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20modern%20city%20skyline%20with%20glass%20office%20towers%20at%20blue%20hour%2C%20soft%20gradient%20sky%2C%20calm%20professional%20corporate%20mood%2C%20subtle%20haze%2C%20muted%20blue%20and%20beige%20tones%2C%204k%20photograph&image_size=landscape_16_9'

// 用户证言区背景：开放式办公空间写实摄影
const testimonialBackgroundImage =
  'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20open%20plan%20office%20workspace%2C%20warm%20ambient%20light%2C%20blurred%20background%2C%20modern%20minimalist%20interior%2C%20neutral%20color%20palette%2C%20professional%20atmosphere%2C%204k%20photograph&image_size=landscape_16_9'

// 使用流程区背景：写字楼大堂写实摄影
const workflowBackgroundImage =
  'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=photorealistic%20modern%20corporate%20lobby%20atrium%2C%20clean%20geometric%20architecture%2C%20soft%20daylight%2C%20marble%20floor%2C%20glass%20walls%2C%20neutral%20tones%2C%20professional%20business%20atmosphere%2C%204k%20photograph&image_size=landscape_16_9'

// 功能卡片点击：未登录打开登录弹窗，已登录跳转对应路径
const handleFeatureClick = (feature: { path: string }) => {
  if (authStore.isAuthenticated) {
    router.push(feature.path)
  } else {
    openAuthModal('login')
  }
}

// 使用流程：4 步时间轴数据
const workflows = [
  {
    title: '注册账号',
    desc: '一键创建专属求职工作台，所有数据本地化部署，安全可控。',
    icon: markRaw(UserAddOutlined),
    points: ['邮箱注册，30 秒搞定', '本地化部署，数据不外泄'],
  },
  {
    title: '完善简历',
    desc: '粘贴目标岗位 JD，AI 智能生成结构化简历，一键润色表达。',
    icon: markRaw(EditOutlined),
    points: ['JD 匹配生成内容', 'AI 优化措辞表达', '多份简历灵活切换'],
  },
  {
    title: '模拟面试',
    desc: '多轮深度追问，支持语音与文字双模式，记录随时回看复盘。',
    icon: markRaw(VideoCameraOutlined),
    points: ['多场景面试题库', '语音文字双模式', '面试记录回放'],
  },
  {
    title: '投递追踪',
    desc: '看板式管理投递状态，记录面试轮次，洞察投递漏斗转化。',
    icon: markRaw(FundProjectionScreenOutlined),
    points: ['6 列状态看板', '面试轮次记录', '投递漏斗分析'],
  },
]

// 常见问题：折叠面板数据
const faqs = [
  {
    key: '1',
    q: '我的简历数据安全吗？',
    a: '职途 采用本地化部署模式，所有用户数据均存储在你自己的服务器上，不会上传至任何第三方。管理员账户通过配置文件管理，不入库，JWT 密钥隔离签发，保障账户安全。',
  },
  {
    key: '2',
    q: 'AI 生成的简历质量如何？',
    a: '简历实验室基于目标岗位 JD 进行智能匹配生成，内容更贴合岗位需求。同时提供 AI 润色与优化建议，可对单条经历进行多次迭代，直到满意为止。建议生成后结合个人实际情况微调，效果更佳。',
  },
  {
    key: '3',
    q: '模拟面试支持哪些岗位？',
    a: '模拟面试覆盖前端、后端、算法、产品、运营、设计等主流岗位，每个岗位提供多轮深度追问场景。支持语音与文字双模式输入，面试结束后可查看完整记录并随时回看复盘。',
  },
  {
    key: '4',
    q: '投递看板的数据来源是什么？',
    a: '投递看板的数据来源主要有两类：一是使用者手动新建的投递记录，可填写公司、岗位、薪资、渠道、优先级等字段；二是各大招聘平台（如 Boss 直聘、拉勾等）的投递数据，后续版本将支持批量导入，敬请期待。',
  },
]
</script>

<style scoped>
.home-page {
  font-family: var(--font-sans);
  color: var(--foreground);
  background: var(--background);
  overflow-x: hidden;
  min-height: 100vh;
  position: relative;
}

/* 整页底层背景图：极淡写实摄影，fixed 视差 */
.page-bg-layer {
  position: fixed;
  inset: 0;
  background-image: var(--page-bg);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  opacity: 0.08;
  pointer-events: none;
  z-index: 0;
}

/* 让所有内容浮在背景之上 */
.home-page > * {
  position: relative;
  z-index: 1;
}

/* ========== 顶部导航栏：固定 + 毛玻璃 ========== */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.nav-inner {
  max-width: 1200px;
  height: 100%;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: var(--primary);
  color: var(--primary-foreground);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.brand-name {
  font-weight: 700;
  font-size: 18px;
  letter-spacing: -0.01em;
  color: var(--foreground);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 32px;
}

.nav-links a {
  color: var(--foreground);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.15s ease;
}

.nav-links a:hover {
  color: var(--primary);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* ========== 胶囊按钮：border-radius 980px ========== */
.home-page .btn-primary,
.home-page .btn-secondary,
.home-page .btn-ghost {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-family: inherit;
  font-weight: 600;
  font-size: 15px;
  line-height: 1;
  padding: 10px 20px;
  border-radius: 980px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: transform 0.15s ease, background-color 0.15s ease,
    box-shadow 0.15s ease, opacity 0.15s ease;
  white-space: nowrap;
  text-decoration: none;
}

.home-page .btn-primary {
  background: var(--primary);
  color: var(--primary-foreground);
}

.home-page .btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
  background: var(--brand-600);
}

.home-page .btn-primary:focus-visible {
  outline: 2px solid var(--ring);
  outline-offset: 2px;
}

.home-page .btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.home-page .btn-secondary {
  background: var(--background-200);
  color: var(--foreground);
}

.home-page .btn-secondary:hover {
  transform: translateY(-1px);
  background: var(--background-300);
}

.home-page .btn-ghost {
  background: transparent;
  color: var(--foreground);
  padding: 10px 14px;
}

.home-page .btn-ghost:hover {
  background: var(--background-200);
}

.home-page .btn-lg {
  padding: 14px 28px;
  font-size: 16px;
}

/* ========== Hero 区 ========== */
.hero {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  padding: 96px 0 64px;
  overflow: hidden;
  background: linear-gradient(180deg, var(--background-100) 0%, var(--background) 100%);
}

/* Hero 写实摄影背景层 */
.hero-bg-layer {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  opacity: 0.18;
  z-index: 0;
}

.hero-bg-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    180deg,
    rgba(255, 255, 255, 0.78) 0%,
    rgba(255, 255, 255, 0.92) 70%,
    var(--background) 100%
  );
  z-index: 0;
}

/* 装饰性光晕 */
.hero-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.35;
  pointer-events: none;
  z-index: 0;
}

.hero-blob-1 {
  width: 480px;
  height: 480px;
  background: var(--brand-200);
  top: -120px;
  right: -80px;
}

.hero-blob-2 {
  width: 420px;
  height: 420px;
  background: var(--brand-100);
  bottom: -120px;
  left: -120px;
  opacity: 0.5;
}

.hero-container {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  gap: 48px;
}

.hero-left {
  flex: 1 1 0;
  max-width: 640px;
  min-width: 0;
}

/* ===== Hero 右侧产品预览浮卡 ===== */
.hero-right {
  flex: 0 0 420px;
  position: relative;
  height: 460px;
  max-width: 100%;
}

.preview-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--border);
  border-radius: 16px;
  box-shadow: var(--shadow-lg);
}

.preview-main {
  position: absolute;
  top: 20px;
  left: 0;
  right: 0;
  padding: 0;
  overflow: hidden;
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
}

.preview-cover {
  position: relative;
  display: block;
  width: 100%;
  height: 240px;
  overflow: hidden;
}

.preview-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.preview-cover-overlay {
  position: absolute;
  left: 14px;
  right: 14px;
  bottom: 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
}

.preview-cover-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.preview-cover-label {
  font-size: 12px;
  color: var(--muted-foreground);
}

.preview-cover-value {
  font-size: 13px;
  font-weight: 700;
  color: var(--primary);
}

.preview-cover-tags {
  display: flex;
  gap: 6px;
}

.preview-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.preview-title {
  margin-left: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--muted-foreground);
}

.preview-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-label {
  flex: 0 0 80px;
  font-size: 12px;
  color: var(--muted-foreground);
}

.preview-field {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--foreground);
}

.preview-bar {
  flex: 1;
  height: 6px;
  background: var(--background-200);
  border-radius: 999px;
  overflow: hidden;
}

.preview-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), var(--brand-400));
  border-radius: 999px;
}

.preview-pct {
  font-size: 12px;
  font-weight: 700;
  color: var(--primary);
}

.preview-tag {
  padding: 2px 8px;
  background: var(--brand-50);
  color: var(--primary);
  border-radius: 6px;
  font-size: 11px;
  font-weight: 500;
}

.preview-float {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  animation: float-y 4s ease-in-out infinite;
}

.preview-float-1 {
  top: -10px;
  right: -20px;
  animation-delay: 0s;
}

.preview-float-2 {
  bottom: 60px;
  right: -30px;
  animation-delay: 1s;
}

.preview-float-3 {
  bottom: -10px;
  left: 20px;
  animation-delay: 2s;
}

.preview-float-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.preview-float-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--foreground);
}

.preview-float-sub {
  font-size: 11px;
  color: var(--muted-foreground);
  margin-top: 1px;
}

@keyframes float-y {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

/* 顶部胶囊标签 */
.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 980px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border);
  font-size: 13px;
  font-weight: 500;
  color: var(--foreground);
  margin-bottom: 24px;
}

.hero-badge-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--primary);
  position: relative;
}

.hero-badge-dot::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  background: var(--primary);
  opacity: 0.35;
  animation: hero-pulse 2s ease-out infinite;
}

@keyframes hero-pulse {
  0% { transform: scale(0.6); opacity: 0.5; }
  100% { transform: scale(1.8); opacity: 0; }
}

/* 主标题：响应式缩放，目标 48px / 700 */
.hero-title {
  font-size: clamp(32px, 5vw, 48px);
  line-height: 1.1;
  font-weight: 700;
  letter-spacing: -0.02em;
  margin: 0 0 20px;
  color: var(--foreground);
}

.hero-title-accent {
  color: var(--primary);
}

/* 副标题：18px / muted */
.hero-desc {
  font-size: 18px;
  line-height: 1.6;
  color: var(--muted-foreground);
  max-width: 540px;
  margin: 0 0 32px;
}

/* 数据统计 */
.hero-stats {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 36px;
}

.hero-stat {
  flex: 1;
  min-width: 120px;
  padding: 16px 18px;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border);
}

.hero-stat-value {
  font-size: 24px;
  font-weight: 800;
  color: var(--primary);
  letter-spacing: -0.01em;
}

.hero-stat-label {
  font-size: 13px;
  color: var(--muted-foreground);
  margin-top: 2px;
}

/* 双 CTA 按钮容器 */
.hero-cta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

/* ========== Section 通用 ========== */
.section {
  padding: 96px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.section-head {
  text-align: center;
  max-width: 640px;
  margin: 0 auto 56px;
}

.section-eyebrow {
  display: inline-block;
  font-size: 13px;
  font-weight: 600;
  color: var(--primary);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  margin-bottom: 12px;
}

.section-title {
  font-size: clamp(28px, 3.5vw, 40px);
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--foreground);
  margin: 0 0 12px;
}

.section-sub {
  font-size: 16px;
  color: var(--muted-foreground);
  line-height: 1.6;
  margin: 0;
}

/* ========== 功能入口卡片：3 列网格 ========== */
.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

/* 卡片：白底 + 1px border + 16px 圆角 + shadow-sm，hover 上浮 */
.feature-card {
  display: flex;
  flex-direction: column;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 32px 28px 24px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

/* 卡片顶部装饰条 */
.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary), var(--brand-400));
  opacity: 0;
  transition: opacity 0.2s ease;
}

.feature-card:hover::before {
  opacity: 1;
}

.feature-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
  border-color: var(--brand-200);
}

/* 卡片图标：56x56 渐变圆角方块 */
.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--primary), var(--brand-600));
  color: var(--primary-foreground);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  margin-bottom: 20px;
  box-shadow: 0 8px 20px rgba(0, 122, 255, 0.25);
  transition: transform 0.2s ease;
}

.feature-card:hover .feature-icon {
  transform: scale(1.06);
}

.feature-body {
  flex: 1;
}

.feature-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--foreground);
  margin: 0 0 8px;
  letter-spacing: -0.01em;
}

.feature-desc {
  font-size: 14px;
  color: var(--muted-foreground);
  line-height: 1.65;
  margin: 0 0 4px;
}

/* 卡片底部箭头：hover 时右移 4px */
.feature-arrow {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px dashed var(--border);
  display: inline-flex;
  align-items: center;
  color: var(--primary);
  font-size: 14px;
  font-weight: 600;
  transition: transform 0.2s ease;
}

.feature-arrow::after {
  content: '查看详情';
  margin-right: 6px;
}

.feature-card:hover .feature-arrow {
  transform: translateX(4px);
}

/* 功能卡片特性标签 */
.feature-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin: 12px 0 14px;
}

.feature-tag {
  padding: 3px 10px;
  background: var(--brand-50);
  color: var(--primary);
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.01em;
}

/* 功能卡片亮点列表 */
.feature-highlights {
  list-style: none;
  padding: 14px 14px;
  margin: 14px 0 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: var(--background-50);
  border-radius: 12px;
  border: 1px solid var(--border);
}

.feature-highlights li {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 13px;
  color: var(--foreground);
  line-height: 1.5;
}

.feature-check {
  color: var(--success);
  font-size: 14px;
  flex-shrink: 0;
  margin-top: 2px;
}

/* ========== 核心能力区 ========== */
.section-tint {
  background: var(--background-100);
}

.capability-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.capability-card {
  display: flex;
  gap: 14px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 14px;
  padding: 22px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.capability-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-200);
}

.capability-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--brand-50);
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.capability-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--foreground);
  margin: 0 0 4px;
  letter-spacing: -0.01em;
}

.capability-desc {
  font-size: 13px;
  color: var(--muted-foreground);
  line-height: 1.55;
  margin: 0;
}

/* ========== 数据成果区（深色渐变） ========== */
.section-grad {
  background: linear-gradient(135deg, #0a2540 0%, #1e3a8a 50%, #007aff 100%);
  color: #fff;
  position: relative;
  overflow: hidden;
}

.section-grad::before {
  content: '';
  position: absolute;
  top: -200px;
  right: -200px;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  filter: blur(60px);
  pointer-events: none;
}

.section-grad::after {
  content: '';
  position: absolute;
  bottom: -150px;
  left: -150px;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: rgba(52, 199, 89, 0.12);
  filter: blur(60px);
  pointer-events: none;
}

.section-head-light {
  position: relative;
  z-index: 1;
}

.section-eyebrow-light {
  color: rgba(255, 255, 255, 0.7);
}

.section-title-light {
  color: #fff;
}

.section-sub-light {
  color: rgba(255, 255, 255, 0.75);
}

.achievement-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  position: relative;
  z-index: 1;
}

.achievement-card {
  text-align: center;
  padding: 32px 16px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  transition: transform 0.2s ease, background 0.2s ease;
}

.achievement-card:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.12);
}

.achievement-value {
  font-size: 44px;
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.02em;
  line-height: 1;
  margin-bottom: 8px;
}

.achievement-label {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.achievement-sub {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.65);
  line-height: 1.4;
}

/* ========== 用户证言区 ========== */
.testimonial-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.testimonial-card {
  display: flex;
  flex-direction: column;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 24px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.testimonial-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-200);
}

.testimonial-stars {
  display: flex;
  gap: 2px;
  color: #ffb800;
  font-size: 14px;
  margin-bottom: 14px;
}

.testimonial-quote {
  font-size: 14px;
  line-height: 1.7;
  color: var(--foreground);
  margin: 0 0 18px;
  flex: 1;
  font-style: italic;
}

.testimonial-author {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}

.testimonial-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  border: 1px solid var(--border);
  background: var(--background-200);
}

.testimonial-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--foreground);
}

.testimonial-role {
  font-size: 12px;
  color: var(--muted-foreground);
  margin-top: 1px;
}

/* ========== 使用流程：4 步时间轴 ========== */
.section-alt {
  position: relative;
  background: var(--background-100);
}

/* 带写实摄影背景的 section */
.section-photo {
  position: relative;
  background: var(--background);
}

.section-bg-layer {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  opacity: 0.12;
  z-index: 0;
  pointer-events: none;
}

.section-photo .container,
.section-alt .container {
  position: relative;
  z-index: 1;
}

/* 在带背景的 section 上叠加白色遮罩，保证内容可读 */
.section-photo::after,
.section-alt::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--background) 0%, rgba(255, 255, 255, 0.85) 50%, var(--background) 100%);
  z-index: 0;
  pointer-events: none;
}

.workflow-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.workflow-card {
  position: relative;
  display: flex;
  flex-direction: column;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 28px 24px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.workflow-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-200);
}

/* 序号 + 连接线容器 */
.workflow-index {
  display: flex;
  align-items: center;
  margin-bottom: 18px;
  height: 32px;
}

.workflow-index-num {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--primary);
  color: var(--primary-foreground);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 700;
  flex-shrink: 0;
}

.workflow-index-line {
  flex: 1;
  height: 2px;
  margin-left: 12px;
  background: linear-gradient(
    90deg,
    var(--brand-300) 0%,
    var(--border) 100%
  );
  border-radius: 1px;
}

.workflow-content {
  flex: 1;
}

.workflow-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--brand-50);
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  margin-bottom: 14px;
}

.workflow-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--foreground);
  margin: 0 0 6px;
  letter-spacing: -0.01em;
}

.workflow-desc {
  font-size: 14px;
  color: var(--muted-foreground);
  line-height: 1.6;
  margin: 0 0 12px;
}

.workflow-points {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.workflow-points li {
  position: relative;
  padding-left: 18px;
  font-size: 13px;
  color: var(--muted-foreground);
  line-height: 1.5;
}

.workflow-points li::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 8px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--primary);
  opacity: 0.7;
}

/* ========== 常见问题：折叠面板 ========== */
.container-narrow {
  max-width: 820px;
}

.faq-collapse {
  background: transparent;
}

.faq-collapse :deep(.ant-collapse-item) {
  background: var(--card);
  border: 1px solid var(--border) !important;
  border-radius: 12px !important;
  margin-bottom: 12px;
  overflow: hidden;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.faq-collapse :deep(.ant-collapse-item:hover) {
  border-color: var(--brand-200) !important;
}

.faq-collapse :deep(.ant-collapse-item-active) {
  border-color: var(--brand-200) !important;
  box-shadow: var(--shadow-sm);
}

.faq-collapse :deep(.ant-collapse-header) {
  padding: 18px 22px !important;
  align-items: center !important;
  font-size: 15px;
  font-weight: 600;
  color: var(--foreground) !important;
}

.faq-collapse :deep(.ant-collapse-content-box) {
  padding: 0 22px 18px !important;
}

.faq-collapse :deep(.ant-collapse-arrow) {
  color: var(--muted-foreground) !important;
  font-size: 13px !important;
}

.faq-answer {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--muted-foreground);
}

/* FAQ 底部 CTA */
.faq-cta {
  margin-top: 32px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.faq-cta-text {
  font-size: 14px;
  color: var(--muted-foreground);
}

/* ========== 页脚 ========== */
.footer {
  background: var(--background-800);
  padding: 56px 0 32px;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 40px;
  margin-bottom: 32px;
}

.footer-col {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.footer-col-brand {
  gap: 8px;
}

.footer-brand {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.footer-logo {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  background: var(--primary);
  color: var(--primary-foreground);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.footer-brand span {
  font-weight: 700;
  font-size: 16px;
  color: var(--text-50);
}

.footer-tagline {
  color: var(--text-300);
  font-size: 14px;
  margin: 0;
}

.footer-copy {
  color: var(--text-400);
  font-size: 13px;
  margin: 8px 0 0;
}

.footer-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-50);
  margin: 0 0 10px;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

.footer-col a {
  color: var(--text-300);
  font-size: 14px;
  text-decoration: none;
  cursor: pointer;
  transition: color 0.15s ease;
  padding: 3px 0;
}

.footer-col a:hover {
  color: var(--primary);
}

/* ========== 认证弹窗 ========== */
.auth-modal :deep(.ant-modal-content) {
  border-radius: var(--radius-lg);
  padding: 32px 28px 24px;
  background: var(--card);
  box-shadow: var(--shadow-2xl);
}

.auth-modal :deep(.ant-modal-close) {
  top: 16px;
  right: 16px;
  color: var(--muted-foreground);
}

.auth-modal-body {
  width: 100%;
}

.auth-modal-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--foreground);
  margin: 0 0 4px;
  letter-spacing: -0.01em;
}

.auth-modal-sub {
  font-size: 14px;
  color: var(--muted-foreground);
  margin: 0 0 24px;
}

/* 登录 / 注册 切换 tabs：胶囊分段控件 */
.auth-tabs {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: var(--background-200);
  border-radius: 980px;
  margin-bottom: 24px;
}

.auth-tab {
  flex: 1;
  padding: 10px;
  border: none;
  background: transparent;
  border-radius: 980px;
  font-family: inherit;
  font-weight: 600;
  font-size: 14px;
  color: var(--muted-foreground);
  cursor: pointer;
  transition: all 0.15s ease;
}

.auth-tab-active {
  background: var(--card);
  color: var(--foreground);
  box-shadow: var(--shadow-sm);
}

/* 表单样式覆盖 */
.auth-form :deep(.ant-form-item) {
  margin-bottom: 16px;
}

.auth-form :deep(.ant-form-item-label > label) {
  font-size: 13px;
  font-weight: 500;
  color: var(--foreground);
}

.auth-form :deep(.ant-input-affix-wrapper),
.auth-form :deep(.ant-input) {
  border-radius: var(--radius-md);
  border-color: var(--border);
  background: var(--background-50);
  color: var(--foreground);
}

.auth-form :deep(.ant-input-affix-wrapper:hover) {
  border-color: var(--primary);
}

.auth-form :deep(.ant-input-affix-wrapper:focus),
.auth-form :deep(.ant-input-affix-wrapper-focused) {
  border-color: var(--primary) !important;
  box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.12) !important;
}

.auth-form :deep(.ant-input::placeholder) {
  color: var(--muted-foreground);
}

.input-icon {
  color: var(--muted-foreground);
}

/* 提交按钮：100% 宽度 */
.auth-submit {
  width: 100%;
  margin-top: 6px;
}

.auth-submit.is-loading {
  opacity: 0.75;
}

.auth-switch {
  text-align: center;
  color: var(--muted-foreground);
  font-size: 13px;
}

.auth-switch a {
  color: var(--primary);
  cursor: pointer;
  margin-left: 4px;
  font-weight: 500;
}

.auth-switch a:hover {
  color: var(--brand-600);
}

/* ========== 响应式 ========== */
/* 768-1023px：2 列网格 */
@media (max-width: 1023px) {
  .feature-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .workflow-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  /* 第 2 个卡片后隐藏连接线，避免 2x2 布局错位 */
  .workflow-card:nth-child(2) .workflow-index-line,
  .workflow-card:nth-child(4) .workflow-index-line {
    display: none;
  }

  /* Hero 右侧预览区在小屏隐藏，左内容居中 */
  .hero-right {
    display: none;
  }

  .hero-container {
    justify-content: center;
  }

  .hero-left {
    max-width: 640px;
    flex: 1 1 auto;
  }

  /* 核心能力 2 列 */
  .capability-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  /* 数据成果 2 列 */
  .achievement-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  /* 证言 2 列 */
  .testimonial-grid {
    grid-template-columns: 1fr;
    max-width: 560px;
    margin: 0 auto;
  }

  /* Footer 调整 */
  .footer-grid {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }
}

/* <768px：1 列网格，Hero 标题缩放 */
@media (max-width: 767px) {
  .nav-links {
    display: none;
  }

  .hero {
    padding: 88px 0 56px;
    min-height: auto;
  }

  .hero-stats {
    flex-direction: column;
  }

  .hero-stat {
    min-width: 0;
  }

  .hero-cta {
    flex-direction: column;
    align-items: stretch;
  }

  .hero-cta .btn-lg {
    width: 100%;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .workflow-grid {
    grid-template-columns: 1fr;
  }

  /* 单列布局下隐藏所有连接线 */
  .workflow-index-line {
    display: none;
  }

  .section {
    padding: 72px 0;
  }

  /* 核心能力 1 列 */
  .capability-grid {
    grid-template-columns: 1fr;
  }

  /* 数据成果 2 列（保持紧凑） */
  .achievement-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .achievement-card {
    padding: 24px 12px;
  }

  .achievement-value {
    font-size: 36px;
  }

  /* Footer 单列 */
  .footer-grid {
    grid-template-columns: 1fr;
    gap: 28px;
  }
}

/* ========== 减少动效（无障碍） ========== */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation: none !important;
    transition: none !important;
  }
}
</style>
