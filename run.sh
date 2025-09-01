#!/bin/bash

# Script untuk menjalankan WhatsApp Bot
# Script ini menyediakan menu interaktif untuk berbagai operasi bot

echo "🤖 WhatsApp Bot Runner"
echo "====================="

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fungsi untuk menampilkan menu
show_menu() {
    echo ""
    echo -e "${BLUE}Pilih operasi yang ingin dilakukan:${NC}"
    echo "1. 🚀 Jalankan Bot (Normal)"
    echo "2. 🔍 Jalankan Bot (Debug Mode)"
    echo "3. 🧹 Clean Session (Logout & Reset)"
    echo "4. 📦 Install Dependencies"
    echo "5. 🔧 Build Binary"
    echo "6. 📊 Cek Status Session"
    echo "7. 📚 Buka Dokumentasi"
    echo "8. ❌ Exit"
    echo ""
}

# Fungsi untuk install dependencies
install_deps() {
    echo -e "${YELLOW}📦 Installing dependencies...${NC}"
    go mod tidy
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Dependencies berhasil diinstall!${NC}"
    else
        echo -e "${RED}❌ Gagal install dependencies!${NC}"
    fi
}

# Fungsi untuk clean session
clean_session() {
    echo -e "${YELLOW}🧹 Cleaning session files...${NC}"
    
    # Hapus file session
    rm -f session.db
    rm -f qrcode.png
    
    # Hapus file log jika ada
    rm -f bot.log
    
    echo -e "${GREEN}✅ Session files berhasil dihapus!${NC}"
    echo -e "${BLUE}💡 Anda perlu scan QR code lagi saat menjalankan bot.${NC}"
}

# Fungsi untuk run bot normal
run_bot() {
    echo -e "${GREEN}🚀 Starting WhatsApp Bot...${NC}"
    echo -e "${BLUE}📱 Jika belum login, QR code akan muncul${NC}"
    echo -e "${BLUE}⚡ Tekan Ctrl+C untuk stop bot${NC}"
    echo ""
    
    cd cmd
    go run main.go
    cd ..
}

# Fungsi untuk run bot debug
run_bot_debug() {
    echo -e "${GREEN}🔍 Starting WhatsApp Bot (Debug Mode)...${NC}"
    echo -e "${YELLOW}⚠️ Debug mode akan menampilkan log detail${NC}"
    echo -e "${BLUE}📱 Jika belum login, QR code akan muncul${NC}"
    echo -e "${BLUE}⚡ Tekan Ctrl+C untuk stop bot${NC}"
    echo ""
    
    cd cmd
    LOG_LEVEL=DEBUG go run main.go
    cd ..
}

# Fungsi untuk build binary
build_bot() {
    echo -e "${YELLOW}🔧 Building WhatsApp Bot binary...${NC}"
    
    cd cmd
    go build -o ../whatsapp-bot main.go
    cd ..
    
    if [ -f "whatsapp-bot" ]; then
        echo -e "${GREEN}✅ Binary berhasil dibuild!${NC}"
        echo -e "${BLUE}💡 Jalankan dengan: ./whatsapp-bot${NC}"
        
        # Buat executable
        chmod +x whatsapp-bot
    else
        echo -e "${RED}❌ Gagal build binary!${NC}"
    fi
}

# Fungsi untuk cek status session
check_session() {
    echo -e "${BLUE}📊 Checking session status...${NC}"
    echo ""
    
    if [ -f "session.db" ]; then
        echo -e "${GREEN}✅ Session file ditemukan${NC}"
        
        # Cek ukuran file
        size=$(stat -f%z session.db 2>/dev/null || stat -c%s session.db 2>/dev/null)
        echo -e "${BLUE}📁 Ukuran file: ${size} bytes${NC}"
        
        # Cek tanggal modifikasi
        if command -v stat >/dev/null 2>&1; then
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                modified=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" session.db)
            else
                # Linux
                modified=$(stat -c "%y" session.db | cut -d'.' -f1)
            fi
            echo -e "${BLUE}📅 Last modified: ${modified}${NC}"
        fi
        
        echo -e "${GREEN}🎉 Bot kemungkinan sudah login sebelumnya${NC}"
    else
        echo -e "${YELLOW}⚠️ Session file tidak ditemukan${NC}"
        echo -e "${BLUE}💡 Bot akan meminta QR code saat dijalankan${NC}"
    fi
    
    if [ -f "qrcode.png" ]; then
        echo -e "${GREEN}✅ QR code backup file ditemukan${NC}"
    fi
    
    echo ""
}

# Fungsi untuk buka dokumentasi
open_docs() {
    echo -e "${BLUE}📚 Membuka dokumentasi...${NC}"
    echo ""
    echo "📖 Dokumentasi yang tersedia:"
    echo "1. README.md - Dokumentasi utama"
    echo "2. docs/LEARNING_GUIDE.md - Panduan belajar"
    echo "3. examples/promote_example.go - Contoh fitur promote"
    echo ""
    
    # Coba buka dengan berbagai editor
    if command -v code >/dev/null 2>&1; then
        echo -e "${GREEN}🚀 Membuka dengan VS Code...${NC}"
        code README.md
    elif command -v nano >/dev/null 2>&1; then
        echo -e "${GREEN}📝 Membuka dengan nano...${NC}"
        nano README.md
    elif command -v vim >/dev/null 2>&1; then
        echo -e "${GREEN}📝 Membuka dengan vim...${NC}"
        vim README.md
    else
        echo -e "${YELLOW}⚠️ Editor tidak ditemukan. Silakan buka file manual:${NC}"
        echo "cat README.md"
    fi
}

# Fungsi untuk validasi environment
check_environment() {
    echo -e "${BLUE}🔍 Checking environment...${NC}"
    
    # Cek Go installation
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}❌ Go tidak terinstall!${NC}"
        echo -e "${YELLOW}💡 Install Go dari: https://golang.org/dl/${NC}"
        exit 1
    fi
    
    # Cek Go version
    go_version=$(go version | cut -d' ' -f3)
    echo -e "${GREEN}✅ Go version: ${go_version}${NC}"
    
    # Cek apakah di dalam project directory
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ Tidak ditemukan go.mod file!${NC}"
        echo -e "${YELLOW}💡 Pastikan Anda berada di directory project yang benar${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Environment OK${NC}"
}

# Fungsi untuk menampilkan info project
show_project_info() {
    echo -e "${BLUE}📋 Project Information:${NC}"
    echo "🤖 WhatsApp Bot dengan Whatsmeow"
    echo "📝 Bahasa: Go (Golang)"
    echo "📚 Library: whatsmeow + go-qrcode"
    echo "🎯 Fitur: Visual QR code, Auto-reply, Commands"
    echo ""
    echo -e "${BLUE}📁 Struktur Project:${NC}"
    echo "├── cmd/main.go          # Entry point"
    echo "├── config/config.go     # Konfigurasi"
    echo "├── handlers/            # Event & message handlers"
    echo "├── utils/               # Utilities (QR, logger)"
    echo "├── docs/                # Dokumentasi"
    echo "├── examples/            # Contoh implementasi"
    echo "└── layout/              # Template promote"
    echo ""
}

# Main script
main() {
    # Cek environment dulu
    check_environment
    
    # Tampilkan info project
    show_project_info
    
    # Main loop
    while true; do
        show_menu
        read -p "Masukkan pilihan (1-8): " choice
        
        case $choice in
            1)
                run_bot
                ;;
            2)
                run_bot_debug
                ;;
            3)
                clean_session
                ;;
            4)
                install_deps
                ;;
            5)
                build_bot
                ;;
            6)
                check_session
                ;;
            7)
                open_docs
                ;;
            8)
                echo -e "${GREEN}👋 Goodbye! Happy coding!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ Pilihan tidak valid! Silakan pilih 1-8.${NC}"
                ;;
        esac
        
        echo ""
        read -p "Tekan Enter untuk kembali ke menu..."
    done
}

# Jalankan main function
main