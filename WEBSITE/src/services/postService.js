import api from "./api.js";

export const postService = {
  // 获取帖子列表
  async getPosts(page = 1, limit = 20) {
    return await api.get("/posts", {
      params: { page, limit },
    });
  },

  // 获取单个帖子
  async getPost(id) {
    return await api.get(`/posts/${id}`);
  },

  // 获取帖子回复
  async getPostReplies(id, page = 1, limit = 20) {
    return await api.get(`/posts/${id}/replies`, {
      params: { page, limit },
    });
  },

  // 创建帖子
  async createPost(postData) {
    return await api.post("/posts", postData);
  },

  // 创建回复
  async createReply(postId, replyData) {
    return await api.post(`/posts/${postId}/replies`, replyData);
  },

  // 基础搜索
  async searchPosts(query, page = 1, limit = 20) {
    return await api.get("/search", {
      params: { q: query, page, limit },
    });
  },

  // 高级搜索
  async advancedSearch(searchParams) {
    return await api.get("/search/advanced", {
      params: searchParams,
    });
  },

  // 获取标签列表
  async getTags() {
    return await api.get("/tags");
  },

  // 根据标签获取帖子
  async getPostsByTag(tagName, page = 1, limit = 20) {
    return await api.get(`/tags/${tagName}/posts`, {
      params: { page, limit },
    });
  },

  // 获取统计信息
  async getStats() {
    return await api.get("/stats");
  },
};
