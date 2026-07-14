package fileapis

import (
	"testing"

	filepath "github.com/pafthang/pocketagent/internal/files/path"
)

func TestMemoDocumentIDForFile(t *testing.T) {
	if got := memoDocumentIDForFile("abc123"); got != "file-abc123" {
		t.Fatalf("unexpected memo id: %q", got)
	}
}

func TestRebaseVirtualPath(t *testing.T) {
	cases := []struct {
		oldRoot, newRoot, path, want string
	}{
		{"/docs", "/archive", "/docs", "/archive"},
		{"/docs", "/archive", "/docs/readme.md", "/archive/readme.md"},
		{"/projects/p1/docs", "/projects/p1/archive", "/projects/p1/docs/a/b.txt", "/projects/p1/archive/a/b.txt"},
	}
	for _, tc := range cases {
		if got := rebaseVirtualPath(tc.oldRoot, tc.newRoot, tc.path); got != tc.want {
			t.Fatalf("rebaseVirtualPath(%q,%q,%q) = %q, want %q", tc.oldRoot, tc.newRoot, tc.path, got, tc.want)
		}
	}
}

func TestResolveScopeMatchesPathPackage(t *testing.T) {
	cases := []struct {
		route, path, query, wantProject, wantDir string
	}{
		{"", "/projects/p1", "", "p1", "/projects/p1"},
		{"p1", "deliverables", "", "p1", "/projects/p1/deliverables"},
		{"", "/readme.md", "", "", "/readme.md"},
		{"", "notes", "p2", "p2", "/projects/p2/notes"},
	}
	for _, tc := range cases {
		scope := filepath.ResolveScope(tc.route, tc.path, tc.query)
		if scope.ProjectID != tc.wantProject || scope.DirPath != tc.wantDir {
			t.Fatalf("ResolveScope(%q,%q,%q) = %#v, want project=%q dir=%q",
				tc.route, tc.path, tc.query, scope, tc.wantProject, tc.wantDir)
		}
	}
}