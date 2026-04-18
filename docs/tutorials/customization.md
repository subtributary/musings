# Customization

Musings is designed to be customized by a software developer.
Non-technical users should not be given direct access to customizable files.
If end user customization without a developer middleman is desired, then a sanitizing tool should be used.

The logistics of maintaining a customization of Musings is left to the developer.
One method is to create a fork of Musings source code, then periodically sync from `upstream`.
If edits are limited to intended customization areas, then merge conflicts will be minimized.

## Customization areas

Customization is performed by modifying files within the repository or deployment.
The following areas are intended to be customized:

 - Styles ("/web/src/scss/" and "/web/static/css")
 - Scripts ("/web/src/ts/" and "/web/static/js")
 - Templates ("/web/templates/...")
 - Media ("/web/static/media/...")

Other areas of the codebase should be considered internal and may change without notice.
Customizing these areas may cause breakages when upgrading Musings.

## Editing styles and scripts

### Within source code

If you are editing a published Musings website,
then the styles will be in "/web/static/css",
and the scripts will be in "/web/static/js".

If you are editing Musings source code,
then SCSS styles can be placed in "/web/src/scss",
and TypeScript scripts can be placed in "/web/src/ts".

Changes to "/web/src" will not be reflected until the assets are built.

```bash
npm install # needed on first build

npm run build
# -or-
npm run build:watch # to rebuild after changes
```

The files are served at the base URLs "/_static/css/" and "/_static/js/".

## Editing the templates

Templates are Go Templates with a ".gohtml" file extension.

The base view model for every template is `ViewModel`.
Individual templates may expand this via its `Data` property.
The structures of the base view model and its dependencies are below:

```go
type LanguageOption struct {
    Code      string
    Label     string
    IsCurrent bool
    URL       string
}

type Strings struct {
    SelectLanguage string
    WebsiteLogo    string
}

type ViewModel struct {
    LanguageOptions []LanguageOption
    Language        LanguageOption
    Translations    Strings
    Data            any
}
```

###  index.gohtml

The "index.gohtml" file is the template for the website index, or landing page.

The index is served at "/".

### post.gohtml

The "post.gohtml" file is the template for posts in "/content/".
(See [content.md](content.md) for how to create the content.)

The posts are served at routes that match the directory structure of "/content/".

The view model data for this template is a `PostData` structure:

```go
type PostData struct {
HtmlContent template.HTML
}
```

### Partials

Partial templates are included from other templates to help reduce code duplication.

To include a partial template, use its logical name (without [locale](localization.md) or extension).
For example, to include "/web/templates/partials/head.en.gohtml", use `{{template "partials/head"}}`.

A partial template does not have a view model unless one is passed from its caller.

## Adding media

Additional assets *to be used for customization* should be placed under "/web/static/media/".
Assets used by content should not be placed here,
as they are managed separately from customization assets.
(See [content.md](content.md) for where to put content media.)

The media directory is served at "/_static/media/".
