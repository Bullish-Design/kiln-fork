---
title: Accent Color
description: Override the default accent color of any theme with a named color from its own palette, keeping visual harmony while personalizing your site.
---
# Accent Color

By default, each [[Themes|theme]] comes with its own built-in accent color — the tint used for links, active states, scrollbar hovers, and other interactive elements. The accent color feature lets you override this default with any color from the theme's own palette, keeping visual harmony while personalizing your site.

## How it works

Every theme defines a palette of 7 named colors: `red`, `orange`, `yellow`, `green`, `blue`, `purple`, and `cyan`. When you set an accent color, Kiln replaces the theme's default accent in both light and dark mode variants with the palette color matching that name.

The override happens at build time — the chosen color is baked directly into the CSS custom property `--accent-color`. No client-side JavaScript is involved.

If the color name is not recognized, Kiln logs a warning and keeps the theme's original accent.

## Available colors

| Name     | Description                  |
| :------- | :--------------------------- |
| `red`    | Red from the theme palette   |
| `orange` | Orange from the theme palette |
| `yellow` | Yellow from the theme palette |
| `green`  | Green from the theme palette |
| `blue`   | Blue from the theme palette  |
| `purple` | Purple from the theme palette |
| `cyan`   | Cyan from the theme palette  |

> The actual hex values depend on the chosen [[Themes|theme]] — for example, `red` in Dracula is `#ff5555` while in Nord it's `#bf616a`.

## Configuration

To set an accent color, use the `--accent-color` flag (short form `-a`) when running the [[generate]] command.

**Example: Using blue accent**
```bash
kiln generate --accent-color blue
```

Or set it permanently in your [[Configuration File]]:

```yaml
accent-color: blue
```

CLI flags override the [[Configuration File]] when both are provided.

**Example: Combining theme and accent color**
```bash
kiln generate --theme dracula --accent-color cyan
```

## Default behavior

If no accent color is specified, each theme uses its own built-in accent color (e.g., `#7e6df7` for the default theme, `#ff79c6` for Dracula). No additional configuration is needed.
