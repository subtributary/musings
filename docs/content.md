# Creating content

Website content and routes are defined by the contents of "content/".
Files are served at routes based on their file paths.

Markdown (.md) files are rendered to HTML and served at routes matching their file paths minus the ".md".
Other files, such as images, are served unchanged at a routes matching their file paths.
For example, "content/2026/vacation.md" is rendered as HTML and served at "/2026/vacation",
and "content/2026/hotel-view.jpg" is served at "/2026/hotel-view.jpg".

## Naming conventions

Follow these naming conventions for the best experience:

* File names must be lowercase.
* File names must have an extension.
* File names must not contain spaces or underscores—use `-` to delimit words.
* Directory names must be lowercase.
    * The exception to this is locale region names. (See [localizing.md].)
* Directory names must not contain spaces or underscores—use `-` to delimit words.

If these naming conventions are not followed, Musings may still load, but the website may not work as intended.

## Markdown content format

Markdown (.md) files must follow CommonMark specifications with a few additions:
Musings adds plugin support and metadata to the format.

### Metadata

Metadata can be added in a frontmatter block.
All metadata fields are optional.

```markdown
---
bylines: [by Nathan Belue]
published: 2026-04-30
summary: Example summary text.
thumbnail: image.jpg
---
# Title

content
```

The `bylines` field is an array of strings.
These specify the bylines of the document,
which may include information such as the author, translator, and editor.

The `published` field is the publication date.
It must be in RFC 3339 format, "yyyy-MM-dd HH:mm:ss" format, or "yyyy-MM-dd" format.

The `summary` field is a summary of the content.
The default value is the first non-empty paragraph of the document.
(Stand-alone images are considered paragraphs, and their alt text can be used.)

The `thumbnail` field specifies an image to associate with the document.
It must be a relative or absolute URL *url* to an image file hosted by Musings.
The default index template shows this next to each listed document.

The default index template sorts content by most recent publication date then by title.
If the publication date is omitted, the document is sorted after others that have a publication date.
If the publication date is in the future, the document is not listed on the index.

### Plugins

(Coming soon!)

## Localization

See [localizing](localizing.md) to learn how to localize content.
