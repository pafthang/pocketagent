package decompose

import "testing"

func TestParseSubtasksResponse(t *testing.T) {
	t.Parallel()

	raw := `Here is the plan:
["research AI trends", "summarize findings"]`

	got, err := parseSubtasksResponse(raw, 4)
	if err != nil {
		t.Fatalf("parseSubtasksResponse() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(subtasks) = %d, want 2", len(got))
	}
	if got[0] != "research AI trends" {
		t.Fatalf("got[0] = %q", got[0])
	}
}

func TestParseSubtasksResponseDeduplicates(t *testing.T) {
	t.Parallel()

	got, err := parseSubtasksResponse(`["do thing", "do thing", "next"]`, 4)
	if err != nil {
		t.Fatalf("parseSubtasksResponse() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(subtasks) = %d, want 2", len(got))
	}
}