package docs

import (
	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/files"
	"github.com/gregoryv/web/toc"
)

// Generate all documentation to the given directory.
func Generate(dir string) error {
	index := NewIndexPage()
	index.Filename = "index.html"
	return index.SaveTo(dir)
}

const gregory = "Gregory Vin&ccaron;i&cacute;"

func NewIndexPage() *Page {
	page := NewPage(
		Html(
			Head(
				Meta(Charset("utf-8")),
				Meta(
					Name("viewport"),
					Content("width=device-width, initial-scale=1"),
				),
				Style(Theme()),
			),
			Body(
				NewProjectArticle(),
			),
		),
	)
	return page
}

func NewProjectArticle() interface{} {
	nav := Nav(H4("Table of contents"))
	article := Article(
		H1("Geco - Golang Code Generators"),

		P(`Generating code solves many otherwise tedious and error
		prone repetitive work. This module provides an API and tools
		to generate various helpful structures.`),

		Div(Class("left"),
			nav,
		),
		Div(Class("right"),
			H2("Install"),
			Pre(
				Code(
					"    go get -u ",
					A(Href("https://github.com/gregoryv/geco"), "github.com/gregoryv/geco"),
					"/...",
				),
			),

			H2("About"),

			Img(Src("me_circle.png"), Class("me")),
			P(
				`Written by `, A(Href("https://github.com/gregoryv"), gregory), Br(),
				A(Href("#license"), "MIT License"),
			),
		),
		Br(Attr("clear", "all")),
		H2("License"),
		Pre(files.MustLoad("../LICENSE")),
	)
	toc.MakeTOC(nav, article, "h2", "h3")
	return article
}
