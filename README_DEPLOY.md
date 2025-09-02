# 🚀 Auto Promote Bot - Docker Deployment Guide

## 📋 Prerequisites

1. **VPS/Server** dengan minimal:
   - 1 GB RAM
   - 1 CPU Core
   - 10 GB Storage
   - Ubuntu 20.04+ / CentOS 7+ / Debian 10+

2. **Docker & Docker Compose** terinstall

## 🛠️ Installation Steps

### 1. Install Docker (jika belum ada)

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Logout and login again
```

### 2. Clone Repository

```bash
git clone https://github.com/yourusername/auto-promote-bot.git
cd auto-promote-bot
```

### 3. Deploy Bot

```bash
# Make scripts executable
chmod +x deploy.sh update.sh

# Deploy the bot
./deploy.sh
```

### 4. Scan QR Code

```bash
# View logs to see QR code
docker-compose logs -f whatsapp-bot
```

## 🔧 Management Commands

### View Logs
```bash
docker-compose logs -f
```

### Stop Bot
```bash
docker-compose down
```

### Restart Bot
```bash
docker-compose restart
```

### Update Bot
```bash
./update.sh
```

### Check Status
```bash
docker-compose ps
```

### Access Container Shell
```bash
docker-compose exec whatsapp-bot sh
```

## 📁 File Structure

```
auto-promote-bot/
├── Dockerfile              # Docker build configuration
├── docker-compose.yml      # Docker Compose configuration
├── deploy.sh               # Deployment script
├── update.sh               # Update script
├── data/                   # Persistent data directory
├── session.db              # WhatsApp session (auto-created)
├── promote.db              # Auto promote database (auto-created)
└── qrcode.png              # QR code file (auto-created)
```

## 🔄 Auto-Restart Configuration

Bot akan otomatis restart jika crash berkat `restart: unless-stopped` di docker-compose.yml.

## 📊 Monitoring

### Health Check
```bash
docker-compose ps
```

### Resource Usage
```bash
docker stats auto-promote-bot
```

### Disk Usage
```bash
du -sh data/
```

## 🛡️ Security

1. **Firewall**: Hanya buka port yang diperlukan
2. **Updates**: Selalu update sistem dan Docker
3. **Backup**: Backup file database secara berkala

## 🔧 Troubleshooting

### Bot tidak start
```bash
# Check logs
docker-compose logs

# Check container status
docker-compose ps

# Restart container
docker-compose restart
```

### QR Code tidak muncul
```bash
# Check logs
docker-compose logs -f whatsapp-bot

# Restart bot
docker-compose restart
```

### Database error
```bash
# Stop bot
docker-compose down

# Remove database files
rm -f session.db promote.db

# Start bot again
docker-compose up -d
```

### Update failed
```bash
# Force rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 📱 Usage

1. **Aktivasi Auto Promote**: `.aca` di grup
2. **Nonaktifkan**: `.disableaca` di grup
3. **Status**: `.statuspromo`
4. **Help**: `.help`

## 🎯 Admin Commands

- `.fetchproducts` - Ambil produk dari API
- `.addtemplate` - Tambah template
- `.deleteall` - Hapus semua template
- `.activegroups` - Lihat grup aktif

## 📞 Support

Jika ada masalah:
1. Check logs: `docker-compose logs`
2. Restart bot: `docker-compose restart`
3. Update bot: `./update.sh`

---

**Happy Promoting!** 🚀💰