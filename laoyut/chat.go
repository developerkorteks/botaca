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
			SimpleTemplate{Format: "🔥 VPN {{vpn_type}} cuma {{price}}/bulan 🚀"},
			SimpleTemplate{Format: "Mau internet stabil? Coba {{vpn_type}} VPN, mulai {{price}} aja!"},
		},
	},
	{
		Category: "ssh",
		Templates: []MessageTemplate{
			SimpleTemplate{Format: "⚡ SSH WS Premium cuma {{price}}, stabil & kenceng!"},
			SimpleTemplate{Format: "SSH murah mulai {{price}} aja, cobain sekarang 🔥"},
		},
	},
	{
		Category: "kuota",
		Templates: []MessageTemplate{
			SimpleTemplate{Format: "💡 Kuota Dor XL {{size}} cuma {{price}}!"},
			SimpleTemplate{Format: "Internet hemat: XL {{size}} cuma {{price}} 👌"},
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
