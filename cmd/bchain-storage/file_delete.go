package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gwaylib/errors"
)

func init() {
	// TODO: rollback
	RegisterHandle("/file/delete", deleteHandler)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) error {
	fAuth, ok := authWrite(r)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}

	oldName, ok := validHttpFilePath(fAuth.spaceName, r.FormValue("file"))
	if !ok {
		return writeMsg(w, 404, "space not found")
	}

	bakName := uuid.New()
	bakKey := fmt.Sprintf(_leveldb_prefix_del, time.Now().Unix(), fAuth.spaceName, bakName)
	fileKey := fmt.Sprintf(".del.%s", bakName)
	newPath, _ := validHttpFilePath(fAuth.spaceName, fileKey)

	if err := PutLevelDB(bakKey, oldName); err != nil {
		return errors.As(err)
	}

	if err := os.Rename(oldName, newPath); err != nil {
		// if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return errors.As(err)
		}
	}
	log.Warnf("Delete file:%s, from:%s", oldName, r.RemoteAddr)
	return nil
}
