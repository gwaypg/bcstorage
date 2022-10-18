package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/move", moveHandler)
}

func moveHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authFile(r, true)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	newPath := r.FormValue("new")
	if !validHttpFilePath(newPath) {
		return writeMsg(w, 403, "error filepath")
	}
	rootPath := _rootPathFlag
	file := r.FormValue("file")
	oldName := filepath.Join(rootPath, fAuth.space, file)
	newName := filepath.Join(rootPath, fAuth.space, newPath)
	if err := os.Rename(oldName, newName); err != nil {
		return errors.As(err)
	}
	log.Infof("Rename file %s to %s, from %s", oldName, newName, r.RemoteAddr)
	return nil
}
