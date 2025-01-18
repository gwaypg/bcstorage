package main

import (
	"os"
	"strconv"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
	"github.com/gwaylib/log/proto"
)

var log *logger.Logger

func init() {
	level := proto.LevelInfo
	logLevel, err := strconv.Atoi(os.Getenv("BCSTORAGE_LOG_LEVEL"))
	if err == nil {
		level = proto.Level(logLevel)
	}
	log = logger.New(&proto.Context{"bc-storage", "1.0.0", logger.HostName}, "server", level, stdio.New(os.Stdout, os.Stderr))
}
