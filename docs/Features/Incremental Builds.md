---
title: "Incremental Builds"
description: "How Kiln's dev command detects file changes and determines which pages to rebuild using a four-stage pipeline of watcher, mtime store, dependency graph, and changeset computation."
---
# Incremental Builds

Kiln's [`dev`](../Commands/dev.md) command uses an incremental build pipeline to detect file changes and rebuild only what is necessary. The pipeline consists of four components working together: a **filesystem watcher** that captures raw events, a **modification time store** that identifies what actually changed, a **dependency graph** that maps relationships between notes, and a **changeset computation** that combines everything into a minimal set of files to rebuild.

## Architecture

The four components form a linear pipeline. Each stage feeds its output into the next:

```
Filesystem Watcher
  │  (write, create, remove, rename events via fsnotify)
  ▼
Mtime Store
  │  (compares current mtimes against stored values)
  │  → changed files  (new or modified)
  │  → removed files  (previously tracked, now gone)
  ▼
Dependency Graph
  │  (maps wikilink and markdown link relationships)
  │  → forward edges  (source → targets it links to)
  │  → reverse edges  (target → sources that link to it)
  ▼
Changeset Computation
  │  (expands changed/removed lists using reverse edges)
  ▼
ChangeSet { Rebuild, Remove }
```

## Filesystem Watcher

The watcher uses [fsnotify](https://github.com/fsnotify/fsnotify) to monitor the input directory for filesystem events.

- **Recursive monitoring** — on startup, all subdirectories under the input directory are added to the watcher
- **Automatic discovery** — when a new subdirectory is created, it is automatically added to the watch list
- **Skipped paths** — dotfiles (`.obsidian`, `.git`) and `_hidden_` prefixed paths are ignored by both the initial walk and the event handler
- **Relevant events** — only write, create, remove, and rename operations trigger the rebuild timer; other events (chmod, etc.) are discarded
- **Debouncing** — rapid changes reset a debounce timer (default 300ms). Multiple saves within the debounce window collapse into a single rebuild, avoiding redundant work during fast editing sessions

## Modification Time Store

The mtime store tracks the last-known modification time of every file in the vault.

- On each rebuild trigger, it walks the input directory and compares the current modification time of every file against the previously stored value
- It reports two lists:
  - **Changed** — files that are new (not previously tracked) or whose modification time differs from the stored value
  - **Removed** — files that were previously tracked but are no longer present on disk
- After the walk, stored entries for removed files are deleted so they are not reported again
- Dotfiles and `_hidden_` prefixed paths are skipped, matching the vault scan behaviour

## Dependency Graph

The dependency graph records which notes link to which other notes. It is built from the vault's parsed wikilinks and markdown links after the initial scan.

- **Forward edges** — map each source file (by relative path) to the set of target names it links to
- **Reverse edges** — map each target name back to the set of source files that reference it
- Both `[[wikilink]]` and `[text](path.md)` syntax are parsed
- Target names are normalised: aliases (`|`) and anchors (`#`) are stripped, path prefixes are reduced to the base filename, and the result is lowercased
- External links (`https://`, `http://`, `mailto:`) are ignored
- When a file is removed, its forward edges are cleaned up from the reverse index via `RemoveSource`, so stale links do not trigger unnecessary rebuilds

## Changeset Computation

The changeset computation combines the changed/removed lists from the mtime store with the dependency graph to produce the final set of files that need action.

- For each **changed** file, the file itself and all its dependents (files that link *to* it) are added to the rebuild set
- For each **removed** file, its dependents are added to the rebuild set and the file itself is added to the remove set
- **Deduplication** — a file that appears as both directly changed and as a dependent of another changed file is only included once in the rebuild set
- **Exclusion** — removed files are excluded from the rebuild set, since there is no point rendering a deleted file
- The result is a `ChangeSet` with two fields:
  - `Rebuild` — relative paths of files to re-render
  - `Remove` — relative paths of files to delete from the output directory

## Lifecycle

When you run `kiln dev`, the following sequence occurs:

1. **Initial full build** — the entire vault is built exactly as the [Generate Command](../Commands/generate.md) would
2. **Mtime baseline populated** — the mtime store performs its first scan, recording every file's modification time. Since no previous entries exist, all files are reported as "changed" (this output is discarded during the initial build)
3. **Dependency graph built** — wikilinks and markdown links are parsed from every note to populate forward and reverse edges
4. **Watcher starts** — fsnotify begins monitoring the input directory and all its subdirectories
5. **Server starts** — a local HTTP server begins serving the output directory on `localhost:PORT`
6. **Edit → rebuild loop** — when a user edits a file, the sequence is: fsnotify event → debounce timer resets → timer fires after 300ms of inactivity → mtime store walks the directory and reports changed/removed files → changeset is computed using the dependency graph → rebuild is triggered

## Current Limitations

- **Full rebuild on each trigger** — the changeset is computed but not yet used to selectively rebuild individual pages; the current rebuild performs a full site generation
- **Static dependency graph** — the graph is built once at startup and not updated after incremental changes; new links added during editing are not reflected until a restart
- **Limited dependency tracking** — only wikilink and markdown link dependencies are tracked. Changes to templates, layouts, or theme files require a full restart of the `dev` command

See the [Dev Command](../Commands/dev.md) page for usage details and available flags.
