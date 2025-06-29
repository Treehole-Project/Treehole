<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <Navbar />

    <!-- 路由视图 -->
    <router-view />
  </div>
</template>

<script setup>
import Navbar from "./components/Navbar.vue";
</script>
