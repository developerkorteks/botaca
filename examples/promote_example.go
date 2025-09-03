// Package examples - Contoh implementasi fitur promote
// File ini menunjukkan cara mengimplementasikan fitur promote grup
// yang terintegrasi dengan template di folder layout/
package examples

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

// PromoteHandler menangani fitur promote member grup
type PromoteHandler struct {
	client *whatsmeow.Client
}

// NewPromoteHandler membuat handler baru untuk fitur promote
func NewPromoteHandler(client *whatsmeow.Client) *PromoteHandler {
	return &PromoteHandler{
		client: client,
	}
}

// HandlePromoteCommand menangani command /promote di grup
// Command ini hanya bisa digunakan oleh admin grup
func (p *PromoteHandler) HandlePromoteCommand(evt *events.Message, messageText string) string {
	// STEP 1: Validasi bahwa ini adalah grup
	if evt.Info.Chat.Server != types.GroupServer {
		return "❌ Command /promote hanya bisa digunakan di grup!"
	}

	// STEP 2: Ambil informasi grup
	groupInfo, err := p.client.GetGroupInfo(evt.Info.Chat)
	if err != nil {
		return "❌ Gagal mendapatkan informasi grup. Pastikan bot masih ada di grup."
	}

	// STEP 3: Cek apakah bot adalah admin grup
	botJID := p.client.Store.ID
	if !p.isBotAdmin(groupInfo, botJID) {
		return "❌ Bot harus menjadi admin grup untuk menggunakan fitur promote!"
	}

	// STEP 4: Cek apakah pengirim command adalah admin
	if !p.isUserAdmin(groupInfo, evt.Info.Sender) {
		return "❌ Hanya admin yang bisa menggunakan command /promote!"
	}

	// STEP 5: Parse command untuk mendapatkan target user
	targetJID, err := p.parsePromoteTarget(evt, messageText)
	if err != nil {
		return fmt.Sprintf("❌ %s", err.Error())
	}

	// STEP 6: Cek apakah target user ada di grup
	if !p.isUserInGroup(groupInfo, targetJID) {
		return "❌ User tidak ditemukan di grup ini."
	}

	// STEP 7: Cek apakah target user sudah admin
	if p.isUserAdmin(groupInfo, targetJID) {
		return "❌ User sudah menjadi admin grup."
	}

	// STEP 8: Lakukan promote
	err = p.promoteUser(evt.Info.Chat, targetJID)
	if err != nil {
		return fmt.Sprintf("❌ Gagal mempromote user: %v", err)
	}

	// STEP 9: Return success message
	return fmt.Sprintf("✅ Berhasil mempromote @%s menjadi admin grup! 🎉", targetJID.User)
}

// HandleDemoteCommand menangani command /demote di grup
func (p *PromoteHandler) HandleDemoteCommand(evt *events.Message, messageText string) string {
	// Implementasi serupa dengan promote, tapi untuk demote
	// ... (implementasi lengkap bisa ditambahkan nanti)
	return "🔧 Fitur demote sedang dalam pengembangan."
}

// isBotAdmin mengecek apakah bot adalah admin grup
func (p *PromoteHandler) isBotAdmin(groupInfo *types.GroupInfo, botJID *types.JID) bool {
	for _, participant := range groupInfo.Participants {
		if participant.JID.User == botJID.User {
			return participant.IsAdmin || participant.IsSuperAdmin
		}
	}
	return false
}

// isUserAdmin mengecek apakah user adalah admin grup
func (p *PromoteHandler) isUserAdmin(groupInfo *types.GroupInfo, userJID types.JID) bool {
	for _, participant := range groupInfo.Participants {
		if participant.JID.User == userJID.User {
			return participant.IsAdmin || participant.IsSuperAdmin
		}
	}
	return false
}

// isUserInGroup mengecek apakah user ada di grup
func (p *PromoteHandler) isUserInGroup(groupInfo *types.GroupInfo, userJID types.JID) bool {
	for _, participant := range groupInfo.Participants {
		if participant.JID.User == userJID.User {
			return true
		}
	}
	return false
}

// parsePromoteTarget mengambil target user dari command atau reply
func (p *PromoteHandler) parsePromoteTarget(evt *events.Message, messageText string) (types.JID, error) {
	// CARA 1: Cek apakah ini reply ke pesan user lain
	if evt.Message.GetExtendedTextMessage() != nil && 
	   evt.Message.GetExtendedTextMessage().GetContextInfo() != nil {
		contextInfo := evt.Message.GetExtendedTextMessage().GetContextInfo()
		
		// Jika ada quoted message (reply)
		if contextInfo.GetStanzaId() != "" && contextInfo.GetParticipant() != "" {
			targetJID, err := types.ParseJID(contextInfo.GetParticipant())
			if err != nil {
				return types.JID{}, fmt.Errorf("format JID tidak valid")
			}
			return targetJID, nil
		}
		
		// Jika ada mention dalam pesan
		mentions := contextInfo.GetMentionedJid()
		if len(mentions) > 0 {
			targetJID, err := types.ParseJID(mentions[0])
			if err != nil {
				return types.JID{}, fmt.Errorf("format JID mention tidak valid")
			}
			return targetJID, nil
		}
	}

	// CARA 2: Parse nomor dari command text
	parts := strings.Fields(messageText)
	if len(parts) >= 2 {
		// Format: /promote 628123456789 atau /promote @628123456789
		phoneNumber := strings.TrimPrefix(parts[1], "@")
		
		// Validasi format nomor (harus angka dan dimulai dengan kode negara)
		if !p.isValidPhoneNumber(phoneNumber) {
			return types.JID{}, fmt.Errorf("format nomor tidak valid. Gunakan format: 628123456789")
		}
		
		targetJID := types.NewJID(phoneNumber, types.DefaultUserServer)
		return targetJID, nil
	}

	return types.JID{}, fmt.Errorf("silakan mention user, reply pesan user, atau tulis nomor. Format: /promote @user atau /promote 628123456789")
}

// isValidPhoneNumber validasi sederhana untuk nomor telepon
func (p *PromoteHandler) isValidPhoneNumber(phone string) bool {
	// Cek apakah semua karakter adalah angka
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	// Cek panjang nomor (minimal 10, maksimal 15 digit)
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	
	return true
}

// promoteUser melakukan promote user menjadi admin
func (p *PromoteHandler) promoteUser(groupJID types.JID, targetJID types.JID) error {
	// Gunakan whatsmeow API untuk promote user
	_, err := p.client.UpdateGroupParticipants(
		groupJID,
		[]types.JID{targetJID},
		whatsmeow.ParticipantChangePromote,
	)
	
	return err
}

// demoteUser melakukan demote admin menjadi member biasa
func (p *PromoteHandler) demoteUser(groupJID types.JID, targetJID types.JID) error {
	// Gunakan whatsmeow API untuk demote user
	_, err := p.client.UpdateGroupParticipants(
		groupJID,
		[]types.JID{targetJID},
		whatsmeow.ParticipantChangeDemote,
	)
	
	return err
}

// GetPromoteHelp mengembalikan help message untuk fitur promote
func (p *PromoteHandler) GetPromoteHelp() string {
	return `🔧 *Fitur Promote/Demote*

📋 *Commands yang tersedia:*
• /promote @user - Promote user menjadi admin
• /promote 628123456789 - Promote dengan nomor
• /demote @user - Demote admin menjadi member

💡 *Cara penggunaan:*
1. **Mention user:** /promote @username
2. **Reply pesan:** Reply pesan user dengan /promote  
3. **Tulis nomor:** /promote 628123456789

⚠️ *Syarat:*
• Bot harus menjadi admin grup
• Anda harus menjadi admin grup
• Target user harus ada di grup

📞 *Contoh:*
/promote @john_doe
/promote 628123456789`
}

// INTEGRASI DENGAN FOLDER LAYOUT
// Fungsi-fungsi di bawah ini menunjukkan cara mengintegrasikan
// dengan template yang sudah Anda buat di folder layout/

// LoadPromoteTemplate memuat template promote dari folder layout/
func (p *PromoteHandler) LoadPromoteTemplate(templateName string) (string, error) {
	// Contoh implementasi untuk load template
	// Anda bisa sesuaikan dengan struktur template di folder layout/
	
	// Baca file template (implementasi sederhana)
	// content, err := ioutil.ReadFile(templatePath)
	// if err != nil {
	//     return "", err
	// }
	
	// return string(content), nil
	
	// Untuk sekarang, return template default
	return p.getDefaultPromoteTemplate(), nil
}

// getDefaultPromoteTemplate template default untuk promote
func (p *PromoteHandler) getDefaultPromoteTemplate() string {
	return `🎉 *SELAMAT!* 🎉

👤 User: @{USER}
🏆 Status: Dipromote menjadi Admin
📅 Tanggal: {DATE}
👑 Oleh: @{PROMOTER}

✨ Selamat bergabung dengan tim admin!
🤝 Mari kita jaga grup ini bersama-sama.

#Promote #NewAdmin #TeamWork`
}

// GeneratePromoteMessage generate pesan promote dengan template
func (p *PromoteHandler) GeneratePromoteMessage(targetUser, promoterUser string) string {
	template := p.getDefaultPromoteTemplate()
	
	// Replace placeholder dengan data actual
	template = strings.ReplaceAll(template, "{USER}", targetUser)
	template = strings.ReplaceAll(template, "{PROMOTER}", promoterUser)
	template = strings.ReplaceAll(template, "{DATE}", "2024-01-01") // Ganti dengan tanggal actual
	
	return template
}

// SendPromoteNotification kirim notifikasi promote dengan template
func (p *PromoteHandler) SendPromoteNotification(groupJID types.JID, targetJID, promoterJID types.JID) error {
	// Generate pesan dengan template
	message := p.GeneratePromoteMessage(targetJID.User, promoterJID.User)
	
	// Buat pesan dengan mention
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &message,
			ContextInfo: &waProto.ContextInfo{
				MentionedJID: []string{targetJID.String(), promoterJID.String()},
			},
		},
	}
	
	// Kirim pesan
	_, err := p.client.SendMessage(context.Background(), groupJID, msg)
	return err
}