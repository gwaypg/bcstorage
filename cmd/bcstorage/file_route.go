package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gwaylib/errors"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request) error

var handles = map[string]HandleFunc{}

func RegisterHandle(path string, handle HandleFunc) {
	_, ok := handles[path]
	if ok {
		panic("already registered:" + path)
	}
	handles[path] = handle
}

func writeMsg(w http.ResponseWriter, code int, msg string) error {
	w.WriteHeader(code)
	if _, err := w.Write([]byte(msg)); err != nil {
		return errors.As(err)
	}
	return nil
}
func writeJson(w http.ResponseWriter, code int, obj interface{}) error {
	output, err := json.MarshalIndent(obj, "", "	")
	if err != nil {
		return errors.As(err)
	}

	w.WriteHeader(code)
	if _, err := w.Write(output); err != nil {
		return errors.As(err)
	}
	return nil
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Infof("from:%s,method:%s,path:%+v", r.RemoteAddr, r.Method, r.URL.Path)

	// for public read
	if strings.HasPrefix(r.URL.Path, "/public") {
		//log.Infof("from:%s,method:%s,path:%+v", r.RemoteAddr, r.Method, r.URL.Path)
		paths := strings.Split(r.URL.Path, "/")
		if len(paths) < 4 {
			writeMsg(w, 404, "error paths")
			return
		}
		//log.Info("paths:", paths, paths[2])
		userSpace, ok := _userMap.GetSpace(paths[2])
		if !ok {
			writeMsg(w, 404, fmt.Sprintf("no userspace '%s'", paths[2]))
			return
		}
		if userSpace.Private {
			writeMsg(w, 401, "unauth")
			return
		}
		to, err := validHttpFilePath(paths[2], filepath.Join(paths[3:]...))
		if err != nil {
			//log.Info(errors.As(err))
			writeMsg(w, 404, errors.As(err).Code())
			return
		}

		log.Info("server file", to, r.URL.Path)

		http.ServeFile(w, r, to)
		return
	}

	// route handler
	handle, ok := handles[r.URL.Path]
	if !ok {
		writeMsg(w, 404, "Not found")
		return
	}

	if err := handle(w, r); err != nil {
		writeMsg(w, 500, err.Error())
		return
	}
	return
}
