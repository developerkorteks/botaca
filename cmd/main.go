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
	"github.com/nabilulilalbab/promote/database"
	"github.com/nabilulilalbab/promote/handlers"
	"github.com/nabilulilalbab/promote/services"
	"github.com/nabilulilalbab/promote/utils"
	
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// WhatsApp Bot dengan struktur yang rapi dan mudah dipelajari
// File ini adalah entry point utama aplikasi
func main() {
	// STEP 1: Load konfigurasi
	// Konfigurasi berisi semua pengaturan bot seperti database path, auto reply, dll
	cfg := config.NewConfig()
	promoteCfg := config.NewPromoteConfig()
	
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
	
	// STEP 7: Setup Auto Promote System (jika diaktifkan)
	var autoPromoteService *services.AutoPromoteService
	var templateService *services.TemplateService
	var promoteCommandHandler *handlers.PromoteCommandHandler
	var adminCommandHandler *handlers.AdminCommandHandler
	
	if promoteCfg.EnableAutoPromote {
		logger.Info("Initializing Auto Promote System...")
		
		// Setup database untuk auto promote
		promoteDB, promoteRepo, err := database.InitializeDatabase(promoteCfg.PromoteDatabasePath)
		if err != nil {
			logger.Errorf("Failed to initialize promote database: %v", err)
			os.Exit(1)
		}
		defer promoteDB.Close()
		
		// Setup services
		templateService = services.NewTemplateService(promoteRepo, logger)
		autoPromoteService = services.NewAutoPromoteService(client, promoteRepo, logger)
		apiProductService := services.NewAPIProductService(templateService, logger)
		
		// Setup command handlers
		promoteCommandHandler = handlers.NewPromoteCommandHandler(autoPromoteService, templateService, logger)
		adminCommandHandler = handlers.NewAdminCommandHandler(autoPromoteService, templateService, apiProductService, logger, promoteCfg.AdminNumbers)
		
		logger.Success("Auto Promote System initialized!")
	}
	
	// STEP 8: Setup handlers untuk menangani pesan dan event
	// Message handler menangani pesan masuk
	messageHandler := handlers.NewMessageHandler(client, cfg.AutoReplyPersonal, cfg.AutoReplyGroup)
	
	// Set auto promote handlers jika tersedia
	if promoteCommandHandler != nil && adminCommandHandler != nil {
		messageHandler.SetAutoPromoteHandlers(promoteCommandHandler, adminCommandHandler)
		logger.Info("Auto Promote handlers attached to message handler")
	}
	
	// Event handler menangani semua event WhatsApp (koneksi, pesan, dll)
	eventHandler := handlers.NewEventHandler(client, messageHandler)
	
	// STEP 9: Daftarkan event handler ke client
	client.AddEventHandler(eventHandler.HandleEvent)
	
	// STEP 10: Connect ke WhatsApp
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
	
	// STEP 11: Start Auto Promote Scheduler (jika diaktifkan)
	if autoPromoteService != nil {
		logger.Info("Starting Auto Promote Scheduler...")
		autoPromoteService.StartScheduler()
		
		// Log konfigurasi auto promote
		logger.Infof("Auto Promote Config: %d admin(s), %d hour interval", 
			len(promoteCfg.AdminNumbers), promoteCfg.AutoPromoteInterval)
	}
	
	// STEP 12: Bot siap digunakan
	logger.Success("Bot berhasil terhubung ke WhatsApp!")
	logger.Info("Bot siap menerima pesan...")
	
	if promoteCfg.EnableAutoPromote {
		logger.Success("🚀 Auto Promote System is READY!")
		logger.Info("Commands: .promote, .disablepromote, .promotehelp")
	}
	
	logger.Info("Tekan Ctrl+C untuk menghentikan bot")
	
	// STEP 13: Wait for interrupt signal (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	// STEP 14: Graceful shutdown
	logger.Info("Menghentikan bot...")
	
	// Stop auto promote scheduler jika berjalan
	if autoPromoteService != nil {
		logger.Info("Stopping Auto Promote Scheduler...")
		autoPromoteService.StopScheduler()
	}
	
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