import { ref } from "vue";
import { useMessage as naiveMessage } from "naive-ui";

export function useMessage() {
  const message = naiveMessage();
  return {
    success: (content) => message.success(content),
    error: (content) => message.error(content),
    warning: (content) => message.warning(content),
    info: (content) => message.info(content),
  };
}

export function useConfirm() {
  const message = naiveMessage();
  return {
    confirm: (content, onOk, onCancel) => {
      message.warning(content, {
        positiveText: "确认",
        negativeText: "取消",
        onPositiveClick: onOk,
        onNegativeClick: onCancel,
      });
    },
  };
}

export function useLoading() {
  const loading = ref(false);
  return {
    loading,
    start: () => {
      loading.value = true;
    },
    stop: () => {
      loading.value = false;
    },
  };
}

export function usePagination() {
  const page = ref(1);
  const size = ref(10);
  const total = ref(0);

  return {
    page,
    size,
    total,
    setPage: (p) => {
      page.value = p;
    },
    setSize: (s) => {
      size.value = s;
    },
    reset: () => {
      page.value = 1;
      total.value = 0;
    },
  };
}
