# Publish Script

The `scripts/publish.sh` script builds all projects and consolidates their binaries and assets into "./publish".

## Syntax

Syntax: `./scripts/publish.sh`

## Internal Flow

// todo

Build:
1. Execute `./scripts/build.sh all`.
2. If build fails, exit with error.

Create publish folder
1. Create directory `./publish`.
2. Create directories `./publish/admin`, `./publish/api`, and `./publish/web`.

Copy `admin` files:
1. Copy `./bin/admin` to `./publish/admin/bin/`
2. Copy `./`

Copy `api` files:

Copy `web` files:

