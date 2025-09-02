# 📋 Planning & Improvement - Auto Promote Feature (Business/Sales)

## 🎯 Tujuan Fitur Auto Promote

Mengembangkan sistem auto promote untuk **JUALAN/BISNIS** yang dapat:
1. Diaktifkan dengan command `.promote` di grup manapun (termasuk grup orang lain)
2. Dihentikan dengan command `.disablepromote`
3. Mengirim pesan promosi jualan setiap 4 jam sekali
4. Menggunakan template promosi bisnis yang dinamis
5. Random selection dari 5-10 template promosi jualan
6. CRUD operations untuk template promosi bisnis

## 🏗️ Struktur Implementasi

### 1. Database Schema
```sql
-- Tabel untuk menyimpan status auto promote per grup
CREATE TABLE auto_promote_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_jid TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    started_at DATETIME,
    last_promote_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel untuk menyimpan template promote
CREATE TABLE promote_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Komponen Utama

#### A. AutoPromoteManager
- Mengelola status auto promote per grup
- Scheduler untuk mengirim pesan setiap 4 jam
- Integration dengan database

#### B. PromoteTemplateManager
- CRUD operations untuk template
- Random selection template
- Template validation

#### C. Command Handlers
- `.promote` - Aktivasi auto promote
- `.disablepromote` - Deaktivasi auto promote
- `.addtemplate` - Tambah template (admin only)
- `.listtemplates` - List semua template
- `.edittemplate` - Edit template (admin only)
- `.deletetemplate` - Hapus template (admin only)

### 3. File Structure
```
promote/
├── database/
│   ├── models.go          # Database models
│   ├── migrations.go      # Database migrations
│   └── repository.go      # Database operations
├── services/
│   ├── auto_promote.go    # Auto promote service
│   ├── template.go        # Template management
│   └── scheduler.go       # Scheduling service
├── handlers/
│   ├── promote_commands.go # Command handlers
│   └── admin_commands.go   # Admin-only commands
└── utils/
    └── permissions.go      # Permission checking
```

## 🔄 Flow Implementasi

### Phase 1: Database & Models
1. ✅ Setup database schema
2. ✅ Create models dan repository
3. ✅ Database migrations

### Phase 2: Core Services
1. ✅ AutoPromoteManager implementation
2. ✅ PromoteTemplateManager implementation
3. ✅ Scheduler service

### Phase 3: Command Handlers
1. ✅ Basic commands (.promote, .disablepromote)
2. ✅ Admin commands (template management)
3. ✅ Permission validation

### Phase 4: Integration
1. ✅ Integrate dengan message handler
2. ✅ Testing dan debugging
3. ✅ Documentation

## 🛡️ Security & Permissions

### Admin Commands
- Hanya admin grup yang bisa menggunakan command promote
- Bot harus menjadi admin untuk menjalankan auto promote
- Validation untuk semua input

### Rate Limiting
- Auto promote hanya setiap 4 jam
- Cooldown untuk command spam prevention

## 📝 Template System

### Default Templates Promosi Bisnis (5-10 template)
1. **Produk Unggulan** - Promosi produk terlaris
2. **Diskon & Promo** - Penawaran khusus dan diskon
3. **Testimoni Customer** - Review dan feedback pelanggan
4. **Flash Sale** - Promosi terbatas waktu
5. **Produk Baru** - Launching produk terbaru
6. **Bundle Package** - Paket hemat dan bundling
7. **Free Ongkir** - Promosi gratis ongkos kirim
8. **Cashback & Reward** - Program loyalitas customer
9. **Limited Stock** - Stok terbatas, beli sekarang
10. **Contact Info** - Informasi kontak dan cara order

### Template Variables
- `{GROUP_NAME}` - Nama grup
- `{MEMBER_COUNT}` - Jumlah member
- `{DATE}` - Tanggal saat ini
- `{TIME}` - Waktu saat ini
- `{ADMIN}` - Nama admin yang mengaktifkan

## 🔧 Configuration

### Environment Variables
```env
AUTO_PROMOTE_INTERVAL=4h
MAX_TEMPLATES=20
ADMIN_ONLY_TEMPLATE_MANAGEMENT=true
```

### Default Settings
- Interval: 4 jam
- Max templates per grup: 20
- Auto cleanup inactive groups: 30 hari

## 📊 Monitoring & Logging

### Metrics
- Total active groups
- Messages sent per day
- Template usage statistics
- Error rates

### Logging
- Auto promote activations/deactivations
- Template CRUD operations
- Scheduler events
- Error tracking

## 🚀 Future Enhancements

1. **Advanced Scheduling**
   - Custom intervals per grup
   - Time-based scheduling (jam tertentu)
   - Day-specific scheduling

2. **Analytics**
   - Message engagement tracking
   - Template performance metrics
   - Group activity correlation

3. **Advanced Templates**
   - Rich media support (images, videos)
   - Interactive buttons
   - Conditional content

4. **Multi-language Support**
   - Template localization
   - Language detection
   - Auto-translation

## ✅ Success Criteria

1. ✅ Auto promote dapat diaktifkan/dinonaktifkan per grup
2. ✅ Pesan terkirim setiap 4 jam setelah aktivasi
3. ✅ Random selection dari template yang tersedia
4. ✅ Admin dapat mengelola template (CRUD)
5. ✅ Sistem stabil dan tidak spam
6. ✅ Error handling yang baik
7. ✅ Logging dan monitoring yang memadai

## 🐛 Known Issues & Solutions

### Issue 1: Memory Usage
- **Problem**: Scheduler bisa consume memory jika banyak grup
- **Solution**: Efficient goroutine management, cleanup inactive groups

### Issue 2: Database Locks
- **Problem**: Concurrent access ke database
- **Solution**: Proper transaction handling, connection pooling

### Issue 3: WhatsApp Rate Limits
- **Problem**: Terlalu banyak pesan bisa kena rate limit
- **Solution**: Implement backoff strategy, message queuing

## 📚 References

- [WhatsApp Business API Limits](https://developers.facebook.com/docs/whatsapp/api/rate-limits)
- [Go Scheduler Patterns](https://pkg.go.dev/time#Ticker)
- [SQLite Best Practices](https://www.sqlite.org/lang.html)