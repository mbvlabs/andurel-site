package docs

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

//go:embed content/*/*.md
var content embed.FS

type Heading struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Level int    `json:"level"`
}

type Link struct {
	Title string
	URL   string
}

type Document struct {
	Version     string
	Section     string
	Slug        string
	Title       string
	Description string
	HTML        string
	Source      string
	Headings    []Heading
	Previous    *Link
	Next        *Link
}

func (d *Document) URL() string {
	return "/docs/" + d.Version + "/" + d.Slug
}

type SearchResult struct {
	Title       string    `json:"title"`
	Section     string    `json:"section"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Headings    []Heading `json:"headings"`
	Text        string    `json:"text"`
}

type Site struct {
	versions  []Version
	documents map[string]*Document
	ordered   []*Document
	search    []SearchResult
}

func New() (*Site, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	site := &Site{
		versions:  Catalog,
		documents: make(map[string]*Document),
	}
	expectedFiles := make(map[string]bool)

	for _, version := range Catalog {
		for _, section := range version.Sections {
			for _, page := range section.Pages {
				name := fmt.Sprintf("content/%s/%s.md", version.Name, page.Slug)
				expectedFiles[name] = true
				source, err := content.ReadFile(name)
				if err != nil {
					return nil, fmt.Errorf("documentation: read %s: %w", name, err)
				}
				document, err := render(markdown, version.Name, section.Title, page, source)
				if err != nil {
					return nil, fmt.Errorf("documentation: %s: %w", name, err)
				}
				key := version.Name + "/" + page.Slug
				if _, exists := site.documents[key]; exists {
					return nil, fmt.Errorf("documentation: duplicate page %s", key)
				}
				site.documents[key] = document
				site.ordered = append(site.ordered, document)
			}
		}
	}

	files, err := fs.Glob(content, "content/*/*.md")
	if err != nil {
		return nil, fmt.Errorf("documentation: list content: %w", err)
	}
	for _, name := range files {
		if !expectedFiles[name] {
			return nil, fmt.Errorf("documentation: %s is missing from the catalog", name)
		}
	}

	for i, document := range site.ordered {
		if i > 0 && site.ordered[i-1].Version == document.Version {
			document.Previous = &Link{Title: site.ordered[i-1].Title, URL: site.ordered[i-1].URL()}
		}
		if i+1 < len(site.ordered) && site.ordered[i+1].Version == document.Version {
			document.Next = &Link{Title: site.ordered[i+1].Title, URL: site.ordered[i+1].URL()}
		}
		site.search = append(site.search, SearchResult{
			Title:       document.Title,
			Section:     document.Section,
			Description: document.Description,
			URL:         document.URL(),
			Headings:    document.Headings,
			Text:        document.Source,
		})
	}

	return site, nil
}

func render(
	markdown goldmark.Markdown,
	version, section string,
	page Page,
	source []byte,
) (*Document, error) {
	root := markdown.Parser().Parse(text.NewReader(source))
	headings := make([]Heading, 0)
	var title string

	if err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		heading, ok := node.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		value := nodeText(heading, source)
		if heading.Level == 1 {
			if title != "" {
				return ast.WalkStop, fmt.Errorf("must contain exactly one H1")
			}
			title = value
			return ast.WalkSkipChildren, nil
		}
		if heading.Level == 2 || heading.Level == 3 {
			id, _ := heading.AttributeString("id")
			headings = append(
				headings,
				Heading{ID: string(id.([]byte)), Text: value, Level: heading.Level},
			)
		}
		return ast.WalkSkipChildren, nil
	}); err != nil {
		return nil, err
	}
	if title != page.Title {
		return nil, fmt.Errorf("H1 %q does not match catalog title %q", title, page.Title)
	}

	var html bytes.Buffer
	if err := markdown.Renderer().Render(&html, source, root); err != nil {
		return nil, fmt.Errorf("render Markdown: %w", err)
	}

	return &Document{
		Version:     version,
		Section:     section,
		Slug:        page.Slug,
		Title:       page.Title,
		Description: page.Description,
		HTML:        html.String(),
		Source:      strings.TrimSpace(string(source)),
		Headings:    headings,
	}, nil
}

func nodeText(node ast.Node, source []byte) string {
	var value strings.Builder
	_ = ast.Walk(node, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || child == node {
			return ast.WalkContinue, nil
		}
		if textNode, ok := child.(*ast.Text); ok {
			value.Write(textNode.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return value.String()
}

func (s *Site) Versions() []Version {
	return s.versions
}

func (s *Site) Find(version, slug string) (*Document, bool) {
	document, ok := s.documents[version+"/"+slug]
	return document, ok
}

func (s *Site) LatestURL() string {
	if len(s.ordered) == 0 {
		return "/docs"
	}
	for _, document := range s.ordered {
		if document.Slug == "installation" {
			return document.URL()
		}
	}
	return s.ordered[0].URL()
}

func (s *Site) VersionURL(version string) (string, bool) {
	for _, document := range s.ordered {
		if document.Version == version {
			return document.URL(), true
		}
	}
	return "", false
}

func (s *Site) Documents() []*Document {
	return s.ordered
}

func (s *Site) SearchIndex() []SearchResult {
	return s.search
}
