package atom

import "encoding/xml"

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Author struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

type Category struct {
	Term string `xml:"term,attr"`
}

type Content struct {
	Type string `xml:"type,attr"`
}

type Entry struct {
	Title      string     `xml:"title"`
	Link       Link       `xml:"link"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Categories []Category `xml:"category"`
	Content    string     `xml:"content"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Links   []Link   `xml:"link"`
	Updated string   `xml:"updated"`
	ID      string   `xml:"id"`
	Author  Author   `xml:"author"`
}
