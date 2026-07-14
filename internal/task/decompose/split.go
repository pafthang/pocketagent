package decompose

import (
	"strings"
)

func fallbackSplit(task string) []string {
	lower := strings.ToLower(task)
	for _, sep := range []string{" and ", "; ", " then "} {
		if strings.Contains(lower, sep) {
			parts := strings.Split(task, sep)
			out := make([]string, 0, len(parts))
			for _, part := range parts {
				if trimmed := strings.TrimSpace(part); trimmed != "" {
					out = append(out, trimmed)
				}
			}
			if len(out) > 1 {
				return out
			}
		}
	}
	return []string{strings.TrimSpace(task)}
}