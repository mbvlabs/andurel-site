package controllers

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"andurel-site/docs"
	"andurel-site/router/routes"

	"github.com/labstack/echo/v5"
)

func TestIndexNowServesKeyFile(t *testing.T) {
	wantPath := "/" + routes.IndexNowKey + ".txt"
	if routes.IndexNow.Path() != wantPath {
		t.Fatalf("IndexNow path = %q, want %q", routes.IndexNow.Path(), wantPath)
	}

	request := httptest.NewRequest(http.MethodGet, routes.IndexNow.Path(), nil)
	recorder := httptest.NewRecorder()
	ctx := echo.New().NewContext(request, recorder)

	if err := (Assets{}).IndexNow(ctx); err != nil {
		t.Fatalf("IndexNow: %v", err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Body.String(); got != routes.IndexNowKey {
		t.Fatalf("body = %q, want %q", got, routes.IndexNowKey)
	}
	if contentType := recorder.Header().Get("Content-Type"); !strings.Contains(
		strings.ToLower(contentType),
		"text/plain",
	) {
		t.Fatalf("Content-Type = %q, want text/plain", contentType)
	}
}

func TestCreateSitemapIncludesPublicPages(t *testing.T) {
	site, err := docs.New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}

	body, err := createSitemap("https://andurel.com", site)
	if err != nil {
		t.Fatalf("create sitemap: %v", err)
	}

	var sitemap Sitemap
	if err := xml.Unmarshal([]byte(body), &sitemap); err != nil {
		t.Fatalf("parse sitemap: %v", err)
	}

	locs := make(map[string]URL, len(sitemap.URL))
	for _, entry := range sitemap.URL {
		if _, exists := locs[entry.Loc]; exists {
			t.Fatalf("duplicate loc %q", entry.Loc)
		}
		locs[entry.Loc] = entry
	}

	home := "https://andurel.com/"
	if entry, ok := locs[home]; !ok {
		t.Fatal("expected homepage in sitemap")
	} else if entry.Priority != "1.0" {
		t.Fatalf("homepage priority = %q, want 1.0", entry.Priority)
	}

	latestIntro := "https://andurel.com/docs/latest/introduction"
	if _, ok := locs[latestIntro]; !ok {
		t.Fatal("expected latest introduction page in sitemap")
	}

	nested := "https://andurel.com/docs/latest/inertia-props"
	if _, ok := locs[nested]; !ok {
		t.Fatal("expected nested inertia docs page in sitemap")
	}

	legacy := "https://andurel.com/docs/1.5.2/introduction"
	if entry, ok := locs[legacy]; !ok {
		t.Fatal("expected legacy docs page in sitemap")
	} else if entry.Priority != "0.6" {
		t.Fatalf("legacy page priority = %q, want 0.6", entry.Priority)
	}

	wantCount := 1 + len(site.Documents())
	if len(sitemap.URL) != wantCount {
		t.Fatalf("sitemap urls = %d, want %d", len(sitemap.URL), wantCount)
	}

	for _, document := range site.Documents() {
		loc := sitemapLoc("https://andurel.com", document.URL())
		if _, ok := locs[loc]; !ok {
			t.Fatalf("missing documentation page %q", loc)
		}
	}

	for _, loc := range locs {
		if strings.Contains(loc.Loc, "/users/") || strings.Contains(loc.Loc, "/assets/") {
			t.Fatalf("unexpected private loc %q", loc.Loc)
		}
	}

	if strings.Contains(body, routes.Sitemap.URL()) {
		t.Fatal("sitemap should not list itself")
	}
}

func TestSitemapLocJoinsBaseAndPath(t *testing.T) {
	got := sitemapLoc("https://andurel.com/", "/docs/latest/introduction")
	if got != "https://andurel.com/docs/latest/introduction" {
		t.Fatalf("sitemapLoc = %q", got)
	}
}
