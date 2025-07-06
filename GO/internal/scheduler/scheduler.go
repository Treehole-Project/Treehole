package scheduler

import (
	"fmt"
	"log"
	"os"
	"sync"
	"treehole/internal/scraper"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron           *cron.Cron
	scraperService *scraper.Service
	isRunning      bool      // 标记是否有爬取任务正在运行
	mutex          sync.Mutex // 保护 isRunning 状态的互斥锁
}

// New 创建新的调度器
func New(scraperService *scraper.Service) *Scheduler {
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:           c,
		scraperService: scraperService,
		isRunning:      false,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	if os.Getenv("INBOUND_SYNC_ENABLED") != "true" {
		log.Println("INBOUND_SYNC_ENABLED is not set to true, scheduler will not start")
		return
	}
	// 从环境变量获取同步间隔，默认为每30分钟
	cronSpec := os.Getenv("SYNC_CRON")
	if cronSpec == "" {
		cronSpec = "0 */30 * * * *" // 每30分钟执行一次
	}

	sourceURL := os.Getenv("SOURCE_URL")
	if sourceURL == "" {
		log.Println("Warning: SOURCE_URL not set, scheduler will not run automatic sync")
		return
	}

	// 添加定时同步任务
	_, err := s.cron.AddFunc(cronSpec, func() {
		// 检查是否已有任务在运行
		s.mutex.Lock()
		if s.isRunning {
			log.Println("Previous sync job is still running, skipping this execution")
			s.mutex.Unlock()
			return
		}
		s.isRunning = true
		s.mutex.Unlock()

		// 在函数结束时重置运行状态
		defer func() {
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}()

		log.Println("Starting scheduled sync...")
		if err := s.scraperService.ScrapeData(); err != nil {
			log.Printf("Scheduled sync failed: %v", err)
		} else {
			log.Println("Scheduled sync completed successfully")
		}
	})

	if err != nil {
		log.Printf("Failed to add sync job: %v", err)
		return
	}

	s.cron.Start()
	log.Printf("Scheduler started with cron spec: %s", cronSpec)
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// AddJob 添加自定义任务
func (s *Scheduler) AddJob(spec string, job func()) error {
	_, err := s.cron.AddFunc(spec, job)
	return err
}

// IsRunning 检查是否有同步任务正在运行
func (s *Scheduler) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isRunning
}

// TriggerSync 手动触发同步（如果没有任务在运行）
func (s *Scheduler) TriggerSync() error {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		return fmt.Errorf("sync job is already running")
	}
	s.isRunning = true
	s.mutex.Unlock()

	// 在函数结束时重置运行状态
	defer func() {
		s.mutex.Lock()
		s.isRunning = false
		s.mutex.Unlock()
	}()

	log.Println("Starting manual sync...")
	if err := s.scraperService.ScrapeData(); err != nil {
		log.Printf("Manual sync failed: %v", err)
		return err
	}

	log.Println("Manual sync completed successfully")
	return nil
}
