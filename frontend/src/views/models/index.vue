<script setup>
import { h } from "vue";
import { formatDate } from "@/utils/date";
import {
  NButton,
  NSpace,
  NTag,
  NInput,
  NForm,
  NFormItem,
  NIcon,
  NPopconfirm,
  NTooltip,
  useMessage,
  NInputNumber,
} from "naive-ui";
import {
  TrashOutline,
  CheckmarkOutline,
  RefreshOutline,
  CreateOutline,
  StarOutline,
  ShieldCheckmarkOutline,
} from "@vicons/ionicons5";
import {
  getModelList,
  createModel,
  updateModel,
  deleteModel,
  setDefaultModel,
  verifyModel,
} from "@/services/model";
import { useCurd } from "@/composables/useCurd";
import CurdPage from "@/components/common/CurdPage.vue";

const message = useMessage();

const defaultForm = {
  name: "",
  provider: "",
  api_url: "",
  api_key: "",
  timeout: 60,
  params: {
    temperature: 0.3,
    max_tokens: 4000,
    top_p: 0.9,
  },
};

const {
  list: models,
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
  fetchList: getModelList,
  createItem: createModel,
  updateItem: updateModel,
  deleteItem: deleteModel,
  defaultForm,
  beforeSubmit: (data) => {
    if (!data.name || !data.api_url) {
      throw new Error("请填写完整信息");
    }
    if (modalMode.value === "create" && !data.api_key) {
      throw new Error("请填写API密钥");
    }
    return data;
  },
});

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "名称", key: "name", ellipsis: { tooltip: true } },
  { title: "提供商", key: "provider", width: 100, ellipsis: { tooltip: true } },
  { title: "API地址", key: "api_url", ellipsis: { tooltip: true } },
  { title: "调用次数", key: "call_count", width: 100 },
  {
    title: "默认",
    key: "is_default",
    width: 70,
    render(row) {
      return row.is_default
        ? h(NIcon, { color: "#18a058" }, { default: () => h(CheckmarkOutline) })
        : "";
    },
  },
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
    width: 170,
    render(row) {
      return formatDate(row.created_at);
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 220,
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
          !row.is_default
            ? h(
                NTooltip,
                { trigger: "hover" },
                {
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        quaternary: true,
                        onClick: () => handleSetDefault(row.id),
                      },
                      {
                        icon: () =>
                          h(NIcon, null, { default: () => h(StarOutline) }),
                      },
                    ),
                  default: () => "设为默认",
                },
              )
            : null,
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
                    onClick: () => handleVerify(row.id),
                  },
                  {
                    icon: () =>
                      h(NIcon, null, {
                        default: () => h(ShieldCheckmarkOutline),
                      }),
                  },
                ),
              default: () => "验证配置",
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
                          disabled: row.is_default,
                        },
                        {
                          icon: () =>
                            h(NIcon, null, { default: () => h(TrashOutline) }),
                        },
                      ),
                    default: () => "删除",
                  },
                ),
              default: () =>
                row.is_default ? "默认模型不能删除" : "确定要删除该模型吗？",
            },
          ),
        ],
      });
    },
  },
];

async function handleSetDefault(id) {
  try {
    await setDefaultModel(id);
    message.success("设置成功");
    fetchData();
  } catch (e) {
    message.error("设置失败");
  }
}

async function handleVerify(id) {
  try {
    await verifyModel(id);
    message.success("验证成功，API配置正确");
  } catch (e) {
    message.error("验证失败，" + (e.message || "请检查API配置"));
  }
}
</script>

<template>
  <CurdPage
    title="AI模型"
    v-model:page="page"
    v-model:page-size="size"
    v-model:show-modal="showModal"
    :loading="loading"
    :columns="columns"
    :data="models"
    :item-count="total"
    :modal-title="modalMode === 'create' ? '添加AI模型' : '编辑AI模型'"
    :submitting="submitting"
    @search="handleSearch"
    @add="handleAdd"
    @submit="handleSubmit"
  >
    <template #search>
      <n-input
        v-model:value="searchParams.keyword"
        placeholder="搜索模型名称"
        clearable
        style="width: 300px"
        @keyup.enter="handleSearch"
      />
    </template>

    <template #form>
      <n-form ref="formRef" :model="form" label-placement="top">
        <n-form-item label="模型名称" path="name" required>
          <n-input v-model:value="form.name" placeholder="如：GPT-4" />
        </n-form-item>
        <n-form-item label="提供商" path="provider">
          <n-input v-model:value="form.provider" placeholder="如：OpenAI" />
        </n-form-item>
        <n-form-item label="API地址" path="api_url" required>
          <n-input
            v-model:value="form.api_url"
            placeholder="https://api.openai.com/v1/chat/completions"
          />
        </n-form-item>
        <n-form-item
          label="API密钥"
          path="api_key"
          :required="modalMode === 'create'"
        >
          <n-input
            v-model:value="form.api_key"
            type="password"
            :placeholder="
              modalMode === 'create' ? 'API Key' : 'API密钥（留空则不修改）'
            "
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="超时时间(秒)" path="timeout">
          <n-input-number
            v-model:value="form.timeout"
            :min="10"
            :max="300"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="Temperature" path="params.temperature">
          <n-input-number
            v-model:value="form.params.temperature"
            :min="0"
            :max="2"
            :step="0.1"
            :precision="1"
            style="width: 100%"
          />
        </n-form-item>
      </n-form>
    </template>
  </CurdPage>
</template>

