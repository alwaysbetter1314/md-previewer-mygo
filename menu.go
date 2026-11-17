package main

import "github.com/egoist/mygo/ui"

func (a *app) shortcuts(c *ui.Context) {
	c.OnShortcut(ui.Cmd, ui.KeyN, a.newFile)
	c.OnShortcut(ui.Cmd, ui.KeyO, a.openFile)
	c.OnShortcut(ui.Cmd, ui.KeyS, func() { a.saveFile(false) })
	c.OnShortcut(ui.Cmd|ui.Shift, ui.KeyS, func() { a.saveFile(true) })
	c.OnShortcut(ui.Cmd, ui.KeyP, a.togglePreviewOnly)
	c.OnShortcut(0, ui.KeyF1, func() { a.helpOpen = true })
	c.OnShortcut(0, ui.KeyF11, func() {
		if a.win != nil {
			a.win.SetFullScreen(!a.win.IsFullScreen())
		}
	})
	c.OnShortcut(0, ui.KeyEscape, func() {
		if a.previewOnly {
			a.togglePreviewOnly()
		}
	})
	if !a.editor.Focused(c) {
		c.OnShortcut(0, ui.KeyPageUp, func() { a.pageScroll(false) })
		c.OnShortcut(0, ui.KeyPageDown, func() { a.pageScroll(true) })
	}
}

func menuButton(c *ui.Context, label string, build func(*ui.Menu)) {
	b := ui.ButtonBase(c).Height(30).Padding(0, 12).Label(label).Radius(4)
	if b.Hovered() {
		b.Background(c.Theme().SurfaceHover)
	}
	b.Children(func() { ui.Text(c, label) })
	b.PopoverMenuButton(build)
}

func (a *app) menuBar(c *ui.Context) {
	svc := c.Services()
	ui.Row(c).Height(38).Shrink(0).Padding(0, 8).Gap(2).AlignItems(ui.Center).Children(func() {
		menuButton(c, "文件", func(m *ui.Menu) {
			if m.Item("新建").Shortcut(ui.Cmd, ui.KeyN).Chosen() {
				a.newFile()
			}
			if m.Item("打开").Shortcut(ui.Cmd, ui.KeyO).Chosen() {
				a.openFile()
			}
			if m.Item("保存").Shortcut(ui.Cmd, ui.KeyS).Chosen() {
				a.saveFile(false)
			}
			if m.Item("另存为").Shortcut(ui.Cmd|ui.Shift, ui.KeyS).Chosen() {
				a.saveFile(true)
			}
			m.Separator()
			if m.Item("退出").Chosen() && a.win != nil {
				a.win.Close()
			}
		})
		menuButton(c, "编辑", func(m *ui.Menu) {
			if m.Item("撤销").Shortcut(ui.Cmd, ui.KeyZ).Chosen() {
				a.undo()
			}
			if m.Item("重做").Shortcut(ui.Cmd, ui.KeyY).Chosen() {
				a.redo()
			}
			m.Separator()
			if m.Item("复制全部").Chosen() {
				svc.WriteClipboard(a.text)
				a.status = "已复制"
			}
		})
		menuButton(c, "视图", func(m *ui.Menu) {
			if m.Item("纯预览").Shortcut(ui.Cmd, ui.KeyP).Checked(a.previewOnly).Chosen() {
				a.togglePreviewOnly()
			}
			if m.Item("专注模式").Checked(a.sidebarHidden && a.previewHidden).Chosen() {
				a.toggleFocus()
			}
			if m.Item("同步滚动").Checked(!a.syncOff).Chosen() {
				a.syncOff = !a.syncOff
			}
			if m.Item("显示侧栏").Checked(!a.sidebarHidden).Chosen() {
				a.sidebarHidden = !a.sidebarHidden
			}
			if m.Item("显示预览").Checked(!a.previewHidden).Chosen() {
				a.previewHidden = !a.previewHidden
			}
		})
		menuButton(c, "主题", func(m *ui.Menu) {
			if m.Item("明亮").Checked(!a.dark).Chosen() {
				a.setTheme(false)
			}
			if m.Item("暗黑").Checked(a.dark).Chosen() {
				a.setTheme(true)
			}
		})
		menuButton(c, "帮助", func(m *ui.Menu) {
			if m.Item("快捷键").Shortcut(0, ui.KeyF1).Chosen() {
				a.helpOpen = true
			}
			m.Separator()
			if m.Item("关于").Chosen() {
				a.aboutOpen = true
			}
		})
	})
}

func (a *app) setTheme(dark bool) {
	a.dark = dark
	a.data.Theme = "light"
	if dark {
		a.data.Theme = "dark"
	}
	if err := a.persist(); err != nil {
		a.status = err.Error()
	}
}

func (a *app) togglePreviewOnly() {
	a.previewOnly = !a.previewOnly
	if a.previewOnly {
		a.savedSidebarHidden, a.savedSidebarFiles = a.sidebarHidden, a.sidebarFiles
		a.sidebarHidden, a.sidebarFiles = false, false
		a.status = "纯预览"
	} else {
		a.sidebarHidden, a.sidebarFiles = a.savedSidebarHidden, a.savedSidebarFiles
		a.status = "编辑预览"
	}
	a.skipSync = 2
}

func (a *app) toggleFocus() {
	a.sidebarHidden = !(a.sidebarHidden && a.previewHidden)
	a.previewHidden = a.sidebarHidden
	if a.sidebarHidden {
		a.status = "专注模式"
	} else {
		a.status = "编辑预览"
	}
	a.skipSync = 2
}
