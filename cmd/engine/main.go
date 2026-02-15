package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
		updateCandidates(candidates)
	}
}

func updateCandidates(candidates candidatesList) {
	fmt.Println(candidates)
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
