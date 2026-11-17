# md-previewer-mygo

A standalone [MyGo](https://github.com/egoist/mygo) version of md-previewer. It edits Markdown source beside a live GitHub Flavored Markdown preview and supports file recovery, recent files, and an outline sidebar.

## Features

- Open, save, save as, create documents, and confirm unsaved changes before switching or closing.
- Live GFM preview, including tables, task lists, strikethrough, autolinks, code blocks, and images.
- Clickable document outline, recent files, and restoration of the last open document and theme.
- Formatting toolbar, undo/redo, task checkbox updates, and synchronized editor/preview scrolling.
- Preview-only, focus, sidebar, and light/dark theme modes.
- PageUp/PageDown scrolling and a 16 KB preview limit for documents at least 512 KB.

## Run

```powershell
cd md-previewer-mygo
go run .            # start without a file, loading the bundled sample
go run . [file.md]  # start and open a file
```

Requirements: Go 1.27.4 or newer and the MyGo module pinned in `go.mod` (v0.9.2).

A window opens with the sidebar, Markdown source editor, and live preview. Pass a path to open a document immediately; the last document and theme are restored on the next launch.

## Shortcuts

| Shortcut | Action |
| --- | --- |
| `Ctrl+N` | New document |
| `Ctrl+O` | Open file |
| `Ctrl+S` | Save |
| `Ctrl+Shift+S` | Save as |
| `Ctrl+P` | Toggle preview-only mode |
| `Ctrl+Z` / `Ctrl+Y` | Undo / redo (editor focus) |
| `PageUp` / `PageDown` | Scroll the preview one screen, when the editor is not focused (scrolls the editor instead when the preview is hidden) |
| `F1` | Shortcut help |
| `F11` | Toggle fullscreen |
| `Esc` | Leave preview-only mode |

## Documentation

See [软件操作说明书.md](软件操作说明书.md) for the full Chinese user manual covering every menu, dialog, view mode, and configuration file.

## Development

Run its checks from this directory:

```powershell
go test ./...
go tool mygo vet .
```

Build a production app with `go tool mygo build`.

## License

Released under the [MIT License](LICENSE).
