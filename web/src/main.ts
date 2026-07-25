import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import Antd from 'ant-design-vue'
import App from './App.vue'
import router from './router'
import './assets/styles/index.css'
import 'ant-design-vue/dist/reset.css'

// 创建Vue应用实例
const app = createApp(App)

// 创建Pinia实例并注册持久化插件
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// 注册Pinia
app.use(pinia)

// 注册路由
app.use(router)

// 注册Ant Design Vue
app.use(Antd)

// 注册全局错误处理器
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误:', err)
  console.error('错误信息:', info)
}

// 挂载应用
app.mount('#app')
