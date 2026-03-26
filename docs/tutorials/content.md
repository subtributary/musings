# Creating content

Musings is a publishing tool for end users to publish content.

## Content directory and routes

The "/content/" directory stores the served content.
The internal directory structure is not restricted,
so the user can organize the files (and routes) as they wish.

The routes for the content files match their paths under "/content/",
except MarkDown ("*.md) files have their extensions trimmed.
For examples, "/content/hello.md" will be at the URL "/hello",
and "/content/hello.png" will be at the URL "/hello.png".

## File naming conventions

Files and subdirectories should follow these conventions:

1. Names should be lowercase.
2. Names should not contain spaces; use "-" to separate words.
3. File names should not contain periods.

The last convention is due to our [l10n](localization.md) naming rules,
which are confusing to use with files that contain periods.

## Localization

Every single file in "/content/" can be localized with a default locale fallback.

To localize a file, create a new one with the locale like an extension before the file extension.
For example, create "hello.fr.md" for the French localization of "hello.md".

See [localization.md](localization.md) for more details.
