package posts

import (
	"path"

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

// VersionURLFunc returns a cache-breaking version of a local URL.
type VersionURLFunc func(currentPath, target string) string

// versionAssetsTransformer appends a version string to the query string of any
// link or image that points to a local asset, if their modified time is set.
type versionAssetsTransformer struct {
	versionURL VersionURLFunc
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
		panic("Unexpected: error walking AST:" + err.Error())
	}
}

// versionAssetURL appends a version timestamp if the modified time is set.
func (t *versionAssetsTransformer) versionAssetURL(pc parser.Context, dest []byte) []byte {
	targetURL := string(dest)

	postPath, _ := getName(pc)
	postPath = path.Dir(postPath)

	return []byte(t.versionURL(postPath, targetURL))
}
