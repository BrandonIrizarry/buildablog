package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/rss"
)

func (cfg config) getRSS(w http.ResponseWriter, r *http.Request) {
	siteTitle := "Biome of Ideas"
	siteURL := cfg.SiteURL
	genre := (*new(posts.Frontmatter)).Genre()

	ps, err := readers.AllArticles[posts.Frontmatter](cfg.PublishedDir(genre), nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []rss.Item
	for _, post := range ps {
		date := post.Frontmatter.Date
		link := fmt.Sprintf("%s/posts/%s", siteURL, date.Format(time.DateOnly))
		pubDate := date.Format(time.RFC1123)

		item := rss.Item{
			Title:   post.Frontmatter.Title,
			Link:    link,
			GUID:    link,
			PubDate: pubDate,
			Description: rss.Description{
				Type: "html",
				Text: post.Content,
			},
		}

		items = append(items, item)
	}

	image := rss.Image{
		Title:  siteTitle,
		Link:   siteURL,
		URL:    fmt.Sprintf("%s/static/bitmap.png", siteURL),
		Width:  80,
		Height: 80,
	}

	rssChannel := rss.Channel{
		Title:       siteTitle,
		Link:        siteURL,
		Description: "My personal website and blog",
		Language:    "en-us",
		Image:       image,
		Items:       items,
	}

	mainRSS := rss.RSS{
		Channel: rssChannel,
		Version: "2.0",
	}

	// Marshal the data to XML
	feed, err := xml.MarshalIndent(mainRSS, "", strings.Repeat(" ", 4))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, xml.Header+string(feed))
}
