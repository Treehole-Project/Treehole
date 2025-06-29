<template>
  <article
    class="card p-6 mb-4 hover:shadow-lg transition-shadow duration-200 cursor-pointer"
    @click="goToPost"
  >
    <!-- 帖子标题 -->
    <h2
      class="text-lg font-semibold text-gray-900 mb-3 truncate hover:text-blue-600 transition-colors"
      :title="post.title"
    >
      {{ post.title }}
    </h2>

    <!-- 帖子内容预览 -->
    <p class="text-gray-600 mb-4 line-clamp-3">
      {{ post.content }}
    </p>

    <!-- 帖子元信息 -->
    <div
      class="flex flex-wrap items-center justify-between text-sm text-gray-500"
    >
      <div class="flex items-center space-x-4 mb-2 sm:mb-0">
        <!-- 作者 -->
        <span class="flex items-center">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
              clip-rule="evenodd"
            />
          </svg>
          {{ post.author }}
        </span>

        <!-- 发布时间 -->
        <span
          class="flex items-center"
          :title="formatFullDateTime(post.created_at)"
        >
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
              clip-rule="evenodd"
            />
          </svg>
          {{ formatDate(post.created_at) }}
        </span>
      </div>

      <div class="flex items-center space-x-4">
        <!-- 查看数量 -->
        <span class="flex items-center">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
            <path
              fill-rule="evenodd"
              d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
              clip-rule="evenodd"
            />
          </svg>
          {{ post.view_count || 0 }}
        </span>

        <!-- 点赞数量 -->
        <span class="flex items-center">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z"
            />
          </svg>
          {{ post.like_num || 0 }}
        </span>

        <!-- 评论数量 -->
        <span class="flex items-center">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z"
              clip-rule="evenodd"
            />
          </svg>
          {{ post.reply_count || 0 }}
        </span>
      </div>
    </div>

    <!-- 标签 -->
    <div v-if="post.tag && post.tag !== '未分析'" class="mt-3">
      <span
        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
      >
        {{ post.tag }}
      </span>
    </div>
  </article>
</template>

<script setup>
import { useRouter } from "vue-router";
import { formatDate, formatFullDateTime } from "../utils/dateUtils.js";

const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
});

const router = useRouter();

const goToPost = () => {
  router.push(`/post/${props.post.id}`);
};
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 确保图标大小正确 */
svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}
</style>
