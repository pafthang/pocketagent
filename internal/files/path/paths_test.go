package path

import "testing"

func TestParseBrowsePathProject(t *testing.T) {
	pid, dir := ParseBrowsePath("/projects/abc/docs")
	if pid != "abc" || dir != "/projects/abc/docs" {
		t.Fatalf("unexpected: %q %q", pid, dir)
	}
}

func TestParseBrowsePathProjectRoot(t *testing.T) {
	pid, dir := ParseBrowsePath("/projects/abc")
	if pid != "abc" || dir != "/projects/abc" {
		t.Fatalf("unexpected: %q %q", pid, dir)
	}
}

func TestJoinPath(t *testing.T) {
	if JoinPath("/", "readme.md") != "/readme.md" {
		t.Fatal("join root failed")
	}
}