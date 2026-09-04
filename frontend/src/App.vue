<template><a-config-provider :theme="theme"><router-view /></a-config-provider></template>
<script setup lang="ts">
const theme = { token: { colorPrimary: '#5b6cf9', borderRadius: 10, colorBgLayout: '#f4f6fb', fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "PingFang SC", sans-serif' } }
</script>
