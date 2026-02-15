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

	// Read the current published data into a slice of
	// [types.PublishData]. By now publishedFile should already
	// exist on disk.
	f, err := os.OpenFile(publishedFile, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer f.Close()

	fileContent, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("can't read file: %w", err)
	}

	var entries []types.PublishData
	if err := json.Unmarshal(fileContent, &entries); err != nil {
		return fmt.Errorf("can't unmarshal: %w", err)
	}

	// GNU Make here looks like it reserves stdout for itself, so
	// let's use stderr.
	fmt.Fprintln(os.Stderr, entries)

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
