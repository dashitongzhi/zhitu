<template>
  <a-modal
    :open="open"
    :width="640"
    title="个人资料"
    :mask-closable="false"
    wrap-class-name="user-profile-modal"
    @update:open="(v: boolean) => emit('update:open', v)"
    @ok="handleSave"
    @cancel="handleCancel"
    ok-text="保存并同步"
    cancel-text="取消"
  >
    <div class="up-intro">
      <UserCircle :size="14" :stroke-width="2" />
      <span>填写的基础信息会自动同步到简历实验室，无需重复输入</span>
    </div>

    <div class="up-avatar-row">
      <div class="up-avatar">
        <img v-if="form.avatar" :src="form.avatar" alt="头像" />
        <span v-else class="up-avatar-placeholder">{{ form.name?.charAt(0).toUpperCase() || '?' }}</span>
      </div>
      <div class="up-avatar-meta">
        <div class="up-avatar-name">{{ form.name || '尚未填写姓名' }}</div>
        <div class="up-avatar-tip">建议上传正方形头像，将显示在简历与导航栏</div>
      </div>
    </div>

    <a-form layout="vertical" :model="form" class="up-form">
      <div class="up-row-2">
        <a-form-item label="姓名" required>
          <a-input v-model:value="form.name" placeholder="请输入姓名" :maxlength="20" />
        </a-form-item>
        <a-form-item label="性别">
          <a-select v-model:value="form.gender" placeholder="请选择">
            <a-select-option value="">未设置</a-select-option>
            <a-select-option value="male">男</a-select-option>
            <a-select-option value="female">女</a-select-option>
          </a-select>
        </a-form-item>
      </div>

      <div class="up-row-2">
        <a-form-item label="手机号">
          <a-input v-model:value="form.phone" placeholder="138-0000-0000" :maxlength="20" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model:value="form.email" placeholder="example@email.com" :maxlength="60" />
        </a-form-item>
      </div>

      <div class="up-row-2">
        <a-form-item label="所在城市">
          <a-input v-model:value="form.city" placeholder="如：上海市浦东新区" :maxlength="50" />
        </a-form-item>
        <a-form-item label="出生日期">
          <a-date-picker
            v-model:value="birthDateObj"
            style="width: 100%"
            placeholder="请选择"
            value-format="YYYY-MM-DD"
            :disabled-date="(d: any) => d && d.isAfter(new Date())"
          />
        </a-form-item>
      </div>

      <a-form-item label="GitHub / 个人主页">
        <a-input v-model:value="form.github" placeholder="github.com/yourname" :maxlength="100" />
      </a-form-item>

      <a-form-item label="自我介绍">
        <a-textarea
          v-model:value="form.self_introduction"
          :rows="4"
          placeholder="一句话介绍自己，会同步到简历的个人简介模块"
          :maxlength="300"
          show-count
        />
      </a-form-item>
    </a-form>

    <div class="up-footer-tip">
      <CheckCircle :size="12" :stroke-width="2" />
      <span>保存后，简历实验室的「个人信息」字段会自动更新</span>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import dayjs, { type Dayjs } from 'dayjs'
import { message } from 'ant-design-vue'
import { UserCircle, CheckCircle } from 'lucide-vue-next'
import { useUserProfileStore, type UserProfileBasic } from '@/stores/userProfile'

interface Props {
  open: boolean
}
const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'saved', data: UserProfileBasic): void
}>()

const userProfileStore = useUserProfileStore()

const form = reactive<UserProfileBasic>({
  name: '',
  phone: '',
  email: '',
  city: '',
  github: '',
  gender: '',
  birth_date: '',
  self_introduction: '',
  avatar: '',
})

// 出生日 Dayjs 对象
const birthDateObj = computed<Dayjs | ''>({
  get: () => (form.birth_date ? dayjs(form.birth_date) : ''),
  set: (v) => {
    form.birth_date = v ? (v as Dayjs).format('YYYY-MM-DD') : ''
  },
})

// 弹窗打开时同步 store 数据
watch(
  () => props.open,
  (v) => {
    if (v) {
      Object.assign(form, userProfileStore.basic)
    }
  },
  { immediate: true }
)

const handleSave = () => {
  if (!form.name?.trim()) {
    message.warning('请填写姓名')
    return
  }
  userProfileStore.updateBasic({ ...form })
  emit('saved', { ...form })
  emit('update:open', false)
  message.success('已保存，并同步到简历实验室')
}

const handleCancel = () => {
  emit('update:open', false)
}
</script>

<style scoped>
.up-intro {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: linear-gradient(135deg, rgba(0, 122, 255, 0.06), rgba(0, 122, 255, 0.02));
  border: 1px solid rgba(0, 122, 255, 0.15);
  border-radius: 8px;
  font-size: 13px;
  color: #007aff;
  margin-bottom: 20px;
}

.up-avatar-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: #f7f8fa;
  border-radius: 10px;
  margin-bottom: 20px;
}
.up-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #007aff, #5e5ce6);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  overflow: hidden;
  flex-shrink: 0;
}
.up-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.up-avatar-placeholder {
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif;
}
.up-avatar-name {
  font-size: 16px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 4px;
}
.up-avatar-tip {
  font-size: 12px;
  color: #86868b;
}

.up-form :deep(.ant-form-item) {
  margin-bottom: 16px;
}
.up-form :deep(.ant-form-item-label > label) {
  font-weight: 600;
  color: #1a1a1a;
  font-size: 13px;
}
.up-form :deep(.ant-input),
.up-form :deep(.ant-select-selector),
.up-form :deep(.ant-picker) {
  border-radius: 8px;
}

.up-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
@media (max-width: 520px) {
  .up-row-2 {
    grid-template-columns: 1fr;
  }
}

.up-footer-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: rgba(52, 199, 89, 0.08);
  border-radius: 8px;
  font-size: 12px;
  color: #248a3d;
  margin-top: 12px;
}
</style>
