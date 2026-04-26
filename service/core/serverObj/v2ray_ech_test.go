package serverObj

import (
	"net/url"
	"strings"
	"testing"
)

func TestParseVlessURLReadsEchConfigList(t *testing.T) {
	link := "vless://b29d5bf8-2a27-4f92-a6f4-53d1558e6ee1@125.138.95.125:10000?type=ws&security=tls&sni=ss.batch1.workers.dev&fp=chrome&echConfigList=cloudflare-ech.com%2Bhttps%3A%2F%2Fdns.alidns.com%2Fdns-query&host=ss.batch1.workers.dev&path=%2F%3Fproxyip%3D125.138.95.125%3A10000#test"
	v, err := ParseVlessURL(link)
	if err != nil {
		t.Fatalf("ParseVlessURL error: %v", err)
	}
	if v.EchConfigList != "cloudflare-ech.com+https://dns.alidns.com/dns-query" {
		t.Fatalf("unexpected echConfigList: %q", v.EchConfigList)
	}
}

func TestParseVlessURLReadsLegacyEchAliasAndNumericAllowInsecure(t *testing.T) {
	link := "vless://b29d5bf8-2a27-4f92-a6f4-53d1558e6ee1@118.45.236.190:16922?encryption=none&security=tls&sni=ss.batch1.workers.dev&fp=chrome&insecure=1&allowInsecure=0&ech=cloudflare-ech.com%2Bhttps%3A%2F%2Fdns.alidns.com%2Fdns-query&type=ws&host=ss.batch1.workers.dev&path=%2F%3Fip%3D118.45.236.190%3A16922#ech-kr-118.45.236.190"
	v, err := ParseVlessURL(link)
	if err != nil {
		t.Fatalf("ParseVlessURL error: %v", err)
	}
	if v.EchConfigList != "cloudflare-ech.com+https://dns.alidns.com/dns-query" {
		t.Fatalf("unexpected legacy ech alias value: %q", v.EchConfigList)
	}
	if !v.AllowInsecure {
		t.Fatal("expected numeric insecure=1 to set AllowInsecure")
	}
}

func TestVlessConfigurationIncludesEchConfigList(t *testing.T) {
	v := &V2Ray{
		Ps:            "test",
		Add:           "125.138.95.125",
		Port:          "10000",
		ID:            "b29d5bf8-2a27-4f92-a6f4-53d1558e6ee1",
		Net:           "ws",
		Host:          "ss.batch1.workers.dev",
		SNI:           "ss.batch1.workers.dev",
		Path:          "/?proxyip=125.138.95.125:10000",
		TLS:           "tls",
		Fingerprint:   "chrome",
		EchConfigList: "cloudflare-ech.com+https://dns.alidns.com/dns-query",
		Protocol:      "vless",
	}
	cfg, err := v.Configuration(PriorInfo{Tag: "proxy"})
	if err != nil {
		t.Fatalf("Configuration error: %v", err)
	}
	if cfg.CoreOutbound.StreamSettings == nil || cfg.CoreOutbound.StreamSettings.TLSSettings == nil {
		t.Fatalf("missing tls settings: %#v", cfg.CoreOutbound.StreamSettings)
	}
	if cfg.CoreOutbound.StreamSettings.TLSSettings.EchConfigList != v.EchConfigList {
		t.Fatalf("unexpected echConfigList in tlsSettings: %#v", cfg.CoreOutbound.StreamSettings.TLSSettings)
	}
}

func TestVlessExportToURLIncludesEchConfigList(t *testing.T) {
	v := &V2Ray{
		Ps:            "test",
		Add:           "125.138.95.125",
		Port:          "10000",
		ID:            "b29d5bf8-2a27-4f92-a6f4-53d1558e6ee1",
		Net:           "ws",
		Host:          "ss.batch1.workers.dev",
		SNI:           "ss.batch1.workers.dev",
		Path:          "/?proxyip=125.138.95.125:10000",
		TLS:           "tls",
		Fingerprint:   "chrome",
		EchConfigList: "cloudflare-ech.com+https://dns.alidns.com/dns-query",
		Protocol:      "vless",
	}
	link := v.ExportToURL()
	if !strings.Contains(link, "echConfigList=") {
		t.Fatalf("exported link missing echConfigList: %s", link)
	}
	u, err := url.Parse(link)
	if err != nil {
		t.Fatalf("url parse error: %v", err)
	}
	if got := u.Query().Get("echConfigList"); got != v.EchConfigList {
		t.Fatalf("unexpected echConfigList in exported link: %q", got)
	}
}
