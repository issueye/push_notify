<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getPushDetail } from "@/services/push";
import { NCard, NButton, NIcon, NResult, NSpin, NScrollbar, NTag, NSpace, NDivider } from "naive-ui";
import { ArrowBackOutline } from "@vicons/ionicons5";
import MarkdownIt from "markdown-it";
import { formatDate } from "@/utils/date";

const route = useRoute();
const router = useRouter();
const md = new MarkdownIt();

const pushId = route.query.id;
const loading = ref(true);
const pushDetail = ref(null);
const error = ref(null);

async function fetchDetail() {
  if (!pushId) {
    error.value = "未提供记录 ID";
    loading.value = false;
    return;
  }

  loading.value = true;
  try {
    const res = await getPushDetail(pushId);
    pushDetail.value = res.data || res;
  } catch (err) {
    error.value = "获取详情失败：" + (err.message || "未知错误");
  } finally {
    loading.value = false;
  }
}

function handleBack() {
  router.back();
}

onMounted(fetchDetail);
</script>

<template>
  <div class="review-detail">
    <div class="mb-6 flex items-center gap-4">
      <n-button quaternary circle @click="handleBack">
        <template #icon>
          <n-icon><ArrowBackOutline /></n-icon>
        </template>
      </n-button>
      <h1 class="text-2xl font-bold">代码审查详情</h1>
    </div>

    <n-spin :show="loading">
      <div v-if="error" class="mt-20">
        <n-result status="error" title="加载失败" :description="error">
          <template #footer>
            <n-button @click="handleBack">返回列表</n-button>
          </template>
        </n-result>
      </div>

      <div v-else-if="pushDetail" class="flex flex-col gap-4">
        <n-card title="基本信息">
          <n-space vertical size="large">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="flex flex-col gap-1">
                <span class="text-gray-500 text-sm">仓库名称</span>
                <span class="font-medium">{{ pushDetail.repo?.name || '-' }}</span>
              </div>
              <div class="flex flex-col gap-1">
                <span class="text-gray-500 text-sm">分支</span>
                <n-tag size="small" type="info">{{ pushDetail.branch || 'main' }}</n-tag>
              </div>
              <div class="flex flex-col gap-1">
                <span class="text-gray-500 text-sm">提交 ID</span>
                <span class="font-mono text-xs">{{ pushDetail.commit_id }}</span>
              </div>
              <div class="flex flex-col gap-1">
                <span class="text-gray-500 text-sm">审查状态</span>
                <n-tag :type="pushDetail.codeview_status === 'success' ? 'success' : 'warning'">
                  {{ pushDetail.codeview_status }}
                </n-tag>
              </div>
            </div>
            <n-divider style="margin: 8px 0" />
            <div class="flex flex-col gap-1">
              <span class="text-gray-500 text-sm">提交信息</span>
              <p class="mt-1">{{ pushDetail.commit_msg }}</p>
            </div>
          </n-space>
        </n-card>

        <n-card title="审查建议" :segmented="{ content: true }">
          <div v-if="pushDetail.codeview_result" class="markdown-body">
            <div v-html="md.render(pushDetail.codeview_result)"></div>
          </div>
          <n-result
            v-else
            status="info"
            title="暂无结果"
            description="AI 尚未完成对本次提交的审查，或审查已跳过。"
          />
        </n-card>
      </div>
    </n-spin>
  </div>
</template>

<style scoped>
.review-detail {
  max-width: 1000px;
  margin: 0 auto;
}

.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
  font-size: 16px;
  line-height: 1.6;
  word-wrap: break-word;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body :deep(h1) { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
.markdown-body :deep(h2) { font-size: 1.25em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
.markdown-body :deep(h3) { font-size: 1.1em; }

.markdown-body :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
}

.markdown-body :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.markdown-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.markdown-body :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e1;
  margin: 0 0 16px 0;
}
</style>
