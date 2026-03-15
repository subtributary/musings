package templates

import (
	"fmt"
	"strings"
)

// PagesPath is the path of page templates relative to the templates directory.
const PagesPath = "."

// PartialsPath is the path of partial templates relative to the templates directory.
const PartialsPath = "partials"

// ExtractNameAndLocale extracts the name and locale from a filename per our conventions.
// The filename is assumed to be like "{name}.{locale}.gohtml" or "{name}.gohtml".
//
// The name may contain periods only if the locale is also specified.
//
// If the locale is specified, it is normalized to lowercase; otherwise, it is set to "".
func ExtractNameAndLocale(filename string) (string, string) {
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 {
		panic(fmt.Sprintf("invalid filename: %s", filename))
	}

	penultimateDot := strings.LastIndex(filename[:lastDot], ".")
	if penultimateDot != -1 {
		locale := filename[penultimateDot+1 : lastDot]
		return filename[:penultimateDot], strings.ToLower(locale)
	}

	return filename[:lastDot], ""
}
