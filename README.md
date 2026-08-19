# Musings

> [!NOTE]
> Musings is currently being reorganized into multiple repos.
> It probably won't work right until this move is complete.

Musings is a publishing CMS focused on text-heavy content like blogs, stories, and news articles.
It features first-class localization, markdown content, and a lightweight plugin mechanism.

For detailed documentation, see the "docs" directory.

## Building from source

The backend is compiled via `go`:

```shell
go build -o bin/server ./cmd/server
```

The frontend is compiled via `npm`:

```shell
npm install # needed on first build
npm run build
```

All the above commands should be run from the repository root. (e.g. ~/dev/musings)


## Running the server

After building Musings, run it directly:

```
bin/server
```

During frontend customization, add the `--live-templates` flag for convenience.

By default, Musings listens at ":8080".  This can be changed via the `--bind` parameter or the `MUSINGS_BIND` environment variable.


## Customizing

The frontend of Musings is designed to be customized.
To avoid breaking changes, limit customizations to the "web" directory.

* `web/src/scss/` - Styles, transpiled to `web/static/css/` and served at the URL `/_static/css/`.
* `web/src/ts/` - Scripts, transpiled to `web/static/js/` and served at the URL `/_static/js/`.
* `web/static/media/` - Images and other media, served at the URL `/_static/media/`.
* `templates/partials/` - HTML templates used by other templates.
* `templates/index.gohtml` - HTML template used to render the index.
* `templates/post.gohtml` - HTML template used to render a post.

For more information on customizing Musings, see [docs/customizing.md].

## Deploying

Most of the files in the repository are not needed to run Musings.
The files and directories that are needed are:

* bin/server
* content/
* data/
* web/static/
* web/templates/

When updating an existing installation, omit the "content" directory to avoid overwriting user content, and review the files in the "data" directory to ensure they match the desired user configuration.

For more information on deploying Musings, see [docs/deploying.md].
