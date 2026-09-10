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

func TestNestedInertiaPagesLoadInReadingOrder(t *testing.T) {
	site, err := New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}

	overview, ok := site.Find("latest", "inertia")
	if !ok {
		t.Fatal("expected latest/inertia to exist")
	}
	if overview.Parent != nil {
		t.Fatalf("overview parent = %+v, want nil", overview.Parent)
	}
	if overview.Next == nil || overview.Next.URL != "/docs/latest/inertia-renderer" {
		t.Fatalf("overview next = %+v, want renderer", overview.Next)
	}

	props, ok := site.Find("latest", "inertia-props")
	if !ok {
		t.Fatal("expected latest/inertia-props to exist")
	}
	if props.Parent == nil || props.Parent.URL != "/docs/latest/inertia" {
		t.Fatalf("props parent = %+v, want inertia overview", props.Parent)
	}
	if props.Previous == nil || props.Previous.URL != "/docs/latest/inertia-vite" {
		t.Fatalf("props previous = %+v, want vite", props.Previous)
	}
	if !strings.Contains(props.HTML, "FromStruct") {
		t.Fatal("expected props page to document FromStruct")
	}

	generators, ok := site.Find("latest", "inertia-generators")
	if !ok {
		t.Fatal("expected latest/inertia-generators to exist")
	}
	if generators.Next == nil || generators.Next.URL != "/docs/latest/hypermedia" {
		t.Fatalf("generators next = %+v, want hypermedia", generators.Next)
	}
}
