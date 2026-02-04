package markdown

import (
	"bytes"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Renderer wraps goldmark markdown parser with HTML sanitization
type Renderer struct {
	md     goldmark.Markdown
	policy *bluemonday.Policy
}

// New creates a new markdown renderer with common extensions
func New() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	// UGCPolicy allows GitHub-compatible HTML tags while blocking scripts
	policy := bluemonday.UGCPolicy()

	// Allow checkbox inputs rendered by goldmark's task list extension
	policy.AllowAttrs("type").
		Matching(regexp.MustCompile(`(?i)^checkbox$`)).
		OnElements("input")
	policy.AllowAttrs("checked", "disabled").OnElements("input")

	return &Renderer{md: md, policy: policy}
}

// Render converts markdown to HTML with sanitization
func (r *Renderer) Render(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return r.policy.SanitizeBytes(buf.Bytes()), nil
}
