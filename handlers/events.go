// Package handlers - File events.go
// File ini menangani semua event dari WhatsApp seperti koneksi, disconnection, dll
package handlers

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// EventHandler adalah struktur yang menangani semua event WhatsApp
type EventHandler struct {
	// client adalah instance WhatsApp client
	client *whatsmeow.Client
	
	// messageHandler untuk menangani pesan masuk
	messageHandler *MessageHandler
}

// NewEventHandler membuat handler baru untuk event WhatsApp
// Parameter:
// - client: WhatsApp client yang sudah terhubung
// - messageHandler: Handler khusus untuk pesan
func NewEventHandler(client *whatsmeow.Client, messageHandler *MessageHandler) *EventHandler {
	return &EventHandler{
		client:         client,
		messageHandler: messageHandler,
	}
}

// HandleEvent adalah fungsi utama yang menangani semua event dari WhatsApp
// Fungsi ini akan dipanggil setiap kali ada event baru (pesan, koneksi, dll)
func (h *EventHandler) HandleEvent(evt interface{}) {
	// Switch statement untuk menangani berbagai jenis event
	switch v := evt.(type) {
	
	case *events.Message:
		// Event pesan masuk - delegate ke message handler
		fmt.Println("📨 Event: Pesan masuk")
		h.messageHandler.HandleMessage(v)
		
	case *events.Connected:
		// Event ketika bot berhasil terhubung ke WhatsApp
		fmt.Println("✅ Event: Terhubung ke WhatsApp")
		h.handleConnected(v)
		
	case *events.Disconnected:
		// Event ketika bot terputus dari WhatsApp
		fmt.Println("❌ Event: Terputus dari WhatsApp")
		h.handleDisconnected(v)
		
	case *events.LoggedOut:
		// Event ketika bot di-logout dari WhatsApp
		fmt.Println("🚪 Event: Logged out dari WhatsApp")
		h.handleLoggedOut(v)
		
	case *events.StreamReplaced:
		// Event ketika ada session lain yang login dengan nomor yang sama
		fmt.Println("🔄 Event: Stream replaced (ada login dari device lain)")
		h.handleStreamReplaced(v)
		
	case *events.Receipt:
		// Event receipt (pesan terkirim, dibaca, dll) - biasanya tidak perlu ditangani
		// Uncomment jika ingin log receipt
		// fmt.Printf("📋 Event: Receipt - %s dari %s\n", v.Type, v.SourceString())
		
	case *events.Presence:
		// Event presence (online, offline, typing) - biasanya tidak perlu ditangani
		// Uncomment jika ingin log presence
		// fmt.Printf("👁️ Event: Presence - %s dari %s\n", v.Presence, v.From.String())
		
	case *events.GroupInfo:
		// Event perubahan info grup (nama, deskripsi, dll)
		fmt.Printf("👥 Event: Group info changed - %s\n", v.JID.String())
		h.handleGroupInfo(v)
		
	case *events.JoinedGroup:
		// Event ketika bot ditambahkan ke grup
		fmt.Printf("🎉 Event: Bot ditambahkan ke grup - %s\n", v.JID.String())
		h.handleJoinedGroup(v)
		
	// case *events.LeftGroup:
	//	// Event ketika bot dikeluarkan dari grup (API berubah di versi terbaru)
	//	fmt.Printf("👋 Event: Bot dikeluarkan dari grup - %s\n", v.JID.String())
	//	h.handleLeftGroup(v)
		
	default:
		// Event lain yang tidak ditangani khusus
		// Uncomment untuk debugging jika ingin melihat semua event
		// fmt.Printf("🔍 Event tidak ditangani: %T\n", v)
	}
}

// handleConnected menangani event ketika bot terhubung
func (h *EventHandler) handleConnected(evt *events.Connected) {
	fmt.Println("🎉 Bot berhasil terhubung ke WhatsApp!")
	fmt.Printf("📱 Device: %s\n", h.client.Store.ID.String())
	fmt.Println("💬 Bot siap menerima pesan...")
}

// handleDisconnected menangani event ketika bot terputus
func (h *EventHandler) handleDisconnected(evt *events.Disconnected) {
	fmt.Println("⚠️ Bot terputus dari WhatsApp")
	
	// Bot akan otomatis mencoba reconnect berkat whatsmeow
	fmt.Println("🔄 Bot akan mencoba reconnect otomatis...")
}

// handleLoggedOut menangani event ketika bot di-logout
func (h *EventHandler) handleLoggedOut(evt *events.LoggedOut) {
	fmt.Println("🚪 Bot telah di-logout dari WhatsApp")
	
	// Log alasan logout jika ada
	if evt.Reason != events.ConnectFailureLoggedOut {
		fmt.Printf("📝 Alasan logout: %v\n", evt.Reason)
	}
	
	fmt.Println("⚠️ Anda perlu scan QR code lagi untuk login")
	fmt.Println("💡 Restart bot untuk mendapatkan QR code baru")
}

// handleStreamReplaced menangani event ketika ada session lain yang login
func (h *EventHandler) handleStreamReplaced(evt *events.StreamReplaced) {
	fmt.Println("🔄 Session digantikan oleh login dari device lain")
	fmt.Println("⚠️ Bot akan disconnect untuk menghindari konflik")
	
	// Biasanya bot akan otomatis disconnect setelah event ini
}

// handleGroupInfo menangani event perubahan info grup
func (h *EventHandler) handleGroupInfo(evt *events.GroupInfo) {
	fmt.Printf("👥 Info grup berubah: %s\n", evt.JID.String())
	
	// Anda bisa menambahkan logic khusus di sini, misalnya:
	// - Log perubahan nama grup
	// - Notifikasi admin jika ada perubahan penting
	// - Update database internal jika ada
}

// handleJoinedGroup menangani event ketika bot ditambahkan ke grup
func (h *EventHandler) handleJoinedGroup(evt *events.JoinedGroup) {
	fmt.Printf("🎉 Bot ditambahkan ke grup: %s\n", evt.JID.String())
	
	// Anda bisa menambahkan logic khusus di sini, misalnya:
	// - Kirim pesan perkenalan ke grup
	// - Log grup baru ke database
	// - Notifikasi admin
	
	// Contoh: kirim pesan perkenalan (uncomment jika ingin diaktifkan)
	/*
	welcomeMsg := `👋 Halo! Saya adalah WhatsApp Bot.
	
🤖 Ketik /help untuk melihat command yang tersedia.
📝 Bot hanya akan merespon command atau mention untuk menghindari spam.

Terima kasih telah menambahkan saya ke grup! 🎉`
	
	h.messageHandler.sendMessage(evt.JID, welcomeMsg)
	*/
}

// handleLeftGroup menangani event ketika bot dikeluarkan dari grup
// Fungsi ini di-comment karena events.LeftGroup tidak tersedia di versi whatsmeow terbaru
// func (h *EventHandler) handleLeftGroup(evt *events.LeftGroup) {
//	fmt.Printf("👋 Bot dikeluarkan dari grup: %s\n", evt.JID.String())
//	
//	// Anda bisa menambahkan logic khusus di sini, misalnya:
//	// - Log grup yang ditinggalkan
//	// - Cleanup data terkait grup
//	// - Notifikasi admin
// }