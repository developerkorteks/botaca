// Package handlers berisi semua handler untuk menangani pesan dan event WhatsApp
// File ini khusus menangani pesan masuk dari chat personal dan grup
package handlers

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// MessageHandler adalah struktur yang menangani semua pesan masuk
type MessageHandler struct {
	// client adalah instance WhatsApp client untuk mengirim pesan
	client *whatsmeow.Client

	// autoReplyPersonal menentukan apakah bot membalas chat personal
	autoReplyPersonal bool

	// autoReplyGroup menentukan apakah bot membalas chat grup
	autoReplyGroup bool

	// Auto Promote handlers
	promoteCommandHandler *PromoteCommandHandler
	adminCommandHandler   *AdminCommandHandler
}

// NewMessageHandler membuat handler baru untuk pesan
// Parameter:
// - client: WhatsApp client yang sudah terhubung
// - autoReplyPersonal: true jika ingin auto reply di chat personal
// - autoReplyGroup: true jika ingin auto reply di grup (hati-hati spam!)
func NewMessageHandler(client *whatsmeow.Client, autoReplyPersonal, autoReplyGroup bool) *MessageHandler {
	return &MessageHandler{
		client:            client,
		autoReplyPersonal: autoReplyPersonal,
		autoReplyGroup:    autoReplyGroup,
	}
}

// SetAutoPromoteHandlers mengatur handlers untuk auto promote
func (h *MessageHandler) SetAutoPromoteHandlers(promoteHandler *PromoteCommandHandler, adminHandler *AdminCommandHandler) {
	h.promoteCommandHandler = promoteHandler
	h.adminCommandHandler = adminHandler
}

// HandleMessage adalah fungsi utama untuk menangani pesan masuk
// Fungsi ini akan dipanggil setiap kali ada pesan baru
func (h *MessageHandler) HandleMessage(evt *events.Message) {
	// STEP 1: Skip pesan dari diri sendiri
	// Ini penting untuk menghindari bot membalas pesannya sendiri (infinite loop)
	if evt.Info.IsFromMe {
		return
	}

	// STEP 2: Ambil teks dari pesan
	// WhatsApp memiliki beberapa tipe pesan, kita hanya proses yang teks
	messageText := h.getMessageText(evt.Message)
	if messageText == "" {
		// Jika bukan pesan teks (misal gambar, voice note), skip
		return
	}

	// STEP 3: Identifikasi jenis chat (personal atau grup)
	isGroup := evt.Info.Chat.Server == types.GroupServer
	chatType := "personal"
	if isGroup {
		chatType = "group"
	}

	// STEP 4: Log informasi pesan untuk debugging
	sender := evt.Info.Sender.User // Nomor pengirim (tanpa @s.whatsapp.net)
	fmt.Printf("📨 Pesan masuk [%s]: %s\n", chatType, messageText)
	fmt.Printf("👤 Dari: %s\n", sender)

	// Jika grup, tampilkan nama grup juga
	if isGroup {
		fmt.Printf("👥 Grup: %s\n", evt.Info.Chat.User)
	}

	// STEP 5: Proses pesan berdasarkan jenis chat
	if isGroup {
		h.handleGroupMessage(evt, messageText)
	} else {
		h.handlePersonalMessage(evt, messageText)
	}
}

// handlePersonalMessage menangani pesan dari chat personal (1 on 1)
func (h *MessageHandler) handlePersonalMessage(evt *events.Message, messageText string) {
	fmt.Println("💬 Memproses pesan personal...")

	// Cek apakah ini adalah command (dimulai dengan / atau .)
	if strings.HasPrefix(messageText, "/") || strings.HasPrefix(messageText, ".") {
		h.handleCommand(evt, messageText)
		return
	}

	// Bot tidak memberikan auto reply untuk non-admin
	// Hanya merespon command auto promote dari admin
}

// handleGroupMessage menangani pesan dari grup
func (h *MessageHandler) handleGroupMessage(evt *events.Message, messageText string) {
	fmt.Println("👥 Memproses pesan grup...")

	// BOT DIAM TOTAL DI GRUP - TIDAK ADA RESPONSE APAPUN
	// Bot hanya akan mengirim auto promote sesuai scheduler
	// Semua kontrol dilakukan melalui chat personal dengan admin

	// Log untuk monitoring (tanpa response)
	fmt.Printf("👥 Grup: %s | Pesan: %s | Action: IGNORED\n",
		evt.Info.Chat.User, h.truncateString(messageText, 30))

	// Bot tidak memberikan response apapun di grup
	return
}

// handleCommand menangani command yang dimulai dengan /
func (h *MessageHandler) handleCommand(evt *events.Message, messageText string) {
	// Ubah ke lowercase untuk case-insensitive commands
	lowerText := strings.ToLower(strings.TrimSpace(messageText))

	var response string

	// Cek apakah ini auto promote command terlebih dahulu
	if h.isAutoPromoteCommand(lowerText) {
		response = h.handleAutoPromoteCommand(evt, messageText)
	} else {
		// Tidak ada response untuk command yang tidak dikenal
		return
	}

	// Kirim response jika ada
	if response != "" {
		h.sendMessage(evt.Info.Chat, response)
	}
}

// sendAutoReply mengirim balasan otomatis
func (h *MessageHandler) sendAutoReply(chatJID types.JID, originalMessage string, isGroup bool) {
	var response string

	if isGroup {
		// Response untuk grup lebih formal dan tidak terlalu sering
		response = `👋 *AUTO-REPLY*

Terima kasih! Saya adalah bot otomatis.
Ketik */help* untuk bantuan.`
	} else {
		// Response untuk personal bisa lebih personal
		response = `👋 *AUTO-REPLY*

✅ Terima kasih atas pesannya!
Saya adalah bot otomatis yang siap membantu.

Ketik */help* untuk melihat command yang tersedia.`
	}

	h.sendMessage(chatJID, response)
}

// getMessageText mengekstrak teks dari berbagai tipe pesan WhatsApp
func (h *MessageHandler) getMessageText(msg *waProto.Message) string {
	// Pesan teks biasa
	if msg.GetConversation() != "" {
		return msg.GetConversation()
	}

	// Pesan teks dengan format (bold, italic, dll) atau reply
	if msg.GetExtendedTextMessage() != nil {
		return msg.GetExtendedTextMessage().GetText()
	}

	// Jika bukan teks, return empty string
	return ""
}

// isBotMentioned mengecek apakah bot di-mention dalam pesan grup
func (h *MessageHandler) isBotMentioned(msg *waProto.Message) bool {
	// Cek di extended text message (yang biasanya berisi mention)
	if msg.GetExtendedTextMessage() != nil && msg.GetExtendedTextMessage().GetContextInfo() != nil {
		mentions := msg.GetExtendedTextMessage().GetContextInfo().GetMentionedJid()
		botJID := h.client.Store.ID.String()

		// Cek apakah JID bot ada dalam daftar mention
		for _, mention := range mentions {
			if mention == botJID {
				return true
			}
		}
	}

	return false
}

// sendMessage mengirim pesan ke chat tertentu
func (h *MessageHandler) sendMessage(chatJID types.JID, text string) {
	// Buat struktur pesan WhatsApp
	msg := &waProto.Message{
		Conversation: &text,
	}

	// Kirim pesan menggunakan client
	_, err := h.client.SendMessage(context.Background(), chatJID, msg)
	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan: %v\n", err)
		return
	}

	// Log pesan yang terkirim
	fmt.Printf("✅ Pesan terkirim: %s\n", h.truncateString(text, 50))
}

// Helper functions untuk pesan informatif

func (h *MessageHandler) getHelpMessage() string {
	return `📋 *BANTUAN WHATSAPP BOT*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
          *COMMAND TERSEDIA*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤖 *BASIC COMMANDS*

• */start*
  _Mulai bot_

• */help*
  _Bantuan lengkap_

• */ping*
  _Test koneksi bot_

• */info*
  _Informasi tentang bot_

• */status*
  _Status bot saat ini_

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 *TIPS PENGGUNAAN*

• *Chat Personal:* Bot membalas semua pesan
• *Di Grup:* Bot hanya respon command/mention
• *Command:* Ketik tanpa parameter untuk info

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📞 *SUPPORT:* Hubungi admin jika ada masalah`
}

func (h *MessageHandler) getInfoMessage() string {
	return `ℹ️ *INFORMASI BOT*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           *DETAIL SISTEM*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤖 *SPESIFIKASI*
📝 *Nama:* WhatsApp Bot
💻 *Bahasa:* Go (Golang)
📚 *Library:* whatsmeow + go-qrcode
✨ *Versi:* 1.0.0
🎯 *Fitur:* Visual QR, Auto-reply, Commands

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔧 *KONFIGURASI AKTIF*
• *Auto Reply Personal:* Aktif
• *Auto Reply Group:* Tidak aktif
• *Session:* Tersimpan otomatis
• *QR Code:* Visual display

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 Bot ini dibuat untuk pembelajaran dan automasi WhatsApp`
}

func (h *MessageHandler) getStatusMessage() string {
	return fmt.Sprintf(`📊 *STATUS BOT*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           *SISTEM STATUS*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔋 *STATUS UTAMA*
✅ *Status:* Online dan aktif
🔗 *Koneksi:* Terhubung ke WhatsApp
💾 *Session:* Tersimpan di database
🤖 *Bot ID:* %s

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚙️ *KONFIGURASI*
📱 *Auto Reply Personal:* %v
👥 *Auto Reply Group:* %v

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🟢 *Semua sistem berjalan normal!*`,
		h.client.Store.ID.User,
		h.autoReplyPersonal,
		h.autoReplyGroup)
}

// truncateString memotong string jika terlalu panjang untuk logging
func (h *MessageHandler) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// isAutoPromoteCommand mengecek apakah pesan adalah command auto promote
func (h *MessageHandler) isAutoPromoteCommand(messageText string) bool {
	if h.promoteCommandHandler == nil {
		return false
	}
	return h.promoteCommandHandler.IsPromoteCommand(messageText)
}

// handleAutoPromoteCommand menangani command auto promote
func (h *MessageHandler) handleAutoPromoteCommand(evt *events.Message, messageText string) string {
	lowerText := strings.ToLower(strings.TrimSpace(messageText))

	// Cek apakah ini admin command
	adminCommands := []string{
		// Group Management Commands
		".listgroups", ".enablegroup", ".enablemulti", ".disablegroup", ".groupstatus", ".testgroup",
		// Template Management Commands
		".addtemplate", ".edittemplate", ".deletetemplate", ".templatestats", ".promotestats", ".activegroups", ".fetchproducts", ".productstats", ".deleteall", ".deletemulti"}
	for _, cmd := range adminCommands {
		if strings.HasPrefix(lowerText, cmd) {
			if h.adminCommandHandler != nil {
				return h.adminCommandHandler.HandleAdminCommands(evt, messageText)
			}
			return `❌ *AKSES DITOLAK*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
			        *TIDAK ADA IZIN*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🚫 Command ini hanya bisa digunakan oleh admin.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 *INFORMASI*
• Hanya admin yang memiliki akses.
• Hubungi admin untuk bantuan.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔒 *Akses terbatas untuk keamanan sistem*`
		}
	}

	// Cek apakah ini template command yang juga perlu admin access
	templateCommands := []string{".listtemplates", ".alltemplates", ".previewtemplate", ".help"}
	for _, cmd := range templateCommands {
		if strings.HasPrefix(lowerText, cmd) {
			// Semua command auto promote sekarang hanya untuk admin
			if h.adminCommandHandler != nil {
				// Cek apakah user adalah admin
				if !h.isUserAdmin(evt.Info.Sender.User) {
					return `❌ *AKSES DITOLAK*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
			        *TIDAK ADA IZIN*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🚫 Command ini hanya bisa digunakan oleh admin.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 *INFORMASI*
• Hanya admin yang memiliki akses.
• Hubungi admin untuk bantuan.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔒 *Akses terbatas untuk keamanan sistem*`
				}
				return h.promoteCommandHandler.HandlePromoteCommands(evt, messageText)
			}
			return "" // Tidak ada response jika handler tidak tersedia
		}
	}

	return "" // Tidak ada response untuk command yang tidak dikenal
}

// isUserAdmin mengecek apakah user adalah admin
func (h *MessageHandler) isUserAdmin(userNumber string) bool {
	if h.adminCommandHandler == nil {
		return false
	}
	return h.adminCommandHandler.IsUserAdmin(userNumber)
}
