package main

import "math"

func (a *app) syncScroll() {
	ed, pv := &a.editorScroll, &a.previewScroll
	if a.skipSync > 0 {
		a.skipSync--
	} else if !a.syncOff && !a.previewOnly && !a.previewHidden {
		edMoved := math.Abs(float64(ed.Y-a.lastEditorY)) > 0.5
		pvMoved := math.Abs(float64(pv.Y-a.lastPreviewY)) > 0.5
		switch {
		case edMoved && !pvMoved && ed.MaxY > 0 && pv.MaxY > 0:
			pv.Y = min(pv.MaxY, max(0, ed.Y/ed.MaxY*pv.MaxY))
		case pvMoved && !edMoved && pv.MaxY > 0 && ed.MaxY > 0:
			ed.Y = min(ed.MaxY, max(0, pv.Y/pv.MaxY*ed.MaxY))
		}
	}
	a.lastEditorY, a.lastPreviewY = ed.Y, pv.Y
}

func (a *app) pageScroll(down bool) {
	direction, height, scroll := float32(-1), a.previewHeight, &a.previewScroll
	if down {
		direction = 1
	}
	if a.previewHidden && !a.previewOnly {
		height, scroll = a.editorHeight, &a.editorScroll
	}
	scroll.Y = pageOffset(scroll.Y, scroll.MaxY, height, direction)
}

func pageOffset(current, maxY, height, direction float32) float32 {
	if height <= 0 || maxY <= 0 || direction == 0 {
		return current
	}
	return min(maxY, max(0, current+height*direction))
}

func (a *app) scrollTop() {
	a.editorScroll.Y, a.previewScroll.Y = 0, 0
	a.lastEditorY, a.lastPreviewY = 0, 0
	a.skipSync = 2
}
