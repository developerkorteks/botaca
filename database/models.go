// Package database berisi model dan struktur data untuk fitur auto promote
package database

import (
	"time"
)

// AutoPromoteGroup menyimpan status auto promote per grup
type AutoPromoteGroup struct {
	ID            int       `json:"id" db:"id"`
	GroupJID      string    `json:"group_jid" db:"group_jid"`           // JID grup WhatsApp
	IsActive      bool      `json:"is_active" db:"is_active"`           // Status aktif/tidak
	StartedAt     *time.Time `json:"started_at" db:"started_at"`        // Waktu mulai auto promote
	LastPromoteAt *time.Time `json:"last_promote_at" db:"last_promote_at"` // Waktu terakhir kirim promosi
	CreatedAt     time.Time `json:"created_at" db:"created_at"`         // Waktu dibuat
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`         // Waktu diupdate
}

// PromoteTemplate menyimpan template promosi bisnis
type PromoteTemplate struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`         // Judul template (misal: "Produk Unggulan")
	Content   string    `json:"content" db:"content"`     // Isi template promosi
	Category  string    `json:"category" db:"category"`   // Kategori (produk, diskon, testimoni, dll)
	IsActive  bool      `json:"is_active" db:"is_active"` // Status aktif/tidak
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PromoteLog menyimpan log pengiriman promosi untuk tracking
type PromoteLog struct {
	ID         int       `json:"id" db:"id"`
	GroupJID   string    `json:"group_jid" db:"group_jid"`     // JID grup tujuan
	TemplateID int       `json:"template_id" db:"template_id"` // ID template yang digunakan
	Content    string    `json:"content" db:"content"`         // Isi pesan yang dikirim
	SentAt     time.Time `json:"sent_at" db:"sent_at"`         // Waktu pengiriman
	Success    bool      `json:"success" db:"success"`         // Status berhasil/gagal
	ErrorMsg   *string   `json:"error_msg" db:"error_msg"`     // Pesan error jika gagal
}

// PromoteStats menyimpan statistik promosi untuk monitoring
type PromoteStats struct {
	ID              int       `json:"id" db:"id"`
	Date            string    `json:"date" db:"date"`                         // Tanggal (YYYY-MM-DD)
	TotalGroups     int       `json:"total_groups" db:"total_groups"`         // Total grup aktif
	TotalMessages   int       `json:"total_messages" db:"total_messages"`     // Total pesan terkirim
	SuccessMessages int       `json:"success_messages" db:"success_messages"` // Pesan berhasil
	FailedMessages  int       `json:"failed_messages" db:"failed_messages"`   // Pesan gagal
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// DefaultPromoteTemplates berisi template default untuk promosi bisnis
var DefaultPromoteTemplates = []PromoteTemplate{
	{
		Title:    "Produk Unggulan",
		Category: "produk",
		Content: `🔥 *PRODUK UNGGULAN HARI INI* 🔥

✨ Dapatkan produk terbaik dengan kualitas premium!
💎 Harga terjangkau, kualitas terjamin
🚀 Stok terbatas, jangan sampai kehabisan!

📱 *Order sekarang:*
💬 WhatsApp: 08123456789
🛒 Link: bit.ly/produk-unggulan

#ProdukUnggulan #KualitasPremium #OrderSekarang`,
		IsActive: true,
	},
	{
		Title:    "Diskon & Promo",
		Category: "diskon",
		Content: `🎉 *PROMO SPESIAL HARI INI* 🎉

💥 DISKON hingga 50% untuk semua produk!
⏰ Promo terbatas hanya sampai {DATE}
🎁 Bonus gratis untuk pembelian minimal 100k

🛍️ *Jangan lewatkan kesempatan emas ini!*
📞 Order: 08123456789
💳 Pembayaran mudah & aman

#PromoSpesial #Diskon50Persen #TerbatasWaktu`,
		IsActive: true,
	},
	{
		Title:    "Testimoni Customer",
		Category: "testimoni",
		Content: `⭐ *TESTIMONI CUSTOMER SETIA* ⭐

💬 "Produknya bagus banget, sesuai ekspektasi!"
👤 - Bu Sarah, Jakarta

💬 "Pelayanan ramah, pengiriman cepat!"
👤 - Pak Budi, Surabaya

💬 "Harga murah, kualitas juara!"
👤 - Mbak Siti, Bandung

🙏 Terima kasih kepercayaannya!
📱 Order: 08123456789

#TestimoniCustomer #KepuasanPelanggan #Terpercaya`,
		IsActive: true,
	},
	{
		Title:    "Flash Sale",
		Category: "flashsale",
		Content: `⚡ *FLASH SALE ALERT!* ⚡

🔥 HANYA 2 JAM LAGI!
💰 Harga super murah, stok terbatas!
⏰ Berakhir pukul 23:59 WIB

🎯 *Yang tersisa:*
• Produk A: 5 pcs tersisa
• Produk B: 3 pcs tersisa
• Produk C: 8 pcs tersisa

💨 BURUAN ORDER SEBELUM KEHABISAN!
📱 WhatsApp: 08123456789

#FlashSale #StokTerbatas #BuruanOrder`,
		IsActive: true,
	},
	{
		Title:    "Produk Baru",
		Category: "produk_baru",
		Content: `🆕 *LAUNCHING PRODUK TERBARU!* 🆕

🎊 Kami bangga memperkenalkan inovasi terbaru!
✨ Fitur canggih, desain modern
🏆 Kualitas terbaik di kelasnya

🎁 *PROMO LAUNCHING:*
• Diskon 30% untuk 100 pembeli pertama
• Gratis ongkir seluruh Indonesia
• Garansi resmi 1 tahun

📱 Pre-order: 08123456789
🚀 Jadilah yang pertama memilikinya!

#ProdukBaru #Launching #PreOrder`,
		IsActive: true,
	},
	{
		Title:    "Bundle Package",
		Category: "bundle",
		Content: `📦 *PAKET HEMAT BUNDLE!* 📦

💡 Beli 1 dapat 3? Why not!
🎯 Hemat hingga 40% dari harga normal
🎁 Bonus eksklusif untuk paket lengkap

📋 *Paket yang tersedia:*
• Paket A: 3 produk = 150k (normal 250k)
• Paket B: 5 produk = 200k (normal 350k)
• Paket C: 10 produk = 350k (normal 600k)

💰 Makin banyak makin hemat!
📱 Order: 08123456789

#BundlePackage #PaketHemat #MakinBanyakMakinHemat`,
		IsActive: true,
	},
	{
		Title:    "Free Ongkir",
		Category: "ongkir",
		Content: `🚚 *GRATIS ONGKIR SELURUH INDONESIA!* 🚚

🎉 Tanpa minimum pembelian!
📦 Pengiriman aman & terpercaya
⏰ Estimasi 1-3 hari kerja

🌟 *Keuntungan lainnya:*
• Packing aman & rapi
• Asuransi pengiriman
• Tracking number real-time
• Customer service 24/7

🛒 Order sekarang juga!
📱 WhatsApp: 08123456789

#GratisOngkir #PengirimanAman #OrderSekarang`,
		IsActive: true,
	},
	{
		Title:    "Cashback & Reward",
		Category: "cashback",
		Content: `💰 *PROGRAM CASHBACK & REWARD!* 💰

🎁 Belanja makin untung dengan reward points!
💎 Tukar poin dengan produk gratis
🔄 Cashback langsung ke rekening

🏆 *Benefit member:*
• Cashback 5% setiap pembelian
• Poin reward setiap transaksi
• Diskon eksklusif member
• Akses produk limited edition

👑 Daftar member sekarang!
📱 WhatsApp: 08123456789

#CashbackReward #MemberExclusive #BelanjaMakinUntung`,
		IsActive: true,
	},
	{
		Title:    "Limited Stock",
		Category: "limited",
		Content: `⚠️ *STOK TERBATAS - SEGERA HABIS!* ⚠️

🔥 Produk favorite hampir sold out!
📊 Sisa stok: 7 pcs saja
⏰ Kemungkinan habis dalam 24 jam

😱 *Jangan sampai menyesal!*
• Produk best seller #1
• Rating 5 bintang dari customer
• Sudah terjual 500+ pcs bulan ini

🏃‍♂️ BURUAN ORDER SEBELUM KEHABISAN!
📱 WhatsApp: 08123456789

#StokTerbatas #BestSeller #BuruanOrder`,
		IsActive: true,
	},
	{
		Title:    "Contact Info",
		Category: "contact",
		Content: `📞 *HUBUNGI KAMI UNTUK ORDER!* 📞

🛒 *Cara Order:*
1️⃣ WhatsApp: 08123456789
2️⃣ Telegram: @tokoonline
3️⃣ Instagram: @toko.online
4️⃣ Website: www.tokoonline.com

💳 *Pembayaran:*
• Transfer Bank (BCA, Mandiri, BRI)
• E-wallet (OVO, DANA, GoPay)
• COD (area tertentu)

⏰ *Jam Operasional:*
Senin-Sabtu: 08:00-22:00 WIB
Minggu: 10:00-20:00 WIB

#ContactInfo #CaraOrder #JamOperasional`,
		IsActive: true,
	},
}