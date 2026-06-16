# Customizing

The frontend of Musings can be customized, but it is a technical process not meant for the end user.
If end user customization is required, a user-friendly tool complete with sanitization should be used.
This document details the technical process of customizing Musings.

## Customization areas

Customization is performed by modifying files within the repository or deployment.
The following areas are intended to be customized:

- Styles ("web/src/scss/" and "web/static/css/")
- Scripts ("web/src/ts/" and "web/static/js/")
- Other static assets ("web/static/\<type\>")
- Templates ("web/templates/")

Other areas of the codebase should be considered internal and should not be edited per deployment.

## Editing styles and scripts

Styles and scripts in source code are in "web/src/scss/" and "web/src/ts/".
These are transpiled to "web/static/css/" and "web/static/js/", respectively.
These are served at "/_static/css/" and "/_static/js/", respectively.
`npm` is used to build these:

```bash
npm install # needed on first build

npm run build
# -or-
npm run build:watch # to rebuild after changes
```

Styles and scripts in a deployed Musings website are in "web/static/css/" and "web/static/js/".
A file here may be overwritten if it shares a name with a file in "web/src/scss/" (for styles) or "web/src/ts/" (for scripts).
(A solution for this is proposed in [Musings issue #43](https://github.com/subtributary/musings/issues/43).)


## Adding media and other assets

Additional assets *to be used for customization* should be placed in a subdirectory of "web/static/".
The subdirectory should be named after the asset type: "fonts", "images", etc.

Additional assets *to be used for customization* should be placed under "web/static/".
Assets used by content should not be placed here
because they are managed separately from customization assets.
(See [content.md](content.md) for where to put content media.)

The subdirectories are served at "/_static/\<name\>/".


## Editing the templates

Templates are Go Templates with a ".gohtml" file extension in the "/web/templates" directory.

The main templates are:

* index.gohtml
* post.gohtml

All templates start with this view model:

```go

type LocaleOption struct {
    Tag         string
    NativeName  string
    Direction   string
    WritingMode string
    IsCurrent   bool
    URL         string
}

type ViewModel struct {
    LocaleOption []LocaleOption
    Language     LocaleOption
    Translations map[string]string
    Data         []any
}
```

The structure of the `Data` property depends on the view.

###  index.gohtml

The "index.gohtml" file is the template for the website index, or home page.
It is served at "/" (or "/q=" when searching).

The view model's `Data` property is a `SearchResult` with this structure:

```go
type IndexedPost struct {
    Path      string
    Bylines   []string
    Published *time.Time
    Summary   string
    Title     string
}

type SearchResults struct {
    Query   string
    Results []posts.IndexedPost
}
```

When not responding to a search,
the posts in the view model are sorted by publication date with the most recent one first.
Posts with no publication date are listed first, sorted by title then path.
Posts with a future publication date are omitted.

When responding to a search,
the posts in the view model are sorted by how well they match the search query.
Non-matches and posts with a future publication date are omitted.

### post.gohtml

The "post.gohtml" file is the template for posts in "content/".
(See [content.md](content.md) for how to create the content.)

Posts are served at routes that match the structure of the "content" directory.
For example, "content/en/hello.md" is served at "/en/hello.md".

The view model's `Data` property is a `ParsedPost` with this structure:

```go
type ParsedPost struct {
    Bylines   []string
    Content   template.HTML
    Published time.Time
    Summary   string
    Title     string
}
```

### partials

Partial templates to support the main templates are in the "templates/partials" directory.
A partial's name is its file name without the extension and with "partials/" prefixed;
for example, "/templates/partials/head.gohtml" is named "partials/head".
The name is used to include the partial in another template.

The view model for a partial is whatever is passed to it from another template.
