package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/rss"
)

func (cfg rssConfig) getRSS(w http.ResponseWriter, r *http.Request) {
	siteTitle := "Biome of Ideas"
	var siteURL string

	// Use flagLocal to set the correct siteURL for
	// purposes of testing the RSS feed locally with
	// something like newsboat.
	if cfg.flagLocal {
		siteURL = "http://localhost:3030"
	} else {
		siteURL = "https://brandonirizarry.xyz"
	}

	ps, err := readers.AllPosts(nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []rss.Item
	for _, post := range ps {
		link := fmt.Sprintf("%s/posts/%s", siteURL, post.Date.Format(time.DateOnly))
		pubDate := post.Date.Format(time.RFC1123)

		item := rss.Item{
			Title:   post.Title,
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
