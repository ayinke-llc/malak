package malak

import (
	"bytes"
	"strings"

	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const ErrHeaderNotFound = malakError("could not find title of update")

func getFirstHeader(markdown UpdateContent) (string, error) {
	doc := goldmark.New().
		Parser().
		Parse(text.
			NewReader([]byte(markdown)))

	var title string

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			if heading.Level == 2 {
				var buf bytes.Buffer
				for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					if t, ok := c.(*ast.Text); ok {
						buf.Write(t.Segment.Value([]byte(markdown)))
					}
				}
				title = strings.TrimSpace(buf.String())
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})

	if util.IsStringEmpty(title) {
		return "", ErrHeaderNotFound
	}

	return strings.TrimSpace(title), nil
}
