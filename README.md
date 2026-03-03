# Musings
Musings is a simple, low-maintenance, self-hostable publishing CMS for blog/news-style sites with first-class localization and a lightweight plugin mechanism.

## Folder Structure

...

## Initial Setup

Musings stores its content in a separate Git repository. This could be configured manually, but it's easiest to use our included script:

```shell
source .env && scripts/init-content.sh
```

The script will also set the environment variable `CONTENT_REPO_DIR` in the ".env" file. Verify that all the settings therein are as you want them:

```shell
cat .env
```

Sensible defaults are preset, so it is unlikely that they will need changed for development work.

## Building Binaries

The backend is written in Go. The frontend is written in CSS, JavaScript, React, SCSS, and TypeScript. For building, just use our script:

```shell
./scripts/build.sh {projects}
```

For the `projects` parameter, specify one or more of "admin", "api", or "web", or specify "all" for all three. The output will be in "./bin".

## Customizing

To avoid breaking changes, limit customizations to the "custom" folder and use the sdk in "./pkg/sdk".

## Deploying

The binaries and their dependencies are spread across directories, so run this script to consolidate them into "./publish":

```shell
./scripts.publish.sh
```

This will build all projects, then copy them and their dependencies into "./publish" with this directory structure:

 - admin/
   - bin/
   - public/
   - config.json
   - .env
 - api/
   - bin/
   - config.json
   - .env
 - web/
   - bin/
   - public/
   - config.json
   - .env

These can be hosted on subdomains (recommended) or on the same domain. If hosting on the same domain, "~/_admin" and "~/_api" are the recommended paths for `admin` and `api`, respectfully.

> [!WARNING]
> Hosting `admin` and `web` on the same host will allow them to share cookies, which is risky.

By default, plugins are disabled to reduce XSS risks. Please do not enable them if hosting everything on the same domain. If you do want to enable them, edit the plugin whitelist in "config.json", and run the website with the `--enable-plugins` parameter.
