package console

import (
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) servePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := os.ReadFile(filepath.Join(s.webDir, name+".html"))
		if err != nil {
			writeError(w, http.StatusNotFound, "page not found")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	}
}
