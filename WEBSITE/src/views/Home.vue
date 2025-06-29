<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="bg-white border-b">
      <div class="max-w-4xl mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Treehole Project</h1>
        <p class="text-gray-600">
          完全开源的匿名化交流平台,由某不知名BIT学生创建.
        </p>
      </div>
    </div>

    <!-- 主要内容 -->
    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- 搜索栏 -->
      <div class="mb-8">
        <div class="relative">
          <input
            v-model="searchQuery"
            @keyup.enter="handleSearch"
            type="text"
            placeholder="搜索帖子..."
            class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
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
          <button
            @click="handleSearch"
            class="absolute inset-y-0 right-0 pr-3 flex items-center"
          >
            <span
              class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
            >
              搜索
            </span>
          </button>
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

      <!-- 发帖功能 -->
      <div class="mb-8">
        <div class="card p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900">发布新帖</h2>
            <button
              @click="togglePostForm"
              class="text-blue-600 hover:text-blue-700 text-sm font-medium transition-colors"
            >
              {{ showPostForm ? "收起" : "展开" }}
            </button>
          </div>

          <div v-show="showPostForm" class="space-y-4">
            <!-- 用户名输入 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                用户名
              </label>
              <input
                v-model="newPost.username"
                type="text"
                placeholder="请输入用户名"
                class="input-field"
              />
            </div>

            <!-- 标题输入 -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                标题
              </label>
              <input
                v-model="newPost.title"
                type="text"
                placeholder="请输入帖子标题"
                class="input-field"
              />
            </div>

            <!-- 内容输入 -->
            <div>
              <div class="flex justify-between items-center mb-2">
                <label class="block text-sm font-medium text-gray-700">
                  内容
                </label>
                <span class="text-xs text-gray-500">
                  {{ newPost.content.length }}/1000 字
                </span>
              </div>
              <textarea
                v-model="newPost.content"
                @keydown.ctrl.enter="submitPost"
                placeholder="分享你的想法... (Ctrl+Enter 快速发布)"
                rows="4"
                maxlength="1000"
                class="input-field resize-none"
              ></textarea>
            </div>

            <!-- 操作按钮 -->
            <div class="flex justify-end space-x-3">
              <button @click="clearPostForm" class="btn-secondary">清空</button>
              <button
                @click="submitPost"
                :disabled="!canSubmitPost || submittingPost"
                class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ submittingPost ? "发布中..." : "发布帖子" }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 成功提示 -->
      <div
        v-if="successMessage"
        class="mb-6 bg-green-50 border border-green-200 rounded-lg p-4"
      >
        <div class="flex">
          <svg
            class="h-5 w-5 text-green-400 mr-2"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clip-rule="evenodd"
            />
          </svg>
          <p class="text-green-800">{{ successMessage }}</p>
        </div>
      </div>

      <!-- 帖子列表 -->
      <div class="space-y-6">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center py-12">
          <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"
          ></div>
        </div>

        <!-- 帖子卡片 -->
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
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无帖子</h3>
          <p class="mt-1 text-sm text-gray-500">
            目前还没有任何帖子，快来发布第一个吧！
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
    </main>

    <!-- Footer -->
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import PostCard from "../components/PostCard.vue";
import Pagination from "../components/Pagination.vue";
import Footer from "../components/Footer.vue";
import { postService } from "../services/postService.js";

const router = useRouter();
const route = useRoute();

// 响应式数据
const posts = ref([]);
const loading = ref(false);
const error = ref("");
const pagination = ref(null);
const searchQuery = ref("");
const showPostForm = ref(false);
const newPost = ref({
  username: "",
  title: "",
  content: "",
});
const submittingPost = ref(false);
const successMessage = ref("");

// 获取帖子列表
const fetchPosts = async (page = 1) => {
  try {
    loading.value = true;
    error.value = "";
    const response = await postService.getPosts(page, 20);
    posts.value = response.posts || [];
    pagination.value = response.pagination;
    
    // 更新URL中的page参数
    if (page > 1) {
      router.replace({ 
        name: 'Home', 
        query: { 
          ...route.query, 
          page: page.toString() 
        } 
      });
    } else {
      // 如果是第一页，移除page参数
      const newQuery = { ...route.query };
      delete newQuery.page;
      router.replace({ 
        name: 'Home', 
        query: newQuery 
      });
    }
  } catch (err) {
    error.value = "获取帖子列表失败，请稍后重试";
    console.error("Error fetching posts:", err);
  } finally {
    loading.value = false;
  }
};

// 处理分页变化
const handlePageChange = (page) => {
  fetchPosts(page);
  // 滚动到顶部
  window.scrollTo({ top: 0, behavior: "smooth" });
};

// 处理搜索
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      name: "SearchResults",
      query: { q: searchQuery.value.trim() },
    });
  }
};

// 切换发帖表单显示
const togglePostForm = () => {
  showPostForm.value = !showPostForm.value;
};

// 清空发帖表单
const clearPostForm = () => {
  newPost.value = {
    username: "",
    title: "",
    content: "",
  };
  successMessage.value = "";
};

// 提交发帖
const submitPost = async () => {
  if (!canSubmitPost.value || submittingPost.value) return;

  submittingPost.value = true;
  successMessage.value = "";
  error.value = "";

  try {
    await postService.createPost(newPost.value);
    successMessage.value = "帖子发布成功！";
    clearPostForm();
    showPostForm.value = false; // 发布成功后收起表单

    // 重新获取第一页数据，显示新发布的帖子
    await fetchPosts(1);

    // 滚动到顶部
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (err) {
    error.value = "发布帖子失败，请稍后重试";
    console.error("Error submitting post:", err);
  } finally {
    submittingPost.value = false;
  }
};

// 计算属性 - 检查是否可以提交帖子
const canSubmitPost = computed(() => {
  return (
    newPost.value.username.trim() &&
    newPost.value.title.trim() &&
    newPost.value.content.trim()
  );
});

// 组件挂载时获取数据
onMounted(() => {
  // 从URL参数中获取页码，默认为1
  const currentPage = parseInt(route.query.page) || 1;
  fetchPosts(currentPage);
});

// 监听路由变化，处理浏览器前进后退
watch(() => route.query.page, (newPage) => {
  const page = parseInt(newPage) || 1;
  if (pagination.value && page !== pagination.value.page) {
    fetchPosts(page);
  }
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

/* 输入框样式 */
.input-field {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  background-color: #fff;
  font-size: 0.875rem;
  color: #111827;
  transition: border-color 0.2s;
}

.input-field:focus {
  border-color: #3b82f6;
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* 按钮样式 */
.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  background-color: #3b82f6;
  color: #fff;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover {
  background-color: #2563eb;
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  background-color: #f3f4f6;
  color: #111827;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-secondary:hover {
  background-color: #e5e7eb;
}

/* 卡片样式 */
.card {
  background-color: #fff;
  border-radius: 0.375rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
</style>
