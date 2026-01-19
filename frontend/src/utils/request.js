import axios from "axios";
import { useUserStore } from "@/stores/user";
import { createDiscreteApi } from "naive-ui";

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "/api/v1",
  timeout: 30000,
});

// 在拦截器外使用 discrete API
const { message } = createDiscreteApi(["message"]);

request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore();
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

request.interceptors.response.use(
  (response) => {
    const { code, message: msg, data } = response.data;
    if (code === 200) {
      return data;
    }
    message.error(msg || "请求失败");
    return Promise.reject(new Error(msg));
  },
  (error) => {
    const { response } = error;
    if (response) {
      switch (response.status) {
        case 401:
          message.error("登录已过期，请重新登录");
          useUserStore().logout();
          window.location.href = "/login";
          break;
        case 403:
          message.error("没有权限访问");
          break;
        case 404:
          message.error("请求的资源不存在");
          break;
        case 500:
          message.error("服务器错误");
          break;
        default:
          message.error(response.data?.message || "请求失败");
      }
    } else {
      message.error("网络连接失败");
    }
    return Promise.reject(error);
  },
);

export default request;

export const $get = (url, params, config) =>
  request.get(url, { params, ...config });
export const $post = (url, data, config) => request.post(url, data, config);
export const $put = (url, data, config) => request.put(url, data, config);
export const $delete = (url, params, config) =>
  request.delete(url, { params, ...config });
