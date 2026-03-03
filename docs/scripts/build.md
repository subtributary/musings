# Build Script

The `scripts/build.sh` script builds all or specific projects.

## Syntax

Syntax: `./scripts/build.sh {targets...}`

The `targets` parameter is one or more of "admin", "api", or "web", or it is "all".

## Internal Flow

Normalization:
1. If `targets` is "all", set to ["admin", "api", "web"].

Admin:
1. If `targets` does not include "admin", skip.
2. Transpile "custom/theme/admin*.scss" into "web/static/admin*.css"
3. Transpile "custom/theme/admin*.js" into "web/static/admin*.js"
4. Build "cmd/admin/main.go"

API:
1. If `targets` does not include "api", skip.
2. Build "cmd/api/main.go"

Web
1. If `targets` does not include "web", skip.
2. Transpile "custom/theme/web*.scss" into "web/static/web*.css"
3. Transpile "custom/theme/web*.ts" into "web/static/web*.js"
