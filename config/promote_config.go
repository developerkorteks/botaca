// Package config - Konfigurasi khusus untuk fitur auto promote
package config

import (
	"fmt"
	"os"
	"strings"
)

// PromoteConfig berisi konfigurasi untuk fitur auto promote
type PromoteConfig struct {
	// DatabasePath untuk database auto promote (terpisah dari session)
	PromoteDatabasePath string

	// AdminNumbers adalah daftar nomor WhatsApp admin yang bisa mengelola template
	AdminNumbers []string

	// AutoPromoteInterval dalam jam (default: 4 jam)
	AutoPromoteInterval int

	// MaxTemplatesPerCategory maksimal template per kategori
	MaxTemplatesPerCategory int

	// EnableAutoPromote mengaktifkan/nonaktifkan fitur auto promote
	EnableAutoPromote bool

	// LogAutoPromote mengaktifkan logging detail untuk auto promote
	LogAutoPromote bool
}

// NewPromoteConfig membuat konfigurasi default untuk auto promote
func NewPromoteConfig() *PromoteConfig {
	return &PromoteConfig{
		// Database terpisah untuk auto promote
		PromoteDatabasePath: getEnvOrDefault("PROMOTE_DB_PATH", "data/promote.db"),

		// Admin numbers dari environment variable (pisahkan dengan koma)
		AdminNumbers: getAdminNumbers(),

		// Interval default 4 jam
		AutoPromoteInterval: getEnvIntOrDefault("AUTO_PROMOTE_INTERVAL", 4),

		// Maksimal 20 template per kategori
		MaxTemplatesPerCategory: getEnvIntOrDefault("MAX_TEMPLATES_PER_CATEGORY", 20),

		// Auto promote diaktifkan secara default
		EnableAutoPromote: getEnvBoolOrDefault("ENABLE_AUTO_PROMOTE", true),

		// Logging detail diaktifkan
		LogAutoPromote: getEnvBoolOrDefault("LOG_AUTO_PROMOTE", true),
	}
}

// getAdminNumbers mengambil daftar nomor admin dari environment variable
func getAdminNumbers() []string {
	adminEnv := os.Getenv("ADMIN_NUMBERS")
	if adminEnv == "" {
		// Default admin numbers (ganti dengan nomor Anda)
		return []string{
			"6285117557905", // Ganti dengan nomor admin utama
			"6285150588080", // Tambahkan nomor admin lain jika perlu
		}
	}

	// Split berdasarkan koma dan bersihkan spasi
	numbers := strings.Split(adminEnv, ",")
	var cleanNumbers []string

	for _, number := range numbers {
		clean := strings.TrimSpace(number)
		if clean != "" {
			cleanNumbers = append(cleanNumbers, clean)
		}
	}

	return cleanNumbers
}

// getEnvIntOrDefault mengambil nilai integer dari environment variable
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple parsing (bisa diperbaiki dengan strconv.Atoi)
		switch value {
		case "1":
			return 1
		case "2":
			return 2
		case "3":
			return 3
		case "4":
			return 4
		case "6":
			return 6
		case "8":
			return 8
		case "12":
			return 12
		case "24":
			return 24
		default:
			return defaultValue
		}
	}
	return defaultValue
}

// IsAdmin mengecek apakah nomor adalah admin
func (c *PromoteConfig) IsAdmin(phoneNumber string) bool {
	for _, admin := range c.AdminNumbers {
		if admin == phoneNumber {
			return true
		}
	}
	return false
}

// AddAdmin menambahkan nomor admin baru
func (c *PromoteConfig) AddAdmin(phoneNumber string) {
	// Cek apakah sudah ada
	if c.IsAdmin(phoneNumber) {
		return
	}

	c.AdminNumbers = append(c.AdminNumbers, phoneNumber)
}

// RemoveAdmin menghapus nomor admin
func (c *PromoteConfig) RemoveAdmin(phoneNumber string) {
	var newAdmins []string
	for _, admin := range c.AdminNumbers {
		if admin != phoneNumber {
			newAdmins = append(newAdmins, admin)
		}
	}
	c.AdminNumbers = newAdmins
}

// GetAdminList mendapatkan daftar admin dalam format string
func (c *PromoteConfig) GetAdminList() string {
	if len(c.AdminNumbers) == 0 {
		return "Tidak ada admin terdaftar"
	}

	var result strings.Builder
	result.WriteString("👑 **Admin Auto Promote:**\n")

	for i, admin := range c.AdminNumbers {
		result.WriteString(fmt.Sprintf("%d. +%s\n", i+1, admin))
	}

	return result.String()
}

// ValidateConfig memvalidasi konfigurasi
func (c *PromoteConfig) ValidateConfig() []string {
	var errors []string

	if c.PromoteDatabasePath == "" {
		errors = append(errors, "Database path tidak boleh kosong")
	}

	if len(c.AdminNumbers) == 0 {
		errors = append(errors, "Minimal harus ada 1 admin")
	}

	if c.AutoPromoteInterval < 1 || c.AutoPromoteInterval > 24 {
		errors = append(errors, "Interval auto promote harus antara 1-24 jam")
	}

	if c.MaxTemplatesPerCategory < 1 {
		errors = append(errors, "Maksimal template per kategori minimal 1")
	}

	return errors
}

// GetConfigInfo mendapatkan informasi konfigurasi dalam format string
func (c *PromoteConfig) GetConfigInfo() string {
	return fmt.Sprintf(`⚙️ **KONFIGURASI AUTO PROMOTE**

📁 **Database:** %s
👑 **Admin:** %d orang
⏰ **Interval:** %d jam
📝 **Max Template/Kategori:** %d
🤖 **Status:** %s
📊 **Logging:** %s

%s

💡 **Environment Variables:**
• PROMOTE_DB_PATH - Path database
• ADMIN_NUMBERS - Nomor admin (pisah koma)
• AUTO_PROMOTE_INTERVAL - Interval jam
• ENABLE_AUTO_PROMOTE - true/false
• LOG_AUTO_PROMOTE - true/false`,
		c.PromoteDatabasePath,
		len(c.AdminNumbers),
		c.AutoPromoteInterval,
		c.MaxTemplatesPerCategory,
		getBoolText(c.EnableAutoPromote),
		getBoolText(c.LogAutoPromote),
		c.GetAdminList())
}

// getBoolText mengkonversi boolean ke teks
func getBoolText(value bool) string {
	if value {
		return "Aktif ✅"
	}
	return "Tidak Aktif ❌"
}

// UpdateConfig memperbarui konfigurasi dari environment variables
func (c *PromoteConfig) UpdateConfig() {
	c.PromoteDatabasePath = getEnvOrDefault("PROMOTE_DB_PATH", c.PromoteDatabasePath)
	c.AdminNumbers = getAdminNumbers()
	c.AutoPromoteInterval = getEnvIntOrDefault("AUTO_PROMOTE_INTERVAL", c.AutoPromoteInterval)
	c.MaxTemplatesPerCategory = getEnvIntOrDefault("MAX_TEMPLATES_PER_CATEGORY", c.MaxTemplatesPerCategory)
	c.EnableAutoPromote = getEnvBoolOrDefault("ENABLE_AUTO_PROMOTE", c.EnableAutoPromote)
	c.LogAutoPromote = getEnvBoolOrDefault("LOG_AUTO_PROMOTE", c.LogAutoPromote)
}
