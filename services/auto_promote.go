// Package services - Auto promote service untuk mengelola promosi otomatis
package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"

	"github.com/nabilulilalbab/promote/database"
	"github.com/nabilulilalbab/promote/utils"
)

// AutoPromoteService mengelola fitur auto promote
type AutoPromoteService struct {
	client     *whatsmeow.Client
	repository database.Repository
	logger     *utils.Logger
	scheduler  *SchedulerService
	isRunning  bool
	interval   time.Duration // Interval auto promote dalam durasi
}

// NewAutoPromoteService membuat service baru
func NewAutoPromoteService(client *whatsmeow.Client, repo database.Repository, logger *utils.Logger) *AutoPromoteService {
	// Inisialisasi random seed sekali saja
	rand.Seed(time.Now().UnixNano())

	service := &AutoPromoteService{
		client:     client,
		repository: repo,
		logger:     logger,
		isRunning:  false,
		interval:   4 * time.Hour, // Default 4 jam
	}
	
	// Inisialisasi scheduler
	service.scheduler = NewSchedulerService(service.processScheduledPromotes, logger)
	
	return service
}

// SetInterval mengatur interval auto promote
func (s *AutoPromoteService) SetInterval(hours int) {
	s.interval = time.Duration(hours) * time.Hour
	s.logger.Infof("Auto promote interval set to %d hours", hours)
}

// StartAutoPromote mengaktifkan auto promote untuk grup tertentu
func (s *AutoPromoteService) StartAutoPromote(groupJID string) error {
	s.logger.Infof("Starting auto promote for group: %s", groupJID)
	
	// Cek apakah grup sudah ada di database
	group, err := s.repository.GetAutoPromoteGroup(groupJID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	
	// Jika grup belum ada, buat baru
	if group == nil {
		group, err = s.repository.CreateAutoPromoteGroup(groupJID)
		if err != nil {
			return fmt.Errorf("failed to create group: %v", err)
		}
	}
	
	// Jika sudah aktif, return error
	if group.IsActive {
		return fmt.Errorf("auto promote sudah aktif untuk grup ini")
	}
	
	// Aktifkan auto promote
	now := time.Now()
	group.IsActive = true
	group.StartedAt = &now
	
	err = s.repository.UpdateAutoPromoteGroup(group)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}
	
	// Start scheduler jika belum berjalan
	if !s.isRunning {
		s.StartScheduler()
	}
	
	s.logger.Successf("Auto promote activated for group: %s", groupJID)
	return nil
}

// StopAutoPromote menghentikan auto promote untuk grup tertentu
func (s *AutoPromoteService) StopAutoPromote(groupJID string) error {
	s.logger.Infof("Stopping auto promote for group: %s", groupJID)
	
	// Cek apakah grup ada di database
	group, err := s.repository.GetAutoPromoteGroup(groupJID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	
	if group == nil {
		return fmt.Errorf("grup tidak ditemukan dalam sistem auto promote")
	}
	
	if !group.IsActive {
		return fmt.Errorf("auto promote tidak aktif untuk grup ini")
	}
	
	// Nonaktifkan auto promote
	group.IsActive = false
	group.StartedAt = nil
	
	err = s.repository.UpdateAutoPromoteGroup(group)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}
	
	s.logger.Successf("Auto promote deactivated for group: %s", groupJID)
	return nil
}

// GetGroupStatus mendapatkan status auto promote untuk grup
func (s *AutoPromoteService) GetGroupStatus(groupJID string) (*database.AutoPromoteGroup, error) {
	return s.repository.GetAutoPromoteGroup(groupJID)
}

// StartScheduler memulai scheduler untuk auto promote
func (s *AutoPromoteService) StartScheduler() {
	if s.isRunning {
		return
	}
	
	s.logger.Info("Starting auto promote scheduler...")
	s.logger.Infof("Scheduler will run every %v", s.interval)
	s.scheduler.Start(s.interval)
	s.isRunning = true
	s.logger.Successf("Auto promote scheduler started with %v interval!", s.interval)
}

// StopScheduler menghentikan scheduler
func (s *AutoPromoteService) StopScheduler() {
	if !s.isRunning {
		return
	}
	
	s.logger.Info("Stopping auto promote scheduler...")
	s.scheduler.Stop()
	s.isRunning = false
	s.logger.Success("Auto promote scheduler stopped!")
}

// processScheduledPromotes memproses promosi terjadwal dengan error handling yang robust
func (s *AutoPromoteService) processScheduledPromotes() {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("Scheduler panic recovered: %v", r)
		}
	}()

	s.logger.Info("Processing scheduled promotes...")
	
	// Ambil semua grup yang aktif dengan retry mechanism
	activeGroups, err := s.getActiveGroupsWithRetry(3)
	if err != nil {
		s.logger.Errorf("Failed to get active groups after retries: %v", err)
		return
	}
	
	if len(activeGroups) == 0 {
		s.logger.Info("No active groups for auto promote")
		return
	}
	
	s.logger.Infof("Found %d active groups", len(activeGroups))
	
	// Ambil template aktif dengan retry mechanism
	templates, err := s.getActiveTemplatesWithRetry(3)
	if err != nil {
		s.logger.Errorf("Failed to get templates after retries: %v", err)
		return
	}
	
	if len(templates) == 0 {
		s.logger.Warning("No active templates available")
		return
	}
	
	s.logger.Infof("Found %d active templates", len(templates))
	
	// Proses setiap grup dengan error handling individual
	successCount := 0
	failCount := 0
	skippedCount := 0
	
	for _, group := range activeGroups {
		// Cek apakah sudah waktunya untuk promote (4 jam sejak terakhir)
		if s.shouldSkipGroup(&group) {
			skippedCount++
			s.logger.Debugf("Skipping group %s (not yet time)", group.GroupJID)
			continue
		}
		
		// Kirim promosi dengan retry mechanism
		err := s.sendPromoteToGroupWithRetry(group.GroupJID, templates, 2)
		if err != nil {
			s.logger.Errorf("Failed to send promote to group %s after retries: %v", group.GroupJID, err)
			failCount++
		} else {
			successCount++
			
			// Update last promote time dengan error handling
			now := time.Now()
			group.LastPromoteAt = &now
			updateErr := s.repository.UpdateAutoPromoteGroup(&group)
			if updateErr != nil {
				s.logger.Errorf("Failed to update group %s last promote time: %v", group.GroupJID, updateErr)
			}
		}
	}
	
	s.logger.Infof("Scheduled promotes completed: %d success, %d failed, %d skipped", successCount, failCount, skippedCount)
	
	// Update statistik dengan error handling
	today := time.Now().Format("2006-01-02")
	statsErr := s.repository.UpdateStats(today, len(activeGroups), successCount+failCount, successCount, failCount)
	if statsErr != nil {
		s.logger.Errorf("Failed to update stats: %v", statsErr)
	}
}

// shouldSkipGroup mengecek apakah grup harus dilewati
func (s *AutoPromoteService) shouldSkipGroup(group *database.AutoPromoteGroup) bool {
	// Jika belum pernah kirim promosi, kirim sekarang
	if group.LastPromoteAt == nil {
		return false
	}
	
	// Cek apakah sudah mencapai interval yang ditentukan sejak promosi terakhir
	intervalAgo := time.Now().Add(-s.interval)
	return group.LastPromoteAt.After(intervalAgo)
}

// sendPromoteToGroup mengirim promosi ke grup tertentu
func (s *AutoPromoteService) sendPromoteToGroup(groupJID string, templates []database.PromoteTemplate) error {
	// Pilih template secara random
	template := s.selectRandomTemplate(templates)
	
	// Parse JID grup
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %v", err)
	}
	
	// Proses template (replace variables)
	content := s.processTemplate(template.Content, jid)
	
	// Kirim pesan
	err = s.sendMessage(jid, content)
	
	// Log hasil
	log := &database.PromoteLog{
		GroupJID:   groupJID,
		TemplateID: template.ID,
		Content:    content,
		SentAt:     time.Now(),
		Success:    err == nil,
	}
	
	if err != nil {
		errorMsg := err.Error()
		log.ErrorMsg = &errorMsg
	}
	
	s.repository.CreateLog(log)
	
	return err
}

// selectRandomTemplate memilih template secara random
func (s *AutoPromoteService) selectRandomTemplate(templates []database.PromoteTemplate) database.PromoteTemplate {
	if len(templates) == 0 {
		// Return empty template jika tidak ada
		return database.PromoteTemplate{}
	}
	
	// Pilih index random
	index := rand.Intn(len(templates))
	return templates[index]
}

// processTemplate memproses template dengan mengganti variables
func (s *AutoPromoteService) processTemplate(content string, groupJID types.JID) string {
	now := time.Now()
	
	// Replace variables yang tersedia
	replacements := map[string]string{
		"{DATE}":     now.Format("2006-01-02"),
		"{TIME}":     now.Format("15:04"),
		"{DAY}":      getDayName(now.Weekday()),
		"{MONTH}":    getMonthName(now.Month()),
		"{YEAR}":     fmt.Sprintf("%d", now.Year()),
		"{GROUP_ID}": groupJID.User,
	}
	
	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	return result
}

// sendMessage mengirim pesan ke grup
func (s *AutoPromoteService) sendMessage(groupJID types.JID, content string) error {
	// Buat pesan WhatsApp
	msg := &waProto.Message{
		Conversation: &content,
	}
	
	// Kirim pesan
	_, err := s.client.SendMessage(context.Background(), groupJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	
	s.logger.Infof("Promote message sent to group: %s", groupJID.String())
	return nil
}

// SendManualPromote mengirim promosi manual (untuk testing)
func (s *AutoPromoteService) SendManualPromote(groupJID string) error {
	// Ambil template aktif
	templates, err := s.repository.GetActiveTemplates()
	if err != nil {
		return fmt.Errorf("failed to get templates: %v", err)
	}
	
	if len(templates) == 0 {
		return fmt.Errorf("no active templates available")
	}
	
	// Kirim promosi
	return s.sendPromoteToGroup(groupJID, templates)
}

// GetActiveGroupsCount mendapatkan jumlah grup aktif
func (s *AutoPromoteService) GetActiveGroupsCount() (int, error) {
	groups, err := s.repository.GetActiveGroups()
	if err != nil {
		return 0, err
	}
	return len(groups), nil
}

// GetActiveGroups mendapatkan daftar grup aktif
func (s *AutoPromoteService) GetActiveGroups() ([]database.AutoPromoteGroup, error) {
	return s.repository.GetActiveGroups()
}

// getActiveGroupsWithRetry mengambil grup aktif dengan retry mechanism
func (s *AutoPromoteService) getActiveGroupsWithRetry(maxRetries int) ([]database.AutoPromoteGroup, error) {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		groups, err := s.repository.GetActiveGroups()
		if err == nil {
			return groups, nil
		}
		
		lastErr = err
		s.logger.Warningf("Retry %d/%d getting active groups failed: %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
		}
	}
	
	return nil, lastErr
}

// getActiveTemplatesWithRetry mengambil template aktif dengan retry mechanism
func (s *AutoPromoteService) getActiveTemplatesWithRetry(maxRetries int) ([]database.PromoteTemplate, error) {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		templates, err := s.repository.GetActiveTemplates()
		if err == nil {
			return templates, nil
		}
		
		lastErr = err
		s.logger.Warningf("Retry %d/%d getting active templates failed: %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
		}
	}
	
	return nil, lastErr
}

// sendPromoteToGroupWithRetry mengirim promosi dengan retry mechanism
func (s *AutoPromoteService) sendPromoteToGroupWithRetry(groupJID string, templates []database.PromoteTemplate, maxRetries int) error {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		err := s.sendPromoteToGroup(groupJID, templates)
		if err == nil {
			return nil
		}
		
		lastErr = err
		s.logger.Warningf("Retry %d/%d sending promote to %s failed: %v", i+1, maxRetries, groupJID, err)
		
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second) // Longer backoff for network issues
		}
	}
	
	return lastErr
}

// Helper functions

func getDayName(day time.Weekday) string {
	days := []string{
		"Minggu", "Senin", "Selasa", "Rabu", 
		"Kamis", "Jumat", "Sabtu",
	}
	return days[day]
}

func getMonthName(month time.Month) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return months[month]
}