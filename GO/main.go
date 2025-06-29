package main

import (
	"log"
	"os"

	"treehole/internal/api"
	"treehole/internal/config"
	"treehole/internal/database"
	"treehole/internal/scheduler"
	"treehole/internal/scraper"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化爬虫
	scraperService := scraper.NewService(db, cfg)

	// 启动定时任务
	scheduler := scheduler.New(scraperService)
	scheduler.Start()
	defer scheduler.Stop()

	// 启动 API 服务器
	router := api.SetupRouter(db, scraperService)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
