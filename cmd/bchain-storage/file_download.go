package main

import (
	"net/http"
	"path/filepath"
)

func init() {
	RegisterHandle("/file/download", downloadHandler)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authFile(r, false)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}
	rootPath := _rootPathFlag
	to := filepath.Join(rootPath, fAuth.space, r.FormValue("file"))
	http.ServeFile(w, r, to)
	return nil
}
