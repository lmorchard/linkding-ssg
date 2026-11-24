package templates

import (
	"strings"
	"testing"
	"time"

	"github.com/lmorchard/linkding-ssg/pkg/linkding"
)

func TestRenderNotesMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "empty",
			input:    "",
			contains: nil,
		},
		{
			name:     "whitespace only",
			input:    "   \n\t  ",
			contains: nil,
		},
		{
			name:     "plain text",
			input:    "This is a simple personal comment.",
			contains: []string{"<p>This is a simple personal comment.</p>"},
		},
		{
			name:     "markdown with emphasis and links",
			input:    "My note with *emphasis* and [link](https://example.com).",
			contains: []string{"<em>emphasis</em>", `<a href="https://example.com">link</a>`},
		},
		{
			name:     "markdown with list and code",
			input:    "Notes:\n- point one\n- point two with `code`",
			contains: []string{"<ul>", "<li>point one</li>", "<code>code</code>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(renderNotesMarkdown(tt.input))
			if len(tt.contains) == 0 {
				if got != "" {
					t.Errorf("expected empty string for %q, got %q", tt.input, got)
				}
				return
			}
			for _, expected := range tt.contains {
				if !strings.Contains(got, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, got)
				}
			}
		})
	}
}

func TestBookmarkTemplates_RenderNotes(t *testing.T) {
	set, err := LoadTemplates("", "", "Test Site", nil)
	if err != nil {
		t.Fatalf("LoadTemplates failed: %v", err)
	}

	bm := linkding.Bookmark{
		ID:          42,
		URL:         "https://example.com/test",
		Title:       "Test Bookmark",
		Description: "A test description",
		Notes:       "My personal commentary with **bold** and [a link](https://example.com/target).",
		DateCreated: time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC),
	}

	// 1. Test single bookmark permalink template
	var permalinkBuf strings.Builder
	err = set.BookmarkTemplate.ExecuteTemplate(&permalinkBuf, "base", map[string]any{
		"Bookmark": bm,
	})
	if err != nil {
		t.Fatalf("BookmarkTemplate execution failed: %v", err)
	}
	renderedPermalink := permalinkBuf.String()

	if !strings.Contains(renderedPermalink, `<div class="bookmark-notes">`) {
		t.Errorf("expected permalink to contain .bookmark-notes div")
	}
	if !strings.Contains(renderedPermalink, `<strong>bold</strong>`) {
		t.Errorf("expected permalink notes to render markdown strong tag, got: %s", renderedPermalink)
	}
	if !strings.Contains(renderedPermalink, `<a href="https://example.com/target">a link</a>`) {
		t.Errorf("expected permalink notes to render markdown link, got: %s", renderedPermalink)
	}

	// 2. Test index / card template
	var indexBuf strings.Builder
	err = set.IndexTemplate.ExecuteTemplate(&indexBuf, "base", map[string]any{
		"Recent": []linkding.Bookmark{bm},
		"Total":  1,
	})
	if err != nil {
		t.Fatalf("IndexTemplate execution failed: %v", err)
	}
	renderedIndex := indexBuf.String()

	if !strings.Contains(renderedIndex, `<div class="bookmark-notes">`) {
		t.Errorf("expected index card to contain .bookmark-notes div")
	}
	if !strings.Contains(renderedIndex, `<strong>bold</strong>`) {
		t.Errorf("expected index card notes to render markdown strong tag, got: %s", renderedIndex)
	}
	if !strings.Contains(renderedIndex, `<a href="https://example.com/target">a link</a>`) {
		t.Errorf("expected index card notes to render markdown link, got: %s", renderedIndex)
	}
}
