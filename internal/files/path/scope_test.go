package path

import "testing"

func TestResolveScopeSpaceRoot(t *testing.T) {
	scope := ResolveScope("", "", "")
	if scope.ProjectID != "" || scope.DirPath != "/" {
		t.Fatalf("unexpected: %#v", scope)
	}
}

func TestResolveScopeAbsoluteProjectPath(t *testing.T) {
	scope := ResolveScope("", "/projects/abc/docs", "")
	if scope.ProjectID != "abc" || scope.DirPath != "/projects/abc/docs" {
		t.Fatalf("unexpected: %#v", scope)
	}
}

func TestResolveScopeProjectRouteRelative(t *testing.T) {
	scope := ResolveScope("abc", "docs", "")
	if scope.ProjectID != "abc" || scope.DirPath != "/projects/abc/docs" {
		t.Fatalf("unexpected: %#v", scope)
	}
}

func TestResolveScopeProjectRouteRoot(t *testing.T) {
	scope := ResolveScope("abc", "", "")
	if scope.ProjectID != "abc" || scope.DirPath != "/projects/abc" {
		t.Fatalf("unexpected: %#v", scope)
	}
}

func TestResolveScopeQueryProjectID(t *testing.T) {
	scope := ResolveScope("", "notes", "abc")
	if scope.ProjectID != "abc" || scope.DirPath != "/projects/abc/notes" {
		t.Fatalf("unexpected: %#v", scope)
	}
}

func TestBuildProjectPath(t *testing.T) {
	if got := BuildProjectPath("abc", "docs/readme.md"); got != "/projects/abc/docs/readme.md" {
		t.Fatalf("unexpected: %q", got)
	}
}