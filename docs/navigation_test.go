package docs

import "testing"

func TestFindLatestInstallation(t *testing.T) {
	_, section, page, ok := Find(LatestVersion, "installation")
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
	if _, _, _, ok := Find(LatestVersion, "missing"); ok {
		t.Fatal("expected missing page to be absent")
	}
}

func TestLatestURL(t *testing.T) {
	if got, want := LatestURL(), "/docs/latest/introduction"; got != want {
		t.Fatalf("LatestURL() = %q, want %q", got, want)
	}
}

func TestNavigationOrder(t *testing.T) {
	versions := Navigation()
	want := []string{LatestVersion, LatestRelease, V152Release, HeadVersion}
	if len(versions) != len(want) {
		t.Fatalf("versions = %d, want %d", len(versions), len(want))
	}
	for i, name := range want {
		if versions[i].Name != name {
			t.Fatalf("versions[%d] = %q, want %q", i, versions[i].Name, name)
		}
	}
}

func TestNavigationIncludesLatestIntroduction(t *testing.T) {
	versions := Navigation()
	if len(versions) == 0 {
		t.Fatal("expected at least one version")
	}
	if versions[0].Name != LatestVersion {
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
	url, ok := VersionURL(V152Release)
	if !ok {
		t.Fatal("expected 1.5.2 to exist")
	}
	if url != "/docs/1.5.2/introduction" {
		t.Fatalf("VersionURL(1.5.2) = %q, want /docs/1.5.2/introduction", url)
	}

	url, ok = VersionURL(LatestRelease)
	if !ok {
		t.Fatal("expected 1.5.5 to exist")
	}
	if url != "/docs/1.5.5/introduction" {
		t.Fatalf("VersionURL(1.5.5) = %q, want /docs/1.5.5/introduction", url)
	}

	url, ok = VersionURL(HeadVersion)
	if !ok {
		t.Fatal("expected head to exist")
	}
	if url != "/docs/head/whats-new" {
		t.Fatalf("VersionURL(head) = %q, want /docs/head/whats-new", url)
	}
}

func TestFindNestedInertiaPage(t *testing.T) {
	_, section, page, ok := Find(HeadVersion, "inertia-props")
	if !ok {
		t.Fatal("expected head/inertia-props to exist")
	}
	if section.Title != "Inertia" {
		t.Fatalf("section = %q, want Inertia", section.Title)
	}
	if page.Title != "Props" {
		t.Fatalf("title = %q, want Props", page.Title)
	}
}

func TestLatestDoesNotIncludeHeadPages(t *testing.T) {
	if _, _, _, ok := Find(LatestVersion, "inertia-props"); ok {
		t.Fatal("expected latest to follow the stable v1 catalog")
	}
}

func TestMovedFromLatest(t *testing.T) {
	dest, ok := MovedFromLatest("inertia-props")
	if !ok {
		t.Fatal("expected v2-only latest URLs to move to head")
	}
	if dest != "/docs/head/inertia-props" {
		t.Fatalf("redirect = %q, want /docs/head/inertia-props", dest)
	}
	if _, ok := MovedFromLatest("installation"); ok {
		t.Fatal("expected shared latest pages to stay on latest")
	}
}

func TestNavigationNestsInertiaChildren(t *testing.T) {
	versions := Navigation()
	head := findNavVersion(t, versions, HeadVersion)

	var inertia NavigationPage
	for _, section := range head.Sections {
		if section.Title != "Inertia" {
			continue
		}
		for _, page := range section.Pages {
			if page.Slug == "inertia" {
				inertia = page
			}
		}
	}
	if inertia.Slug == "" {
		t.Fatal("expected Inertia / Inertia in head navigation")
	}
	if inertia.URL != "/docs/head/inertia" {
		t.Fatalf("inertia URL = %q, want /docs/head/inertia", inertia.URL)
	}

	want := []string{
		"inertia-renderer",
		"inertia-pages",
		"inertia-props",
		"inertia-shared",
		"inertia-vite",
		"inertia-ssr",
		"inertia-typescript-sync",
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
		if inertia.Children[i].URL != "/docs/head/"+slug {
			t.Fatalf("children[%d].URL = %q", i, inertia.Children[i].URL)
		}
	}
}

func findNavVersion(t *testing.T, versions []NavigationVersion, name string) NavigationVersion {
	t.Helper()
	for _, version := range versions {
		if version.Name == name {
			return version
		}
	}
	t.Fatalf("expected navigation version %q", name)
	return NavigationVersion{}
}
