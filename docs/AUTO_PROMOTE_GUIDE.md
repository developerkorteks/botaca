# 🚀 Auto Promote System - Panduan Lengkap

## 📋 Daftar Isi
1. [Pengenalan](#pengenalan)
2. [Instalasi & Setup](#instalasi--setup)
3. [Konfigurasi](#konfigurasi)
4. [Commands User](#commands-user)
5. [Commands Admin](#commands-admin)
6. [Template System](#template-system)
7. [Troubleshooting](#troubleshooting)

## 🎯 Pengenalan

Auto Promote System adalah fitur untuk mengirim promosi bisnis secara otomatis ke grup-grup WhatsApp setiap 4 jam sekali. Sistem ini dirancang untuk:

- ✅ Promosi bisnis/jualan otomatis
- ✅ Template promosi yang bervariasi
- ✅ Management template oleh admin
- ✅ Random selection untuk menghindari spam
- ✅ Kontrol per grup (aktif/nonaktif)

## 🛠️ Instalasi & Setup

### 1. Persiapan
```bash
# Clone repository
git clone <repository-url>
cd promote

# Install dependencies
go mod tidy

# Build aplikasi
go build -o bot cmd/main.go
```

### 2. Konfigurasi Environment Variables
Buat file `.env` atau set environment variables:

```bash
# Database paths
DB_PATH=session.db
PROMOTE_DB_PATH=promote.db

# Auto promote settings
ENABLE_AUTO_PROMOTE=true
AUTO_PROMOTE_INTERVAL=4
LOG_AUTO_PROMOTE=true

# Admin numbers (pisahkan dengan koma)
ADMIN_NUMBERS=628123456789,628987654321

# Bot settings
LOG_LEVEL=INFO
AUTO_REPLY_PERSONAL=true
AUTO_REPLY_GROUP=false
```

### 3. Jalankan Bot
```bash
# Jalankan bot
./bot

# Atau langsung dengan go run
go run cmd/main.go
```

## ⚙️ Konfigurasi

### Admin Numbers
Ganti nomor admin di `config/promote_config.go`:
```go
return []string{
    "628123456789", // Nomor admin utama (GANTI DENGAN NOMOR ANDA)
    "628987654321", // Nomor admin kedua (opsional)
}
```

### Interval Auto Promote
Default: 4 jam. Bisa diubah via environment variable:
```bash
AUTO_PROMOTE_INTERVAL=6  # 6 jam
```

## 👥 Commands User

### Basic Commands
```
.promote              - Aktifkan auto promote di grup
.disablepromote       - Nonaktifkan auto promote di grup
.statuspromo          - Cek status auto promote grup
.testpromo            - Test kirim promosi manual
.promotehelp          - Bantuan lengkap auto promote
```

### Template Commands
```
.listtemplates        - Lihat daftar template promosi
.previewtemplate [ID] - Preview template berdasarkan ID
```

### Contoh Penggunaan
```
User: .promote
Bot: ✅ AUTO PROMOTE DIAKTIFKAN! 🚀
     Promosi akan dikirim setiap 4 jam...

User: .statuspromo
Bot: 📊 STATUS AUTO PROMOTE
     ✅ Status: Aktif
     📅 Dimulai: 2024-01-01 10:00
     ⏰ Promosi Terakhir: 2024-01-01 14:00

User: .disablepromote
Bot: 🛑 AUTO PROMOTE DINONAKTIFKAN!
```

## 👑 Commands Admin

### Template Management
```
.addtemplate "Judul" "Kategori" "Konten"
.edittemplate [ID] "Judul" "Kategori" "Konten"
.deletetemplate [ID]
.templatestats
```

### System Management
```
.promotestats         - Statistik auto promote
.activegroups         - Lihat grup yang aktif
```

### Contoh Admin Commands
```
Admin: .addtemplate "Flash Sale" "diskon" "🔥 FLASH SALE! Diskon 50% hari ini! Order: 08123456789"
Bot: ✅ TEMPLATE BERHASIL DIBUAT!
     🆔 ID: 11
     🏷️ Judul: Flash Sale

Admin: .edittemplate 11 "Super Flash Sale" "diskon" "⚡ SUPER FLASH SALE! Diskon 70%!"
Bot: ✅ TEMPLATE BERHASIL DIUPDATE!

Admin: .templatestats
Bot: 📊 STATISTIK TEMPLATE
     📝 Total Template: 12
     ✅ Aktif: 10
     ❌ Tidak Aktif: 2
```

## 📝 Template System

### Template Variables
Template mendukung variables yang akan diganti otomatis:
- `{DATE}` - Tanggal saat ini (2024-01-01)
- `{TIME}` - Waktu saat ini (14:30)
- `{DAY}` - Hari (Senin, Selasa, dll)
- `{MONTH}` - Bulan (Januari, Februari, dll)
- `{YEAR}` - Tahun (2024)

### Kategori Template Default
1. **produk** - Promosi produk unggulan
2. **diskon** - Penawaran diskon dan promo
3. **testimoni** - Review dan testimoni customer
4. **flashsale** - Flash sale dan promosi terbatas
5. **produk_baru** - Launching produk baru
6. **bundle** - Paket bundling hemat
7. **ongkir** - Promosi gratis ongkir
8. **cashback** - Program cashback dan reward
9. **limited** - Stok terbatas
10. **contact** - Informasi kontak dan cara order

### Contoh Template
```
🔥 *PRODUK UNGGULAN HARI INI* 🔥

✨ Dapatkan produk terbaik dengan kualitas premium!
💎 Harga terjangkau, kualitas terjamin
🚀 Stok terbatas, jangan sampai kehabisan!

📱 *Order sekarang:*
💬 WhatsApp: 08123456789
🛒 Link: bit.ly/produk-unggulan

#ProdukUnggulan #KualitasPremium #OrderSekarang
```

### Tips Template
- ✅ Gunakan emoji untuk menarik perhatian
- ✅ Sertakan call-to-action yang jelas
- ✅ Tambahkan kontak/link order
- ✅ Gunakan hashtag untuk branding
- ✅ Maksimal 4000 karakter
- ❌ Jangan terlalu spam
- ❌ Hindari konten menyesatkan

## 🔧 Troubleshooting

### Bot Tidak Merespon Commands
1. Pastikan bot sudah terhubung ke WhatsApp
2. Cek apakah auto promote diaktifkan (`ENABLE_AUTO_PROMOTE=true`)
3. Restart bot jika perlu

### Auto Promote Tidak Jalan
1. Cek status dengan `.statuspromo`
2. Pastikan ada template aktif (`.listtemplates`)
3. Cek log bot untuk error
4. Pastikan interval sudah lewat (4 jam)

### Template Tidak Muncul
1. Cek apakah template aktif
2. Pastikan kategori benar
3. Cek dengan `.previewtemplate [ID]`

### Admin Commands Tidak Bisa
1. Pastikan nomor Anda terdaftar sebagai admin
2. Cek konfigurasi `ADMIN_NUMBERS`
3. Restart bot setelah ubah konfigurasi

### Database Error
1. Pastikan file database bisa ditulis
2. Cek permission folder
3. Hapus file database dan restart (akan create ulang)

## 📊 Monitoring

### Log Files
Bot akan menampilkan log real-time:
```
2024/01/01 10:00:00 [INFO] Auto Promote System initialized!
2024/01/01 10:00:00 [SUCCESS] 🚀 Auto Promote System is READY!
2024/01/01 14:00:00 [INFO] Processing scheduled promotes...
2024/01/01 14:00:00 [INFO] Promote message sent to group: xxx
```

### Database
- `session.db` - Session WhatsApp
- `promote.db` - Data auto promote (grup, template, log)

### Backup
Backup file database secara berkala:
```bash
cp promote.db promote_backup_$(date +%Y%m%d).db
```

## 🚨 Peringatan Penting

1. **Rate Limiting**: WhatsApp memiliki rate limit. Jangan terlalu agresif.
2. **Spam Prevention**: Gunakan interval minimal 4 jam.
3. **Group Permission**: Pastikan bot tidak di-kick dari grup.
4. **Content Policy**: Ikuti kebijakan WhatsApp tentang konten.
5. **Backup Data**: Selalu backup database template.

## 🆘 Support

Jika mengalami masalah:
1. Cek dokumentasi ini
2. Lihat log error di terminal
3. Restart bot
4. Hubungi developer

## 📈 Tips Sukses

1. **Konten Berkualitas**: Buat template yang menarik dan informatif
2. **Timing**: Sesuaikan waktu promosi dengan target audience
3. **Variasi**: Gunakan banyak template untuk variasi
4. **Monitoring**: Pantau performa dan feedback
5. **Update**: Update template secara berkala

---

**Happy Promoting!** 🚀💰