package files

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"

	"golang.org/x/text/language"
)

type wantEntry struct {
	name  string
	isDir bool
}

func TestOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fs           fstest.MapFS
		tag          language.Tag
		expectedData []byte
	}{
		{
			name: "localization is optional",
			fs: fstest.MapFS{
				"name.md": {Data: []byte("regular")},
			},
			tag:          language.Und,
			expectedData: []byte("regular"),
		},
		{
			name: "uses localized variant",
			fs: fstest.MapFS{
				"name.md":    {Data: []byte("original")},
				"name.en.md": {Data: []byte("English")},
			},
			tag:          language.English,
			expectedData: []byte("English"),
		},
		{
			name: "uses parent locale when no exact match",
			fs: fstest.MapFS{
				"name.md":    {Data: []byte("original")},
				"name.en.md": {Data: []byte("English")},
			},
			tag:          language.AmericanEnglish,
			expectedData: []byte("English"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := NewLocalizedFS(tt.fs, tt.tag)
			content, err := fs.ReadFile(subject, "name.md")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(content, tt.expectedData) {
				t.Errorf("got %s, want %s", content, tt.expectedData)
			}
		})
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fs   fstest.MapFS
		tags []language.Tag
		want map[language.Tag][]wantEntry
	}{
		{
			name: "localization is optional",
			fs:   fstest.MapFS{"regular.md": {}},
			tags: []language.Tag{language.Und},
			want: map[language.Tag][]wantEntry{
				language.Und: {
					{name: "regular.md", isDir: false},
				},
			},
		},
		{
			name: "directories cannot be localized",
			fs:   fstest.MapFS{"dir.en.ext/": {}},
			tags: []language.Tag{language.English},
			want: map[language.Tag][]wantEntry{
				language.English: {
					{name: "dir.en.ext", isDir: true},
				},
			},
		},
		{
			name: "files can be localized",
			fs: fstest.MapFS{
				"en-only.en.md": {},
				"hello.md":      {},
				"hello.en.md":   {},
				"hello.ko.md":   {},
			},
			tags: []language.Tag{
				language.English,
				language.French,
				language.Korean,
			},
			want: map[language.Tag][]wantEntry{
				language.English: {
					{name: "en-only.md", isDir: false},
					{name: "hello.md", isDir: false},
				},
				language.French: {
					{name: "hello.md", isDir: false},
				},
				language.Korean: {
					{name: "hello.md", isDir: false},
				},
			},
		},
		{
			name: "unconfigured locales are ignored",
			fs:   fstest.MapFS{"ignored.fr.md": {}},
			tags: []language.Tag{language.English},
			want: map[language.Tag][]wantEntry{
				language.English: {},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			groupedByLocale, err := Scan(tt.fs, tt.tags)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for tag, wantEntries := range tt.want {
				gotEntries, ok := groupedByLocale[tag]
				if !ok {
					t.Fatalf("locale group %q does not exist", tag)
				}

				if len(gotEntries) != len(wantEntries) {
					t.Fatalf("locale %q: got %d files, expected %d", tag, len(gotEntries), len(wantEntries))
				}

				for i, wanted := range wantEntries {
					got := gotEntries[i]

					if got.Name() != wanted.name {
						t.Errorf("locale %q: got name %q, expected %q", tag, got.Name(), wanted.name)
					}
					if got.IsDir() != wanted.isDir {
						t.Errorf("locale %q: got isDir=%v, expected %v", tag, got.IsDir(), wanted.isDir)
					}
				}
			}
		})
	}
}
