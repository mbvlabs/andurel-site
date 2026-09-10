package docs

import "testing"

func TestFindLatestInstallation(t *testing.T) {
	_, section, page, ok := Find("latest", "installation")
	if !ok {
		t.Fatal("expected latest/installation to exist")
	}
	if section.Title != "Getting Started" {
		t.Fatalf("section = %q, want Getting Started", section.Title)
	}
	if page.Title != "Installation" {
		t.Fatalf("title = %q, want Installation", page.Title)
	}
}

func TestFindUnknownPage(t *testing.T) {
	if _, _, _, ok := Find("latest", "missing"); ok {
		t.Fatal("expected missing page to be absent")
	}
}

func TestLatestURL(t *testing.T) {
	if got, want := LatestURL(), "/docs/latest/introduction"; got != want {
		t.Fatalf("LatestURL() = %q, want %q", got, want)
	}
}

func TestNavigationIncludesCurrentMasterSections(t *testing.T) {
	versions := Navigation()
	if len(versions) == 0 {
		t.Fatal("expected at least one version")
	}
	if versions[0].Name != "latest" {
		t.Fatalf("first version = %q, want latest", versions[0].Name)
	}

	var found bool
	for _, section := range versions[0].Sections {
		if section.Title != "Getting Started" {
			continue
		}
		for _, page := range section.Pages {
			if page.Slug == "introduction" && page.URL == "/docs/latest/introduction" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("expected Getting Started / introduction in latest navigation")
	}
}

func TestVersionURL(t *testing.T) {
	url, ok := VersionURL("1.5.2")
	if !ok {
		t.Fatal("expected 1.5.2 to exist")
	}
	if url != "/docs/1.5.2/introduction" {
		t.Fatalf("VersionURL(1.5.2) = %q, want /docs/1.5.2/introduction", url)
	}
}

func TestFindNestedInertiaPage(t *testing.T) {
	_, section, page, ok := Find("latest", "inertia-props")
	if !ok {
		t.Fatal("expected latest/inertia-props to exist")
	}
	if section.Title != "Framework Packages" {
		t.Fatalf("section = %q, want Framework Packages", section.Title)
	}
	if page.Title != "Props" {
		t.Fatalf("title = %q, want Props", page.Title)
	}
}

func TestNavigationNestsInertiaChildren(t *testing.T) {
	versions := Navigation()
	if len(versions) == 0 {
		t.Fatal("expected at least one version")
	}

	var inertia NavigationPage
	for _, section := range versions[0].Sections {
		if section.Title != "Framework Packages" {
			continue
		}
		for _, page := range section.Pages {
			if page.Slug == "inertia" {
				inertia = page
			}
		}
	}
	if inertia.Slug == "" {
		t.Fatal("expected Framework Packages / Inertia in latest navigation")
	}
	if inertia.URL != "/docs/latest/inertia" {
		t.Fatalf("inertia URL = %q, want /docs/latest/inertia", inertia.URL)
	}

	want := []string{
		"inertia-renderer",
		"inertia-pages",
		"inertia-vite",
		"inertia-props",
		"inertia-shared",
		"inertia-ssr",
		"inertia-diagnostics",
		"inertia-generators",
	}
	if len(inertia.Children) != len(want) {
		t.Fatalf("inertia children = %d, want %d", len(inertia.Children), len(want))
	}
	for i, slug := range want {
		if inertia.Children[i].Slug != slug {
			t.Fatalf("children[%d].Slug = %q, want %q", i, inertia.Children[i].Slug, slug)
		}
		if inertia.Children[i].URL != "/docs/latest/"+slug {
			t.Fatalf("children[%d].URL = %q", i, inertia.Children[i].URL)
		}
	}
}
