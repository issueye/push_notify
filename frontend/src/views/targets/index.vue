<script setup>
import { ref, h } from "vue";
import { formatDate } from "@/utils/date";
import {
  NButton,
  NSpace,
  NTag,
  NInput,
  NSelect,
  NForm,
  NFormItem,
  NIcon,
  NPopconfirm,
  NTooltip,
  useMessage,
  NRadioGroup,
  NRadio,
} from "naive-ui";
import {
  TrashOutline,
  RefreshOutline,
  CreateOutline,
} from "@vicons/ionicons5";
import {
  getTargetList,
  createTarget,
  updateTarget,
  deleteTarget,
  testTarget,
} from "@/services/target";
import { useCurd } from "@/composables/useCurd";
import CurdPage from "@/components/common/CurdPage.vue";

const message = useMessage();

const targetTypeOptions = [
  { label: "全部类型", value: null },
  { label: "钉钉", value: "dingtalk" },
  { label: "Webhook", value: "webhook" },
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

const {
  list: targets,
  loading,
  total,
  page,
  size,
  searchParams,
  showModal,
  modalMode,
  submitting,
  form,
  formRef,
  fetchData,
  handleSearch,
  handleAdd,
  handleEdit,
  handleSubmit,
  handleDelete,
} = useCurd({
  fetchList: getTargetList,
  createItem: createTarget,
  updateItem: updateTarget,
  deleteItem: deleteTarget,
  defaultForm,
  beforeSubmit: (data) => {
    if (!data.name) {
      throw new Error("请填写名称");
    }
    if (data.type === "dingtalk" && !data.config.access_token) {
      throw new Error("请填写AccessToken");
    }
    if (data.type === "webhook" && !data.config.webhook_url) {
      throw new Error("请填写Webhook URL");
    }
    return data;
  },
});

const testingId = ref(null);

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
</script>

<template>
  <CurdPage
    title="推送目标"
    v-model:page="page"
    v-model:page-size="size"
    v-model:show-modal="showModal"
    :loading="loading"
    :columns="columns"
    :data="targets"
    :item-count="total"
    :modal-title="modalMode === 'create' ? '添加推送目标' : '编辑推送目标'"
    :submitting="submitting"
    @search="handleSearch"
    @add="handleAdd"
    @submit="handleSubmit"
  >
    <template #search>
      <n-input
        v-model:value="searchParams.keyword"
        placeholder="搜索名称"
        clearable
        style="width: 300px"
        @keyup.enter="handleSearch"
      />
      <n-select
        v-model:value="searchParams.type"
        :options="targetTypeOptions"
        placeholder="选择类型"
        clearable
        style="width: 150px"
        @update:value="handleSearch"
      />
    </template>

    <template #form>
      <n-form ref="formRef" :model="form" label-placement="left" label-width="110">
        <n-form-item label="名称" path="name" required>
          <n-input v-model:value="form.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="类型" path="type" required>
          <n-radio-group v-model:value="form.type">
            <n-radio value="dingtalk">钉钉</n-radio>
          </n-radio-group>
        </n-form-item>
        <template v-if="form.type === 'dingtalk'">
          <n-form-item label="Webhook URL" path="config.webhook_url" required>
            <n-input
              v-model:value="form.config.webhook_url"
              placeholder="https://example.com/webhook"
            />
          </n-form-item>
          <n-form-item label="AccessToken" path="config.access_token" required>
            <n-input
              v-model:value="form.config.access_token"
              placeholder="钉钉机器人AccessToken"
            />
          </n-form-item>
          <n-form-item label="Secret" path="config.secret">
            <n-input
              v-model:value="form.config.secret"
              placeholder="钉钉机器人Secret"
            />
          </n-form-item>
        </template>
      </n-form>
    </template>
  </CurdPage>
</template>
