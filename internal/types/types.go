package types

import "time"

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
	//
	// DEPRECATED: when generating the /archives page, only list
	// posts with the proper time format (symbolic links.) Hence
	// final inclusion of the date field will ideally signal the
	// publish (idea: a 'make publish' action can generate the
	// appropriate symlinks for posts.)
	Publish bool `toml:"publish"`

	// Date is the date of the post. It looks like TOML will
	// accept time.DateOnly as a value, so let's try it.
	Date time.Time `toml:"date"`
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
