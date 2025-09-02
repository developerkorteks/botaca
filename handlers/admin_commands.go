// Package handlers - Admin command handlers untuk mengelola template dan sistem
package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow/types/events"

	"github.com/nabilulilalbab/promote/services"
	"github.com/nabilulilalbab/promote/utils"
)

// AdminCommandHandler menangani command admin untuk auto promote
type AdminCommandHandler struct {
	autoPromoteService *services.AutoPromoteService
	templateService    *services.TemplateService
	logger             *utils.Logger
	adminNumbers       []string // Daftar nomor admin yang bisa menggunakan command admin
}

// NewAdminCommandHandler membuat handler baru
func NewAdminCommandHandler(
	autoPromoteService *services.AutoPromoteService,
	templateService *services.TemplateService,
	logger *utils.Logger,
	adminNumbers []string,
) *AdminCommandHandler {
	return &AdminCommandHandler{
		autoPromoteService: autoPromoteService,
		templateService:    templateService,
		logger:             logger,
		adminNumbers:       adminNumbers,
	}
}

// isAdmin mengecek apakah user adalah admin
func (h *AdminCommandHandler) isAdmin(userNumber string) bool {
	for _, admin := range h.adminNumbers {
		if admin == userNumber {
			return true
		}
	}
	return false
}

// HandleAddTemplateCommand menangani command .addtemplate
func (h *AdminCommandHandler) HandleAddTemplateCommand(evt *events.Message, args []string) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	// Format: .addtemplate "Judul" "Kategori" "Konten"
	if len(args) < 4 {
		return `❌ *FORMAT SALAH*

📝 **Format:** .addtemplate "Judul" "Kategori" "Konten"

📋 **Contoh:**
.addtemplate "Flash Sale Hari Ini" "flashsale" "🔥 FLASH SALE! Diskon 50% hanya hari ini! Order: 08123456789"

💡 **Tips:**
• Gunakan tanda kutip untuk teks yang mengandung spasi
• Kategori: produk, diskon, testimoni, flashsale, dll
• Konten bisa menggunakan emoji dan formatting WhatsApp`
	}

	// Parse arguments (simplified parsing)
	fullText := strings.Join(args[1:], " ")
	parts := h.parseQuotedArgs(fullText)
	
	if len(parts) < 3 {
		return "❌ Format salah. Gunakan: .addtemplate \"Judul\" \"Kategori\" \"Konten\""
	}

	title := parts[0]
	category := parts[1]
	content := parts[2]

	// Buat template
	template, err := h.templateService.CreateTemplate(title, content, category)
	if err != nil {
		h.logger.Errorf("Failed to create template: %v", err)
		return fmt.Sprintf("❌ Gagal membuat template: %s", err.Error())
	}

	return fmt.Sprintf(`✅ *TEMPLATE BERHASIL DIBUAT!*

🆔 **ID:** %d
🏷️ **Judul:** %s
📂 **Kategori:** %s
✅ **Status:** Aktif

📝 **Konten:**
%s

💡 **Info:**
• Template langsung aktif dan bisa digunakan
• Gunakan .previewtemplate %d untuk preview
• Gunakan .edittemplate %d untuk edit`, 
		template.ID, template.Title, template.Category, template.Content, template.ID, template.ID)
}

// HandleEditTemplateCommand menangani command .edittemplate
func (h *AdminCommandHandler) HandleEditTemplateCommand(evt *events.Message, args []string) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	// Format: .edittemplate [ID] "Judul" "Kategori" "Konten"
	if len(args) < 5 {
		return `❌ *FORMAT SALAH*

📝 **Format:** .edittemplate [ID] "Judul" "Kategori" "Konten"

📋 **Contoh:**
.edittemplate 1 "Promo Terbaru" "diskon" "🎉 Promo spesial! Diskon 30%"

💡 **Tips:**
• Gunakan .listtemplates untuk melihat ID template
• Gunakan tanda kutip untuk teks yang mengandung spasi`
	}

	// Parse ID
	templateID, err := strconv.Atoi(args[1])
	if err != nil {
		return "❌ ID template harus berupa angka"
	}

	// Parse arguments
	fullText := strings.Join(args[2:], " ")
	parts := h.parseQuotedArgs(fullText)
	
	if len(parts) < 3 {
		return "❌ Format salah. Gunakan: .edittemplate [ID] \"Judul\" \"Kategori\" \"Konten\""
	}

	title := parts[0]
	category := parts[1]
	content := parts[2]

	// Update template
	err = h.templateService.UpdateTemplate(templateID, title, content, category, true)
	if err != nil {
		h.logger.Errorf("Failed to update template %d: %v", templateID, err)
		return fmt.Sprintf("❌ Gagal mengupdate template: %s", err.Error())
	}

	return fmt.Sprintf(`✅ *TEMPLATE BERHASIL DIUPDATE!*

🆔 **ID:** %d
🏷️ **Judul:** %s
📂 **Kategori:** %s

📝 **Konten Baru:**
%s

💡 **Info:**
• Template telah diperbarui
• Gunakan .previewtemplate %d untuk preview
• Perubahan langsung berlaku untuk auto promote`, 
		templateID, title, category, content, templateID)
}

// HandleDeleteTemplateCommand menangani command .deletetemplate
func (h *AdminCommandHandler) HandleDeleteTemplateCommand(evt *events.Message, args []string) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	if len(args) < 2 {
		return `❌ *FORMAT SALAH*

📝 **Format:** .deletetemplate [ID]
📋 **Contoh:** .deletetemplate 5

💡 Gunakan .listtemplates untuk melihat ID template`
	}

	// Parse ID
	templateID, err := strconv.Atoi(args[1])
	if err != nil {
		return "❌ ID template harus berupa angka"
	}

	// Ambil info template sebelum dihapus
	template, err := h.templateService.GetTemplateByID(templateID)
	if err != nil {
		return fmt.Sprintf("❌ Gagal mendapatkan template: %s", err.Error())
	}

	if template == nil {
		return fmt.Sprintf("❌ Template dengan ID %d tidak ditemukan", templateID)
	}

	// Hapus template
	err = h.templateService.DeleteTemplate(templateID)
	if err != nil {
		h.logger.Errorf("Failed to delete template %d: %v", templateID, err)
		return fmt.Sprintf("❌ Gagal menghapus template: %s", err.Error())
	}

	return fmt.Sprintf(`🗑️ *TEMPLATE BERHASIL DIHAPUS!*

🆔 **ID:** %d
🏷️ **Judul:** %s
📂 **Kategori:** %s

⚠️ **Peringatan:**
• Template telah dihapus permanen
• Tidak bisa dikembalikan lagi
• Auto promote akan menggunakan template lain yang tersedia

💡 Gunakan .listtemplates untuk melihat template yang tersisa`, 
		templateID, template.Title, template.Category)
}

// HandleTemplateStatsCommand menangani command .templatestats
func (h *AdminCommandHandler) HandleTemplateStatsCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	stats, err := h.templateService.GetTemplateStats()
	if err != nil {
		h.logger.Errorf("Failed to get template stats: %v", err)
		return "❌ Gagal mendapatkan statistik template"
	}

	var result strings.Builder
	result.WriteString("📊 *STATISTIK TEMPLATE*\n\n")
	result.WriteString(fmt.Sprintf("📝 **Total Template:** %d\n", stats["total"]))
	result.WriteString(fmt.Sprintf("✅ **Aktif:** %d\n", stats["active"]))
	result.WriteString(fmt.Sprintf("❌ **Tidak Aktif:** %d\n\n", stats["inactive"]))

	result.WriteString("📂 **Per Kategori:**\n")
	categories := stats["categories"].(map[string]int)
	for category, count := range categories {
		result.WriteString(fmt.Sprintf("• %s: %d template\n", category, count))
	}

	result.WriteString("\n💡 **Commands:**\n")
	result.WriteString("• .addtemplate - Tambah template baru\n")
	result.WriteString("• .edittemplate [ID] - Edit template\n")
	result.WriteString("• .deletetemplate [ID] - Hapus template")

	return result.String()
}

// HandlePromoteStatsCommand menangani command .promotestats
func (h *AdminCommandHandler) HandlePromoteStatsCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	// Ambil jumlah grup aktif
	activeCount, err := h.autoPromoteService.GetActiveGroupsCount()
	if err != nil {
		h.logger.Errorf("Failed to get active groups count: %v", err)
		return "❌ Gagal mendapatkan statistik auto promote"
	}

	return fmt.Sprintf(`📊 *STATISTIK AUTO PROMOTE*

🎯 **Grup Aktif:** %d grup
⏰ **Interval:** Setiap 4 jam
🤖 **Status Scheduler:** Berjalan

📈 **Performa:**
• Total grup terdaftar: %d
• Grup aktif: %d
• Grup tidak aktif: %d

💡 **Info:**
• Statistik diperbarui real-time
• Gunakan .activegroups untuk detail grup
• Scheduler berjalan otomatis`, activeCount, activeCount, activeCount, 0)
}

// HandleActiveGroupsCommand menangani command .activegroups
func (h *AdminCommandHandler) HandleActiveGroupsCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	return `👥 *GRUP AKTIF AUTO PROMOTE*

🔧 Fitur ini sedang dalam pengembangan

💡 **Yang akan ditampilkan:**
• Daftar grup yang menggunakan auto promote
• Status terakhir promosi
• Waktu aktivasi
• Statistik per grup

📞 Hubungi developer untuk informasi lebih lanjut`
}

// parseQuotedArgs memparse argument yang menggunakan tanda kutip
func (h *AdminCommandHandler) parseQuotedArgs(text string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	
	for i, char := range text {
		if char == '"' {
			if inQuotes {
				// End of quoted string
				args = append(args, current.String())
				current.Reset()
				inQuotes = false
			} else {
				// Start of quoted string
				inQuotes = true
			}
		} else if char == ' ' && !inQuotes {
			// Space outside quotes - separator
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
		
		// Handle end of string
		if i == len(text)-1 && current.Len() > 0 {
			args = append(args, current.String())
		}
	}
	
	return args
}

// HandleAdminCommands menangani semua command admin
func (h *AdminCommandHandler) HandleAdminCommands(evt *events.Message, messageText string) string {
	args := strings.Fields(messageText) // Gunakan original text untuk preserve case
	
	if len(args) == 0 {
		return ""
	}
	
	command := strings.ToLower(args[0])
	
	switch command {
	case ".addtemplate":
		return h.HandleAddTemplateCommand(evt, args)
		
	case ".edittemplate":
		return h.HandleEditTemplateCommand(evt, args)
		
	case ".deletetemplate":
		return h.HandleDeleteTemplateCommand(evt, args)
		
	case ".templatestats":
		return h.HandleTemplateStatsCommand(evt)
		
	case ".promotestats":
		return h.HandlePromoteStatsCommand(evt)
		
	case ".activegroups":
		return h.HandleActiveGroupsCommand(evt)
		
	default:
		return ""
	}
}