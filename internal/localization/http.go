package localization

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

type localeKey struct{}

func LocaleFromContext(ctx context.Context) (language.Tag, bool) {
	tag, ok := ctx.Value(localeKey{}).(language.Tag)
	return tag, ok
}

func withLocale(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, localeKey{}, tag)
}

// LocalizedRoute enforces the locale in the URL unless the supported locales
// include [language.Und]. In any case, it also sets the discovered locale in
// the context.
func LocalizedRoute(tags []language.Tag) func(next http.Handler) http.Handler {
	if len(tags) == 0 {
		panic("tags must not be empty")
	}

	matcher := language.NewMatcher(tags)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chiContext := chi.RouteContext(r.Context())
			reqPath := chiContext.RoutePath
			if reqPath == "" {
				reqPath = r.URL.Path
			}

			tag, trailing, err := ParsePath(reqPath)

			if err == nil && slices.Contains(tags, tag) {
				// We have a localized path.
				r = r.WithContext(withLocale(r.Context(), tag))
				chiContext.RoutePath = trailing
				next.ServeHTTP(w, r)

			} else if slices.Contains(tags, language.Und) {
				// We do not have a localized path, and we do not care.
				r = r.WithContext(withLocale(r.Context(), language.Und))
				next.ServeHTTP(w, r)

			} else {
				// We do not have a localized path, so let's fix that.
				lang := r.Header.Get("Accept-Language")
				tag, _ = language.MatchStrings(matcher, lang)
				redirectURL := fmt.Sprintf("/%s%s", tag, reqPath)
				http.Redirect(w, r, redirectURL, http.StatusFound)
			}
		})
	}
}

// ParsePath parses the locale out of the first segment of a path.
// It returns the language tag and the remaining path after that.
// If the locale is invalid, then it returns an error.
//
// All paths should be prefixed with a slash.
func ParsePath(reqPath string) (language.Tag, string, error) {
	reqPath = path.Clean(reqPath)
	segments := strings.SplitN(reqPath, "/", 3)

	tag, err := language.Parse(segments[1])
	if err != nil {
		return language.Und, "", err
	}

	if len(segments) == 3 {
		return tag, "/" + segments[2], nil
	}
	return tag, "/", nil
}
