package ssr

import (
	"os"
	"testing"
)

func TestEmbedPathForBundle(t *testing.T) {
	if got := embedPathForBundle("assets/dist/ssr/ssr.js"); got != "dist/ssr/ssr.js" {
		t.Fatalf("embedPathForBundle = %q, want dist/ssr/ssr.js", got)
	}
}

func TestResolveSSRBundleExtractsFromEmbed(t *testing.T) {
	t.Chdir(t.TempDir())

	path, cleanup, err := resolveSSRBundle("assets/dist/ssr/ssr.js")
	if err != nil {
		t.Fatalf("resolveSSRBundle: %v", err)
	}
	if cleanup == nil {
		t.Fatal("cleanup = nil, want temp-dir cleanup after embed extract")
	}
	defer cleanup()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("extracted bundle: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("extracted bundle is empty")
	}
}
