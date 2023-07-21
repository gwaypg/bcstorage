package main

import (
	"net/http"
	"os"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/move", moveHandler)
}

func moveHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authWrite(r)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	newName, ok := validHttpFilePath(fAuth.spaceName, r.FormValue("new"))
	if !ok {
		return writeMsg(w, 403, "error filepath")
	}
	oldName, ok := validHttpFilePath(fAuth.spaceName, r.FormValue("file"))
	if !ok {
		return writeMsg(w, 404, "file not found")
	}
	if err := os.Rename(oldName, newName); err != nil {
		return errors.As(err)
	}
	log.Infof("Rename file %s to %s, from %s", oldName, newName, r.RemoteAddr)
	return nil
}
