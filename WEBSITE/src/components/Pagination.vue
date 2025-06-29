<template>
  <div class="space-y-4">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-8">
      <div
        class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"
      ></div>
    </div>

    <!-- 分页信息 -->
    <div
      v-if="pagination && !loading"
      class="flex items-center justify-between text-sm text-gray-600 bg-white px-4 py-3 rounded-lg border"
    >
      <span>
        共 {{ pagination.total }} 条记录，第 {{ pagination.page }} /
        {{ pagination.pages }} 页
      </span>
    </div>

    <!-- 分页控件 -->
    <div
      v-if="pagination && pagination.pages > 1"
      class="flex items-center justify-center space-x-2"
    >
      <!-- 上一页 -->
      <button
        @click="goToPage(pagination.page - 1)"
        :disabled="pagination.page <= 1"
        class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        上一页
      </button>

      <!-- 页码 -->
      <template v-for="page in visiblePages" :key="page">
        <button
          v-if="page !== '...'"
          @click="goToPage(page)"
          :class="[
            'px-3 py-2 text-sm font-medium rounded-md',
            page === pagination.page
              ? 'bg-blue-500 text-white'
              : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50',
          ]"
        >
          {{ page }}
        </button>
        <span v-else class="px-3 py-2 text-sm text-gray-500">...</span>
      </template>

      <!-- 下一页 -->
      <button
        @click="goToPage(pagination.page + 1)"
        :disabled="pagination.page >= pagination.pages"
        class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  pagination: {
    type: Object,
    default: null,
  },
  loading: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["page-change"]);

const goToPage = (page) => {
  if (
    page >= 1 &&
    page <= props.pagination.pages &&
    page !== props.pagination.page
  ) {
    emit("page-change", page);
  }
};

// 计算可见的页码
const visiblePages = computed(() => {
  if (!props.pagination) return [];

  const current = props.pagination.page;
  const total = props.pagination.pages;
  const pages = [];

  if (total <= 7) {
    // 如果总页数少于等于7页，显示所有页码
    for (let i = 1; i <= total; i++) {
      pages.push(i);
    }
  } else {
    // 复杂的分页逻辑
    if (current <= 4) {
      // 当前页在前面
      for (let i = 1; i <= 5; i++) {
        pages.push(i);
      }
      pages.push("...");
      pages.push(total);
    } else if (current >= total - 3) {
      // 当前页在后面
      pages.push(1);
      pages.push("...");
      for (let i = total - 4; i <= total; i++) {
        pages.push(i);
      }
    } else {
      // 当前页在中间
      pages.push(1);
      pages.push("...");
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i);
      }
      pages.push("...");
      pages.push(total);
    }
  }

  return pages;
});
</script>
