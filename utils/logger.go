// Package utils - File logger.go
// File ini berisi utility functions untuk logging yang konsisten
package utils

import (
	"fmt"
	"time"
)

// Logger adalah struktur sederhana untuk logging dengan format konsisten
type Logger struct {
	// prefix adalah prefix yang akan ditambahkan di setiap log
	prefix string
	
	// showTimestamp menentukan apakah timestamp ditampilkan
	showTimestamp bool
}

// NewLogger membuat logger baru
// Parameter:
// - prefix: prefix untuk setiap log (contoh: "BOT", "CLIENT", dll)
// - showTimestamp: true jika ingin menampilkan timestamp
func NewLogger(prefix string, showTimestamp bool) *Logger {
	return &Logger{
		prefix:        prefix,
		showTimestamp: showTimestamp,
	}
}

// Info menampilkan log level INFO (informasi umum)
func (l *Logger) Info(message string) {
	l.log("INFO", "ℹ️", message)
}

// Success menampilkan log level SUCCESS (operasi berhasil)
func (l *Logger) Success(message string) {
	l.log("SUCCESS", "✅", message)
}

// Warning menampilkan log level WARNING (peringatan)
func (l *Logger) Warning(message string) {
	l.log("WARNING", "⚠️", message)
}

// Error menampilkan log level ERROR (error/kesalahan)
func (l *Logger) Error(message string) {
	l.log("ERROR", "❌", message)
}

// Debug menampilkan log level DEBUG (untuk debugging)
func (l *Logger) Debug(message string) {
	l.log("DEBUG", "🔍", message)
}

// Event menampilkan log untuk event khusus
func (l *Logger) Event(message string) {
	l.log("EVENT", "📨", message)
}

// log adalah fungsi internal untuk format dan tampilkan log
func (l *Logger) log(level, emoji, message string) {
	var output string
	
	// Tambahkan timestamp jika diaktifkan
	if l.showTimestamp {
		timestamp := time.Now().Format("15:04:05")
		output = fmt.Sprintf("[%s]", timestamp)
	}
	
	// Tambahkan prefix jika ada
	if l.prefix != "" {
		if output != "" {
			output += " "
		}
		output += fmt.Sprintf("[%s]", l.prefix)
	}
	
	// Tambahkan level dan emoji
	if output != "" {
		output += " "
	}
	output += fmt.Sprintf("%s %s: %s", emoji, level, message)
	
	// Print ke console
	fmt.Println(output)
}

// Infof menampilkan log INFO dengan format string
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Successf menampilkan log SUCCESS dengan format string
func (l *Logger) Successf(format string, args ...interface{}) {
	l.Success(fmt.Sprintf(format, args...))
}

// Warningf menampilkan log WARNING dengan format string
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warning(fmt.Sprintf(format, args...))
}

// Errorf menampilkan log ERROR dengan format string
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Debugf menampilkan log DEBUG dengan format string
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Eventf menampilkan log EVENT dengan format string
func (l *Logger) Eventf(format string, args ...interface{}) {
	l.Event(fmt.Sprintf(format, args...))
}