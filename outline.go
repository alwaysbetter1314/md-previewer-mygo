package main

import (
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui/markdown/parser"
)

type heading struct {
	level, offset, block int
	text                 string
}

func (a *app) headings() []heading {
	source := a.previewText()
	if a.outlineValid && a.outlineText == source {
		return a.outline
	}
	doc := parser.Parse(source)
	var headings []heading
	for i := 1; i < doc.Len(); i++ {
		block := doc.Block(i)
		if block.Kind != parser.Heading {
			continue
		}
		text := strings.TrimSpace(parser.PlainText(doc.Inlines(i)))
		headings = append(headings, heading{block.Level, utf8.RuneCountInString(source[:block.Start]), i, text})
	}
	a.outline, a.outlineText, a.outlineValid = headings, source, true
	return headings
}

func (a *app) jumpToHeading(h heading) {
	a.anchorBlock = h.block
	a.md.Anchor = "outline"
	a.md.AnchorVersion++
	a.edit = editAction{jump: true, offset: h.offset}
	a.status = "已定位到：" + h.text
	a.skipSync = 2
}
