package web

import (
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

// VersionURL returns adds cache-busting to a relative URL.
func VersionURL(root *os.Root, currentPath, targetPath string) string {
	// Check if the target is a relative URL before doing anything else.
	parsedURL, err := url.Parse(targetPath)
	if err != nil || parsedURL.IsAbs() || parsedURL.Host != "" {
		return targetPath
	}

	filePath := path.Join(currentPath, targetPath)
	filePath, _ = strings.CutPrefix(filePath, "/")

	info, err := root.Stat(filePath)
	if err != nil {
		return targetPath
	}
	modTime := info.ModTime()

	query := parsedURL.Query()
	query.Set("v", strconv.FormatInt(modTime.Unix(), 16))
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
