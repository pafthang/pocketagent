package memo

import "testing"

func TestConfigListenAddrFromPort(t *testing.T) {
	cfg := &Config{Port: "9090"}
	cfg.normalizeListen()
	if got := cfg.ListenAddr(); got != ":9090" {
		t.Fatalf("ListenAddr() = %q, want :9090", got)
	}
}

func TestConfigListenOverride(t *testing.T) {
	cfg := &Config{Port: "8082", Listen: "0.0.0.0:8082"}
	cfg.normalizeListen()
	if got := cfg.ListenAddr(); got != "0.0.0.0:8082" {
		t.Fatalf("ListenAddr() = %q, want 0.0.0.0:8082", got)
	}
}