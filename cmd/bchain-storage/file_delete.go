package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/delete", deleteHandler)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authFile(r, true)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	rootPath := _rootPathFlag
	path := filepath.Join(rootPath, fAuth.space, r.FormValue("file"))
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return errors.As(err)
		}
	}
	log.Warnf("Remove file:%s, from:%s", path, r.RemoteAddr)
	return nil
}
