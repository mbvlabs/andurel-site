package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	documentation "andurel-site/docs"

	"github.com/labstack/echo/v5"
)

func TestDocumentationHTTPBoundary(t *testing.T) {
	site, err := documentation.New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}
	docs := NewDocs(site)

	t.Run("landing page", func(t *testing.T) {
		status, _, body := runHandler(t, http.MethodGet, "/", Pages{}.Home, nil)
		if status != http.StatusOK || !strings.Contains(body, "Read the documentation") ||
			!strings.Contains(body, "/docs/1.5.2/installation") {
			t.Fatalf("unexpected landing response: status=%d body=%q", status, body)
		}
	})

	t.Run("docs index redirects to latest installation", func(t *testing.T) {
		status, headers, _ := runHandler(t, http.MethodGet, "/docs", docs.Index, nil)
		if status != http.StatusPermanentRedirect ||
			headers.Get("Location") != "/docs/1.5.2/installation" {
			t.Fatalf("unexpected redirect: status=%d location=%q", status, headers.Get("Location"))
		}
	})

	t.Run("known document", func(t *testing.T) {
		status, _, body := runHandler(
			t,
			http.MethodGet,
			"/docs/1.5.2/installation",
			docs.Show,
			map[string]string{
				"version": "1.5.2",
				"slug":    "installation",
			},
		)
		for _, expected := range []string{
			"Installation - andurel",
			`aria-current="page"`,
			`href="#requirements"`,
			`href="/docs/1.5.2/configuration"`,
			`rel="canonical"`,
			`background-color:#282828`,
			`user-select:none`,
			`<span style=`,
		} {
			if !strings.Contains(body, expected) {
				t.Fatalf("document body missing %q", expected)
			}
		}
		if status != http.StatusOK {
			t.Fatalf("document status=%d, want 200", status)
		}
	})

	for _, params := range []map[string]string{
		{"version": "9.x", "slug": "installation"},
		{"version": "1.5.2", "slug": "missing"},
	} {
		t.Run("missing document", func(t *testing.T) {
			status, _, body := runHandler(t, http.MethodGet, "/docs/missing", docs.Show, params)
			if status != http.StatusNotFound || !strings.Contains(body, "Not found") {
				t.Fatalf("unexpected missing response: status=%d body=%q", status, body)
			}
		})
	}

	t.Run("search index", func(t *testing.T) {
		status, headers, body := runHandler(
			t,
			http.MethodGet,
			"/docs/search.json",
			docs.Search,
			nil,
		)
		if status != http.StatusOK ||
			!strings.Contains(headers.Get("Content-Type"), "application/json") {
			t.Fatalf(
				"unexpected search response: status=%d content-type=%q",
				status,
				headers.Get("Content-Type"),
			)
		}
		for _, expected := range []string{"Installation", "/docs/1.5.2/installation", "Requirements"} {
			if !strings.Contains(body, expected) {
				t.Fatalf("search body missing %q", expected)
			}
		}
	})

	t.Run("sitemap", func(t *testing.T) {
		xml, err := createSitemap(site)
		if err != nil {
			t.Fatalf("create sitemap: %v", err)
		}
		for _, expected := range []string{"http://localhost:8080", "/docs/1.5.2/installation", "/docs/1.5.2/configuration"} {
			if !strings.Contains(xml, expected) {
				t.Fatalf("sitemap missing %q", expected)
			}
		}
	})
}

func runHandler(
	t *testing.T,
	method string,
	path string,
	handler echo.HandlerFunc,
	params map[string]string,
) (int, http.Header, string) {
	t.Helper()

	e := echo.New()
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(method, path, nil), recorder)
	pathValues := make(echo.PathValues, 0, len(params))
	for name, value := range params {
		pathValues = append(pathValues, echo.PathValue{Name: name, Value: value})
	}
	ctx.SetPathValues(pathValues)

	if err := handler(ctx); err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	return recorder.Code, recorder.Header(), recorder.Body.String()
}
