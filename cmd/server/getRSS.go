package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/projects"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/rss"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func rssItems[F types.Frontmatter](siteURL string, articles []types.Article[F]) []rss.Item {
	genre := (*new(F)).Genre()

	var items []rss.Item
	for _, article := range articles {
		// Define these up front for readability
		title := article.Frontmatter.GetTitle()
		date := article.Frontmatter.GetDate()
		link := fmt.Sprintf("%s/%s/%s", siteURL, genre, date.Format(time.DateOnly))
		content := article.Content

		item := rss.NewItem(
			title,
			link,
			link,
			content,
			date,
		)

		items = append(items, item)
	}

	return items
}

func (cfg config) getRSS(w http.ResponseWriter, r *http.Request) {
	siteTitle := "Biome of Ideas"
	siteURL := cfg.SiteURL
	var genre string

	// Scan all posts.
	genre = (*new(posts.Frontmatter)).Genre()
	ps, err := readers.AllArticles[posts.Frontmatter](cfg.PublishedDir(genre), nil, cfg.Timezone)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	psItems := rssItems(siteURL, ps)

	// Scan all projects.
	genre = (*new(projects.Frontmatter)).Genre()
	projs, err := readers.AllArticles[projects.Frontmatter](cfg.PublishedDir(genre), nil, cfg.Timezone)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projsItems := rssItems(siteURL, projs)

	// Gather all RSS items in one place and then sort them by
	// date.
	items := []rss.Item{}
	items = append(items, psItems...)
	items = append(items, projsItems...)

	slices.SortFunc(items, func(item1 rss.Item, item2 rss.Item) int {
		return item1.Date().Compare(item2.Date())
	})

	// List RSS items in descending date order.
	slices.Reverse(items)

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
