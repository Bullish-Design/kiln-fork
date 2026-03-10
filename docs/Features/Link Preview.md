---
title: "Link Preview"
description: "Kiln shows a live preview tooltip when you hover over internal links on desktop, letting readers peek at linked content without navigating away."
---
# Link Preview

Kiln shows a **live preview tooltip** when you hover over internal links, letting readers peek at the linked page's content without navigating away. The feature is fully automatic — it activates on every page with no opt-in or configuration needed.

## How it works

When a reader hovers over an internal link on a desktop device, Kiln follows this flow:

1. **Wait** 300 milliseconds to avoid triggering previews on casual mouse movement
2. **Fetch** the target page's full HTML in the background
3. **Extract** the main content area from the fetched page
4. **Render** it inside a floating tooltip positioned near the link

Fetched pages are stored in an **in-memory cache** that persists for the browser session. Once a page has been fetched, hovering over any link that points to it shows the preview instantly without a second network request.

The tooltip is positioned **below the link** by default. If there isn't enough space below (near the bottom of the viewport), it flips **above the link** instead. Horizontal position is constrained so the tooltip never overflows the viewport edges.

## Compatibility

Link Preview requires all three of the following conditions:

- **Desktop viewport** — screen width of 1024 pixels or more
- **Hover-capable device** — a pointer that supports hover (mouse or trackpad, not touch-only)
- **JavaScript enabled** — the feature is implemented entirely in client-side JavaScript

> [!note] Mobile and tablet
> On mobile and touch devices the feature is automatically disabled — readers simply follow links normally. There is nothing to configure; Kiln detects the device capabilities at page load.

## What triggers a preview

Previews appear only on **internal wikilinks** — the links Kiln's markdown renderer outputs with the `internal-link` class. These are the links you write in your notes using [[Wikilinks]] syntax.

The following links are **excluded** and will not trigger a preview:

- **External links** — links pointing outside your site
- **Hash-only links** (`#section`) — anchors within the current page
- **Current page links** — links whose resolved path matches the page you are already on
- **Links with `target="_blank"`** — links configured to open in a new tab

## Visual behavior

- **Short delay** — the tooltip appears after a 300-millisecond hover to prevent accidental triggers
- **Compact size** — maximum 400 pixels wide and 300 pixels tall; content scrolls if it overflows
- **Consistent styling** — the tooltip inherits your page's `.content` styles, so text, headings, code blocks, and tables look the same as they do on the actual page
- **Theme-aware** — follows [[Light-Dark Mode]] automatically via CSS variables
- **Dismissal** — the tooltip closes when you move the mouse away, scroll the page, press **Escape**, click outside the tooltip, or navigate to another page via [[Client Side Navigation]]

## Limitations

- **Desktop only** — the feature is entirely disabled on mobile and touch devices
- **No keyboard navigation or screen reader support** — the tooltip is mouse-driven and not reachable via keyboard focus
- **Full page fetch** — Kiln downloads the complete HTML of the target page and extracts the content area client-side; there is no dedicated lightweight preview endpoint
- **Unbounded session cache** — every fetched page is kept in memory for the duration of the browser session with no eviction strategy
- **No loading indicator** — there is no visual feedback while the target page is being fetched
- **No error feedback** — if a fetch fails (network error, 404), the tooltip simply does not appear
