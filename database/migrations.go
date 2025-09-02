// Package database - migrations untuk auto promote feature
package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations menjalankan semua migrasi database yang diperlukan
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createAutoPromoteGroupsTable,
		createPromoteTemplatesTable,
		createPromoteLogsTable,
		createPromoteStatsTable,
		insertDefaultTemplates,
	}

	for i, migration := range migrations {
		fmt.Printf("Running migration %d/%d...\n", i+1, len(migrations))
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %v", i+1, err)
		}
	}

	fmt.Println("✅ All migrations completed successfully!")
	return nil
}

// SQL untuk membuat tabel auto_promote_groups
const createAutoPromoteGroupsTable = `
CREATE TABLE IF NOT EXISTS auto_promote_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_jid TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    started_at DATETIME,
    last_promote_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auto_promote_groups_jid ON auto_promote_groups(group_jid);
CREATE INDEX IF NOT EXISTS idx_auto_promote_groups_active ON auto_promote_groups(is_active);
`

// SQL untuk membuat tabel promote_templates
const createPromoteTemplatesTable = `
CREATE TABLE IF NOT EXISTS promote_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'general',
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_promote_templates_active ON promote_templates(is_active);
CREATE INDEX IF NOT EXISTS idx_promote_templates_category ON promote_templates(category);
`

// SQL untuk membuat tabel promote_logs
const createPromoteLogsTable = `
CREATE TABLE IF NOT EXISTS promote_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_jid TEXT NOT NULL,
    template_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN DEFAULT TRUE,
    error_msg TEXT,
    FOREIGN KEY (template_id) REFERENCES promote_templates(id)
);

CREATE INDEX IF NOT EXISTS idx_promote_logs_group ON promote_logs(group_jid);
CREATE INDEX IF NOT EXISTS idx_promote_logs_sent_at ON promote_logs(sent_at);
CREATE INDEX IF NOT EXISTS idx_promote_logs_success ON promote_logs(success);
`

// SQL untuk membuat tabel promote_stats
const createPromoteStatsTable = `
CREATE TABLE IF NOT EXISTS promote_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT UNIQUE NOT NULL,
    total_groups INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    success_messages INTEGER DEFAULT 0,
    failed_messages INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_promote_stats_date ON promote_stats(date);
`

// SQL untuk insert template default
const insertDefaultTemplates = `
INSERT OR IGNORE INTO promote_templates (title, content, category, is_active) VALUES
('Produk Unggulan', '🔥 *PRODUK UNGGULAN HARI INI* 🔥

✨ Dapatkan produk terbaik dengan kualitas premium!
💎 Harga terjangkau, kualitas terjamin
🚀 Stok terbatas, jangan sampai kehabisan!

📱 *Order sekarang:*
💬 WhatsApp: 08123456789
🛒 Link: bit.ly/produk-unggulan

#ProdukUnggulan #KualitasPremium #OrderSekarang', 'produk', 1),

('Diskon & Promo', '🎉 *PROMO SPESIAL HARI INI* 🎉

💥 DISKON hingga 50% untuk semua produk!
⏰ Promo terbatas hanya sampai hari ini
🎁 Bonus gratis untuk pembelian minimal 100k

🛍️ *Jangan lewatkan kesempatan emas ini!*
📞 Order: 08123456789
💳 Pembayaran mudah & aman

#PromoSpesial #Diskon50Persen #TerbatasWaktu', 'diskon', 1),

('Testimoni Customer', '⭐ *TESTIMONI CUSTOMER SETIA* ⭐

💬 "Produknya bagus banget, sesuai ekspektasi!"
👤 - Bu Sarah, Jakarta

💬 "Pelayanan ramah, pengiriman cepat!"
👤 - Pak Budi, Surabaya

💬 "Harga murah, kualitas juara!"
👤 - Mbak Siti, Bandung

🙏 Terima kasih kepercayaannya!
📱 Order: 08123456789

#TestimoniCustomer #KepuasanPelanggan #Terpercaya', 'testimoni', 1),

('Flash Sale', '⚡ *FLASH SALE ALERT!* ⚡

🔥 HANYA 2 JAM LAGI!
💰 Harga super murah, stok terbatas!
⏰ Berakhir pukul 23:59 WIB

🎯 *Yang tersisa:*
• Produk A: 5 pcs tersisa
• Produk B: 3 pcs tersisa
• Produk C: 8 pcs tersisa

💨 BURUAN ORDER SEBELUM KEHABISAN!
📱 WhatsApp: 08123456789

#FlashSale #StokTerbatas #BuruanOrder', 'flashsale', 1),

('Produk Baru', '🆕 *LAUNCHING PRODUK TERBARU!* 🆕

🎊 Kami bangga memperkenalkan inovasi terbaru!
✨ Fitur canggih, desain modern
🏆 Kualitas terbaik di kelasnya

🎁 *PROMO LAUNCHING:*
• Diskon 30% untuk 100 pembeli pertama
• Gratis ongkir seluruh Indonesia
• Garansi resmi 1 tahun

📱 Pre-order: 08123456789
🚀 Jadilah yang pertama memilikinya!

#ProdukBaru #Launching #PreOrder', 'produk_baru', 1),

('Bundle Package', '📦 *PAKET HEMAT BUNDLE!* 📦

💡 Beli 1 dapat 3? Why not!
🎯 Hemat hingga 40% dari harga normal
🎁 Bonus eksklusif untuk paket lengkap

📋 *Paket yang tersedia:*
• Paket A: 3 produk = 150k (normal 250k)
• Paket B: 5 produk = 200k (normal 350k)
• Paket C: 10 produk = 350k (normal 600k)

💰 Makin banyak makin hemat!
📱 Order: 08123456789

#BundlePackage #PaketHemat #MakinBanyakMakinHemat', 'bundle', 1),

('Free Ongkir', '🚚 *GRATIS ONGKIR SELURUH INDONESIA!* 🚚

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

#GratisOngkir #PengirimanAman #OrderSekarang', 'ongkir', 1),

('Cashback & Reward', '💰 *PROGRAM CASHBACK & REWARD!* 💰

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

#CashbackReward #MemberExclusive #BelanjaMakinUntung', 'cashback', 1),

('Limited Stock', '⚠️ *STOK TERBATAS - SEGERA HABIS!* ⚠️

🔥 Produk favorit hampir sold out!
📊 Sisa stok: 7 pcs saja
⏰ Kemungkinan habis dalam 24 jam

😱 *Jangan sampai menyesal!*
• Produk best seller #1
• Rating 5 bintang dari customer
• Sudah terjual 500+ pcs bulan ini

🏃‍♂️ BURUAN ORDER SEBELUM KEHABISAN!
📱 WhatsApp: 08123456789

#StokTerbatas #BestSeller #BuruanOrder', 'limited', 1),

('Contact Info', '📞 *HUBUNGI KAMI UNTUK ORDER!* 📞

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

#ContactInfo #CaraOrder #JamOperasional', 'contact', 1);
`