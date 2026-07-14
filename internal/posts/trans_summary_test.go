package posts_test

import (
	"strings"
	"testing"

	"github.com/subtributary/musings/internal/posts"
)

func TestSummaryTransformer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "empty",
			content: "",
			want:    "",
		},
		{
			name:    "one paragraph",
			content: "Hello, world!",
			want:    "Hello, world!",
		},
		{
			name:    "two paragraphs",
			content: "Hello, world!\n\nGoodbye now.",
			want:    "Hello, world!",
		},
		{
			name:    "formatted",
			content: "[Hello](hello.md), _world_! **How** are *you*?",
			want:    "Hello, world! How are you?",
		},
		{
			name:    "with header",
			content: "# Title\n\nHello, world!",
			want:    "Hello, world!",
		},
		{
			name:    "with frontmatter",
			content: "---\nPublished: 2026-06-16\n---\nHello, world!",
			want:    "Hello, world!",
		},
		{
			name:    "with frontmatter and title",
			content: "---\nPublished: 2026-06-16\n---\n# Title\n\nHello!",
			want:    "Hello!",
		},
		{
			name:    "image without alt text then paragraph",
			content: "![](img.png)\n\nHello, world!",
			want:    "Hello, world!",
		},
		{
			name:    "image with alt text then paragraph",
			content: "![Alt text](img.png)\n\nHello, world!",
			want:    "Alt text",
		},
	}

	versionURL := func(_, name string) string { return name }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser := posts.NewParser(versionURL)
			post, err := parser.ParseContent("", []byte(tt.content))
			if err != nil {
				t.Fatalf("error parsing content: %v", err)
			}

			got := strings.TrimSpace(string(post.Summary))
			if got != tt.want {
				t.Errorf("post.Content = %s, want %s", got, tt.want)
			}
		})
	}

}
