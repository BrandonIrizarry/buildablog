package main

import (
	"html/template"
)

type reader func(string, string) (string, error)

type postData struct {
	Title   string `toml:"title"`
	Date    string `toml:"date"`
	Content template.HTML
	Author  author `toml:"author"`
}

type author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}
