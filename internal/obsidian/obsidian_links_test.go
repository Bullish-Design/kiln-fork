// @feature:vault-scan Tests for markdown link extraction in processMarkdown.
package obsidian

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// newTestFile creates a temporary markdown file and returns a File ready for processMarkdown.
func newTestFile(t *testing.T, content string) *File {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return &File{
		Path:      path,
		Name:      "test",
		Ext:       ".md",
		Links:     []string{},
		Backlinks: []string{},
		Tags:      make(map[string]struct{}),
		Embeds:    []string{},
	}
}

func TestProcessMarkdown_MdLinkSibling(t *testing.T) {
	f := newTestFile(t, `[text](./sibling.md)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(f.Links, "[text](./sibling.md)") {
		t.Errorf("expected Links to contain markdown link, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdLinkParent(t *testing.T) {
	f := newTestFile(t, `[text](../parent/note.md)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(f.Links, "[text](../parent/note.md)") {
		t.Errorf("expected Links to contain markdown link, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdLinkRelativeFolder(t *testing.T) {
	f := newTestFile(t, `[text](folder/note.md)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(f.Links, "[text](folder/note.md)") {
		t.Errorf("expected Links to contain markdown link, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdLinkExternalSkipped(t *testing.T) {
	f := newTestFile(t, `[text](https://example.com)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if len(f.Links) != 0 {
		t.Errorf("expected no links for external URL, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdLinkAnchorSkipped(t *testing.T) {
	f := newTestFile(t, `[text](#heading)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if len(f.Links) != 0 {
		t.Errorf("expected no links for pure anchor, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdImageEmbed(t *testing.T) {
	f := newTestFile(t, `![alt](image.png)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(f.Embeds, "[alt](image.png)") {
		t.Errorf("expected Embeds to contain image link, got %v", f.Embeds)
	}
	if !slices.Contains(f.Links, "[alt](image.png)") {
		t.Errorf("expected Links to contain image link, got %v", f.Links)
	}
}

func TestProcessMarkdown_MixedWikiAndMdLinks(t *testing.T) {
	f := newTestFile(t, `[[wiki-note]]
[md-note](./path.md)
![[embed.png]]
![img](photo.jpg)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}

	// Wikilink
	if !slices.Contains(f.Links, "[[wiki-note]]") {
		t.Errorf("expected Links to contain wikilink, got %v", f.Links)
	}
	// Markdown link
	if !slices.Contains(f.Links, "[md-note](./path.md)") {
		t.Errorf("expected Links to contain markdown link, got %v", f.Links)
	}
	// Wikilink embed in both Embeds and Links
	if !slices.Contains(f.Embeds, "[[embed.png]]") {
		t.Errorf("expected Embeds to contain wikilink embed, got %v", f.Embeds)
	}
	if !slices.Contains(f.Links, "[[embed.png]]") {
		t.Errorf("expected Links to contain wikilink embed, got %v", f.Links)
	}
	// Markdown image in both Embeds and Links
	if !slices.Contains(f.Embeds, "[img](photo.jpg)") {
		t.Errorf("expected Embeds to contain markdown image, got %v", f.Embeds)
	}
	if !slices.Contains(f.Links, "[img](photo.jpg)") {
		t.Errorf("expected Links to contain markdown image, got %v", f.Links)
	}
}

func TestProcessMarkdown_MdLinkHttpSkipped(t *testing.T) {
	f := newTestFile(t, `[a](http://example.com) [b](mailto:user@test.com)`)
	if err := f.processMarkdown(); err != nil {
		t.Fatal(err)
	}
	if len(f.Links) != 0 {
		t.Errorf("expected no links for external URLs, got %v", f.Links)
	}
}
