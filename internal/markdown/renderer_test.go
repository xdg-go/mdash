package markdown

import (
	"strings"
	"testing"
)

func TestRender_TaskList(t *testing.T) {
	r := New()

	input := []byte("- [ ] unchecked\n- [x] checked\n")
	got, err := r.Render(input)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	html := string(got)

	// Unchecked checkbox must be present
	if !strings.Contains(html, `<input`) {
		t.Errorf("expected <input> element in output, got:\n%s", html)
	}
	if !strings.Contains(html, `type="checkbox"`) {
		t.Errorf("expected type=\"checkbox\" in output, got:\n%s", html)
	}
	if !strings.Contains(html, `disabled`) {
		t.Errorf("expected disabled attribute in output, got:\n%s", html)
	}

	// Checked checkbox must have checked attribute
	if !strings.Contains(html, `checked`) {
		t.Errorf("expected checked attribute in output, got:\n%s", html)
	}
}

// A list starting at 0 must keep its numbering; the sanitizer used to strip
// the start attribute goldmark emits.
func TestRenderKeepsOrderedListStart(t *testing.T) {
	out, err := New().Render([]byte("0. zero\n1. one\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `<ol start="0">`) {
		t.Errorf("start attribute dropped:\n%s", out)
	}
}
