# Content metadata

Metadata such as bylines and tags can be added to Markdown (*.md) files.
Most of this metadata, except the title,  can be specified in what is called a frontmatter block at the top of the file.
The title is a special case that is detailed later in this document.
A document with all of its metadata defined looks like this:

```markdown
---
bylines: by Nathan Belue
published: 2026-04-30
tags: [apple, banana, cucumber]
---
# Title

content
```


## Frontmatter metadata

Musings supports a subset of [YAML](https://yaml.org/spec/1.2.2/) for metadata in the frontmatter.
The frontmatter must be demarcated by triple dashes (`---`).
The following metadata properties can be set in the frontmatter:

- bylines
- published
- tags

These are detailed in the subsections below.

The YAML parser may change in future versions.
To ensure compatibility, follow these guidelines:

1. Lowercase all property names.
2. Do not set undefined or user-defined properties.
3. Do not use aliases, anchors, or other YAML features not demonstrated herein.

### Property: bylines

Bylines credit those who participated in the creation of content.
The `bylines` property can be a single string value or a list of string values.

These formats are supported:

```yaml
bylines: by Nathan Belue
```

```yaml
bylines: [by Nathan Belue, translated by 네이선]
```

```yaml
bylines:
  - by Nathan Belue
  - translated by 네이선 
```

### Property: published

The publication date is the time when the content is made available to the target audience.
The `published` property is a time in the format "yyyy-MM-dd", "yyyy-MM-dd HH:mm:ss", or [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339).
Here are examples of each:

```yaml
# April 30, 2026
published: 2026-04-30
```

```yaml
# April 30, 2026 at 3:20PM
published: 2026-04-30 15:20:00
```

```yaml
# April 30, 2026, at 3:20PM in Coordinated Universal Time 
published: 2026-04-30T15:20:00Z
```

RFC 3339, which is the last example, is suggested because it can include the time zone offset.
The other formats will adopt an undefined but consistent time zone.

### Property: tags

Tags are used to categorize a document for more effective searching within a corpus.
The `tags` property can be a single string value or a list of string values.

These formats are supported:

```yaml
tags: by Nathan Belue
```

```yaml
tags: [apple, orange]
```

```yaml
tags:
  - apple
  - banana
```

## Special case: title

The title describes the topic or purpose of the document and is usually displayed prominently at the top of the document.

The first top-level H1 heading (`#`) of the content is extracted into the `title` property.
The `title` property cannot be set in the frontmatter.

In the default Musings configuration, the title is displayed at the top of the page, in the title bar, and in search results.
But the specific locations where it is displayed can be customized by a developer.

All formatting is stripped from the title,
so Musings does not support a formatted first top-level heading.
