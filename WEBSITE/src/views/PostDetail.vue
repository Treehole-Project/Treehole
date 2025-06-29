<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center min-h-screen">
      <div
        class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"
      ></div>
    </div>

    <!-- 主要内容 -->
    <div v-else>
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
              <h1
                class="text-3xl font-bold text-gray-900 mb-2"
                :title="post.title"
              >
                {{ truncateTitle(post.title) }}
              </h1>
              <div class="flex items-center space-x-4 text-sm text-gray-600">
                <span class="flex items-center">
                  <svg
                    class="w-4 h-4 mr-1"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                      clip-rule="evenodd"
                    />
                  </svg>
                  {{ post.author }}
                </span>
                <span
                  class="flex items-center"
                  :title="formatFullDateTime(post.created_at)"
                >
                  <svg
                    class="w-4 h-4 mr-1"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                      clip-rule="evenodd"
                    />
                  </svg>
                  {{ formatDate(post.created_at) }}
                </span>
                <span class="flex items-center">
                  <svg
                    class="w-4 h-4 mr-1"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                    <path
                      fill-rule="evenodd"
                      d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                  {{ post.view_count || 0 }} 次查看
                </span>
              </div>
            </div>

            <!-- 点赞按钮 -->
            <div class="flex items-center space-x-2">
              <button
                class="flex items-center px-4 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <svg
                  class="w-4 h-4 mr-2"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z"
                  />
                </svg>
                {{ post.like_num || 0 }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 帖子内容 -->
      <main class="max-w-4xl mx-auto px-4 py-8">
        <!-- 帖子正文 -->
        <div class="bg-white rounded-lg shadow-md p-8 mb-8">
          <div class="prose max-w-none">
            <p class="text-gray-800 leading-relaxed whitespace-pre-wrap">
              {{ post.content }}
            </p>
          </div>

          <!-- 帖子图片 -->
          <div v-if="postImages.length > 0" class="mt-4">
            <div
              v-for="(imageUrl, index) in postImages"
              :key="index"
              class="mb-3 last:mb-0"
            >
              <img
                :src="imageUrl"
                :alt="`帖子图片 ${index + 1}`"
                class="max-w-full h-auto rounded-lg cursor-pointer hover:opacity-90 transition-opacity shadow-md"
                @click="openImageViewer(imageUrl, `帖子图片 ${index + 1}`)"
                loading="lazy"
              />
            </div>
          </div>

          <!-- 标签 -->
          <div v-if="post.tag && post.tag !== '未分析'" class="mt-6">
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
            >
              {{ post.tag }}
            </span>
          </div>
        </div>

        <!-- 评论区 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-semibold text-gray-900">
              评论 ({{ repliesPagination?.total || 0 }})
            </h2>
          </div>

          <!-- 发表评论 -->
          <div class="mb-8 p-4 bg-gray-50 rounded-lg">
            <h3 class="text-sm font-medium text-gray-900 mb-3">发表评论</h3>
            <div class="space-y-3">
              <input
                v-model="newReply.username"
                type="text"
                placeholder="你的昵称"
                class="input-field"
                maxlength="20"
              />
              <textarea
                v-model="newReply.content"
                placeholder="写下你的评论..."
                rows="3"
                class="input-field resize-none"
                maxlength="500"
              ></textarea>
              <div class="flex justify-between items-center">
                <span class="text-xs text-gray-500"
                  >{{ newReply.content.length }}/500</span
                >
                <button
                  @click="submitReply"
                  :disabled="!canSubmitReply"
                  class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  发表评论
                </button>
              </div>
            </div>
          </div>

          <!-- 评论列表 -->
          <div class="space-y-6">
            <!-- 加载状态 -->
            <div v-if="repliesLoading" class="flex justify-center py-8">
              <div
                class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"
              ></div>
            </div>

            <!-- 评论项 -->
            <div v-else-if="replies.length > 0">
              <div
                v-for="reply in replies"
                :key="reply.id"
                class="border-l-4 pl-4 py-4"
                :class="{
                  'border-blue-400': reply.level === 2,
                  'border-orange-300': reply.isOrphan,
                  'border-gray-200': reply.level === 1 && !reply.isOrphan,
                }"
              >
                <!-- 一级评论或孤儿评论 -->
                <div
                  v-if="reply.level === 1 || reply.isOrphan"
                  class="space-y-3"
                >
                  <!-- 孤儿评论提示 -->
                  <div
                    v-if="reply.isOrphan"
                    class="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded mb-2"
                  >
                    💬 回复评论（原评论在其他页面）
                  </div>
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                      <div
                        class="w-8 h-8 rounded-full flex items-center justify-center"
                        :class="
                          reply.isOrphan ? 'bg-orange-200' : 'bg-gray-300'
                        "
                      >
                        <span
                          class="text-sm font-medium"
                          :class="
                            reply.isOrphan ? 'text-orange-700' : 'text-gray-600'
                          "
                        >
                          {{ reply.author.charAt(0).toUpperCase() }}
                        </span>
                      </div>
                      <div>
                        <p class="font-medium text-gray-900">
                          {{ reply.author }}
                        </p>
                        <p class="text-xs text-gray-500">
                          {{ formatDate(reply.created_at) }}
                        </p>
                      </div>
                    </div>
                    <div class="flex items-center space-x-3">
                      <button
                        v-if="!reply.isOrphan"
                        @click="toggleReplyForm(reply.id)"
                        class="text-sm text-blue-600 hover:text-blue-700"
                      >
                        回复
                      </button>
                      <span v-else class="text-xs text-gray-400">
                        （跨页回复）
                      </span>
                      <button
                        class="flex items-center text-sm text-gray-500 hover:text-gray-700"
                      >
                        <svg
                          class="w-4 h-4 mr-1"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path
                            d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z"
                          />
                        </svg>
                        {{ reply.like_num || 0 }}
                      </button>
                    </div>
                  </div>
                  <p class="text-gray-800 leading-relaxed">
                    {{ reply.content }}
                  </p>

                  <!-- 一级评论图片 -->
                  <div v-if="getReplyImages(reply).length > 0" class="mt-3">
                    <div
                      v-for="(imageUrl, index) in getReplyImages(reply)"
                      :key="index"
                      class="mb-2 last:mb-0"
                    >
                      <img
                        :src="imageUrl"
                        :alt="`评论图片 ${index + 1}`"
                        class="max-w-full h-auto rounded-lg cursor-pointer hover:opacity-90 transition-opacity shadow-sm"
                        @click="
                          openImageViewer(imageUrl, `评论图片 ${index + 1}`)
                        "
                        loading="lazy"
                        style="max-height: 300px"
                      />
                    </div>
                  </div>

                  <!-- 回复表单 -->
                  <div
                    v-if="showReplyForm === reply.id && !reply.isOrphan"
                    class="mt-4 p-4 bg-gray-50 rounded-lg"
                  >
                    <div class="space-y-3">
                      <input
                        v-model="subReply.username"
                        type="text"
                        placeholder="你的昵称"
                        class="input-field"
                        maxlength="20"
                      />
                      <textarea
                        v-model="subReply.content"
                        :placeholder="`回复 @${reply.author}:`"
                        rows="2"
                        class="input-field resize-none"
                        maxlength="500"
                      ></textarea>
                      <div class="flex justify-end space-x-2">
                        <button
                          @click="cancelReply"
                          class="btn-secondary text-sm"
                        >
                          取消
                        </button>
                        <button
                          @click="submitSubReply(reply.id)"
                          :disabled="!canSubmitSubReply"
                          class="btn-primary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          回复
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 二级评论 -->
                <div v-else class="ml-8 space-y-3">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                      <div
                        class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center"
                      >
                        <span class="text-xs font-medium text-blue-600">
                          {{ reply.author.charAt(0).toUpperCase() }}
                        </span>
                      </div>
                      <div>
                        <p class="text-sm font-medium text-gray-900">
                          {{ reply.author }}
                        </p>
                        <p class="text-xs text-gray-500">
                          {{ formatDate(reply.created_at) }}
                        </p>
                      </div>
                    </div>
                    <button
                      class="flex items-center text-sm text-gray-500 hover:text-gray-700"
                    >
                      <svg
                        class="w-4 h-4 mr-1"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z"
                        />
                      </svg>
                      {{ reply.like_num || 0 }}
                    </button>
                  </div>
                  <p class="text-sm text-gray-800 leading-relaxed">
                    {{ reply.content }}
                  </p>

                  <!-- 二级评论图片 -->
                  <div v-if="getReplyImages(reply).length > 0" class="mt-2">
                    <div
                      v-for="(imageUrl, index) in getReplyImages(reply)"
                      :key="index"
                      class="mb-2 last:mb-0"
                    >
                      <img
                        :src="imageUrl"
                        :alt="`回复图片 ${index + 1}`"
                        class="max-w-full h-auto rounded-lg cursor-pointer hover:opacity-90 transition-opacity shadow-sm"
                        @click="
                          openImageViewer(imageUrl, `回复图片 ${index + 1}`)
                        "
                        loading="lazy"
                        style="max-height: 250px"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else class="text-center py-8">
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
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">暂无评论</h3>
              <p class="mt-1 text-sm text-gray-500">成为第一个发表评论的人</p>
            </div>
          </div>

          <!-- 评论分页 -->
          <div class="mt-8">
            <Pagination
              :pagination="repliesPagination"
              :loading="repliesLoading"
              @page-change="handleRepliesPageChange"
            />
          </div>
        </div>
      </main>
    </div>

    <!-- 错误提示 -->
    <div
      v-if="error"
      class="fixed bottom-4 right-4 bg-red-500 text-white px-6 py-3 rounded-lg shadow-lg"
    >
      {{ error }}
    </div>

    <!-- Footer -->
    <Footer />

    <!-- 图片查看器 -->
    <ImageViewer
      :is-visible="imageViewer.isVisible"
      :image-url="imageViewer.imageUrl"
      :alt-text="imageViewer.altText"
      @close="closeImageViewer"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import Pagination from "../components/Pagination.vue";
import Footer from "../components/Footer.vue";
import ImageViewer from "../components/ImageViewer.vue";
import { postService } from "../services/postService.js";
import { formatDate, formatFullDateTime } from "../utils/dateUtils.js";

const route = useRoute();
const router = useRouter();

// 截断标题函数
const truncateTitle = (title) => {
  if (!title) return "";
  return title.length > 50 ? title.substring(0, 50) + "..." : title;
};

// 响应式数据
const post = ref({});
const replies = ref([]);
const loading = ref(true);
const repliesLoading = ref(false);
const error = ref("");
const repliesPagination = ref(null);
const showReplyForm = ref(null);

// 新评论数据
const newReply = ref({
  username: "",
  content: "",
});

// 子评论数据
const subReply = ref({
  username: "",
  content: "",
});

// 图片查看器状态
const imageViewer = ref({
  isVisible: false,
  imageUrl: "",
  altText: "",
});

// 计算属性
const canSubmitReply = computed(() => {
  return newReply.value.username.trim() && newReply.value.content.trim();
});

const canSubmitSubReply = computed(() => {
  return subReply.value.username.trim() && subReply.value.content.trim();
});

// 解析图片数组
const parseImages = (imagesStr) => {
  if (!imagesStr || imagesStr === "[]" || imagesStr === "") return [];

  console.log("Input imagesStr:", imagesStr); // 调试日志

  try {
    // 首先尝试作为JSON数组解析
    if (imagesStr.startsWith("[") && imagesStr.endsWith("]")) {
      const parsed = JSON.parse(imagesStr);
      const result = Array.isArray(parsed)
        ? parsed.filter((url) => url && url.trim())
        : [];
      console.log("JSON parsed result:", result);
      return result;
    }

    // 如果不是JSON格式，尝试按逗号分割
    if (imagesStr.includes(",")) {
      const result = imagesStr
        .split(",")
        .map((url) => url.trim())
        .filter((url) => url);
      console.log("Comma split result:", result);
      return result;
    }

    // 单个图片URL
    const result = [imagesStr.trim()].filter((url) => url);
    console.log("Single URL result:", result);
    return result;
  } catch (err) {
    console.error("Error parsing images:", err);

    // 如果JSON解析失败，尝试按逗号分割
    if (imagesStr.includes(",")) {
      const result = imagesStr
        .split(",")
        .map((url) => url.trim())
        .filter((url) => url);
      console.log("Fallback comma split result:", result);
      return result;
    }

    // 最后尝试作为单个URL
    const result = imagesStr.trim() ? [imagesStr.trim()] : [];
    console.log("Fallback single URL result:", result);
    return result;
  }
};

// 计算属性 - 解析帖子图片
const postImages = computed(() => {
  return parseImages(post.value.images);
});

// 获取回复的图片
const getReplyImages = (reply) => {
  return parseImages(reply.images);
};

// 打开图片查看器
const openImageViewer = (imageUrl, altText = "图片") => {
  imageViewer.value = {
    isVisible: true,
    imageUrl,
    altText,
  };
};

// 关闭图片查看器
const closeImageViewer = () => {
  imageViewer.value.isVisible = false;
};

// 获取帖子详情
const fetchPost = async () => {
  try {
    loading.value = true;
    const response = await postService.getPost(route.params.id);
    post.value = response;
  } catch (err) {
    error.value = "获取帖子详情失败";
    console.error("Error fetching post:", err);
  } finally {
    loading.value = false;
  }
};

// 将评论组织成层级结构
const organizeReplies = (repliesArray) => {
  const organized = [];
  const orphanReplies = []; // 存储孤儿评论（找不到父评论的子评论）
  const repliesMap = new Map();

  // 首先将所有评论按ID存储到Map中
  repliesArray.forEach((reply) => {
    repliesMap.set(reply.id, { ...reply, children: [] });
  });

  // 然后处理层级关系
  repliesArray.forEach((reply) => {
    if (reply.parent_id === 0) {
      // 一级评论，直接添加到结果数组
      organized.push(repliesMap.get(reply.id));
    } else {
      // 子评论，尝试添加到父评论的children数组中
      const parent = repliesMap.get(reply.parent_id);
      if (parent) {
        parent.children.push(repliesMap.get(reply.id));
      } else {
        // 找不到父评论，将其标记为孤儿评论
        const orphanReply = repliesMap.get(reply.id);
        orphanReply.isOrphan = true; // 添加孤儿标记
        orphanReply.originalParentId = reply.parent_id; // 保留原始父评论ID
        orphanReplies.push(orphanReply);
      }
    }
  });

  // 将结果扁平化，保持层级顺序
  const flattenReplies = (replies, result = []) => {
    replies.forEach((reply) => {
      result.push(reply);
      if (reply.children && reply.children.length > 0) {
        // 按创建时间排序子评论
        reply.children.sort(
          (a, b) => new Date(a.created_at) - new Date(b.created_at)
        );
        flattenReplies(reply.children, result);
      }
    });
    return result;
  };

  // 合并正常评论和孤儿评论
  const normalReplies = flattenReplies(organized);

  // 将孤儿评论按创建时间排序后添加到结果末尾
  const sortedOrphanReplies = orphanReplies.sort(
    (a, b) => new Date(a.created_at) - new Date(b.created_at)
  );

  return [...normalReplies, ...sortedOrphanReplies];
};

// 获取评论列表
const fetchReplies = async (page = 1) => {
  try {
    repliesLoading.value = true;
    const response = await postService.getPostReplies(
      route.params.id,
      page,
      20
    );

    // 组织评论数据为层级结构
    const organizedReplies = organizeReplies(response.replies || []);
    replies.value = organizedReplies;
    repliesPagination.value = response.pagination;
  } catch (err) {
    console.error("Error fetching replies:", err);
  } finally {
    repliesLoading.value = false;
  }
};

// 处理评论分页变化
const handleRepliesPageChange = (page) => {
  fetchReplies(page);
};

// 提交主评论
const submitReply = async () => {
  try {
    await postService.createReply(route.params.id, {
      content: newReply.value.content,
      username: newReply.value.username,
      parent_id: 0,
    });

    // 清空表单
    newReply.value = { username: "", content: "" };

    // 重新加载评论列表
    await fetchReplies();

    // 更新帖子的回复数量
    if (post.value.reply_count !== undefined) {
      post.value.reply_count++;
    }
  } catch (err) {
    error.value = "发表评论失败，请稍后重试";
    console.error("Error submitting reply:", err);
  }
};

// 提交子评论
const submitSubReply = async (parentId) => {
  try {
    await postService.createReply(route.params.id, {
      content: subReply.value.content,
      username: subReply.value.username,
      parent_id: parentId,
    });

    // 清空表单并关闭回复框
    subReply.value = { username: "", content: "" };
    showReplyForm.value = null;

    // 重新加载评论列表
    await fetchReplies();

    // 更新帖子的回复数量
    if (post.value.reply_count !== undefined) {
      post.value.reply_count++;
    }
  } catch (err) {
    error.value = "发表回复失败，请稍后重试";
    console.error("Error submitting sub reply:", err);
  }
};

// 切换回复表单
const toggleReplyForm = (replyId) => {
  showReplyForm.value = showReplyForm.value === replyId ? null : replyId;
  if (showReplyForm.value === null) {
    subReply.value = { username: "", content: "" };
  }
};

// 取消回复
const cancelReply = () => {
  showReplyForm.value = null;
  subReply.value = { username: "", content: "" };
};

// 返回上一页
const goBack = () => {
  router.go(-1);
};

// 组件挂载时获取数据
onMounted(async () => {
  await fetchPost();
  await fetchReplies();
});
</script>

<style scoped>
/* 修复图标大小 */
svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.icon-sm {
  width: 0.875rem;
  height: 0.875rem;
}

.icon-md {
  width: 1.25rem;
  height: 1.25rem;
}

/* 确保内容不会被图标影响 */
.flex svg {
  margin-right: 0.25rem;
}
</style>
