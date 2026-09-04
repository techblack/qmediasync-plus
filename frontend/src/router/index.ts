import { createRouter, createWebHashHistory } from 'vue-router'
import AdminLayout from '../layouts/AdminLayout.vue'
import { useSessionStore } from '../stores/session'

const routes = [
  { path:'/login', component:()=>import('../views/LoginView.vue'), meta:{public:true} },
  { path:'/', component:AdminLayout, children:[
    { path:'', redirect:'/dashboard' },
    { path:'dashboard', component:()=>import('../views/DashboardView.vue'), meta:{title:'控制台',subtitle:'系统运行状态与任务概览'} },
    { path:'accounts', component:()=>import('../views/accounts/AccountsView.vue'), meta:{title:'网盘账号',subtitle:'管理云盘授权与 Cookie 驱动'} },
    { path:'strm/paths', component:()=>import('../views/strm/PathsView.vue'), meta:{title:'同步目录',subtitle:'管理云盘与本地 STRM 同步配置'} },
    { path:'strm/records', component:()=>import('../views/strm/RecordsView.vue'), meta:{title:'同步记录',subtitle:'查看同步任务执行结果'} },
    { path:'strm/settings', component:()=>import('../views/strm/SettingsView.vue'), meta:{title:'STRM 设置',subtitle:'配置生成规则与元数据策略'} },
    { path:'scrape/paths', component:()=>import('../views/scrape/ScrapeView.vue'), meta:{title:'刮削与整理',subtitle:'管理刮削目录和处理记录'} },
    { path:'scrape/settings', component:()=>import('../views/scrape/SettingsView.vue'), meta:{title:'刮削设置',subtitle:'TMDB、AI 与媒体分类策略'} },
    { path:'tasks/queue', component:()=>import('../views/tasks/QueueView.vue'), meta:{title:'上传下载',subtitle:'管理文件传输队列'} },
    { path:'tasks/files', component:()=>import('../views/tasks/FilesView.vue'), meta:{title:'网盘文件',subtitle:'浏览和管理云盘文件'} },
    { path:'backup/records', component:()=>import('../views/backup/BackupView.vue'), meta:{title:'备份记录',subtitle:'数据库备份历史'} },
    { path:'backup/settings', component:()=>import('../views/backup/SettingsView.vue'), meta:{title:'备份设置',subtitle:'自动备份策略'} },
    { path:'backup/restore', component:()=>import('../views/backup/RestoreView.vue'), meta:{title:'恢复备份',subtitle:'从历史记录或上传文件恢复'} },
    { path:'system/:section', component:()=>import('../views/system/SystemView.vue'), props:true },
  ]},
]
const router = createRouter({ history:createWebHashHistory(), routes })
router.beforeEach(async(to)=>{ const session=useSessionStore(); if(!session.initialized) await session.restore(); if(!to.meta.public&&!session.user)return'/login'; if(to.path==='/login'&&session.user)return'/dashboard' })
export default router
