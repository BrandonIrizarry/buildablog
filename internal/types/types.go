package types

// FrontmatterData is used for unmarshalling blog post frontmatter.
type FrontmatterData struct {
	// Title is the blog post's title.
	Title string `toml:"title"`

	// Summary is a short summary of the blog post.
	Summary string `toml:"summary"`

	// Tags is the list of tags categorizing the blog post.
	Tags []string `toml:"tags"`

	// Publish flags whether to publish the blog post into the
	// publishing file.
	Publish bool `toml:"publish"`
}

// PublishData is used from marshalling/unmarshalling the currently
// published archives.
//
// Note that some fields are repeated from [FrontmatterData]; this is
// done to avoid obfuscating the concerns of a given struct type
// solely for the sake of DRY.
type PublishData struct {
	// Title is the blog post's title.
	Title string `json:"title"`

	// Summary is a short summary of the blog post.
	Summary string `json:"summary"`

	// Tags is the list of tags categorizing the blog post.
	Tags []string `json:"tags"`

	// Created details when the post was created (Unix timestamp.)
	Created int64 `json:"created"`

	// Updated details when the post was last modified (Unix
	// timestamp.) Naturally, when a post is first created, this
	// should match [PublishData.Created].
	Updated int64 `json:"modified"`

	// Slug is the post's given filename minus the file
	// extension. The various templates use this to generate the
	// links to the posts.
	Slug string `json:"slug"`
}

/* The following structs are for handling RSS feeds. Their fields are
   meant to conform to the RSS specification. */

// Channel is used for marshalling data into the blog's RSS feed.
type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Image       RSSImage  `xml:"image"`
	Items       []RSSItem `xml:"item"`
}

// Item is used to enumerate a blog post's mention in the RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// Image is used to display an image when aggregators present a field.
type RSSImage struct {
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	URL    string `xml:"url"`
	Width  int    `xml:"width"`
	Height int    `xml:"height"`
}
