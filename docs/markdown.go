package docs

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	markdownOnce     sync.Once
	markdownRenderer goldmark.Markdown
)

func Markdown() goldmark.Markdown {
	markdownOnce.Do(func() {
		registerDocsLexerAliases()
		markdownRenderer = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithGuessLanguage(true),
					highlighting.WithFormatOptions(
						chromahtml.WithClasses(true),
						chromahtml.WithLineNumbers(true),
						chromahtml.TabWidth(4),
					),
				),
			),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		)
	})
	return markdownRenderer
}

func RenderHTML(source []byte) (string, error) {
	var html bytes.Buffer
	if err := Markdown().Convert(source, &html); err != nil {
		return "", fmt.Errorf("render Markdown: %w", err)
	}
	return html.String(), nil
}

type aliasedLexer struct {
	chroma.Lexer
	name string
}

func (a aliasedLexer) Config() *chroma.Config {
	config := *a.Lexer.Config()
	config.Name = a.name
	config.Aliases = []string{a.name}
	config.Filenames = nil
	config.MimeTypes = nil
	return &config
}

func registerDocsLexerAliases() {
	for alias, source := range map[string]string{
		"templ":  "go-html-template",
		"dotenv": "bash",
		"env":    "bash",
	} {
		if lexers.Get(alias) != nil {
			continue
		}
		lexer := lexers.Get(source)
		if lexer == nil {
			continue
		}
		lexers.Register(aliasedLexer{Lexer: lexer, name: alias})
	}
}
