<template>
  <a-config-provider
    :locale="zhCN"
    :theme="themeConfig"
  >
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { theme } from 'ant-design-vue'

dayjs.locale('zh-cn')

// Pinguo 设计系统主题配置
const themeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: '#007aff',
    borderRadius: 10,
    fontFamily:
      "-apple-system, BlinkMacSystemFont, 'SF Pro Text', 'SF Pro Display', 'PingFang SC', 'DM Sans', ui-sans-serif, system-ui, sans-serif",
  },
}
</script>

<style>
#app {
  min-height: 100vh;
}
</style>
