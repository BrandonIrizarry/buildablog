package types

import (
	"html/template"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/genres/index"
	"github.com/BrandonIrizarry/buildablog/internal/genres/posts"
	"github.com/BrandonIrizarry/buildablog/internal/genres/projects"
)

type Frontmatter interface {
	posts.Frontmatter | projects.Frontmatter | index.Frontmatter
	GetDate() time.Time
	GetTitle() string
	Genre() string
}

type Article[F Frontmatter] struct {
	Frontmatter F
	Content     template.HTML
}
