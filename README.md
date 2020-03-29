# rapidtemplate
Rapid Template is a small, server-agnostic, Markdown to HTML templating application written in Go.

I find myself frequently creating static websites for my projects. This is fine, but keeping track of common HTML between pages can be a chore. There are tons of static site generators such as Jekyll, but I wanted a easy and portable solution. Rapid Template is just a single binary that quickly generates static HTML from Markdown files.

### Setup

1. Place your Markdown files in `pages/` and static assets in `out/`
2. Edit `template.html` to include your header, footer, scripts, and stylesheets.
3. Build and run rapidtemplate. `go run rapidtemplate.go`

Rapid Template will build all the files in `pages/` and watch for changes, updating accordingly.