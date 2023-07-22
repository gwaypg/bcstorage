package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/sys/file/token", tokenHandler)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) error {
	auth, err := authAdmin(r)
	if err != nil {
		return writeMsg(w, 401, errors.As(err).Code())
	}
	if auth.User == adminUser {
		return writeMsg(w, 401, "can not use 'admin' to manage files")
	}
	space := r.FormValue("space")
	if len(space) == 0 {
		space = auth.User
	}
	file := r.FormValue("file")
	if len(file) == 0 {
		return writeMsg(w, 403, "params failed")
	}

	if r.Method == "POST" {
		if !_handler.DelayToken(file) {
			return writeMsg(w, 403, "token has expired")
		}
		return writeMsg(w, 200, "success")
	}

	if r.Method == "DELETE" {
		_handler.DeleteToken(file)
		return writeMsg(w, 200, "success")
	}

	token := uuid.New().String()
	_handler.AddToken(space, file, token)
	return writeMsg(w, 200, token)
}
