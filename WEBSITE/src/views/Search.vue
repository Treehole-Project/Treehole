<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="bg-white border-b">
      <div class="max-w-4xl mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">搜索</h1>
        <p class="text-gray-600">搜索你感兴趣的内容</p>
      </div>
    </div>

    <!-- 主要内容 -->
    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- 基础搜索 -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">基础搜索</h2>
        <div class="flex space-x-4">
          <input
            v-model="basicQuery"
            @keyup.enter="handleBasicSearch"
            type="text"
            placeholder="输入关键词搜索..."
            class="flex-1 input-field"
          />
          <button
            @click="handleBasicSearch"
            :disabled="!basicQuery.trim()"
            class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            搜索
          </button>
        </div>
      </div>

      <!-- 高级搜索 -->
      <div class="bg-white rounded-lg shadow-md p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-lg font-semibold text-gray-900">高级搜索</h2>
          <button
            @click="toggleAdvancedSearch"
            class="text-blue-600 hover:text-blue-700 text-sm font-medium"
          >
            {{ showAdvanced ? "收起" : "展开" }}
          </button>
        </div>

        <div v-show="showAdvanced" class="space-y-6">
          <!-- 搜索条件 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 标题搜索 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                标题关键词
              </label>
              <input
                v-model="advancedQuery.title"
                type="text"
                placeholder="搜索标题中的关键词"
                class="input-field"
              />
            </div>

            <!-- 内容搜索 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                内容关键词
              </label>
              <input
                v-model="advancedQuery.content"
                type="text"
                placeholder="搜索内容中的关键词"
                class="input-field"
              />
            </div>

            <!-- 作者搜索 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                作者
              </label>
              <input
                v-model="advancedQuery.author"
                type="text"
                placeholder="搜索作者用户名"
                class="input-field"
              />
            </div>

            <!-- 标签搜索 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                标签
              </label>
              <select v-model="advancedQuery.tag" class="input-field">
                <option value="">请选择标签</option>
                <option v-for="tag in tags" :key="tag" :value="tag">
                  {{ tag }}
                </option>
              </select>
            </div>

            <!-- 帖子ID -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                帖子ID
              </label>
              <input
                v-model="advancedQuery.post_id"
                type="text"
                placeholder="精确搜索帖子ID"
                class="input-field"
              />
            </div>

            <!-- 评论搜索 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                评论内容
              </label>
              <input
                v-model="advancedQuery.comment"
                type="text"
                placeholder="搜索评论中的关键词"
                class="input-field"
              />
            </div>
          </div>

          <!-- 逻辑关系选择 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              搜索逻辑
            </label>
            <div class="flex space-x-4">
              <label class="flex items-center">
                <input
                  v-model="advancedQuery.logic"
                  type="radio"
                  value="and"
                  class="text-blue-600 focus:ring-blue-500"
                />
                <span class="ml-2 text-sm text-gray-700"
                  >AND（所有条件都满足）</span
                >
              </label>
              <label class="flex items-center">
                <input
                  v-model="advancedQuery.logic"
                  type="radio"
                  value="or"
                  class="text-blue-600 focus:ring-blue-500"
                />
                <span class="ml-2 text-sm text-gray-700"
                  >OR（任一条件满足）</span
                >
              </label>
            </div>
          </div>

          <!-- 搜索按钮 -->
          <div class="flex justify-end space-x-3">
            <button @click="clearAdvancedSearch" class="btn-secondary">
              清空
            </button>
            <button
              @click="handleAdvancedSearch"
              :disabled="!hasAdvancedQuery"
              class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              高级搜索
            </button>
          </div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div
        v-if="error"
        class="mt-6 bg-red-50 border border-red-200 rounded-lg p-4"
      >
        <div class="flex">
          <svg
            class="h-5 w-5 text-red-400 mr-2"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd"
            />
          </svg>
          <p class="text-red-800">{{ error }}</p>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import Footer from "../components/Footer.vue";
import { postService } from "../services/postService.js";

const router = useRouter();

// 响应式数据
const basicQuery = ref("");
const showAdvanced = ref(true);
const tags = ref([]);
const error = ref("");

const advancedQuery = ref({
  title: "",
  content: "",
  author: "",
  tag: "",
  post_id: "",
  comment: "",
  logic: "and",
});

// 计算属性 - 检查是否有高级搜索条件
const hasAdvancedQuery = computed(() => {
  return Object.keys(advancedQuery.value).some(
    (key) => key !== "logic" && advancedQuery.value[key].trim()
  );
});

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await postService.getTags();
    tags.value = response.tags || [];
  } catch (err) {
    console.error("Error fetching tags:", err);
  }
};

// 切换高级搜索显示
const toggleAdvancedSearch = () => {
  showAdvanced.value = !showAdvanced.value;
};

// 处理基础搜索
const handleBasicSearch = () => {
  if (basicQuery.value.trim()) {
    router.push({
      name: "SearchResults",
      query: { q: basicQuery.value.trim() },
    });
  }
};

// 处理高级搜索
const handleAdvancedSearch = () => {
  if (!hasAdvancedQuery.value) return;

  const query = {};

  // 构建查询参数
  Object.keys(advancedQuery.value).forEach((key) => {
    const value = advancedQuery.value[key];
    if (value && value.trim()) {
      query[key] = value.trim();
    }
  });

  // 添加高级搜索标识
  query.advanced = "true";

  router.push({
    name: "SearchResults",
    query,
  });
};

// 清空高级搜索
const clearAdvancedSearch = () => {
  advancedQuery.value = {
    title: "",
    content: "",
    author: "",
    tag: "",
    post_id: "",
    comment: "",
    logic: "and",
  };
};

// 组件挂载时获取标签
onMounted(() => {
  fetchTags();
});
</script>

<style scoped>
/* 修复图标大小 */
svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}
</style>
