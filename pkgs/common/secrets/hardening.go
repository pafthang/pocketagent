package secrets

import (
	"fmt"
	"os"
	"strings"
)

// IsProduction reports whether APP_ENV is production.
func IsProduction() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "production", "prod":
		return true
	default:
		return false
	}
}

var weakSecrets = []string{
	"changeme",
	"changeme123",
	"password",
	"secret",
	"admin",
	"admin123",
	"12345678",
	"test",
	"example",
}

// ValidateRequiredSecret ensures a secret is set and meets production policy.
func ValidateRequiredSecret(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	if !IsProduction() {
		return nil
	}
	return validateProductionSecret(name, value)
}

func validateProductionSecret(name, value string) error {
	if len(value) < 12 {
		return fmt.Errorf("%s must be at least 12 characters in production (APP_ENV=production)", name)
	}
	lower := strings.ToLower(value)
	for _, weak := range weakSecrets {
		if lower == weak || strings.Contains(lower, weak) {
			return fmt.Errorf("%s uses a weak or default value; set a strong secret for production", name)
		}
	}
	return nil
}