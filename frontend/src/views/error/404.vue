<script setup>
import { h } from 'vue'
import { NCard } from 'naive-ui'
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <n-card class="w-96 text-center">
      <div class="text-6xl mb-4">😕</div>
      <h1 class="text-2xl font-bold mb-2">页面不存在</h1>
      <p class="text-gray-500 mb-4">您访问的页面不存在或已被删除</p>
      <router-link to="/dashboard">
        <n-button type="primary">返回首页</n-button>
      </router-link>
    </n-card>
  </div>
</template>
