<script setup>
import { ref, reactive, onMounted } from "vue";
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NRadioGroup,
  NRadio,
  NSwitch,
  NSpace,
  NDivider,
} from "naive-ui";
import { useUserStore } from "@/stores/user";
import { useAppStore } from "@/stores/app";
import { useMessage } from "@/composables/useMessage";

const message = useMessage();
const userStore = useUserStore();
const appStore = useAppStore();

const form = reactive({
  email: "",
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
});

const notifySettings = reactive({
  channels: ["dingtalk", "email"],
  quietHours: false,
});

const themeOptions = [
  { label: "浅色", value: "light" },
  { label: "深色", value: "dark" },
];

const theme = ref(appStore.theme);

function handleUpdateEmail() {
  message.success("邮箱更新成功");
}

function handleChangePassword() {
  if (form.newPassword !== form.confirmPassword) {
    message.error("两次密码不一致");
    return;
  }
  if (form.newPassword.length < 8) {
    message.error("密码长度至少8位");
    return;
  }
  userStore.changePassword(form.oldPassword, form.newPassword).then(() => {
    message.success("密码修改成功");
    form.oldPassword = "";
    form.newPassword = "";
    form.confirmPassword = "";
  });
}

function handleThemeChange(val) {
  appStore.setTheme(val);
}

onMounted(() => {
  if (userStore.userInfo) {
    form.email = userStore.userInfo.email;
  }
});
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">个人设置</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <n-card title="基本信息">
        <n-form label-placement="left" label-width="100">
          <n-form-item label="用户名">
            <n-input :value="userStore.userInfo?.username" disabled />
          </n-form-item>
          <n-form-item label="角色">
            <n-input
              :value="
                userStore.userInfo?.role === 'admin' ? '管理员' : '普通用户'
              "
              disabled
            />
          </n-form-item>
          <n-form-item label="邮箱">
            <n-input v-model:value="form.email" placeholder="请输入邮箱" />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="handleUpdateEmail"
              >更新邮箱</n-button
            >
          </n-form-item>
        </n-form>
      </n-card>

      <n-card title="修改密码">
        <n-form label-placement="left" label-width="100">
          <n-form-item label="原密码">
            <n-input
              v-model:value="form.oldPassword"
              type="password"
              placeholder="请输入原密码"
              show-password-on="click"
            />
          </n-form-item>
          <n-form-item label="新密码">
            <n-input
              v-model:value="form.newPassword"
              type="password"
              placeholder="请输入新密码"
              show-password-on="click"
            />
          </n-form-item>
          <n-form-item label="确认密码">
            <n-input
              v-model:value="form.confirmPassword"
              type="password"
              placeholder="请确认新密码"
              show-password-on="click"
            />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="handleChangePassword"
              >修改密码</n-button
            >
          </n-form-item>
        </n-form>
      </n-card>

      <n-card title="外观设置">
        <n-form label-placement="left" label-width="100">
          <n-form-item label="主题">
            <n-radio-group
              v-model:value="theme"
              @update:value="handleThemeChange"
            >
              <n-radio value="light">浅色</n-radio>
              <n-radio value="dark">深色</n-radio>
            </n-radio-group>
          </n-form-item>
        </n-form>
      </n-card>

      <n-card title="通知设置">
        <n-form label-placement="left" label-width="120">
          <n-form-item label="通知渠道">
            <n-checkbox-group v-model:value="notifySettings.channels">
              <n-space>
                <n-checkbox value="dingtalk">钉钉</n-checkbox>
                <n-checkbox value="email">邮箱</n-checkbox>
              </n-space>
            </n-checkbox-group>
          </n-form-item>
          <n-form-item label="免打扰模式">
            <n-switch v-model:value="notifySettings.quietHours" />
          </n-form-item>
        </n-form>
      </n-card>
    </div>
  </div>
</template>
