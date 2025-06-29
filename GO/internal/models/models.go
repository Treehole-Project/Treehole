package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 树洞帖子模型
type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OriginalID  string         `json:"original_id" gorm:"not null"`
	Title       string         `json:"title"`
	Content     string         `json:"content" gorm:"type:text"`
	Author      string         `json:"author"`
	AuthorID    string         `json:"author_id"` // openid
	IP          string         `json:"ip"`
	LikeNum     int            `json:"like_num" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	ReplyCount  int            `json:"reply_count" gorm:"default:0"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	RadioGroup  string         `json:"radio_group"` // 帖子分组
	CampusGroup string         `json:"campus_group"` // 校区分组
	Region      string         `json:"region"`
	Price       string         `json:"price"`
	Wechat      string         `json:"wechat"`
	Images      string         `json:"images" gorm:"type:text"` // JSON 格式存储图片URL列表
	Cover       string         `json:"cover"`
	State       string         `json:"state" gorm:"default:normal"` // normal, deleted, complaint, chosen, hot
	Tag         string         `json:"tag"` // 标签
	Replies     []Reply        `json:"replies,omitempty"`
}

// Reply 回复模型
type Reply struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	PostID     uint           `json:"post_id" gorm:"not null"`
	OriginalID string         `json:"original_id" gorm:"not null"`
	Content    string         `json:"content" gorm:"type:text"`
	Author     string         `json:"author"`
	AuthorID   string         `json:"author_id"` // openid
	ApplyTo    string         `json:"apply_to"`  // 回复给谁的 openid
	Level      int            `json:"level" gorm:"default:1"` // 回复层级
	ParentID   int            `json:"parent_id" gorm:"default:0"` // 父评论ID (pid)
	LikeNum    int            `json:"like_num" gorm:"default:0"`
	Images     string         `json:"images" gorm:"type:text"` // JSON 格式存储图片URL列表
	Tag        string         `json:"tag"` // 标签
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Post       Post           `json:"-" gorm:"foreignKey:PostID"`
}

// SyncStatus 同步状态模型
type SyncStatus struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	LastSyncTime time.Time `json:"last_sync_time"`
	LastPostID   string    `json:"last_post_id"`
	TotalPosts   int       `json:"total_posts"`
	TotalReplies int       `json:"total_replies"`
	Status       string    `json:"status"` // "success", "error", "running"
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
