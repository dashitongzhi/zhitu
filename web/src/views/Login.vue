<template>
  <div class="login-container">
    <div class="login-card">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <img src="/favicon.svg" alt="职途" class="logo" />
        <h1 class="title">职途</h1>
        <p class="subtitle">智能求职助手平台</p>
      </div>

      <!-- 登录表单 -->
      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleLogin"
        layout="vertical"
        class="login-form"
      >
        <a-form-item label="邮箱" name="email">
          <a-input
            v-model:value="formState.email"
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
            v-model:value="formState.password"
            placeholder="请输入密码"
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
            登录
          </a-button>
        </a-form-item>

        <div class="login-footer">
          <span>还没有账号？</span>
          <router-link to="/register" class="register-link">立即注册</router-link>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { MailOutlined, LockOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 加载状态
const loading = ref(false)

// 表单数据
const formState = reactive({
  email: '',
  password: '',
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
}

// 处理登录
const handleLogin = async () => {
  loading.value = true
  try {
    const success = await authStore.login(formState)
    if (success) {
      // 检查是否有重定向地址
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #999;
}

.login-form {
  margin-top: 24px;
}

.input-icon {
  color: #bfbfbf;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  color: #999;
  font-size: 14px;
}

.register-link {
  color: #1890ff;
  margin-left: 4px;
}

.register-link:hover {
  text-decoration: underline;
}
</style>