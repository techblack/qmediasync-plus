import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { message } from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import './styles/theme.css'
import App from './App.vue'
import router from './router'
import { handleOAuthCallback } from './services/oauth'

const app = createApp(App)
app.config.errorHandler = (error) => message.error(error instanceof Error ? error.message : '页面操作失败')
window.addEventListener('unhandledrejection', (event) => { event.preventDefault(); message.error(event.reason instanceof Error ? event.reason.message : '请求失败') })
app.use(createPinia()).use(router).mount('#app')
handleOAuthCallback().then((handled) => { if (handled) router.replace('/accounts') }).catch((error) => message.error(error instanceof Error ? error.message : '授权确认失败'))
