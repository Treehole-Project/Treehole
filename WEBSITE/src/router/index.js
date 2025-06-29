import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import Search from "../views/Search.vue";
import PostDetail from "../views/PostDetail.vue";
import SearchResults from "../views/SearchResults.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: { title: "Treehole - 首页" },
  },
  {
    path: "/search",
    name: "Search",
    component: Search,
    meta: { title: "Treehole - 搜索" },
  },
  {
    path: "/post/:id",
    name: "PostDetail",
    component: PostDetail,
    meta: { title: "Treehole - 帖子详情" },
  },
  {
    path: "/search-results",
    name: "SearchResults",
    component: SearchResults,
    meta: { title: "Treehole - 搜索结果" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫 - 设置页面标题
router.beforeEach((to, from, next) => {
  document.title = to.meta.title || "树洞网站";
  next();
});

export default router;
