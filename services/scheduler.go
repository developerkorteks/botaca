// Package services - Scheduler service untuk menjalankan auto promote secara berkala
package services

import (
	"sync"
	"time"

	"github.com/nabilulilalbab/promote/utils"
)

// SchedulerService mengelola penjadwalan auto promote
type SchedulerService struct {
	ticker    *time.Ticker
	done      chan bool
	task      func()
	logger    *utils.Logger
	isRunning bool
	mutex     sync.RWMutex
}

// NewSchedulerService membuat scheduler baru
func NewSchedulerService(task func(), logger *utils.Logger) *SchedulerService {
	return &SchedulerService{
		task:      task,
		logger:    logger,
		done:      make(chan bool),
		isRunning: false,
	}
}

// Start memulai scheduler dengan interval tertentu
func (s *SchedulerService) Start(interval time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		s.logger.Warning("Scheduler already running")
		return
	}

	s.logger.Infof("Starting scheduler with interval: %v", interval)
	
	s.ticker = time.NewTicker(interval)
	s.isRunning = true

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.logger.Debug("Scheduler tick - executing task")
				s.executeTask()
			case <-s.done:
				s.logger.Info("Scheduler stopped")
				return
			}
		}
	}()

	s.logger.Success("Scheduler started successfully")
}

// Stop menghentikan scheduler
func (s *SchedulerService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		s.logger.Warning("Scheduler not running")
		return
	}

	s.logger.Info("Stopping scheduler...")
	
	s.ticker.Stop()
	s.done <- true
	s.isRunning = false

	s.logger.Success("Scheduler stopped successfully")
}

// IsRunning mengecek apakah scheduler sedang berjalan
func (s *SchedulerService) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isRunning
}

// executeTask menjalankan task dengan error handling
func (s *SchedulerService) executeTask() {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("Scheduler task panic: %v", r)
		}
	}()

	start := time.Now()
	s.task()
	duration := time.Since(start)
	
	s.logger.Debugf("Task completed in %v", duration)
}

// Restart menghentikan dan memulai ulang scheduler
func (s *SchedulerService) Restart(interval time.Duration) {
	s.Stop()
	time.Sleep(100 * time.Millisecond) // Small delay to ensure clean stop
	s.Start(interval)
}

// GetStatus mendapatkan status scheduler
func (s *SchedulerService) GetStatus() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status := map[string]interface{}{
		"is_running": s.isRunning,
		"has_ticker": s.ticker != nil,
	}

	return status
}