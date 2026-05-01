package posts

import (
	"slices"
	"strings"
	"testing"
	"time"
)

func TestParseContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    ParsedPost
	}{
		{
			name:    "empty",
			content: "",
			want:    ParsedPost{},
		},
		{
			name: "full frontmatter",
			content: "---\n" +
				"bylines: [by Nathan Belue]\n" +
				"published: 2026-04-30\n" +
				"tags: [apple, banana]\n" +
				"---\n",
			want: ParsedPost{
				Bylines:   []string{"by Nathan Belue"},
				Published: parseTime(time.DateOnly, "2026-04-30"),
				Tags:      []string{"apple", "banana"},
			},
		},
		{
			name: "frontmatter cannot overwrite title or content",
			content: "---\ntitle: Apple\ncontent: Banana\n---\n" +
				"# Cucumber\nDate",
			want: ParsedPost{Title: "Cucumber", Content: "<p>Date</p>"},
		},
		{
			name:    "invalid frontmatter",
			content: "--\ntags: [apple\n---",
			wantErr: true,
		},
		{
			name:    "move heading to title",
			content: "# Hello world",
			want:    ParsedPost{Title: "Hello world", Content: ""},
		},
		{
			name:    "formatted heading",
			content: "# Hello *world*",
			want:    ParsedPost{Title: "Hello world", Content: ""},
		},
		{
			name:    "only first top-level h1 becomes title",
			content: "# Title\n\n# Section",
			want:    ParsedPost{Title: "Title", Content: "<h1>Section</h1>"},
		},
		{
			name:    "blockquote h1 is content not title",
			content: "> # Quoted heading",
			want: ParsedPost{
				Title:   "",
				Content: "<blockquote>\n<h1>Quoted heading</h1>\n</blockquote>",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser := NewParser()
			post, err := parser.ParseContent([]byte(tt.content))

			if tt.wantErr && err == nil {
				t.Fatalf("ParseContent(%q): expected error", tt.content)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ParseContent(%q): unexpected error: %v", tt.content, err)
			}

			post.Content = strings.TrimSpace(post.Content)
			if post.Content != tt.want.Content {
				t.Errorf("post.Content: got %v, want %v", post.Content, tt.want.Content)
			}

			// Metadata
			if !slices.Equal(post.Bylines, tt.want.Bylines) {
				t.Errorf("post.Bylines: got %v, want %v", post.Bylines, tt.want.Bylines)
			}
			if post.Published != tt.want.Published {
				t.Errorf("post.Published: got %v, want %v", post.Published, tt.want.Published)
			}
			if !slices.Equal(post.Tags, tt.want.Tags) {
				t.Errorf("post.Tags: got %v, want %v", post.Tags, tt.want.Tags)
			}
			if post.Title != tt.want.Title {
				t.Errorf("post.Title: got %v, want %v", post.Title, tt.want.Title)
			}
		})
	}
}

func TestParseContent_FrontmatterDates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  *time.Time
	}{
		{
			name:  "empty",
			input: "",
			want:  nil,
		},
		{
			name:  "invalid time",
			input: "apple",
			want:  nil,
		},
		{
			name:  "unsupported format",
			input: "2026-04-30T12:00",
			want:  nil,
		},
		{
			name:  "support yyyy-MM-dd",
			input: "2026-04-30",
			want:  parseTime(time.DateOnly, "2026-04-30"),
		},
		{
			name:  "support yyyy-MM-dd HH:mm:ss",
			input: "2026-04-30 12:20:36",
			want:  parseTime(time.DateTime, "2026-04-30 12:20:36"),
		},
		{
			name:  "support RFC 3339",
			input: "2026-04-30T12:20:36-05:00",
			want:  parseTime(time.RFC3339, "2026-04-30T12:20:36-05:00"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			content := "---\npublished: " + tt.input + "\n---"
			parser := NewParser()
			post, err := parser.ParseContent([]byte(content))

			if err != nil {
				t.Fatalf("ParseContent(%q): unexpected error: %v", content, err)
			}

			if tt.want == nil && post.Published != nil {
				t.Errorf("post.Published: got %v, want nil", post.Published)
			} else if tt.want != nil && post.Published == nil {
				t.Errorf("post.Published: got %v, want nil", post.Published)
			} else if tt.want != nil && *post.Published != *tt.want {
				t.Errorf("post.Published: got %v, want %v", post.Published, tt.want)
			}
		})
	}
}

func parseTime(format string, value string) *time.Time {
	if t, err := time.Parse(format, value); err == nil {
		return &t
	}
	return nil
}
