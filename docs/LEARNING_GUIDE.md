# 📚 Panduan Belajar WhatsApp Bot dengan Whatsmeow

Panduan ini akan membantu Anda memahami cara kerja library `whatsmeow` step by step, dari konsep dasar hingga implementasi advanced.

## 🎯 Tujuan Pembelajaran

Setelah mengikuti panduan ini, Anda akan memahami:
1. ✅ Konsep dasar WhatsApp Web API
2. ✅ Cara kerja library whatsmeow
3. ✅ Struktur project yang baik
4. ✅ Cara menangani pesan personal vs grup
5. ✅ Cara membuat command system
6. ✅ Best practices untuk bot WhatsApp

## 📖 Konsep Dasar

### 1. Apa itu WhatsApp Web API?
WhatsApp Web API adalah cara untuk berinteraksi dengan WhatsApp melalui protokol yang sama dengan WhatsApp Web. Library `whatsmeow` mengimplementasikan protokol ini dalam bahasa Go.

### 2. Komponen Utama Whatsmeow

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WhatsApp      │    │   whatsmeow     │    │   Your Bot      │
│   Servers       │◄──►│   Library       │◄──►│   Application   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**Komponen:**
- **Client**: Objek utama untuk koneksi ke WhatsApp
- **Store**: Database untuk menyimpan session dan data
- **Event Handler**: Fungsi yang menangani event dari WhatsApp
- **Message**: Struktur data untuk pesan WhatsApp

## 🔧 Struktur Project Explained

### 1. Mengapa Struktur Ini?

```
cmd/main.go        ← Entry point (seperti main function di C++)
config/config.go   ← Pengaturan (seperti config file)
handlers/          ← Business logic (seperti controller di MVC)
utils/             ← Helper functions (seperti library utility)
```

**Keuntungan struktur ini:**
- ✅ **Separation of Concerns**: Setiap file punya tanggung jawab spesifik
- ✅ **Maintainable**: Mudah diubah dan diperbaiki
- ✅ **Testable**: Mudah untuk unit testing
- ✅ **Scalable**: Mudah ditambah fitur baru

### 2. Flow Aplikasi

```
1. main.go          → Load config, setup components
2. config.go        → Provide configuration
3. qrcode.go        → Generate QR for login
4. events.go        → Handle all WhatsApp events
5. message.go       → Handle incoming messages specifically
6. logger.go        → Log everything for debugging
```

## 💬 Cara Menangani Pesan

### 1. Personal Chat vs Group Chat

**Personal Chat:**
```go
// Di personal chat, bot lebih bebas merespon
if !isGroup {
    // Auto-reply aktif
    // Semua pesan direspon
    // Response lebih personal
}
```

**Group Chat:**
```go
// Di grup, bot harus lebih hati-hati
if isGroup {
    // Hanya respon command atau mention
    // Hindari spam
    // Response lebih formal
}
```

### 2. Mengapa Berbeda?

**Alasan teknis:**
- Grup bisa punya ratusan member
- Auto-reply di grup = spam
- WhatsApp bisa ban bot yang spam

**Alasan praktis:**
- User tidak suka bot yang terlalu aktif di grup
- Admin grup bisa kick bot yang mengganggu

### 3. Implementasi di Kode

```go
func (h *MessageHandler) handleGroupMessage(evt *events.Message, messageText string) {
    // Cek apakah bot di-mention
    isMentioned := h.isBotMentioned(evt.Message)
    
    // Cek apakah ini command
    isCommand := strings.HasPrefix(messageText, "/")
    
    if isCommand {
        // Selalu proses command
        h.handleCommand(evt, messageText)
    } else if isMentioned {
        // Respon jika di-mention
        h.sendAutoReply(evt.Info.Chat, messageText, true)
    } else if h.autoReplyGroup {
        // Hanya jika setting diaktifkan (tidak direkomendasikan)
        h.sendAutoReply(evt.Info.Chat, messageText, true)
    }
    // Jika tidak ada kondisi di atas, bot TIDAK merespon
}
```

## 🤖 Command System

### 1. Konsep Command

Command adalah pesan khusus yang dimulai dengan `/` (seperti `/help`, `/ping`). Bot akan selalu merespon command, baik di personal maupun grup.

### 2. Implementasi Command

```go
func (h *MessageHandler) handleCommand(evt *events.Message, messageText string) {
    // Ubah ke lowercase untuk case-insensitive
    lowerText := strings.ToLower(strings.TrimSpace(messageText))
    
    var response string
    
    switch {
    case lowerText == "/help":
        response = h.getHelpMessage()
    case lowerText == "/ping":
        response = "🏓 Pong!"
    // ... command lainnya
    }
    
    h.sendMessage(evt.Info.Chat, response)
}
```

### 3. Menambah Command Baru

**Step 1:** Tambahkan case baru
```go
case lowerText == "/weather":
    response = "🌤️ Cuaca hari ini cerah!"
```

**Step 2:** Update help message
```go
func (h *MessageHandler) getHelpMessage() string {
    return `Commands:
• /help - Bantuan
• /ping - Test koneksi  
• /weather - Cek cuaca  ← Tambahkan ini
`
}
```

## 🔍 Event Handling

### 1. Jenis-jenis Event

WhatsApp mengirim berbagai event ke bot:

```go
switch v := evt.(type) {
case *events.Message:        // Pesan masuk
case *events.Connected:      // Bot terhubung
case *events.Disconnected:   // Bot terputus
case *events.LoggedOut:      // Bot di-logout
case *events.Receipt:        // Pesan terkirim/dibaca
case *events.Presence:       // Online/offline/typing
case *events.GroupInfo:      // Info grup berubah
// ... dan banyak lagi
}
```

### 2. Event yang Penting

**Message Event:**
```go
case *events.Message:
    // Ini yang paling penting - pesan masuk
    h.messageHandler.HandleMessage(v)
```

**Connection Events:**
```go
case *events.Connected:
    // Bot berhasil terhubung
    fmt.Println("✅ Terhubung ke WhatsApp")

case *events.Disconnected:
    // Bot terputus (akan auto-reconnect)
    fmt.Println("❌ Terputus dari WhatsApp")
```

### 3. Event yang Bisa Diabaikan

```go
case *events.Receipt:
    // Receipt biasanya tidak perlu ditangani
    // Kecuali Anda ingin tracking delivery status

case *events.Presence:
    // Presence (online/offline) biasanya tidak perlu
    // Kecuali untuk fitur khusus
```

## 🔐 Session Management

### 1. Apa itu Session?

Session adalah data login yang tersimpan agar bot tidak perlu scan QR code setiap kali restart.

### 2. Cara Kerja Session

```go
// Cek apakah sudah login sebelumnya
if client.Store.ID == nil {
    // Belum login - perlu QR code
    connectWithQR(client, qrGen, logger)
} else {
    // Sudah login - langsung connect
    client.Connect()
}
```

### 3. File Session

```
session.db  ← File SQLite yang menyimpan session
```

**Isi file session:**
- Device ID
- Encryption keys
- Contact list
- Group info
- dll

### 4. Troubleshooting Session

**Problem:** Bot tidak bisa login
**Solution:** Hapus file session dan scan QR lagi
```bash
rm session.db
go run cmd/main.go
```

## 🎨 QR Code Visual

### 1. Mengapa QR Code Visual?

QR code dari WhatsApp berupa string panjang seperti:
```
2@Qxfir9lkP+53PgGgPKZXQE4VA2fVKm4psu69dPh4LMllT431n2wUAuY00XlKzrLz+/37890Y6+FnFR9D7+QliGGQyoAsd9hqZ+Y=,Fb/94uMM4vO08ZfIplxaXu6hvLbJkpdrTAbl4pGbuEM=...
```

String ini tidak bisa di-scan. Kita perlu mengubahnya menjadi QR code visual.

### 2. Implementasi QR Code

```go
func (q *QRCodeGenerator) GenerateAndDisplay(code string) error {
    // STEP 1: Generate QR code object
    qr, err := qrcode.New(code, qrcode.Medium)
    
    // STEP 2: Convert ke ASCII art
    asciiQR := qr.ToSmallString(false)
    fmt.Println(asciiQR)
    
    // STEP 3: Save sebagai PNG file
    err = qr.WriteFile(256, q.filePath)
    
    return nil
}
```

### 3. Hasil QR Code

```
█████████████████████████████████
████ ▄▄▄▄▄ ██  ▄▀ ▀▀▀▀ ▄ ▄▀▀██▀▀ 
████ █   █ █ █   █▀▀▄ ▀▄▄▄▄█████▄
████ █▄▄▄█ █ ▀██▄ ▄ █▀██ ▀█▀▀  ▀ 
█████████████████████████████████
```

## 🚀 Best Practices

### 1. Error Handling

**Selalu handle error:**
```go
_, err := client.SendMessage(context.Background(), chatJID, msg)
if err != nil {
    logger.Errorf("Gagal mengirim pesan: %v", err)
    return
}
```

### 2. Logging

**Gunakan logger yang konsisten:**
```go
logger.Info("Bot starting...")
logger.Success("Bot connected!")
logger.Warning("QR code timeout")
logger.Error("Connection failed")
```

### 3. Rate Limiting

**Jangan spam pesan:**
```go
// JANGAN:
for i := 0; i < 100; i++ {
    client.SendMessage(...)  // Ini akan kena ban!
}

// LAKUKAN:
time.Sleep(1 * time.Second)  // Kasih jeda antar pesan
client.SendMessage(...)
```

### 4. Graceful Shutdown

**Handle Ctrl+C dengan baik:**
```go
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c

logger.Info("Menghentikan bot...")
client.Disconnect()  // Disconnect dengan benar
```

## 🧪 Testing dan Debugging

### 1. Testing Step by Step

**Step 1:** Test QR code generation
```bash
go run cmd/main.go
# Pastikan QR code muncul dan bisa di-scan
```

**Step 2:** Test personal chat
```
Kirim: "halo"
Expect: Bot membalas dengan greeting
```

**Step 3:** Test commands
```
Kirim: "/ping"
Expect: "🏓 Pong!"
```

**Step 4:** Test group behavior
```
Di grup: Kirim "halo" (tanpa mention)
Expect: Bot TIDAK membalas (jika auto-reply grup = false)

Di grup: Kirim "/ping"
Expect: Bot membalas "🏓 Pong!"
```

### 2. Debug Logging

**Enable debug mode:**
```bash
export LOG_LEVEL=DEBUG
go run cmd/main.go
```

**Debug output:**
```
[15:04:05] [BOT] ℹ️ INFO: Memulai WhatsApp Bot...
[15:04:05] [BOT] ℹ️ INFO: Menginisialisasi database session...
[15:04:05] [BOT] ✅ SUCCESS: Bot berhasil terhubung ke WhatsApp!
[15:04:10] [BOT] 📨 EVENT: Pesan masuk [personal]: halo
```

### 3. Common Issues

**Issue 1:** QR code tidak muncul
```bash
# Solution: Cek dependencies
go mod tidy
```

**Issue 2:** Bot tidak merespon
```bash
# Solution: Cek log untuk error
export LOG_LEVEL=DEBUG
go run cmd/main.go
```

**Issue 3:** Session expired
```bash
# Solution: Hapus session dan login ulang
rm session.db
go run cmd/main.go
```

## 🎓 Latihan

### Latihan 1: Tambah Command Baru
Tambahkan command `/time` yang menampilkan waktu saat ini.

**Hint:**
```go
case lowerText == "/time":
    currentTime := time.Now().Format("15:04:05")
    response = "🕐 Waktu sekarang: " + currentTime
```

### Latihan 2: Custom Auto-Reply
Buat auto-reply yang berbeda berdasarkan kata kunci dalam pesan.

**Hint:**
```go
if strings.Contains(lowerText, "terima kasih") {
    response = "🙏 Sama-sama!"
} else if strings.Contains(lowerText, "selamat pagi") {
    response = "🌅 Selamat pagi juga!"
}
```

### Latihan 3: Group Welcome Message
Buat bot mengirim pesan selamat datang ketika ditambahkan ke grup.

**Hint:** Lihat event `*events.JoinedGroup` di `handlers/events.go`

## 🎯 Next Steps

Setelah memahami dasar-dasar ini, Anda bisa:

1. **Implementasi fitur promote** (sesuai kebutuhan Anda)
2. **Integrasi dengan database** untuk menyimpan data user
3. **Fitur kirim media** (gambar, video, dokumen)
4. **Scheduled messages** untuk pesan otomatis
5. **Web dashboard** untuk monitoring bot

---

**Selamat belajar!** 🚀 Jika ada pertanyaan, cek dokumentasi di setiap file kode atau eksperimen langsung dengan bot.