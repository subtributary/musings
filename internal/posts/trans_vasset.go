package posts

import (
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var nameKey = parser.NewContextKey()

func getName(pc parser.Context) (string, bool) {
	name, ok := pc.Get(nameKey).(string)
	return name, ok
}

func setName(pc parser.Context, name string) {
	pc.Set(nameKey, name)
}

// versionAssetsTransformer appends a version string to the query string of any
// link or image that points to a local asset, if their modified time is set.
type versionAssetsTransformer struct {
	modTime ModTimeFunc
}

func (t *versionAssetsTransformer) Transform(doc *ast.Document, _ text.Reader, pc parser.Context) {
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Link:
			node.Destination = t.versionAssetURL(pc, node.Destination)
		case *ast.Image:
			node.Destination = t.versionAssetURL(pc, node.Destination)
		}

		return ast.WalkContinue, nil
	})

	// We should never get an error since we never return one in the walker.
	if err != nil {
		log.Fatalf("Unexpected error walking AST: %v", err)
	}
}

// versionAssetURL appends a version timestamp if the modified time is set.
func (t *versionAssetsTransformer) versionAssetURL(pc parser.Context, dest []byte) []byte {
	assetURL, err := url.Parse(string(dest))
	if err != nil {
		return dest // Not a URL we can work with; leave it alone.
	}
	if assetURL.IsAbs() || assetURL.Host != "" {
		return dest
	}

	postPath, ok := getName(pc)
	if !ok {
		return dest
	}

	assetPath := t.resolveAssetPath(postPath, assetURL)
	when, ok := t.modTime(assetPath)
	if !ok {
		return dest
	}

	query := assetURL.Query()
	query.Set("v", strconv.FormatInt(when.Unix(), 16))
	assetURL.RawQuery = query.Encode()

	return []byte(assetURL.String())
}

func (t *versionAssetsTransformer) resolveAssetPath(postPath string, assetURL *url.URL) string {
	if strings.HasPrefix(assetURL.Path, "/") {
		return assetURL.Path
	}

	return "/" + path.Join(path.Dir(postPath), assetURL.Path)
}
