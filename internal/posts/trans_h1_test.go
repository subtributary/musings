package posts_test

import (
	"strings"
	"testing"
	"time"

	"github.com/subtributary/musings/internal/posts"
)

func TestRemoveH1Transformer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    posts.ParsedPost
	}{
		{
			name: "empty",
			want: posts.ParsedPost{
				Title:   "",
				Content: "",
			},
		},
		{
			name:    "heading",
			content: "# Hello",
			want: posts.ParsedPost{
				Title:   "Hello",
				Content: "",
			},
		},
		{
			name:    "formatted heading",
			content: "# Hello *world*",
			want: posts.ParsedPost{
				Title:   "Hello world",
				Content: "",
			},
		},
		{
			name:    "two headings",
			content: "# Title\n\n# Section",
			want: posts.ParsedPost{
				Title:   "Title",
				Content: "<h1>Section</h1>",
			},
		},
		{
			name:    "blockquote heading",
			content: "> # Quoted heading",
			want: posts.ParsedPost{
				Title:   "",
				Content: "<blockquote>\n<h1>Quoted heading</h1>\n</blockquote>",
			},
		},
	}

	modTime := func(string) (time.Time, bool) { return time.Time{}, false }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser := posts.NewParser(modTime)
			post, err := parser.ParseContent("", []byte(tt.content))
			if err != nil {
				t.Fatalf("error parsing content: %v", err)
			}

			gotContent := strings.TrimSpace(string(post.Content))
			wantContent := string(tt.want.Content)
			if gotContent != wantContent {
				t.Errorf("post.Content = %v, want %v", gotContent, wantContent)
			}

			if post.Title != tt.want.Title {
				t.Errorf("post.Title = %v, want %v", post.Title, tt.want.Title)
			}
		})
	}
}
