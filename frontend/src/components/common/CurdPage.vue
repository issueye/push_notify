<script setup>
import { NCard, NDataTable, NPagination, NButton, NSpace, NIcon, NModal } from "naive-ui";
import { AddOutline, SearchOutline } from "@vicons/ionicons5";

const props = defineProps({
  title: String,
  loading: Boolean,
  columns: Array,
  data: Array,
  page: Number,
  pageSize: Number,
  itemCount: Number,
  showModal: Boolean,
  modalTitle: String,
  submitting: Boolean,
  scrollX: {
    type: Number,
    default: 1200
  }
});

const emit = defineEmits(["update:page", "update:pageSize", "update:showModal", "search", "add", "submit"]);

const handlePageChange = (p) => {
  emit("update:page", p);
};

const handlePageSizeChange = (s) => {
  emit("update:pageSize", s);
};

const handleCloseModal = () => {
  emit("update:showModal", false);
};
</script>

<template>
  <div class="curd-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">{{ title }}</h1>
      <n-button type="primary" @click="$emit('add')">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        添加{{ title }}
      </n-button>
    </div>

    <!-- 搜索栏 -->
    <n-card class="mb-4" v-if="$slots.search">
      <div class="flex gap-4 items-center">
        <slot name="search"></slot>
        <n-button type="primary" @click="$emit('search')">
          <template #icon>
            <n-icon><SearchOutline /></n-icon>
          </template>
          搜索
        </n-button>
        <slot name="search-actions"></slot>
      </div>
    </n-card>

    <!-- 表格 -->
    <n-card>
      <n-data-table
        :loading="loading"
        :columns="columns"
        :data="data"
        :pagination="false"
        :bordered="true"
        :scroll-x="scrollX"
      />
      <div class="mt-4 flex justify-end">
        <n-pagination
          :page="page"
          :page-size="pageSize"
          :item-count="itemCount"
          show-size-picker
          :page-sizes="[10, 20, 50, 100]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- 弹窗 -->
    <n-modal
      :show="showModal"
      preset="card"
      :title="modalTitle"
      style="width: 600px"
      @update:show="handleCloseModal"
    >
      <slot name="form"></slot>
      <template #footer>
        <div class="flex justify-end gap-2">
          <n-button @click="handleCloseModal">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="$emit('submit')">
            确定
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.curd-page {
  padding: 0;
}
</style>
