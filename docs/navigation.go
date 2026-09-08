package docs

const docsPrefix = "/docs"

type NavigationPage struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type NavigationSection struct {
	Title string           `json:"title"`
	Pages []NavigationPage `json:"pages"`
}

type NavigationVersion struct {
	Name     string              `json:"name"`
	URL      string              `json:"url"`
	Sections []NavigationSection `json:"sections"`
}

func PageURL(version, slug string) string {
	return docsPrefix + "/" + version + "/" + slug
}

func VersionPath(version string) string {
	return docsPrefix + "/" + version
}

func LatestURL() string {
	if len(Catalog) == 0 {
		return docsPrefix
	}

	latest := Catalog[0]
	for _, section := range latest.Sections {
		for _, page := range section.Pages {
			if page.Slug == "introduction" {
				return PageURL(latest.Name, page.Slug)
			}
		}
	}

	url, ok := VersionURL(latest.Name)
	if ok {
		return url
	}

	return docsPrefix
}

func VersionURL(version string) (string, bool) {
	for _, entry := range Catalog {
		if entry.Name != version {
			continue
		}
		for _, section := range entry.Sections {
			if len(section.Pages) == 0 {
				continue
			}
			return PageURL(entry.Name, section.Pages[0].Slug), true
		}
	}

	return "", false
}

func Find(version, slug string) (Version, Section, Page, bool) {
	for _, entry := range Catalog {
		if entry.Name != version {
			continue
		}
		for _, section := range entry.Sections {
			for _, page := range section.Pages {
				if page.Slug == slug {
					return entry, section, page, true
				}
			}
		}
	}

	return Version{}, Section{}, Page{}, false
}

func Navigation() []NavigationVersion {
	versions := make([]NavigationVersion, 0, len(Catalog))
	for _, version := range Catalog {
		url, ok := VersionURL(version.Name)
		if !ok {
			url = VersionPath(version.Name)
		}

		sections := make([]NavigationSection, 0, len(version.Sections))
		for _, section := range version.Sections {
			pages := make([]NavigationPage, 0, len(section.Pages))
			for _, page := range section.Pages {
				pages = append(pages, NavigationPage{
					Slug:        page.Slug,
					Title:       page.Title,
					Description: page.Description,
					URL:         PageURL(version.Name, page.Slug),
				})
			}
			sections = append(sections, NavigationSection{
				Title: section.Title,
				Pages: pages,
			})
		}

		versions = append(versions, NavigationVersion{
			Name:     version.Name,
			URL:      url,
			Sections: sections,
		})
	}

	return versions
}
