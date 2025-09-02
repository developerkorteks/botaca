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
}

// NewAutoPromoteService membuat service baru
func NewAutoPromoteService(client *whatsmeow.Client, repo database.Repository, logger *utils.Logger) *AutoPromoteService {
	service := &AutoPromoteService{
		client:     client,
		repository: repo,
		logger:     logger,
		isRunning:  false,
	}
	
	// Inisialisasi scheduler
	service.scheduler = NewSchedulerService(service.processScheduledPromotes, logger)
	
	return service
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
	s.scheduler.Start(4 * time.Hour) // Setiap 4 jam
	s.isRunning = true
	s.logger.Success("Auto promote scheduler started!")
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

// processScheduledPromotes memproses promosi terjadwal
func (s *AutoPromoteService) processScheduledPromotes() {
	s.logger.Info("Processing scheduled promotes...")
	
	// Ambil semua grup yang aktif
	activeGroups, err := s.repository.GetActiveGroups()
	if err != nil {
		s.logger.Errorf("Failed to get active groups: %v", err)
		return
	}
	
	if len(activeGroups) == 0 {
		s.logger.Info("No active groups for auto promote")
		return
	}
	
	s.logger.Infof("Found %d active groups", len(activeGroups))
	
	// Ambil template aktif
	templates, err := s.repository.GetActiveTemplates()
	if err != nil {
		s.logger.Errorf("Failed to get templates: %v", err)
		return
	}
	
	if len(templates) == 0 {
		s.logger.Warning("No active templates available")
		return
	}
	
	// Proses setiap grup
	successCount := 0
	failCount := 0
	
	for _, group := range activeGroups {
		// Cek apakah sudah waktunya untuk promote (4 jam sejak terakhir)
		if s.shouldSkipGroup(&group) {
			continue
		}
		
		// Kirim promosi
		err := s.sendPromoteToGroup(group.GroupJID, templates)
		if err != nil {
			s.logger.Errorf("Failed to send promote to group %s: %v", group.GroupJID, err)
			failCount++
		} else {
			successCount++
			
			// Update last promote time
			now := time.Now()
			group.LastPromoteAt = &now
			s.repository.UpdateAutoPromoteGroup(&group)
		}
	}
	
	s.logger.Infof("Scheduled promotes completed: %d success, %d failed", successCount, failCount)
	
	// Update statistik
	today := time.Now().Format("2006-01-02")
	s.repository.UpdateStats(today, len(activeGroups), successCount+failCount, successCount, failCount)
}

// shouldSkipGroup mengecek apakah grup harus dilewati
func (s *AutoPromoteService) shouldSkipGroup(group *database.AutoPromoteGroup) bool {
	// Jika belum pernah kirim promosi, kirim sekarang
	if group.LastPromoteAt == nil {
		return false
	}
	
	// Cek apakah sudah 4 jam sejak promosi terakhir
	fourHoursAgo := time.Now().Add(-4 * time.Hour)
	return group.LastPromoteAt.After(fourHoursAgo)
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
	
	// Seed random generator
	rand.Seed(time.Now().UnixNano())
	
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