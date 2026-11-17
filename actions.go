package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui/markdown"
)

type editAction struct {
	prefix, suffix, placeholder string
	linePrefix                  string
	jump                        bool
	offset                      int
}

func (a *app) markEdited(text string) {
	if text == a.text {
		return
	}
	a.text = text
	a.recordEdit()
}

func (a *app) recordEdit() {
	if len(a.history) > 0 && a.history[a.historyAt] == a.text {
		return
	}
	a.revision, a.status = a.revision+1, "Modified"
	a.outlineValid = false
	if a.historyAt+1 < len(a.history) {
		a.history = a.history[:a.historyAt+1]
	}
	a.history = append(a.history, a.text)
	if len(a.history) > 100 {
		a.history = a.history[1:]
	}
	a.historyAt = len(a.history) - 1
}

func (a *app) undo() {
	if a.historyAt <= 0 {
		return
	}
	a.historyAt--
	a.text, a.revision, a.status = a.history[a.historyAt], a.revision+1, "Modified"
	a.outlineValid = false
}

func (a *app) redo() {
	if a.historyAt+1 >= len(a.history) {
		return
	}
	a.historyAt++
	a.text, a.revision, a.status = a.history[a.historyAt], a.revision+1, "Modified"
	a.outlineValid = false
}

func (a *app) confirmThen(next func()) {
	if a.busy {
		return
	}
	if a.text == a.savedText {
		next()
		return
	}
	if a.win == nil {
		return
	}
	a.busy = true
	win, docID, revision := a.win, a.docID, a.revision
	go func() {
		result, err := mygo.Dialog.Message(mygo.MessageOptions{
			Parent: win, Type: mygo.MessageQuestion, Title: "Unsaved changes",
			Message: "Save changes to the current document before continuing?",
			Buttons: []string{"Save", "Discard", "Cancel"}, DefaultButton: 0, CancelButton: 2,
		})
		win.Update(func() {
			if err != nil || result.Button == 2 || result.Button < 0 {
				a.busy = false
				return
			}
			if a.docID != docID || a.revision != revision {
				a.busy = false
				a.confirmThen(next)
				return
			}
			if result.Button == 1 {
				a.busy = false
				next()
				return
			}
			a.saveCurrent(func(ok bool) {
				if ok {
					next()
				}
			})
		})
	}()
}

func (a *app) newFile() {
	a.confirmThen(func() {
		a.setDocument("", defaultMarkdown)
		a.status = "New document"
	})
}

func (a *app) openFile() {
	if a.busy || a.win == nil {
		return
	}
	a.confirmThen(func() { a.openFileDialog() })
}

func (a *app) openFileDialog() {
	if a.busy || a.win == nil {
		return
	}
	a.busy = true
	win, docID, revision := a.win, a.docID, a.revision
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Parent: win, Title: "Open Markdown",
			Filters: []mygo.FileFilter{{Name: "Markdown", Extensions: []string{"md", "markdown"}}},
		})
		if err != nil || len(paths) == 0 {
			win.Update(func() {
				a.busy = false
				if err != nil {
					a.status = err.Error()
				}
			})
			return
		}
		data, err := os.ReadFile(paths[0])
		win.Update(func() {
			a.busy = false
			if a.docID != docID {
				return
			}
			if err != nil {
				a.status = err.Error()
				return
			}
			load := func() {
				a.setDocument(absolutePath(paths[0]), normalizeText(string(data)))
				a.addRecent(a.path)
				a.status = "Opened " + a.displayName()
			}
			if a.revision != revision {
				a.confirmThen(load)
			} else {
				load()
			}
		})
	}()
}

func (a *app) saveFile(saveAs bool) {
	if a.busy || a.win == nil {
		return
	}
	a.saveCurrent(func(bool) {}, saveAs)
}

func (a *app) saveCurrent(done func(bool), saveAs ...bool) {
	if a.win == nil {
		done(false)
		return
	}
	as := len(saveAs) > 0 && saveAs[0]
	path := a.path
	if as || path == "" {
		path = ""
	}
	text, rev, docID, win, name := a.text, a.revision, a.docID, a.win, a.displayName()
	a.busy, a.status = true, "Saving..."
	win.Invalidate()
	go func() {
		var err error
		if path == "" {
			path, err = mygo.Dialog.Save(mygo.SaveDialogOptions{
				Parent: win, Title: "Save Markdown", DefaultPath: name,
				Filters: []mygo.FileFilter{{Name: "Markdown", Extensions: []string{"md", "markdown"}}},
			})
			if err == nil && path != "" && filepath.Ext(path) == "" {
				path += ".md"
			}
		}
		if err == nil && path != "" {
			path = absolutePath(path)
			err = os.WriteFile(path, []byte(strings.ReplaceAll(text, "\n", "\r\n")), 0o644)
		}
		win.Update(func() {
			a.busy = false
			if a.docID != docID {
				done(false)
				return
			}
			if err != nil {
				a.status = err.Error()
				done(false)
				return
			}
			if path == "" {
				a.status = "Ready"
				done(false)
				return
			}
			a.path = path
			a.savedText = text
			a.md.LoadImage = markdown.FileImages(filepath.Dir(path))
			if a.docID == docID && a.revision == rev {
				a.status = "Saved"
			} else {
				a.status = "Saved; newer edits remain"
			}
			a.addRecent(path)
			_ = a.persist()
			done(a.docID == docID && a.revision == rev)
		})
	}()
}

func (a *app) onClose(e *mygo.CloseEvent) {
	if a.busy && !a.closing {
		e.PreventDefault()
		return
	}
	if a.closing || a.text == a.savedText {
		_ = a.persist()
		return
	}
	e.PreventDefault()
	if a.busy || a.win == nil {
		return
	}
	a.busy = true
	win := a.win
	go func() {
		result, err := mygo.Dialog.Message(mygo.MessageOptions{
			Parent: win, Type: mygo.MessageQuestion, Title: "Unsaved changes",
			Message: "Save changes before exiting?", Buttons: []string{"Save", "Discard", "Cancel"}, DefaultButton: 0, CancelButton: 2,
		})
		win.Update(func() {
			if err != nil || result.Button == 2 || result.Button < 0 {
				a.busy = false
				return
			}
			if result.Button == 1 {
				a.busy = false
				a.closing = true
				win.Close()
				return
			}
			a.saveCurrent(func(ok bool) {
				if ok {
					a.closing = true
					win.Close()
				}
			})
		})
	}()
}

func (a *app) toggleTask(line int, checked bool) {
	if line < 1 {
		return
	}
	lines := strings.SplitAfter(a.text, "\n")
	if line > len(lines) {
		return
	}
	old := lines[line-1]
	if checked {
		old = strings.Replace(old, "[ ]", "[x]", 1)
		old = strings.Replace(old, "[X]", "[x]", 1)
	} else {
		old = strings.Replace(old, "[x]", "[ ]", 1)
		old = strings.Replace(old, "[X]", "[ ]", 1)
	}
	lines[line-1] = old
	a.markEdited(strings.Join(lines, ""))
}
