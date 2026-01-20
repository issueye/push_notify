import { createRouter, createWebHashHistory } from "vue-router";
import { useUserStore } from "@/stores/user";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/index.vue"),
    meta: { title: "登录", public: true },
  },
  {
    path: "/",
    component: () => import("@/components/layout/MainLayout.vue"),
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/dashboard/index.vue"),
        meta: { title: "工作台" },
      },
      {
        path: "repos",
        name: "Repos",
        component: () => import("@/views/repos/index.vue"),
        meta: { title: "仓库管理" },
      },
      {
        path: "targets",
        name: "Targets",
        component: () => import("@/views/targets/index.vue"),
        meta: { title: "推送目标" },
      },
      {
        path: "pushes",
        name: "Pushes",
        component: () => import("@/views/pushes/index.vue"),
        meta: { title: "推送记录" },
      },
      {
        path: "templates",
        name: "Templates",
        component: () => import("@/views/templates/index.vue"),
        meta: { title: "消息模板" },
      },
      {
        path: "prompts",
        name: "Prompts",
        component: () => import("@/views/prompts/index.vue"),
        meta: { title: "提示词" },
      },
      {
        path: "models",
        name: "Models",
        component: () => import("@/views/models/index.vue"),
        meta: { title: "AI模型" },
      },
      {
        path: "users",
        name: "Users",
        component: () => import("@/views/users/index.vue"),
        meta: { title: "用户管理", roles: ["admin"] },
      },
      {
        path: "logs/system",
        name: "SystemLog",
        component: () => import("@/views/logs/SystemLog.vue"),
        meta: { title: "系统日志" },
      },
      {
        path: "logs/operations",
        name: "OperationLog",
        component: () => import("@/views/logs/OperationLog.vue"),
        meta: { title: "操作日志" },
      },
      {
        path: "logs/ai-calls",
        name: "AICallLog",
        component: () => import("@/views/logs/AICallLog.vue"),
        meta: { title: "AI调用日志" },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/settings/index.vue"),
        meta: { title: "个人设置" },
      },
    ],
  },
  {
    path: "/403",
    name: "403",
    component: () => import("@/views/error/403.vue"),
    meta: { title: "无权限", public: true },
  },
  {
    path: "/404",
    name: "404",
    component: () => import("@/views/error/404.vue"),
    meta: { title: "页面不存在", public: true },
  },
  { path: "/:pathMatch(.*)*", redirect: "/404" },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title
    ? `${to.meta.title} - Push Notify`
    : "Push Notify";
  const userStore = useUserStore();

  if (to.meta.public) {
    next();
    return;
  }

  if (!userStore.isLoggedIn) {
    next({ name: "Login", query: { redirect: to.fullPath } });
    return;
  }

  if (!userStore.userInfo) {
    try {
      await userStore.getUserProfile();
    } catch {
      next({ name: "Login" });
      return;
    }
  }

  if (to.meta.roles && !to.meta.roles.includes(userStore.roles[0])) {
    next({ name: "403" });
    return;
  }

  next();
});

export default router;
