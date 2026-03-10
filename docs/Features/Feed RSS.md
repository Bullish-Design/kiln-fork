---
title: "Feed RSS"
description: "Kiln automatically generates an RSS 2.0 feed so readers can subscribe to your site with any feed reader. Learn how to enable and customize it."
---
# Feed RSS

Kiln automatically generates an **RSS 2.0 feed** every time you build your site. RSS (Really Simple Syndication) lets readers subscribe to your content using feed readers like Feedly, NetNewsWire, or Thunderbird — so they get notified whenever you publish something new, without having to check your site manually.

## Why it matters

An RSS feed turns your static site into a living stream of updates. Readers who prefer feed aggregators can follow your content alongside dozens of other sources in a single app. For blogs, digital gardens, and documentation sites alike, offering an RSS feed is a low-effort way to keep your audience engaged and coming back.

## Configuration

Because an RSS feed requires **absolute URLs** (e.g., `https://example.com/my-note` instead of just `/my-note`), Kiln needs to know your domain name to generate it.

You must provide the `--url` flag when running the generate command:

```bash
./kiln generate --url "https://your-domain.com"
```

> [!warning] Missing URL Flag
> If you do not provide the `--url` flag, Kiln will skip generating the RSS feed entirely to prevent creating an invalid file with relative links.

You can also set the URL permanently in your [[Configuration File]] so you don't have to pass the flag every time.

## How it works

During the build, Kiln automatically collects metadata from every `.md` file in your vault and assembles them into a single feed. No opt-in or configuration beyond the `--url` flag is needed.

- **All `.md` files** are included automatically
- **Title** is taken from the frontmatter `title` field; falls back to the filename if omitted
- **Description** is taken from the frontmatter `description` field (optional)
- **Publication date** is derived from the file's creation time (birth time)
- Entries are **sorted newest first**, limited to the **50 most recent**
- Output file is `feed.xml` in the root of the output directory

The feed follows the RSS 2.0 specification. The channel title and link are set to your site's base URL.

## Customizing entries

The best way to control how your pages appear in the feed is through **frontmatter fields**. Each note's frontmatter directly maps to the corresponding RSS item:

```yaml
---
title: "My Latest Post"
description: "A short summary that appears in feed readers."
---
```

| Frontmatter field | RSS element   | Fallback              |
| :---------------- | :------------ | :-------------------- |
| `title`           | Item title    | Filename without `.md` |
| `description`     | Item description | Empty (omitted)     |

Adding a descriptive `title` and `description` to your notes ensures they look good in every feed reader. See [[Sitemap xml]] and [[Structured Data (SEO)]] for other features that benefit from the same frontmatter fields.

## Limitations

- **Maximum 50 entries** — the feed includes only the 50 most recent files (hard-coded)
- **All `.md` files included** — there is no way to exclude specific files from the feed
- **No full-text content** — the feed contains titles and descriptions only, not the full body of each note
- **Publication date uses file creation time** — you cannot override it via frontmatter
- **Channel title uses the base URL** — there is no option to set a custom site name for the feed title
