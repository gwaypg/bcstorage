package main

import (
	"fmt"
	"testing"
)

func TestLevelDB(t *testing.T) {
	if err := InitLevelDB("./leveldb"); err != nil {
		t.Fatal(err)
	}

	IterLevelDB("_user", func(key, val []byte) error {
		fmt.Println("_user", string(key), string(val))
		return nil
	})
	IterLevelDB("_space", func(key, val []byte) error {
		fmt.Println("_space", string(key), string(val))
		return nil
	})
	IterLevelDB(_leveldb_prefix_del, func(key, val []byte) error {
		fmt.Println("_del", string(key), string(val))
		return nil
	})
	IterLevelDB("", func(key, val []byte) error {
		fmt.Println("all", string(key), string(val))
		return nil
	})
}

func TestScanLevelDB(t *testing.T) {
	if err := InitLevelDB("./leveldb"); err != nil {
		t.Fatal(err)
	}
	name := "g1"
	space := UserSpace{}
	if err := ScanLevelDB(fmt.Sprintf(_leveldb_prefix_space, name), &space); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", space)
}
