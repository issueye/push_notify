<script setup>
import { ref, onMounted, computed } from 'vue'
import { NCard, NStatistic, NGrid, NGridItem, NSpin } from 'naive-ui'
import { getPushStats } from '@/services/push'
import { useLoading } from '@/composables/useMessage'

// ECharts 引入
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import VChart from 'vue-echarts'

use([CanvasRenderer, PieChart, LineChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

const { loading, start, stop } = useLoading()
const stats = ref({
  today: { total: 0, success: 0, failed: 0 },
  this_week: { total: 0, success: 0, failed: 0 },
  this_month: { total: 0, success: 0, failed: 0 },
  trend: []
})

// 饼图配置 (成功率)
const pieOption = computed(() => {
  const success = stats.value.this_month.success || 0
  const failed = stats.value.this_month.failed || 0
  
  return {
    title: {
      text: '本月推送成功率',
      left: 'center'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      bottom: '5%',
      left: 'center'
    },
    series: [
      {
        name: '推送状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        data: [
          { value: success, name: '成功', itemStyle: { color: '#10b981' } },
          { value: failed, name: '失败', itemStyle: { color: '#ef4444' } }
        ]
      }
    ]
  }
})

// 折线图配置 (趋势)
const lineOption = computed(() => {
  const trend = stats.value.trend || []
  return {
    title: {
      text: '近 7 天推送趋势',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['总计', '成功', '失败'],
      bottom: '0%'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: trend.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '总计',
        type: 'line',
        smooth: true,
        data: trend.map(item => item.total),
        itemStyle: { color: '#6366f1' }
      },
      {
        name: '成功',
        type: 'line',
        smooth: true,
        data: trend.map(item => item.success),
        itemStyle: { color: '#10b981' }
      },
      {
        name: '失败',
        type: 'line',
        smooth: true,
        data: trend.map(item => item.failed),
        itemStyle: { color: '#ef4444' }
      }
    ]
  }
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
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">工作台</h1>
      <n-spin v-if="loading" size="small" />
    </div>

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

    <n-grid :cols="2" :x-gap="16" :y-gap="16" class="mb-6">
      <n-grid-item>
        <n-card class="h-[400px]">
          <v-chart class="h-full" :option="pieOption" autoresize />
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card class="h-[400px]">
          <v-chart class="h-full" :option="lineOption" autoresize />
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-card title="快速开始">
      <div class="grid grid-cols-3 gap-4">
        <router-link to="/repos" class="p-4 border rounded-lg hover:bg-gray-50 text-center transition-colors">
          <div class="text-3xl mb-2">📦</div>
          <div class="font-medium">添加仓库</div>
          <div class="text-sm text-gray-500">配置代码仓库</div>
        </router-link>
        <router-link to="/targets" class="p-4 border rounded-lg hover:bg-gray-50 text-center transition-colors">
          <div class="text-3xl mb-2">🔔</div>
          <div class="font-medium">配置推送</div>
          <div class="text-sm text-gray-500">设置钉钉/邮箱</div>
        </router-link>
        <router-link to="/templates" class="p-4 border rounded-lg hover:bg-gray-50 text-center transition-colors">
          <div class="text-3xl mb-2">📝</div>
          <div class="font-medium">消息模板</div>
          <div class="text-sm text-gray-500">自定义消息格式</div>
        </router-link>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.n-card {
  @apply shadow-sm hover:shadow-md transition-shadow;
}
</style>
