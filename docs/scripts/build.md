# Build Script

The `scripts/build.sh` script builds musings.

## Syntax

Syntax: `./scripts/build.sh [targets...]`

Parameters:
 - (optional) `targets` - What to build. One of "all" (default), "assets", or "bin".

## Internal Flow

Build binaries:
1. If `targets` is neither `bin` nor `all`, skip.
2. Build "./cmd/musings/main.go" to "./bin/musings".

Build assets:
1. Transpile "./custom/theme/*.scss" into "./web/static/*.css"
2. Minimize "./custom/theme/*.js" into "./web/static/*.js"
