package views

import (
	"encoding/json"
	"strings"

	"github.com/a-h/templ"
	"github.com/mbvlabs/andurel/pkg/inertia"
)

const (
	siteName        = "Andurel"
	siteTagline     = "Space-grade Go framework for humans and agents"
	siteDescription = "Andurel is the web development framework for Go. Everything you and your agents need to build robust, performant applications."
	twitterSite     = "@mbvlabs"
	twitterCreator  = "@mbvisti"
	defaultOgImage  = "https://media.andurel.com/andurel-og.png"
	defaultOgAlt    = "Andurel wordmark"
	themeColor      = "#090b0d"
	githubURL       = "https://github.com/mbvlabs/andurel"
	xURL            = "https://x.com/mbvlabs"
	organizationURL = "https://mbvlabs.com"
)

type pageSEO struct {
	Title       string
	Description string
	Canonical   string
	Robots      string
	Type        string
	JSONLD      string
}

func stringProp(props map[string]any, key string) string {
	value, _ := props[key].(string)
	return value
}

func inertiaPageSEO(data inertia.RootData) pageSEO {
	canonical := buildCanonicalURL(headDefaults.baseURL, data.Page.URL)
	seo := pageSEO{
		Title:       siteName + " — " + siteTagline,
		Description: siteDescription,
		Canonical:   canonical,
		Robots:      "index, follow",
		Type:        "website",
	}

	component := data.Page.Component
	if strings.HasPrefix(component, "Auth/") || strings.HasPrefix(component, "Errors/") {
		seo.Robots = "noindex, nofollow"
	}

	switch {
	case component == "Home":
		seo.JSONLD = marshalJSONLD(homeJSONLD(canonical))
	case component == "Documentation/Show":
		title := stringProp(data.Page.Props, "title")
		description := stringProp(data.Page.Props, "description")
		if title != "" {
			seo.Title = title + " · Andurel Docs"
		}
		if description != "" {
			seo.Description = description
		}
		seo.Type = "article"
		seo.JSONLD = marshalJSONLD(docsJSONLD(seo, data.Page.Props))
	case component == "Auth/Login":
		seo.Title = "Log in · " + siteName
		seo.Description = "Log in to your Andurel account."
	case component == "Auth/Registration":
		seo.Title = "Create an account · " + siteName
		seo.Description = "Create an Andurel account."
	case component == "Auth/ConfirmEmail":
		seo.Title = "Verify your email · " + siteName
		seo.Description = "Verify your email address to finish creating your Andurel account."
	case component == "Auth/ResetPasswordRequest":
		seo.Title = "Reset password · " + siteName
		seo.Description = "Request a password reset for your Andurel account."
	case component == "Auth/ResetPassword":
		seo.Title = "Choose a new password · " + siteName
		seo.Description = "Choose a new password for your Andurel account."
	case component == "Errors/NotFound":
		seo.Title = "Page not found · " + siteName
		seo.Description = "The page you are looking for could not be found."
	case component == "Errors/BadRequest":
		seo.Title = "Bad request · " + siteName
		seo.Description = "The request made was invalid."
	case component == "Errors/InternalError":
		seo.Title = "Something went wrong · " + siteName
		seo.Description = "The application hit an unexpected error."
	}

	return seo
}

func marshalJSONLD(value any) string {
	if value == nil {
		return ""
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func jsonLDScript(json string) templ.Component {
	return templ.Raw(`<script type="application/ld+json" head-key="json-ld">` + json + `</script>`)
}

func organizationNode() map[string]any {
	return map[string]any{
		"@type":  "Organization",
		"@id":    organizationURL + "#organization",
		"name":   "MBV Labs",
		"url":    organizationURL,
		"sameAs": []string{githubURL, xURL},
	}
}

func homeJSONLD(canonical string) map[string]any {
	org := organizationNode()
	return map[string]any{
		"@context": "https://schema.org",
		"@graph": []any{
			org,
			map[string]any{
				"@type":       "WebSite",
				"@id":         canonical + "#website",
				"name":        siteName,
				"url":         canonical,
				"description": siteDescription,
				"publisher":   map[string]any{"@id": org["@id"]},
				"inLanguage":  "en-US",
			},
			map[string]any{
				"@type":               "SoftwareApplication",
				"name":                siteName,
				"applicationCategory": "DeveloperApplication",
				"operatingSystem":     "Linux, macOS",
				"programmingLanguage": "Go",
				"url":                 canonical,
				"description":         siteDescription,
				"downloadUrl":         githubURL,
				"author":              map[string]any{"@id": org["@id"]},
				"publisher":           map[string]any{"@id": org["@id"]},
				"offers": map[string]any{
					"@type":         "Offer",
					"price":         "0",
					"priceCurrency": "USD",
				},
			},
		},
	}
}

func docsJSONLD(seo pageSEO, props map[string]any) map[string]any {
	org := organizationNode()
	version := stringProp(props, "currentVersion")
	section := stringProp(props, "currentSection")
	headline := stringProp(props, "title")
	if headline == "" {
		headline = seo.Title
	}

	elements := []map[string]any{
		{"@type": "ListItem", "position": 1, "name": siteName, "item": buildCanonicalURL(headDefaults.baseURL, "/")},
		{"@type": "ListItem", "position": 2, "name": "Docs", "item": buildCanonicalURL(headDefaults.baseURL, "/docs/latest/introduction")},
	}
	if version != "" {
		elements = append(elements, map[string]any{
			"@type":    "ListItem",
			"position": 3,
			"name":     version,
			"item":     buildCanonicalURL(headDefaults.baseURL, "/docs/"+version),
		})
	}
	if section != "" {
		elements = append(elements, map[string]any{
			"@type":    "ListItem",
			"position": len(elements) + 1,
			"name":     section,
		})
	}
	elements = append(elements, map[string]any{
		"@type":    "ListItem",
		"position": len(elements) + 1,
		"name":     headline,
		"item":     seo.Canonical,
	})

	return map[string]any{
		"@context": "https://schema.org",
		"@graph": []any{
			org,
			map[string]any{
				"@type":               "TechArticle",
				"headline":            headline,
				"description":         seo.Description,
				"url":                 seo.Canonical,
				"inLanguage":          "en-US",
				"isAccessibleForFree": true,
				"author":              map[string]any{"@id": org["@id"]},
				"publisher":           map[string]any{"@id": org["@id"]},
				"mainEntityOfPage":    seo.Canonical,
			},
			map[string]any{
				"@type":           "BreadcrumbList",
				"itemListElement": elements,
			},
		},
	}
}
