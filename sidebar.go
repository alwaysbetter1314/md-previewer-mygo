package main

import (
	"path/filepath"

	"github.com/egoist/mygo/ui"
)

func (a *app) sidebar(c *ui.Context) {
	ui.Column(c).Width(250).Shrink(0).MinHeight(0).Children(func() {
		ui.Row(c).Height(52).Shrink(0).AlignItems(ui.Stretch).Children(func() {
			a.sidebarTab(c, "文件", a.sidebarFiles, func() { a.sidebarFiles = true })
			a.sidebarTab(c, "大纲", !a.sidebarFiles, func() { a.sidebarFiles = false })
		})
		separator(c)
		ui.Scroll(c).Grow(1).MinHeight(0).Children(func() {
			ui.Column(c).FillWidth().Padding(16, 18).Gap(2).Children(func() {
				if a.sidebarFiles {
					a.fileList(c)
				} else {
					a.outlineList(c)
				}
			})
		})
	})
}

func (a *app) sidebarTab(c *ui.Context, label string, selected bool, action func()) {
	b := ui.ButtonBase(c).Grow(1).Label(label).OnClick(action)
	b.Children(func() {
		weight := 400
		if selected {
			weight = 600
		}
		ui.Text(c, label).FontWeight(weight)
	})
	if selected {
		b.BorderWidth(0, 0, 2, 0).BorderColor(c.Theme().Text)
	}
}

func (a *app) outlineList(c *ui.Context) {
	heads := a.headings()
	if len(heads) == 0 {
		ui.Text(c, "文档里还没有标题").FontSize(13).TextColor(c.Theme().TextMuted)
		return
	}
	for _, h := range heads {
		h := h
		b := ui.ButtonBase(c.Key(h.block)).FillWidth().Height(32).Shrink(0).Justify(ui.Start).Padding(0, 3+float32(h.level-1)*15).Label(h.text).OnClick(func() { a.jumpToHeading(h) })
		if b.Hovered() {
			b.Background(c.Theme().SurfaceHover)
		}
		b.Children(func() {
			t := ui.Text(c, h.text).FontSize(14).NoWrap().Ellipsis("…")
			if h.level == 1 {
				t.FontWeight(600)
			}
		})
	}
}

func (a *app) fileList(c *ui.Context) {
	ui.Text(c, "当前文件").FontSize(12).FontWeight(600).TextColor(c.Theme().TextMuted).Padding(6, 0)
	fileRow(c, "current", a.displayName(), a.scrollTop)
	separator(c)
	ui.Text(c, "打开历史").FontSize(12).FontWeight(600).TextColor(c.Theme().TextMuted).Padding(10, 0, 4)
	if len(a.data.Recent) == 0 {
		ui.Text(c, "还没有打开记录").FontSize(13).TextColor(c.Theme().TextMuted)
	}
	for _, item := range a.data.Recent {
		path, name := item.Path, item.Name
		if name == "" {
			name = filepath.Base(path)
		}
		fileRow(c, path, name, func() { a.openRecent(path) })
	}
}

func fileRow(c *ui.Context, key, name string, action func()) {
	b := ui.ButtonBase(c.Key(key)).FillWidth().Height(34).Shrink(0).Justify(ui.Start).Gap(8).Padding(0, 4).Label(name).Tooltip(name).OnClick(action)
	if b.Hovered() {
		b.Background(c.Theme().SurfaceHover)
	}
	b.Children(func() {
		ui.Icon(c, icons["file"]).Size(15, 15).TextColor(c.Theme().TextMuted)
		ui.Text(c, name).FontSize(13).NoWrap().Ellipsis("…")
	})
}

func (a *app) openRecent(path string) {
	a.confirmThen(func() {
		if err := a.loadPath(path); err != nil {
			a.status = err.Error()
			return
		}
		a.addRecent(a.path)
		a.status = "Opened " + a.displayName()
	})
}
