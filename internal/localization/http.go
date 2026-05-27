package localization

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

type localeKey struct{}

func LocaleFromContext(ctx context.Context) language.Tag {
	if value, ok := ctx.Value(localeKey{}).(language.Tag); ok {
		return value
	}
	panic("context is not set; LocalizedRoute middleware needs to be used")
}

func withLocale(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, localeKey{}, tag)
}

// LocalizedRoute enforces the locale in the URL unless the supported locales
// include [language.Und]. In any case, it also sets the discovered locale in
// the context. Paths starting with "/_" are not localized or redirected.
func LocalizedRoute(tags []language.Tag) func(next http.Handler) http.Handler {
	if len(tags) == 0 {
		// Tags are fixed at startup. But just in case, fail quick:
		panic("tags must contain at least one tag")
	}

	matcher := language.NewMatcher(tags)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chiContext := chi.RouteContext(r.Context())
			reqPath := r.URL.Path

			if !strings.HasPrefix(reqPath, "/") {
				log.Printf("expected slash prefix on path, got %q", reqPath)
				reqPath = "/" + reqPath
			}

			// Reserved and system paths are not localized.
			if strings.HasPrefix(reqPath, "/_") {
				r = r.WithContext(withLocale(r.Context(), language.Und))
				chiContext.RoutePath = reqPath
				next.ServeHTTP(w, r)
				return
			}

			// Get the localized path info; (Und, reqPath) if not localized.
			tag, trailing := ParsePath(reqPath)

			// Do not redirect if localized path or localization disabled.
			if slices.Contains(tags, tag) {
				r = r.WithContext(withLocale(r.Context(), tag))
				chiContext.RoutePath = trailing
				next.ServeHTTP(w, r)
				return
			}

			// `tag` is invalid, so replace it with a good match.
			lang := r.Header.Get("Accept-Language")
			tag, i := language.MatchStrings(matcher, lang)
			tag = tags[i] // See <https://github.com/golang/go/issues/24211>.

			// Redirect the user to the localized path we found for them.
			redirectURL := *r.URL
			redirectURL.Path = "/" + tag.String() + reqPath
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		})
	}
}

// ParsePath parses the locale out of the first segment of a path.
// It returns the language tag and the remaining path after that.
// If the locale is invalid or missing, then it returns (Und, reqPath).
func ParsePath(reqPath string) (language.Tag, string) {
	reqPath = strings.TrimLeft(reqPath, "/")

	// Parse `loc/trailing/etc` into `["loc", "trailing/etc"]`.
	segments := strings.SplitN(reqPath, "/", 2)

	tag, err := language.Parse(segments[0])
	if err != nil {
		return language.Und, "/" + reqPath
	}

	if len(segments) > 1 {
		return tag, "/" + segments[1]
	}
	return tag, "/"
}
