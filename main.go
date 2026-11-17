package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/egoist/mygo/ui/markdown/parser"
)

const defaultMarkdown = "# md-previewer\n\nA Markdown editor with a live preview.\n\n## GitHub Flavored Markdown\n\n- [x] Tables and task lists\n- [ ] **Bold**, *italic*, and ~~strikethrough~~\n\n| Editor | Preview |\n| --- | --- |\n| Markdown | GFM |\n"

type app struct {
	text, savedText, path, status string
	docID, revision               uint64
	split                         float32
	win                           *mygo.Window
	store                         *stateStore
	data                          storedState

	previewOnly, previewHidden, sidebarHidden, sidebarFiles, syncOff bool
	savedSidebarHidden, savedSidebarFiles                            bool
	editorScroll, previewScroll                                      ui.ScrollState
	editorHeight, previewHeight                                      float32
	lastEditorY, lastPreviewY                                        float32
	skipSync                                                         int
	editor                                                           ui.Handle
	edit                                                             editAction
	history                                                          []string
	historyAt                                                        int
	outline                                                          []heading
	outlineText                                                      string
	outlineValid                                                     bool
	anchorBlock                                                      int
	md                                                               markdown.Options
	themeLight, themeDark                                            *ui.Theme
	busy, closing                                                    bool
	dark, helpOpen, aboutOpen                                        bool
}

func newApp(path string) *app {
	a := &app{
		status: "Ready", split: 50, docID: 1,
		themeLight: appTheme(false), themeDark: appTheme(true),
	}
	a.setDocument("", defaultMarkdown)
	if path != "" {
		if err := a.loadPath(path); err != nil {
			a.status = err.Error()
		}
	}
	return a
}

func (a *app) loadPath(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	a.setDocument(absolutePath(path), normalizeText(string(data)))
	return nil
}

func (a *app) setDocument(path, text string) {
	a.docID++
	a.revision++
	a.path, a.text, a.savedText = path, text, text
	a.history, a.historyAt = []string{text}, 0
	a.edit = editAction{jump: true, offset: 0}
	a.outlineValid = false
	a.editorScroll.Y, a.previewScroll.Y = 0, 0
	a.lastEditorY, a.lastPreviewY = 0, 0
	a.skipSync = 0
	dir := "."
	if path != "" {
		dir = filepath.Dir(path)
	}
	a.md = markdown.Options{OnTask: a.toggleTask, LoadImage: markdown.FileImages(dir)}
	a.md.AnchorBlock = func(_ *parser.Document, name string) int {
		if name == "outline" {
			return a.anchorBlock
		}
		return -1
	}
}

func (a *app) displayName() string {
	if a.path == "" {
		return "Untitled.md"
	}
	return filepath.Base(a.path)
}

func (a *app) statusLabel() string {
	if a.text != a.savedText && a.status != "Saving..." {
		return "Modified"
	}
	if a.status == "Modified" {
		return "Saved"
	}
	return a.status
}

func normalizeText(text string) string { return strings.ReplaceAll(text, "\r\n", "\n") }

func absolutePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func lineCount(text string) string {
	return strconv.Itoa(strings.Count(text, "\n")+1) + " lines"
}

func main() {
	a := newApp("")
	a.restore()
	if len(os.Args) > 1 && strings.TrimSpace(os.Args[1]) != "" {
		if err := a.loadPath(strings.TrimSpace(os.Args[1])); err != nil {
			a.status = err.Error()
		} else {
			a.addRecent(a.path)
		}
	}
	mygo.App.WhenReady(func() {
		a.win = mygo.NewWindow(mygo.WindowOptions{
			Title: "md-previewer", Width: 1220, Height: 820,
			MinWidth: 720, MinHeight: 500, StateKey: "main", Content: ui.View(a.view),
		})
		a.win.OnClose(a.onClose)
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
