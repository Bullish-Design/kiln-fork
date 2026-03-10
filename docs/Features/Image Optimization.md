---
title: "Image Optimization"
description: "Kiln automatically optimizes images at build time, generating responsive variants in modern formats like AVIF and WebP for faster page loads."
---
# Image Optimization

Kiln automatically optimizes images at build time, generating **responsive variants** in modern formats. Every `.png`, `.jpg`, and `.jpeg` in your vault is resized at multiple breakpoints and encoded into AVIF and WebP alongside the original format. The result is a `<picture>` element that lets browsers pick the smallest file they support — no configuration required.

## How it works

During a build, Kiln processes each optimizable image through a simple pipeline:

1. **Decode** the source image
2. **Resize** to each breakpoint that is smaller than the original width (1200px, 800px, 400px)
3. **Encode** every resized version into three formats — AVIF (best compression), WebP (good compression), and the original format (fallback)
4. **Generate HTML** using a `<picture>` element with `<source>` tags for each format

Each image can produce up to **9 variants** (3 breakpoints times 3 formats). All encoding uses a fixed quality of **80**. If the original image is narrower than a given breakpoint, that breakpoint is skipped.

Format priority in the generated HTML follows the order **AVIF > WebP > original**, so browsers that support AVIF will download the smallest file first.

## Supported formats

| Input format        | Optimized? | Output formats produced        |
| :------------------ | :--------- | :----------------------------- |
| `.png`              | Yes        | AVIF, WebP, PNG                |
| `.jpg` / `.jpeg`    | Yes        | AVIF, WebP, JPEG               |
| `.gif`              | No         | Copied as-is                   |
| `.svg`              | No         | Copied as-is                   |
| `.webp`             | No         | Copied as-is                   |

Images that are not optimized are still copied to the output directory — they just skip the resize-and-encode pipeline.

## Requirements

Kiln shells out to two external CLI tools for modern format encoding:

- **`cwebp`** (from the [libwebp](https://developers.google.com/speed/webp/docs/cwebp) package) — converts images to WebP
- **`avifenc`** (from the [libavif](https://github.com/AOMediaCodec/libavif) project) — converts images to AVIF

Install them on common systems:

```bash
# Debian / Ubuntu
sudo apt install webp libavif-bin

# macOS (Homebrew)
brew install webp libavif

# Nix
nix-shell -p libwebp libavif
```

> [!tip] Graceful degradation
> If either tool is missing, Kiln silently skips that format and falls back to copying the original image. Your build will not fail — you simply won't get the optimized variants for the missing encoder.

## The generated HTML

For each optimized image, Kiln emits a `<picture>` element wrapping multiple `<source>` tags and a fallback `<img>`. Here is a representative example for a PNG with all three breakpoints:

```html
<figure class="img-figure">
  <picture>
    <source type="image/avif"
            srcset="/img/photo-400w.avif 400w, /img/photo-800w.avif 800w, /img/photo-1200w.avif 1200w"
            sizes="min(65ch, 100vw)">
    <source type="image/webp"
            srcset="/img/photo-400w.webp 400w, /img/photo-800w.webp 800w, /img/photo-1200w.webp 1200w"
            sizes="min(65ch, 100vw)">
    <source type="image/png"
            srcset="/img/photo-400w.png 400w, /img/photo-800w.png 800w, /img/photo-1200w.png 1200w"
            sizes="min(65ch, 100vw)">
    <img src="/img/photo.png" alt="A photo" sizes="min(65ch, 100vw)" loading="lazy">
  </picture>
  <button class="img-expand-btn" aria-label="View full size" type="button">...</button>
</figure>
```

Key details:

- **`srcset`** lists each variant with its width descriptor (e.g., `400w`)
- **`sizes`** is set to `min(65ch, 100vw)`, matching the content column width
- **`loading="lazy"`** on the fallback `<img>` defers off-screen images
- An **expand button** (lightbox) is added to every image, allowing readers to view the full-size version

Images that have no optimized variants (unsupported format, or encoders missing) receive a plain `<img loading="lazy">` tag with the expand button.

## Usage

No configuration is needed. Reference images in your markdown the way you normally would:

```markdown
![A scenic photo](photo.jpg)
```

Or with [[Wikilinks]] syntax:

```markdown
![[photo.jpg]]
![[photo.jpg|Alt text goes here]]
```

Kiln detects supported formats and produces the optimized `<picture>` output automatically.

## Limitations

- **External tools required** — WebP and AVIF encoding depend on `cwebp` and `avifenc` being available on your PATH
- **Fixed breakpoints** — the resize widths (1200, 800, 400) are hard-coded and not configurable
- **Fixed quality** — encoding quality is set to 80 and cannot be changed
- **GIF and SVG are not optimized** — these formats are copied as-is
- **Increased output size** — each optimizable image can produce up to 9 variant files, increasing total disk usage
- **Longer builds** — image encoding adds build time proportional to the number and size of images in your vault
