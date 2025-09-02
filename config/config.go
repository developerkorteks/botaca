// Package config berisi konfigurasi untuk WhatsApp Bot
// File ini mengatur semua pengaturan dasar bot seperti database path, log level, dll
package config

import (
	"os"
)

// Config adalah struktur yang menyimpan semua konfigurasi bot
type Config struct {
	// DatabasePath adalah lokasi file database SQLite untuk menyimpan session WhatsApp
	DatabasePath string
	
	// LogLevel menentukan level logging (DEBUG, INFO, WARN, ERROR)
	LogLevel string
	
	// QRCodePath adalah lokasi file QR code PNG akan disimpan
	QRCodePath string
	
	// AutoReplyPersonal menentukan apakah bot otomatis membalas chat personal
	AutoReplyPersonal bool
	
	// AutoReplyGroup menentukan apakah bot otomatis membalas chat grup
	// PENTING: Set false jika Anda ada di banyak grup untuk menghindari spam
	AutoReplyGroup bool
}

// NewConfig membuat konfigurasi default untuk bot
// Fungsi ini akan dipanggil saat bot pertama kali dijalankan
func NewConfig() *Config {
	return &Config{
		// Database akan disimpan di folder data
		DatabasePath: getEnvOrDefault("DB_PATH", "data/session.db"),
		
		// Log level default adalah INFO (tidak terlalu verbose)
		LogLevel: getEnvOrDefault("LOG_LEVEL", "INFO"),
		
		// QR code akan disimpan di folder data
		QRCodePath: getEnvOrDefault("QR_PATH", "data/qrcode.png"),
		
		// Auto reply untuk chat personal diaktifkan
		AutoReplyPersonal: getEnvBoolOrDefault("AUTO_REPLY_PERSONAL", true),
		
		// Auto reply untuk grup DIMATIKAN untuk menghindari spam
		// Anda bisa mengubah ini ke true jika ingin bot membalas di grup
		AutoReplyGroup: getEnvBoolOrDefault("AUTO_REPLY_GROUP", false),
	}
}

// getEnvOrDefault mengambil nilai dari environment variable atau menggunakan default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolOrDefault mengambil nilai boolean dari environment variable atau menggunakan default
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}