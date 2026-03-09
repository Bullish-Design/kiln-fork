// @feature:rss Tests for RSS feed generation.
package rss

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestBuildFeedXML_FullFeed(t *testing.T) {
	pub1 := time.Date(2024, 7, 20, 14, 0, 0, 0, time.UTC)
	pub2 := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	params := FeedParams{
		Title:       "My Blog",
		Link:        "https://example.com",
		Description: "A blog about things",
	}
	items := []ItemParams{
		{
			Title:       "First Post",
			Link:        "https://example.com/blog/first-post",
			Description: "This is the first post",
			PubDate:     pub1,
			GUID:        "https://example.com/blog/first-post",
		},
		{
			Title:       "Second Post",
			Link:        "https://example.com/blog/second-post",
			Description: "This is the second post",
			PubDate:     pub2,
			GUID:        "https://example.com/blog/second-post",
		},
	}

	got, err := BuildFeedXML(params, items)
	if err != nil {
		t.Fatalf("BuildFeedXML returned error: %v", err)
	}

	if !strings.HasPrefix(got, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("expected XML declaration header")
	}

	var rss rssDoc
	if err := xml.Unmarshal([]byte(got), &rss); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if rss.Version != "2.0" {
		t.Errorf("rss version = %q, want 2.0", rss.Version)
	}
	if rss.Channel.Title != "My Blog" {
		t.Errorf("channel title = %q, want My Blog", rss.Channel.Title)
	}
	if rss.Channel.Link != "https://example.com" {
		t.Errorf("channel link = %q, want https://example.com", rss.Channel.Link)
	}
	if rss.Channel.Description != "A blog about things" {
		t.Errorf("channel description = %q, want A blog about things", rss.Channel.Description)
	}
	if rss.Channel.LastBuildDate == "" {
		t.Error("expected lastBuildDate to be set")
	}

	if len(rss.Channel.Items) != 2 {
		t.Fatalf("item count = %d, want 2", len(rss.Channel.Items))
	}

	item0 := rss.Channel.Items[0]
	if item0.Title != "First Post" {
		t.Errorf("item[0] title = %q, want First Post", item0.Title)
	}
	if item0.Link != "https://example.com/blog/first-post" {
		t.Errorf("item[0] link = %q, want https://example.com/blog/first-post", item0.Link)
	}
	if item0.Description != "This is the first post" {
		t.Errorf("item[0] description = %q, want This is the first post", item0.Description)
	}
	if item0.PubDate != pub1.Format(time.RFC1123Z) {
		t.Errorf("item[0] pubDate = %q, want %q", item0.PubDate, pub1.Format(time.RFC1123Z))
	}
	if item0.GUID.Value != "https://example.com/blog/first-post" {
		t.Errorf("item[0] guid = %q, want https://example.com/blog/first-post", item0.GUID.Value)
	}
	if item0.GUID.IsPermaLink != "true" {
		t.Errorf("item[0] guid isPermaLink = %q, want true", item0.GUID.IsPermaLink)
	}

	item1 := rss.Channel.Items[1]
	if item1.Title != "Second Post" {
		t.Errorf("item[1] title = %q, want Second Post", item1.Title)
	}
	if item1.PubDate != pub2.Format(time.RFC1123Z) {
		t.Errorf("item[1] pubDate = %q, want %q", item1.PubDate, pub2.Format(time.RFC1123Z))
	}
}

func TestBuildFeedXML_EmptyItems(t *testing.T) {
	params := FeedParams{
		Title:       "Empty Blog",
		Link:        "https://example.com",
		Description: "Nothing here yet",
	}

	got, err := BuildFeedXML(params, nil)
	if err != nil {
		t.Fatalf("BuildFeedXML returned error: %v", err)
	}

	var rss rssDoc
	if err := xml.Unmarshal([]byte(got), &rss); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if rss.Channel.Title != "Empty Blog" {
		t.Errorf("channel title = %q, want Empty Blog", rss.Channel.Title)
	}
	if len(rss.Channel.Items) != 0 {
		t.Errorf("item count = %d, want 0", len(rss.Channel.Items))
	}
}

func TestBuildFeedXML_EmptyDescription(t *testing.T) {
	params := FeedParams{
		Title:       "My Blog",
		Link:        "https://example.com",
		Description: "A blog",
	}
	items := []ItemParams{
		{
			Title:   "No Desc Post",
			Link:    "https://example.com/no-desc",
			PubDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			GUID:    "https://example.com/no-desc",
		},
	}

	got, err := BuildFeedXML(params, items)
	if err != nil {
		t.Fatalf("BuildFeedXML returned error: %v", err)
	}

	if strings.Contains(got, "<description></description>") {
		t.Error("expected empty description to be omitted from item")
	}

	var rss rssDoc
	if err := xml.Unmarshal([]byte(got), &rss); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("item count = %d, want 1", len(rss.Channel.Items))
	}
	if rss.Channel.Items[0].Description != "" {
		t.Errorf("item description = %q, want empty", rss.Channel.Items[0].Description)
	}
}

func TestBuildFeedXML_SpecialCharacters(t *testing.T) {
	params := FeedParams{
		Title:       "Blog & News",
		Link:        "https://example.com",
		Description: "Posts about <things> & stuff",
	}
	items := []ItemParams{
		{
			Title:       "A < B & C > D",
			Link:        "https://example.com/special",
			Description: "Use <html> & escape \"quotes\"",
			PubDate:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			GUID:        "https://example.com/special",
		},
	}

	got, err := BuildFeedXML(params, items)
	if err != nil {
		t.Fatalf("BuildFeedXML returned error: %v", err)
	}

	// The XML should be parseable (meaning special chars are properly escaped)
	var rss rssDoc
	if err := xml.Unmarshal([]byte(got), &rss); err != nil {
		t.Fatalf("XML with special characters should be valid, got error: %v", err)
	}

	if rss.Channel.Title != "Blog & News" {
		t.Errorf("channel title = %q, want %q", rss.Channel.Title, "Blog & News")
	}
	if rss.Channel.Description != "Posts about <things> & stuff" {
		t.Errorf("channel description = %q, want %q", rss.Channel.Description, "Posts about <things> & stuff")
	}
	if rss.Channel.Items[0].Title != "A < B & C > D" {
		t.Errorf("item title = %q, want %q", rss.Channel.Items[0].Title, "A < B & C > D")
	}
	if rss.Channel.Items[0].Description != `Use <html> & escape "quotes"` {
		t.Errorf("item description = %q, want %q", rss.Channel.Items[0].Description, `Use <html> & escape "quotes"`)
	}
}
