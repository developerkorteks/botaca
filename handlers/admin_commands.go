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
	apiProductService  *services.APIProductService
	logger             *utils.Logger
	adminNumbers       []string // Daftar nomor admin yang bisa menggunakan command admin
}

// NewAdminCommandHandler membuat handler baru
func NewAdminCommandHandler(
	autoPromoteService *services.AutoPromoteService,
	templateService *services.TemplateService,
	apiProductService *services.APIProductService,
	logger *utils.Logger,
	adminNumbers []string,
) *AdminCommandHandler {
	return &AdminCommandHandler{
		autoPromoteService: autoPromoteService,
		templateService:    templateService,
		apiProductService:  apiProductService,
		logger:             logger,
		adminNumbers:       adminNumbers,
	}
}

// isAdmin mengecek apakah user adalah admin dengan validasi ketat
func (h *AdminCommandHandler) isAdmin(userNumber string) bool {
	// Validasi input
	if userNumber == "" {
		h.logger.Warning("Empty user number provided for admin check")
		return false
	}
	
	// Log attempt untuk security monitoring
	h.logger.Debugf("Admin check for user: %s", userNumber)
	
	// Cek apakah user ada dalam daftar admin
	for _, admin := range h.adminNumbers {
		if admin == userNumber {
			h.logger.Infof("Admin access granted for: %s", userNumber)
			return true
		}
	}
	
	// Log unauthorized attempt
	h.logger.Warningf("Unauthorized admin attempt from: %s", userNumber)
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

	// Ambil daftar grup aktif dari service
	activeGroups, err := h.autoPromoteService.GetActiveGroups()
	if err != nil {
		h.logger.Errorf("Failed to get active groups: %v", err)
		return "❌ Gagal mendapatkan daftar grup aktif"
	}

	if len(activeGroups) == 0 {
		return `👥 *GRUP AKTIF AUTO PROMOTE*

❌ Tidak ada grup yang menggunakan auto promote

💡 **Info:**
• Gunakan .promote di grup untuk mengaktifkan
• Auto promote akan muncul di sini setelah diaktifkan
• Gunakan .promotestats untuk statistik umum`
	}

	var result strings.Builder
	result.WriteString("👥 *GRUP AKTIF AUTO PROMOTE*\n\n")
	result.WriteString(fmt.Sprintf("📊 **Total Grup Aktif:** %d\n\n", len(activeGroups)))

	for i, group := range activeGroups {
		if i >= 20 { // Batasi tampilan maksimal 20 grup
			result.WriteString(fmt.Sprintf("... dan %d grup lainnya\n", len(activeGroups)-20))
			break
		}

		// Format group JID untuk tampilan
		groupDisplay := h.formatGroupJID(group.GroupJID)
		
		result.WriteString(fmt.Sprintf("**%d.** 👥 %s\n", i+1, groupDisplay))
		
		if group.StartedAt != nil {
			result.WriteString(fmt.Sprintf("📅 Dimulai: %s\n", group.StartedAt.Format("2006-01-02 15:04")))
		}
		
		if group.LastPromoteAt != nil {
			result.WriteString(fmt.Sprintf("⏰ Promosi Terakhir: %s\n", group.LastPromoteAt.Format("2006-01-02 15:04")))
		} else {
			result.WriteString("⏰ Promosi Terakhir: Belum ada\n")
		}
		
		result.WriteString(fmt.Sprintf("✅ Status: Aktif\n\n"))
	}

	result.WriteString("💡 **Commands:**\n")
	result.WriteString("• .promotestats - Statistik detail\n")
	result.WriteString("• .testpromo - Test promosi manual")

	return result.String()
}

// formatGroupJID memformat group JID untuk tampilan yang lebih readable
func (h *AdminCommandHandler) formatGroupJID(groupJID string) string {
	// Ambil hanya bagian ID grup (sebelum @g.us)
	if strings.Contains(groupJID, "@g.us") {
		parts := strings.Split(groupJID, "@")
		if len(parts) > 0 {
			return fmt.Sprintf("Grup-%s", parts[0][len(parts[0])-8:]) // 8 digit terakhir
		}
	}
	return groupJID
}

// HandleFetchProductsCommand menangani command .fetchproducts
func (h *AdminCommandHandler) HandleFetchProductsCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	if h.apiProductService == nil {
		return "❌ API Product service tidak tersedia"
	}

	h.logger.Info("Admin requesting product fetch from API...")
	
	result, err := h.apiProductService.FetchProductsAndCreateTemplates()
	if err != nil {
		h.logger.Errorf("Failed to fetch products: %v", err)
		return fmt.Sprintf("❌ Gagal mengambil produk dari API: %s", err.Error())
	}

	return result
}

// HandleProductStatsCommand menangani command .productstats
func (h *AdminCommandHandler) HandleProductStatsCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	if h.apiProductService == nil {
		return "❌ API Product service tidak tersedia"
	}

	result, err := h.apiProductService.GetProductStats()
	if err != nil {
		h.logger.Errorf("Failed to get product stats: %v", err)
		return fmt.Sprintf("❌ Gagal mendapatkan statistik produk: %s", err.Error())
	}

	return result
}

// HandleDeleteAllTemplatesCommand menangani command .deleteall
func (h *AdminCommandHandler) HandleDeleteAllTemplatesCommand(evt *events.Message) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	// Ambil semua template
	templates, err := h.templateService.GetAllTemplates()
	if err != nil {
		return fmt.Sprintf("❌ Gagal mendapatkan template: %s", err.Error())
	}

	if len(templates) == 0 {
		return "❌ Tidak ada template untuk dihapus"
	}

	// Hapus semua template
	deletedCount := 0
	var errors []string

	for _, template := range templates {
		err := h.templateService.DeleteTemplate(template.ID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("ID %d: %v", template.ID, err))
		} else {
			deletedCount++
		}
	}

	var result strings.Builder
	result.WriteString("🗑️ *HAPUS SEMUA TEMPLATE*\n\n")
	result.WriteString(fmt.Sprintf("✅ **Berhasil dihapus:** %d template\n", deletedCount))
	
	if len(errors) > 0 {
		result.WriteString(fmt.Sprintf("❌ **Gagal dihapus:** %d template\n", len(errors)))
	}

	result.WriteString("\n⚠️ **PERINGATAN:**\n")
	result.WriteString("• Semua template telah dihapus permanen\n")
	result.WriteString("• Auto promote akan berhenti jika tidak ada template\n")
	result.WriteString("• Gunakan .fetchproducts untuk isi ulang template")

	return result.String()
}

// HandleDeleteMultipleTemplatesCommand menangani command .deletemulti [ID1,ID2,ID3]
func (h *AdminCommandHandler) HandleDeleteMultipleTemplatesCommand(evt *events.Message, args []string) string {
	// Cek admin permission
	if !h.isAdmin(evt.Info.Sender.User) {
		return "❌ Command ini hanya bisa digunakan oleh admin"
	}

	if len(args) < 2 {
		return `❌ *FORMAT SALAH*

📝 **Format:** .deletemulti [ID1,ID2,ID3]
📋 **Contoh:** .deletemulti 1,5,8,12

💡 **Tips:**
• Pisahkan ID dengan koma tanpa spasi
• Gunakan .alltemplates untuk melihat ID
• Maksimal 20 ID sekaligus`
	}

	// Parse ID dari argument
	idsStr := strings.Join(args[1:], "")
	idStrings := strings.Split(idsStr, ",")
	
	if len(idStrings) > 20 {
		return "❌ Maksimal 20 template sekaligus"
	}

	var ids []int
	for _, idStr := range idStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Sprintf("❌ ID tidak valid: %s", idStr)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return "❌ Tidak ada ID yang valid"
	}

	// Hapus template berdasarkan ID
	deletedCount := 0
	var errors []string
	var deletedTitles []string

	for _, id := range ids {
		// Ambil info template sebelum dihapus
		template, err := h.templateService.GetTemplateByID(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("ID %d: tidak ditemukan", id))
			continue
		}

		err = h.templateService.DeleteTemplate(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("ID %d: %v", id, err))
		} else {
			deletedCount++
			deletedTitles = append(deletedTitles, fmt.Sprintf("ID %d: %s", id, template.Title))
		}
	}

	var result strings.Builder
	result.WriteString("🗑️ *HAPUS MULTIPLE TEMPLATE*\n\n")
	result.WriteString(fmt.Sprintf("✅ **Berhasil dihapus:** %d template\n", deletedCount))
	
	if len(errors) > 0 {
		result.WriteString(fmt.Sprintf("❌ **Gagal dihapus:** %d template\n\n", len(errors)))
	}

	if len(deletedTitles) > 0 {
		result.WriteString("📋 **Template yang dihapus:**\n")
		for _, title := range deletedTitles {
			result.WriteString(fmt.Sprintf("• %s\n", title))
		}
	}

	return result.String()
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
		
	case ".fetchproducts":
		return h.HandleFetchProductsCommand(evt)
		
	case ".productstats":
		return h.HandleProductStatsCommand(evt)
		
	case ".deleteall":
		return h.HandleDeleteAllTemplatesCommand(evt)
		
	case ".deletemulti":
		return h.HandleDeleteMultipleTemplatesCommand(evt, args)
		
	default:
		return ""
	}
}