<template>
  <div class="change-password-page">
    <a-card title="修改密码" class="pwd-card" :bordered="false">
      <a-form
        ref="formRef"
        :model="form"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item label="当前密码" name="old_password">
          <a-input-password
            v-model:value="form.old_password"
            placeholder="请输入当前密码"
            autocomplete="current-password"
          />
        </a-form-item>
        <a-form-item label="新密码" name="new_password">
          <a-input-password
            v-model:value="form.new_password"
            placeholder="至少 8 位"
            autocomplete="new-password"
          />
        </a-form-item>
        <a-form-item label="确认新密码" name="confirm_password">
          <a-input-password
            v-model:value="form.confirm_password"
            placeholder="请再次输入新密码"
            autocomplete="new-password"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit" :loading="submitting">
              提交
            </a-button>
            <a-button @click="handleBack">返回</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message, type FormInstance } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

// 确认密码校验
const validateConfirm = async (_rule: Rule, value: string): Promise<void> => {
  if (!value) {
    return Promise.reject('请再次输入新密码')
  }
  if (value !== form.new_password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const rules: Record<string, Rule[]> = {
  old_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能少于 8 位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, validator: validateConfirm, trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  submitting.value = true
  const ok = await authStore.changePassword({
    old_password: form.old_password,
    new_password: form.new_password,
  })
  submitting.value = false
  if (ok) {
    message.success('密码修改成功')
    // 清空表单
    form.old_password = ''
    form.new_password = ''
    form.confirm_password = ''
    // 返回上一页
    router.back()
  }
}

const handleBack = () => {
  router.back()
}
</script>

<style scoped>
.change-password-page {
  max-width: 560px;
  margin: 0 auto;
}

.pwd-card {
  border-radius: 8px;
}
</style>
