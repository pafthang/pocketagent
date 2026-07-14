package secrets

import (
	"os"
	"testing"
)

func TestValidateRequiredSecret(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	if err := ValidateRequiredSecret("TEST", "short"); err != nil {
		t.Fatalf("dev should allow short secret: %v", err)
	}

	t.Setenv("APP_ENV", "production")
	if err := ValidateRequiredSecret("TEST", ""); err == nil {
		t.Fatal("expected empty secret error")
	}
	if err := ValidateRequiredSecret("TEST", "short"); err == nil {
		t.Fatal("expected min length error in production")
	}
	if err := ValidateRequiredSecret("TEST", "changeme123456"); err == nil {
		t.Fatal("expected weak secret error")
	}
	if err := ValidateRequiredSecret("TEST", "xK9!mP2vQw7nLz4"); err != nil {
		t.Fatalf("strong secret should pass: %v", err)
	}
}

func TestIsProduction(t *testing.T) {
	for _, env := range []string{"production", "prod", "PRODUCTION"} {
		t.Setenv("APP_ENV", env)
		if !IsProduction() {
			t.Fatalf("expected production for %q", env)
		}
	}
	t.Setenv("APP_ENV", "dev")
	if IsProduction() {
		t.Fatal("expected non-production")
	}
	os.Unsetenv("APP_ENV")
	if IsProduction() {
		t.Fatal("expected non-production when unset")
	}
}