package types

import "html/template"

type FrontmatterData struct {
	Title   string   `toml:"title"`
	Summary string   `toml:"summary"`
	Publish bool     `toml:"publish"`
	Tags    []string `toml:"tags"`
}

type PostData struct {
	FrontmatterData
	Content template.HTML
}

// PublishData represents a blog post in the archive file
// (published.json).
type PublishData struct {
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

	// Title is (by convention) the proper title of the post.
	Title string `json:"title"`

	// Summary is (by convention) a blurb summary of the post,
	// used in the archives listing.
	Summary string `json:"summary"`

	// Tags: a post's list of tags.
	Tags []string `json:"tags"`
}
