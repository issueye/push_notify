<script setup>
import { ref, onMounted } from 'vue'
import { NCard, NStatistic, NGrid, NGridItem, NSpin } from 'naive-ui'
import { getPushStats } from '@/services/push'
import { useLoading } from '@/composables/useMessage'

const { loading, start, stop } = useLoading()
const stats = ref({
  today: { total: 0, success: 0, failed: 0 },
  this_week: { total: 0, success: 0, failed: 0 },
  this_month: { total: 0, success: 0, failed: 0 }
})

async function fetchStats() {
  start()
  try {
    const data = await getPushStats()
    stats.value = data
  } catch (e) {
    console.error(e)
  } finally {
    stop()
  }
}

onMounted(fetchStats)
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">工作台</h1>

    <n-grid :cols="3" :x-gap="16" :y-gap="16" class="mb-6">
      <n-grid-item>
        <n-card>
          <n-statistic label="今日推送">
            <template #prefix>
              <span class="text-primary-600">{{ stats.today.total || 0 }}</span>
            </template>
            <template #suffix>
              <span class="text-green-600 text-sm">成功 {{ stats.today.success || 0 }}</span>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card>
          <n-statistic label="本周推送">
            <template #prefix>
              <span class="text-primary-600">{{ stats.this_week.total || 0 }}</span>
            </template>
            <template #suffix>
              <span class="text-green-600 text-sm">成功 {{ stats.this_week.success || 0 }}</span>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card>
          <n-statistic label="本月推送">
            <template #prefix>
              <span class="text-primary-600">{{ stats.this_month.total || 0 }}</span>
            </template>
            <template #suffix>
              <span class="text-green-600 text-sm">成功 {{ stats.this_month.success || 0 }}</span>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-card title="快速开始">
      <div class="grid grid-cols-3 gap-4">
        <router-link to="/repos" class="p-4 border rounded-lg hover:bg-gray-50 text-center">
          <div class="text-3xl mb-2">📦</div>
          <div class="font-medium">添加仓库</div>
          <div class="text-sm text-gray-500">配置代码仓库</div>
        </router-link>
        <router-link to="/targets" class="p-4 border rounded-lg hover:bg-gray-50 text-center">
          <div class="text-3xl mb-2">🔔</div>
          <div class="font-medium">配置推送</div>
          <div class="text-sm text-gray-500">设置钉钉/邮箱</div>
        </router-link>
        <router-link to="/templates" class="p-4 border rounded-lg hover:bg-gray-50 text-center">
          <div class="text-3xl mb-2">📝</div>
          <div class="font-medium">消息模板</div>
          <div class="text-sm text-gray-500">自定义消息格式</div>
        </router-link>
      </div>
    </n-card>
  </div>
</template>
