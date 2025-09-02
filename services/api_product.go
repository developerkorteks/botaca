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
		return "❌ Tidak ada produk yang ditemukan dari API", nil
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
			errors = append(errors, fmt.Sprintf("Gagal membuat template group %d: %v", i/groupSize+1, err))
			continue
		}
		
		createdCount++
		s.logger.Infof("Created template group %d with %d products", i/groupSize+1, len(productGroup))
	}

	// Buat response
	var result strings.Builder
	result.WriteString("🛒 *UPDATE PRODUK DARI API*\n\n")
	result.WriteString(fmt.Sprintf("✅ **Berhasil:** %d template group dibuat\n", createdCount))
	result.WriteString(fmt.Sprintf("📦 **Total Produk:** %d (digroup per %d)\n", len(products), groupSize))
	
	if len(errors) > 0 {
		result.WriteString(fmt.Sprintf("❌ **Gagal:** %d group\n", len(errors)))
	}
	
	if createdCount > 0 {
		result.WriteString("\n💡 **Info:**\n")
		result.WriteString("• Template produk sudah digroup dan ditambahkan\n")
		result.WriteString("• Setiap template berisi 15 produk\n")
		result.WriteString("• Auto promote akan random pilih group template")
	}

	return result.String(), nil
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
	
	template.WriteString(fmt.Sprintf(`🛒 *PAKET DATA GROUP %d*

🔥 *PROMO TERBATAS!* 
Stok menipis, buruan order!

━━━━━━━━━━━━━━━━━━━━━━

`, groupNum))

	// Tambahkan daftar produk (hanya yang ada, tidak paksa 15)
	for _, product := range products {
		// Validasi data produk
		if product.PackageNameShort == "" || product.PackageHarga == "" {
			continue // Skip produk dengan data kosong
		}
		
		template.WriteString(fmt.Sprintf("📱 **%s**\n", product.PackageNameShort))
		template.WriteString(fmt.Sprintf("💰 %s (Harga / Jasa DOR, baca deskripsi bot dor)\n\n", product.PackageHarga))
	}

	// Tambahkan teknik FOMO dan promosi
	template.WriteString(`━━━━━━━━━━━━━━━━━━━━━━

⚡ *KENAPA PILIH KAMI?*

🤖 *Malas lihat Telegram yang ribet?*
   WhatsApp aja cukup! Simple & user-friendly

📱 *Pengen privasi lebih?*
   Telegram kami siap melayani dengan fitur lengkap!

━━━━━━━━━━━━━━━━━━━━━━

🛡️ *JAMINAN KEPERCAYAAN:*
✅ Paket RESMI = GARANSI PENUH 
⚠️ Paket DOR = TANPA GARANSI
💰 Harga tertera = Harga / Jasa DOR
📖 Baca deskripsi bot untuk detail paket DOR
💯 Transparansi total untuk kepercayaan Anda!

━━━━━━━━━━━━━━━━━━━━━━

🌐 **VPN INJECT TERSEDIA:**

📱 *Android VPN:*
🇸🇬 Server SG: Rp 8.000/bulan
   • Max 2 IP • 900GB Bandwidth

📺 *STB VPN:*  
🇸🇬 Server SG: Rp 8.000/bulan
   • Max 1 IP • 900GB Bandwidth
🇮🇩 Server Indo: Rp 15.000/bulan
   • Max 1 IP • 900GB Bandwidth

🖥️ *PC/Laptop VPN:*
🇮🇩 Server Indo: Rp 10.000/bulan
   • Max 3 IP • 900GB Bandwidth

━━━━━━━━━━━━━━━━━━━━━━

🎯 *ORDER SEKARANG:*

📱 *WhatsApp:*
   wa.me/6287786388052

🤖 *Telegram:*
   https://t.me/grnstoreofficial_bot

⏰ *JAM OPERASIONAL:*
   🟢 BUKA: 01:00 - 23:00 WIB
   🔴 TUTUP: 23:00 - 01:00 WIB
   📞 Respon cepat di jam operasional

━━━━━━━━━━━━━━━━━━━━━━

⏰ *BURUAN!* Stok terbatas!
🔥 *FOMO ALERT:* Yang ragu pasti nyesal!

#PaketData #VPN #GRNStore #OrderSekarang`)

	return template.String()
}

// generateProductTemplate membuat template promosi untuk produk individual (backup)
func (s *APIProductService) generateProductTemplate(product Product) string {
	// Potong deskripsi jika terlalu panjang
	description := product.PackageDescription
	if len(description) > 200 {
		description = description[:200] + "..."
	}
	
	template := fmt.Sprintf(`📱 *%s*

💰 **Harga:** %s
📝 **Detail:** %s

🎯 **Order Sekarang:**
📱 WhatsApp: wa.me/6287786388052
🤖 Telegram: https://t.me/grnstoreofficial_bot

⚡ *Stok terbatas, buruan order!*
🔥 *Jangan sampai nyesal kemudian!*

#PaketData #GRNStore #OrderSekarang #%s`,
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
	result.WriteString(fmt.Sprintf("📦 **Total Paket:** %d\n", len(products)))
	result.WriteString(fmt.Sprintf("⏰ **Dengan Daily Limit:** %d\n", dailyLimitCount))
	result.WriteString(fmt.Sprintf("🔓 **Tanpa Login:** %d\n", noLoginCount))
	result.WriteString(fmt.Sprintf("🔐 **Perlu Login:** %d\n\n", len(products)-noLoginCount))
	
	result.WriteString("💡 **Info:**\n")
	result.WriteString("• Semua paket dari API GRN Store\n")
	result.WriteString("• Data diambil real-time dari server\n")
	result.WriteString("• Gunakan .fetchproducts untuk update template")

	return result.String(), nil
}