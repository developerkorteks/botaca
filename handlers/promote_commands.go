// Package handlers - Command handlers untuk fitur auto promote
package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/nabilulilalbab/promote/services"
	"github.com/nabilulilalbab/promote/utils"
)

// PromoteCommandHandler menangani command-command auto promote
type PromoteCommandHandler struct {
	autoPromoteService *services.AutoPromoteService
	templateService    *services.TemplateService
	logger             *utils.Logger
}

// NewPromoteCommandHandler membuat handler baru
func NewPromoteCommandHandler(
	autoPromoteService *services.AutoPromoteService,
	templateService *services.TemplateService,
	logger *utils.Logger,
) *PromoteCommandHandler {
	return &PromoteCommandHandler{
		autoPromoteService: autoPromoteService,
		templateService:    templateService,
		logger:             logger,
	}
}

// HandlePromoteCommand menangani command .promote
func (h *PromoteCommandHandler) HandlePromoteCommand(evt *events.Message) string {
	// Hanya bisa digunakan di grup
	if evt.Info.Chat.Server != types.GroupServer {
		return "❌ Command .promote hanya bisa digunakan di grup!"
	}

	groupJID := evt.Info.Chat.String()
	
	// Aktifkan auto promote
	err := h.autoPromoteService.StartAutoPromote(groupJID)
	if err != nil {
		h.logger.Errorf("Failed to start auto promote for %s: %v", groupJID, err)
		return fmt.Sprintf("❌ Gagal mengaktifkan auto promote: %s", err.Error())
	}

	return `✅ *AUTO PROMOTE DIAKTIFKAN!* 🚀

🎯 **Status:** Aktif untuk grup ini
⏰ **Interval:** Setiap 4 jam sekali
📝 **Template:** Random dari template yang tersedia
🔄 **Mulai:** Promosi pertama akan dikirim dalam 4 jam

💡 **Info:**
• Gunakan .disablepromote untuk menghentikan
• Gunakan .statuspromo untuk cek status
• Gunakan .testpromo untuk test kirim promosi

🎉 Selamat! Auto promote sudah aktif untuk grup ini.`
}

// HandleDisablePromoteCommand menangani command .disablepromote
func (h *PromoteCommandHandler) HandleDisablePromoteCommand(evt *events.Message) string {
	// Hanya bisa digunakan di grup
	if evt.Info.Chat.Server != types.GroupServer {
		return "❌ Command .disablepromote hanya bisa digunakan di grup!"
	}

	groupJID := evt.Info.Chat.String()
	
	// Nonaktifkan auto promote
	err := h.autoPromoteService.StopAutoPromote(groupJID)
	if err != nil {
		h.logger.Errorf("Failed to stop auto promote for %s: %v", groupJID, err)
		return fmt.Sprintf("❌ Gagal menonaktifkan auto promote: %s", err.Error())
	}

	return `🛑 *AUTO PROMOTE DINONAKTIFKAN!*

❌ **Status:** Tidak aktif untuk grup ini
⏹️ **Promosi otomatis dihentikan**

💡 **Info:**
• Gunakan .promote untuk mengaktifkan kembali
• Template dan pengaturan tetap tersimpan
• Anda bisa mengaktifkan kapan saja

👋 Auto promote berhasil dinonaktifkan untuk grup ini.`
}

// HandleStatusPromoCommand menangani command .statuspromo
func (h *PromoteCommandHandler) HandleStatusPromoCommand(evt *events.Message) string {
	// Hanya bisa digunakan di grup
	if evt.Info.Chat.Server != types.GroupServer {
		return "❌ Command .statuspromo hanya bisa digunakan di grup!"
	}

	groupJID := evt.Info.Chat.String()
	
	// Ambil status grup
	group, err := h.autoPromoteService.GetGroupStatus(groupJID)
	if err != nil {
		h.logger.Errorf("Failed to get group status for %s: %v", groupJID, err)
		return "❌ Gagal mendapatkan status grup"
	}

	if group == nil {
		return `📊 *STATUS AUTO PROMOTE*

❌ **Status:** Tidak terdaftar
💡 **Info:** Grup ini belum pernah menggunakan auto promote

🚀 Gunakan .promote untuk mengaktifkan auto promote`
	}

	// Format status
	status := "❌ Tidak Aktif"
	if group.IsActive {
		status = "✅ Aktif"
	}

	var startedInfo string
	if group.StartedAt != nil {
		startedInfo = group.StartedAt.Format("2006-01-02 15:04")
	} else {
		startedInfo = "Belum pernah"
	}

	var lastPromoteInfo string
	if group.LastPromoteAt != nil {
		lastPromoteInfo = group.LastPromoteAt.Format("2006-01-02 15:04")
	} else {
		lastPromoteInfo = "Belum pernah"
	}

	// Ambil jumlah template aktif
	templates, _ := h.templateService.GetActiveTemplates()
	templateCount := len(templates)

	return fmt.Sprintf(`📊 *STATUS AUTO PROMOTE*

🎯 **Status:** %s
📅 **Dimulai:** %s
⏰ **Promosi Terakhir:** %s
📝 **Template Tersedia:** %d template

💡 **Commands:**
• .promote - Aktifkan auto promote
• .disablepromote - Nonaktifkan auto promote
• .testpromo - Test kirim promosi
• .listtemplates - Lihat template`, status, startedInfo, lastPromoteInfo, templateCount)
}

// HandleTestPromoCommand menangani command .testpromo
func (h *PromoteCommandHandler) HandleTestPromoCommand(evt *events.Message) string {
	// Hanya bisa digunakan di grup
	if evt.Info.Chat.Server != types.GroupServer {
		return "❌ Command .testpromo hanya bisa digunakan di grup!"
	}

	groupJID := evt.Info.Chat.String()
	
	// Kirim promosi manual
	err := h.autoPromoteService.SendManualPromote(groupJID)
	if err != nil {
		h.logger.Errorf("Failed to send manual promote for %s: %v", groupJID, err)
		return fmt.Sprintf("❌ Gagal mengirim test promosi: %s", err.Error())
	}

	return `🧪 *TEST PROMOSI BERHASIL!*

✅ Promosi test telah dikirim ke grup ini
🎲 Template dipilih secara random
📝 Ini adalah contoh bagaimana auto promote bekerja

💡 **Info:**
• Test ini tidak mempengaruhi jadwal auto promote
• Auto promote tetap berjalan sesuai interval 4 jam
• Gunakan .statuspromo untuk cek status`
}

// HandleListTemplatesCommand menangani command .listtemplates
func (h *PromoteCommandHandler) HandleListTemplatesCommand(evt *events.Message) string {
	templates, err := h.templateService.GetActiveTemplates()
	if err != nil {
		h.logger.Errorf("Failed to get templates: %v", err)
		return "❌ Gagal mendapatkan daftar template"
	}

	if len(templates) == 0 {
		return `📝 *DAFTAR TEMPLATE PROMOSI*

❌ Tidak ada template aktif yang tersedia

💡 **Info:**
• Admin belum menambahkan template promosi
• Gunakan .addtemplate untuk menambah template (admin only)
• Contoh: .addtemplate "Promo Hari Ini" "diskon" "🔥 Diskon 50%!"
• Gunakan .alltemplates untuk melihat semua template`
	}

	var result strings.Builder
	result.WriteString("📝 *DAFTAR TEMPLATE PROMOSI*\n\n")
	result.WriteString(fmt.Sprintf("📊 **Total:** %d template aktif\n\n", len(templates)))

	for i, template := range templates {
		if i >= 15 { // Batasi tampilan maksimal 15 template
			result.WriteString(fmt.Sprintf("... dan %d template lainnya\n", len(templates)-15))
			break
		}

		result.WriteString(fmt.Sprintf("🆔 **ID: %d** - %s\n", template.ID, template.Title))
		result.WriteString(fmt.Sprintf("📂 Kategori: %s\n", template.Category))
		result.WriteString(fmt.Sprintf("📅 Dibuat: %s\n", template.CreatedAt.Format("2006-01-02")))
		result.WriteString(fmt.Sprintf("✅ Status: %s\n\n", getTemplateStatusText(template.IsActive)))
	}

	result.WriteString("💡 **Commands:**\n")
	result.WriteString("• .previewtemplate [ID] - Preview template\n")
	result.WriteString("• .alltemplates - Lihat semua template (aktif & nonaktif)\n")
	result.WriteString("• .addtemplate - Tambah template (admin)\n")
	result.WriteString("• .edittemplate [ID] - Edit template (admin)\n")
	result.WriteString("• .deletetemplate [ID] - Hapus template (admin)\n\n")
	result.WriteString("📋 **Contoh:** .previewtemplate 1 atau .deletetemplate 5")

	return result.String()
}

// getTemplateStatusText helper function untuk status template
func getTemplateStatusText(isActive bool) string {
	if isActive {
		return "Aktif ✅"
	}
	return "Tidak Aktif ❌"
}

// HandlePreviewTemplateCommand menangani command .previewtemplate [ID]
func (h *PromoteCommandHandler) HandlePreviewTemplateCommand(evt *events.Message, args []string) string {
	if len(args) < 2 {
		return `❌ *FORMAT SALAH*

📝 **Format:** .previewtemplate [ID]
📋 **Contoh:** .previewtemplate 1

💡 Gunakan .listtemplates untuk melihat daftar template`
	}

	// Parse ID template
	templateID, err := strconv.Atoi(args[1])
	if err != nil {
		return "❌ ID template harus berupa angka"
	}

	// Preview template
	preview, err := h.templateService.PreviewTemplate(templateID)
	if err != nil {
		h.logger.Errorf("Failed to preview template %d: %v", templateID, err)
		return fmt.Sprintf("❌ Gagal preview template: %s", err.Error())
	}

	return preview
}

// HandleAllTemplatesCommand menangani command .alltemplates
func (h *PromoteCommandHandler) HandleAllTemplatesCommand(evt *events.Message) string {
	templates, err := h.templateService.GetAllTemplates()
	if err != nil {
		h.logger.Errorf("Failed to get all templates: %v", err)
		return "❌ Gagal mendapatkan daftar template"
	}

	if len(templates) == 0 {
		return `📝 *SEMUA TEMPLATE PROMOSI*

❌ Database template masih kosong

💡 **Cara Menambah Template (Admin Only):**
• .addtemplate "Judul" "Kategori" "Konten"
• Contoh: .addtemplate "Flash Sale" "diskon" "🔥 FLASH SALE! Diskon 70% hari ini! Order: 08123456789"

📋 **Kategori yang Disarankan:**
• produk, diskon, testimoni, flashsale, bundle, ongkir, cashback, limited, contact`
	}

	var result strings.Builder
	result.WriteString("📝 *SEMUA TEMPLATE PROMOSI*\n\n")
	result.WriteString(fmt.Sprintf("📊 **Total:** %d template\n\n", len(templates)))

	activeCount := 0
	inactiveCount := 0

	for _, template := range templates {
		if template.IsActive {
			activeCount++
		} else {
			inactiveCount++
		}

		statusIcon := "✅"
		if !template.IsActive {
			statusIcon = "❌"
		}

		result.WriteString(fmt.Sprintf("%s **ID: %d** - %s\n", statusIcon, template.ID, template.Title))
		result.WriteString(fmt.Sprintf("📂 Kategori: %s\n", template.Category))
		result.WriteString(fmt.Sprintf("📅 Dibuat: %s\n", template.CreatedAt.Format("2006-01-02")))
		result.WriteString(fmt.Sprintf("✅ Status: %s\n\n", getTemplateStatusText(template.IsActive)))
	}

	result.WriteString(fmt.Sprintf("📊 **Ringkasan:**\n"))
	result.WriteString(fmt.Sprintf("• Aktif: %d template\n", activeCount))
	result.WriteString(fmt.Sprintf("• Tidak Aktif: %d template\n\n", inactiveCount))

	result.WriteString("💡 **Commands Admin:**\n")
	result.WriteString("• .deletetemplate [ID] - Hapus template\n")
	result.WriteString("• .edittemplate [ID] - Edit template\n")
	result.WriteString("• .previewtemplate [ID] - Preview template")

	return result.String()
}

// HandleHelpCommand menangani command .help
func (h *PromoteCommandHandler) HandleHelpCommand(evt *events.Message) string {
	return `📋 *BANTUAN AUTO PROMOTE SYSTEM*

━━━━━━━━━━━━━━━━━━━━━━

🎯 *COMMANDS USER:*

• *.promote*
  Aktifkan auto promote di grup

• *.disablepromote*
  Nonaktifkan auto promote di grup

• *.statuspromo*
  Cek status auto promote grup

• *.testpromo*
  Test kirim promosi manual

• *.listtemplates*
  Lihat template aktif

• *.alltemplates*
  Lihat semua template (aktif & nonaktif)

• *.previewtemplate [ID]*
  Preview template berdasarkan ID
  Contoh: .previewtemplate 5

• *.help*
  Bantuan ini

━━━━━━━━━━━━━━━━━━━━━━

👑 *COMMANDS ADMIN:*

📝 *Template Management:*
• *.addtemplate "Judul" "Kategori" "Konten"*
  Tambah template baru

• *.edittemplate [ID] "Judul" "Kategori" "Konten"*
  Edit template existing

• *.deletetemplate [ID]*
  Hapus template berdasarkan ID

• *.deletemulti [ID1,ID2,ID3]*
  Hapus multiple template
  Contoh: .deletemulti 1,5,8,12

• *.deleteall*
  Hapus semua template

📊 *API & Statistics:*
• *.fetchproducts*
  Ambil produk dari API (digroup per 15)

• *.productstats*
  Statistik produk dari API

• *.templatestats*
  Statistik template

• *.promotestats*
  Statistik auto promote

• *.activegroups*
  Lihat grup yang aktif auto promote

━━━━━━━━━━━━━━━━━━━━━━

💡 *CARA PENGGUNAAN:*

1️⃣ *Aktivasi Auto Promote:*
   • Ketik .promote di grup
   • Bot akan kirim promosi setiap 4 jam
   • Template dipilih random

2️⃣ *Management Template (Admin):*
   • .fetchproducts untuk ambil dari API
   • .addtemplate untuk tambah manual
   • .listtemplates untuk lihat semua

3️⃣ *Monitoring:*
   • .statuspromo untuk cek status grup
   • .activegroups untuk lihat semua grup aktif
   • .templatestats untuk statistik

━━━━━━━━━━━━━━━━━━━━━━

⚙️ *FITUR SISTEM:*

🤖 *Auto Promote:*
• Interval: 4 jam sekali
• Random template selection
• Per-group control
• Anti-spam protection

📝 *Template System:*
• Support variables: {DATE}, {TIME}, dll
• Kategori: produk, diskon, testimoni, dll
• CRUD operations lengkap
• API integration

🛡️ *Admin Control:*
• Permission-based commands
• Bulk operations
• Real-time statistics
• Error handling

━━━━━━━━━━━━━━━━━━━━━━

❓ *BUTUH BANTUAN?*
Hubungi admin atau gunakan commands di atas

🚀 *Happy Promoting!*`
}

// HandlePromoteHelpCommand menangani command .promotehelp
func (h *PromoteCommandHandler) HandlePromoteHelpCommand(evt *events.Message) string {
	return `📋 *BANTUAN AUTO PROMOTE*

🤖 **Fitur Auto Promote:**
Sistem otomatis untuk mengirim promosi bisnis setiap 4 jam

🎯 **Commands Utama:**
• .promote - Aktifkan auto promote di grup
• .disablepromote - Nonaktifkan auto promote
• .statuspromo - Cek status auto promote
• .testpromo - Test kirim promosi manual

📝 **Commands Template:**
• .listtemplates - Lihat daftar template
• .previewtemplate [ID] - Preview template
• .addtemplate - Tambah template (admin only)
• .edittemplate [ID] - Edit template (admin only)
• .deletetemplate [ID] - Hapus template (admin only)

⚙️ **Commands Admin:**
• .templatestats - Statistik template
• .promotestats - Statistik auto promote
• .activegroups - Lihat grup aktif

💡 **Cara Kerja:**
1. Aktifkan dengan .promote di grup
2. Bot akan kirim promosi setiap 4 jam
3. Template dipilih random dari yang tersedia
4. Nonaktifkan kapan saja dengan .disablepromote

🎲 **Template System:**
• 10+ template promosi bisnis siap pakai
• Random selection untuk variasi
• Admin bisa tambah/edit template
• Support variables: {DATE}, {TIME}, dll

❓ **Butuh bantuan?**
Hubungi admin atau gunakan command di atas`
}

// IsPromoteCommand mengecek apakah pesan adalah command auto promote
func (h *PromoteCommandHandler) IsPromoteCommand(messageText string) bool {
	lowerText := strings.ToLower(strings.TrimSpace(messageText))
	
	promoteCommands := []string{
		".promote",
		".disablepromote", 
		".statuspromo",
		".testpromo",
		".listtemplates",
		".alltemplates",
		".previewtemplate",
		".promotehelp",
		".addtemplate",
		".edittemplate", 
		".deletetemplate",
		".templatestats",
		".promotestats",
		".activegroups",
		".fetchproducts",
		".productstats",
		".deleteall",
		".deletemulti",
		".help",
	}
	
	for _, cmd := range promoteCommands {
		if strings.HasPrefix(lowerText, cmd) {
			return true
		}
	}
	
	return false
}

// HandlePromoteCommands menangani semua command auto promote
func (h *PromoteCommandHandler) HandlePromoteCommands(evt *events.Message, messageText string) string {
	lowerText := strings.ToLower(strings.TrimSpace(messageText))
	args := strings.Fields(lowerText)
	
	if len(args) == 0 {
		return ""
	}
	
	command := args[0]
	
	switch command {
	case ".promote":
		return h.HandlePromoteCommand(evt)
		
	case ".disablepromote":
		return h.HandleDisablePromoteCommand(evt)
		
	case ".statuspromo":
		return h.HandleStatusPromoCommand(evt)
		
	case ".testpromo":
		return h.HandleTestPromoCommand(evt)
		
	case ".listtemplates":
		return h.HandleListTemplatesCommand(evt)
		
	case ".alltemplates":
		return h.HandleAllTemplatesCommand(evt)
		
	case ".previewtemplate":
		return h.HandlePreviewTemplateCommand(evt, args)
		
	case ".promotehelp":
		return h.HandlePromoteHelpCommand(evt)
		
	case ".help":
		return h.HandleHelpCommand(evt)
		
	// Admin commands (akan diimplementasi di file terpisah)
	case ".addtemplate", ".edittemplate", ".deletetemplate":
		return "🔧 Command admin template sedang dalam pengembangan"
		
	case ".templatestats", ".promotestats", ".activegroups":
		return "📊 Command statistik sedang dalam pengembangan"
		
	default:
		return ""
	}
}