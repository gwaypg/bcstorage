package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/truncate", truncateHandler)
}

func truncateHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authWrite(r)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	size, err := strconv.ParseInt(r.FormValue("size"), 10, 64)
	if err != nil {
		return writeMsg(w, 403, "file size failed")
	}

	path, ok := validHttpFilePath(fAuth.spaceName, r.FormValue("file"))
	if !ok {
		return writeMsg(w, 404, "file not found")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return errors.As(err, path)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644) // nolint
	if err != nil {
		return errors.As(err, path)
	}
	defer f.Close()

	log.Infof("Trucate %s, size:%d, from:%s", path, size, r.RemoteAddr)
	if err := f.Truncate(size); err != nil {
		return errors.As(err, path)
	}
	return nil
}
