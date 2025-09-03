// Package services - Group Manager service untuk mengelola grup WhatsApp
package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/nabilulilalbab/promote/database"
	"github.com/nabilulilalbab/promote/utils"
)

// GroupInfo berisi informasi grup yang diikuti bot
type GroupInfo struct {
	ID          int    `json:"id"`
	JID         string `json:"jid"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	MemberCount int    `json:"member_count"`
	Description string `json:"description"`
}

// GroupManagerService mengelola grup-grup yang diikuti bot
type GroupManagerService struct {
	client     *whatsmeow.Client
	repository database.Repository
	logger     *utils.Logger
}

// NewGroupManagerService membuat service baru
func NewGroupManagerService(client *whatsmeow.Client, repo database.Repository, logger *utils.Logger) *GroupManagerService {
	// Inisialisasi random seed untuk template selection
	rand.Seed(time.Now().UnixNano())

	return &GroupManagerService{
		client:     client,
		repository: repo,
		logger:     logger,
	}
}

// GetAllJoinedGroups mengambil semua grup yang diikuti bot dari WhatsApp
func (s *GroupManagerService) GetAllJoinedGroups() ([]GroupInfo, error) {
	s.logger.Info("Getting all joined groups from WhatsApp...")

	// Ambil semua grup dari WhatsApp client
	groups, err := s.client.GetJoinedGroups()
	if err != nil {
		s.logger.Errorf("Failed to get joined groups: %v", err)
		return nil, fmt.Errorf("failed to get joined groups: %v", err)
	}

	var groupInfos []GroupInfo

	for i, group := range groups {
		// group adalah *types.GroupInfo, bukan JID
		groupJID := group.JID

		// Cek status auto promote dari database
		dbGroup, err := s.repository.GetAutoPromoteGroup(groupJID.String())
		isActive := false
		if err == nil && dbGroup != nil {
			isActive = dbGroup.IsActive
		}

		// Format nama grup
		groupName := "Unnamed Group"
		if group.Name != "" {
			groupName = group.Name
		}

		groupInfos = append(groupInfos, GroupInfo{
			ID:          i + 1,
			JID:         groupJID.String(),
			Name:        groupName,
			IsActive:    isActive,
			MemberCount: len(group.Participants),
			Description: group.Topic,
		})
	}

	s.logger.Infof("Found %d joined groups", len(groupInfos))
	return groupInfos, nil
}

// GetGroupByID mengambil info grup berdasarkan ID (index)
func (s *GroupManagerService) GetGroupByID(groupID int) (*GroupInfo, error) {
	groups, err := s.GetAllJoinedGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.ID == groupID {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("grup dengan ID %d tidak ditemukan", groupID)
}

// EnableAutoPromoteForGroup mengaktifkan auto promote untuk grup tertentu
func (s *GroupManagerService) EnableAutoPromoteForGroup(groupID int) error {
	// Ambil info grup
	groupInfo, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	s.logger.Infof("Enabling auto promote for group: %s (%s)", groupInfo.Name, groupInfo.JID)

	// Cek apakah grup sudah ada di database
	dbGroup, err := s.repository.GetAutoPromoteGroup(groupInfo.JID)
	if err != nil {
		return fmt.Errorf("failed to get group from database: %v", err)
	}

	// Jika grup belum ada, buat baru
	if dbGroup == nil {
		dbGroup, err = s.repository.CreateAutoPromoteGroup(groupInfo.JID)
		if err != nil {
			return fmt.Errorf("failed to create group in database: %v", err)
		}
	}

	// Jika sudah aktif, return error
	if dbGroup.IsActive {
		return fmt.Errorf("auto promote sudah aktif untuk grup %s", groupInfo.Name)
	}

	// Aktifkan auto promote
	dbGroup.IsActive = true
	now := time.Now()
	dbGroup.StartedAt = &now

	err = s.repository.UpdateAutoPromoteGroup(dbGroup)
	if err != nil {
		return fmt.Errorf("failed to update group in database: %v", err)
	}

	s.logger.Successf("Auto promote enabled for group: %s", groupInfo.Name)
	return nil
}

// DisableAutoPromoteForGroup menonaktifkan auto promote untuk grup tertentu
func (s *GroupManagerService) DisableAutoPromoteForGroup(groupID int) error {
	// Ambil info grup
	groupInfo, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	s.logger.Infof("Disabling auto promote for group: %s (%s)", groupInfo.Name, groupInfo.JID)

	// Cek apakah grup ada di database
	dbGroup, err := s.repository.GetAutoPromoteGroup(groupInfo.JID)
	if err != nil {
		return fmt.Errorf("failed to get group from database: %v", err)
	}

	if dbGroup == nil {
		return fmt.Errorf("grup %s tidak ditemukan dalam sistem auto promote", groupInfo.Name)
	}

	if !dbGroup.IsActive {
		return fmt.Errorf("auto promote tidak aktif untuk grup %s", groupInfo.Name)
	}

	// Nonaktifkan auto promote
	dbGroup.IsActive = false
	dbGroup.StartedAt = nil

	err = s.repository.UpdateAutoPromoteGroup(dbGroup)
	if err != nil {
		return fmt.Errorf("failed to update group in database: %v", err)
	}

	s.logger.Successf("Auto promote disabled for group: %s", groupInfo.Name)
	return nil
}

// GetGroupStatus mengambil status auto promote untuk grup tertentu
func (s *GroupManagerService) GetGroupStatus(groupID int) (*GroupInfo, *database.AutoPromoteGroup, error) {
	// Ambil info grup
	groupInfo, err := s.GetGroupByID(groupID)
	if err != nil {
		return nil, nil, err
	}

	// Ambil status dari database
	dbGroup, err := s.repository.GetAutoPromoteGroup(groupInfo.JID)
	if err != nil {
		return groupInfo, nil, fmt.Errorf("failed to get group status: %v", err)
	}

	return groupInfo, dbGroup, nil
}

// SendTestPromoteToGroup mengirim test promosi ke grup tertentu
func (s *GroupManagerService) SendTestPromoteToGroup(groupID int) error {
	// Ambil info grup
	groupInfo, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	s.logger.Infof("Sending test promote to group: %s (%s)", groupInfo.Name, groupInfo.JID)

	// Ambil template aktif
	templates, err := s.repository.GetActiveTemplates()
	if err != nil {
		return fmt.Errorf("failed to get templates: %v", err)
	}

	if len(templates) == 0 {
		return fmt.Errorf("no active templates available")
	}

	// Pilih template secara random untuk test (sama seperti auto promote)
	template := s.selectRandomTemplate(templates)

	// Parse JID grup
	jid, err := types.ParseJID(groupInfo.JID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %v", err)
	}

	// Proses template (replace variables)
	content := s.processTemplate(template.Content, jid)

	// Kirim pesan promosi natural (tanpa embel-embel test)
	err = s.sendMessage(jid, content)
	if err != nil {
		return fmt.Errorf("failed to send test message: %v", err)
	}

	s.logger.Successf("Test promote sent to group: %s", groupInfo.Name)
	return nil
}

// processTemplate memproses template dengan mengganti variables
func (s *GroupManagerService) processTemplate(content string, groupJID types.JID) string {
	now := time.Now()

	// Replace variables yang tersedia
	replacements := map[string]string{
		"{DATE}":     now.Format("2006-01-02"),
		"{TIME}":     now.Format("15:04"),
		"{DAY}":      s.getDayName(now.Weekday()),
		"{MONTH}":    s.getMonthName(now.Month()),
		"{YEAR}":     fmt.Sprintf("%d", now.Year()),
		"{GROUP_ID}": groupJID.User,
	}

	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// getDayName helper function
func (s *GroupManagerService) getDayName(day time.Weekday) string {
	days := []string{
		"Minggu", "Senin", "Selasa", "Rabu",
		"Kamis", "Jumat", "Sabtu",
	}
	return days[day]
}

// getMonthName helper function
func (s *GroupManagerService) getMonthName(month time.Month) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return months[month]
}

// sendMessage mengirim pesan ke grup
func (s *GroupManagerService) sendMessage(groupJID types.JID, content string) error {
	// Buat pesan WhatsApp
	msg := &waProto.Message{
		Conversation: &content,
	}

	// Kirim pesan
	_, err := s.client.SendMessage(context.Background(), groupJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	s.logger.Infof("Message sent to group: %s", groupJID.String())
	return nil
}

// selectRandomTemplate memilih template secara random
func (s *GroupManagerService) selectRandomTemplate(templates []database.PromoteTemplate) database.PromoteTemplate {
	if len(templates) == 0 {
		// Return empty template jika tidak ada
		return database.PromoteTemplate{}
	}

	// Pilih index random
	index := rand.Intn(len(templates))
	return templates[index]
}

// Helper functions - menggunakan yang sudah ada di auto_promote.go
