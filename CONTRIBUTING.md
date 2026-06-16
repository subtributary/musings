# How to contribute

General Subtributary contribution guidelines are in the organization documentation.

This document covers Musings-specific conventions.

## Conventions

### Paths

In documentation, it should be clear whether a file path is a directory or file.
Unless a file path is obviously a directory based on context or surrounding words,
append a slash ('/') to its name; if it is obvious, then omit the slash.
This rule does not apply to non-file paths in documentation.

In code, when a path is returned from an internal package's API,
it should return a URL-style path with a leading slash ('/').
For example, `/en/hello.md`.
This rule does not necessarily apply to public APIs.

Paths used inside a package may use any appropriate representation.

Use `path`, not `filepath` for URL-style paths.
Use `filepath` only for operating-system filesystem paths.
