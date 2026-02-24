package rss

// Channel is used for marshalling data into the blog's RSS feed.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Image       Image  `xml:"image"`
	Items       []Item `xml:"item"`
}

// Item is used to enumerate a blog post's mention in the RSS feed.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// Image is used to display an image when aggregators present a field.
type Image struct {
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	URL    string `xml:"url"`
	Width  int    `xml:"width"`
	Height int    `xml:"height"`
}
