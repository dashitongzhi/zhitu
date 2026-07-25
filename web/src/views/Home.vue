<template>
  <div class="home-page">
    <!-- 顶部导航栏 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="brand">
          <img src="/favicon.svg" alt="职途AI" class="brand-logo" />
          <span class="brand-name">职途AI</span>
        </div>
        <nav class="nav-links">
          <a href="#features">功能特性</a>
          <a href="#workflow">使用流程</a>
          <a href="#faq">常见问题</a>
        </nav>
        <div class="nav-actions">
          <a-button type="text" @click="scrollToAuth">登录</a-button>
          <a-button type="primary" @click="switchToRegister">免费注册</a-button>
        </div>
      </div>
    </header>

    <!-- Hero 区 -->
    <section class="hero">
      <div class="hero-bg">
        <div class="blob blob-1"></div>
        <div class="blob blob-2"></div>
        <div class="blob blob-3"></div>
      </div>

      <div class="hero-container">
        <!-- 左侧文案 -->
        <div class="hero-content">
          <div class="hero-badge">
            <span class="badge-dot"></span>
            AI 驱动的智能求职助手
          </div>
          <h1 class="hero-title">
            让每一次求职<br />
            都<span class="highlight">精准而高效</span>
          </h1>
          <p class="hero-desc">
            职途AI 帮你管理个人档案、智能生成与润色简历、模拟真实面试场景、追踪投递进度。
            从准备到入职，一站式陪伴你的求职之旅。
          </p>

          <div class="hero-stats">
            <div class="stat-item">
              <div class="stat-num">4+</div>
              <div class="stat-label">核心模块</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <div class="stat-num">AI</div>
              <div class="stat-label">智能生成</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <div class="stat-num">100%</div>
              <div class="stat-label">本地化部署</div>
            </div>
          </div>

          <div class="hero-cta">
            <a-button type="primary" size="large" @click="switchToRegister">
              立即开始使用
              <template #icon>
                <ArrowRightOutlined />
              </template>
            </a-button>
            <a-button size="large" ghost @click="scrollToFeatures">
              了解更多
            </a-button>
          </div>
        </div>

        <!-- 右侧登录/注册卡片 -->
        <div class="auth-card" id="auth-card">
          <div class="auth-card-header">
            <div class="auth-tabs">
              <button
                class="auth-tab"
                :class="{ active: activeTab === 'login' }"
                @click="activeTab = 'login'"
              >
                登录
              </button>
              <button
                class="auth-tab"
                :class="{ active: activeTab === 'register' }"
                @click="activeTab = 'register'"
              >
                注册
              </button>
              <div
                class="auth-tab-indicator"
                :class="{ 'indicator-right': activeTab === 'register' }"
              ></div>
            </div>
          </div>

          <!-- 登录表单 -->
          <a-form
            v-if="activeTab === 'login'"
            :model="loginForm"
            :rules="loginRules"
            @finish="handleLogin"
            layout="vertical"
            class="auth-form"
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
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <LockOutlined class="input-icon" />
                </template>
              </a-input-password>
            </a-form-item>

            <a-form-item>
              <a-button
                type="primary"
                html-type="submit"
                size="large"
                block
                :loading="loading"
              >
                登录
              </a-button>
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
            @finish="handleRegister"
            layout="vertical"
            class="auth-form"
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
                placeholder="请输入密码（不少于6位）"
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
              <a-button
                type="primary"
                html-type="submit"
                size="large"
                block
                :loading="loading"
              >
                注册
              </a-button>
            </a-form-item>

            <div class="auth-switch">
              已有账号？
              <a @click="activeTab = 'login'">立即登录</a>
            </div>
          </a-form>
        </div>
      </div>
    </section>

    <!-- 功能特性区 -->
    <section class="features" id="features">
      <div class="section-inner">
        <div class="section-header">
          <div class="section-tag">核心功能</div>
          <h2 class="section-title">一站式求职解决方案</h2>
          <p class="section-desc">
            从个人档案到投递跟踪，覆盖求职全流程，让 AI 为你的每一步赋能
          </p>
        </div>

        <div class="feature-grid">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="feature-card"
          >
            <div class="feature-icon" :style="{ background: feature.gradient }">
              <component :is="feature.icon" />
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-desc">{{ feature.desc }}</p>
            <ul class="feature-list">
              <li v-for="item in feature.points" :key="item">
                <CheckCircleFilled class="check-icon" />
                <span>{{ item }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- 使用流程区 -->
    <section class="workflow" id="workflow">
      <div class="section-inner">
        <div class="section-header">
          <div class="section-tag">使用流程</div>
          <h2 class="section-title">四步开启你的求职之旅</h2>
          <p class="section-desc">简单几步，让 AI 帮你打造专属求职策略</p>
        </div>

        <div class="workflow-steps">
          <div
            v-for="(step, index) in steps"
            :key="step.title"
            class="step-item"
          >
            <div class="step-num">{{ String(index + 1).padStart(2, '0') }}</div>
            <div class="step-line" v-if="index < steps.length - 1"></div>
            <div class="step-content">
              <h3 class="step-title">{{ step.title }}</h3>
              <p class="step-desc">{{ step.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="faq" id="faq">
      <div class="section-inner">
        <div class="section-header">
          <div class="section-tag">常见问题</div>
          <h2 class="section-title">关于职途AI</h2>
        </div>

        <a-collapse :bordered="false" class="faq-collapse">
          <a-collapse-panel
            v-for="(item, index) in faqs"
            :key="index"
            :header="item.q"
          >
            <p class="faq-answer">{{ item.a }}</p>
          </a-collapse-panel>
        </a-collapse>
      </div>
    </section>

    <!-- CTA 区 -->
    <section class="cta-section">
      <div class="cta-inner">
        <h2 class="cta-title">准备好开启智能求职之旅了吗？</h2>
        <p class="cta-desc">注册即可免费使用全部功能</p>
        <a-button
          type="primary"
          size="large"
          class="cta-btn"
          @click="switchToRegister"
        >
          立即免费注册
        </a-button>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="footer-inner">
        <div class="footer-brand">
          <img src="/favicon.svg" alt="职途AI" class="footer-logo" />
          <span>职途AI</span>
        </div>
        <p class="footer-text">让求职更简单 · 智能求职助手平台</p>
        <p class="footer-copy">© 2024 职途AI. All rights reserved.</p>
      </div>
    </footer>
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
  CheckCircleFilled,
  ProfileOutlined,
  FileTextOutlined,
  MessageOutlined,
  SendOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const authStore = useAuthStore()

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
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
}

// 注册表单
const registerForm = reactive({
  email: '',
  nickname: '',
  password: '',
  confirmPassword: '',
})

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
    { min: 2, max: 20, message: '昵称长度为2-20个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

// 处理登录
const handleLogin = async () => {
  loading.value = true
  try {
    const success = await authStore.login(loginForm)
    if (success) {
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

// 滚动到认证卡片
const scrollToAuth = () => {
  const el = document.getElementById('auth-card')
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

// 滚动到功能区
const scrollToFeatures = () => {
  const el = document.getElementById('features')
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

// 切换到注册 tab 并滚动
const switchToRegister = () => {
  activeTab.value = 'register'
  scrollToAuth()
}

// 功能特性数据
const features = [
  {
    title: '个人档案管理',
    desc: '完整记录你的教育背景、工作经历、项目经验等信息，AI 实时评估档案完成度。',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    icon: markRaw(ProfileOutlined),
    points: ['多维信息结构化存储', '完成度智能评估', '简历生成数据源'],
  },
  {
    title: '智能简历生成',
    desc: '基于个人档案与目标岗位 JD，AI 自动生成、润色简历，并提供专业评分建议。',
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    icon: markRaw(FileTextOutlined),
    points: ['针对 JD 定制生成', 'AI 润色与优化', '多份简历管理'],
  },
  {
    title: 'AI 模拟面试',
    desc: '支持语音/文字多种模式，模拟真实企业面试场景，每次面试形成完整记录。',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    icon: markRaw(MessageOutlined),
    points: ['多轮次深度追问', '语音+文字双模式', '面试记录留存'],
  },
  {
    title: '投递看板追踪',
    desc: '可视化管理求职投递全流程，从投递到 offer，状态流转一目了然。',
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    icon: markRaw(SendOutlined),
    points: ['看板式状态管理', '面试轮次记录', '投递漏斗分析'],
  },
]

// 流程步骤
const steps = [
  {
    title: '创建账号',
    desc: '邮箱注册即可开始使用，无需复杂配置',
  },
  {
    title: '完善档案',
    desc: '填写个人教育、工作、项目等信息，提升档案完成度',
  },
  {
    title: '生成简历',
    desc: '选择目标岗位，AI 自动生成定制化简历并智能润色',
  },
  {
    title: '模拟面试 & 投递',
    desc: '通过 AI 面试演练提升能力，追踪每一次投递进度',
  },
]

// 常见问题
const faqs = [
  {
    q: '职途AI 是免费的吗？',
    a: '本平台为本地化部署版本，注册后即可使用全部核心功能，包括档案管理、简历生成、模拟面试和投递看板。',
  },
  {
    q: 'AI 功能如何工作？',
    a: '职途AI 接入 OpenAI 兼容的 LLM 接口，结合你的个人档案和目标岗位信息，智能生成简历内容、模拟面试问题，并提供专业的评分建议。',
  },
  {
    q: '我的数据安全吗？',
    a: '所有数据存储在本地 SQLite 数据库中，简历、面试记录等文件也保存在本地服务器，不会上传到第三方。请妥善保管服务器访问权限。',
  },
  {
    q: '如何开始使用？',
    a: '点击页面右侧的注册按钮创建账号，登录后进入控制台即可开始完善个人档案、生成简历、模拟面试和投递追踪。',
  },
]
</script>

<style scoped>
.home-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC',
    'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  color: #1a1a2e;
  overflow-x: hidden;
  /* 覆盖全局蓝色主色为紫色主题，影响 var(--primary-color) 引用 */
  --primary-color: #667eea;
}

/* ========== 顶部导航栏 ========== */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.nav-inner {
  max-width: 1280px;
  height: 100%;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-logo {
  width: 32px;
  height: 32px;
}

.brand-name {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-links {
  display: flex;
  gap: 32px;
}

.nav-links a {
  color: #555;
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  transition: color 0.2s;
}

.nav-links a:hover {
  color: #667eea;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 导航栏主按钮：使用紫色渐变，覆盖 Ant Design 默认蓝色 */
.nav-actions :deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.nav-actions :deep(.ant-btn-primary):hover,
.nav-actions :deep(.ant-btn-primary):focus {
  background: linear-gradient(135deg, #5a72e0 0%, #6a3f95 100%);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

/* ========== Hero 区 ========== */
.hero {
  position: relative;
  min-height: 100vh;
  padding: 120px 32px 80px;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7ff 0%, #ede9fe 50%, #f5f3ff 100%);
}

.hero-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 12s ease-in-out infinite;
}

.blob-1 {
  width: 500px;
  height: 500px;
  background: #667eea;
  top: -150px;
  left: -100px;
}

.blob-2 {
  width: 400px;
  height: 400px;
  background: #f093fb;
  top: 30%;
  right: -100px;
  animation-delay: -4s;
}

.blob-3 {
  width: 350px;
  height: 350px;
  background: #4facfe;
  bottom: -100px;
  left: 40%;
  animation-delay: -8s;
}

@keyframes float {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(40px, -40px) scale(1.1);
  }
  66% {
    transform: translate(-30px, 30px) scale(0.95);
  }
}

.hero-container {
  position: relative;
  max-width: 1280px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 440px;
  gap: 64px;
  align-items: center;
}

.hero-content {
  max-width: 600px;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 100px;
  font-size: 13px;
  color: #667eea;
  font-weight: 500;
  margin-bottom: 24px;
}

.badge-dot {
  width: 8px;
  height: 8px;
  background: #667eea;
  border-radius: 50%;
  box-shadow: 0 0 12px #667eea;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.hero-title {
  font-size: 56px;
  line-height: 1.2;
  font-weight: 800;
  margin: 0 0 24px;
  letter-spacing: -1px;
  color: #1a1a2e;
}

.highlight {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-desc {
  font-size: 17px;
  line-height: 1.7;
  color: #555;
  margin: 0 0 32px;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 40px;
  padding: 20px 28px;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.8);
  border-radius: 16px;
  max-width: fit-content;
}

.stat-item {
  text-align: center;
}

.stat-num {
  font-size: 26px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-label {
  font-size: 12px;
  color: #888;
  margin-top: 2px;
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: rgba(0, 0, 0, 0.1);
}

.hero-cta {
  display: flex;
  gap: 12px;
}

.hero-cta :deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  height: 48px;
  padding: 0 28px;
  font-size: 15px;
  font-weight: 500;
}

.hero-cta :deep(.ant-btn-primary):hover,
.hero-cta :deep(.ant-btn-primary):focus {
  background: linear-gradient(135deg, #5a72e0 0%, #6a3f95 100%);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.5);
  transform: translateY(-1px);
}

.hero-cta :deep(.ant-btn-background-ghost) {
  border-color: #667eea;
  color: #667eea;
  height: 48px;
  padding: 0 28px;
  font-size: 15px;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(6px);
}

.hero-cta :deep(.ant-btn-background-ghost):hover,
.hero-cta :deep(.ant-btn-background-ghost):focus {
  border-color: #5a72e0;
  color: #5a72e0;
  background: rgba(255, 255, 255, 0.85);
}

/* ========== 认证卡片 ========== */
.auth-card {
  position: relative;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border-radius: 24px;
  padding: 36px 32px;
  box-shadow: 0 32px 80px rgba(102, 126, 234, 0.18),
    0 12px 32px rgba(31, 38, 135, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.6);
  overflow: hidden;
}

.auth-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
}

.auth-card-header {
  margin-bottom: 24px;
}

.auth-tabs {
  position: relative;
  display: flex;
  background: #f5f5f7;
  border-radius: 12px;
  padding: 4px;
}

.auth-tab {
  flex: 1;
  position: relative;
  z-index: 2;
  padding: 10px 0;
  background: transparent;
  border: none;
  font-size: 15px;
  font-weight: 600;
  color: #888;
  cursor: pointer;
  transition: color 0.3s;
}

.auth-tab.active {
  color: #fff;
}

.auth-tab-indicator {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.auth-tab-indicator.indicator-right {
  transform: translateX(100%);
}

.auth-form {
  animation: fadeIn 0.4s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.auth-form :deep(.ant-form-item) {
  margin-bottom: 16px;
}

/* 独立输入框（无前缀图标）与外层 affix 包裹容器统一高度 */
.auth-form :deep(.ant-input-affix-wrapper),
.auth-form :deep(> * .ant-input:not(.ant-input-affix-wrapper .ant-input)) {
  border-radius: 10px;
}

.auth-form :deep(.ant-input-affix-wrapper) {
  height: 44px;
  display: flex;
  align-items: center;
  padding-top: 0;
  padding-bottom: 0;
}

/* 独立输入框（如无 prefix 的场景）保持固定高度 */
.auth-form :deep(.ant-input:not(.ant-input-affix-wrapper .ant-input)) {
  height: 44px;
  border-radius: 10px;
}

/* affix 包裹内部真实 input 不再设固定高度，交由 flex 垂直居中，修复 placeholder 错位 */
.auth-form :deep(.ant-input-affix-wrapper > .ant-input) {
  height: auto;
  line-height: 1.5;
  background: transparent;
}

/* 输入框统一玻璃质感底色与柔和边框 */
.auth-form :deep(.ant-input-affix-wrapper),
.auth-form :deep(.ant-input:not(.ant-input-affix-wrapper .ant-input)) {
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(102, 126, 234, 0.16);
  transition: border-color 0.25s, box-shadow 0.25s, background 0.25s;
}

.auth-form :deep(.ant-input-affix-wrapper .ant-input-prefix) {
  margin-inline-end: 10px;
}

.auth-form :deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  height: 44px;
  border-radius: 10px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.auth-form :deep(.ant-btn-primary):hover,
.auth-form :deep(.ant-btn-primary):focus {
  background: linear-gradient(135deg, #5a72e0 0%, #6a3f95 100%);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

/* 输入框聚焦态：覆盖全局蓝色为紫色主题 */
.auth-form :deep(.ant-input:focus),
.auth-form :deep(.ant-input-focused),
.auth-form :deep(.ant-input-affix-wrapper:focus),
.auth-form :deep(.ant-input-affix-wrapper-focused) {
  border-color: #667eea !important;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) !important;
}

.auth-form :deep(.ant-input-affix-wrapper):hover {
  border-color: #667eea !important;
}

.input-icon {
  color: #bfbfbf;
}

.auth-switch {
  text-align: center;
  color: #999;
  font-size: 13px;
}

.auth-switch a {
  color: #667eea;
  cursor: pointer;
  margin-left: 4px;
}

.auth-switch a:hover {
  color: #5a72e0;
  text-decoration: underline;
}

/* ========== 通用 section 样式 ========== */
.section-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 100px 32px;
}

.section-header {
  text-align: center;
  margin-bottom: 64px;
}

.section-tag {
  display: inline-block;
  padding: 6px 16px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  font-size: 13px;
  font-weight: 600;
  border-radius: 100px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 40px;
  font-weight: 800;
  margin: 0 0 16px;
  letter-spacing: -0.5px;
  color: #1a1a2e;
}

.section-desc {
  font-size: 16px;
  color: #777;
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.7;
}

/* ========== 功能特性区 ========== */
.features {
  background: #fff;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.feature-card {
  padding: 36px;
  background: linear-gradient(180deg, #ffffff 0%, #fafbff 100%);
  border: 1px solid #f0f0f5;
  border-radius: 20px;
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 48px rgba(102, 126, 234, 0.12);
  border-color: rgba(102, 126, 234, 0.2);
}

.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 26px;
  margin-bottom: 20px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

.feature-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 12px;
  color: #1a1a2e;
}

.feature-desc {
  font-size: 15px;
  color: #666;
  line-height: 1.7;
  margin: 0 0 20px;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
  font-size: 14px;
  color: #555;
}

.check-icon {
  color: #52c41a;
  font-size: 16px;
  flex-shrink: 0;
}

/* ========== 使用流程区 ========== */
.workflow {
  background: linear-gradient(180deg, #fafbff 0%, #f5f7ff 100%);
}

.workflow-steps {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 32px;
  position: relative;
}

.step-item {
  position: relative;
  padding: 32px;
  background: #fff;
  border-radius: 20px;
  border: 1px solid #f0f0f5;
  transition: all 0.3s;
}

.step-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 40px rgba(102, 126, 234, 0.1);
}

.step-num {
  font-size: 36px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 16px;
  line-height: 1;
}

.step-line {
  position: absolute;
  top: 48px;
  right: -20px;
  width: 40px;
  height: 2px;
  background: linear-gradient(90deg, #667eea, transparent);
  z-index: 1;
}

.step-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 10px;
  color: #1a1a2e;
}

.step-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.7;
  margin: 0;
}

/* ========== FAQ ========== */
.faq {
  background: #fff;
}

.faq :deep(.ant-collapse) {
  max-width: 800px;
  margin: 0 auto;
  background: transparent;
}

.faq :deep(.ant-collapse-item) {
  margin-bottom: 16px;
  border: 1px solid #f0f0f5 !important;
  border-radius: 16px !important;
  overflow: hidden;
  background: #fff;
  transition: all 0.3s;
}

.faq :deep(.ant-collapse-item:hover) {
  border-color: rgba(102, 126, 234, 0.3) !important;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.08);
}

.faq :deep(.ant-collapse-header) {
  padding: 20px 24px !important;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  align-items: center !important;
}

.faq :deep(.ant-collapse-arrow) {
  color: #667eea !important;
}

.faq :deep(.ant-collapse-item-active) .ant-collapse-header {
  color: #667eea !important;
}

.faq :deep(.ant-collapse-content) {
  padding: 0 24px 20px !important;
  border-top: 1px solid #f5f5f5;
  background: #fff !important;
}

.faq-answer {
  margin: 16px 0 0;
  color: #666;
  line-height: 1.7;
  font-size: 14px;
}

/* ========== CTA 区 ========== */
.cta-section {
  padding: 80px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.cta-section::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(
    circle at 20% 50%,
    rgba(255, 255, 255, 0.15) 0%,
    transparent 50%
  );
}

.cta-inner {
  position: relative;
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
}

.cta-title {
  font-size: 36px;
  font-weight: 800;
  color: #fff;
  margin: 0 0 16px;
}

.cta-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0 0 32px;
}

.cta-btn {
  height: 48px !important;
  padding: 0 36px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
  background: #fff !important;
  color: #667eea !important;
  border: none !important;
  border-radius: 100px !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2) !important;
}

.cta-btn:hover,
.cta-btn:focus {
  transform: translateY(-2px);
  background: #fff !important;
  color: #5a72e0 !important;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.25) !important;
}

/* ========== 页脚 ========== */
.footer {
  background: #1a1a2e;
  padding: 48px 32px;
  color: rgba(255, 255, 255, 0.7);
}

.footer-inner {
  max-width: 1280px;
  margin: 0 auto;
  text-align: center;
}

.footer-brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.footer-logo {
  width: 28px;
  height: 28px;
}

.footer-brand span {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.footer-text {
  margin: 0 0 8px;
  font-size: 14px;
}

.footer-copy {
  margin: 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

/* ========== 响应式 ========== */
@media (max-width: 1024px) {
  .hero-container {
    grid-template-columns: 1fr;
    gap: 48px;
  }

  .auth-card {
    max-width: 440px;
    margin: 0 auto;
  }

  .hero-title {
    font-size: 44px;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .workflow-steps {
    grid-template-columns: repeat(2, 1fr);
  }

  .step-line {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .hero {
    padding: 100px 20px 60px;
  }

  .hero-title {
    font-size: 34px;
  }

  .hero-stats {
    padding: 16px 20px;
    gap: 16px;
  }

  .stat-num {
    font-size: 20px;
  }

  .hero-cta {
    flex-direction: column;
  }

  .section-inner {
    padding: 64px 20px;
  }

  .section-title {
    font-size: 28px;
  }

  .workflow-steps {
    grid-template-columns: 1fr;
  }

  .cta-title {
    font-size: 26px;
  }
}
</style>
