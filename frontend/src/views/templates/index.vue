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
} from "naive-ui";
import {
  AddOutline,
  TrashOutline,
  RefreshOutline,
  CreateOutline,
  SearchOutline,
  LockClosedOutline,
  LockOpenOutline,
} from "@vicons/ionicons5";
import {
  getTemplateList,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  updateTemplateStatus,
  getTemplateDetail,
  generateTemplate,
} from "@/services/template";
import { getModelList } from "@/services/model";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const templates = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchType = ref(null);
const showModal = ref(false);
const submitting = ref(false);
const modalMode = ref("create");
const generating = ref(false);
const models = ref([]);
const modelOptions = ref([]);
const selectedModelId = ref(null);

const typeOptions = [
  { label: "全部类型", value: null },
  { label: "钉钉", value: "dingtalk" },
];

const sceneOptions = [
  { label: "代码提交通知", value: "commit_notify" },
  { label: "审查结果通知", value: "review_notify" },
];

const defaultForm = {
  name: "",
  type: "dingtalk",
  scene: "commit_notify",
  title: "",
  content: "",
};

const form = reactive({ ...defaultForm });

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "名称", key: "name" },
  {
    title: "类型",
    key: "type",
    width: 80,
    render() {
      return h(NTag, { type: "info", size: "small" }, () => "钉钉");
    },
  },
  {
    title: "场景",
    key: "scene",
    width: 120,
    render(row) {
      const sceneMap = {
        commit_notify: "代码提交",
        review_notify: "审查结果",
      };
      return sceneMap[row.scene] || row.scene;
    },
  },
  { title: "版本", key: "version", width: 70 },
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
            NTooltip,
            { trigger: "hover" },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    quaternary: true,
                    onClick: () => toggleStatus(row),
                  },
                  {
                    icon: () =>
                      h(NIcon, null, {
                        default: () =>
                          row.status === "active"
                            ? h(LockClosedOutline)
                            : h(LockOpenOutline),
                      }),
                  },
                ),
              default: () => (row.status === "active" ? "禁用" : "启用"),
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
              default: () => "确定要删除该模板吗？",
            },
          ),
        ],
      });
    },
  },
];

async function fetchTemplates() {
  loading.value = true;
  try {
    const res = await getTemplateList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
      type: searchType.value,
    });
    templates.value = res.data?.list || res.list || [];
    total.value = res.data?.total || res.total || 0;
  } catch (e) {
    message.error("获取模板列表失败");
    templates.value = [];
  } finally {
    loading.value = false;
  }
}

async function fetchModels() {
  try {
    const data = await getModelList({ page: 1, size: 100 });
    models.value = data.list || [];
    modelOptions.value = models.value.map((m) => ({
      label: m.name,
      value: m.id,
    }));
    if (models.value.length > 0) {
      selectedModelId.value = models.value[0].id;
    }
  } catch (e) {
    console.error("Failed to fetch models", e);
  }
}

async function handleGenerate() {
  if (!form.name || !form.title) {
    message.warning("请先填写模板名称和标题，以便AI更准确地生成内容");
    return;
  }

  generating.value = true;
  try {
    const res = await generateTemplate({
      name: form.name,
      type: form.type,
      scene: form.scene,
      title: form.title,
      model_id: selectedModelId.value,
    });

    if (res && res.content) {
      form.content = res.content;
      message.success("模板生成成功");
    }
  } catch (e) {
    message.error("模板生成失败");
  } finally {
    generating.value = false;
  }
}

function handleAdd() {
  modalMode.value = "create";
  Object.assign(form, defaultForm);
  if (modelOptions.value.length === 0) {
    fetchModels();
  }
  showModal.value = true;
}

async function handleEdit(row) {
  modalMode.value = "edit";
  form.id = row.id;
  form.name = row.name;
  form.type = row.type;
  form.scene = row.scene;
  form.title = row.title;
  form.content = row.content;

  if (modelOptions.value.length === 0) {
    fetchModels();
  }

  showModal.value = true;
}

async function handleSubmit() {
  if (!form.name || !form.title || !form.content) {
    message.warning("请填写完整信息");
    return;
  }
  submitting.value = true;
  try {
    if (modalMode.value === "create") {
      await createTemplate(form);
      message.success("创建成功");
    } else {
      await updateTemplate(form.id, form);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchTemplates();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    submitting.value = false;
  }
}

async function toggleStatus(row) {
  const newStatus = row.status === "active" ? "inactive" : "active";
  try {
    await updateTemplateStatus(row.id, newStatus);
    row.status = newStatus;
    message.success(newStatus === "active" ? "已启用" : "已禁用");
  } catch (e) {
    message.error("操作失败");
  }
}

async function handleDelete(id) {
  try {
    await deleteTemplate(id);
    message.success("删除成功");
    fetchTemplates();
  } catch (e) {
    message.error("删除失败");
  }
}

watch([page, searchKeyword, searchType], fetchTemplates);
onMounted(fetchTemplates);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">消息模板</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon
          ><n-icon><AddOutline /></n-icon
        ></template>
        添加模板
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索模板名称"
          clearable
          style="width: 300px"
          @keyup.enter="fetchTemplates"
        >
          <template #prefix
            ><n-icon><SearchOutline /></n-icon
          ></template>
        </n-input>
        <n-select
          v-model:value="searchType"
          :options="typeOptions"
          style="width: 150px"
        />
        <n-button @click="fetchTemplates">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="templates"
        :loading="loading"
        :pagination="false"
        :bordered="true"
        :scroll-x="1500"
      />
      <div class="mt-4 flex justify-end">
        <n-pagination
          v-model:page="page"
          v-model:page-size="size"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50, 100]"
          @update:page="fetchTemplates"
          @update:page-size="fetchTemplates"
        />
      </div>
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加模板' : '编辑模板'"
      style="width: 700px"
    >
      <n-form :model="form" label-placement="top">
        <n-form-item label="模板名称" required>
          <n-input v-model:value="form.name" placeholder="请输入模板名称" />
        </n-form-item>
        <n-form-item label="模板类型">
          <NRadioGroup v-model:value="form.type">
            <NRadio value="dingtalk">钉钉</NRadio>
          </NRadioGroup>
        </n-form-item>
        <n-form-item label="使用场景" required>
          <n-select v-model:value="form.scene" :options="sceneOptions" />
        </n-form-item>
        <n-form-item label="模板标题" required>
          <n-input
            v-model:value="form.title"
            placeholder="消息标题，支持变量替换"
          />
        </n-form-item>
        <n-form-item label="模板内容" required>
          <div class="flex flex-col w-full gap-2">
            <div class="flex items-center gap-2 mb-2 bg-gray-50 p-3 rounded">
              <span class="text-sm text-gray-500 whitespace-nowrap"
                >AI 辅助生成:</span
              >
              <n-select
                v-model:value="selectedModelId"
                :options="modelOptions"
                placeholder="选择模型"
                size="small"
                style="width: 200px"
              />
              <n-button
                type="info"
                size="small"
                :loading="generating"
                @click="handleGenerate"
                :disabled="!form.name || !form.title"
              >
                <template #icon
                  ><n-icon><CreateOutline /></n-icon
                ></template>
                一键生成内容
              </n-button>
              <span class="text-xs text-gray-400 ml-2">需先填写名称和标题</span>
            </div>
            <n-input
              v-model:value="form.content"
              type="textarea"
              :rows="8"
              placeholder="消息内容，支持变量替换，如 {{.RepoName}}, {{.CommitMsg}} 等"
            />
          </div>
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
