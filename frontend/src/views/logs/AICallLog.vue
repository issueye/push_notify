<script setup>
import { ref, onMounted, watch, h } from "vue";
import { formatDate } from "@/utils/date";
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NInput,
  NDatePicker,
  useMessage,
} from "naive-ui";
import { SearchOutline, DownloadOutline } from "@vicons/ionicons5";
import { getAICallLogs } from "@/services/log";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const logs = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchTimeRange = ref(null);

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "模型ID", key: "model_id", width: 80 },
  {
    title: "状态",
    key: "level",
    width: 80,
    render(row) {
      return h(
        NTag,
        { type: row.level === "info" ? "success" : "error", size: "small" },
        () => (row.level === "info" ? "成功" : "失败"),
      );
    },
  },
  { title: "详情", key: "message", ellipsis: { tooltip: true } },
  {
    title: "时间",
    key: "created_at",
    width: 170,
    render(row) {
      return formatDate(row.created_at);
    },
  },
];

async function fetchLogs() {
  loading.value = true;
  const params = {
    page: page.value,
    size: size.value,
    keyword: searchKeyword.value,
  };
  if (searchTimeRange.value) {
    params.start_time = new Date(searchTimeRange.value[0]).toISOString();
    params.end_time = new Date(searchTimeRange.value[1]).toISOString();
  }
  try {
    const data = await getAICallLogs(params);
    logs.value = data.list;
    total.value = data.pagination.total;
  } catch (e) {
    message.error("获取日志失败");
  } finally {
    loading.value = false;
  }
}

watch([page, searchKeyword, searchTimeRange], fetchLogs);
onMounted(fetchLogs);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">AI调用日志</h1>
      <n-button @click="message.info('导出功能')">
        <template #icon
          ><n-icon><DownloadOutline /></n-icon
        ></template>
        导出
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索日志内容"
          clearable
          style="width: 250px"
        >
          <template #prefix
            ><n-icon><SearchOutline /></n-icon
          ></template>
        </n-input>
        <n-date-picker
          v-model:value="searchTimeRange"
          type="daterange"
          clearable
          style="width: 300px"
        />
        <n-button @click="fetchLogs">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="false"
        :bordered="true"
      />
    </n-card>
  </div>
</template>
