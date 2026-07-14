package tools

import (
	"encoding/json"
	"strings"
)

// ParseArgs converts tool argument JSON strings to a map.
func ParseArgs(raw string) map[string]interface{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]interface{}{}
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &args); err == nil && args != nil {
		return args
	}
	return map[string]interface{}{}
}

// ArgString returns the first non-empty string value for the given keys.
func ArgString(args map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		v, ok := args[key]
		if !ok {
			continue
		}
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}