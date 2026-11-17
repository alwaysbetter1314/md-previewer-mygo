package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

func appTheme(dark bool) *ui.Theme {
	t := ui.LightTheme()
	if dark {
		t = ui.DarkTheme()
		t.Background, t.Surface = ui.Hex("#1e1e1e"), ui.Hex("#2a2a2a")
		t.Border, t.Text = ui.Hex("#444444"), ui.Hex("#e5e5e5")
	} else {
		t.Background, t.Surface = ui.Hex("#ffffff"), ui.Hex("#f5f5f5")
		t.Border, t.Text = ui.Hex("#e5e5e5"), ui.Hex("#1f1f1f")
	}
	t.Radius, t.FontSize = 5, 14
	return t
}

func (a *app) view(c *ui.Context) {
	_, height := c.Size()
	a.editorHeight = max(0, height-158)
	a.previewHeight = max(0, height-123)
	if bounds := a.editorScroll.Bounds(c); bounds.H > 0 {
		a.editorHeight = bounds.H
	}
	if bounds := a.previewScroll.Bounds(c); bounds.H > 0 {
		a.previewHeight = bounds.H
	}
	if a.dark {
		c.SetTheme(a.themeDark)
	} else {
		c.SetTheme(a.themeLight)
	}
	a.syncScroll()
	a.shortcuts(c)
	ui.Column(c).Fill().Gap(0).Children(func() {
		a.menuBar(c)
		separator(c)
		ui.Row(c).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
			if !a.sidebarHidden {
				a.sidebar(c)
				verticalSeparator(c)
			}
			switch {
			case a.previewOnly:
				a.previewPane(c)
			case a.previewHidden:
				a.editorPane(c)
			default:
				ui.SplitPercent(c, &a.split, func() { a.editorPane(c) }, func() { a.previewPane(c) }).Grow(1).MinWidth(0).MinHeight(0)
			}
		})
		separator(c)
		a.statusBar(c)
	})
	if a.helpOpen {
		ui.AlertDialog(c, &a.helpOpen, "快捷键", "Ctrl+N  新建\nCtrl+O  打开\nCtrl+S  保存\nCtrl+Shift+S  另存为\nCtrl+P  纯预览\nCtrl+Z / Y  撤销 / 重做\nF11  全屏 / 还原\nF1  快捷键\nEsc  退出纯预览", "关闭")
	}
	if a.aboutOpen {
		ui.AlertDialog(c, &a.aboutOpen, "关于", "md-previewer\nbuilt with MyGo\n支持 GitHub Flavored Markdown 预览。", "关闭")
	}
}

func separator(c *ui.Context) { ui.Box(c).Height(1).FillWidth().Shrink(0).Background(c.Theme().Border) }
func verticalSeparator(c *ui.Context) {
	ui.Box(c).Width(1).FillHeight().Shrink(0).Background(c.Theme().Border)
}

func tool(c *ui.Context, name, title string, fn func()) {
	b := ui.ButtonBase(c).Size(30, 30).Shrink(0).Radius(5).Label(title).Tooltip(title).KeepFocus().OnClick(fn)
	if b.Hovered() {
		b.Background(c.Theme().SurfaceHover)
	}
	b.Children(func() { ui.Icon(c, icons[name]).Size(17, 17).TextColor(c.Theme().Text) })
}

func (a *app) previewPane(c *ui.Context) {
	ui.Column(c).Grow(1).MinWidth(0).MinHeight(0).Children(func() {
		ui.Row(c).Height(50).Shrink(0).Padding(0, 16).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "预览").FontWeight(600)
			ui.Spacer(c)
			if a.previewOnly {
				tool(c, "split", "退出纯预览", a.togglePreviewOnly)
			}
			svc := c.Services()
			tool(c, "copy", "复制全部", func() { svc.WriteClipboard(a.text); a.status = "已复制" })
		})
		separator(c)
		ui.Scroll(c).Grow(1).MinHeight(0).TrackScroll(&a.previewScroll).Children(func() {
			ui.Column(c).FillWidth().Padding(28, 32).Children(func() {
				markdown.View(c, a.previewText(), &a.md)
			})
		})
	})
}

func (a *app) statusBar(c *ui.Context) {
	ui.Row(c).Height(32).Shrink(0).Padding(0, 16).Gap(14).AlignItems(ui.Center).Children(func() {
		ui.Text(c, lineCount(a.text)).FontSize(12).TextColor(c.Theme().TextMuted)
		mode := "编辑预览"
		if a.previewOnly {
			mode = "纯预览"
		}
		ui.Text(c, mode).FontSize(12).TextColor(c.Theme().TextMuted)
		ui.Spacer(c)
		ui.Text(c, a.statusLabel()).FontSize(12).TextColor(c.Theme().TextMuted)
	})
}
