package main

import "github.com/egoist/mygo/ui"

func icon(path string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">` + path + `</svg>`))
}

var icons = map[string]*ui.SVG{
	"undo":   icon(`<path d="M3 7v6h6"/><path d="M3 13c3-5 12-7 17-1 1 1 1 3 1 5"/>`),
	"redo":   icon(`<path d="M21 7v6h-6"/><path d="M21 13C18 8 9 6 4 12c-1 1-1 3-1 5"/>`),
	"bold":   icon(`<path d="M6 12h7a4 4 0 0 1 0 8H6zM6 4h6a4 4 0 0 1 0 8H6z"/>`),
	"italic": icon(`<path d="M19 4h-9M14 20H5M15 4 9 20"/>`),
	"code":   icon(`<path d="m8 5-6 7 6 7M16 5l6 7-6 7"/>`),
	"link":   icon(`<path d="M10 13a5 5 0 0 0 7.5.5l3-3a5 5 0 0 0-7-7l-1.7 1.7"/><path d="M14 11a5 5 0 0 0-7.5-.5l-3 3a5 5 0 0 0 7 7l1.7-1.7"/>`),
	"image":  icon(`<rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-4-4L5 21"/>`),
	"list":   icon(`<path d="M9 6h12M9 12h12M9 18h12M4 6h.01M4 12h.01M4 18h.01"/>`),
	"tasks":  icon(`<path d="m3 6 1.5 1.5L7 5M10 6h11M10 12h11M10 18h11"/><rect x="3" y="10" width="4" height="4" rx=".5"/><rect x="3" y="16" width="4" height="4" rx=".5"/>`),
	"eye":    icon(`<path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12z"/><circle cx="12" cy="12" r="2.5"/>`),
	"split":  icon(`<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M12 4v16"/>`),
	"copy":   icon(`<rect x="8" y="8" width="12" height="12" rx="2"/><path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2"/>`),
	"plus":   icon(`<path d="M12 5v14M5 12h14"/>`),
	"file":   icon(`<path d="M6 3h8l4 4v14H6zM14 3v5h4"/>`),
}
