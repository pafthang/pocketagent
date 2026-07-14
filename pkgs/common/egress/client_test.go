package egress

import "testing"

func TestValidateHostAllowlisted(t *testing.T) {
	cfg := Config{
		Enabled: true,
		Allowlist: map[string]struct{}{
			"api.tavily.com": {},
		},
	}

	if err := validateHost(cfg, "api.tavily.com"); err != nil {
		t.Fatalf("expected allowed host: %v", err)
	}
	if err := validateHost(cfg, "evil.example"); err == nil {
		t.Fatal("expected denied host")
	}
}

func TestValidateURLScheme(t *testing.T) {
	if err := ValidateURL("file:///etc/passwd"); err == nil {
		t.Fatal("expected scheme rejection")
	}
}