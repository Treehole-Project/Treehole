package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"treehole/internal/models"
	"treehole/internal/scraper"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB, scraperService *scraper.Service) *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建处理器
	handler := &Handler{
		db:             db,
		scraperService: scraperService,
	}

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 帖子相关路由
		api.GET("/posts", handler.GetPosts)
		api.GET("/posts/:id", handler.GetPost)
		api.GET("/posts/:id/replies", handler.GetPostReplies)
		api.POST("/posts", handler.CreatePost)
		api.POST("/posts/:id/replies", handler.CreateReply)

		// 搜索路由
		api.GET("/search", handler.SearchPosts)
		api.GET("/search/advanced", handler.AdvancedSearch)
		// api.GET("/search/users", handler.SearchUsers)
		// api.GET("/search/comments", handler.SearchComments)

		// 用户相关路由
		// api.GET("/users/:user_id/posts", handler.GetUserPosts)
		// api.GET("/users/:user_id/replies", handler.GetUserReplies)

		// 标签路由
		api.GET("/tags", handler.GetTags)
		api.GET("/tags/:name/posts", handler.GetPostsByTag)

		// 统计路由
		api.GET("/stats", handler.GetStats)

		// 同步相关路由
		// api.POST("/sync", handler.TriggerSync)
		api.GET("/sync/status", handler.GetSyncStatus)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 静态文件托管
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/", "./dist/index.html")
	r.StaticFile("/index.html", "./dist/index.html")
	
	// SPA 路由支持 - 如果路由不匹配，返回 index.html
	r.NoRoute(func(c *gin.Context) {
		// 如果请求的是 API 路径，返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		// 其他路径返回 index.html，支持前端路由
		c.File("./dist/index.html")
	})

	return r
}

// Handler API 处理器
type Handler struct {
	db             *gorm.DB
	scraperService *scraper.Service
}

// GetPosts 获取帖子列表
func (h *Handler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 获取总数
	h.db.Model(&models.Post{}).Count(&total)

	// 获取帖子列表
	if err := h.db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetPost 获取单个帖子
func (h *Handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := h.db.Where("id = ? OR original_id = ?", id, id).
		First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 增加浏览次数
	h.db.Model(&post).Update("view_count", gorm.Expr("view_count + ?", 1))

	c.JSON(http.StatusOK, post)
}

// GetPostReplies 获取帖子回复
func (h *Handler) GetPostReplies(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	// 先找到帖子
	var post models.Post
	if err := h.db.Where("id = ? OR original_id = ?", id, id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var replies []models.Reply
	var total int64

	h.db.Unscoped().Model(&models.Reply{}).Where("post_id = ?", post.ID).Count(&total)

	if err := h.db.Unscoped().Where("post_id = ?", post.ID).
		Order("created_at asc").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": replies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// CreatePostRequest 创建帖子的请求结构
type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	UserName string `json:"username" binding:"required"`
}

// CreateReplyRequest 创建回复的请求结构
type CreateReplyRequest struct {
	Content  string `json:"content" binding:"required"`
	UserName string `json:"username" binding:"required"`
	ParentID int    `json:"parent_id"` // 父评论ID，如果是0则回复帖子本身
}

// CreatePost 创建帖子
func (h *Handler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建本地帖子记录
	post := models.Post{
		Title:       req.Title,
		Content:     req.Content,
		Author:      req.UserName,
		OriginalID: "0", // 默认原始ID
		AuthorID:    "automatic"+time.Now().Format("20060102150405"), // 默认openid
		RadioGroup:  "radio40",   // 默认分组
		CampusGroup: "2",         // 默认校区
		Region:      "0",         // 默认地区
		LikeNum:     0,
		ReplyCount:  0,
		ViewCount:   0,
		Tag:         "未分析",
		State:       "normal",
		Images:      "[]",
		Cover:       "[]",
		CreatedAt:   time.Now(),
	}

	// 保存到本地数据库
	if err := h.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post locally"})
		return
	}

	// 同步到主站
	go func() {
		if os.Getenv("SYNC_ENABLED") != "true" {
			return
		}
		if err := h.scraperService.SyncPostToMainSite(post); err != nil {
			log.Printf("Failed to sync post to main site: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// CreateReply 创建回复
func (h *Handler) CreateReply(c *gin.Context) {
	postID := c.Param("id")
	
	var req CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 先找到帖子
	var post models.Post
	if err := h.db.Where("id = ? OR original_id = ?", postID, postID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建回复记录
	reply := models.Reply{
		PostID:    post.ID,
		Content:   req.Content,
		Author:    req.UserName,
		AuthorID:  "automatic"+time.Now().Format("20060102150405"), // 默认openid
		ApplyTo:   "automatic",
		Level:     1,           // 默认层级
		ParentID:  req.ParentID,
		LikeNum:   0,
		Tag:       "未分析",
		Images:    "[]",
		CreatedAt: time.Now(),
	}

	// 如果有父评论ID，设置层级为2
	if req.ParentID > 0 {
		reply.Level = 2
	}

	// 保存到本地数据库
	if err := h.db.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply locally"})
		return
	}

	// 更新帖子的评论数
	h.db.Model(&post).Update("reply_count", gorm.Expr("reply_count + ?", 1))

	// 同步到主站
	go func() {
		if os.Getenv("SYNC_ENABLED") != "true" {
			return
		}
		if err := h.scraperService.SyncReplyToMainSite(post, reply); err != nil {
			log.Printf("Failed to sync reply to main site: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reply created successfully",
		"reply":   reply,
	})
}

// SearchPosts 搜索帖子
func (h *Handler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 分割关键词
	keywords := splitKeywords(query)
	if len(keywords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid keywords provided"})
		return
	}

	// 构建多关键词搜索条件
	var titleConditions []string
	var contentConditions []string
	var args []interface{}

	for _, keyword := range keywords {
		titleConditions = append(titleConditions, "title LIKE ?")
		contentConditions = append(contentConditions, "content LIKE ?")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	// 标题或内容包含所有关键词
	whereClause := "((" + strings.Join(titleConditions, " AND ") + ") OR (" + strings.Join(contentConditions, " AND ") + "))"

	h.db.Model(&models.Post{}).
		Where(whereClause, args...).
		Count(&total)

	if err := h.db.Where(whereClause, args...).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
		"query": query,
		"keywords": keywords,
	})
}

// GetTags 获取标签列表 (从帖子中提取唯一标签)
func (h *Handler) GetTags(c *gin.Context) {
	var tags []string
	
	// 从帖子中获取所有唯一标签
	rows, err := h.db.Model(&models.Post{}).Select("DISTINCT tag").Where("tag != ''").Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err == nil && tag != "" {
			tags = append(tags, tag)
		}
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// GetPostsByTag 根据标签获取帖子
func (h *Handler) GetPostsByTag(c *gin.Context) {
	tagName := c.Param("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 直接通过tag字段查询
	h.db.Model(&models.Post{}).Where("tag = ?", tagName).Count(&total)

	if err := h.db.Where("tag = ?", tagName).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
		"tag": tagName,
	})
}

// GetStats 获取统计信息
func (h *Handler) GetStats(c *gin.Context) {
	var totalPosts int64
	var totalReplies int64
	var totalTags int64

	h.db.Model(&models.Post{}).Count(&totalPosts)
	h.db.Model(&models.Reply{}).Count(&totalReplies)
	
	// 统计唯一标签数量
	h.db.Model(&models.Post{}).Select("COUNT(DISTINCT tag)").Where("tag != ''").Scan(&totalTags)

	// 获取最新帖子
	var latestPost models.Post
	h.db.Order("created_at desc").First(&latestPost)

	c.JSON(http.StatusOK, gin.H{
		"total_posts":   totalPosts,
		"total_replies": totalReplies,
		"total_tags":    totalTags,
		"latest_post":   latestPost,
	})
}

// TriggerSync 触发同步
func (h *Handler) TriggerSync(c *gin.Context) {
	// 异步执行同步
	go func() {
		if err := h.scraperService.ScrapeData(); err != nil {
			// 这里可以添加日志记录
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Sync started"})
}

// GetSyncStatus 获取同步状态
func (h *Handler) GetSyncStatus(c *gin.Context) {
	status, err := h.scraperService.GetLastSyncStatus()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No sync history found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// AdvancedSearch 高级搜索
func (h *Handler) AdvancedSearch(c *gin.Context) {
	// 搜索参数
	title := c.Query("title")       // 标题关键词
	content := c.Query("content")   // 内容关键词
	author := c.Query("author")     // 作者用户名
	authorID := c.Query("author_id") // 作者ID (openid)
	postID := c.Query("post_id")    // 帖子ID
	originalID := c.Query("original_id") // 原始ID
	comment := c.Query("comment")   // 评论内容（搜索评论并返回对应帖子）
	tag := c.Query("tag")           // 标签
	state := c.Query("state")       // 状态
	radioGroup := c.Query("radio_group") // 分组
	logic := c.DefaultQuery("logic", "and") // 逻辑关系：and 或 or
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 构建查询条件
	db := h.db.Model(&models.Post{})
	
	// 存储所有条件
	var conditions []string
	var args []interface{}

	// 标题搜索
	if title != "" {
		keywords := splitKeywords(title)
		if len(keywords) > 0 {
			titleCondition, titleArgs := buildMultiKeywordCondition("title", keywords)
			conditions = append(conditions, titleCondition)
			args = append(args, titleArgs...)
		}
	}

	// 内容搜索
	if content != "" {
		keywords := splitKeywords(content)
		if len(keywords) > 0 {
			contentCondition, contentArgs := buildMultiKeywordCondition("content", keywords)
			conditions = append(conditions, contentCondition)
			args = append(args, contentArgs...)
		}
	}

	// 作者搜索
	if author != "" {
		keywords := splitKeywords(author)
		if len(keywords) > 0 {
			authorCondition, authorArgs := buildMultiKeywordCondition("author", keywords)
			conditions = append(conditions, authorCondition)
			args = append(args, authorArgs...)
		}
	}

	// 作者ID搜索
	if authorID != "" {
		conditions = append(conditions, "author_id = ?")
		args = append(args, authorID)
	}

	// 帖子ID搜索
	if postID != "" {
		conditions = append(conditions, "id = ?")
		args = append(args, postID)
	}

	// 原始ID搜索
	if originalID != "" {
		conditions = append(conditions, "original_id = ?")
		args = append(args, originalID)
	}

	// 标签搜索
	if tag != "" {
		conditions = append(conditions, "tag = ?")
		args = append(args, tag)
	}

	// 状态搜索
	if state != "" {
		conditions = append(conditions, "state = ?")
		args = append(args, state)
	}

	// 分组搜索
	if radioGroup != "" {
		conditions = append(conditions, "radio_group = ?")
		args = append(args, radioGroup)
	}

	// 评论搜索（通过评论内容找到对应的帖子）
	if comment != "" {
		keywords := splitKeywords(comment)
		if len(keywords) > 0 {
			// 为每个关键词构建子查询
			var commentConditions []string
			for _, keyword := range keywords {
				subQuery := h.db.Unscoped().Model(&models.Reply{}).
					Select("DISTINCT post_id").
					Where("content LIKE ?", "%"+keyword+"%")
				commentConditions = append(commentConditions, "id IN (?)")
				args = append(args, subQuery)
			}
			// 所有关键词都必须匹配（AND关系）
			if len(commentConditions) == 1 {
				conditions = append(conditions, commentConditions[0])
			} else {
				conditions = append(conditions, "("+strings.Join(commentConditions, " AND ")+")")
			}
		}
	}

	// 根据逻辑关系连接条件
	if len(conditions) > 0 {
		var whereClause string
		if logic == "or" {
			whereClause = "(" + conditions[0]
			for i := 1; i < len(conditions); i++ {
				whereClause += " OR " + conditions[i]
			}
			whereClause += ")"
		} else { // and
			whereClause = conditions[0]
			for i := 1; i < len(conditions); i++ {
				whereClause += " AND " + conditions[i]
			}
		}
		
		db = db.Where(whereClause, args...)
	}

	// 获取总数
	db.Count(&total)

	// 获取结果
	if err := db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
		"search_params": gin.H{
			"title":       title,
			"content":     content,
			"author":      author,
			"author_id":   authorID,
			"post_id":     postID,
			"original_id": originalID,
			"comment":     comment,
			"tag":         tag,
			"state":       state,
			"radio_group": radioGroup,
			"logic":       logic,
		},
	})
}

// SearchUsers 搜索用户（按用户名或ID）
func (h *Handler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	type UserInfo struct {
		Author   string `json:"author"`
		AuthorID string `json:"author_id"`
		PostCount int64 `json:"post_count"`
		ReplyCount int64 `json:"reply_count"`
	}

	var users []UserInfo

	// 按用户名或ID搜索，并聚合统计信息
	searchPattern := "%" + query + "%"
	
	rows, err := h.db.Raw(`
		SELECT 
			author,
			author_id,
			COUNT(*) as post_count,
			(SELECT COUNT(*) FROM replies WHERE replies.author_id = posts.author_id) as reply_count
		FROM posts 
		WHERE (author LIKE ? OR author_id LIKE ?) 
			AND deleted_at IS NULL
		GROUP BY author_id, author
		ORDER BY post_count DESC
		LIMIT ? OFFSET ?
	`, searchPattern, searchPattern, limit, offset).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.Author, &user.AuthorID, &user.PostCount, &user.ReplyCount); err == nil {
			users = append(users, user)
		}
	}

	// 获取总数
	var total int64
	h.db.Model(&models.Post{}).
		Select("DISTINCT author_id").
		Where("author LIKE ? OR author_id LIKE ?", searchPattern, searchPattern).
		Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
		"query": query,
	})
}

// GetUserPosts 获取指定用户的帖子
func (h *Handler) GetUserPosts(c *gin.Context) {
	userID := c.Param("user_id") // 可以是 author_id 或 author
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 同时按 author_id 和 author 搜索
	db := h.db.Model(&models.Post{}).Where("author_id = ? OR author = ?", userID, userID)
	
	db.Count(&total)

	if err := db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取用户信息
	var userInfo struct {
		Author   string `json:"author"`
		AuthorID string `json:"author_id"`
	}
	
	if len(posts) > 0 {
		userInfo.Author = posts[0].Author
		userInfo.AuthorID = posts[0].AuthorID
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"user_info": userInfo,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUserReplies 获取指定用户的回复
func (h *Handler) GetUserReplies(c *gin.Context) {
	userID := c.Param("user_id") // 可以是 author_id 或 author
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var replies []models.Reply
	var total int64

	// 同时按 author_id 和 author 搜索
	db := h.db.Unscoped().Model(&models.Reply{}).Where("author_id = ? OR author = ?", userID, userID)
	
	db.Count(&total)

	if err := db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取用户信息
	var userInfo struct {
		Author   string `json:"author"`
		AuthorID string `json:"author_id"`
	}
	
	if len(replies) > 0 {
		userInfo.Author = replies[0].Author
		userInfo.AuthorID = replies[0].AuthorID
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": replies,
		"user_info": userInfo,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// SearchComments 搜索评论
func (h *Handler) SearchComments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	author := c.Query("author")     // 作者用户名
	authorID := c.Query("author_id") // 作者ID (openid)
	postID := c.Query("post_id")    // 限制在某个帖子内搜索
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var replies []models.Reply
	var total int64

	// 构建查询条件
	db := h.db.Unscoped().Model(&models.Reply{})

	// 内容搜索
	keywords := splitKeywords(query)
	if len(keywords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid keywords provided"})
		return
	}

	// 构建多关键词搜索条件
	contentCondition, contentArgs := buildMultiKeywordCondition("content", keywords)
	db = db.Where(contentCondition, contentArgs...)

	// 作者搜索
	if author != "" {
		authorKeywords := splitKeywords(author)
		if len(authorKeywords) > 0 {
			authorCondition, authorArgs := buildMultiKeywordCondition("author", authorKeywords)
			db = db.Where(authorCondition, authorArgs...)
		}
	}

	// 作者ID搜索
	if authorID != "" {
		db = db.Where("author_id = ?", authorID)
	}

	// 限制在某个帖子内搜索
	if postID != "" {
		// 先找到帖子的内部ID
		var post models.Post
		if err := h.db.Where("id = ? OR original_id = ?", postID, postID).First(&post).Error; err == nil {
			db = db.Where("post_id = ?", post.ID)
		}
	}

	db.Count(&total)

	if err := db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 为每个回复添加对应的帖子信息
	type ReplyWithPost struct {
		models.Reply
		PostTitle    string `json:"post_title"`
		PostOriginalID string `json:"post_original_id"`
	}

	var repliesWithPost []ReplyWithPost
	for _, reply := range replies {
		var post models.Post
		if err := h.db.Where("id = ?", reply.PostID).First(&post).Error; err == nil {
			repliesWithPost = append(repliesWithPost, ReplyWithPost{
				Reply:          reply,
				PostTitle:      post.Title,
				PostOriginalID: post.OriginalID,
			})
		} else {
			repliesWithPost = append(repliesWithPost, ReplyWithPost{
				Reply:          reply,
				PostTitle:      "未知帖子",
				PostOriginalID: "",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": repliesWithPost,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
		"search_params": gin.H{
			"query":     query,
			"author":    author,
			"author_id": authorID,
			"post_id":   postID,
		},
	})
}

// splitKeywords 将搜索词按空格分割并去除空白
func splitKeywords(query string) []string {
	keywords := strings.Fields(strings.TrimSpace(query))
	var result []string
	for _, keyword := range keywords {
		if keyword != "" {
			result = append(result, keyword)
		}
	}
	return result
}

// buildMultiKeywordCondition 构建多关键词搜索条件
func buildMultiKeywordCondition(field string, keywords []string) (string, []interface{}) {
	if len(keywords) == 0 {
		return "", nil
	}
	
	var conditions []string
	var args []interface{}
	
	for _, keyword := range keywords {
		conditions = append(conditions, field+" LIKE ?")
		args = append(args, "%"+keyword+"%")
	}
	
	// 所有关键词都必须匹配（AND关系）
	condition := "(" + strings.Join(conditions, " AND ") + ")"
	return condition, args
}
