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

## Deployment

Thanks to Go's relatively painless package system, deployment is as
simple as

`go install github.com/BrandonIrizarry/buildablog/cmd/server@<latest commit>`

A `systemd` service, `buildablog.service`, ensures that the Buildablog
server will restart on reboot:

```desktop
[Unit]
Description=buildablog
After=network.target

[Service]
User=bci
ExecStart=/home/bci/go/bin/server
WorkingDirectory=/home/bci

[Install]
WantedBy=multi-user.target
```

Updating the live version of Buildablog therefore amounts to
performing a `go install` with the latest commit, followed by
restarting the buildablog service. Note that this is currently a
manual process.

The service file isn't a part of the repo. When testing Builablog
locally, I restart the server using a Makefile, which is included in
the repo.

## Serving Requests

Requests are served on `localhost` on a VPS. Nginx then serves the
content to the Web via reverse-proxy. The site is currently viewable
at <https://brandonirizarry.xyz>.

## Directory Layout

The blog content itself lives in a separate directory, which the SSG
knows about through an environment variable. The blog itself has a
peculiar layout which the server expects to see:

```dircolors
blog/
    assets/
        image1.png
        image2.jpg
        etc.
    index/
        site-front-page.md
    posts/
        post1.md
        post2.md
        etc.
    projects/
        project1.md
        project2.md
        etc.
```

## Frontmatter

Generics are used to support handling a variety of frontmatter layouts
without much code duplication. The project refers to these various
layouts as *genres*, since they ultimately define the end-purpose of
the corresponding post.

For example, blog posts use a frontmatter section that looks like
this:

```toml
+++
title = "Adding a CGit Subdomain To My Site"
tags = ["linux", "nginx", "certbot", "cgit"]
summary = "Setting up CGit on my VPS."
date = 2026-03-06
+++
```

I decided to add the concept of a *project* post, which is more or
less like a blog post except that it's meant to showcase an entry in
my projects portfolio. As such, it uses a different set of frontmatter
fields than blog posts, which the example below demonstrates:

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

The frontmatter interface, used as the generic type in this case,
looks like this:

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

As implied by this example, adding a new frontmatter type is a matter
of adding the requisite struct type, implementing a few interface
functions on it, and then adding it as a supported type.

## Articles

The concept of an *article* subsumes the various genres of post: blog
post, project post, and whatever else you define! It simply wraps the
generic frontmatter type with the content itself, which is always of
type `template.HTML`:

```go
type Article[F Frontmatter] struct {
	Frontmatter F
	Content     template.HTML
}
```

The various server REST endpoints, at their core, simply unmarshal
post content into `Article` structs one way or another, and then feed
these to the corresponding Go template.

## Publishing

A post is served whenever its date frontmatter field has been filled
out. Internally, the SSG looks for a non-zero value of the date's
corresponding `time.Time` value.

This is helpful for quickly viewing a draft post locally, so that I
can, for example, verify that CSS styling is being applied correctly.

## RSS

RSS turned out to be surprisingly easy to implement, once you learn
the ropes. These two resources were helpful in pointing me in the
right direction:

1. [How to Create an RSS Feed](https://www.wikihow.com/Create-an-RSS-Feed)

2. [Build Your Own RSS Feed Generator in Go](https://www.youtube.com/watch?v=b2E1JpC38Pg) 

## Configuration

The SSG can be configured with an `.env` file at the project
top-level. An example `.env` might look like this:

```bash
BLOGDIR=/home/bci/blog
SITEURL="https://brandonirizarry.xyz"
PORT="3030"
```

## Frontend

As mentioned earlier, when a request for a page is received, the SSG
fetches some Markdown, converts it into a Go data structure, and then
feeds it into the appropriate Go template. The template in turn is
then styled by some hand-written plain CSS, much of it taken from
other blogs on the Web, especially <https://maurycyz.com/>.

### BLOGDIR

Used to identify the root directory of the user's blog content.

### SITEURL

Used mainly for testing the generated RSS feed locally.

### PORT

Used to specify the port on which to launch the SSG server.

# Motivation

My blog was initially generated using [Eleventy](https://www.11ty.dev/). When revisting my
blog after a hiatus, I had actually forgotten how to deploy my blog! I
used this as an excuse to try out a different SSG. When searching my
Web for one, I ran into Calhoun's tutorial and found it really
interesting, and so starting chiseling away at it. I had already
become of the opinion that writing my own SSG is only marginally more
complicated than figuring out how to use an existing one, and comes
with the benefit of infinitely flexible customization. So I plunged
head on, and now I have a blog again, whose inner workings and
conventions I understand 360.

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

