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

	document, ok := site.Find(LatestVersion, "installation")
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
	if !strings.Contains(document.HTML, "/docs/latest/frontend-options") {
		t.Fatal("expected latest alias to rewrite shared 1.5.5 links")
	}

	views, ok := site.Find(LatestVersion, "views")
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

	pinned, ok := site.Find(LatestRelease, "installation")
	if !ok {
		t.Fatal("expected 1.5.5/installation to exist")
	}
	if !strings.Contains(pinned.HTML, "/docs/1.5.5/frontend-options") {
		t.Fatal("expected 1.5.5 pages to keep pinned version links")
	}

	legacy, ok := site.Find(V152Release, "introduction")
	if !ok {
		t.Fatal("expected 1.5.2/introduction to exist")
	}
	if legacy.Title != "Introduction" {
		t.Fatalf("legacy title = %q, want Introduction", legacy.Title)
	}

	head, ok := site.Find(HeadVersion, "installation")
	if !ok {
		t.Fatal("expected head/installation to exist")
	}
	if !strings.Contains(head.HTML, "@master") {
		t.Fatal("expected head installation to document the master CLI")
	}

	if _, ok := site.Find(LatestVersion, "missing"); ok {
		t.Fatal("expected missing page to be absent")
	}
}

func TestNestedInertiaPagesLoadInReadingOrder(t *testing.T) {
	site, err := New()
	if err != nil {
		t.Fatalf("load documentation: %v", err)
	}

	overview, ok := site.Find(HeadVersion, "inertia")
	if !ok {
		t.Fatal("expected head/inertia to exist")
	}
	if overview.Parent != nil {
		t.Fatalf("overview parent = %+v, want nil", overview.Parent)
	}
	if overview.Next == nil || overview.Next.URL != "/docs/head/inertia-renderer" {
		t.Fatalf("overview next = %+v, want renderer", overview.Next)
	}

	props, ok := site.Find(HeadVersion, "inertia-props")
	if !ok {
		t.Fatal("expected head/inertia-props to exist")
	}
	if props.Parent == nil || props.Parent.URL != "/docs/head/inertia" {
		t.Fatalf("props parent = %+v, want inertia overview", props.Parent)
	}
	if props.Previous == nil || props.Previous.URL != "/docs/head/inertia-vite" {
		t.Fatalf("props previous = %+v, want vite", props.Previous)
	}
	if !strings.Contains(props.HTML, "FromStruct") {
		t.Fatal("expected props page to document FromStruct")
	}

	generators, ok := site.Find(HeadVersion, "inertia-generators")
	if !ok {
		t.Fatal("expected head/inertia-generators to exist")
	}
	if generators.Next == nil || generators.Next.URL != "/docs/head/hypermedia" {
		t.Fatalf("generators next = %+v, want hypermedia", generators.Next)
	}
}
