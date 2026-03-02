package types

import (
	"html/template"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/projects"
)

type Frontmatter interface {
	posts.Frontmatter | projects.Frontmatter
	GetDate() time.Time
	Genre() string
}

type Article[F Frontmatter] struct {
	Frontmatter F
	Content     template.HTML
}
