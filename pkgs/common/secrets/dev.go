package secrets

import (
	"os"
	"strings"
)

// Local dev bootstrap credentials (see .env.example). Not used when APP_ENV=production.
const (
	DevSuperuserEmail    = "admin@example.com"
	DevSuperuserPassword = "changeme123"
	DevMemoServiceToken  = "dev-memo-token"
)

// ResolvePocketBaseAdmin fills admin credentials for non-production when unset.
func ResolvePocketBaseAdmin(email, password string) (string, string) {
	return resolveDevCredentials(email, password, "POCKETBASE_ADMIN_EMAIL", "POCKETBASE_ADMIN_PASSWORD")
}

// ResolvePocketBaseSuperuser fills bootstrap superuser credentials for non-production when unset.
func ResolvePocketBaseSuperuser(email, password string) (string, string) {
	return resolveDevCredentials(email, password, "POCKETBASE_SUPERUSER_EMAIL", "POCKETBASE_SUPERUSER_PASSWORD")
}

// ResolveMemoServiceToken fills memo service token for non-production when unset.
func ResolveMemoServiceToken(token string) string {
	if v := strings.TrimSpace(os.Getenv("MEMO_SERVICE_TOKEN")); v != "" {
		return v
	}
	if IsProduction() {
		return token
	}
	if strings.TrimSpace(token) == "" {
		return DevMemoServiceToken
	}
	return token
}

func resolveDevCredentials(email, password, emailEnv, passwordEnv string) (string, string) {
	if v := strings.TrimSpace(os.Getenv(emailEnv)); v != "" {
		email = v
	}
	if v := strings.TrimSpace(os.Getenv(passwordEnv)); v != "" {
		password = v
	}
	if IsProduction() {
		return email, password
	}
	if strings.TrimSpace(email) == "" {
		email = DevSuperuserEmail
	}
	if strings.TrimSpace(password) == "" {
		password = DevSuperuserPassword
	}
	return email, password
}