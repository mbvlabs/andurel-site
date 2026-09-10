package views

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/server"
)

func TestInertiaPageSEOHome(t *testing.T) {
	ConfigureHead("andurel-site", "https://andurel.com")

	seo := inertiaPageSEO(inertia.RootData{
		Page: inertia.Page{Component: "Home", URL: "/"},
	})
	if seo.Title != siteName+" — "+siteTagline {
		t.Fatalf("title = %q", seo.Title)
	}
	if seo.Canonical != "https://andurel.com/" {
		t.Fatalf("canonical = %q", seo.Canonical)
	}
	if seo.Type != "website" || seo.Robots != "index, follow" {
		t.Fatalf("type/robots = %q %q", seo.Type, seo.Robots)
	}
	if !strings.Contains(seo.JSONLD, `"SoftwareApplication"`) {
		t.Fatalf("expected software application json-ld, got %s", seo.JSONLD)
	}
}

func TestInertiaPageSEODocs(t *testing.T) {
	ConfigureHead("andurel-site", "https://andurel.com")

	seo := inertiaPageSEO(inertia.RootData{
		Page: inertia.Page{
			Component: "Documentation/Show",
			URL:       "/docs/latest/introduction",
			Props: map[string]any{
				"title":          "Introduction",
				"description":    "Meet Andurel v2.",
				"currentVersion": "latest",
				"currentSection": "Getting Started",
			},
		},
	})
	if seo.Title != "Introduction · Andurel Docs" {
		t.Fatalf("title = %q", seo.Title)
	}
	if seo.Description != "Meet Andurel v2." {
		t.Fatalf("description = %q", seo.Description)
	}
	if seo.Type != "article" {
		t.Fatalf("type = %q", seo.Type)
	}
	if !strings.Contains(seo.JSONLD, `"TechArticle"`) || !strings.Contains(seo.JSONLD, `"BreadcrumbList"`) {
		t.Fatalf("expected article breadcrumbs json-ld, got %s", seo.JSONLD)
	}
}

func TestInertiaPageSEOAuthNoIndex(t *testing.T) {
	ConfigureHead("andurel-site", "https://andurel.com")

	seo := inertiaPageSEO(inertia.RootData{
		Page: inertia.Page{Component: "Auth/Login", URL: "/users/sign-in"},
	})
	if seo.Robots != "noindex, nofollow" {
		t.Fatalf("robots = %q", seo.Robots)
	}
	if seo.Title != "Log in · Andurel" {
		t.Fatalf("title = %q", seo.Title)
	}
}

func TestRootRendersSEOTags(t *testing.T) {
	ConfigureHead("andurel-site", "https://andurel.com")

	component := Root(inertia.RootData{
		Page:        inertia.Page{Component: "Home", URL: "/"},
		ContainerID: "app",
		ProjectName: "andurel-site",
		Environment: server.DevEnvironment,
		PageJSON:    []byte(`{"component":"Home","props":{},"url":"/","version":"1"}`),
	})

	var buf strings.Builder
	if err := component.Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	html := buf.String()
	for _, needle := range []string{
		`property="og:site_name"`,
		`name="twitter:card"`,
		`rel="canonical"`,
		`type="application/ld+json"`,
		`SoftwareApplication`,
	} {
		if !strings.Contains(html, needle) {
			t.Fatalf("root html missing %q", needle)
		}
	}
}

func TestRootDoesNotPanicWithoutPageJSON(t *testing.T) {
	ConfigureHead("andurel-site", "https://andurel.com")

	component := Root(inertia.RootData{
		Page:        inertia.Page{Component: "Errors/NotFound", URL: "/missing"},
		ContainerID: "app",
		Environment: server.DevEnvironment,
		PageJSON:    []byte(`{"component":"Errors/NotFound","props":{},"url":"/missing","version":"1"}`),
	})
	if err := component.Render(context.Background(), io.Discard); err != nil {
		t.Fatal(err)
	}
}
