// @feature:imgopt
package imgopt

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// newTestImage creates an NRGBA image with the given dimensions.
func newTestImage(w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.NRGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

// savePNG writes an image to disk as PNG for test fixtures.
func savePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
}

func TestResize(t *testing.T) {
	src := newTestImage(1600, 900)
	got := Resize(src, 800)
	bounds := got.Bounds()
	if bounds.Dx() != 800 || bounds.Dy() != 450 {
		t.Errorf("expected 800x450, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResize_SmallerThanTarget(t *testing.T) {
	src := newTestImage(400, 300)
	got := Resize(src, 800)
	bounds := got.Bounds()
	if bounds.Dx() != 400 || bounds.Dy() != 300 {
		t.Errorf("expected 400x300 (no-op), got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestDecodeImage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")
	savePNG(t, path, newTestImage(64, 48))

	img, format, err := DecodeImage(path)
	if err != nil {
		t.Fatalf("DecodeImage: %v", err)
	}
	if format != "png" {
		t.Errorf("expected format png, got %s", format)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 48 {
		t.Errorf("expected 64x48, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestProcessImage(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "photo.png")
	savePNG(t, srcPath, newTestImage(1600, 900))

	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	result, err := ProcessImage(srcPath, outDir, "/images", "photo", []int{800, 400})
	if err != nil {
		t.Fatalf("ProcessImage: %v", err)
	}

	if result.Original != srcPath {
		t.Errorf("expected Original=%s, got %s", srcPath, result.Original)
	}

	// At minimum we expect PNG variants for each breakpoint.
	pngCount := 0
	for _, v := range result.Variants {
		if v.Format == "png" {
			pngCount++
			if _, err := os.Stat(v.OutPath); err != nil {
				t.Errorf("variant file missing: %s", v.OutPath)
			}
		}
	}
	if pngCount < 2 {
		t.Errorf("expected at least 2 PNG variants, got %d", pngCount)
	}

	// Check that widths match breakpoints.
	widths := map[int]bool{}
	for _, v := range result.Variants {
		widths[v.Width] = true
	}
	for _, bp := range []int{800, 400} {
		if !widths[bp] {
			t.Errorf("missing variant for breakpoint %d", bp)
		}
	}
}

func TestIsOptimizable(t *testing.T) {
	tests := []struct {
		ext  string
		want bool
	}{
		{".png", true},
		{".jpg", true},
		{".jpeg", true},
		{".PNG", true},
		{".JPG", true},
		{".gif", false},
		{".svg", false},
		{".webp", false},
		{".avif", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsOptimizable(tt.ext); got != tt.want {
			t.Errorf("IsOptimizable(%q) = %v, want %v", tt.ext, got, tt.want)
		}
	}
}
