package store

import (
	"os"
	"path/filepath"
	"strings"
)

func (s *Store) List(prefix string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dir := filepath.Join(s.root, filepath.FromSlash(prefix))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		out = append(out, strings.TrimSuffix(entry.Name(), ".json"))
	}
	return out
}

func (s *Store) Count(prefix string) int {
	return len(s.List(prefix))
}
