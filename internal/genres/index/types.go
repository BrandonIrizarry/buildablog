package index

import "time"

type Frontmatter struct {
	Title string    `toml:"title"`
	Date  time.Time `toml:"date"`
}

func (f Frontmatter) GetDate() time.Time {
	return f.Date
}

func (f Frontmatter) GetTitle() string {
	return f.Title
}

func (f Frontmatter) Genre() string {
	return "index"
}
