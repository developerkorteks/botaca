package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	
	"github.com/nabilulilalbab/promote/config"
	"github.com/nabilulilalbab/promote/handlers"
	"github.com/nabilulilalbab/promote/utils"
	
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// WhatsApp Bot dengan struktur yang rapi dan mudah dipelajari
// File ini adalah entry point utama aplikasi
func main() {
	// STEP 1: Load konfigurasi
	// Konfigurasi berisi semua pengaturan bot seperti database path, auto reply, dll
	cfg := config.NewConfig()
	
	// STEP 2: Setup logger
	// Logger untuk menampilkan informasi dengan format yang rapi
	logger := utils.NewLogger("BOT", true)
	logger.Info("Memulai WhatsApp Bot...")
	
	// STEP 3: Setup QR code generator
	// QR code generator untuk menampilkan QR code visual di terminal
	qrGen := utils.NewQRCodeGenerator(cfg.QRCodePath)
	
	// STEP 4: Setup database untuk session WhatsApp
	// Database SQLite untuk menyimpan session agar tidak perlu login berulang
	logger.Info("Menginisialisasi database session...")
	dbLog := waLog.Noop
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:"+cfg.DatabasePath+"?_foreign_keys=on", dbLog)
	if err != nil {
		logger.Errorf("Gagal membuat database: %v", err)
		os.Exit(1)
	}
	
	// STEP 5: Ambil device store dari database
	// Device store berisi informasi device WhatsApp yang tersimpan
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		logger.Errorf("Gagal mendapatkan device store: %v", err)
		os.Exit(1)
	}
	
	// STEP 6: Buat WhatsApp client
	// Client adalah objek utama untuk berinteraksi dengan WhatsApp
	logger.Info("Membuat WhatsApp client...")
	clientLog := waLog.Stdout("Client", cfg.LogLevel, true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	
	// STEP 7: Setup handlers untuk menangani pesan dan event
	// Message handler menangani pesan masuk
	messageHandler := handlers.NewMessageHandler(client, cfg.AutoReplyPersonal, cfg.AutoReplyGroup)
	
	// Event handler menangani semua event WhatsApp (koneksi, pesan, dll)
	eventHandler := handlers.NewEventHandler(client, messageHandler)
	
	// STEP 8: Daftarkan event handler ke client
	client.AddEventHandler(eventHandler.HandleEvent)
	
	// STEP 9: Connect ke WhatsApp
	if client.Store.ID == nil {
		// Belum login, perlu scan QR code
		logger.Warning("Belum login, memerlukan QR code...")
		err = connectWithQR(client, qrGen, logger)
		if err != nil {
			logger.Errorf("Gagal connect dengan QR: %v", err)
			os.Exit(1)
		}
	} else {
		// Sudah login sebelumnya, langsung connect
		logger.Info("Sudah login sebelumnya, connecting...")
		err = client.Connect()
		if err != nil {
			logger.Errorf("Gagal connect: %v", err)
			os.Exit(1)
		}
	}
	
	// STEP 10: Bot siap digunakan
	logger.Success("Bot berhasil terhubung ke WhatsApp!")
	logger.Info("Bot siap menerima pesan...")
	logger.Info("Tekan Ctrl+C untuk menghentikan bot")
	
	// STEP 11: Wait for interrupt signal (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	// STEP 12: Graceful shutdown
	logger.Info("Menghentikan bot...")
	client.Disconnect()
	logger.Success("Bot berhasil dihentikan. Sampai jumpa!")
}

// connectWithQR menangani proses koneksi dengan QR code
// Fungsi ini akan menampilkan QR code dan menunggu user untuk scan
func connectWithQR(client *whatsmeow.Client, qrGen *utils.QRCodeGenerator, logger *utils.Logger) error {
	// Dapatkan channel untuk menerima QR code dari WhatsApp
	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		return err
	}
	
	// Mulai proses koneksi
	err = client.Connect()
	if err != nil {
		return err
	}
	
	// Loop untuk menangani event QR code
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			// QR code baru diterima, tampilkan ke user
			logger.Info("QR code diterima, menampilkan...")
			err = qrGen.GenerateAndDisplay(evt.Code)
			if err != nil {
				logger.Errorf("Gagal menampilkan QR code: %v", err)
				// Tetap lanjut, tampilkan QR code sebagai text
				logger.Infof("QR Code (text): %s", evt.Code)
			}
			
		case "success":
			// Login berhasil
			logger.Success("QR code berhasil di-scan! Login berhasil.")
			return nil
			
		case "timeout":
			// QR code timeout, akan generate yang baru
			logger.Warning("QR code timeout, generating QR code baru...")
			
		case "error":
			// Error dalam proses login
			logger.Error("Error dalam proses login QR code")
			return fmt.Errorf("QR code login error")
			
		default:
			// Event lain
			logger.Debugf("QR code event: %s", evt.Event)
		}
	}
	
	return nil
}