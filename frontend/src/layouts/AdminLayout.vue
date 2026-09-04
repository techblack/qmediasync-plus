<template>
  <a-layout class="shell">
    <a-layout-sider v-model:collapsed="collapsed" collapsible :width="244" class="sider">
      <div class="brand"><span class="logo">Q</span><strong v-if="!collapsed">QMediaSync</strong></div>
      <a-menu theme="dark" mode="inline" :selectedKeys="[route.path]" @click="navigate">
        <a-menu-item-group v-for="group in menus" :key="group.title" :title="group.title">
          <a-menu-item v-for="item in group.items" :key="item.path"><component :is="item.icon"/><span>{{item.label}}</span></a-menu-item>
        </a-menu-item-group>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="topbar"><a-button class="mobile-trigger" type="text" @click="collapsed=!collapsed"><MenuOutlined/></a-button><div class="grow"/><a-dropdown><a-button type="text"><a-avatar size="small">{{session.user?.username?.[0]?.toUpperCase()}}</a-avatar><span class="username">{{session.user?.username}}</span><DownOutlined/></a-button><template #overlay><a-menu><a-menu-item @click="router.push('/system/user')">账号设置</a-menu-item><a-menu-divider/><a-menu-item danger @click="signOut">退出登录</a-menu-item></a-menu></template></a-dropdown></a-layout-header>
      <a-layout-content class="content"><div class="page-head"><div><h1>{{route.meta.title||sectionTitle}}</h1><p>{{route.meta.subtitle||'QMediaSync 系统管理'}}</p></div><a-button @click="reload"><ReloadOutlined/>刷新</a-button></div><router-view :key="route.fullPath"/></a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'; import { useRoute,useRouter } from 'vue-router'; import { AppstoreOutlined,CloudOutlined,SwapOutlined,DatabaseOutlined,SettingOutlined,FileSearchOutlined,CloudUploadOutlined,FolderOpenOutlined,KeyOutlined,BellOutlined,ApiOutlined,ToolOutlined,FileTextOutlined,QuestionCircleOutlined,MenuOutlined,DownOutlined,ReloadOutlined } from '@ant-design/icons-vue'; import { useSessionStore } from '../stores/session'
const route=useRoute(),router=useRouter(),session=useSessionStore(),collapsed=ref(window.innerWidth<760)
const menus=[{title:'概览',items:[['/dashboard','控制台',AppstoreOutlined],['/accounts','网盘账号',CloudOutlined]]},{title:'STRM',items:[['/strm/paths','同步目录',SwapOutlined],['/strm/records','同步记录',FileTextOutlined],['/strm/settings','STRM 设置',SettingOutlined]]},{title:'刮削整理',items:[['/scrape/paths','目录与记录',FileSearchOutlined],['/scrape/settings','刮削设置',SettingOutlined]]},{title:'文件任务',items:[['/tasks/queue','上传下载',CloudUploadOutlined],['/tasks/files','网盘文件',FolderOpenOutlined]]},{title:'备份恢复',items:[['/backup/records','备份记录',DatabaseOutlined],['/backup/settings','备份设置',SettingOutlined],['/backup/restore','恢复备份',DatabaseOutlined]]},{title:'系统',items:[['/system/api-keys','API Key',KeyOutlined],['/system/notifications','通知管理',BellOutlined],['/system/emby','Emby',ApiOutlined],['/system/rate','接口速率',ApiOutlined],['/system/network','网络与线程',SettingOutlined],['/system/update','版本更新',CloudUploadOutlined],['/system/maintenance','维护工具',ToolOutlined],['/system/logs','运行日志',FileTextOutlined],['/system/user','用户设置',SettingOutlined],['/system/help','使用帮助',QuestionCircleOutlined]]}].map(g=>({...g,items:g.items.map(i=>({path:i[0] as string,label:i[1] as string,icon:i[2]}))}))
const sectionTitle=computed(()=>menus.flatMap(g=>g.items).find(i=>i.path===route.path)?.label||'系统管理'); const navigate=({key}:{key:string})=>router.push(key); const reload=()=>router.replace({...route,query:{...route.query,_:Date.now()}}); async function signOut(){await session.logout();router.push('/login')}
</script>
<style scoped>.shell{min-height:100vh}.sider{position:sticky!important;top:0;height:100vh;overflow:auto}.brand{height:70px;display:flex;align-items:center;gap:11px;padding:0 20px;color:#fff;font-size:17px}.logo{display:grid;place-items:center;width:34px;height:34px;border-radius:10px;background:linear-gradient(135deg,#6e7cff,#8f62df);font-size:20px}.topbar{height:64px;padding:0 26px;background:#fff;display:flex;align-items:center;border-bottom:1px solid #edf0f5}.content{padding:28px 34px 44px;max-width:1600px;width:100%;margin:auto}.mobile-trigger{display:none}.username{margin:0 8px}@media(max-width:760px){.sider{position:fixed!important;z-index:20}.sider.ant-layout-sider-collapsed{transform:translateX(-80px)}.mobile-trigger{display:inline-flex}.content{padding:22px 16px}.username{display:none}}</style>
