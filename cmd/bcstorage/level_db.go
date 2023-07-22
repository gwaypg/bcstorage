package main

import (
	"encoding/json"

	"github.com/gwaylib/errors"
	"github.com/syndtr/goleveldb/leveldb"
	lerrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	_levelDB *leveldb.DB
)

const (
	_leveldb_prefix_user  = "_user.%s"      // username
	_leveldb_prefix_space = "_space.%s"     // spacename
	_leveldb_prefix_del   = "_del.%d.%s.%s" // timestamp.spacename.uuid
)

func CloseLevelDB() error {
	if _levelDB != nil {
		return _levelDB.Close()
	}
	return nil
}

func InitLevelDB(path string) error {
	if _levelDB != nil {
		return errors.New("level db has inited").As(path)
	}

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		if !lerrors.IsCorrupted(err) {
			return errors.As(err)
		}
		db, err = leveldb.RecoverFile(path, nil)
		if err != nil {
			return errors.As(err)
		}
	}
	_levelDB = db
	return nil
}

func PutLevelDB(key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.As(err, key, val)
	}
	if err := _levelDB.Put([]byte(key), data, nil); err != nil {
		return errors.As(err)
	}
	return nil
}
func ScanLevelDB(key string, val interface{}) error {
	data, err := _levelDB.Get([]byte(key), nil)
	if err != nil {
		return errors.As(err, key)
	}
	if err := json.Unmarshal(data, val); err != nil {
		return errors.As(err, key, string(data))
	}
	return nil
}

func DelLevelDB(key string) error {
	if err := _levelDB.Delete([]byte(key), nil); err != nil {
		return errors.As(err, key)
	}
	return nil
}

func IterLevelDB(prefix string, cb func(key, val []byte) error) error {
	iter := _levelDB.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		if err := cb(iter.Key(), iter.Value()); err != nil {
			return errors.As(err)
		}
	}
	return errors.As(iter.Error())
}
