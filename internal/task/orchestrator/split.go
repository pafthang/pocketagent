package orchestrator

import (
	"fmt"
	"strings"
)

func subtaskCorrelationID(parentCorrID string, index int) string {
	return fmt.Sprintf("%s-%d", parentCorrID, index)
}

func buildFinalAnswer(results map[int]string) string {
	if len(results) == 0 {
		return "No subtask results received"
	}

	maxIdx := -1
	for i := range results {
		if i > maxIdx {
			maxIdx = i
		}
	}

	var sb strings.Builder
	sb.WriteString("Parallel execution result:\n")
	for i := 0; i <= maxIdx; i++ {
		r, ok := results[i]
		if !ok {
			r = "(no result)"
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r))
	}
	return sb.String()
}

func truncateMeta(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}