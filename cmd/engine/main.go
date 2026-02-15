package main

import (
	"encoding/json"
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

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

type candidatesList []string

// String implements the [flag.Value] interface.
func (cl *candidatesList) String() string {
	var printed strings.Builder

	for _, c := range *cl {
		printed.WriteString(" " + c)
	}

	return printed.String()
}

// Set implements the [flag.Value] interface.
func (cl *candidatesList) Set(value string) error {
	if len(*cl) > 0 {
		return errors.New("publishall flag already set")
	}

	for c := range strings.SplitSeq(value, ",") {
		*cl = append(*cl, c)
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		candidate  string
		candidates candidatesList
	)

	flag.StringVar(&candidate, "publish", "", "Publish this post")
	flag.Var(&candidates, "candidates", "List of blog posts to update")
	flag.Parse()

	if candidate != "" {
		publish(candidate)
	}

	if len(candidates) > 0 {
		if err := updateCandidates(candidates); err != nil {
			log.Fatal(err)
		}
	}
}

func updateCandidates(candidates candidatesList) error {
	log.Printf("Candidates are: %v", candidates)

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

	log.Printf("Current contents of %s: %s", publishedFile, fileContent)

	var alreadyPublished []types.PublishData
	if err := json.Unmarshal(fileContent, &alreadyPublished); err != nil {
		return fmt.Errorf("can't unmarshal: %w", err)
	}

	//  By now, the existing published entries have been loaded
	//  into memory. We trust that Make has provided us only with
	//  those candidates that have been recently edited.
	candidateSlugSet := make(map[string]struct{})
	for _, c := range candidates {
		slug := strings.TrimSuffix(filepath.Base(c), ".md")
		candidateSlugSet[slug] = struct{}{}
	}

	log.Printf("Candidate set is now: %v", candidateSlugSet)

	var putItBack []types.PublishData
	for _, p := range alreadyPublished {
		if _, ok := candidateSlugSet[p.Slug]; ok {
			// UPDATED POSTS
			log.Printf("Caught %s as having been edited", p.Slug)
			data, err := readers.ReadPage(p.Slug, "posts")
			if err != nil {
				return fmt.Errorf("can't read candidate slug '%s' (updated post): %w", p.Slug, err)
			}

			// Right now we support only editing the
			// title, though of course I plan on adding
			// more stuff here soon.
			newPublishedData := types.PublishData{
				Date:  p.Date,
				Slug:  p.Slug,
				Title: data.Title,
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

	// Here should appear only new posts.
	// NEW POSTS:
	for slug := range candidateSlugSet {
		data, err := readers.ReadPage(slug, "posts")
		if err != nil {
			return fmt.Errorf("can't read candidate slug '%s' (new post): %w", slug, err)
		}

		newPublishedData := types.PublishData{
			Date:  time.Now().Format(time.DateOnly),
			Slug:  slug,
			Title: data.Title,
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

// publish publishes the given candidate.
//
// This marshals the title and description into a file, which the
// /archives endpoint reads and uses to define its content.
func publish(candidate string) {
	// We want to take advantage of terminal auto-completion,
	// since blog slugs are often long and bespoke, often closely
	// mirroring their actual titles. Because of this, candidate
	// is read in as a relative path starting from
	// [constants.ContentDirName]. But now we need to do extra
	// work to split this apart.
	dir, file := filepath.Split(candidate)
	if !strings.HasPrefix(dir, constants.ContentDirName) {
		log.Fatalf("'%s' isn't inside %s", dir, constants.ContentDirName)
	}

	label := strings.TrimPrefix(dir, constants.ContentDirName)
	slug := strings.TrimSuffix(file, ".md")

	data, err := readers.ReadPage(slug, label)
	if err != nil {
		log.Fatalf("couldn't read content/posts/%s: %v", candidate, err)
	}

	if !data.Publish {
		log.Printf("'%s' is a draft, and so will not be published", candidate)
		os.Exit(0)
	}

	b, err := json.Marshal(types.PublishData{
		Date:    time.Now().Format(time.DateOnly),
		Slug:    slug,
		Title:   data.Title,
		Summary: data.Summary,
	})
	if err != nil {
		log.Fatalf("couldn't marshal title and summary: %v", err)
	}

	f, err := os.OpenFile("published", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("couldn't open 'published' file: %v", err)
	}
	defer f.Close()

	// Append a newline for human-readability purposes, as well as
	// for the 'GET /archives' endpoint to be able to read this
	// file line by line.
	b = append(b, '\n')

	_, err = f.Write(b)
	if err != nil {
		log.Fatalf("couldn't write to 'published' file: %v", err)
	}
}
