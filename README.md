# Introduction

Buildablog is a static site generator for hosting a weblog. It's
implemented as an HTTP server which serves a collection of Markdown
blog posts as HTML static content.

More specifically, it allows the author to publish content in a
variety of formats, though a combination of various frontmatter
layouts and Go templates. For example, there is one for ordinary blog
posts, and one for software projects.

It takes Jon Calhoun's [Building a Blog Exercise](https://www.calhoun.io/building-a-blog-exercise/) as its starting
point, though it has evolved way past that point.

# How It Works

Requests are served on `localhost` (for example, `localhost:3030`) on
a VPS. Nginx then serves the content to the Web via reverse-proxy; the
site is currently viewable at <https://brandonirizarry.xyz>.

## Directory Layout

The blog itself lives in a separate directory on the same VPS
filesystem, which the SSG knows about through an environment variable
(see below.) The blog itself has a peculiar layout which the server
expects to see:

```
blog/
    assets/
        <site-wide images appearing in the various posts>
    index/
        <the site's front page>
    posts/
        <blog posts>
    projects/
        <project posts>
```

## Frontmatter and Publishing

A post is served whenever its date frontmatter field has been filled
out. Internally, the SSG looks for a non-zero value of the date's
corresponding `time.Time` value.

Generics are used heavily to support handling a variety of frontmatter
layouts (represented as structs) without much code duplication. For
example, a post frontmatter section looks like this:

```toml
+++
title = "Adding a CGit Subdomain To My Site"
tags = ["linux", "nginx", "certbot", "cgit"]
summary = "Setting up CGit on my VPS."
date = 2026-03-06
+++
```

A project frontmatter section looks like this:

```toml
+++
name = "buildablog"
title = "Building My Own SSG"
host_url = "https://github.com/BrandonIrizarry/buildablog"
synopsis = "The SSG used to build this site."
stack = ["Go", "HTML", "CSS"]
thumbnail = "assets/github-white.svg"
date = 2026-03-01
+++
```

Adding a new frontmatter type is a matter of adding the requisite
struct type, implementing a few interface functions on it, and then
adding it as a supported type:

```go
type Frontmatter interface {
    // Our registered frontmatter types.
	posts.Frontmatter | projects.Frontmatter | index.Frontmatter
    
    // A few basic interface methods.
	GetDate() time.Time
	GetTitle() string
	Genre() string
}
```

## Configuration

The SSG can be configured with an `.env` file at the project
top-level. An example `.env` might look like this:

```bash
BLOGDIR=/home/bci/blog
SITEURL="https://brandonirizarry.xyz"
PORT="3030"
```

### BLOGDIR

Used to identify the root directory of the user's blog content.

### SITEURL

Used mainly for testing the generated RSS feed locally.

### PORT

Used to specify the port on which to launch the SSG server.


# Copyright and Licensing

Copyright © 2026 Brandon C. Irizarry

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Affero General Public License for more details.

You should have received a [copy](./LICENSE) of the GNU Affero General Public
License along with this program.  If not, see
<https://www.gnu.org/licenses/>.

