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

Requests are served on `localhost` on a VPS. Nginx then serves the
content to the Web via reverse-proxy; the site is currently viewable
at <https://brandonirizarry.xyz>.

The blog itself lives in a separate directory on the same VPS
filesystem. The blog itself has a peculiar layout which the server
expects to see:

```
blog/
    assets/
        <site-wide images appearing in the various posts>
    .git
    .gitignore
    index/
        <the site's front page>
    posts/
        <blog posts>
    projects/
        <project posts>
```

A post is served whenever its date frontmatter field has been filled
out. Internally, the SSG looks for a non-zero value of the date's
corresponding `time.Time` value.

# Configuration

The SSG can be configured with an `.env` file at the project
top-level. An example `.env` might look like this:

```bash
BLOGDIR=/home/bci/blog
SITEURL="https://brandonirizarry.xyz"
PORT="3030"
```

## BLOGDIR

Used to identify the root directory of the user's blog content.

## SITEURL

This is used mainly for testing the generated RSS feed locally.

## PORT

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

