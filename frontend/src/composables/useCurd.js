import { ref, reactive, onMounted, watch } from "vue";
import { useMessage } from "naive-ui";

export function useCurd({
  fetchList,
  createItem,
  updateItem,
  deleteItem,
  defaultForm = {},
  beforeSubmit = (data) => data,
  afterFetch = (data) => data,
}) {
  const message = useMessage();
  
  // 列表数据
  const list = ref([]);
  const loading = ref(false);
  const total = ref(0);
  const page = ref(1);
  const size = ref(10);
  
  // 搜索参数
  const searchParams = reactive({});
  
  // 弹窗表单
  const showModal = ref(false);
  const modalMode = ref("create"); // 'create' | 'edit'
  const submitting = ref(false);
  const form = reactive({ ...defaultForm });
  const formRef = ref(null);

  // 获取列表
  const fetchData = async () => {
    loading.value = true;
    try {
      const params = {
        page: page.value,
        size: size.value,
        ...searchParams,
      };
      const res = await fetchList(params);
      const data = afterFetch(res);
      list.value = data.list || [];
      total.value = data.pagination?.total || data.total || 0;
    } catch (e) {
      message.error("获取数据失败");
      console.error(e);
    } finally {
      loading.value = false;
    }
  };

  // 搜索
  const handleSearch = () => {
    page.value = 1;
    fetchData();
  };

  // 重置搜索
  const handleReset = () => {
    Object.keys(searchParams).forEach((key) => {
      searchParams[key] = null;
    });
    handleSearch();
  };

  // 新增
  const handleAdd = () => {
    modalMode.value = "create";
    Object.assign(form, defaultForm);
    showModal.value = true;
  };

  // 编辑
  const handleEdit = (row) => {
    modalMode.value = "edit";
    Object.assign(form, row);
    showModal.value = true;
  };

  // 提交
  const handleSubmit = async () => {
    formRef.value?.validate(async (errors) => {
      if (errors) return;
      
      submitting.value = true;
      try {
        const submitData = beforeSubmit({ ...form });
        if (modalMode.value === "create") {
          await createItem(submitData);
          message.success("创建成功");
        } else {
          await updateItem(form.id, submitData);
          message.success("更新成功");
        }
        showModal.value = false;
        fetchData();
      } catch (e) {
        message.error(modalMode.value === "create" ? "创建失败" : "更新失败");
      } finally {
        submitting.value = false;
      }
    });
  };

  // 删除
  const handleDelete = async (id) => {
    try {
      await deleteItem(id);
      message.success("删除成功");
      fetchData();
    } catch (e) {
      message.error("删除失败");
    }
  };

  // 监听分页变化
  watch([page, size], () => {
    fetchData();
  });

  onMounted(fetchData);

  return {
    list,
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
    handleReset,
    handleAdd,
    handleEdit,
    handleSubmit,
    handleDelete,
  };
}
