package client

import "testing"

func TestDLQStreamStatsNilJetStream(t *testing.T) {
	_, err := StreamStats(nil)
	if err == nil {
		t.Fatal("expected error for nil jetstream")
	}
}