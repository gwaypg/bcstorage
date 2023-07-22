package main

import (
	"net/http"
	"syscall"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/capacity", capacityHandler)
}

func capacityHandler(w http.ResponseWriter, r *http.Request) error {
	_, fAuth, err := authWrite(r)
	if err != nil {
		return writeMsg(w, 401, errors.As(err).Code())
	}

	// implement the df -h
	root, err := _userMap.SpacePath(fAuth.spaceName)
	if err != nil {
		return writeMsg(w, 404, "space not found")
	}
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(root, &fs); err != nil {
		return errors.As(err, root)
	}

	return writeJson(w, 200, fs)
}
