package docs

import (
	"strings"
	"testing"
)

func TestSiteLoadsCatalogContent(t *testing.T) {
	site, err := New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}

	document, ok := site.Find("latest", "installation")
	if !ok {
		t.Fatal("expected latest/installation to exist")
	}
	if document.Title != "Installation" {
		t.Fatalf("title = %q, want Installation", document.Title)
	}
	if document.Section != "Getting Started" {
		t.Fatalf("section = %q, want Getting Started", document.Section)
	}
	if !strings.Contains(document.HTML, "Requirements") {
		t.Fatal("expected rendered HTML to include Requirements")
	}
	if !strings.Contains(document.HTML, `class="chroma"`) {
		t.Fatal("expected highlighted code to use chroma classes")
	}
	if strings.Contains(document.HTML, "background-color:#282828") {
		t.Fatal("highlighted code still uses gruvbox inline colors")
	}

	views, ok := site.Find("latest", "views")
	if !ok {
		t.Fatal("expected latest/views to exist")
	}
	if !strings.Contains(views.HTML, `class="chroma"`) {
		t.Fatal("expected templ examples to use the same highlighter as Go examples")
	}
	if document.Previous == nil || document.Previous.URL != "/docs/latest/introduction" {
		t.Fatalf("previous = %+v, want introduction", document.Previous)
	}
	if document.Next == nil || document.Next.URL != "/docs/latest/configuration" {
		t.Fatalf("next = %+v, want configuration", document.Next)
	}

	legacy, ok := site.Find("1.5.2", "introduction")
	if !ok {
		t.Fatal("expected 1.5.2/introduction to exist")
	}
	if legacy.Title != "Introduction" {
		t.Fatalf("legacy title = %q, want Introduction", legacy.Title)
	}

	if _, ok := site.Find("latest", "missing"); ok {
		t.Fatal("expected missing page to be absent")
	}
}
