// Package services - Template management service untuk mengelola template promosi
package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/nabilulilalbab/promote/database"
	"github.com/nabilulilalbab/promote/utils"
)

// TemplateService mengelola template promosi
type TemplateService struct {
	repository database.Repository
	logger     *utils.Logger
}

// NewTemplateService membuat service baru
func NewTemplateService(repo database.Repository, logger *utils.Logger) *TemplateService {
	return &TemplateService{
		repository: repo,
		logger:     logger,
	}
}

// GetAllTemplates mendapatkan semua template
func (s *TemplateService) GetAllTemplates() ([]database.PromoteTemplate, error) {
	templates, err := s.repository.GetAllTemplates()
	if err != nil {
		s.logger.Errorf("Failed to get all templates: %v", err)
		return nil, err
	}

	s.logger.Infof("Retrieved %d templates", len(templates))
	return templates, nil
}

// GetActiveTemplates mendapatkan template yang aktif
func (s *TemplateService) GetActiveTemplates() ([]database.PromoteTemplate, error) {
	templates, err := s.repository.GetActiveTemplates()
	if err != nil {
		s.logger.Errorf("Failed to get active templates: %v", err)
		return nil, err
	}

	s.logger.Infof("Retrieved %d active templates", len(templates))
	return templates, nil
}

// GetTemplateByID mendapatkan template berdasarkan ID
func (s *TemplateService) GetTemplateByID(id int) (*database.PromoteTemplate, error) {
	template, err := s.repository.GetTemplateByID(id)
	if err != nil {
		s.logger.Errorf("Failed to get template %d: %v", id, err)
		return nil, err
	}

	if template == nil {
		return nil, fmt.Errorf("template dengan ID %d tidak ditemukan", id)
	}

	return template, nil
}

// CreateTemplate membuat template baru
func (s *TemplateService) CreateTemplate(title, content, category string) (*database.PromoteTemplate, error) {
	// Validasi input
	if err := s.validateTemplate(title, content, category); err != nil {
		return nil, err
	}

	template := &database.PromoteTemplate{
		Title:    strings.TrimSpace(title),
		Content:  strings.TrimSpace(content),
		Category: strings.ToLower(strings.TrimSpace(category)),
		IsActive: true,
	}

	err := s.repository.CreateTemplate(template)
	if err != nil {
		s.logger.Errorf("Failed to create template: %v", err)
		return nil, fmt.Errorf("gagal membuat template: %v", err)
	}

	s.logger.Successf("Template created: %s (ID: %d)", template.Title, template.ID)
	return template, nil
}

// UpdateTemplate mengupdate template yang ada
func (s *TemplateService) UpdateTemplate(id int, title, content, category string, isActive bool) error {
	// Cek apakah template ada
	existing, err := s.repository.GetTemplateByID(id)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("template dengan ID %d tidak ditemukan", id)
	}

	// Validasi input
	if err := s.validateTemplate(title, content, category); err != nil {
		return err
	}

	// Update template
	existing.Title = strings.TrimSpace(title)
	existing.Content = strings.TrimSpace(content)
	existing.Category = strings.ToLower(strings.TrimSpace(category))
	existing.IsActive = isActive

	err = s.repository.UpdateTemplate(existing)
	if err != nil {
		s.logger.Errorf("Failed to update template %d: %v", id, err)
		return fmt.Errorf("gagal mengupdate template: %v", err)
	}

	s.logger.Successf("Template updated: %s (ID: %d)", existing.Title, existing.ID)
	return nil
}

// DeleteTemplate menghapus template
func (s *TemplateService) DeleteTemplate(id int) error {
	// Cek apakah template ada
	existing, err := s.repository.GetTemplateByID(id)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("template dengan ID %d tidak ditemukan", id)
	}

	err = s.repository.DeleteTemplate(id)
	if err != nil {
		s.logger.Errorf("Failed to delete template %d: %v", id, err)
		return fmt.Errorf("gagal menghapus template: %v", err)
	}

	s.logger.Successf("Template deleted: %s (ID: %d)", existing.Title, existing.ID)
	return nil
}

// ToggleTemplateStatus mengaktifkan/menonaktifkan template
func (s *TemplateService) ToggleTemplateStatus(id int) error {
	template, err := s.repository.GetTemplateByID(id)
	if err != nil {
		return err
	}

	if template == nil {
		return fmt.Errorf("template dengan ID %d tidak ditemukan", id)
	}

	// Toggle status
	template.IsActive = !template.IsActive

	err = s.repository.UpdateTemplate(template)
	if err != nil {
		s.logger.Errorf("Failed to toggle template %d status: %v", id, err)
		return fmt.Errorf("gagal mengubah status template: %v", err)
	}

	status := "dinonaktifkan"
	if template.IsActive {
		status = "diaktifkan"
	}

	s.logger.Successf("Template %s: %s (ID: %d)", status, template.Title, template.ID)
	return nil
}

// GetTemplatesByCategory mendapatkan template berdasarkan kategori
func (s *TemplateService) GetTemplatesByCategory(category string) ([]database.PromoteTemplate, error) {
	allTemplates, err := s.repository.GetAllTemplates()
	if err != nil {
		return nil, err
	}

	var filtered []database.PromoteTemplate
	categoryLower := strings.ToLower(category)

	for _, template := range allTemplates {
		if template.Category == categoryLower {
			filtered = append(filtered, template)
		}
	}

	s.logger.Infof("Found %d templates in category: %s", len(filtered), category)
	return filtered, nil
}

// GetTemplateCategories mendapatkan daftar kategori yang tersedia
func (s *TemplateService) GetTemplateCategories() ([]string, error) {
	templates, err := s.repository.GetAllTemplates()
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]bool)
	for _, template := range templates {
		categoryMap[template.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}

// PreviewTemplate memformat template untuk preview
func (s *TemplateService) PreviewTemplate(templateID int) (string, error) {
	template, err := s.repository.GetTemplateByID(templateID)
	if err != nil {
		return "", err
	}

	if template == nil {
		return "", fmt.Errorf("template tidak ditemukan")
	}

	// Proses template dengan sample data
	preview := s.processTemplateForPreview(template.Content)

	return fmt.Sprintf(`📋 *PREVIEW TEMPLATE*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	          *DETAIL TEMPLATE*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🏷️ *Judul:* %s
📂 *Kategori:* %s
📈 *Status:* %s

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	          *KONTEN PREVIEW*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%s

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	          *INFORMASI*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 Variabel dinamis seperti *{DATE}* dan *{TIME}* akan diganti saat promosi dikirim.`,
		template.Title,
		template.Category,
		getStatusText(template.IsActive),
		preview), nil
}

// validateTemplate memvalidasi input template
func (s *TemplateService) validateTemplate(title, content, category string) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	category = strings.TrimSpace(category)

	if title == "" {
		return fmt.Errorf("judul template tidak boleh kosong")
	}

	if len(title) > 100 {
		return fmt.Errorf("judul template maksimal 100 karakter")
	}

	if content == "" {
		return fmt.Errorf("konten template tidak boleh kosong")
	}

	if len(content) > 4000 {
		return fmt.Errorf("konten template maksimal 4000 karakter")
	}

	if category == "" {
		return fmt.Errorf("kategori template tidak boleh kosong")
	}

	if len(category) > 50 {
		return fmt.Errorf("kategori template maksimal 50 karakter")
	}

	return nil
}

// processTemplateForPreview memproses template untuk preview
func (s *TemplateService) processTemplateForPreview(content string) string {
	now := time.Now()

	replacements := map[string]string{
		"{DATE}":  now.Format("2006-01-02"),
		"{TIME}":  now.Format("15:04"),
		"{DAY}":   getDayName(now.Weekday()),
		"{MONTH}": getMonthName(now.Month()),
		"{YEAR}":  fmt.Sprintf("%d", now.Year()),
	}

	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// getStatusText mengkonversi boolean status ke teks
func getStatusText(isActive bool) string {
	if isActive {
		return "Aktif ✅"
	}
	return "Tidak Aktif ❌"
}

// GetTemplateStats mendapatkan statistik template
func (s *TemplateService) GetTemplateStats() (map[string]interface{}, error) {
	templates, err := s.repository.GetAllTemplates()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":      len(templates),
		"active":     0,
		"inactive":   0,
		"categories": make(map[string]int),
	}

	categoryCount := make(map[string]int)

	for _, template := range templates {
		if template.IsActive {
			stats["active"] = stats["active"].(int) + 1
		} else {
			stats["inactive"] = stats["inactive"].(int) + 1
		}

		categoryCount[template.Category]++
	}

	stats["categories"] = categoryCount

	return stats, nil
}
