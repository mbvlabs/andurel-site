package docs

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func ExtractHeadings(source []byte) ([]Heading, error) {
	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	root := markdown.Parser().Parse(text.NewReader(source))
	headings := make([]Heading, 0)

	if err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		heading, ok := node.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		if heading.Level != 2 && heading.Level != 3 {
			return ast.WalkSkipChildren, nil
		}

		idAttr, hasID := heading.AttributeString("id")
		if !hasID {
			return ast.WalkSkipChildren, nil
		}
		idBytes, ok := idAttr.([]byte)
		if !ok {
			return ast.WalkSkipChildren, fmt.Errorf("heading id is not bytes")
		}

		headings = append(headings, Heading{
			ID:    string(idBytes),
			Text:  nodeText(heading, source),
			Level: heading.Level,
		})
		return ast.WalkSkipChildren, nil
	}); err != nil {
		return nil, err
	}

	return headings, nil
}
