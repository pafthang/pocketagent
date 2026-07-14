package auth

import "testing"

func TestHashTokenDeterministic(t *testing.T) {
	a := HashToken("test-token")
	b := HashToken("test-token")
	if a != b || a == "" {
		t.Fatalf("expected stable non-empty hash, got %q and %q", a, b)
	}
	if HashToken("other") == a {
		t.Fatal("expected different hash for different token")
	}
}

func TestGenerateTokenUnique(t *testing.T) {
	t1, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	t2, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	if t1 == "" || t2 == "" || t1 == t2 {
		t.Fatalf("expected unique tokens, got %q %q", t1, t2)
	}
}