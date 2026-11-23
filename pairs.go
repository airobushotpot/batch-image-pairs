// Package imagepairs matches source photographs with edited exports by filename.
package imagepairs

import (
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Match contains slash-separated paths relative to each input filesystem.
// Status is paired, missing, extra or ambiguous.
type Match struct {
	Name string
	Originals []string
	Edits []string
	Status string
}

// Compare recursively indexes image filenames. A read error returns no partial report.
func Compare(originals, edits fs.FS) ([]Match, error) {
	left, err := index(originals)
	if err != nil { return nil, err }
	right, err := index(edits)
	if err != nil { return nil, err }
	keys := make([]string, 0, len(left)+len(right))
	for name := range left { keys = append(keys, name) }
	for name := range right {
		if _, exists := left[name]; !exists { keys = append(keys, name) }
	}
	sort.Strings(keys)
	report := make([]Match, 0, len(keys))
	for _, name := range keys {
		m := Match{Name: name, Originals: left[name], Edits: right[name]}
		switch {
		case len(m.Originals) > 1 || len(m.Edits) > 1: m.Status = "ambiguous"
		case len(m.Originals) == 0: m.Status = "extra"
		case len(m.Edits) == 0: m.Status = "missing"
		default: m.Status = "paired"
		}
		report = append(report, m)
	}
	return report, nil
}

func index(tree fs.FS) (map[string][]string, error) {
	items := make(map[string][]string)
	err := fs.WalkDir(tree, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil { return err }
		if !entry.Type().IsRegular() { return nil }
		ext := path.Ext(entry.Name())
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png", ".webp", ".avif", ".tif", ".tiff", ".heic":
			name := strings.ToLower(strings.TrimSuffix(entry.Name(), ext))
			items[name] = append(items[name], p)
		}
		return nil
	})
	return items, err
}
