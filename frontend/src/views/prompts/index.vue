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
  NInputNumber,
} from "naive-ui";
import {
  AddOutline,
  TrashOutline,
  RefreshOutline,
  PlayOutline,
  CreateOutline,
  SearchOutline,
} from "@vicons/ionicons5";
import {
  getPromptList,
  createPrompt,
  updatePrompt,
  deletePrompt,
  testPrompt,
} from "@/services/prompt";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total } = usePagination();

const prompts = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchType = ref(null);
const showModal = ref(false);
const showTestModal = ref(false);
const submitting = ref(false);
const modalMode = ref("create");

const typeOptions = [
  { label: "全部类型", value: null },
  { label: "CODEVIEW", value: "codeview" },
  { label: "消息提示", value: "message" },
];

const defaultForm = {
  name: "",
  type: "codeview",
  scene: "",
  language: "",
  content: "",
};

const form = reactive({ ...defaultForm });

const testData = ref(
  '// 测试代码\nfunction hello() {\n  console.log("Hello World");\n}',
);
const currentPrompt = ref(null);
const testResult = ref("");

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "名称", key: "name" },
  {
    title: "类型",
    key: "type",
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          type: row.type === "codeview" ? "info" : "warning",
          size: "small",
        },
        () => (row.type === "codeview" ? "CODEVIEW" : "消息"),
      );
    },
  },
  { title: "场景", key: "scene", width: 120 },
  { title: "语言", key: "language", width: 80 },
  { title: "版本", key: "version", width: 70 },
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
    width: 200,
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
                    onClick: () => openTest(row),
                  },
                  {
                    icon: () =>
                      h(NIcon, null, { default: () => h(PlayOutline) }),
                  },
                ),
              default: () => "测试提示词",
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
              default: () => "确定要删除该提示词吗？",
            },
          ),
        ],
      });
    },
  },
];

async function fetchPrompts() {
  loading.value = true;
  try {
    const res = await getPromptList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
      type: searchType.value,
    });
    prompts.value = res.data?.list || res.list || [];
    total.value = res.data?.total || res.total || 0;
  } catch (e) {
    message.error("获取提示词列表失败");
    prompts.value = [];
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
  form.type = row.type;
  form.scene = row.scene || "";
  form.language = row.language || "";
  form.content = row.content;
  showModal.value = true;
}

async function handleSubmit() {
  if (!form.name || !form.content) {
    message.warning("请填写完整信息");
    return;
  }
  submitting.value = true;
  try {
    if (modalMode.value === "create") {
      await createPrompt(form);
      message.success("创建成功");
    } else {
      await updatePrompt(form.id, form);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchPrompts();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    submitting.value = false;
  }
}

function openTest(row) {
  currentPrompt.value = row;
  testData.value =
    '// 测试代码\nfunction hello() {\n  console.log("Hello World");\n}';
  testResult.value = "";
  showTestModal.value = true;
}

async function handleTest() {
  try {
    const res = await testPrompt(currentPrompt.value.id, {
      file_content: testData.value,
    });
    message.success("测试完成");
    testResult.value = res.data?.result || res.result || JSON.stringify(res);
  } catch (e) {
    message.error("测试失败");
    testResult.value = e.message || "测试失败";
  }
}

async function handleDelete(id) {
  try {
    await deletePrompt(id);
    message.success("删除成功");
    fetchPrompts();
  } catch (e) {
    message.error("删除失败");
  }
}

watch([page, searchKeyword, searchType], fetchPrompts);
onMounted(fetchPrompts);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">提示词管理</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon
          ><n-icon><AddOutline /></n-icon
        ></template>
        添加提示词
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索提示词名称"
          clearable
          style="width: 300px"
          @keyup.enter="fetchPrompts"
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
        <n-button @click="fetchPrompts">搜索</n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="prompts"
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
          @update:page="fetchPrompts"
          @update:page-size="fetchPrompts"
        />
      </div>
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加提示词' : '编辑提示词'"
      style="width: 700px"
    >
      <n-form :model="form" label-placement="top">
        <n-form-item label="名称" required>
          <n-input v-model:value="form.name" placeholder="提示词名称" />
        </n-form-item>
        <n-form-item label="类型" required>
          <NRadioGroup v-model:value="form.type">
            <NRadio value="codeview">CODEVIEW</NRadio>
            <NRadio value="message">消息提示</NRadio>
          </NRadioGroup>
        </n-form-item>
        <n-form-item label="场景">
          <n-input
            v-model:value="form.scene"
            placeholder="如：代码规范检查、安全检查等"
          />
        </n-form-item>
        <n-form-item label="适用语言">
          <n-input
            v-model:value="form.language"
            placeholder="如：Go、Python、JavaScript"
          />
        </n-form-item>
        <n-form-item label="提示词内容" required>
          <n-input
            v-model:value="form.content"
            type="textarea"
            :rows="10"
            placeholder="提示词模板，支持变量如 {{.FileContent}}, {{.FileName}}"
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

    <n-modal
      v-model:show="showTestModal"
      preset="card"
      title="测试提示词"
      style="width: 700px"
    >
      <p class="mb-4 text-gray-600">提示词：{{ currentPrompt?.name }}</p>
      <n-form-item label="测试数据">
        <n-input v-model:value="testData" type="textarea" :rows="8" />
      </n-form-item>
      <n-form-item v-if="testResult" label="测试结果">
        <n-input
          v-model:value="testResult"
          type="textarea"
          :rows="4"
          readonly
        />
      </n-form-item>
      <div class="flex justify-end gap-2 mt-4">
        <n-button @click="showTestModal = false">关闭</n-button>
        <n-button type="primary" @click="handleTest">开始测试</n-button>
      </div>
    </n-modal>
  </div>
</template>
