package main

import (
	"strings"
	"unicode/utf8"
)

const largeDocumentLimit = 512 * 1024

func (a *app) previewText() string {
	if len(a.text) < largeDocumentLimit {
		return a.text
	}
	const previewBytes = 16 * 1024
	text := a.text
	if len(text) > previewBytes {
		text = text[:previewBytes]
		if cut := strings.LastIndexByte(text, '\n'); cut > 0 {
			text = text[:cut]
		} else {
			cut := previewBytes
			for cut > 0 && !utf8.RuneStart(a.text[cut]) {
				cut--
			}
			text = a.text[:cut]
		}
	}
	return text + "\n\n> Large file preview is limited to the first 16 KB."
}
