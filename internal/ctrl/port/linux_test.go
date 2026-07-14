//go:build linux

package port

import "testing"

func TestPortHexLE(t *testing.T) {
	tests := map[int]string{
		80:   "5000",
		4222: "7E10",
		8080: "901F",
		8090: "9A1F",
		8081: "911F",
		8082: "921F",
		8083: "931F",
	}
	for port, want := range tests {
		if got := portHexLE(port); got != want {
			t.Fatalf("portHexLE(%d) = %q, want %q", port, got, want)
		}
	}
}

func TestPortHexBE(t *testing.T) {
	if got := portHexBE(8090); got != "1F9A" {
		t.Fatalf("portHexBE(8090) = %q, want 1F9A", got)
	}
}