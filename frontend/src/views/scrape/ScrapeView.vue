<template>
  <div class="toolbar">
    <a-segmented v-model:value="tab" :options="[{label:'刮削目录',value:'paths'},{label:'处理记录',value:'records'}]"/>
    <a-button v-if="tab==='paths'" type="primary" @click="openCreate"><PlusOutlined/>新增刮削目录</a-button>
    <div class="grow"/>
    <template v-if="tab==='records'"><a-input-search v-model:value="keyword" placeholder="搜索文件" style="width:220px" @search="load"/><a-button @click="retryFailed">重试失败</a-button><a-popconfirm title="确定清理失败记录？" @confirm="clearFailed"><a-button danger>清理失败</a-button></a-popconfirm></template>
  </div>
  <a-card class="panel-card">
    <a-table v-if="tab==='paths'" :loading="loading" :data-source="paths" :columns="pathColumns" row-key="id" :scroll="{x:900}">
      <template #bodyCell="{column,record}">
        <template v-if="column.key==='path'"><strong>{{record.source_path}}</strong><div class="muted">目标：{{record.dest_path}}</div></template>
        <template v-if="column.key==='cron'"><a-switch :checked="record.enable_cron" @change="toggleCron(record)"/> <span class="muted">{{record.cron_expression}}</span></template>
        <template v-if="column.key==='status'"><a-badge :status="record.is_running?'processing':'success'" :text="record.is_running===2?'运行中':record.is_running===1?'等待中':'空闲'"/></template>
        <template v-if="column.key==='action'"><a-space><a-button size="small" @click="start(record)"><PlayCircleOutlined/></a-button><a-button size="small" @click="stop(record)"><PauseCircleOutlined/></a-button><a-button size="small" @click="edit(record)"><EditOutlined/></a-button><a-popconfirm title="确定删除？" @confirm="remove(record)"><a-button size="small" danger><DeleteOutlined/></a-button></a-popconfirm></a-space></template>
      </template>
    </a-table>
    <a-table v-else :loading="loading" :data-source="records" :columns="recordColumns" row-key="id">
      <template #bodyCell="{column,record}"><template v-if="column.key==='file'"><strong>{{record.file_name}}</strong><div class="muted">{{record.path}}</div></template><template v-if="column.key==='status'"><a-tag>{{record.status}}</a-tag></template><template v-if="column.key==='action'"><a-button size="small" @click="rescrape(record)"><RedoOutlined/>重新刮削</a-button></template></template>
    </a-table>
  </a-card>
  <a-modal v-model:open="visible" :title="form.id?'编辑刮削目录':'新增刮削目录'" width="820px" @ok="save">
    <a-form layout="vertical">
      <a-row :gutter="16"><a-col :span="12"><a-form-item label="来源类型"><a-select v-model:value="form.source_type" :disabled="!!form.id" :options="sources"/></a-form-item></a-col><a-col :span="12"><a-form-item label="网盘账号"><a-select v-model:value="form.account_id" :disabled="form.source_type==='local'" :options="accountOptions"/></a-form-item></a-col></a-row>
      <a-form-item label="来源目录"><a-input v-model:value="form.source_path"/></a-form-item><a-form-item v-if="form.source_type!=='local'" label="来源目录 ID"><a-input v-model:value="form.source_path_id"/></a-form-item>
      <a-form-item label="目标目录"><a-input v-model:value="form.dest_path"/></a-form-item><a-form-item v-if="form.source_type!=='local'" label="目标目录 ID"><a-input v-model:value="form.dest_path_id"/></a-form-item>
      <a-row :gutter="16"><a-col :span="12"><a-form-item label="媒体类型"><a-select v-model:value="form.media_type" :options="mediaTypes"/></a-form-item></a-col><a-col :span="12"><a-form-item label="处理方式"><a-select v-model:value="form.scrape_type" :options="scrapeTypes"/></a-form-item></a-col><a-col :span="12"><a-form-item label="整理方式"><a-select v-model:value="form.rename_type" :options="renameTypes"/></a-form-item></a-col><a-col :span="12"><a-form-item label="最大线程"><a-input-number v-model:value="form.max_threads" :min="1" style="width:100%"/></a-form-item></a-col></a-row>
      <a-form-item label="文件夹模板"><a-input v-model:value="form.folder_name_template"/></a-form-item><a-form-item label="文件模板"><a-input v-model:value="form.file_name_template"/></a-form-item>
      <a-space wrap><a-switch v-model:checked="form.enable_cron"/>定时刮削<a-switch v-model:checked="form.enable_category"/>启用分类<a-switch v-model:checked="form.enable_fanart_tv"/>Fanart<a-switch v-model:checked="form.exclude_no_image_actor"/>排除无图片演员<a-switch v-model:checked="form.force_delete_source_path"/>强制清理源目录</a-space>
      <a-form-item v-if="form.enable_cron" label="Cron" style="margin-top:16px"><a-input v-model:value="form.cron_expression"/></a-form-item>
    </a-form>
  </a-modal>
</template>
<script setup lang="ts">
import{computed,onMounted,reactive,ref,watch}from'vue';import{message}from'ant-design-vue';import{PlusOutlined,PlayCircleOutlined,PauseCircleOutlined,DeleteOutlined,RedoOutlined,EditOutlined}from'@ant-design/icons-vue';import{del,get,post}from'../../services/http';import type{Account,AnyMap,ScrapePath}from'../../types';
const tab=ref('paths'),keyword=ref(''),loading=ref(true),visible=ref(false),paths=ref<ScrapePath[]>([]),records=ref<AnyMap[]>([]),accounts=ref<Account[]>([]);
const defaults={id:0,source_type:'local',account_id:0,media_type:'movie',source_path:'',source_path_id:'',dest_path:'',dest_path_id:'',scrape_type:'scrape_and_rename',rename_type:'move',max_threads:5,enable_cron:false,cron_expression:'',enable_category:false,enable_fanart_tv:false,exclude_no_image_actor:false,force_delete_source_path:false,folder_name_template:'{{title}} ({{year}})',file_name_template:'{{title}} ({{year}})',enable_ai:'off'};const form=reactive({...defaults});
const pathColumns=[{title:'目录',key:'path'},{title:'媒体',dataIndex:'media_type'},{title:'方式',dataIndex:'scrape_type'},{title:'定时',key:'cron'},{title:'状态',key:'status'},{title:'操作',key:'action'}],recordColumns=[{title:'文件',key:'file'},{title:'类型',dataIndex:'type'},{title:'状态',key:'status'},{title:'操作',key:'action'}];
const sources=[['local','本地目录'],['115','115 网盘'],['baidupan','百度网盘'],['openlist','OpenList']].map(([value,label])=>({value,label})),mediaTypes=[{value:'movie',label:'电影'},{value:'tvshow',label:'电视剧'}],scrapeTypes=[{value:'scrape_and_rename',label:'刮削并整理'},{value:'only_scrape',label:'仅刮削'},{value:'only_rename',label:'仅整理'}],renameTypes=[['move','移动'],['copy','复制'],['hard_symlink','硬链接'],['soft_symlink','软链接']].map(([value,label])=>({value,label}));
const accountOptions=computed(()=>accounts.value.filter(a=>a.source_type===form.source_type).map(a=>({value:a.id,label:a.name||a.username})));
async function load(){loading.value=true;try{const[p,r,a]=await Promise.all([get<ScrapePath[]>('/scrape/pathes'),get<AnyMap>(`/scrape/records?page=1&pageSize=100&name=${encodeURIComponent(keyword.value)}`),get<Account[]>('/account/list')]);paths.value=p;records.value=r.list||[];accounts.value=a}finally{loading.value=false}}
function openCreate(){Object.assign(form,defaults);visible.value=true}async function edit(row:ScrapePath){Object.assign(form,await get<AnyMap>(`/scrape/pathes/${row.id}`));visible.value=true}async function save(){await post('/scrape/pathes',form);visible.value=false;message.success('刮削目录已保存');load()}
const start=(r:ScrapePath)=>post('/scrape/pathes/start',{id:r.id}),stop=(r:ScrapePath)=>post('/scrape/pathes/stop',{id:r.id});async function toggleCron(r:ScrapePath){await post('/scrape/pathes/toggle-cron',{id:r.id});load()}async function remove(r:ScrapePath){await del(`/scrape/pathes/${r.id}`);load()}const rescrape=(r:AnyMap)=>post('/scrape/re-scrape',{id:r.id,tmdb_id:0,season:0,episode:0}),retryFailed=()=>post('/scrape/rename-failed',{}),clearFailed=()=>post('/scrape/clear-failed',{});watch(tab,load);onMounted(load)
</script>
