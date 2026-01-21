<script setup>
import { ref, h } from "vue";
import { formatDate } from "@/utils/date";
import {
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
  useMessage,
  NRadioGroup,
  NRadio,
} from "naive-ui";
import {
  TrashOutline,
  PlayOutline,
  CreateOutline,
} from "@vicons/ionicons5";
import {
  getPromptList,
  createPrompt,
  updatePrompt,
  deletePrompt,
  testPrompt,
} from "@/services/prompt";
import { useCurd } from "@/composables/useCurd";
import CurdPage from "@/components/common/CurdPage.vue";

const message = useMessage();

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

const {
  list: prompts,
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
  fetchList: getPromptList,
  createItem: createPrompt,
  updateItem: updatePrompt,
  deleteItem: deletePrompt,
  defaultForm,
  beforeSubmit: (data) => {
    if (!data.name || !data.content) {
      throw new Error("请填写完整信息");
    }
    return data;
  },
});

const showTestModal = ref(false);
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
</script>

<template>
  <CurdPage
    title="提示词管理"
    v-model:page="page"
    v-model:page-size="size"
    v-model:show-modal="showModal"
    :loading="loading"
    :columns="columns"
    :data="prompts"
    :item-count="total"
    :modal-title="modalMode === 'create' ? '添加提示词' : '编辑提示词'"
    :submitting="submitting"
    @search="handleSearch"
    @add="handleAdd"
    @submit="handleSubmit"
  >
    <template #search>
      <n-input
        v-model:value="searchParams.keyword"
        placeholder="搜索提示词名称"
        clearable
        style="width: 300px"
        @keyup.enter="handleSearch"
      />
      <n-select
        v-model:value="searchParams.type"
        :options="typeOptions"
        placeholder="选择类型"
        clearable
        style="width: 150px"
        @update:value="handleSearch"
      />
    </template>

    <template #form>
      <n-form ref="formRef" :model="form" label-placement="top">
        <n-form-item label="名称" path="name" required>
          <n-input v-model:value="form.name" placeholder="提示词名称" />
        </n-form-item>
        <n-form-item label="类型" path="type" required>
          <n-radio-group v-model:value="form.type">
            <n-radio value="codeview">CODEVIEW</n-radio>
            <n-radio value="message">消息提示</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="场景" path="scene">
          <n-input
            v-model:value="form.scene"
            placeholder="如：代码规范检查、安全检查等"
          />
        </n-form-item>
        <n-form-item label="适用语言" path="language">
          <n-input
            v-model:value="form.language"
            placeholder="如：Go、Python、JavaScript"
          />
        </n-form-item>
        <n-form-item label="提示词内容" path="content" required>
          <n-input
            v-model:value="form.content"
            type="textarea"
            :rows="10"
            placeholder="提示词模板，支持变量如 {{.FileContent}}, {{.FileName}}"
          />
        </n-form-item>
      </n-form>
    </template>
  </CurdPage>

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
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="showTestModal = false">关闭</n-button>
        <n-button type="primary" @click="handleTest">开始测试</n-button>
      </div>
    </template>
  </n-modal>
</template>
