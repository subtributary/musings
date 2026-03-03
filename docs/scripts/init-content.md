# Init Content Script

The `scripts/init-content.sh` script creates and initializes the content repository for Musings.

## Syntax

Syntax: `./scripts/init-content.sh [directory]`

The `directory` parameter is optional. If omitted, the script will read it from `CONTENT_REPO_DIR` or prompt for it.

## Internal Flow

Check Sanity:
1. If `git` is not installed, exit with error.

Fetch parameters:
1. Read the directory from the parameter.
2. If the directory is not set, read it from `CONTENT_REPO_DIR`.
3. If the directory is not set, prompt the user for it.

Prepare repository:
1. If the directory does not exist, create it.
2. If the directory is not empty, exit with error.
3. If the directory is not a repository, execute `git init` in it.
4. If the directory does not have the `main` and `draft` worktrees, add them.

Update configs:
1. If "./.env" exists, update its `CONTENT_REPO_DIR` variable to point to the directory.
