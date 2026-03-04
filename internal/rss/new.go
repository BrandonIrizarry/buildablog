package rss

import (
	"html/template"
	"time"
)

func NewItem(title, link, guid string, text template.HTML, date time.Time) Item {
	pubDate := date.Format(time.RFC1123)

	return Item{
		Title:   title,
		Link:    link,
		GUID:    guid,
		PubDate: pubDate,
		Description: Description{
			Type: "html",
			Text: text,
		},

		date: date,
	}
}
