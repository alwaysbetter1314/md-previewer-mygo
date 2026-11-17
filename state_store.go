package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type storedRecent struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type storedDocument struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
}

type storedState struct {
	Current storedDocument `json:"current"`
	Recent  []storedRecent `json:"recent,omitempty"`
	Theme   string         `json:"theme,omitempty"`
}

type stateStore struct{ path string }

func openStateStore() (*stateStore, storedState) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, storedState{}
	}
	path := filepath.Join(dir, "github.com", "goedui", "goed-md-previewer", "state.json")
	s := &stateStore{path: path}
	var data storedState
	raw, err := os.ReadFile(path)
	if err == nil && len(raw) > 0 {
		_ = json.Unmarshal(raw, &data)
	}
	return s, data
}

func (a *app) restore() {
	s, data := openStateStore()
	a.store, a.data = s, data
	a.dark = data.Theme == "dark"
	if data.Current.Path != "" {
		if err := a.loadPath(data.Current.Path); err == nil {
			a.addRecent(a.path)
			return
		}
	}
	if data.Current.Text != "" {
		a.setDocument("", data.Current.Text)
	}
}

func (a *app) persist() error {
	if a.store == nil {
		return nil
	}
	a.data.Current = storedDocument{Path: a.path, Name: a.displayName()}
	if a.path == "" {
		a.data.Current.Text = a.savedText
	}
	raw, err := json.MarshalIndent(a.data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(a.store.path), 0o755); err != nil {
		return err
	}
	tmp := a.store.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, a.store.path)
}

func (a *app) addRecent(path string) {
	if path == "" {
		return
	}
	abs := absolutePath(path)
	key := pathKey(abs)
	out := make([]storedRecent, 0, 10)
	out = append(out, storedRecent{Path: abs, Name: filepath.Base(abs)})
	seen := map[string]bool{key: true}
	for _, item := range a.data.Recent {
		itemKey := pathKey(item.Path)
		if item.Path == "" || seen[itemKey] {
			continue
		}
		seen[itemKey] = true
		out = append(out, item)
		if len(out) == 10 {
			break
		}
	}
	a.data.Recent = out
	_ = a.persist()
}

func pathKey(path string) string {
	abs := filepath.Clean(absolutePath(path))
	if runtime.GOOS == "windows" && !strings.HasPrefix(abs, `\\?\`) {
		abs = strings.ToLower(abs)
	}
	return filepath.ToSlash(abs)
}
