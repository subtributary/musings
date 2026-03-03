# Scope

Musings is a lightweight publishing CMS with first-class localization. It is designed for blog-like websites out of the box, but advanced users can customize it into other kinds of websites.

This document defines the scope so that Musings does not grow too complex. Additions to the scope must be carefully considered.


## Executables

These scripts and binaries will be included in the compiled project:

 - bin/api      - API for interfacing programatically.
 - bin/admin    - Content management system website.
 - bin/web      - Public website.
 - scripts/TBD  - {useful developer scripts}

Musings can be run simply with `all`, which combines `api` ("/api"), `admin` ("/admin"), and `web` ("/") into a monolith application. Or, the components can be run separately for more flexible hosting options.

## Features

### Alerts for Content Problems

Alerts can be attached to content (pages and posts) and will be presented in the UI. These alerts can be added by external tools.

A tool that checks that updates were made consistently across locales will be included.

### API Endpoints

An API is exposed for external tools to manage the content.

### Content Editor

The blog content of the website is stored in Markdown files that follow the CommonMark specification.

Documentation may recommend against using some rarely-implemented CommonMark features, but the CMS will not prevent the user from using them.

Content metadata (title, tags, and summary) can be edited too.

### Custom Routes

Page routes can be added to render a post at a non-post URL. For example, post `1` can be rendered at <http://example.com/en/custom> _instead_ of at <http://example.com/en/posts/1>.

These changes are made by editing "content/routes.json". The UI does not present a way to edit routes. This limitation ensures that only a system administrator can change the routes.

### Customizable Presentation

The architecture separates presentation from core logic, so customizing the website via code is easy.

The file structure is:
 - bin/         - binaries
 - core/        - source code for the core layers
    - {API used by `custom` code}
    - {internal files}
 - content/     - git-backed content
 - custom/      - customizable components
    - pages/
        - admin/
            - _Layout.??
            - _Head.{locale}.??
        - Index.??      - list of posts
        - Post.??       - framing for the post render
        - PostPage.??   - serves a post at a custom route
    - posts/
        - {id}/
            - metadata.json
            - {locale}.md
    - theme/
        - base.css
        - base.{locale}.css
        - base.js
    - config.json
 - scripts/ - useful scripts
 - .env
 - LICENSE.md

It is not recommended to modify anything outside of the "custom" folder.

### Git-backed Content

Content is stored in a local Git repository, and Git features may be exposed in user interfaces.

The repository structure is:
 - media/{year}/{month}/    - Ignored by default.
    - {id}.{size}.ext
 - posts/{id}/
    - metadata.json
    - {locale}.md
 - .gitignore
 - routes.json  - custom routes for posts

Unpublished content is on the `draft` branch, and published content is on the `main` branch.

### Image Manager

Editors may upload and browse images. Image metadata may be edited.

Uploaded images are resized to various widths for thumbnails and full-width renders.

### Localized URLs

Every URL (except for API URLs) includes the locale of the content. This cannot be disabled.

 * Admin URLs are of the form <http://example.com/_admin/{locale}/{page}>
 * Page URLs are of the form <http://example.com/{locale}/{page}>
 * Post URLs are of the form <http://example.com/{locale}/posts/{id}/{slug}>
 * Asset URLs are of the form <http://example.com/_include/{asset}.ext> or <http://example.com/_include/{asset}.{locale}.ext>

If the locale is omitted, then it is deduced from the `Accept-Language` header or chosen indiscriminately, and the browser is redirected to a URL that includes the locale.

### Lightweight Plugins

When a post is rendered, text substitution plugins are executed on it. These replace text matching a regular expression with other text or the output of a CLI command. It is required but not enforced that CLI output is deterministic.

When a post is saved, review plugins are executed via CLI and can access the content via the API.

Manual tools may be used for management via the API.

None of these are included in this repository, but they may be made available separately by Subtributary Software or other software providers.

### Post Search

Posts can be searched for, and the search includes titles, tags, and summaries.

### Risky Features Disabled by Default

Features that interface with anything outside of the Musings software will be disabled by default. This precaution enables safely hosting Musings blogs out of the box for customers.


## Non-Features

### No "/_" Route Prefixes

Routes beginning with "/_" are reserved and must not be configured in "content/routes.json". However, this is not enforced.

The routes "/_admin" and "/_api" are available to sysadmins for hosting those applications.

### No HTML or WYSIWYG Editor

The CMS only supports MarkDown. This is enough.

Accessory files (such as CSS and JS files) are not editable via the CMS, so there is no need to support editing those formats.

### No Image Editor

While an image manager is included, it does not support editing image content. It only supports what was explicitely stated in the features section.

### No Plugin Manager

Plugins are managed via the "config.json" file. There is no UI for this. This limitation ensures that risky plugin changes are performed by a server administrator only.

### No User Management

Users must be managed by a third-party application such as Authentik or FusionAuth.  Musings supports OIDC and that's it.


