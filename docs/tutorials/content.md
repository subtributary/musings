# Creating content

Musings is a publishing tool that allows end users to publish content.

## Content directory and routes

The "/content/" directory stores the served content.
The internal directory structure is not restricted,
allowing you to organize content and routes as desired.

The routes for the content files match their paths under "/content/",
but Markdown ("*.md") files have their extensions trimmed.
For examples, "/content/hello.md" will be at the URL "/hello",
and "/content/hello.png" will be at the URL "/hello.png".

If a Markdown file and other file would resolve to the same route,
then the behavior is undefined and should be avoided.
For example, "/content/hello.png.md" and "/content/hello.png" resolve to the same route.

## File naming conventions

Files and subdirectories should follow these conventions:

1. Names should be lowercase.
2. Names should not contain spaces; use "-" to separate words.
3. File names (excluding the extension) should not contain periods if localization is enabled.

The last convention is due to our [localization](localization.md) naming rules,
which are confusing to use with files that contain periods.

## Localization

Files "/content/" can have localized variants.

To localize a file, create a new one with the locale as a second extension before the file extension.
For example, create "hello.fr.md" as the French localization of "hello.md".

See [localization.md](localization.md) for more details.
