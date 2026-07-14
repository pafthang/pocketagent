package health

import "testing"

func TestCheckMemoStoreUp(t *testing.T) {
	status := checkMemoStore(func() error { return nil })
	if status.Status != "up" {
		t.Fatalf("expected up, got %#v", status)
	}
}

func TestCheckMemoStoreDown(t *testing.T) {
	status := checkMemoStore(func() error { return errTestDown })
	if status.Status != "down" || status.Error == "" {
		t.Fatalf("expected down with error, got %#v", status)
	}
}

var errTestDown = &healthTestError{"store unavailable"}

type healthTestError struct{ msg string }

func (e *healthTestError) Error() string { return e.msg }