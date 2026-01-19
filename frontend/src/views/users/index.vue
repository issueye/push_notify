<script setup>
import { ref, onMounted, h, watch } from "vue";
import { formatDate } from "@/utils/date";
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NInput,
  NSelect,
  useMessage,
  NModal,
  NForm,
  NFormItem,
  NPopconfirm,
  NIcon,
} from "naive-ui";
import {
  AddOutline,
  TrashOutline,
  LockOpenOutline,
  LockClosedOutline,
  RefreshOutline,
} from "@vicons/ionicons5";
import {
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  resetUserPassword,
  lockUser,
} from "@/services/user";
import { usePagination } from "@/composables/useMessage";

const message = useMessage();
const { page, size, total, updatePagination } = usePagination();

const users = ref([]);
const loading = ref(false);
const searchKeyword = ref("");
const searchRole = ref(null);

const showModal = ref(false);
const modalMode = ref("create"); // 'create' | 'edit'
const currentUser = ref({ id: null, username: "", email: "", role: "user" });
const formLoading = ref(false);

const roleOptions = [
  { label: "全部角色", value: null },
  { label: "管理员", value: "admin" },
  { label: "普通用户", value: "user" },
];

const columns = [
  { title: "ID", key: "id", width: 60 },
  { title: "用户名", key: "username" },
  { title: "邮箱", key: "email" },
  {
    title: "角色",
    key: "role",
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          type: row.role === "admin" ? "error" : "info",
          size: "small",
        },
        () => (row.role === "admin" ? "管理员" : "普通用户"),
      );
    },
  },
  {
    title: "状态",
    key: "status",
    width: 90,
    render(row) {
      return h(
        NTag,
        {
          type: row.status === "active" ? "success" : "default",
          size: "small",
        },
        () => (row.status === "active" ? "正常" : "锁定"),
      );
    },
  },
  {
    title: "最后登录",
    key: "last_login_at",
    width: 170,
    render(row) {
      return formatDate(row.last_login_at);
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
    width: 180,
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              onClick: () => handleResetPassword(row.id),
            },
            () => "重置密码",
          ),
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              onClick: () => handleToggleStatus(row),
            },
            () => (row.status === "active" ? "锁定" : "解锁"),
          ),
        ],
      });
    },
  },
];

async function fetchUsers() {
  loading.value = true;
  try {
    const res = await getUserList({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value,
      role: searchRole.value,
    });
    console.log("res", res);
    users.value = res.list || [];
    total.value = res.pagination.total || 0;
  } catch (e) {
    message.error("获取用户列表失败");
    users.value = [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalMode.value = "create";
  currentUser.value = {
    id: null,
    username: "",
    email: "",
    role: "user",
    password: "",
  };
  showModal.value = true;
}

function handleEdit(row) {
  modalMode.value = "edit";
  currentUser.value = { ...row, password: "" };
  showModal.value = true;
}

async function handleSubmit() {
  if (!currentUser.value.username || !currentUser.value.email) {
    message.warning("请填写用户名和邮箱");
    return;
  }

  formLoading.value = true;
  try {
    if (modalMode.value === "create") {
      await createUser(currentUser.value);
      message.success("创建成功");
    } else {
      await updateUser(currentUser.value.id, currentUser.value);
      message.success("更新成功");
    }
    showModal.value = false;
    fetchUsers();
  } catch (e) {
    message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
  } finally {
    formLoading.value = false;
  }
}

async function handleDelete(id) {
  try {
    await deleteUser(id);
    message.success("删除成功");
    fetchUsers();
  } catch (e) {
    message.error("删除失败");
  }
}

async function handleResetPassword(id) {
  try {
    const res = await resetUserPassword(id);
    message.success(`密码已重置为: ${res.data.new_password || "TempP@ss123"}`);
  } catch (e) {
    message.error("重置密码失败");
  }
}

async function handleToggleStatus(row) {
  try {
    await lockUser(row.id, row.status === "active");
    message.success(row.status === "active" ? "用户已锁定" : "用户已解锁");
    fetchUsers();
  } catch (e) {
    message.error("操作失败");
  }
}

watch([page, searchKeyword, searchRole], fetchUsers);
onMounted(fetchUsers);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">用户管理</h1>
      <n-button type="primary" @click="handleAdd">
        <template #icon
          ><n-icon><AddOutline /></n-icon
        ></template>
        添加用户
      </n-button>
    </div>

    <n-card class="mb-4">
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索用户名/邮箱"
          clearable
          style="width: 300px"
          @keyup.enter="fetchUsers"
        />
        <n-select
          v-model:value="searchRole"
          :options="roleOptions"
          style="width: 150px"
        />
        <n-button @click="fetchUsers">
          <template #icon
            ><n-icon><RefreshOutline /></n-icon
          ></template>
          搜索
        </n-button>
      </div>
    </n-card>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="users"
        :loading="loading"
        :pagination="false"
        :bordered="true"
      />
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '添加用户' : '编辑用户'"
      style="width: 500px"
    >
      <n-form :model="currentUser" label-placement="left" label-width="80">
        <n-form-item label="用户名" required>
          <n-input
            v-model:value="currentUser.username"
            placeholder="请输入用户名"
          />
        </n-form-item>
        <n-form-item label="邮箱" required>
          <n-input v-model:value="currentUser.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="角色">
          <n-select
            v-model:value="currentUser.role"
            :options="[
              { label: '普通用户', value: 'user' },
              { label: '管理员', value: 'admin' },
            ]"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item v-if="modalMode === 'create'" label="密码" required>
          <n-input
            v-model:value="currentUser.password"
            type="password"
            placeholder="请输入密码"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="formLoading" @click="handleSubmit">
            {{ modalMode === "create" ? "创建" : "保存" }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>
