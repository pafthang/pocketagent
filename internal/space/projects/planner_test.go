package projects

import (
	"context"
	"testing"
)

func TestParseGoalWithoutOllama(t *testing.T) {
	analysis, err := ParseGoal(context.Background(), nil, "", "Build a REST API for todos")
	if err != nil {
		t.Fatalf("ParseGoal: %v", err)
	}
	if analysis["domain"] == nil {
		t.Fatal("expected domain in analysis")
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Fatalf("truncate short: %q", got)
	}
	if got := truncate("hello world", 5); got != "hello..." {
		t.Fatalf("truncate long: %q", got)
	}
}
