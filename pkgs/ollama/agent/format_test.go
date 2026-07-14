package agent

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/ollama/api"
)

func TestFormatToolArguments(t *testing.T) {
	if got := api.FormatToolArguments(map[string]interface{}{"query": "test"}); got != "test" {
		t.Fatalf("expected query extraction, got %q", got)
	}
}