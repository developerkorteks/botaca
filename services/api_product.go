// Package services - API Product service untuk mengambil produk dari API
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nabilulilalbab/promote/utils"
)

// APIProductService mengelola pengambilan produk dari API
type APIProductService struct {
	templateService *TemplateService
	logger          *utils.Logger
	apiBaseURL      string
}

// ProductResponse struktur response dari API sesuai dokumentasi
type ProductResponse struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Success    bool      `json:"success"`
	Data       []Product `json:"data"`
}

// Product struktur produk sesuai API dokumentasi
type Product struct {
	PackageCode        string `json:"package_code"`
	PackageName        string `json:"package_name"`
	PackageNameShort   string `json:"package_name_alias_short"`
	PackageDescription string `json:"package_description"`
	PackageHargaInt    int    `json:"package_harga_int"`
	PackageHarga       string `json:"package_harga"`
	HaveDailyLimit     bool   `json:"have_daily_limit"`
	NoNeedLogin        bool   `json:"no_need_login"`
}

// NewAPIProductService membuat service baru
func NewAPIProductService(templateService *TemplateService, logger *utils.Logger) *APIProductService {
	return &APIProductService{
		templateService: templateService,
		logger:          logger,
		apiBaseURL:      "https://grn-store.vercel.app/api", // URL API sesuai dokumentasi
	}
}

// FetchProductsAndCreateTemplates mengambil produk dari API dan membuat template
func (s *APIProductService) FetchProductsAndCreateTemplates() (string, error) {
	s.logger.Info("Fetching products from API...")

	// Ambil data dari API
	products, err := s.fetchProductsFromAPI()
	if err != nil {
		s.logger.Errorf("Failed to fetch products: %v", err)
		return "", fmt.Errorf("gagal mengambil data produk: %v", err)
	}

	if len(products) == 0 {
		return `ℹ️ *TIDAK ADA PRODUK*

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	          *API KOSONG*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🚫 Tidak ada produk yang ditemukan dari API saat ini.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 *KEMUNGKINAN PENYEBAB*
• Server API sedang dalam maintenance.
• Belum ada produk yang ditambahkan di API.
• Terjadi kesalahan filter pada API.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔄 *Coba lagi nanti atau hubungi admin API*`, nil
	}

	// Group produk per 15 dan buat template gabungan
	createdCount := 0
	var errors []string
	groupSize := 15

	for i := 0; i < len(products); i += groupSize {
		end := i + groupSize
		if end > len(products) {
			end = len(products)
		}

		productGroup := products[i:end]
		templateContent := s.generateGroupedProductTemplate(productGroup, i/groupSize+1)
		templateTitle := fmt.Sprintf("Paket Group %d (%d Produk)", i/groupSize+1, len(productGroup))

		_, err := s.templateService.CreateTemplate(templateTitle, templateContent, "produk_api_group")
		if err != nil {
			s.logger.Errorf("Failed to create template group %d: %v", i/groupSize+1, err)
			errors = append(errors, fmt.Sprintf("Group %d: %v", i/groupSize+1, err))
			continue
		}

		createdCount++
		s.logger.Infof("Created template group %d with %d products", i/groupSize+1, len(productGroup))
	}

	// Buat response
	var result strings.Builder
	result.WriteString("🛒 *UPDATE PRODUK DARI API*\n\n")

	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	result.WriteString("          *HASIL FETCH API*\n")
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	result.WriteString("📊 *STATISTIK IMPORT*\n")
	result.WriteString(fmt.Sprintf("✅ *Berhasil:* %d template group\n", createdCount))
	result.WriteString(fmt.Sprintf("📦 *Total Produk:* %d\n", len(products)))
	result.WriteString(fmt.Sprintf("📋 *Per Group:* %d produk\n", groupSize))

	if len(errors) > 0 {
		result.WriteString(fmt.Sprintf("❌ *Gagal:* %d group\n", len(errors)))
		result.WriteString("\n🔍 *Detail Error:*\n")
		for i, errMsg := range errors {
			if i < 3 { // Tampilkan maksimal 3 error pertama
				result.WriteString(fmt.Sprintf("• %s\n", errMsg))
			}
		}
		if len(errors) > 3 {
			result.WriteString(fmt.Sprintf("• ... dan %d error lainnya\n", len(errors)-3))
		}
	}

	if createdCount > 0 {
		result.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		result.WriteString("📋 *INFORMASI SISTEM*\n")
		result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		result.WriteString("• Template produk sudah digroup\n")
		result.WriteString("• Setiap template berisi 15 produk\n")
		result.WriteString("• Auto promote pilih random group\n")
		result.WriteString("• Format WhatsApp sudah optimized\n")

		result.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		result.WriteString("🎮 *COMMANDS SELANJUTNYA*\n")
		result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		result.WriteString("• *.listtemplates*\n")
		result.WriteString("  _Lihat template yang dibuat_\n\n")
		result.WriteString("• *.templatestats*\n")
		result.WriteString("  _Statistik semua template_\n\n")
		result.WriteString("• *.testgroup [ID]*\n")
		result.WriteString("  _Test kirim ke grup_")
	}

	return result.String(), nil
}

// Helper function untuk max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// fetchProductsFromAPI mengambil data produk dari API
func (s *APIProductService) fetchProductsFromAPI() ([]Product, error) {
	// Buat HTTP client dengan timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Buat request sesuai dokumentasi API
	req, err := http.NewRequest("GET", "https://grnstore.domcloud.dev/api/user/products?limit=200", nil)
	if err != nil {
		return nil, err
	}

	// Set headers sesuai dokumentasi
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", "nadia-admin-2024-secure-key")
	req.Header.Set("User-Agent", "WhatsApp-Bot/1.0")

	// Kirim request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var productResp ProductResponse
	err = json.Unmarshal(body, &productResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if !productResp.Success {
		return nil, fmt.Errorf("API error: %s", productResp.Message)
	}

	return productResp.Data, nil
}

// generateGroupedProductTemplate membuat template promosi untuk group produk
func (s *APIProductService) generateGroupedProductTemplate(products []Product, groupNum int) string {
	var template strings.Builder

	template.WriteString(fmt.Sprintf(`🔥 *VPN PREMIUM CATALOG* 🔥

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
🌐 *VPN SERVICES*
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🚀 *PROTOCOLS:*
• Trojan GRPC/WS • VMess GRPC/WS
• VLess GRPC/WS • SSH WebSocket
• Multipath • Wildcard

🌍 *SERVERS:*
🇮🇩 ID: wa.me/6287786388052
🇸🇬 SG: t.me/grnstoreofficial_bot

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🛒 *PAKET DATA GROUP %d*

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
⚡ *PROMO TERBATAS!*
_Stok menipis, buruan order!_
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

📋 *DAFTAR PAKET:*

`, groupNum))

	// Tambahkan daftar produk dengan format yang lebih ringkas
	for i, product := range products {
		// Validasi data produk
		if product.PackageNameShort == "" || product.PackageHarga == "" {
			continue // Skip produk dengan data kosong
		}

		template.WriteString(fmt.Sprintf("📱 *%s* - %s\n", product.PackageNameShort, product.PackageHarga))

		if i < len(products)-1 {
			template.WriteString("\n")
		}
	}

	// Tambahkan informasi singkat dan contact
	template.WriteString(`

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

✅ *RESMI* = GARANSI PENUH 
⚠️ *DOR* = TANPA GARANSI
💰 *Harga* = Harga/Jasa DOR

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
🔥 *VPN PREMIUM FEATURES*
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🌐 *VPN FEATURES:*
• ⚡ High Speed • 🔒 Military Encryption
• 🌍 Multi Server • 📱 All Device
• 🛡️ No Log • 🔄 24/7 Reconnect

🚀 *ADVANCED PROTOCOLS:*
• Trojan-GRPC (Ultra Fast)
• VMess-WS (Stable) 
• VLess-GRPC (Low Latency)
• SSH-WS (Bypass DPI)
• Multipath Custom • Wildcard

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
📞 *ORDER CENTER*
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🇮🇩 *VPN ID:* wa.me/6287786388052
🇸🇬 *VPN SG:* t.me/grnstoreofficial_bot

🛒 *PAKET DATA:*
📱 wa.me/6287786388052
🤖 t.me/grnstoreofficial_bot

👨‍💼 *ADMIN:*
📱 wa.me/6287786388052
📱 wa.me/6285117557905

👥 *GROUP:* chat.whatsapp.com/IeIXOndIoFr0apnlKzghUC

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🟢 *BUKA:* 01:00 - 23:00 WIB
⏰ *BURUAN ORDER!* Stok terbatas!

#PaketData #VPNPremium #GRNStore`)

	return template.String()
}

// generateProductTemplate membuat template promosi untuk produk individual (backup)
func (s *APIProductService) generateProductTemplate(product Product) string {
	// Potong deskripsi jika terlalu panjang
	description := product.PackageDescription
	if len(description) > 200 {
		description = description[:200] + "..."
	}

	template := fmt.Sprintf(`🔥 *VPN PREMIUM CATALOG* 🔥

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
🌐 *VPN SERVICES*
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🚀 *PROTOCOLS:*
• Trojan GRPC/WS • VMess GRPC/WS
• VLess GRPC/WS • SSH WebSocket
• Multipath • Wildcard

🌍 *SERVERS:*
🇮🇩 ID: wa.me/6287786388052
🇸🇬 SG: t.me/grnstoreofficial_bot

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

📱 *%s*

💰 *Harga:* %s
📝 *Detail:* %s

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
📞 *ORDER CENTER*
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

🇮🇩 *VPN ID:* wa.me/6287786388052
🇸🇬 *VPN SG:* t.me/grnstoreofficial_bot

🛒 *PAKET DATA:*
📱 wa.me/6287786388052
🤖 t.me/grnstoreofficial_bot

👨‍💼 *ADMIN:*
📱 wa.me/6287786388052
📱 wa.me/6285117557905

👥 *GROUP:* chat.whatsapp.com/IeIXOndIoFr0apnlKzghUC

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬

⚡ *Stok terbatas, buruan order!*
🔥 *Jangan sampai nyesal kemudian!*

#PaketData #VPNPremium #GRNStore #%s`,
		product.PackageNameShort,
		product.PackageHarga,
		description,
		product.PackageCode)

	return template
}

// formatPrice memformat harga ke format Rupiah
func (s *APIProductService) formatPrice(price float64) string {
	if price < 1000 {
		return fmt.Sprintf("Rp %.0f", price)
	} else if price < 1000000 {
		return fmt.Sprintf("Rp %.0fK", price/1000)
	} else {
		return fmt.Sprintf("Rp %.1fJT", price/1000000)
	}
}

// UpdateAPIBaseURL mengupdate URL API
func (s *APIProductService) UpdateAPIBaseURL(newURL string) {
	s.apiBaseURL = newURL
	s.logger.Infof("API Base URL updated to: %s", newURL)
}

// GetProductStats mendapatkan statistik produk dari API
func (s *APIProductService) GetProductStats() (string, error) {
	products, err := s.fetchProductsFromAPI()
	if err != nil {
		return "", err
	}

	dailyLimitCount := 0
	noLoginCount := 0

	for _, product := range products {
		if product.HaveDailyLimit {
			dailyLimitCount++
		}
		if product.NoNeedLogin {
			noLoginCount++
		}
	}

	var result strings.Builder
	result.WriteString("📊 *STATISTIK PRODUK API*\n\n")
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	result.WriteString("           *RINGKASAN PRODUK*\n")
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	result.WriteString(fmt.Sprintf("📦 *Total Paket Tersedia:* %d\n", len(products)))
	result.WriteString(fmt.Sprintf("⏰ *Paket dengan Limit Harian:* %d\n", dailyLimitCount))
	result.WriteString(fmt.Sprintf("🔓 *Paket Tanpa Login:* %d\n", noLoginCount))
	result.WriteString(fmt.Sprintf("🔐 *Paket Perlu Login:* %d\n", len(products)-noLoginCount))

	result.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	result.WriteString("           *INFORMASI TAMBAHAN*\n")
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	result.WriteString("• Semua paket bersumber dari API GRN Store.\n")
	result.WriteString("• Data statistik ini diambil secara real-time.\n")
	result.WriteString("• Gunakan *.fetchproducts* untuk memperbarui template.")

	return result.String(), nil
}
