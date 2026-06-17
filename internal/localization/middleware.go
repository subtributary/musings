package localization

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
)

type localeKey struct{}

// LocaleFromContext returns the locale set by the LocalizedRoute middleware.
// If localization is disabled, UndLocale is returned.
func LocaleFromContext(ctx context.Context) Locale {
	if value, ok := ctx.Value(localeKey{}).(Locale); ok {
		return value
	}
	return UndLocale
}

func withLocale(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, localeKey{}, locale)
}

// LocalizedRoute enforces localized routes only if localization is enabled.
// See LocalizedRouteMiddleware for details of localized route handling.
func LocalizedRoute(locales []Locale) func(next http.Handler) http.Handler {
	// If localization is disabled, do nothing.
	if len(locales) == 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	m := NewLocalizedRouteMiddleware(locales)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r, ok := m.Handle(w, r)
			if !ok {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LocalizedRouteMiddleware enforces URLs that have the locale as the first
// path segment. The locale can be read via LocaleFromContext, and the chi
// context's RoutePath is updated to be the trailing path prefixed with "/".
//
// If the URL is not localized, then the response is a redirect to a localized
// URL that has a configured locale that is best suited per the request.
//
// If the path begins with "/_", then the middleware has no effect.
type LocalizedRouteMiddleware struct {
	locales []Locale
	matcher LocaleMatcher
}

func NewLocalizedRouteMiddleware(locales []Locale) *LocalizedRouteMiddleware {
	return &LocalizedRouteMiddleware{
		locales: locales,
		matcher: NewLocaleMatcher(locales),
	}
}

// Handle processes and acts on the request. It returns and updated request and
// whether the next middleware should be called.
func (m *LocalizedRouteMiddleware) Handle(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	chiContext := chi.RouteContext(r.Context())

	// Reserved and system paths are not localized.
	if strings.HasPrefix(r.URL.Path, "/_") {
		return r, true
	}

	// Get the locale the user wants and the trailing path.
	locale, trailing := ExtractLocale(m.locales, r.URL.Path)
	if locale.Tag == Und {
		accept := r.Header.Get("Accept-Language")
		locale = m.matcher.Choose(accept)
	}

	// Ensure the URL is the canonical localized URL.
	locPath := "/" + path.Join(locale.Tag, trailing)
	if r.URL.Path != locPath {
		redirectURL := *r.URL
		redirectURL.Path = locPath
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		return nil, false
	}

	// Everything is good, store locale in context and update path in context.
	r = r.WithContext(withLocale(r.Context(), locale))
	chiContext.RoutePath = trailing
	return r, true
}
