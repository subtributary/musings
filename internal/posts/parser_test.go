package posts_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/subtributary/musings/internal/posts"
)

func TestParseContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    posts.ParsedPost
		wantErr bool
	}{
		{
			name: "empty",
			want: posts.ParsedPost{},
		},
		{
			name:    "full frontmatter",
			content: "---\nbylines: [by Nathan]\npublished: 2026-04-30\n---\n",
			want: posts.ParsedPost{
				Bylines:   []string{"by Nathan"},
				Published: parseTime(time.DateOnly, "2026-04-30"),
			},
		},
		{
			name:    "frontmatter with invalid syntax",
			content: "---\ntags: [apple\n---",
			wantErr: true,
		},
		{
			name:    "frontmatter with invalid time",
			content: "---\npublished: apple\n---\n",
			wantErr: true,
		},
		{
			name:    "frontmatter with unsupported time",
			content: "---\npublished: 2026-04-30T12:00\n---\n",
			wantErr: true,
		},
		{
			name:    "frontmatter with yyyy-MM-dd time",
			content: "---\npublished: 2026-04-30\n---\n",
			want: posts.ParsedPost{
				Published: parseTime(time.DateOnly, "2026-04-30"),
			},
		},
		{
			name:    "frontmatter with yyyy-MM-dd HH:mm:ss time",
			content: "---\npublished: 2026-04-30 12:20:36\n---\n",
			want: posts.ParsedPost{
				Published: parseTime(time.DateTime, "2026-04-30 12:20:36"),
			},
		},
		{
			name:    "frontmatter with RFC 3339 time",
			content: "---\npublished: 2026-04-30T12:20:36-05:00\n---\n",
			want: posts.ParsedPost{
				Published: parseTime(time.RFC3339, "2026-04-30T12:20:36-05:00"),
			},
		},
		{
			name:    "frontmatter does not overwrite title or content",
			content: "---\ntitle: Apple\ncontent: Banana\n---\n# Carot\nDate",
			want: posts.ParsedPost{
				Title:   "Carot",
				Content: "<p>Date</p>",
			},
		},
	}

	modTime := func(string) (time.Time, bool) { return time.Time{}, false }

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser := posts.NewParser(modTime)
			post, err := parser.ParseContent("", []byte(tt.content))
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("failed to parse content: %v", err)
			}

			gotContent := strings.TrimSpace(string(post.Content))
			wantContent := string(tt.want.Content)
			if gotContent != wantContent {
				t.Errorf("post.Content: got %v, want %v", post.Content, tt.want.Content)
			}

			if !slices.Equal(post.Bylines, tt.want.Bylines) {
				t.Errorf("post.Bylines: got %v, want %v", post.Bylines, tt.want.Bylines)
			}

			if !post.Published.Equal(tt.want.Published) {
				t.Errorf("post.Published = %v, want %v", post.Published, tt.want.Published)
			}

			// Moving the heading to the title is tested in "trans_h1_test.go".
			// Here we only care that our other features didn't mess it up.
			if post.Title != tt.want.Title {
				t.Errorf("post.Title: got %v, want %v", post.Title, tt.want.Title)
			}
		})
	}
}

func parseTime(format string, value string) time.Time {
	if t, err := time.Parse(format, value); err == nil {
		return t
	}
	return time.Time{}
}
