<script setup>
import { computed, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  NLayout,
  NLayoutSider,
  NLayoutContent,
  NMenu,
  NIcon,
  NBadge,
} from "naive-ui";
import { useAppStore } from "@/stores/app";
import { useUserStore } from "@/stores/user";
import {
  HomeOutline,
  GitBranchOutline,
  NotificationsOutline,
  SendOutline,
  DocumentTextOutline,
  BulbOutline,
  HardwareChipOutline,
  PeopleOutline,
  FileTrayFullOutline,
  SettingsOutline,
  LogOutOutline,
} from "@vicons/ionicons5";

const router = useRouter();
const route = useRoute();
const appStore = useAppStore();
const userStore = useUserStore();

const collapsed = computed(() => appStore.sidebarCollapsed);

const menuOptions = [
  {
    label: "工作台",
    key: "/dashboard",
    icon: () => h(NIcon, null, { default: () => h(HomeOutline) }),
  },
  {
    label: "仓库管理",
    key: "/repos",
    icon: () => h(NIcon, null, { default: () => h(GitBranchOutline) }),
  },
  {
    label: "推送目标",
    key: "/targets",
    icon: () => h(NIcon, null, { default: () => h(NotificationsOutline) }),
  },
  {
    label: "推送记录",
    key: "/pushes",
    icon: () => h(NIcon, null, { default: () => h(SendOutline) }),
  },
  {
    label: "消息模板",
    key: "/templates",
    icon: () => h(NIcon, null, { default: () => h(DocumentTextOutline) }),
  },
  {
    label: "提示词",
    key: "/prompts",
    icon: () => h(NIcon, null, { default: () => h(BulbOutline) }),
  },
  {
    label: "AI模型",
    key: "/models",
    icon: () => h(NIcon, null, { default: () => h(HardwareChipOutline) }),
  },
  ...(userStore.isAdmin
    ? [
        {
          label: "用户管理",
          key: "/users",
          icon: () => h(NIcon, null, { default: () => h(PeopleOutline) }),
        },
      ]
    : []),
  {
    label: "日志管理",
    key: "logs",
    icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }),
    children: [
      { label: "系统日志", key: "/logs/system" },
      { label: "操作日志", key: "/logs/operations" },
      { label: "AI调用日志", key: "/logs/ai-calls" },
    ],
  },
  {
    label: "个人设置",
    key: "/settings",
    icon: () => h(NIcon, null, { default: () => h(SettingsOutline) }),
  },
];

function handleMenuUpdate(key, keyPath) {
  router.push(key);
}

function handleLogout() {
  userStore.logout();
}

function toggleSidebar() {
  appStore.toggleSidebar();
}
</script>

<template>
  <n-layout has-sider class="h-screen">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger
      @collapse="toggleSidebar"
      @expand="toggleSidebar"
      :native-scrollbar="false"
    >
      <div
        class="h-16 flex items-center justify-center border-b border-gray-200 gap-2"
      >
        <img src="@/assets/logo.svg" class="w-8 h-8" alt="Logo" />
        <h1 v-if="!collapsed" class="text-lg font-bold text-primary-600">
          Push Notify
        </h1>
      </div>
      <n-menu
        :options="menuOptions"
        :value="route.path"
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        @update:value="handleMenuUpdate"
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header
        bordered
        class="h-16 flex items-center justify-between px-6"
      >
        <div class="flex items-center">
          <span class="text-gray-600"
            >欢迎，{{ userStore.userInfo?.username || "用户" }}</span
          >
        </div>
        <div class="flex items-center gap-4">
          <n-button quaternary circle @click="handleLogout">
            <template #icon>
              <n-icon><LogOutOutline /></n-icon>
            </template>
          </n-button>
        </div>
      </n-layout-header>
      <n-layout-content
        class="p-6 bg-gray-50"
        style="height: calc(100% - 65px)"
        :native-scrollbar="false"
      >
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>
