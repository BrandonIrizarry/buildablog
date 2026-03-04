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
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func rssItems[F types.Frontmatter](siteURL string, articles []types.Article[F]) []rss.Item {
	genre := (*new(posts.Frontmatter)).Genre()

	var items []rss.Item
	for _, article := range articles {
		date := article.Frontmatter.GetDate()
		pubDate := date.Format(time.RFC1123)
		link := fmt.Sprintf("%s/%s/%s", siteURL, genre, date.Format(time.DateOnly))

		item := rss.Item{
			Title:   article.Frontmatter.GetTitle(),
			Link:    link,
			GUID:    link,
			PubDate: pubDate,
			Description: rss.Description{
				Type: "html",
				Text: article.Content,
			},
		}

		items = append(items, item)
	}

	return items
}

func (cfg config) getRSS(w http.ResponseWriter, r *http.Request) {
	siteTitle := "Biome of Ideas"
	siteURL := cfg.SiteURL
	genre := (*new(posts.Frontmatter)).Genre()

	ps, err := readers.AllArticles[posts.Frontmatter](cfg.PublishedDir(genre), nil, cfg.Timezone)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := rssItems(siteURL, ps)

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
