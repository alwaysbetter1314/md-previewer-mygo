package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

func TestNewAppLoadsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.md")
	if err := os.WriteFile(path, []byte("# Loaded document\r\n\r\nBody"), 0o600); err != nil {
		t.Fatal(err)
	}

	a := newApp(path)
	if a.path != absolutePath(path) {
		t.Fatalf("path = %q, want %q", a.path, absolutePath(path))
	}
	if a.text != "# Loaded document\n\nBody" || a.text != a.savedText {
		t.Fatalf("loaded text = %q, saved text = %q", a.text, a.savedText)
	}
}

func TestViewEditsMarkdownAndUpdatesPreview(t *testing.T) {
	a := newApp("")
	tt := ui.NewTester(a.view, 1000, 700)
	for _, label := range []string{"Markdown source", "Preview", "GitHub Flavored Markdown"} {
		if !tt.HasText(label) {
			t.Fatalf("view does not show %q: %q", label, tt.Texts())
		}
	}

	field, ok := tt.AccessByRole(ui.RoleTextField, "Markdown source")
	if !ok {
		t.Fatal("Markdown source text field is missing")
	}
	tt.ClickAt(field.Bounds.X+5, field.Bounds.Y+field.Bounds.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("# Edited document")

	if !strings.Contains(a.text, "# Edited document") {
		t.Fatalf("editor text was not updated: %q", a.text)
	}
	if a.status != "Modified" || !tt.HasText("Modified") || !tt.HasText("Edited document") {
		t.Fatalf("edit state = %q, visible text = %q", a.status, tt.Texts())
	}
}

func TestTaskToggleUndoRedoAndOutline(t *testing.T) {
	a := newApp("")
	a.setDocument("", "# Guide\n\n## Tasks\n- [ ] Ship it\n")
	got := a.headings()
	if len(got) != 2 || got[0].text != "Guide" || got[1].text != "Tasks" || got[1].level != 2 {
		t.Fatalf("headings = %#v", got)
	}

	a.toggleTask(4, true)
	if !strings.Contains(a.text, "- [x] Ship it") || a.status != "Modified" {
		t.Fatalf("task edit state = %q, status = %q", a.text, a.status)
	}
	a.undo()
	if !strings.Contains(a.text, "- [ ] Ship it") {
		t.Fatalf("undo did not restore unchecked task: %q", a.text)
	}
	if a.statusLabel() != "Saved" {
		t.Fatalf("undo to the saved document should show Saved, got %q", a.statusLabel())
	}
	a.redo()
	if !strings.Contains(a.text, "- [x] Ship it") {
		t.Fatalf("redo did not restore checked task: %q", a.text)
	}
}

func TestLargeDocumentsUseBoundedPreview(t *testing.T) {
	a := newApp("")
	a.text = "# Start\n" + strings.Repeat("content\n", 80*1024)
	preview := a.previewText()
	if len(preview) > 16*1024+100 || !strings.Contains(preview, "# Start") || !strings.Contains(preview, "first 16 KB") {
		t.Fatalf("large preview is not bounded or missing its notice: len=%d", len(preview))
	}
	a.text = "# Small\n\n" + strings.Repeat("x", 20*1024)
	if got := a.previewText(); got != a.text {
		t.Fatalf("documents below the large-file threshold should preview in full: len=%d", len(got))
	}
	a.text = strings.Repeat("界", largeDocumentLimit/3+1)
	if got := a.previewText(); !utf8.ValidString(got) {
		t.Fatalf("large preview split a UTF-8 character: %q", got[len(got)-32:])
	}
}

func TestLargeDocumentOutlineMatchesVisiblePreview(t *testing.T) {
	a := newApp("")
	a.text = "# Visible\n\n" + strings.Repeat("body\n", 4000) + "# Hidden\n\n" + strings.Repeat("tail\n", 110000)
	got := a.headings()
	if len(got) != 1 || got[0].text != "Visible" {
		t.Fatalf("large document outline should match its preview prefix: %#v", got)
	}
}

func TestPageOffsetClampsToViewportAndDocument(t *testing.T) {
	if got := pageOffset(0, 1200, 400, 1); got != 400 {
		t.Fatalf("PageDown offset = %v, want 400", got)
	}
	if got := pageOffset(1000, 1200, 400, 1); got != 1200 {
		t.Fatalf("PageDown should stop at the end: %v", got)
	}
	if got := pageOffset(100, 1200, 400, -1); got != 0 {
		t.Fatalf("PageUp should stop at the start: %v", got)
	}
}

func TestRecentFilesAreUniqueAndLimited(t *testing.T) {
	a := newApp("")
	dir := t.TempDir()
	for i := range 12 {
		a.addRecent(filepath.Join(dir, "note-"+string(rune('a'+i))+".md"))
	}
	if len(a.data.Recent) != 10 {
		t.Fatalf("recent file count = %d, want 10", len(a.data.Recent))
	}
	a.addRecent(a.data.Recent[0].Path)
	if len(a.data.Recent) != 10 || a.data.Recent[0].Path == "" || a.data.Recent[0].Path != filepath.Join(dir, "note-l.md") {
		t.Fatalf("recent list did not keep its size and current item: %#v", a.data.Recent)
	}
}

func TestPreviewOnlyRestoresSidebarState(t *testing.T) {
	a := newApp("")
	a.sidebarHidden, a.sidebarFiles = true, false
	a.togglePreviewOnly()
	if !a.previewOnly || a.sidebarHidden || a.sidebarFiles {
		t.Fatalf("preview-only state = %#v", a)
	}
	a.togglePreviewOnly()
	if a.previewOnly || !a.sidebarHidden || a.sidebarFiles {
		t.Fatalf("restored view state = %#v", a)
	}
}

func TestPageDownScrollsPreviewAndSynchronizesEditor(t *testing.T) {
	a := newApp("")
	a.text = strings.Repeat("A paragraph for scrolling.\n\n", 300)
	tt := ui.NewTester(a.view, 1000, 700)
	if a.previewScroll.MaxY <= 0 || a.previewHeight <= 0 {
		t.Fatalf("preview scroll area not laid out: state=%#v height=%.1f", a.previewScroll, a.previewHeight)
	}

	tt.Key(0, ui.KeyPageDown)
	if a.previewScroll.Y <= 0 {
		t.Fatal("PageDown did not scroll the preview")
	}
	if a.editorScroll.Y <= 0 {
		t.Fatal("PageDown did not synchronize the editor")
	}
}

func TestEditorPageDownKeepsNativeCaretPagingAndSyncsPreview(t *testing.T) {
	a := newApp("")
	a.text = strings.Repeat("A paragraph for scrolling.\n\n", 300)
	tt := ui.NewTester(a.view, 1000, 700)
	if a.editorScroll.Y != 0 || a.previewScroll.Y != 0 {
		t.Fatalf("new document should start at the top: editor=%.1f preview=%.1f", a.editorScroll.Y, a.previewScroll.Y)
	}
	field, ok := tt.AccessByRole(ui.RoleTextField, "Markdown source")
	if !ok {
		t.Fatal("Markdown source text field is missing")
	}
	tt.ClickAt(field.Bounds.X+field.Bounds.W/2, field.Bounds.Y+field.Bounds.H/2)
	tt.Key(0, ui.KeyPageDown)
	if a.editorScroll.Y <= 0 {
		t.Fatal("PageDown did not move the focused editor")
	}
	if a.previewScroll.Y <= 0 {
		t.Fatal("editor PageDown did not synchronize the preview")
	}
}
