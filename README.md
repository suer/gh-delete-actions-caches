# gh-delete-actions-caches

A GitHub CLI extension to bulk delete GitHub Actions caches by key prefix or ref prefix.

## Installation

```bash
gh extension install suer/gh-delete-actions-caches
```

To upgrade:

```bash
gh extension upgrade suer/gh-delete-actions-caches
```

## Usage

```bash
gh delete-actions-caches --key-prefix <prefix>
gh delete-actions-caches --ref-prefix <prefix>
gh delete-actions-caches --key-prefix <prefix> --ref-prefix <prefix>
```

At least one of `--key-prefix` or `--ref-prefix` is required.

### Flags

- `-k, --key-prefix`: Delete caches whose key starts with this prefix
- `-r, --ref-prefix`: Delete caches whose ref starts with this prefix
- `-d, --dryrun`: List caches to be deleted without actually deleting
- `-R, --repo`: Target repository (`OWNER/REPO`). Defaults to the repository of the current directory
- `--version`: Show version

### Examples

Delete caches by key prefix:

```bash
gh delete-actions-caches --key-prefix "Linux-build"
```

Delete caches by ref prefix:

```bash
gh delete-actions-caches --ref-prefix "refs/heads/feature/"
```

Delete caches matching both key prefix and ref prefix (AND condition):

```bash
gh delete-actions-caches --key-prefix "Linux-build" --ref-prefix "refs/heads/main"
```

Preview caches to be deleted without actually deleting (dry run):

```bash
gh delete-actions-caches --key-prefix "Linux-build" -d
```

Operate on a specific repository:

```bash
gh delete-actions-caches --key-prefix "Linux-build" -R owner/repo
```

## GitHub Actions

This repository can also be used as an action.

```yaml
      - name: Delete actions caches
        uses: suer/gh-delete-actions-caches@master
        with:
          key-prefix: buildkit-blob-1-
```

The job needs the `actions: write` permission (`actions: read` is enough for `dryrun: true`).

```yaml
jobs:
  delete-caches:
    runs-on: ubuntu-latest
    permissions:
      actions: write
    steps:
      - name: Delete actions caches
        uses: suer/gh-delete-actions-caches@master
        with:
          key-prefix: buildkit-blob-1-
          ref-prefix: refs/heads/main
          dryrun: false
```

### Inputs

- `key-prefix`: Delete caches whose key starts with this prefix
- `ref-prefix`: Delete caches whose ref starts with this prefix
- `dryrun`: List caches to be deleted without actually deleting (default: `false`)

At least one of `key-prefix` or `ref-prefix` is required. Caches are always deleted from the repository running the workflow.

### Token

`${{ github.token }}` is used by default. To use another token, set `GH_TOKEN` in the workflow:

```yaml
      - name: Delete actions caches
        uses: suer/gh-delete-actions-caches@master
        with:
          key-prefix: buildkit-blob-1-
        env:
          GH_TOKEN: ${{ secrets.MY_PAT }}
```

The action installs the extension with `gh`, so it requires a runner with the GitHub CLI available (GitHub-hosted runners have it preinstalled). Referencing a tag such as `@v0.0.2` instead of `@master` also pins the CLI version the action installs.

## Output Format

### Dry run

```
3 cache(s) to be deleted:
  key: Linux-build-abc123, ref: refs/heads/main
  key: Linux-build-def456, ref: refs/heads/feature/foo
  key: Linux-build-ghi789, ref: refs/heads/main
```

### Normal run

```
Deleted: Linux-build-abc123 (ref: refs/heads/main)
Deleted: Linux-build-def456 (ref: refs/heads/feature/foo)
Deleted: Linux-build-ghi789 (ref: refs/heads/main)
3 cache(s) deleted.
```

## For developers

To build and install locally:

```bash
go build .
gh extension install .
```
