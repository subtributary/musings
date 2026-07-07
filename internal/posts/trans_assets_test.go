package posts_test

import (
	"strings"
	"testing"

	"github.com/subtributary/musings/internal/posts"
)

func TestVersionAssetsTransformer(t *testing.T) {
	t.Parallel()

	const version = "--version--"
	versionURL := func(_, name string) string {
		return name + "?v=" + version
	}

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
			name:    "image with no path",
			content: "![](asset.png)",
			want:    `<p><img src="asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "image with relative path",
			content: "![](../asset.png)",
			want:    `<p><img src="../asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "image with absolute path",
			content: "![](/asset.png)",
			want:    `<p><img src="/asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "link with no path",
			content: "[file](file.txt)",
			want:    `<p><a href="file.txt?v=` + version + `">file</a></p>`,
		},
		{
			name:    "link with relative path",
			content: "[file](../file.txt)",
			want:    `<p><a href="../file.txt?v=` + version + `">file</a></p>`,
		},
		{
			name:    "link with absolute path",
			content: "[file](/file.txt)",
			want:    `<p><a href="/file.txt?v=` + version + `">file</a></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser := posts.NewParser(versionURL)
			post, err := parser.ParseContent("", []byte(tt.content))
			if err != nil {
				t.Fatalf("error parsing content: %v", err)
			}

			got := strings.TrimSpace(string(post.Content))
			if got != tt.want {
				t.Errorf("post.Content = %v, want %v", got, tt.want)
			}
		})
	}
}
