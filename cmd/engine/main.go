package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
	"github.com/google/uuid"
)

const atomFeedFile = "atom.xml"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var candidates candidatesList

	// The reset flag for now lets us test writing the XML
	// file. It looks like it won't be of much use later on
	// though.
	var reset bool

	flag.Var(&candidates, "candidates", "List of blog posts to update")
	flag.BoolVar(&reset, "reset", false, "Start over with a blank atom feed file.")
	flag.Parse()

	if len(candidates) > 0 {
		if err := updateCandidates(candidates); err != nil {
			log.Fatal(err)
		}
	}

	// Handle the "reset" flag.
	if reset {
		if err := bootstrapAtomXMLFile(); err != nil {
			log.Fatal(err)
		}
	}
}

type AtomLink struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type AtomAuthor struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

type AtomCategory struct {
	Term string `xml:"term,attr"`
}

type AtomContent struct {
	Type string `xml:"type,attr"`
}

type AtomEntry struct {
	Title      string         `xml:"title"`
	Link       AtomLink       `xml:"link"`
	ID         string         `xml:"id"`
	Updated    string         `xml:"updated"`
	Categories []AtomCategory `xml:"category"`
	Content    string         `xml:"content"`
}

type AtomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Title   string     `xml:"title"`
	Links   []AtomLink `xml:"link"`
	Updated string     `xml:"updated"`
	ID      string     `xml:"id"`
	Author  AtomAuthor `xml:"author"`
}

// bootstrapAtomXMLFile sets up atom.xml for the first time with the
// site's metadata.
func bootstrapAtomXMLFile() error {
	atomFeed := AtomFeed{
		Title: "brandonirizarry.xyz",
		Links: []AtomLink{
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: "https://brandonirizarry.xyz",
			},

			{
				Rel:  "self",
				Type: "application/atom+xml",
				Href: "https://brandonirizarry.xyz/feed",
			},
		},
		Updated: time.Now().Format(time.RFC3339),
		ID:      uuid.New().URN(),
		Author: AtomAuthor{
			Name:  "Brandon Irizarry",
			URI:   "https://brandonirizarry.xyz",
			Email: "brandon.irizarry@gmail.com",
		},
	}

	xmlBytes, err := xml.MarshalIndent(atomFeed, "", strings.Repeat(" ", 4))
	if err != nil {
		return fmt.Errorf("can't marshal atom feed struct: %w", err)
	}

	if err := os.WriteFile(atomFeedFile, xmlBytes, 0644); err != nil {
		return fmt.Errorf("can't write atom feed file '%s': %w", atomFeedFile, err)
	}

	return nil
}

// updateAtomXML updates the Atom feed for this blog. This XML file is
// meant to replace the eariler "publish.json" file used.
func updateAtomXML(candidates candidatesList) error {
	log.Printf("Candidates are: %v", candidates)

	// We trust that Make has provided us only with those
	// candidates that have been recently edited.
	candidateSlugSet := make(map[string]struct{})
	for _, c := range candidates {
		slug := strings.TrimSuffix(filepath.Base(c), ".md")
		candidateSlugSet[slug] = struct{}{}
	}

	log.Printf("Candidate set is now: %v", candidateSlugSet)

	// If the atom feed file doesn't exist, create a new one.
	if _, err := os.Stat(atomFeedFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := bootstrapAtomXMLFile(); err != nil {
				return fmt.Errorf("can't write new %s: %w", atomFeedFile, err)
			}
		} else {
			return fmt.Errorf("can't stat %s: %w", atomFeedFile, err)
		}
	}

	return nil
}

// updateCandidates revises the contents of published.json to reflect
// any updates to a blog-post's front matter.
func updateCandidates(candidates candidatesList) error {
	log.Printf("Candidates are: %v", candidates)

	// We trust that Make has provided us only with those
	// candidates that have been recently edited.
	candidateSlugSet := make(map[string]struct{})
	for _, c := range candidates {
		slug := strings.TrimSuffix(filepath.Base(c), ".md")
		candidateSlugSet[slug] = struct{}{}
	}

	log.Printf("Candidate set is now: %v", candidateSlugSet)

	const publishedFile = "published.json"

	// If publishedFile doesn't exist, create a new one whose sole
	// contents are "[]". This makes it a valid JSON data
	// structure which we can unmarshal later on.
	if _, err := os.Stat(publishedFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.WriteFile(publishedFile, []byte("[]"), 0644); err != nil {
				return fmt.Errorf("can't write new %s: %w", publishedFile, err)
			}
		} else {
			return fmt.Errorf("can't stat %s: %w", publishedFile, err)
		}
	}

	log.Printf("Publishing file %s exists", publishedFile)

	// Read the current published data into a slice of
	// [types.PublishData]. By now publishedFile should already
	// exist on disk.
	f, err := os.Open(publishedFile)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer f.Close()

	log.Printf("Opened publishing file %s successfully", publishedFile)

	fileContent, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("can't read file: %w", err)
	}

	log.Printf("Successfully read %s", publishedFile)

	// alreadyPublished represents the current contents of
	// published.json.
	var alreadyPublished []types.PublishData
	if err := json.Unmarshal(fileContent, &alreadyPublished); err != nil {
		return fmt.Errorf("can't unmarshal: %w", err)
	}

	var putItBack []types.PublishData
	for _, p := range alreadyPublished {
		if _, ok := candidateSlugSet[p.Slug]; ok {
			// UPDATED POSTS
			log.Printf("Caught %s as having been edited", p.Slug)
			data, _, err := readers.ReadPage(p.Slug, "posts")
			if err != nil {
				return fmt.Errorf("can't read candidate slug '%s' (updated post): %w", p.Slug, err)
			}

			if !data.Publish {
				// This means the post was demoted to
				// "unpublished" status, either by
				// setting the frontmatter 'publish'
				// bool to false, or else by removing
				// that line entirely.
				log.Printf("Revoked post '%s' (%s)", p.Slug, data.Title)
				continue
			}

			// Right now we support only editing the
			// title, though of course I plan on adding
			// more stuff here soon.
			newPublishedData := types.PublishData{
				// These two fields are simply
				// propagated across each update.
				Created: p.Created,
				Slug:    p.Slug,

				// The remaining fields are either
				// modified by necessity, or else are
				// suitable to have been modified by
				// the user.
				Updated: time.Now().Unix(),
				Title:   data.Title,
				Summary: data.Summary,
				Tags:    data.Tags,
			}

			putItBack = append(putItBack, newPublishedData)

			// Remove this slug, since the ones which will
			// remain at the end are precisely newly
			// published files.
			delete(candidateSlugSet, p.Slug)
		} else {
			// INERT POSTS
			log.Printf("The file having slug %s wasn't edited recently", p.Slug)
			putItBack = append(putItBack, p)
		}
	}

	// candidateSlugSet should now only refer to those candidates
	// that weren't already listed in published.json.
	//
	// NEW POSTS:
	for slug := range candidateSlugSet {
		data, _, err := readers.ReadPage(slug, "posts")
		if err != nil {
			return fmt.Errorf("can't read candidate slug '%s' (new post): %w", slug, err)
		}

		if !data.Publish {
			log.Printf("post '%s' (%s) isn't marked for publishing; skip", slug, data.Title)
			continue
		}

		newPublishedData := types.PublishData{
			Created: time.Now().Unix(),
			Updated: time.Now().Unix(),
			Slug:    slug,
			Title:   data.Title,
			Summary: data.Summary,
			Tags:    data.Tags,
		}

		putItBack = append(putItBack, newPublishedData)
	}

	newContent, err := json.MarshalIndent(putItBack, "", strings.Repeat(" ", 4))
	if err != nil {
		return fmt.Errorf("can't marshal updated published content: %w", err)
	}

	if err := os.WriteFile(publishedFile, newContent, 0644); err != nil {
		return fmt.Errorf("can't write updated published content to")
	}

	return nil
}
