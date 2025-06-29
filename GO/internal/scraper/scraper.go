package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"treehole/internal/config"
	"treehole/internal/database"
	"treehole/internal/models"

	"gorm.io/gorm"
)

// Service 爬虫服务
type Service struct {
	db         *gorm.DB
	client     *http.Client // 用于抓取数据的客户端（不使用代理）
	syncClient *http.Client // 用于同步到主站的客户端（使用代理）
	baseURL    string
	config     *config.Config
	saveMux    sync.Mutex // 保护数据库写入操作的互斥锁
}

// APIResponse 通用 API 响应结构
type APIResponse struct {
	TaskList    []TaskData    `json:"taskList"`
	CommentList []CommentData `json:"commentList"`
}

// TaskData 帖子数据结构
type TaskData struct {
	ID           int    `json:"id"`
	IP           string `json:"ip"`
	Content      string `json:"content"`
	Price        string `json:"price"`
	Title        string `json:"title"`
	Wechat       string `json:"wechat"`
	OpenID       string `json:"openid"`
	Avatar       string `json:"avatar"`
	CampusGroup  string `json:"campusGroup"`
	CommentNum   int    `json:"commentNum"`
	WatchNum     int    `json:"watchNum"`
	LikeNum      int    `json:"likeNum"`
	RadioGroup   string `json:"radioGroup"`
	Images       string `json:"img"`
	Cover        string `json:"cover"`
	IsDelete     int    `json:"is_delete"`
	IsComplaint  int    `json:"is_complaint"`
	Region       string `json:"region"`
	UserName     string `json:"userName"`
	CTime        string `json:"c_time"`
	CommentTime  string `json:"comment_time"`
	Choose       int    `json:"choose"`
	Hot          int    `json:"hot"`
}

// CommentData 评论数据结构
type CommentData struct {
	ID          int           `json:"id"`
	OpenID      string        `json:"openid"`
	ApplyTo     string        `json:"applyTo"`
	Avatar      string        `json:"avatar"`
	Comment     string        `json:"comment"`
	PK          int           `json:"pk"`
	UserName    string        `json:"userName"`
	CTime       string        `json:"c_time"`
	Images      string        `json:"img"`
	Level       interface{}   `json:"level"` // 可能是字符串或数字
	PID         int           `json:"pid"`
	LikeNum     int           `json:"like_num"`
	CommentList []CommentData `json:"commentList"`
}

// NewService 创建新的爬虫服务
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	// 创建用于抓取数据的HTTP客户端（不使用代理）
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建用于同步到主站的HTTP客户端
	syncClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 如果启用了代理，只为同步客户端配置代理
	if cfg.ProxyEnabled && cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			log.Printf("Failed to parse proxy URL: %v", err)
		} else {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			syncClient.Transport = transport
			log.Printf("Proxy enabled for sync operations: %s", cfg.ProxyURL)
		}
	}

	return &Service{
		db:         db,
		client:     client,
		syncClient: syncClient,
		baseURL:    "https://www.yqtech.ltd:8802",
		config:     cfg,
		saveMux:    sync.Mutex{},
	}
}

// ScrapeData 抓取数据 - 主要的同步方法，增加事务处理
func (s *Service) ScrapeData() error {
	log.Println("Starting data synchronization...")

	var totalPosts, totalReplies int
	var errors []string

	// 使用事务记录同步状态
	err := database.SafeTransaction(s.db, func(tx *gorm.DB) error {
		// 记录同步开始
		syncStatus := &models.SyncStatus{
			LastSyncTime: time.Now(),
			Status:       "running",
		}
		return tx.Create(syncStatus).Error
	})

	if err != nil {
		log.Printf("Failed to create sync status: %v", err)
	}

	// 1. 抓取新帖子（从本地最大ID到最新ID）
	if err := s.scrapeNewPosts(&totalPosts, &errors); err != nil {
		log.Printf("Error scraping new posts: %v", err)
		errors = append(errors, fmt.Sprintf("New posts error: %v", err))
	}

	// 2. 抓取有新回复的帖子
	if err := s.scrapeNewReplies(&totalReplies, &errors); err != nil {
		log.Printf("Error scraping new replies: %v", err)
		errors = append(errors, fmt.Sprintf("New replies error: %v", err))
	}

	// 更新同步状态
	err = database.WithRetry(s.db, func(db *gorm.DB) error {
		syncStatus := &models.SyncStatus{
			LastSyncTime: time.Now(),
			Status:       "success",
			TotalPosts:   totalPosts,
			TotalReplies: totalReplies,
		}
		
		if len(errors) > 0 {
			syncStatus.Status = "error"
			syncStatus.ErrorMessage = strings.Join(errors, "; ")
		}
		
		return db.Create(syncStatus).Error
	})

	if err != nil {
		log.Printf("Failed to update sync status: %v", err)
	}

	log.Printf("Sync completed. Posts: %d, Replies: %d, Errors: %d", totalPosts, totalReplies, len(errors))
	return nil
}

// scrapeNewPosts 抓取新帖子
func (s *Service) scrapeNewPosts(totalPosts *int, errors *[]string) error {
	// 获取本地最大ID
	localMaxID := s.getLocalMaxPostID()
	if localMaxID == "" || localMaxID == "0" {
		localMaxID = "300003"
	} else {
		if id, err := strconv.Atoi(localMaxID); err != nil || id < 300003 { // 277506
			localMaxID = "300003"
		}
	}
	
	// 获取远程最大ID
	remoteMaxID, err := s.getMaxIndex()
	if err != nil {
		return fmt.Errorf("failed to get max index: %v", err)
	}

	log.Printf("Syncing posts from %s to %s", localMaxID, remoteMaxID)

	// 从本地最大ID+1开始，到远程最大ID为止
	startID, _ := strconv.Atoi(localMaxID)
	endID, _ := strconv.Atoi(remoteMaxID)

	for id := startID + 1; id <= endID; id++ {
		postID := strconv.Itoa(id)
		
		// 获取帖子信息
		post, err := s.getTask(postID)
		if err != nil {
			log.Printf("Failed to get post %s: %v", postID, err)
			*errors = append(*errors, fmt.Sprintf("Post %s: %v", postID, err))
			continue
		}

		if post == nil {
			continue // 帖子不存在或已删除
		}

		// 保存帖子
		if err := s.savePost(post); err != nil {
			log.Printf("Failed to save post %s: %v", postID, err)
			*errors = append(*errors, fmt.Sprintf("Save post %s: %v", postID, err))
			continue
		}

		// 获取并保存评论
		if err := s.scrapePostComments(postID); err != nil {
			log.Printf("Failed to scrape comments for post %s: %v", postID, err)
			*errors = append(*errors, fmt.Sprintf("Comments %s: %v", postID, err))
		}

		*totalPosts++
		
		// 添加延迟避免请求过快，减少数据库压力
		// time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// scrapeNewReplies 抓取有新回复的帖子
func (s *Service) scrapeNewReplies(totalReplies *int, errors *[]string) error {
	// 获取有新回复的帖子列表
	newReplyPosts, err := s.getNewReplyPosts()
	if err != nil {
		return fmt.Errorf("failed to get new reply posts: %v", err)
	}

	log.Printf("Found %d posts with new replies", len(newReplyPosts))

	// 获取上次检查新回复的时间
	lastCheckTime := s.getLastReplyCheckTime()

	for _, post := range newReplyPosts {
		// 检查帖子的最后回复时间是否在我们上次检查之后
		commentTime := s.parseTime(post.CommentTime)
		
		log.Printf("Post ID: %d, Title: %s, Last Reply Time: %s, Last Check Time: %s", post.ID, post.Title, commentTime.Format(time.RFC3339), lastCheckTime.Format(time.RFC3339))
		if !lastCheckTime.IsZero() && commentTime.Before(lastCheckTime) {
			continue // 这个帖子的回复我们已经处理过了
		}

		// 重新抓取这个帖子的所有评论
		if err := s.scrapePostComments(strconv.Itoa(post.ID)); err != nil {
			log.Printf("Failed to scrape comments for post %d: %v", post.ID, err)
			*errors = append(*errors, fmt.Sprintf("Comments %d: %v", post.ID, err))
			continue
		}

		*totalReplies++
		// time.Sleep(200 * time.Millisecond)
	}

	// 更新最后检查时间
	s.updateLastReplyCheckTime()

	return nil
}

// API 相关方法

// get 发送 GET 请求
func (s *Service) get(url string) (*APIResponse, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// getMaxIndex 获取最大帖子ID
func (s *Service) getMaxIndex() (string, error) {
	url := fmt.Sprintf("%s/gettaskbyType?length=0&radioGroup=%%5B%%22radio4%%22%%2C%%22radio40%%22%%2C%%22radio41%%22%%2C%%22radio42%%22%%2C%%22radio43%%22%%5D&type=0", s.baseURL)
	resp, err := s.get(url)
	if err != nil {
		return "", err
	}

	if len(resp.TaskList) == 0 {
		return "0", nil
	}

	return strconv.Itoa(resp.TaskList[0].ID), nil
}

// getTask 获取单个帖子
func (s *Service) getTask(pk string) (*TaskData, error) {
	url := fmt.Sprintf("%s/gettaskbyId?pk=%s", s.baseURL, pk)
	resp, err := s.get(url)
	if err != nil {
		return nil, err
	}

	if len(resp.TaskList) == 0 {
		return nil, nil // 帖子不存在
	}

	return &resp.TaskList[0], nil
}

// getComments 获取帖子的所有评论
func (s *Service) getComments(pk string) ([]CommentData, error) {
	var allComments []CommentData
	length := 0

	for {
		url := fmt.Sprintf("%s/getCommentByType?length=%d&pk=%s&type=0", s.baseURL, length, pk)
		resp, err := s.get(url)
		if err != nil {
			return nil, err
		}

		if len(resp.CommentList) == 0 {
			break
		}

		allComments = append(allComments, resp.CommentList...)
		length += len(resp.CommentList)
	}

	return allComments, nil
}

// getNewReplyPosts 获取有新回复的帖子
func (s *Service) getNewReplyPosts() ([]TaskData, error) {
	url := fmt.Sprintf("%s/gettaskbyType?length=0&radioGroup=%%5B%%22radio4%%22%%2C%%22radio40%%22%%2C%%22radio41%%22%%2C%%22radio42%%22%%2C%%22radio43%%22%%5D&type=1", s.baseURL)
	resp, err := s.get(url)
	if err != nil {
		return nil, err
	}

	return resp.TaskList, nil
}

// getUserPosts 获取指定用户的最新帖子（只需要第一个结果）
func (s *Service) getUserPosts(openid string) ([]TaskData, error) {
	// 只获取第一个结果，offset=0
	url := fmt.Sprintf("%s/gettaskbyOpenId?openid=%s&length=0", s.baseURL, openid)
	resp, err := s.get(url)
	if err != nil {
		return nil, err
	}

	return resp.TaskList, nil
}

// getUserComments 获取指定用户的最新评论（只需要第一个结果）
func (s *Service) getUserComments(openid string) ([]CommentData, error) {
	// 只获取第一个结果，offset=0
	url := fmt.Sprintf("%s/getCommentByOpenid?openid=%s&length=0", s.baseURL, openid)
	resp, err := s.get(url)
	if err != nil {
		return nil, err
	}

	return resp.CommentList, nil
}

// 数据库相关方法

// getLocalMaxPostID 获取本地数据库中最大的帖子ID
func (s *Service) getLocalMaxPostID() string {
	var post models.Post
	err := database.WithRetry(s.db, func(db *gorm.DB) error {
		return db.Order("CAST(original_id AS INTEGER) DESC").First(&post).Error
	})
	if err != nil {
		return "0" // 如果没有帖子，从0开始
	}
	return post.OriginalID
}

// savePost 保存帖子到数据库 - 使用重试机制
func (s *Service) savePost(taskData *TaskData) error {
	s.saveMux.Lock()
	defer s.saveMux.Unlock()

	return database.WithRetry(s.db, func(db *gorm.DB) error {
		// 检查帖子是否已存在
		var existingPost models.Post
		result := db.Where("original_id = ?", strconv.Itoa(taskData.ID)).First(&existingPost)
		
		createdAt := s.parseTime(taskData.CTime)

		if result.Error == gorm.ErrRecordNotFound {
			// 新帖子
			post := models.Post{
				OriginalID:  strconv.Itoa(taskData.ID),
				Title:       taskData.Title,
				Content:     taskData.Content,
				Author:      taskData.UserName,
				AuthorID:    taskData.OpenID,
				IP:          taskData.IP,
				LikeNum:     taskData.LikeNum,
				CreatedAt:   createdAt,
				ReplyCount:  taskData.CommentNum,
				ViewCount:   taskData.WatchNum,
				RadioGroup:  taskData.RadioGroup,
				CampusGroup: taskData.CampusGroup,
				Region:      taskData.Region,
				Price:       taskData.Price,
				Wechat:      taskData.Wechat,
				Images:      s.formatImages(taskData.Images),
				Cover:       s.formatImages(taskData.Cover),
				State:       s.formatState(taskData.IsDelete, taskData.IsComplaint, taskData.Choose, taskData.Hot),
				Tag:         "未分析",
			}

			if err := db.Create(&post).Error; err != nil {
				return err
			}
			log.Printf("Created new post: %d - %s", taskData.ID, taskData.Title)
		} else if result.Error == nil {
			// 更新现有帖子
			existingPost.Title = taskData.Title
			existingPost.Content = taskData.Content
			existingPost.Author = taskData.UserName
			existingPost.AuthorID = taskData.OpenID
			existingPost.IP = taskData.IP
			existingPost.LikeNum = taskData.LikeNum
			existingPost.ReplyCount = taskData.CommentNum
			existingPost.ViewCount = taskData.WatchNum
			existingPost.RadioGroup = taskData.RadioGroup
			existingPost.CampusGroup = taskData.CampusGroup
			existingPost.Region = taskData.Region
			existingPost.Price = taskData.Price
			existingPost.Wechat = taskData.Wechat
			existingPost.Images = s.formatImages(taskData.Images)
			existingPost.Cover = s.formatImages(taskData.Cover)
			existingPost.State = s.formatState(taskData.IsDelete, taskData.IsComplaint, taskData.Choose, taskData.Hot)
			existingPost.Tag = "未分析"
			existingPost.UpdatedAt = time.Now()

			if err := db.Save(&existingPost).Error; err != nil {
				return err
			}
			log.Printf("Updated post: %d - %s", taskData.ID, taskData.Title)
		} else {
			return result.Error
		}

		return nil
	})
}

// scrapePostComments 抓取帖子的评论 - 使用批量保存
func (s *Service) scrapePostComments(postID string) error {
	comments, err := s.getComments(postID)
	if err != nil {
		return err
	}

	// 获取帖子的数据库ID
	var post models.Post
	if err := s.db.Where("original_id = ?", postID).First(&post).Error; err != nil {
		return err
	}

	// 收集所有评论（包括嵌套评论）
	var allReplies []models.Reply
	for _, comment := range comments {
		reply := s.buildReply(comment, post.ID)
		if reply != nil {
			allReplies = append(allReplies, *reply)
		}
		
		// 处理嵌套评论
		for _, nestedComment := range comment.CommentList {
			nestedReply := s.buildReply(nestedComment, post.ID)
			if nestedReply != nil {
				allReplies = append(allReplies, *nestedReply)
			}
		}
	}

	// 批量保存评论
	if len(allReplies) > 0 {
		err := database.WithRetry(s.db, func(db *gorm.DB) error {
			// 分批插入，每批50条
			return database.BatchInsert(db, allReplies, 50)
		})
		
		if err != nil {
			log.Printf("Failed to batch save comments for post %s: %v", postID, err)
			// 如果批量插入失败，尝试逐个插入
			for _, reply := range allReplies {
				s.saveCommentSingle(reply)
			}
		}
	}

	log.Printf("Saved %d comments for post %s", len(allReplies), postID)
	return nil
}

// buildReply 构建回复对象
func (s *Service) buildReply(comment CommentData, postID uint) *models.Reply {
	// 检查是否已存在
	var existingReply models.Reply
	if err := s.db.Where("original_id = ?", strconv.Itoa(comment.ID)).First(&existingReply).Error; err == nil {
		return nil // 已存在
	}

	// 获取父评论ID
	parentID := 0
	if comment.PID > 0 {
		var parentReply models.Reply
		if err := s.db.Where("original_id = ?", strconv.Itoa(comment.PID)).First(&parentReply).Error; err == nil {
			parentID = int(parentReply.ID)
		}
	}

	// 使用原子更新操作递增评论数，避免数据库busy
	err := database.WithRetry(s.db, func(db *gorm.DB) error {
		return db.Model(&models.Post{}).Where("id = ?", postID).Update("reply_count", gorm.Expr("reply_count + ?", 1)).Error
	})
	
	if err != nil {
		log.Printf("Failed to increment reply count for post %d: %v", postID, err)
		// 不返回nil，继续创建回复对象
	}

	return &models.Reply{
		PostID:     postID,
		OriginalID: strconv.Itoa(comment.ID),
		Content:    comment.Comment,
		Author:     comment.UserName,
		AuthorID:   comment.OpenID,
		ApplyTo:    comment.ApplyTo,
		Level:      s.parseLevelToInt(comment.Level),
		ParentID:   parentID,
		LikeNum:    comment.LikeNum,
		Images:     s.formatImages(comment.Images),
		Tag:        "未分析",
		CreatedAt:  s.parseTime(comment.CTime),
	}
}

// saveCommentSingle 保存单条评论（作为批量保存的后备方案）
func (s *Service) saveCommentSingle(reply models.Reply) {
	err := database.WithRetry(s.db, func(db *gorm.DB) error {
		return db.Create(&reply).Error
	})
	
	if err != nil {
		log.Printf("Failed to save reply %s: %v", reply.OriginalID, err)
	}
}

// parseTime 解析时间字符串，并统一处理时区
func (s *Service) parseTime(timeStr string) time.Time {
	// 常见时间格式
	formats := []string{
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
		"01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, time.Local); err == nil {
			return t
		}
	}

	// 如果解析失败，返回当前时间
	log.Printf("Warning: Failed to parse time: %s", timeStr)
	return time.Now()
}

// parseLevelToInt 解析 level 字段为整数
func (s *Service) parseLevelToInt(level interface{}) int {
	switch v := level.(type) {
	case string:
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}
		return 1 // 默认值
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 1 // 默认值
	}
}

// getLastReplyCheckTime 获取上次检查新回复的时间
func (s *Service) getLastReplyCheckTime() time.Time {
	var syncStatus models.SyncStatus
	err := database.WithRetry(s.db, func(db *gorm.DB) error {
		return db.Where("status != ?", "running").Order("created_at desc").First(&syncStatus).Error
	})
	if err != nil {
		return time.Time{} // 返回零值时间
	}
	return syncStatus.LastSyncTime
}

// updateLastReplyCheckTime 更新最后检查新回复的时间
func (s *Service) updateLastReplyCheckTime() {
	// 这个时间会在同步状态中自动更新，这里可以添加额外的逻辑
}

// updatePostIDAfterSync 同步帖子后更新本地ID
func (s *Service) updatePostIDAfterSync(post models.Post) error {
	// 获取该用户在主站的帖子（只需要第一个结果）
	remotePosts, err := s.getUserPosts(post.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to get user posts: %v", err)
	}

	// 由于AuthorID是唯一的，直接使用第一个结果
	if len(remotePosts) > 0 {
		remotePost := remotePosts[0] // 最新的帖子就是刚刚同步的
		
		// 更新本地帖子的original_id
		err := database.WithRetry(s.db, func(db *gorm.DB) error {
			return db.Model(&post).Update("original_id", strconv.Itoa(remotePost.ID)).Error
		})
		
		if err != nil {
			return fmt.Errorf("failed to update post original_id: %v", err)
		}
		log.Printf("Updated post ID: %d -> %d", post.ID, remotePost.ID)
		return nil
	}

	log.Printf("Warning: No remote posts found for user %s", post.AuthorID)
	return nil
}

// updateReplyIDAfterSync 同步回复后更新本地ID
func (s *Service) updateReplyIDAfterSync(reply models.Reply) error {
	// 获取该用户在主站的评论（只需要第一个结果）
	remoteComments, err := s.getUserComments(reply.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to get user comments: %v", err)
	}

	// 由于AuthorID是唯一的，直接使用第一个结果
	if len(remoteComments) > 0 {
		remoteComment := remoteComments[0] // 最新的评论就是刚刚同步的
		
		var existingReply models.Reply
		if err := s.db.Where("original_id = ?", strconv.Itoa(remoteComment.ID)).First(&existingReply).Error; err == nil {
			// 已有对应 original_id，删除当前 reply
			if err := database.WithRetry(s.db, func(db *gorm.DB) error {
				return db.Delete(&reply).Error
			}); err != nil {
				log.Printf("Failed to delete duplicate reply %d: %v", reply.ID, err)
			}
			return nil
		}
		// 更新本地回复的original_id
		err := database.WithRetry(s.db, func(db *gorm.DB) error {
			return db.Model(&reply).Update("original_id", strconv.Itoa(remoteComment.ID)).Error
		})
		
		if err != nil {
			return fmt.Errorf("failed to update reply original_id: %v", err)
		}
		log.Printf("Updated reply ID: %d -> %d", reply.ID, remoteComment.ID)
		return nil
	}

	log.Printf("Warning: No remote comments found for user %s", reply.AuthorID)
	return nil
}

// GetLastSyncStatus 获取最后同步状态
func (s *Service) GetLastSyncStatus() (*models.SyncStatus, error) {
	var status models.SyncStatus
	err := database.WithRetry(s.db, func(db *gorm.DB) error {
		return db.Order("created_at desc").First(&status).Error
	})
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// 辅助函数

// formatState 将多个状态字段转换为单一状态
func (s *Service) formatState(isDelete, isComplaint, choose, hot int) string {
	if isDelete == 1 {
		return "deleted"
	}
	if isComplaint == 1 {
		return "complaint"
	}
	if choose == 1 {
		return "chosen"
	}
	if hot == 1 {
		return "hot"
	}
	return "normal"
}

// formatImages 格式化图片字段，修复 [object Object] 问题
func (s *Service) formatImages(imgStr string) string {
	if imgStr == "" || imgStr == "[]" || imgStr == "[object Object]" {
		return "[]"
	}
	
	// 以逗号分隔
	parts := strings.Split(imgStr, ",")
	var cleaned []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && part != "[object Object]" {
			cleaned = append(cleaned, part)
		}
	}
	// 转换为 JSON 数组字符串
	jsonArray, err := json.Marshal(cleaned)
	if err != nil {
		return "[]"
	}
	return string(jsonArray)
}

// SyncPostToMainSite 同步帖子到主站
func (s *Service) SyncPostToMainSite(post models.Post) error {
	// 构建请求URL
	timeStr := url.QueryEscape(post.CreatedAt.Format("2006/01/02 15:04:05"))
	content := url.QueryEscape(post.Content)
	title := url.QueryEscape(post.Title)
	userName := url.QueryEscape(post.Author)

	syncURL := fmt.Sprintf("%s/addtask?c_time=%s&content=%s&price=&title=%s&wechat=&avatar=http%%3A%%2F%%2Fyqtech.ltd%%2Fanimal%%2F4.png&radioGroup=radio40&campusGroup=2&userName=%s&img=%%5B%%5D&cover=%%5B%%5D&region=0&likeNum=0&commentNum=0&watchNum=%d&openid=%s",
		s.baseURL, timeStr, content, title, userName, post.ViewCount, post.AuthorID)

	// 发送请求到主站（使用代理客户端）
	resp, err := s.syncClient.Get(syncURL)
	if err != nil {
		log.Printf("Failed to sync post to main site: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Post synced to main site successfully, ID: %d", post.ID)
	
	// 等待一小段时间确保主站处理完成
	// time.Sleep(2 * time.Second)
	
	// 更新本地帖子ID
	if err := s.updatePostIDAfterSync(post); err != nil {
		log.Printf("Failed to update post ID after sync: %v", err)
		// 不返回错误，因为同步已经成功
	}
	
	return nil
}

// SyncReplyToMainSite 同步回复到主站
func (s *Service) SyncReplyToMainSite(post models.Post, reply models.Reply) error {
	// 构建请求URL
	timeStr := url.QueryEscape(reply.CreatedAt.Format("2006/01/02 15:04:05"))
	content := url.QueryEscape(reply.Content)
	userName := url.QueryEscape(reply.Author)

	// 使用帖子的original_id作为pk
	pk := post.OriginalID
	if pk == "" {
		pk = fmt.Sprintf("%d", post.ID)
	}

	pid := 0
	if reply.ParentID > 0 {
		// 根据pid找对应originalid
		var parentReply models.Reply
		err := s.db.First(&parentReply, reply.ParentID).Error
		if err != nil {
			log.Printf("Parent reply not found for reply %d, skipping sync.", reply.ID)
			return nil // 如果没找到父评论，则不同步
		}
		pid, _ = strconv.Atoi(parentReply.OriginalID)
	}

	syncURL := fmt.Sprintf("%s/addcomment?c_time=%s&openid=%s&pk=%s&comment=%s&userName=%s&avatar=http%%3A%%2F%%2Fyqtech.ltd%%2Fanimal%%2F4.png&applyTo=%s&img=%%5B%%5D&level=%d&pid=%d",
		s.baseURL, timeStr, reply.AuthorID, pk, content, userName, reply.ApplyTo, reply.Level, pid)

	// 发送请求到主站（使用代理客户端）
	resp, err := s.client.Get(syncURL)
	if err != nil {
		log.Printf("Failed to sync reply to main site: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Reply synced to main site successfully, ID: %d", reply.ID)
	
	// 等待一小段时间确保主站处理完成
	// time.Sleep(2 * time.Second)
	
	// 更新本地回复ID
	if err := s.updateReplyIDAfterSync(reply); err != nil {
		log.Printf("Failed to update reply ID after sync: %v", err)
		// 不返回错误，因为同步已经成功
	}
	
	return nil
}

