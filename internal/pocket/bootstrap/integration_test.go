package bootstrap

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestRunCreatesExpectedCollections(t *testing.T) {
	dir := t.TempDir()
	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: dir})
	if err := app.Bootstrap(); err != nil {
		t.Fatalf("bootstrap app: %v", err)
	}

	if err := Run(app, Config{
		SuperuserEmail:    "admin@test.local",
		SuperuserPassword: "testpassword123",
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	expected := ExpectedCollections()
	for _, name := range expected {
		if _, err := app.FindCollectionByNameOrId(name); err != nil {
			t.Errorf("missing collection %q: %v", name, err)
		}
	}

	collections, err := app.FindAllCollections()
	if err != nil {
		t.Fatalf("FindAllCollections: %v", err)
	}
	if len(collections) < len(expected) {
		t.Fatalf("got %d collections, want at least %d", len(collections), len(expected))
	}
}