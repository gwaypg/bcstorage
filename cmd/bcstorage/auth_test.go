package main

import (
	"fmt"
	"testing"
)

func TestUserMap(t *testing.T) {
	if err := InitLevelDB("./leveldb"); err != nil {
		panic(err)
	}
	defer CloseLevelDB()

	userSpace, err := _userMap.GetSpace("g1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", userSpace)
}
