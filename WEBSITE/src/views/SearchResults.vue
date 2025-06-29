<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="bg-white border-b">
      <div class="max-w-4xl mx-auto px-4 py-8">
        <div class="flex items-center justify-between">
          <div>
            <button
              @click="goBack"
              class="flex items-center text-gray-600 hover:text-gray-900 mb-4"
            >
              <svg
                class="w-5 h-5 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 19l-7-7 7-7"
                />
              </svg>
              返回
            </button>
            <h1 class="text-3xl font-bold text-gray-900 mb-2">搜索结果</h1>
            <p class="text-gray-600">
              <template v-if="isAdvancedSearch"> 高级搜索结果 </template>
              <template v-else> "{{ searchQuery }}" 的搜索结果 </template>
            </p>
          </div>

          <!-- 搜索框 -->
          <div class="hidden md:block w-80">
            <div class="relative">
              <input
                v-model="newSearchQuery"
                @keyup.enter="handleNewSearch"
                type="text"
                placeholder="搜索帖子..."
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <div
                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
              >
                <svg
                  class="h-5 w-5 text-gray-400"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- 搜索条件显示 -->
      <div
        v-if="isAdvancedSearch && Object.keys(searchParams).length > 0"
        class="bg-white rounded-lg shadow-md p-4 mb-6"
      >
        <h3 class="text-sm font-medium text-gray-900 mb-3">搜索条件：</h3>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="(value, key) in displaySearchParams"
            :key="key"
            class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
          >
            {{ getFieldName(key) }}: {{ value }}
          </span>
        </div>
      </div>

      <!-- 错误提示 -->
      <div
        v-if="error"
        class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4"
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

      <!-- 搜索结果 -->
      <div class="space-y-6">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center py-12">
          <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"
          ></div>
        </div>

        <!-- 搜索结果列表 -->
        <template v-else-if="posts.length > 0">
          <PostCard v-for="post in posts" :key="post.id" :post="post" />
        </template>

        <!-- 空状态 -->
        <div v-else class="text-center py-12">
          <svg
            class="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">未找到相关结果</h3>
          <p class="mt-1 text-sm text-gray-500">
            尝试使用不同的关键词或
            <router-link to="/search" class="text-blue-600 hover:text-blue-700">
              高级搜索
            </router-link>
          </p>
        </div>
      </div>

      <!-- 分页组件 -->
      <div class="mt-8">
        <Pagination
          :pagination="pagination"
          :loading="loading"
          @page-change="handlePageChange"
        />
      </div>

      <!-- 相关操作 -->
      <div class="mt-8 flex justify-center space-x-4">
        <router-link to="/search" class="btn-secondary"> 新搜索 </router-link>
        <router-link to="/" class="btn-primary"> 返回首页 </router-link>
      </div>
    </main>

    <!-- Footer -->
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import PostCard from "../components/PostCard.vue";
import Pagination from "../components/Pagination.vue";
import Footer from "../components/Footer.vue";
import { postService } from "../services/postService.js";

const route = useRoute();
const router = useRouter();

// 响应式数据
const posts = ref([]);
const loading = ref(false);
const error = ref("");
const pagination = ref(null);
const newSearchQuery = ref("");
const searchParams = ref({});

// 计算属性
const searchQuery = computed(() => route.query.q || "");
const isAdvancedSearch = computed(() => route.query.advanced === "true");

const displaySearchParams = computed(() => {
  const params = { ...searchParams.value };
  delete params.advanced;
  delete params.page;
  delete params.limit;
  return Object.fromEntries(
    Object.entries(params).filter(([key, value]) => value && value.trim())
  );
});

// 字段名映射
const getFieldName = (key) => {
  const fieldNames = {
    title: "标题",
    content: "内容",
    author: "作者",
    tag: "标签",
    post_id: "帖子ID",
    comment: "评论",
    logic: "逻辑",
  };
  return fieldNames[key] || key;
};

// 执行搜索
const performSearch = async (page = 1) => {
  try {
    loading.value = true;
    error.value = "";

    let response;

    if (isAdvancedSearch.value) {
      // 高级搜索
      const params = { ...searchParams.value, page, limit: 20 };
      delete params.advanced;
      response = await postService.advancedSearch(params);
    } else {
      // 基础搜索
      response = await postService.searchPosts(searchQuery.value, page, 20);
    }

    posts.value = response.posts || [];
    pagination.value = response.pagination;
  } catch (err) {
    error.value = "搜索失败，请稍后重试";
    console.error("Error performing search:", err);
  } finally {
    loading.value = false;
  }
};

// 处理分页变化
const handlePageChange = (page) => {
  const query = { ...route.query, page };
  router.push({ name: "SearchResults", query });
};

// 处理新搜索
const handleNewSearch = () => {
  if (newSearchQuery.value.trim()) {
    router.push({
      name: "SearchResults",
      query: { q: newSearchQuery.value.trim() },
    });
  }
};

// 返回上一页
const goBack = () => {
  router.go(-1);
};

// 监听路由变化
watch(
  () => route.query,
  (newQuery) => {
    searchParams.value = { ...newQuery };
    performSearch(parseInt(newQuery.page) || 1);
  },
  { immediate: true }
);

// 组件挂载时初始化
onMounted(() => {
  searchParams.value = { ...route.query };
  newSearchQuery.value = searchQuery.value;
});
</script>

<style scoped>
/* 修复图标大小 */
svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.animate-spin {
  width: 3rem;
  height: 3rem;
}
</style>
