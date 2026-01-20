<script setup>
import { ref, reactive, onMounted, watch, h } from "vue";
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
  useMessage,
  NDynamicInput,
  NGrid,
  NGridItem,
  NSelect as NSelectOption, // Alias if needed, but NSelect is fine
  NDivider,
} from "naive-ui";
import {
  AddOutline,
  CreateOutline,
  TrashOutline,
  RefreshOutline,
  SearchOutline,
} from "@vicons/ionicons5";
import {
  getRepoList,
  getRepoDetail,
  createRepo,
  updateRepo,
  deleteRepo,
  testWebhook,
  getRepoTargets,
} from "@/services/repo";
import { getTargetList } from "@/services/target";
import { getModelList } from "@/services/model";
import { getTemplateList } from "@/services/template";
import { usePagination, useConfirm } from "@/composables/useMessage";

const message = useMessage();
const { confirm } = useConfirm();
const { page, size, total, setPage } = usePagination();

const repos = ref([]);
const targets = ref([]);
const targetOptions = ref([]);
const models = ref([]);
const modelOptions = ref([]);
const templates = ref([]);
const commitTemplateOptions = ref([]);
const reviewTemplateOptions = ref([]);
const loading = ref(false);
const targetLoading = ref(false);
const modelLoading = ref(false);
const templateLoading = ref(false);
const searchKeyword = ref("");
const showModal = ref(false);
const formRef = ref(null);
const submitting = ref(false);
const modalMode = ref("create");

const repoTypeOptions = [
  { label: "GitHub", value: "github" },
  { label: "GitLab", value: "gitlab" },
  { label: "Gitee", value: "gitee" },
];

const defaultForm = {
  name: "",
  url: "",
  type: "github",
  access_token: "",
  webhook_secret: "",
  target_ids: [],
  model_id: null,
  commit_template_id: null,
  review_templates: [], // [{ template_id: 1, language: 'Go' }]
};

const languageOptions = [
  { label: "默认", value: "default" },
  { label: "Go", value: "Go" },
  { label: "Java", value: "Java" },
  { label: "Python", value: "Python" },
  { label: "JavaScript", value: "JavaScript" },
  { label: "TypeScript", value: "TypeScript" },
  { label: "Vue", value: "Vue" },
  { label: "PHP", value: "PHP" },
  { label: "Rust", value: "Rust" },
  { label: "C++", value: "C++" },
];

const form = reactive({ ...defaultForm });

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "名称", key: "name", width: 200, ellipsis: { tooltip: true } },
  { title: "类型", key: "type", width: 100 },
  {
    title: "关联模型",
    key: "model",
    width: 150,
    render(row) {
      if (!row.model_id) return "-";
      const model = models.value.find((m) => m.id === row.model_id);
      return model ? model.name : `模型 #${row.model_id}`;
    },
  },
  {
    title: "Webhook URL",
    key: "webhook_url",
    minWidth: 250,
    ellipsis: { tooltip: true },
  },
  {
    title: "创建时间",
    key: "created_at",
    width: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return formatDate(row.created_at);
    },
  },
  {
    title: "状态",
    key: "status",
    width: 100,
    fixed: "right",
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
    title: "操作",
    key: "actions",
    width: 180,
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
            {
              icon: () => h(NIcon, null, { default: () => h(CreateOutline) }),
            },
          ),
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              onClick: () => handleTest(row.id),
            },
            {
              icon: () => h(NIcon, null, { default: () => h(RefreshOutline) }),
            },
          ),
          h(
            NPopconfirm,
            { onPositiveClick: () => handleDelete(row.id) },
            {
              trigger: () =>
                h(
                  NButton,
                  { size: "small", quaternary: true, type: "error" },
                  {
                    icon: () =>
                      h(NIcon, null, { default: () => h(TrashOutline) }),
                  },
                ),
              default: () => "确定要删除该仓库吗？",
            },
          ),
        ],
      });
    },
  },
];

async function fetchRepos() {
  loading.value = true;
  try {
    const data = await getRepoList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
    });
    repos.value = data.list;
    total.value = data.pagination.total;
  } catch (e) {
    message.error("获取仓库列表失败");
  } finally {
    loading.value = false;
  }
}

async function fetchTargets() {
  targetLoading.value = true;
  try {
    const data = await getTargetList({ page: 1, size: 100 });
    targets.value = data.list || [];
    targetOptions.value = (data.list || []).map((t) => ({
      label: t.name,
      value: t.id,
    }));
  } catch (e) {
    console.error("获取推送目标失败", e);
  } finally {
    targetLoading.value = false;
  }
}

async function fetchModels() {
  modelLoading.value = true;
  try {
    const data = await getModelList({ page: 1, size: 100 });
    models.value = data.list || [];
    modelOptions.value = (data.list || []).map((m) => ({
      label: m.name,
      value: m.id,
    }));
  } catch (e) {
    console.error("获取模型列表失败", e);
  } finally {
    modelLoading.value = false;
  }
}

async function fetchTemplates() {
  templateLoading.value = true;
  try {
    const data = await getTemplateList({ page: 1, size: 100 });
    templates.value = data.list || [];

    commitTemplateOptions.value = templates.value
      .filter((t) => t.scene === "commit_notify")
      .map((t) => ({ label: t.name, value: t.id }));

    reviewTemplateOptions.value = templates.value
      .filter((t) => t.scene === "review_notify")
      .map((t) => ({ label: t.name, value: t.id }));
  } catch (e) {
    console.error("获取模板列表失败", e);
  } finally {
    templateLoading.value = false;
  }
}

function handleAdd() {
  modalMode.value = "create";
  Object.assign(form, defaultForm);
  // 重置 review_templates 为空数组，确保 dynamic input 正常工作
  form.review_templates = [];
  if (targetOptions.value.length === 0) {
    fetchTargets();
  }
  if (commitTemplateOptions.value.length === 0) {
    fetchTemplates();
  }
  showModal.value = true;
}

async function handleEdit(row) {
  modalMode.value = "edit";
  form.id = row.id;
  form.name = row.name;
  form.url = row.url;
  form.type = row.type;
  form.access_token = "";
  form.webhook_secret = row.webhook_secret || "";
  form.target_ids = [];
  form.model_id = row.model_id || null;
  form.commit_template_id = row.commit_template_id || null;
  form.review_templates = [];

  if (row.review_templates && row.review_templates.length > 0) {
    form.review_templates = row.review_templates.map((rt) => ({
      template_id: rt.template_id,
      language: rt.language,
    }));
  }

  // 确保列表已加载
  if (targetOptions.value.length === 0) {
    await fetchTargets();
  }
  if (modelOptions.value.length === 0) {
    await fetchModels();
  }
  if (commitTemplateOptions.value.length === 0) {
    await fetchTemplates();
  }

  // 获取已绑定的推送目标
  try {
    const targets = await getRepoTargets(row.id);
    form.target_ids = targets.map((t) => t.id);
  } catch (e) {
    console.error("获取推送目标失败", e);
  }

  showModal.value = true;
}

async function handleSubmit() {
  if (!form.name || !form.url) {
    message.warning("请填写完整信息");
    return;
  }
  submitting.value = true;
  try {
    if (modalMode.value === "create") {
      await createRepo(form);
      message.success("创建成功");
    } else {
      await updateRepo(form.id, form);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchRepos();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    submitting.value = false;
  }
}

async function handleTest(id) {
  try {
    await testWebhook(id);
    message.success("测试成功");
  } catch (e) {
    message.error("测试失败");
  }
}

async function handleDelete(id) {
  try {
    await deleteRepo(id);
    message.success("删除成功");
    fetchRepos();
  } catch (e) {
    message.error("删除失败");
  }
}

watch([page, searchKeyword], () => fetchRepos());

onMounted(() => {
  fetchRepos();
  fetchTargets();
  fetchModels();
  fetchTemplates();
});
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">仓库管理</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon
          ><n-icon><AddOutline /></n-icon
        ></template>
        添加仓库
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索仓库名称"
          clearable
          style="width: 300px"
          @keyup.enter="fetchRepos"
        >
          <template #prefix
            ><n-icon><SearchOutline /></n-icon
          ></template>
        </n-input>
        <n-button @click="fetchRepos">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="repos"
        :loading="loading"
        :pagination="false"
        :bordered="true"
        :scroll-x="1500"
      />
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加仓库' : '编辑仓库'"
      style="width: 600px"
    >
      <n-form :model="form" label-placement="left" label-width="100">
        <n-form-item label="仓库名称" required>
          <n-input v-model:value="form.name" placeholder="请输入仓库名称" />
        </n-form-item>
        <n-form-item label="仓库地址" required>
          <n-input
            v-model:value="form.url"
            placeholder="https://github.com/xxx/xxx"
          />
        </n-form-item>
        <n-form-item label="仓库类型" required>
          <n-select v-model:value="form.type" :options="repoTypeOptions" />
        </n-form-item>
        <n-form-item label="访问令牌">
          <n-input
            v-model:value="form.access_token"
            type="password"
            show-password-on="click"
            placeholder="请输入访问令牌 (AccessToken)，留空不修改"
          />
        </n-form-item>
        <n-form-item label="关联模型">
          <n-select
            v-model:value="form.model_id"
            clearable
            :options="modelOptions"
            :loading="modelLoading"
            placeholder="选择关联的AI模型（用于CodeView）"
          />
        </n-form-item>
        <n-form-item label="推送目标">
          <n-select
            v-model:value="form.target_ids"
            multiple
            :options="targetOptions"
            :loading="targetLoading"
            placeholder="选择推送目标"
          />
        </n-form-item>

        <n-divider title-placement="left">模板配置</n-divider>

        <n-form-item label="提交通知模板">
          <n-select
            v-model:value="form.commit_template_id"
            clearable
            :options="commitTemplateOptions"
            :loading="templateLoading"
            placeholder="选择代码提交通知模板（留空使用系统默认）"
          />
        </n-form-item>

        <n-form-item label="审查通知模板">
          <n-dynamic-input
            v-model:value="form.review_templates"
            :on-create="() => ({ template_id: null, language: 'default' })"
          >
            <template #default="{ value }">
              <div style="display: flex; gap: 8px; width: 100%">
                <n-select
                  v-model:value="value.template_id"
                  :options="reviewTemplateOptions"
                  placeholder="选择模板"
                  style="width: 60%"
                />
                <n-select
                  v-model:value="value.language"
                  :options="languageOptions"
                  placeholder="适用语言"
                  style="width: 40%"
                />
              </div>
            </template>
          </n-dynamic-input>
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
