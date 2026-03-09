// @feature:rss RSS 2.0 feed generation for blog-style vaults.
package rss

import (
	"encoding/xml"
	"time"
)

// FeedParams holds the channel-level configuration for the RSS feed.
type FeedParams struct {
	Title       string // Site name (from config)
	Link        string // Base URL of the site
	Description string // Site description (use SiteName + " RSS Feed" as fallback)
}

// ItemParams holds the data for a single RSS item.
type ItemParams struct {
	Title       string    // Item title
	Link        string    // Absolute URL (baseURL + webPath)
	Description string    // Plain-text description (stripped markdown, truncated)
	PubDate     time.Time // Created time from the file
	GUID        string    // Same as Link (permalink)
}

// rssDoc is the root XML element for RSS 2.0.
type rssDoc struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description,omitempty"`
	PubDate     string  `xml:"pubDate"`
	GUID        rssGUID `xml:"guid"`
}

type rssGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// BuildFeedXML generates a complete RSS 2.0 XML document as a string.
// It takes channel-level params and a slice of items.
// Items should be pre-sorted by PubDate descending (newest first).
// Returns the XML string and any error from marshaling.
func BuildFeedXML(params FeedParams, items []ItemParams) (string, error) {
	rssItems := make([]rssItem, 0, len(items))
	for _, item := range items {
		rssItems = append(rssItems, rssItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     item.PubDate.Format(time.RFC1123Z),
			GUID: rssGUID{
				IsPermaLink: "true",
				Value:       item.GUID,
			},
		})
	}

	doc := rssDoc{
		Version: "2.0",
		Channel: rssChannel{
			Title:         params.Title,
			Link:          params.Link,
			Description:   params.Description,
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         rssItems,
		},
	}

	data, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}

	return `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(data), nil
}
