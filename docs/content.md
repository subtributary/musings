# Creating content

Musings is a publishing tool with a focus on text-heavy content like blogs, stories, and news websites.

## Content directory and routes

The "content" directory stores the served content.
The internal directory structure is not restricted,
except for localization rules,
allowing you to organize content and routes as desired.

The routes for the content files match their paths under "content",
but Markdown ("*.md") files have their extensions trimmed.
For example, "content/hello.md" will be at the URL "/hello",
and "content/hello.png" will be at the URL "/hello.png".

If a Markdown file and another file resolve to the same route,
then the behavior is undefined and should be avoided.
For example, "content/hello.png.md" and "content/hello.png" resolve to the same route.

If localization is enabled, 
then the immediate subdirectories of the "content" directory must be the configured locales.
The routes will still match the directory structure. 


## Content format

The content must follow CommonMark specifications with a few additions.
Musings adds plugin support (coming soon!) and metadata to the format.

### Plugins

(To be implemented.)

### Metadata

Metadata such as bylines and publication date can be added in a frontmatter block.

```markdown
---
bylines: [by Nathan Belue]
published: 2026-04-30
---
# Title

content
```

The `bylines` field is an array of strings.
These specify the bylines of the document,
which may include information such as the author, translator, and editor.

The `published` field is the publication date.
It must be in RFC 3339 format, "yyyy-MM-dd HH:mm:ss" format, or "yyyy-MM-dd" format.

By default, content is sorted by publication date on the index.
If a publication date is omitted for a document, that document is sorted above others that have a publication date.
If a publication date is in the future, that document is not listed on the index.

## Localization

See [localizing](localizing.md) for how to localize content.
