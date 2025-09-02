# 🤖 WhatsApp Bot dengan Auto Promote System

Bot WhatsApp yang dibuat dengan Go menggunakan library `whatsmeow`. Project ini dirancang dengan struktur yang rapi dan mudah dipelajari, lengkap dengan dokumentasi detail untuk setiap bagian kode.

## 📁 Struktur Project

```
promote/
├── cmd/
│   └── main.go              # Entry point aplikasi
├── config/
│   └── config.go            # Konfigurasi bot (database, auto-reply, dll)
├── handlers/
│   ├── events.go            # Handler untuk event WhatsApp
│   └── message.go           # Handler untuk pesan masuk
├── utils/
│   ├── logger.go            # Utility untuk logging
│   └── qrcode.go            # Utility untuk QR code visual
├── docs/
│   └── (dokumentasi akan ditambahkan)
├── examples/
│   └── (contoh penggunaan akan ditambahkan)
├── layout/
│   └── (template promote yang sudah Anda buat)
├── go.mod                   # Dependencies Go
├── go.sum                   # Checksum dependencies
└── README.md               # Dokumentasi ini
```

## 🚀 Cara Menjalankan

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Jalankan Bot
```bash
cd cmd
go run main.go
```

### 3. Scan QR Code
- QR code akan muncul di terminal secara visual
- Scan dengan WhatsApp: Settings > Linked Devices > Link a Device
- File `qrcode.png` juga akan tersimpan sebagai backup

## ⚙️ Konfigurasi

Bot dapat dikonfigurasi melalui file `config/config.go` atau environment variables:

| Setting | Default | Deskripsi |
|---------|---------|-----------|
| `DB_PATH` | `session.db` | Lokasi file database session |
| `LOG_LEVEL` | `INFO` | Level logging (DEBUG, INFO, WARN, ERROR) |
| `QR_PATH` | `qrcode.png` | Lokasi file QR code PNG |
| `AUTO_REPLY_PERSONAL` | `true` | Auto reply di chat personal |
| `AUTO_REPLY_GROUP` | `false` | Auto reply di grup (TIDAK DIREKOMENDASIKAN) |

### Contoh Environment Variables:
```bash
export AUTO_REPLY_GROUP=false
export LOG_LEVEL=DEBUG
go run cmd/main.go
```

## 💬 Cara Bot Bekerja

### Chat Personal
- ✅ **Auto-reply aktif**: Bot akan membalas semua pesan
- ✅ **Commands**: Bot merespon command seperti `/help`, `/ping`, dll
- ✅ **Friendly**: Response lebih personal dan ramah

### Chat Grup
- ❌ **Auto-reply TIDAK aktif** (untuk menghindari spam)
- ✅ **Commands**: Bot tetap merespon command
- ✅ **Mention**: Bot merespon jika di-mention
- ⚠️ **Bisa diubah**: Set `AUTO_REPLY_GROUP=true` jika diperlukan

## 🤖 Commands yang Tersedia

| Command | Deskripsi | Contoh |
|---------|-----------|--------|
| `/start` | Memulai bot | `/start` |
| `/help` | Menampilkan bantuan | `/help` |
| `/ping` | Test koneksi bot | `/ping` |
| `/info` | Informasi tentang bot | `/info` |
| `/status` | Status bot saat ini | `/status` |
| `/promote` | Promote member grup (coming soon) | `/promote @user` |

## 📚 Penjelasan Kode

### 1. Entry Point (`cmd/main.go`)
```go
// STEP 1: Load konfigurasi
cfg := config.NewConfig()

// STEP 2: Setup logger  
logger := utils.NewLogger("BOT", true)

// STEP 3: Setup QR code generator
qrGen := utils.NewQRCodeGenerator(cfg.QRCodePath)

// ... dan seterusnya
```

**Penjelasan**: File main.go adalah entry point yang mengatur semua komponen bot secara berurutan.

### 2. Konfigurasi (`config/config.go`)
```go
type Config struct {
    DatabasePath      string  // Lokasi database session
    LogLevel         string  // Level logging
    QRCodePath       string  // Lokasi QR code file
    AutoReplyPersonal bool   // Auto reply personal chat
    AutoReplyGroup   bool    // Auto reply grup chat
}
```

**Penjelasan**: Struktur konfigurasi yang memudahkan pengaturan bot tanpa mengubah kode.

### 3. Message Handler (`handlers/message.go`)
```go
func (h *MessageHandler) HandleMessage(evt *events.Message) {
    // STEP 1: Skip pesan dari diri sendiri
    if evt.Info.IsFromMe {
        return
    }
    
    // STEP 2: Ambil teks dari pesan
    messageText := h.getMessageText(evt.Message)
    
    // STEP 3: Identifikasi jenis chat
    isGroup := evt.Info.Chat.Server == types.GroupServer
    
    // STEP 4: Proses berdasarkan jenis chat
    if isGroup {
        h.handleGroupMessage(evt, messageText)
    } else {
        h.handlePersonalMessage(evt, messageText)
    }
}
```

**Penjelasan**: Handler ini membedakan antara chat personal dan grup, lalu memproses sesuai aturan yang berbeda.

### 4. Event Handler (`handlers/events.go`)
```go
func (h *EventHandler) HandleEvent(evt interface{}) {
    switch v := evt.(type) {
    case *events.Message:
        h.messageHandler.HandleMessage(v)
    case *events.Connected:
        fmt.Println("✅ Terhubung ke WhatsApp")
    case *events.Disconnected:
        fmt.Println("❌ Terputus dari WhatsApp")
    // ... event lainnya
    }
}
```

**Penjelasan**: Handler ini menangani semua event dari WhatsApp, tidak hanya pesan tapi juga koneksi, disconnection, dll.

### 5. QR Code Generator (`utils/qrcode.go`)
```go
func (q *QRCodeGenerator) GenerateAndDisplay(code string) error {
    // STEP 1: Generate QR code object
    qr, err := qrcode.New(code, qrcode.Medium)
    
    // STEP 2: Tampilkan sebagai ASCII art
    fmt.Println(qr.ToSmallString(false))
    
    // STEP 3: Simpan sebagai file PNG
    err = qr.WriteFile(256, q.filePath)
    
    return nil
}
```

**Penjelasan**: Utility ini mengubah string QR code menjadi visual yang bisa di-scan langsung di terminal.

## 🔧 Cara Menambah Fitur

### 1. Menambah Command Baru
Edit file `handlers/message.go`, tambahkan case baru di fungsi `handleCommand`:

```go
case lowerText == "/mycommand":
    response = "Response untuk command baru!"
```

### 2. Menambah Auto-Reply Custom
Edit bagian default di fungsi `sendAutoReply`:

```go
responses := []string{
    "Response custom 1",
    "Response custom 2", 
    "Response custom 3",
}
```

### 3. Menambah Event Handler
Edit file `handlers/events.go`, tambahkan case baru di fungsi `HandleEvent`:

```go
case *events.NewEventType:
    h.handleNewEvent(v)
```

## 🛡️ Keamanan dan Best Practices

### 1. Rate Limiting
- Bot sudah mengatur auto-reply grup = false untuk menghindari spam
- Jangan mengubah setting ini kecuali benar-benar diperlukan

### 2. Error Handling
- Semua fungsi penting sudah dilengkapi error handling
- Log error akan muncul dengan format yang jelas

### 3. Session Management
- Session tersimpan otomatis di database SQLite
- Hapus file `session.db` jika ingin logout dan login ulang

### 4. Logging
- Gunakan logger yang sudah disediakan untuk konsistensi
- Level logging bisa diatur sesuai kebutuhan (DEBUG untuk development)

## 🔍 Debugging

### 1. Enable Debug Logging
```bash
export LOG_LEVEL=DEBUG
go run cmd/main.go
```

### 2. Cek File Session
```bash
ls -la session.db  # Cek apakah file session ada
```

### 3. Cek QR Code File
```bash
ls -la qrcode.png  # Cek apakah QR code tersimpan
```

## 🚀 Deployment

### 1. Build Binary
```bash
cd cmd
go build -o whatsapp-bot main.go
./whatsapp-bot
```

### 2. Dengan Docker (Opsional)
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bot cmd/main.go
CMD ["./bot"]
```

## 📞 Support

Jika ada pertanyaan atau masalah:
1. Cek dokumentasi di setiap file kode
2. Lihat log error untuk debugging
3. Pastikan dependencies sudah terinstall dengan benar

## 🎯 Roadmap

- [ ] Implementasi fitur promote grup
- [ ] Integrasi dengan template di folder `layout/`
- [ ] Fitur kirim media (gambar, video, dokumen)
- [ ] Database untuk menyimpan data user
- [ ] Web dashboard untuk monitoring
- [ ] Plugin system untuk extend functionality

---

**Happy Coding!** 🚀