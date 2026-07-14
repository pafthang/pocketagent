package guard

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common/secrets"
)

// ErrPromptRejected indicates a prompt failed injection checks.
var ErrPromptRejected = errors.New("prompt rejected by security policy")

const defaultPromptMaxLen = 32_000

// Mode controls whether violations block or only warn.
type Mode string

const (
	ModeBlock Mode = "block"
	ModeWarn  Mode = "warn"
)

// Config configures prompt injection defenses.
type Config struct {
	Enabled bool
	Mode    Mode
	MaxLen  int
}

var injectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous|prior|above)\s+(instructions|prompts|rules)`),
	regexp.MustCompile(`(?i)disregard\s+(all\s+)?(previous|prior|your)\s+(instructions|prompts|rules)`),
	regexp.MustCompile(`(?i)forget\s+(everything|all)\s+(you\s+)?(know|learned|were told)`),
	regexp.MustCompile(`(?i)you\s+are\s+now\s+(a|an|in)\s+`),
	regexp.MustCompile(`(?i)new\s+system\s+prompt`),
	regexp.MustCompile(`(?i)override\s+(the\s+)?system\s+(prompt|instructions)`),
	regexp.MustCompile(`(?i)jailbreak`),
	regexp.MustCompile(`(?i)\bdan\s+mode\b`),
	regexp.MustCompile(`(?i)do\s+anything\s+now`),
	regexp.MustCompile(`(?i)<\s*/?\s*system\s*>`),
	regexp.MustCompile(`(?i)\[INST\]|\[/INST\]`),
	regexp.MustCompile(`(?i)reveal\s+(the\s+)?(system|hidden)\s+(prompt|instructions)`),
}

// LoadConfig reads guard settings from the environment.
func LoadConfig() Config {
	cfg := Config{
		Enabled: true,
		Mode:    ModeBlock,
		MaxLen:  defaultPromptMaxLen,
	}

	if v := strings.TrimSpace(os.Getenv("PROMPT_GUARD_ENABLED")); v != "" {
		cfg.Enabled = v == "1" || strings.EqualFold(v, "true")
	} else if !secrets.IsProduction() {
		cfg.Mode = ModeWarn
	}

	if mode := strings.ToLower(strings.TrimSpace(os.Getenv("PROMPT_GUARD_MODE"))); mode != "" {
		switch Mode(mode) {
		case ModeWarn:
			cfg.Mode = ModeWarn
		default:
			cfg.Mode = ModeBlock
		}
	}

	if v := strings.TrimSpace(os.Getenv("PROMPT_MAX_LENGTH")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxLen = n
		}
	}

	return cfg
}

// Check validates user-provided prompt text.
func Check(cfg Config, text string) error {
	if !cfg.Enabled {
		return nil
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	if cfg.MaxLen > 0 && len(text) > cfg.MaxLen {
		return fmt.Errorf("%w: exceeds max length %d", ErrPromptRejected, cfg.MaxLen)
	}

	if matches := DetectInjectionPatterns(text); len(matches) > 0 {
		msg := fmt.Errorf("%w: matched %s", ErrPromptRejected, strings.Join(matches, ", "))
		if cfg.Mode == ModeWarn {
			return nil
		}
		return msg
	}

	return nil
}

// DetectInjectionPatterns returns human-readable labels for matched heuristics.
func DetectInjectionPatterns(text string) []string {
	found := make([]string, 0)
	seen := make(map[string]struct{})

	for _, pattern := range injectionPatterns {
		if loc := pattern.FindStringIndex(text); loc != nil {
			label := pattern.String()
			if _, ok := seen[label]; ok {
				continue
			}
			seen[label] = struct{}{}
			found = append(found, truncateMatch(text[loc[0]:loc[1]]))
		}
	}
	return found
}

// Prompt enforces prompt policy and logs warn-mode matches.
func Prompt(log *slog.Logger, cfg Config, text string) error {
	if err := Check(cfg, text); err != nil {
		return err
	}
	if cfg.Enabled && cfg.Mode == ModeWarn {
		if matches := DetectInjectionPatterns(text); len(matches) > 0 && log != nil {
			log.Warn("prompt injection heuristic matched", "matches", matches)
		}
	}
	return nil
}

func truncateMatch(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= 48 {
		return s
	}
	return s[:48] + "..."
}