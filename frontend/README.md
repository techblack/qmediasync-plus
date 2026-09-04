# QMediaSync Web

基于 Vue 3、TypeScript、Ant Design Vue、Vue Router 和 Pinia 的管理前端。生产构建会直接输出到项目根目录的 `web_statics/`，由 Go 服务托管。

## 模块结构

```text
src/
  main.ts              # Vue 应用入口
  App.vue              # Ant Design 全局主题
  router/              # 路由与登录守卫
  stores/              # Pinia 状态仓库
  services/            # 类型安全的 API 客户端
  types/               # 后端接口类型模型
  layouts/             # 管理后台整体布局
  views/               # 按业务领域拆分的 Vue 单文件组件
  styles/              # 全局设计令牌与响应式样式
```

```bash
cd frontend
npm install
npm run dev
```

本地开发默认将 `/api` 代理到 `http://127.0.0.1:12333`，也可以临时指定其他服务：

```bash
VITE_API_PROXY=http://127.0.0.1:8115 npm run dev
```

生成生产资源：

```bash
npm run build
```
