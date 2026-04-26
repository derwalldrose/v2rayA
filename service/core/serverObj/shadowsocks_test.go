package serverObj

import "testing"

func TestPluginMuxFollowsPriorInfo(t *testing.T) {
	if mux := getPluginMux(PriorInfo{MuxEnabled: false, MuxConcurrency: 8}); mux != nil {
		t.Fatalf("expected nil mux when MuxEnabled=false, got %#v", mux)
	}

	mux := getPluginMux(PriorInfo{MuxEnabled: true, MuxConcurrency: 8})
	if mux == nil {
		t.Fatal("expected mux when MuxEnabled=true")
	}
	if !mux.Enabled || mux.Concurrency != 8 {
		t.Fatalf("unexpected mux: %#v", mux)
	}

	mux = getPluginMux(PriorInfo{MuxEnabled: true, MuxConcurrency: 0})
	if mux == nil {
		t.Fatal("expected mux when MuxEnabled=true and MuxConcurrency invalid")
	}
	if mux.Concurrency != 1 {
		t.Fatalf("expected fallback concurrency 1, got %#v", mux)
	}
}

func TestV2rayPluginTransportUsesPriorInfoMux(t *testing.T) {
	s := &Shadowsocks{
		Name:     "edgetunnel",
		Server:   "www.shopify.com",
		Port:     2052,
		Password: "secret",
		Cipher:   "aes-128-gcm",
		Plugin: Sip003{
			Name: "v2ray-plugin",
			Opts: Sip003Opts{
				Host: "ss.batch1.workers.dev",
				Path: "/?enc=aes-128-gcm&proxyip=test:50001",
				Impl: "transport",
			},
		},
	}

	cfg, err := s.Configuration(PriorInfo{Tag: "proxy", MuxEnabled: true, MuxConcurrency: 16})
	if err != nil {
		t.Fatalf("Configuration returned error: %v", err)
	}
	if cfg.CoreOutbound.Mux == nil {
		t.Fatal("expected mux to be present")
	}
	if cfg.CoreOutbound.Mux.Concurrency != 16 {
		t.Fatalf("expected concurrency 16, got %#v", cfg.CoreOutbound.Mux)
	}

	cfg, err = s.Configuration(PriorInfo{Tag: "proxy", MuxEnabled: false, MuxConcurrency: 16})
	if err != nil {
		t.Fatalf("Configuration returned error: %v", err)
	}
	if cfg.CoreOutbound.Mux != nil {
		t.Fatalf("expected mux to be omitted when disabled, got %#v", cfg.CoreOutbound.Mux)
	}
}

func TestParseSip003OptsUnescapesEscapedEqualsInPath(t *testing.T) {
	opts := ParseSip003Opts("mode=websocket;host=ss.batch1.workers.dev;path=/?enc\\=aes-128-gcm&proxyip\\=kr.090401.xyz:50001;mux=0")
	if opts.Path != "/?enc=aes-128-gcm&proxyip=kr.090401.xyz:50001" {
		t.Fatalf("unexpected parsed path: %q", opts.Path)
	}
}
