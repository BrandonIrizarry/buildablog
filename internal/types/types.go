package types

import (
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/projects"
)

type Frontmatter interface {
	posts.Frontmatter | projects.Frontmatter
	GetDate() time.Time
}
