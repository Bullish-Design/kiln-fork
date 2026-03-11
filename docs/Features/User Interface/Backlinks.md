---
title: Backlinks
description: Kiln automatically detects and displays backlinks in the right sidebar, showing every page in your vault that links to the current page.
---
# Backlinks

Kiln scans your vault for [[Wikilinks]] and standard Markdown links, then builds a reverse index. For every page, the **Backlinks** panel in the **Right Sidebar** lists all other pages that link to it. This helps readers discover related content and navigate your knowledge graph in reverse.

## Why it is useful

- In a wiki or digital garden, backlinks reveal hidden relationships between notes.
- They answer "What pages reference this concept?" without manual curation.
- They complement the [[Local Graph]] by providing a clickable list of incoming connections.

## How it works

During the build process, Kiln constructs a complete reverse-link index for your vault.

1. **Scan outgoing links:** Kiln scans every note for outgoing links—both `[[wikilinks]]` and `[markdown](links.md)` style.
2. **Record reverse references:** For each link found, a reverse reference is recorded on the target page.
3. **Build the list:** The result is a list of all pages that point to the current page.
4. **Exclude self-links:** A page cannot appear in its own backlinks.
5. **Deduplicate:** Duplicate links from the same source page are deduplicated, so each source appears only once.

## Behavior by Layout

### Default Layout

- Backlinks appear in the **Right Sidebar**, below the [[Table of Contents]] and [[Local Graph]].
- Each backlink is a clickable link that navigates to the source page.
- Navigation uses [[Client Side Navigation]] (HTMX) for instant page transitions without a full reload.

### Simple Layout

- Backlinks are accessed via a floating action button in the header bar.
- Clicking the button opens a panel overlay listing all incoming links.
- The button only appears when the current page has at least one backlink.

## Disabling Backlinks

If your site does not need backlinks, you can hide them using the `--disable-backlinks` flag:

```bash
kiln generate --disable-backlinks
```

Or set it permanently in your `kiln.yaml` configuration file:

```yaml
disable-backlinks: true
```

When disabled, the backlinks panel is hidden from both the default and simple layouts. See [[Configuration File]] for all available options.
