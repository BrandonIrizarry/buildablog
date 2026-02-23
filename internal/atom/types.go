package atom

import "encoding/xml"

type AtomLink struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type AtomAuthor struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

type AtomCategory struct {
	Term string `xml:"term,attr"`
}

type AtomContent struct {
	Type string `xml:"type,attr"`
}

type AtomEntry struct {
	Title      string         `xml:"title"`
	Link       AtomLink       `xml:"link"`
	ID         string         `xml:"id"`
	Updated    string         `xml:"updated"`
	Categories []AtomCategory `xml:"category"`
	Content    string         `xml:"content"`
}

type AtomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Title   string     `xml:"title"`
	Links   []AtomLink `xml:"link"`
	Updated string     `xml:"updated"`
	ID      string     `xml:"id"`
	Author  AtomAuthor `xml:"author"`
}
