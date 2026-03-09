// @feature:wikilinks Tests for wikilink resolution and nil pointer safety.
package markdown

import (
	"errors"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func newTestResolver(index map[string][]*obsidian.File) *IndexResolver {
	return &IndexResolver{
		Index:     index,
		SourceMap: map[string]string{},
		Links:     []obsidian.GraphLink{},
	}
}

func TestFindFile_NotInIndex(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})

	file, anchor, err := r.FindFile([]byte("nonexistent"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
	if anchor != "" {
		t.Errorf("expected empty anchor, got %q", anchor)
	}
}

func TestFindFile_NotInIndexWithAnchor(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})

	file, anchor, err := r.FindFile([]byte("nonexistent#heading"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
	if anchor != "#heading" {
		t.Errorf("expected anchor '#heading', got %q", anchor)
	}
}

func TestFindFile_EmptyCandidateSlice(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{
		"empty": {},
	})

	file, _, err := r.FindFile([]byte("empty"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
}

func TestFindFile_SingleCandidate(t *testing.T) {
	expected := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {expected},
	})

	file, anchor, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
	if anchor != "" {
		t.Errorf("expected empty anchor, got %q", anchor)
	}
}

func TestFindFile_PathBasedLookupNoMatch(t *testing.T) {
	candidate := &obsidian.File{
		Name:    "note",
		RelPath: "other/note.md",
		Path:    "/vault/other/note.md",
		Ext:     ".md",
		WebPath: "/other/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {candidate},
	})

	// Path-based lookup where the path doesn't match — should fallback to first candidate
	file, _, err := r.FindFile([]byte("folder/note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != candidate {
		t.Errorf("expected fallback to first candidate, got %+v", file)
	}
}

func TestFindFile_PathBasedLookupMatch(t *testing.T) {
	c1 := &obsidian.File{
		Name:    "note",
		RelPath: "folder/note.md",
		Path:    "/vault/folder/note.md",
		Ext:     ".md",
		WebPath: "/folder/note",
	}
	c2 := &obsidian.File{
		Name:    "note",
		RelPath: "other/note.md",
		Path:    "/vault/other/note.md",
		Ext:     ".md",
		WebPath: "/other/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {c1, c2},
	})

	file, _, err := r.FindFile([]byte("folder/note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != c1 {
		t.Errorf("expected c1, got %+v", file)
	}
}

func TestFindFile_ShortestPathFallback(t *testing.T) {
	long := &obsidian.File{
		Name:    "note",
		RelPath: "a/b/c/note.md",
		Path:    "/vault/a/b/c/note.md",
		Ext:     ".md",
		WebPath: "/a/b/c/note",
	}
	short := &obsidian.File{
		Name:    "note",
		RelPath: "a/note.md",
		Path:    "/vault/a/note.md",
		Ext:     ".md",
		WebPath: "/a/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {long, short},
	})

	file, _, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != short {
		t.Errorf("expected shortest path candidate, got %+v", file)
	}
}

func TestFindFile_RootFilePreferred(t *testing.T) {
	nested := &obsidian.File{
		Name:    "note",
		RelPath: "folder/note.md",
		Path:    "/vault/folder/note.md",
		Ext:     ".md",
		WebPath: "/folder/note",
	}
	root := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {nested, root},
	})

	file, _, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != root {
		t.Errorf("expected root file, got %+v", file)
	}
}

func TestFindFile_CaseInsensitiveFallback(t *testing.T) {
	expected := &obsidian.File{
		Name:    "MyNote",
		RelPath: "MyNote.md",
		Path:    "/vault/MyNote.md",
		Ext:     ".md",
		WebPath: "/mynote",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"mynote": {expected},
	})

	file, _, err := r.FindFile([]byte("MyNote"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
}

func TestFindFile_StripsMdExtension(t *testing.T) {
	expected := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {expected},
	})

	file, _, err := r.FindFile([]byte("note.md"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
}

func TestResolveMarkdownLink(t *testing.T) {
	tests := []struct {
		name          string
		currentSource string
		dest          string
		want          string
	}{
		{
			name:          "relative sibling link",
			currentSource: "/folder/current",
			dest:          "./sibling.md",
			want:          "/folder/sibling",
		},
		{
			name:          "parent traversal link",
			currentSource: "/folder/current",
			dest:          "../other/note.md",
			want:          "/other/note",
		},
		{
			name:          "absolute path",
			currentSource: "/folder/current",
			dest:          "/folder/note.md",
			want:          "/folder/note",
		},
		{
			name:          "external URL https",
			currentSource: "/folder/current",
			dest:          "https://example.com",
			want:          "https://example.com",
		},
		{
			name:          "anchor-only link",
			currentSource: "/folder/current",
			dest:          "#heading",
			want:          "#heading",
		},
		{
			name:          "relative link with anchor",
			currentSource: "/folder/current",
			dest:          "./note.md#heading",
			want:          "/folder/note#heading",
		},
		{
			name:          "non-md file",
			currentSource: "/folder/current",
			dest:          "./image.png",
			want:          "/folder/image.png",
		},
		{
			name:          "external URL http",
			currentSource: "/folder/current",
			dest:          "http://example.com/page",
			want:          "http://example.com/page",
		},
		{
			name:          "mailto link",
			currentSource: "/folder/current",
			dest:          "mailto:user@example.com",
			want:          "mailto:user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestResolver(map[string][]*obsidian.File{})
			r.CurrentSource = tt.currentSource

			got := r.ResolveMarkdownLink(tt.dest)
			if got != tt.want {
				t.Errorf("ResolveMarkdownLink(%q) = %q, want %q", tt.dest, got, tt.want)
			}
		})
	}
}

func TestResolveMarkdownLink_RecordsGraphLink(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})
	r.CurrentSource = "/folder/current"

	r.ResolveMarkdownLink("./sibling.md")

	if len(r.Links) != 1 {
		t.Fatalf("expected 1 graph link, got %d", len(r.Links))
	}
	if r.Links[0].Source != "/folder/current" {
		t.Errorf("expected source '/folder/current', got %q", r.Links[0].Source)
	}
	if r.Links[0].Target != "/folder/sibling" {
		t.Errorf("expected target '/folder/sibling', got %q", r.Links[0].Target)
	}
}

func TestResolveMarkdownLink_NoGraphLinkForExternal(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})
	r.CurrentSource = "/folder/current"

	r.ResolveMarkdownLink("https://example.com")

	if len(r.Links) != 0 {
		t.Errorf("expected no graph links for external URL, got %d", len(r.Links))
	}
}

func TestResolveMarkdownLink_NoGraphLinkForAnchor(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})
	r.CurrentSource = "/folder/current"

	r.ResolveMarkdownLink("#heading")

	if len(r.Links) != 0 {
		t.Errorf("expected no graph links for anchor-only link, got %d", len(r.Links))
	}
}
