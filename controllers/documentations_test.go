package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"andurel-site/docs"

	"github.com/labstack/echo/v5"
)

func TestShowRedirectsV2LatestURLsToHead(t *testing.T) {
	site, err := docs.New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}

	handler := Documentations{site: site}
	request := httptest.NewRequest(http.MethodGet, "/docs/latest/inertia-props", nil)
	recorder := httptest.NewRecorder()
	ctx := echo.New().NewContext(request, recorder)
	ctx.SetPathValues(echo.PathValues{
		{Name: "version", Value: docs.LatestVersion},
		{Name: "slug", Value: "inertia-props"},
	})

	if err := handler.Show(ctx); err != nil {
		t.Fatalf("Show: %v", err)
	}
	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMovedPermanently)
	}
	if got := recorder.Header().Get("Location"); got != "/docs/head/inertia-props" {
		t.Fatalf("Location = %q, want /docs/head/inertia-props", got)
	}
}
