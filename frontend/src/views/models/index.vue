<script setup>
import { ref, onMounted, reactive, h, watch } from "vue";
import { formatDate } from "@/utils/date";
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NInput,
  NModal,
  NForm,
  NFormItem,
  NIcon,
  NPopconfirm,
  useMessage,
  NInputNumber,
} from "naive-ui";
import {
  AddOutline,
  TrashOutline,
  CheckmarkOutline,
  RefreshOutline,
  CreateOutline,
  SearchOutline,
} from "@vicons/ionicons5";
import {
  getModelList,
  createModel,
  updateModel,
  deleteModel,
  setDefaultModel,
  verifyModel,
} from "@/services/model";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const models = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const showModal = ref(false);
const submitting = ref(false);
const modalMode = ref("create");

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

const form = reactive({ ...defaultForm });

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
    width: 300,
    fixed: "right",
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              onClick: () => handleEdit(row),
            },
            () => "编辑",
          ),
          !row.is_default
            ? h(
                NButton,
                {
                  size: "small",
                  quaternary: true,
                  onClick: () => handleSetDefault(row.id),
                },
                () => "设为默认",
              )
            : null,
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              onClick: () => handleVerify(row.id),
            },
            () => "验证",
          ),
          h(
            NPopconfirm,
            { onPositiveClick: () => handleDelete(row.id) },
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
                  () => "删除",
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

async function fetchModels() {
  loading.value = true;
  try {
    const res = await getModelList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
    });
    models.value = res.data?.list || res.list || [];
    total.value = res.data?.total || res.total || 0;
  } catch (e) {
    message.error("获取模型列表失败");
    models.value = [];
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
  form.provider = row.provider || "";
  form.api_url = row.api_url;
  form.api_key = "";
  form.timeout = row.timeout || 60;
  form.params = row.params || {
    temperature: 0.3,
    max_tokens: 4000,
    top_p: 0.9,
  };
  showModal.value = true;
}

async function handleSubmit() {
  if (!form.name || !form.api_url) {
    message.warning("请填写完整信息");
    return;
  }
  if (modalMode.value === "create" && !form.api_key) {
    message.warning("请填写API密钥");
    return;
  }
  submitting.value = true;
  try {
    if (modalMode.value === "create") {
      await createModel(form);
      message.success("创建成功");
    } else {
      await updateModel(form.id, form);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchModels();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    submitting.value = false;
  }
}

async function handleSetDefault(id) {
  try {
    await setDefaultModel(id);
    message.success("设置成功");
    fetchModels();
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

async function handleDelete(id) {
  try {
    await deleteModel(id);
    message.success("删除成功");
    fetchModels();
  } catch (e) {
    message.error("删除失败");
  }
}

watch([page, searchKeyword], fetchModels);
onMounted(fetchModels);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">AI模型</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon>
          <n-icon>
            <AddOutline />
          </n-icon>
        </template>
        添加模型
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索模型名称"
          clearable
          style="width: 300px"
          @keyup.enter="fetchModels"
        >
          <template #prefix
            ><n-icon><SearchOutline /></n-icon
          ></template>
        </n-input>
        <n-button @click="fetchModels">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="models"
        :loading="loading"
        :pagination="false"
        :bordered="true"
        :scroll-x="1100"
      />
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加AI模型' : '编辑AI模型'"
      style="width: 600px"
    >
      <n-form :model="form" label-placement="top">
        <n-form-item label="模型名称" required>
          <n-input v-model:value="form.name" placeholder="如：GPT-4" />
        </n-form-item>
        <n-form-item label="提供商">
          <n-input v-model:value="form.provider" placeholder="如：OpenAI" />
        </n-form-item>
        <n-form-item label="API地址" required>
          <n-input
            v-model:value="form.api_url"
            placeholder="https://api.openai.com/v1/chat/completions"
          />
        </n-form-item>
        <n-form-item
          :label="
            modalMode === 'create'
              ? 'API密钥（必填）'
              : 'API密钥（留空则不修改）'
          "
          :required="modalMode === 'create'"
        >
          <n-input
            v-model:value="form.api_key"
            type="password"
            placeholder="API Key"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="超时时间(秒)">
          <NInputNumber
            v-model:value="form.timeout"
            :min="10"
            :max="300"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="Temperature">
          <NInputNumber
            v-model:value="form.params.temperature"
            :min="0"
            :max="2"
            :step="0.1"
            :precision="1"
            style="width: 100%"
          />
        </n-form-item>
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
