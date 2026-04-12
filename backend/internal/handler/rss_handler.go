package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// RSSHandler handles RSS feed generation.
type RSSHandler struct {
	postService *service.PostService
	siteURL     string
}

// NewRSSHandler creates a new RSSHandler.
func NewRSSHandler(svc *service.PostService, siteURL string) *RSSHandler {
	return &RSSHandler{postService: svc, siteURL: siteURL}
}

type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func excerpt(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "…"
}

// Feed handles GET /api/v1/feed.xml.
func (h *RSSHandler) Feed(w http.ResponseWriter, r *http.Request) {
	posts, _, err := h.postService.ListPublished(model.ListOptions{Page: 1, PerPage: 20})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        fmt.Sprintf("%s/blog/%s", h.siteURL, p.Slug),
			Description: excerpt(p.Content, 300),
			PubDate:     p.CreatedAt.Format(time.RFC1123Z),
			GUID:        fmt.Sprintf("%s/blog/%s", h.siteURL, p.Slug),
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         "neo",
			Link:          h.siteURL,
			Description:   "personal web garden",
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(feed)
}
