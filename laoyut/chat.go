package layout

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MessageTemplate interface {
	Render(data map[string]string) string
}

type SimpleTemplate struct {
	Format string
}

func (t SimpleTemplate) Render(data map[string]string) string {
	msg := t.Format
	for k, v := range data {
		placeholder := fmt.Sprintf("{{%s}}", k)
		msg = strings.ReplaceAll(msg, placeholder, v)
	}
	return msg
}

type TemplateSection struct {
	Category  string
	Templates []MessageTemplate
}

var sections = []TemplateSection{
	{
		Category: "vpn",
		Templates: []MessageTemplate{
			SimpleTemplate{Format: "🔥 *VPN PREMIUM {{vpn_type}}* {{price}}/bulan 🚀\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n🌐 *PROTOCOLS:*\n• Trojan GRPC/WS • VMess GRPC/WS\n• VLess GRPC/WS • SSH WebSocket\n• Multipath • Wildcard\n\n🌍 *SERVERS:*\n🇮🇩 ID: wa.me/6287786388052 \n🇸🇬 SG: t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
			SimpleTemplate{Format: "⚡ *VPN {{vpn_type}} PREMIUM* {{price}} aja! 🔥\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n🚀 *FEATURES:*\n• ⚡ High Speed • 🔒 Military Encryption\n• 🌍 Multi Server • 📱 All Device\n• 🛡️ No Log • 🔄 24/7 Reconnect\n\n📱 *ORDER:*\n🇮🇩 wa.me/6287786388052\n🇸🇬 t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
		},
	},
	{
		Category: "ssh",
		Templates: []MessageTemplate{
			SimpleTemplate{Format: "⚡ *SSH WS PREMIUM* {{price}} stabil & kenceng! 🚀\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n🔒 *SSH FEATURES:*\n• WebSocket Support • Bypass DPI\n• High Speed • Stable Connection\n• Multi Port • SSL/TLS Encryption\n\n📱 *ORDER:*\n🇮🇩 wa.me/6287786388052\n🇸🇬 t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
			SimpleTemplate{Format: "🔥 *SSH MURAH* {{price}} aja, cobain sekarang! ⚡\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n🌐 *ADVANTAGES:*\n• WebSocket Protocol • Anti Blokir\n• Speed Unlimited • Server Stabil\n• Support All Device • 24/7 Online\n\n💬 *CONTACT:*\n🇮🇩 wa.me/6287786388052\n🇸🇬 t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
		},
	},
	{
		Category: "kuota",
		Templates: []MessageTemplate{
			SimpleTemplate{Format: "💡 *KUOTA DOR XL {{size}}* {{price}}! 🔥\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n📱 *PAKET DATA:*\n• Kuota {{size}} • Harga {{price}}\n• Proses Cepat • Garansi Masuk\n• Support 24/7\n\n📞 *ORDER:*\n📱 wa.me/6287786388052\n🤖 t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
			SimpleTemplate{Format: "🚀 *INTERNET HEMAT* XL {{size}} {{price}} 👌\n\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\n✅ *BENEFITS:*\n• Harga Terjangkau • Kuota Besar\n• Proses Otomatis • Respon Cepat\n• Terpercaya\n\n💬 *CONTACT:*\n📱 wa.me/6287786388052\n🤖 t.me/grnstoreofficial_bot\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬"},
		},
	},
}

func GetRandomMessage(category string, data map[string]string) (string, error) {
	rand.Seed(time.Now().UnixNano())

	for _, section := range sections {
		if section.Category == category {
			tmpl := section.Templates[rand.Intn(len(section.Templates))]
			return tmpl.Render(data), nil
		}
	}

	return "", fmt.Errorf("category %s not found", category)
}
