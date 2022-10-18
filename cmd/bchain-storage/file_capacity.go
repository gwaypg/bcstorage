package main

import (
	"net/http"
	"path/filepath"
	"syscall"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/capacity", capacityHandler)
}

func capacityHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authFile(r, false)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	// implement the df -h
	root := filepath.Join(_rootPathFlag, fAuth.space)
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(root, &fs); err != nil {
		return errors.As(err, root)
	}

	return writeJson(w, 200, fs)
}
