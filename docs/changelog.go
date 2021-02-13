package docs

import . "github.com/gregoryv/web"

func NewChangelog() *Element {
	return Wrap(
		H1("Changelog"),

		P(`All notable changes to this project will be documented in
		this file. Project adheres to semantic versioning(v2).`),

		Section(
			H2("[unreleased]"),
			Ul(
				Li("Add TypeUnderTest generator"),
			),
		),
	)
}
