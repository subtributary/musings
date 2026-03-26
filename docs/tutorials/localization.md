# Localization

Many files--including content, templates, and assets--can be created with localized variants.
When a user visits the website, Musings will detect their locale then use the most specific localized versions of files to serve them the website.

## How it works

Suppose these files exist:

 * /content/hello.md
 * /content/hello.fr.md
 * /content/hello.fr-ca.md

Musings considers these to localized variants of the exact same file.
When "hello.md" is needed, it will choose a variant based on the user's locale.

 * For Candian French users, "hello.fr-ca.md" will be used.
 * For every other French user, "hello.fr.md" will be used.
 * For everyone else, "hello.md" will be used.

## What can be localized?

Every single file operation in routed through our localizer,
so every single file that Musings opens can be localized.
There are no exceptions.

## What are the locales?

The list of locales, or a link to one, will be added in a future version.
