// Package utils - File qrcode.go
// File ini berisi utility functions untuk menangani QR code
// Termasuk generate QR code visual di terminal dan save ke file
package utils

import (
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
)

// QRCodeGenerator adalah struktur untuk generate QR code
type QRCodeGenerator struct {
	// filePath adalah lokasi file PNG QR code akan disimpan
	filePath string
}

// NewQRCodeGenerator membuat generator QR code baru
// Parameter:
// - filePath: lokasi file PNG untuk menyimpan QR code (contoh: "qrcode.png")
func NewQRCodeGenerator(filePath string) *QRCodeGenerator {
	return &QRCodeGenerator{
		filePath: filePath,
	}
}

// GenerateAndDisplay membuat QR code dari text dan menampilkannya di terminal
// Juga menyimpan QR code sebagai file PNG untuk backup
// Parameter:
// - code: string QR code yang diterima dari WhatsApp
func (q *QRCodeGenerator) GenerateAndDisplay(code string) error {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📱 SCAN QR CODE INI DENGAN WHATSAPP ANDA")
	fmt.Println(strings.Repeat("=", 60))
	
	// STEP 1: Generate QR code object
	// qrcode.Medium adalah level error correction (Low, Medium, High, Highest)
	// Medium cukup untuk kebanyakan kasus dan tidak terlalu besar
	qr, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		fmt.Printf("❌ Gagal generate QR code: %v\n", err)
		fmt.Printf("QR Code (text): %s\n", code)
		return err
	}
	
	// STEP 2: Tampilkan QR code sebagai ASCII art di terminal
	// ToSmallString(false) menghasilkan QR code dengan karakter ASCII
	// Parameter false berarti tidak invert warna (hitam = █, putih = spasi)
	asciiQR := qr.ToSmallString(false)
	fmt.Println(asciiQR)
	
	// STEP 3: Simpan QR code sebagai file PNG untuk backup
	// Size 256x256 pixel cukup untuk scan dengan mudah
	err = qr.WriteFile(256, q.filePath)
	if err != nil {
		fmt.Printf("⚠️ Gagal menyimpan QR code ke file: %v\n", err)
	} else {
		fmt.Printf("💾 QR code juga disimpan sebagai '%s'\n", q.filePath)
	}
	
	// STEP 4: Tampilkan instruksi cara scan
	q.displayInstructions()
	
	return nil
}

// displayInstructions menampilkan instruksi cara scan QR code
func (q *QRCodeGenerator) displayInstructions() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📲 CARA SCAN QR CODE:")
	fmt.Println("")
	fmt.Println("1. 📱 Buka WhatsApp di HP Android/iPhone")
	fmt.Println("2. ⚙️  Pergi ke Settings (Pengaturan)")
	fmt.Println("3. 🔗 Pilih 'Linked Devices' (Perangkat Tertaut)")
	fmt.Println("4. ➕ Tap 'Link a Device' (Tautkan Perangkat)")
	fmt.Println("5. 📷 Scan QR code di atas dengan kamera HP")
	fmt.Println("")
	fmt.Println("💡 TIPS:")
	fmt.Printf("   • Jika QR di terminal tidak jelas, buka file '%s'\n", q.filePath)
	fmt.Println("   • Pastikan layar cukup terang untuk scan")
	fmt.Println("   • Jarak ideal: 15-30 cm dari layar")
	fmt.Println("")
	fmt.Println(strings.Repeat("=", 60))
}

// GenerateQRFile hanya menyimpan QR code ke file tanpa menampilkan di terminal
// Berguna jika Anda hanya ingin file PNG saja
func (q *QRCodeGenerator) GenerateQRFile(code string) error {
	qr, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("gagal generate QR code: %w", err)
	}
	
	err = qr.WriteFile(256, q.filePath)
	if err != nil {
		return fmt.Errorf("gagal menyimpan QR code: %w", err)
	}
	
	fmt.Printf("💾 QR code disimpan sebagai '%s'\n", q.filePath)
	return nil
}

// GetQRAsString mengembalikan QR code sebagai string ASCII
// Berguna jika Anda ingin menampilkan QR code di tempat lain
func (q *QRCodeGenerator) GetQRAsString(code string) (string, error) {
	qr, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("gagal generate QR code: %w", err)
	}
	
	return qr.ToSmallString(false), nil
}