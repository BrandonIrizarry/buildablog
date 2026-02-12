package main

import (
	"html/template"
)

type frontmatterData struct {
	Title  string `toml:"title"`
	Date   string `toml:"date"`
	Author author `toml:"author"`
	Style  style  `toml:"style"`
}

type author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type style struct {
	Syntax string `toml:"syntax"`
}

type postData struct {
	frontmatterData
	Content template.HTML
}
