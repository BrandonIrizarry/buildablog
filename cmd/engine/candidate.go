package main

import (
	"errors"
	"strings"
)

// A candidate is the relative path to a blog post, starting from the
// project top-level. For example,
//
// content/posts/eleventy.md
//
// A separate type is needed here for the candidate list so that we
// can pass it in as a command-line argument. The list elements are
// separated by a comma (,), no spaces.
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
		return errors.New("flag already set")
	}

	for c := range strings.SplitSeq(value, ",") {
		*cl = append(*cl, c)
	}

	return nil
}
