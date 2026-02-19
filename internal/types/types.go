package types

// frontmatter fields
type FrontmatterData struct {
	Title   string   `toml:"title"`
	Summary string   `toml:"summary"`
	Tags    []string `toml:"tags"`
	Publish bool     `toml:"publish"`
}

// inside published.json
type PublishData struct {
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`

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
