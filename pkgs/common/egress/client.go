package egress

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pafthang/pocketagent/pkgs/common/secrets"
)

const httpTimeout = 20 * time.Second

// Config controls outbound HTTP host allowlisting.
type Config struct {
	Enabled   bool
	Allowlist map[string]struct{}
}

var (
	once       sync.Once
	cfg        Config
	httpClient *http.Client
)

// LoadConfig parses EGRESS_ALLOWLIST (comma-separated hostnames).
func LoadConfig() Config {
	raw := strings.TrimSpace(os.Getenv("EGRESS_ALLOWLIST"))
	hosts := parseHostList(raw)
	if len(hosts) == 0 {
		hosts = defaultHosts()
	}

	enabled := raw != "" || secrets.IsProduction()
	if v := strings.TrimSpace(os.Getenv("EGRESS_ALLOWLIST_ENABLED")); v != "" {
		enabled = v == "1" || strings.EqualFold(v, "true")
	}

	allow := make(map[string]struct{}, len(hosts))
	for _, host := range hosts {
		allow[normalizeHost(host)] = struct{}{}
	}

	return Config{Enabled: enabled, Allowlist: allow}
}

func defaultHosts() []string {
	hosts := []string{
		"api.duckduckgo.com",
		"google.serper.dev",
		"api.tavily.com",
		"localhost",
		"127.0.0.1",
	}
	for _, envKey := range []string{"OLLAMA_URL", "POCKETBASE_URL", "MEMO_URL", "SPACE_URL"} {
		if host := hostFromEnvURL(os.Getenv(envKey)); host != "" {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func parseHostList(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func hostFromEnvURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return normalizeHost(u.Hostname())
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

// ValidateHost returns an error when the host is not allowlisted.
func ValidateHost(host string) error {
	return validateHost(activeConfig(), host)
}

func validateHost(c Config, host string) error {
	if !c.Enabled {
		return nil
	}
	host = normalizeHost(host)
	if host == "" {
		return fmt.Errorf("egress denied: empty host")
	}
	if _, ok := c.Allowlist[host]; ok {
		return nil
	}
	return fmt.Errorf("egress denied for host %q", host)
}

// ValidateURL validates an outbound HTTP(S) URL against the allowlist.
func ValidateURL(rawURL string) error {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return fmt.Errorf("egress denied: unsupported scheme %q", parsed.Scheme)
	}
	return ValidateHost(parsed.Hostname())
}

// HTTPClient returns a shared HTTP client with allowlist enforcement.
func HTTPClient() *http.Client {
	once.Do(initClient)
	return httpClient
}

func activeConfig() Config {
	once.Do(initClient)
	return cfg
}

func initClient() {
	cfg = LoadConfig()
	httpClient = &http.Client{
		Timeout: httpTimeout,
		Transport: &allowlistTransport{
			Base:   http.DefaultTransport,
			Config: cfg,
		},
	}
}

type allowlistTransport struct {
	Base   http.RoundTripper
	Config Config
}

func (t *allowlistTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := ValidateHost(req.URL.Hostname()); err != nil {
		return nil, err
	}
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}