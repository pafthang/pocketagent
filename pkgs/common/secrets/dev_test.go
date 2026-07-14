package secrets

import (
	"os"
	"testing"
)

func TestResolveMemoServiceTokenDevDefault(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	os.Unsetenv("MEMO_SERVICE_TOKEN")

	token := ResolveMemoServiceToken("")
	if token != DevMemoServiceToken {
		t.Fatalf("unexpected token: %q", token)
	}
}

func TestResolvePocketBaseAdminDevDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	os.Unsetenv("POCKETBASE_ADMIN_EMAIL")
	os.Unsetenv("POCKETBASE_ADMIN_PASSWORD")

	email, password := ResolvePocketBaseAdmin("", "")
	if email != DevSuperuserEmail || password != DevSuperuserPassword {
		t.Fatalf("unexpected defaults: %s / %s", email, password)
	}
}