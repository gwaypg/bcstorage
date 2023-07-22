package main

import (
	"net/http"

	"github.com/gwaylib/errors"
)

func init() {
	RegisterHandle("/file/download", downloadHandler)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) error {
	authPath, _, err := authWrite(r)
	if err != nil {
		return writeMsg(w, 401, errors.As(err).Code())
	}
	http.ServeFile(w, r, authPath)
	return nil
}
