<template>
  <!-- 图片查看器模态框 -->
  <div
    v-if="isVisible"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75"
    @click="close"
  >
    <div class="relative max-w-4xl max-h-full p-4">
      <!-- 关闭按钮 -->
      <button
        @click="close"
        class="absolute top-2 right-2 z-10 p-2 text-white bg-black bg-opacity-50 rounded-full hover:bg-opacity-75 transition-opacity"
      >
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>

      <!-- 图片 -->
      <img
        :src="imageUrl"
        :alt="altText"
        class="max-w-full max-h-full object-contain rounded-lg cursor-pointer"
        @click.stop
        @load="onImageLoad"
        @error="onImageError"
      />

      <!-- 加载状态 -->
      <div
        v-if="loading"
        class="absolute inset-0 flex items-center justify-center"
      >
        <div
          class="animate-spin rounded-full h-8 w-8 border-b-2 border-white"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch } from "vue";

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false,
  },
  imageUrl: {
    type: String,
    default: "",
  },
  altText: {
    type: String,
    default: "图片",
  },
});

const emit = defineEmits(["close"]);

const loading = ref(false);

const close = () => {
  emit("close");
};

const onImageLoad = () => {
  loading.value = false;
};

const onImageError = () => {
  loading.value = false;
};

// 监听图片URL变化，显示加载状态
watch(
  () => props.imageUrl,
  () => {
    if (props.imageUrl) {
      loading.value = true;
    }
  }
);

// 监听是否显示，重置加载状态
watch(
  () => props.isVisible,
  (newVal) => {
    if (newVal && props.imageUrl) {
      loading.value = true;
    }
  }
);
</script>

<style scoped>
/* 防止页面滚动 */
.fixed {
  overflow: hidden;
}
</style>
