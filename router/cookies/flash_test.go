package cookies

import (
	"net/http"
	"testing"

	"github.com/mbvlabs/andurel/pkg/server"
)

func TestBuildFlashSessionNameSlugsProjectName(t *testing.T) {
	got := buildFlashSessionName("Andurel Site", "development")
	want := "andurel-site_dev_flash_key"
	if got != want {
		t.Fatalf("buildFlashSessionName(dev) = %q, want %q", got, want)
	}
	if err := (&http.Cookie{Name: got, Value: "x"}).Valid(); err != nil {
		t.Fatalf("dev flash cookie name %q is invalid: %v", got, err)
	}

	got = buildFlashSessionName("Andurel Site", server.ProdEnvironment)
	want = "andurel-site_flash_key"
	if got != want {
		t.Fatalf("buildFlashSessionName(prod) = %q, want %q", got, want)
	}
	if err := (&http.Cookie{Name: got, Value: "x"}).Valid(); err != nil {
		t.Fatalf("prod flash cookie name %q is invalid: %v", got, err)
	}
}
