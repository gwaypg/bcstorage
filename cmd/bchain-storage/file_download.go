package main

import (
	"net/http"
)

func init() {
	RegisterHandle("/file/download", downloadHandler)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authWrite(r)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}
	to, ok := validHttpFilePath(fAuth.spaceName, r.FormValue("file"))
	if !ok {
		return writeMsg(w, 404, "file not found")
	}
	http.ServeFile(w, r, to)
	return nil
}
