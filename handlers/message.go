// Package handlers berisi semua handler untuk menangani pesan dan event WhatsApp
// File ini khusus menangani pesan masuk dari chat personal dan grup
package handlers

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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
	
	// Jika bukan command dan auto reply personal diaktifkan
	if h.autoReplyPersonal {
		h.sendAutoReply(evt.Info.Chat, messageText, false)
	}
}

// handleGroupMessage menangani pesan dari grup
func (h *MessageHandler) handleGroupMessage(evt *events.Message, messageText string) {
	fmt.Println("👥 Memproses pesan grup...")
	
	// PENTING: Untuk grup, kita hanya merespon jika:
	// 1. Pesan adalah command (dimulai dengan /)
	// 2. Bot di-mention (@bot)
	// 3. Auto reply grup diaktifkan (tidak direkomendasikan)
	
	// Cek apakah bot di-mention dalam pesan
	isMentioned := h.isBotMentioned(evt.Message)
	
	// Cek apakah ini adalah command (/ atau .)
	isCommand := strings.HasPrefix(messageText, "/") || strings.HasPrefix(messageText, ".")
	
	if isCommand {
		// Selalu proses command di grup
		h.handleCommand(evt, messageText)
	} else if isMentioned {
		// Jika bot di-mention, balas meskipun auto reply grup dimatikan
		h.sendAutoReply(evt.Info.Chat, messageText, true)
	} else if h.autoReplyGroup {
		// Hanya auto reply jika diaktifkan (HATI-HATI: bisa spam!)
		h.sendAutoReply(evt.Info.Chat, messageText, true)
	}
	// Jika tidak ada kondisi di atas yang terpenuhi, bot tidak akan membalas
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
		// Daftar command yang tersedia
		switch {
		case lowerText == "/start":
			response = "🤖 *WhatsApp Bot Aktif!*\n\n✨ Bot siap melayani Anda.\nKetik /help untuk melihat command yang tersedia."
			
		case lowerText == "/help":
			response = h.getHelpMessage()
			
		case lowerText == "/ping":
			response = "🏓 Pong! Bot aktif dan berjalan dengan baik."
			
		case lowerText == "/info":
			response = h.getInfoMessage()
			
		case lowerText == "/status":
			response = h.getStatusMessage()
			
		case strings.HasPrefix(lowerText, "/promote"):
			// Command promote untuk grup (akan diimplementasi nanti)
			response = "🔧 Fitur promote sedang dalam pengembangan."
			
		default:
			response = "❓ Command tidak dikenal. Ketik /help untuk melihat command yang tersedia."
		}
	}
	
	// Kirim response
	h.sendMessage(evt.Info.Chat, response)
}

// sendAutoReply mengirim balasan otomatis
func (h *MessageHandler) sendAutoReply(chatJID types.JID, originalMessage string, isGroup bool) {
	var response string
	
	if isGroup {
		// Response untuk grup lebih formal dan tidak terlalu sering
		response = "👋 Terima kasih! Saya adalah bot otomatis. Ketik /help untuk bantuan."
	} else {
		// Response untuk personal bisa lebih personal
		responses := []string{
			"✅ Terima kasih atas pesannya! Ketik /help untuk melihat command yang tersedia.",
			"🤖 Pesan diterima! Saya adalah bot otomatis yang siap membantu.",
			"👍 Got it! Kirim /help untuk melihat apa yang bisa saya lakukan.",
		}
		
		// Pilih response berdasarkan panjang pesan untuk variasi
		responseIndex := len(originalMessage) % len(responses)
		response = responses[responseIndex]
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
	return `📋 *Bantuan WhatsApp Bot*

🤖 *Command yang tersedia:*
• /start - Mulai bot
• /help - Bantuan ini
• /ping - Test koneksi bot
• /info - Informasi tentang bot
• /status - Status bot saat ini
• /promote - Promote member grup (coming soon)

💡 *Tips:*
• Di chat personal: Bot akan membalas semua pesan
• Di grup: Bot hanya merespon command atau mention
• Ketik command tanpa parameter untuk info lebih lanjut

📞 *Support:* Hubungi admin jika ada masalah`
}

func (h *MessageHandler) getInfoMessage() string {
	return `ℹ️ *Informasi Bot*

🤖 Nama: WhatsApp Bot
📝 Bahasa: Go (Golang)
📚 Library: whatsmeow + go-qrcode
✨ Versi: 1.0.0
🎯 Fitur: Visual QR code, Auto-reply, Commands

🔧 *Konfigurasi Saat Ini:*
• Auto Reply Personal: Aktif
• Auto Reply Group: Tidak aktif (recommended)
• Session: Tersimpan otomatis

Bot ini dibuat untuk pembelajaran dan automasi WhatsApp.`
}

func (h *MessageHandler) getStatusMessage() string {
	return fmt.Sprintf(`📊 *Status Bot*

✅ Status: Online dan aktif
🔗 Koneksi: Terhubung ke WhatsApp
💾 Session: Tersimpan di database
🤖 Bot ID: %s
📱 Auto Reply Personal: %v
👥 Auto Reply Group: %v

🟢 Semua sistem berjalan normal!`, 
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
	adminCommands := []string{".addtemplate", ".edittemplate", ".deletetemplate", ".templatestats", ".promotestats", ".activegroups"}
	for _, cmd := range adminCommands {
		if strings.HasPrefix(lowerText, cmd) {
			if h.adminCommandHandler != nil {
				return h.adminCommandHandler.HandleAdminCommands(evt, messageText)
			}
			return "❌ Admin commands tidak tersedia"
		}
	}
	
	// Handle regular promote commands
	if h.promoteCommandHandler != nil {
		return h.promoteCommandHandler.HandlePromoteCommands(evt, messageText)
	}
	
	return "❌ Auto promote commands tidak tersedia"
}