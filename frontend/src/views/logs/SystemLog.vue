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
  NSelect,
  NDatePicker,
  NPagination,
  useMessage,
} from "naive-ui";
import { SearchOutline, DownloadOutline } from "@vicons/ionicons5";
import { getSystemLogs } from "@/services/log";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const logs = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchLevel = ref(null);
const searchTimeRange = ref(null);

const levelOptions = [
  { label: "全部级别", value: null },
  { label: "DEBUG", value: "debug" },
  { label: "INFO", value: "info" },
  { label: "WARN", value: "warn" },
  { label: "ERROR", value: "error" },
];

const columns = [
  { title: "ID", key: "id", width: 60 },
  {
    title: "级别",
    key: "level",
    width: 80,
    render(row) {
      const colors = {
        debug: "default",
        info: "info",
        warn: "warning",
        error: "error",
      };
      return h(
        NTag,
        { type: colors[row.level] || "default", size: "small" },
        () => row.level?.toUpperCase(),
      );
    },
  },
  { title: "模块", key: "module", width: 100 },
  { title: "消息", key: "message", ellipsis: { tooltip: true } },
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
    level: searchLevel.value,
  };
  if (searchTimeRange.value) {
    params.start_time = new Date(searchTimeRange.value[0]).toISOString();
    params.end_time = new Date(searchTimeRange.value[1]).toISOString();
  }
  try {
    const data = await getSystemLogs(params);
    logs.value = data.list;
    total.value = data.pagination.total;
  } catch (e) {
    message.error("获取日志失败");
  } finally {
    loading.value = false;
  }
}

function handleExport() {
  message.info("导出功能");
}

watch([page, searchKeyword, searchLevel, searchTimeRange], fetchLogs);
onMounted(fetchLogs);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">系统日志</h1>
      <n-button @click="handleExport">
        <template #icon
          ><n-icon><DownloadOutline /></n-icon
        ></template>
        导出
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4 flex-wrap">
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
        <n-select
          v-model:value="searchLevel"
          :options="levelOptions"
          style="width: 150px"
        />
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
      <div class="mt-4 flex justify-end">
        <n-pagination
          v-model:page="page"
          v-model:page-size="size"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50, 100]"
          @update:page="fetchLogs"
          @update:page-size="fetchLogs"
        />
      </div>
    </n-card>
  </div>
</template>
