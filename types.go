package main

import (
	"html/template"
)

type reader func(string, string) (string, error)

type frontmatterData struct {
	Title  string `toml:"title"`
	Date   string `toml:"date"`
	Author author `toml:"author"`
}

type author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type postData struct {
	frontmatterData
	Content template.HTML
}
