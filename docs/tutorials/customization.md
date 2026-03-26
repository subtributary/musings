# Customization

Musings is designed to be customized by a software developer.
It is recommended to *not* give non-technical users access to the customizable files.
If end user customization without a developer middleman is desired, then a sanitizing tool should be used.

The logistics of maintaining a customization of Musings is left to the developer.
One method is to create a fork of Musings source code, then periodically sync from `upstream`.
If edits are limited to areas intented to be customized, then merge conflicts will be minimized.

## Customization areas

Customization is performed by modifying files within the repository or deployment.
The following areas are intended to be customized:

- Styles and scripts ("/web/src/...")
- Templates ("/web/templates/...")
- Media ("/web/static/media/...")

Other areas of the codebase should be considered internal and may change without notice.

## Editing styles and scripts

Styles and scripts are written in SCSS and TypeScript, respectively.
The source code is located in "/web/src/", and the transpiled files are located in "./web/static/".
The files will need to be built initially and after every change.

To build, first install all needed packages.

```bash
npm install
```

Then build the files.

```bash
npm run build
# -or-
npm run build:watch
```

Transpiled files are served at the base URLs "/_static/css/" and "/_static/js/".

Transpiled files follow our [l10n](localization.md) rules, so localized variants can be created.

## Editing the templates

Templates are Go Templates with a ".gohtml" file extension..

Templates, including partials, follow our [l10n](localization.md) rules, so localized variants can be created.

###  index.gohtml

The "index.gohtml" file is the template for the website index, or landing page.

The index is served at "/".

There is (currently) no view model for this template.

### post.gohtml

The "post.gohtml" file is the template for posts in "/content/".
(See [content.md](content.md) for how to create the content.)

The posts are served at routes that match the directory structure of "/content/".

The view model for this template is a `PostData` structure:

```go
type PostData struct {
  HtmlContent template.HTML
}
```

### Partials

Partial templates are included from other templates to help reduce code duplication.

To include a partial template, use "partials/filename" without its locale and extension as its name.
For example, to include "/web/templates/partials/head.en.gohtml", use `{{template "partials/head"}}`.

A partial template does not have a view model unless one is passed from its caller.

## Adding media

Additional assets *to be used for customization* should be placed under "/web/static/media/".
Assets used by content should not be placed here.
(See [content.md](content.md).)

The media directory is served at "/_static/media/".
