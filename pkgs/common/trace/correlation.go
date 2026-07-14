package trace

import (
	"strconv"
	"strings"
)

// RootCorrelationID strips the subtask index suffix (e.g. task-123-0 → task-123).
func RootCorrelationID(corrID string) string {
	if i := strings.LastIndex(corrID, "-"); i > 0 {
		if _, err := strconv.Atoi(corrID[i+1:]); err == nil {
			return corrID[:i]
		}
	}
	return corrID
}

// SubtaskIndex returns the subtask index for a subtask correlation ID, or -1.
func SubtaskIndex(parentCorrID, subCorrID string) int {
	prefix := parentCorrID + "-"
	if !strings.HasPrefix(subCorrID, prefix) {
		return -1
	}
	idx, err := strconv.Atoi(subCorrID[len(prefix):])
	if err != nil {
		return -1
	}
	return idx
}