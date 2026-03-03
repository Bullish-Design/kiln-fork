// @feature:ogimage Tests for OG image generation.
package ogimage

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func defaultConfig() ImageConfig {
	return ImageConfig{
		Title:       "Hello World",
		Description: "A short description of the page.",
		SiteName:    "My Site",
		AccentColor: "#6c63ff",
		BgColor:     "#1e1e2e",
		TextColor:   "#cdd6f4",
	}
}

func TestGenerateOGImage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "og.png")

	cfg := defaultConfig()
	if err := GenerateOGImage(cfg, out); err != nil {
		t.Fatalf("GenerateOGImage returned error: %v", err)
	}

	f, err := os.Open(out)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 1200 || bounds.Dy() != 630 {
		t.Errorf("expected 1200x630, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateTwitterImage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "twitter.png")

	cfg := defaultConfig()
	if err := GenerateTwitterImage(cfg, out); err != nil {
		t.Fatalf("GenerateTwitterImage returned error: %v", err)
	}

	f, err := os.Open(out)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 1200 || bounds.Dy() != 600 {
		t.Errorf("expected 1200x600, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateOGImage_LongTitle(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "long.png")

	cfg := defaultConfig()
	cfg.Title = strings.Repeat("A very long title ", 10) // well over 80 chars

	if err := GenerateOGImage(cfg, out); err != nil {
		t.Fatalf("GenerateOGImage with long title returned error: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}
	if fi.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestGenerateOGImage_EmptyTitle(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "empty.png")

	cfg := defaultConfig()
	cfg.Title = ""

	if err := GenerateOGImage(cfg, out); err != nil {
		t.Fatalf("GenerateOGImage with empty title returned error: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}
	if fi.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestGenerateOGImage_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "nested", "path")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	out := filepath.Join(subdir, "og.png")

	cfg := defaultConfig()
	if err := GenerateOGImage(cfg, out); err != nil {
		t.Fatalf("GenerateOGImage returned error: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("expected file at %s, got error: %v", out, err)
	}
	if fi.Size() == 0 {
		t.Error("output file is empty")
	}
}
