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
  useMessage,
} from "naive-ui";
import { RefreshOutline, SearchOutline } from "@vicons/ionicons5";
import { getPushList, retryPush, batchRetry } from "@/services/push";
import { usePagination, useConfirm } from "@/composables/useMessage";

const message = useMessage();
const { confirm } = useConfirm();
const { page, size, total } = usePagination();

const pushes = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchStatus = ref(null);
const searchTimeRange = ref(null);
const selectedRowKeys = ref([]);

const statusOptions = [
  { label: "全部状态", value: null },
  { label: "成功", value: "success" },
  { label: "失败", value: "failed" },
  { label: "待推送", value: "pending" },
];

const columns = [
  { type: "selection", width: 50 },
  { title: "ID", key: "id", width: 60 },
  { title: "仓库", key: "repo_name", width: 120 },
  { title: "目标", key: "target_name", width: 100 },
  { title: "提交信息", key: "commit_msg", ellipsis: { tooltip: true } },
  {
    title: "状态",
    key: "status",
    width: 90,
    render(row) {
      const type =
        { success: "success", failed: "error", pending: "warning" }[
          row.status
        ] || "default";
      const text =
        { success: "成功", failed: "失败", pending: "待推送" }[row.status] ||
        row.status;
      return h(NTag, { type, size: "small" }, () => text);
    },
  },
  {
    title: "推送时间",
    key: "pushed_at",
    width: 170,
    render(row) {
      return formatDate(row.pushed_at);
    },
  },
  {
    title: "创建时间",
    key: "created_at",
    width: 170,
    render(row) {
      return formatDate(row.created_at);
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 100,
    render(row) {
      return h(
        NButton,
        { size: "small", quaternary: true, onClick: () => handleRetry(row.id) },
        () => "重试",
      );
    },
  },
];

async function fetchPushes() {
  loading.value = true;
  const params = {
    page: page.value,
    size: size.value,
    keyword: searchKeyword.value,
    status: searchStatus.value,
  };
  if (searchTimeRange.value) {
    params.start_time = new Date(searchTimeRange.value[0]).toISOString();
    params.end_time = new Date(searchTimeRange.value[1]).toISOString();
  }
  try {
    const data = await getPushList(params);
    pushes.value = data.list;
    total.value = data.pagination.total;
  } catch (e) {
    message.error("获取推送记录失败");
  } finally {
    loading.value = false;
  }
}

async function handleRetry(id) {
  try {
    await retryPush(id);
    message.success("重试已提交");
    fetchPushes();
  } catch (e) {
    message.error("重试失败");
  }
}

function handleBatchRetry() {
  if (selectedRowKeys.value.length === 0) {
    message.warning("请选择要重试的记录");
    return;
  }
  confirm(`确定要重试选中的 ${selectedRowKeys.value.length} 条记录吗？`, () => {
    batchRetry(selectedRowKeys.value).then(() => {
      message.success("批量重试已提交");
      selectedRowKeys.value = [];
      fetchPushes();
    });
  });
}

watch([page, searchKeyword, searchStatus, searchTimeRange], fetchPushes);
onMounted(fetchPushes);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">推送记录</h1>
      <n-button
        type="primary"
        :disabled="selectedRowKeys.length === 0"
        @click="handleBatchRetry"
      >
        批量重试 ({{ selectedRowKeys.length }})
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4 flex-wrap">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索提交信息"
          clearable
          style="width: 250px"
        >
          <template #prefix
            ><n-icon><SearchOutline /></n-icon
          ></template>
        </n-input>
        <n-select
          v-model:value="searchStatus"
          :options="statusOptions"
          style="width: 150px"
        />
        <n-date-picker
          v-model:value="searchTimeRange"
          type="daterange"
          clearable
          style="width: 300px"
        />
        <n-button @click="fetchPushes">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="pushes"
        :loading="loading"
        v-model:checked-row-keys="selectedRowKeys"
        :pagination="false"
        :bordered="true"
      />
    </n-card>
  </div>
</template>
