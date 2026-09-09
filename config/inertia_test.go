package config

import "testing"

func TestInertiaSSRListenAllowsUnspecifiedBind(t *testing.T) {
	t.Setenv("INERTIA_SSR_LISTEN", "http://0.0.0.0:13714")
	t.Setenv("INERTIA_SSR_URL", "http://ssr-worker:13714")

	cfg, err := NewInertia()
	if err != nil {
		t.Fatalf("NewInertia: %v", err)
	}
	if cfg.SSRBindHost() != "0.0.0.0" {
		t.Fatalf("SSRBindHost = %q, want 0.0.0.0", cfg.SSRBindHost())
	}
	if cfg.SSRBindPort() != "13714" {
		t.Fatalf("SSRBindPort = %q, want 13714", cfg.SSRBindPort())
	}
	if got := cfg.SSRHealthConfig().URL; got != "http://127.0.0.1:13714" {
		t.Fatalf("SSRHealthConfig.URL = %q, want http://127.0.0.1:13714", got)
	}
	if cfg.SSRClientConfig().URL != "http://ssr-worker:13714" {
		t.Fatalf("SSRClientConfig.URL = %q, want http://ssr-worker:13714", cfg.SSRClientConfig().URL)
	}
}

func TestInertiaSSRListenRejectsServiceHostname(t *testing.T) {
	t.Setenv("INERTIA_SSR_LISTEN", "http://ssr-worker:13714")

	if _, err := NewInertia(); err == nil {
		t.Fatal("NewInertia() error = nil, want listen hostname rejection")
	}
}
