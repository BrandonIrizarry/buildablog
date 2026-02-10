package main

import "html/template"

type reader func(string) (string, error)

type postData struct {
	Title, Author string
	Content       template.HTML
}
