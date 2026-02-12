package types

import (
	"html/template"
)

type FrontmatterData struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Date    string `toml:"date"`
	Author  author `toml:"author"`
	Style   style  `toml:"style"`
}

type author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type style struct {
	Syntax string `toml:"syntax"`
}

type PostData struct {
	FrontmatterData
	Content template.HTML
}
