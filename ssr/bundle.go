package ssr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"andurel-site/assets"
)

func resolveSSRBundle(bundle string) (string, func(), error) {
	if _, err := os.Stat(bundle); err == nil {
		return bundle, nil, nil
	}

	embedPath := embedPathForBundle(bundle)
	data, err := fs.ReadFile(assets.Files, embedPath)
	if err != nil {
		return "", nil, fmt.Errorf("inertia SSR bundle %q: %w", bundle, err)
	}

	dir, err := os.MkdirTemp("", "inertia-ssr-")
	if err != nil {
		return "", nil, fmt.Errorf("inertia SSR bundle %q: %w", bundle, err)
	}

	cleanup := func() { _ = os.RemoveAll(dir) }
	dest := filepath.Join(dir, filepath.Base(bundle))
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("inertia SSR bundle %q: %w", bundle, err)
	}

	if mapData, mapErr := fs.ReadFile(assets.Files, embedPath+".map"); mapErr == nil {
		_ = os.WriteFile(dest+".map", mapData, 0o644)
	}

	return dest, cleanup, nil
}

func embedPathForBundle(bundle string) string {
	return strings.TrimPrefix(filepath.ToSlash(bundle), "assets/")
}
