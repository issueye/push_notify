<script setup>
import { ref, onMounted, watch, reactive, h } from "vue";
import { formatDate } from "@/utils/date";
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NInput,
  NSelect,
  NModal,
  NForm,
  NFormItem,
  NIcon,
  NPopconfirm,
  NTooltip,
  NPagination,
  useMessage,
  NRadioGroup,
  NRadio,
  NInputNumber,
} from "naive-ui";
import {
  AddOutline,
  TrashOutline,
  RefreshOutline,
  SearchOutline,
  CreateOutline,
  LinkOutline,
} from "@vicons/ionicons5";
import {
  getTargetList,
  createTarget,
  updateTarget,
  deleteTarget,
  testTarget,
} from "@/services/target";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const targets = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchType = ref(null);
const showModal = ref(false);
const submitting = ref(false);
const modalMode = ref("create");
const testingId = ref(null);

const targetTypeOptions = [
  { label: "全部类型", value: null },
  { label: "钉钉", value: "dingtalk" },
  { label: "Webhook", value: "webhook" },
];

const scopeOptions = [
  { label: "全局", value: "global" },
  { label: "指定仓库", value: "repo" },
];

const defaultForm = {
  name: "",
  type: "dingtalk",
  scope: "global",
  config: {
    access_token: "",
    secret: "",
    webhook_url: "",
    method: "POST",
    headers: {},
  },
};

const form = reactive({ ...defaultForm });

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "名称", key: "name", minWidth: 200 },
  {
    title: "类型",
    key: "type",
    width: 100,
    render(row) {
      const typeMap = {
        dingtalk: { type: "info", text: "钉钉" },
        webhook: { type: "warning", text: "Webhook" },
      };
      const info = typeMap[row.type] || { type: "default", text: row.type };
      return h(NTag, { type: info.type, size: "small" }, () => info.text);
    },
  },
  { title: "推送次数", key: "push_count", width: 100 },
  {
    title: "状态",
    key: "status",
    width: 80,
    render(row) {
      return h(
        NTag,
        {
          type: row.status === "active" ? "success" : "default",
          size: "small",
        },
        () => (row.status === "active" ? "启用" : "禁用"),
      );
    },
  },
  {
    title: "创建时间",
    key: "created_at",
    width: 180,
    render(row) {
      return formatDate(row.created_at);
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 170,
    fixed: "right",
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NTooltip,
            { trigger: "hover" },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    quaternary: true,
                    loading: testingId.value === row.id,
                    onClick: () => handleTest(row.id),
                  },
                  {
                    icon: () =>
                      h(NIcon, null, { default: () => h(RefreshOutline) }),
                  },
                ),
              default: () => "测试推送",
            },
          ),
          h(
            NTooltip,
            { trigger: "hover" },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    quaternary: true,
                    onClick: () => handleEdit(row),
                  },
                  {
                    icon: () =>
                      h(NIcon, null, { default: () => h(CreateOutline) }),
                  },
                ),
              default: () => "编辑",
            },
          ),
          h(
            NPopconfirm,
            { onPositiveClick: () => handleDelete(row.id) },
            {
              trigger: () =>
                h(
                  NTooltip,
                  { trigger: "hover" },
                  {
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          quaternary: true,
                          type: "error",
                        },
                        {
                          icon: () =>
                            h(NIcon, null, { default: () => h(TrashOutline) }),
                        },
                      ),
                    default: () => "删除",
                  },
                ),
              default: () => "确定要删除该推送目标吗？",
            },
          ),
        ],
      });
    },
  },
];

async function fetchTargets() {
  loading.value = true;
  try {
    const res = await getTargetList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
      type: searchType.value,
    });
    targets.value = res.data?.list || res.list || [];
    total.value = res.data?.total || res.total || 0;
  } catch (e) {
    message.error("获取推送目标失败");
    targets.value = [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalMode.value = "create";
  Object.assign(form, defaultForm);
  showModal.value = true;
}

function handleEdit(row) {
  modalMode.value = "edit";
  form.id = row.id;
  form.name = row.name;
  form.type = row.type || "dingtalk";
  form.scope = row.scope || "global";
  form.config = row.config
    ? { ...defaultForm.config, ...row.config }
    : { ...defaultForm.config };
  showModal.value = true;
}

async function handleSubmit() {
  if (!form.name) {
    message.warning("请填写名称");
    return;
  }
  // 验证配置
  if (form.type === "dingtalk" && !form.config.access_token) {
    message.warning("请填写AccessToken");
    return;
  }
  if (form.type === "webhook" && !form.config.webhook_url) {
    message.warning("请填写Webhook URL");
    return;
  }
  submitting.value = true;
  try {
    if (modalMode.value === "create") {
      await createTarget(form);
      message.success("创建成功");
    } else {
      await updateTarget(form.id, form);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchTargets();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    submitting.value = false;
  }
}

async function handleTest(id) {
  testingId.value = id;
  try {
    const res = await testTarget(id);
    if (res.type === "webhook" && res.response) {
      message.success(`测试成功 (状态码: ${res.status_code})`);
    } else {
      message.success("测试消息已发送");
    }
  } catch (e) {
    message.error(e.message || "测试失败");
  } finally {
    testingId.value = null;
  }
}

async function handleDelete(id) {
  try {
    await deleteTarget(id);
    message.success("删除成功");
    fetchTargets();
  } catch (e) {
    message.error("删除失败");
  }
}

watch([page, searchKeyword, searchType], fetchTargets);
onMounted(fetchTargets);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">推送目标</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        添加推送目标
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索名称"
          clearable
          style="width: 300px"
          @keyup.enter="fetchTargets"
        >
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
        <n-select
          v-model:value="searchType"
          :options="targetTypeOptions"
          style="width: 150px"
        />
        <n-button @click="fetchTargets">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="targets"
        :loading="loading"
        :pagination="false"
        :bordered="true"
        scroll-x="1200"
      />
      <div class="mt-4 flex justify-end">
        <n-pagination
          v-model:page="page"
          v-model:page-size="size"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50, 100]"
          @update:page="fetchTargets"
          @update:page-size="fetchTargets"
        />
      </div>
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加推送目标' : '编辑推送目标'"
      style="width: 550px"
    >
      <n-form :model="form" label-placement="left" label-width="110">
        <n-form-item label="名称" required>
          <n-input v-model:value="form.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="类型" required>
          <NRadioGroup v-model:value="form.type">
            <NRadio value="dingtalk">钉钉</NRadio>
          </NRadioGroup>
        </n-form-item>
        <!-- 钉钉配置 -->
        <template v-if="form.type === 'dingtalk'">
          <n-form-item label="Webhook URL" required>
            <n-input
              v-model:value="form.config.webhook_url"
              placeholder="https://example.com/webhook"
            />
          </n-form-item>
          <n-form-item label="AccessToken" required>
            <n-input
              v-model:value="form.config.access_token"
              placeholder="钉钉机器人AccessToken"
            />
          </n-form-item>
          <n-form-item label="Secret">
            <n-input
              v-model:value="form.config.secret"
              placeholder="钉钉机器人Secret"
            />
          </n-form-item>
        </template>
      </n-form>
      <div class="flex justify-end gap-2 mt-4">
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ modalMode === "create" ? "创建" : "保存" }}
        </n-button>
      </div>
    </n-modal>
  </div>
</template>
