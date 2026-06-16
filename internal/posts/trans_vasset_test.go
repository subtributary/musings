package posts_test

import (
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/subtributary/musings/internal/posts"
)

func TestVersionAssetsTransformer(t *testing.T) {
	t.Parallel()

	modified := time.Date(2026, 6, 15, 19, 0, 0, 0, time.UTC)
	version := strconv.FormatInt(modified.Unix(), 16)

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
			name:    "local image, no path",
			content: "![](asset.png)",
			want:    `<p><img src="asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "local image, relative path",
			content: "![](../asset.png)",
			want:    `<p><img src="../asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "local image, root path",
			content: "![](/asset.png)",
			want:    `<p><img src="/asset.png?v=` + version + `" alt=""></p>`,
		},
		{
			name:    "external image",
			content: "![](https://example.com/asset.png)",
			want:    `<p><img src="https://example.com/asset.png" alt=""></p>`,
		},
		{
			name:    "local link, no path",
			content: "[file](file.txt)",
			want:    `<p><a href="file.txt?v=` + version + `">file</a></p>`,
		},
		{
			name:    "local link, relative path",
			content: "[file](../file.txt)",
			want:    `<p><a href="../file.txt?v=` + version + `">file</a></p>`,
		},
		{
			name:    "local link, root path",
			content: "[file](/file.txt)",
			want:    `<p><a href="/file.txt?v=` + version + `">file</a></p>`,
		},
		{
			name:    "external link",
			content: "[file](https://example.com/file.txt)",
			want:    `<p><a href="https://example.com/file.txt">file</a></p>`,
		},
		{
			name:    "link to post, no path",
			content: "[post](post)",
			want:    `<p><a href="post">post</a></p>`,
		},
		{
			name:    "link to post, relative path",
			content: "[post](../post)",
			want:    `<p><a href="../post">post</a></p>`,
		},
		{
			name:    "link to post, root path",
			content: "[post](/post)",
			want:    `<p><a href="/post">post</a></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// The server returns (_, false) for post mod times,
			// so we will do the same here for testing.
			modTime := func(name string) (time.Time, bool) {
				if path.Ext(name) == "" {
					return modified, false
				}
				return modified, true
			}

			parser := posts.NewParser(modTime)
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
