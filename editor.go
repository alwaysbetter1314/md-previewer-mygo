package main

import (
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

func (a *app) editorPane(c *ui.Context) {
	ui.Column(c).Grow(1).MinWidth(0).MinHeight(0).Children(func() {
		ui.Row(c).Height(44).Shrink(0).Padding(0, 10).Gap(5).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Height(30).Padding(0, 10).Gap(7).AlignItems(ui.Center).Radius(5).Background(c.Theme().Surface).Children(func() {
				ui.Icon(c, icons["file"]).Size(16, 16).TextColor(c.Theme().TextMuted)
				ui.Text(c, a.displayName()).FontWeight(600).FontSize(13).NoWrap().Ellipsis("…")
				if a.text != a.savedText {
					ui.Text(c, "•").TextColor(c.Theme().TextMuted)
				}
			})
			tool(c, "plus", "新建", a.newFile)
		})
		separator(c)
		a.editorToolbar(c)
		separator(c)
		svc := c.Services()
		ed := ui.TextAreaBase(c.Key(a.docID), &a.text).
			Grow(1).MinHeight(0).MinWidth(0).Padding(18, 20).
			Font("monospace").FontSize(14).Label("Markdown source").
			Bind(&a.editor).TrackScroll(&a.editorScroll).
			OnChange(a.recordEdit).
			OnTextKey(func(k ui.TextKey) bool {
				undo := k.Mods == ui.Cmd && k.Key == ui.KeyZ
				redo := k.Mods == ui.Cmd && k.Key == ui.KeyY || k.Mods == ui.Cmd|ui.Shift && k.Key == ui.KeyZ
				if !undo && !redo {
					return false
				}
				if a.text != k.Text() {
					a.text = k.Text()
					a.recordEdit()
				}
				if undo {
					a.undo()
				} else {
					a.redo()
				}
				k.SetText(a.text)
				svc.Invalidate()
				return true
			})
		ed.ContextMenu(func(m *ui.Menu) { m.EditItems() })
		a.applyEdit(ed)
	})
}

func (a *app) editorToolbar(c *ui.Context) {
	ui.ScrollHorizontal(c).Height(40).Shrink(0).FillWidth().Children(func() {
		ui.Row(c).Height(40).Padding(0, 10).Gap(2).AlignItems(ui.Center).Children(func() {
			tool(c, "undo", "撤销", a.undo)
			tool(c, "redo", "重做", a.redo)
			toolDivider(c)
			tool(c, "bold", "粗体", func() { a.edit = editAction{prefix: "**", suffix: "**"} })
			tool(c, "italic", "斜体", func() { a.edit = editAction{prefix: "*", suffix: "*"} })
			tool(c, "code", "行内代码", func() { a.edit = editAction{prefix: "`", suffix: "`"} })
			tool(c, "link", "链接", func() { a.edit = editAction{prefix: "[", suffix: "](https://)", placeholder: "链接"} })
			tool(c, "image", "图片", func() { a.edit = editAction{prefix: "![", suffix: "](image.png)", placeholder: "描述"} })
			toolDivider(c)
			tool(c, "list", "列表", func() { a.edit = editAction{linePrefix: "- "} })
			tool(c, "tasks", "任务列表", func() { a.edit = editAction{linePrefix: "- [ ] "} })
			toolDivider(c)
			tool(c, "eye", "纯预览", a.togglePreviewOnly)
			title := "隐藏预览"
			if a.previewHidden {
				title = "显示预览"
			}
			tool(c, "split", title, func() { a.previewHidden = !a.previewHidden; a.skipSync = 2 })
		})
	})
}

func toolDivider(c *ui.Context) {
	ui.Box(c).Width(1).Height(18).Margin(0, 4).Shrink(0).Background(c.Theme().Border)
}

func (a *app) applyEdit(ed ui.Element) {
	op := a.edit
	if op == (editAction{}) {
		return
	}
	a.edit = editAction{}
	if op.jump {
		ed.SetTextSelection(op.offset, op.offset)
		return
	}
	start, end := ed.TextSelection()
	text := []rune(a.text)
	start, end = min(max(start, 0), len(text)), min(max(end, 0), len(text))
	if op.linePrefix != "" {
		lineStart, lineEnd := start, end
		for lineStart > 0 && text[lineStart-1] != '\n' {
			lineStart--
		}
		if start == end {
			for lineEnd < len(text) && text[lineEnd] != '\n' {
				lineEnd++
			}
		}
		lines := strings.Split(string(text[lineStart:lineEnd]), "\n")
		for i := range lines {
			lines[i] = op.linePrefix + lines[i]
		}
		insert := strings.Join(lines, "\n")
		ed.ReplaceText(lineStart, lineEnd, insert)
		ed.SetTextSelection(lineStart, lineStart+utf8.RuneCountInString(insert))
	} else {
		selected := string(text[start:end])
		if start == end {
			selected = op.placeholder
		}
		insert := op.prefix + selected + op.suffix
		ed.ReplaceText(start, end, insert)
		from := start + utf8.RuneCountInString(op.prefix)
		ed.SetTextSelection(from, from+utf8.RuneCountInString(selected))
	}
	a.recordEdit()
	a.editor.Focus()
}
