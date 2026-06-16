# How to contribute

General Subtributary contribution guidelines are in the organization documentation.

This document covers Musings-specific conventions.

## Path conventions

Paths returned from a package's API should be URL-style with a leading '/'.
For example, `/en/hello.txt`.

Paths used inside a package may use any appropriate representation.

Use `path`, not `filepath` for URL-style paths.
Use `filepath` only for operating-system filesystem paths.
