package database

import (
	"strings"
	"time"
	"treehole/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxRetries: 5,
	BaseDelay:  100 * time.Millisecond,
	MaxDelay:   5 * time.Second,
}

// InitDB 初始化数据库连接
func InitDB(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 根据数据库 URL 判断数据库类型
	if strings.Contains(databaseURL, "mysql://") || strings.Contains(databaseURL, "@tcp(") {
		// MySQL 数据库
		db, err = gorm.Open(mysql.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		// SQLite 数据库 (默认) - 使用纯 Go 驱动
		db, err = gorm.Open(sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        databaseURL,
		}, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		
		if err != nil {
			return nil, err
		}

		// 优化 SQLite 配置以支持高并发
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		// 设置连接池参数
		sqlDB.SetMaxOpenConns(1)     // SQLite 只支持一个写连接
		sqlDB.SetMaxIdleConns(1)     // 保持一个空闲连接
		sqlDB.SetConnMaxLifetime(0)  // 连接不过期

		// 执行 SQLite 优化设置
		db.Exec("PRAGMA journal_mode=WAL")           // 启用 WAL 模式提高并发性能
		db.Exec("PRAGMA synchronous=NORMAL")         // 平衡安全性和性能
		db.Exec("PRAGMA cache_size=10000")           // 增加缓存大小
		db.Exec("PRAGMA temp_store=memory")          // 临时文件存储在内存中
		db.Exec("PRAGMA mmap_size=268435456")        // 启用内存映射 (256MB)
		db.Exec("PRAGMA busy_timeout=30000")         // 设置繁忙超时为30秒
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	// 对于SQLite，我们需要更仔细地处理迁移
	// 首先检查表是否存在，如果存在则尝试兼容性迁移
	
	// 检查posts表是否存在
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	
	if count > 0 {
		// 表已存在，执行兼容性检查和迁移
		return migrateExistingTables(db)
	}
	
	// 表不存在，执行标准迁移
	return db.AutoMigrate(
		&models.Post{},
		&models.Reply{},
		&models.SyncStatus{},
	)
}

// migrateExistingTables 迁移现有表
func migrateExistingTables(db *gorm.DB) error {
	// 对于现有的表，我们不执行AutoMigrate，而是检查字段兼容性
	// 如果需要，可以在这里添加ALTER TABLE语句
	
	// 检查posts表结构
	var tableInfo []struct {
		CID        int    `gorm:"column:cid"`
		Name       string `gorm:"column:name"`
		Type       string `gorm:"column:type"`
		NotNull    int    `gorm:"column:notnull"`
		DefaultVal string `gorm:"column:dflt_value"`
		PK         int    `gorm:"column:pk"`
	}
	
	db.Raw("PRAGMA table_info(posts)").Scan(&tableInfo)
	
	// 检查是否存在必要的字段，如果不存在则添加
	hasStateField := false
	for _, field := range tableInfo {
		if field.Name == "state" {
			hasStateField = true
			break
		}
	}
	
	// 如果state字段不存在，添加它
	if !hasStateField {
		if err := db.Exec("ALTER TABLE posts ADD COLUMN state TEXT DEFAULT 'normal'").Error; err != nil {
			return err
		}
	}
	
	// 确保replies表存在
	var replyCount int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='replies'").Scan(&replyCount)
	if replyCount == 0 {
		if err := db.AutoMigrate(&models.Reply{}); err != nil {
			return err
		}
	}
	
	// 确保sync_statuses表存在
	var syncCount int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sync_statuses'").Scan(&syncCount)
	if syncCount == 0 {
		if err := db.AutoMigrate(&models.SyncStatus{}); err != nil {
			return err
		}
	}
	
	return nil
}

// WithRetry 使用重试机制执行数据库操作
func WithRetry(db *gorm.DB, operation func(*gorm.DB) error) error {
	return WithRetryConfig(db, operation, DefaultRetryConfig)
}

// WithRetryConfig 使用自定义重试配置执行数据库操作
func WithRetryConfig(db *gorm.DB, operation func(*gorm.DB) error, config RetryConfig) error {
	var err error
	delay := config.BaseDelay

	for i := 0; i <= config.MaxRetries; i++ {
		err = operation(db)
		if err == nil {
			return nil
		}

		// 检查是否是数据库锁定错误
		if !isDatabaseBusyError(err) {
			return err
		}

		// 最后一次重试失败
		if i == config.MaxRetries {
			break
		}

		// 等待后重试
		time.Sleep(delay)
		
		// 指数退避，但不超过最大延迟
		delay *= 2
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return err
}

// isDatabaseBusyError 检查是否是数据库繁忙错误
func isDatabaseBusyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "database is locked") || 
		   strings.Contains(errStr, "sqlite_busy") ||
		   strings.Contains(errStr, "database locked")
}

// BatchInsert 批量插入数据
func BatchInsert[T any](db *gorm.DB, items []T, batchSize int) error {
	if len(items) == 0 {
		return nil
	}

	// 分批处理
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		err := WithRetry(db, func(db *gorm.DB) error {
			return db.CreateInBatches(batch, len(batch)).Error
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// SafeTransaction 安全的事务操作
func SafeTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return WithRetry(db, func(db *gorm.DB) error {
		return db.Transaction(fn)
	})
}

// 添加数据库连接安全配置
func configureDatabaseSecurity(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	// 设置连接池安全参数
	sqlDB.SetMaxOpenConns(25)                 // 限制最大连接数
	sqlDB.SetMaxIdleConns(5)                  // 限制空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // 空闲连接最大时间
	
	return nil
}
