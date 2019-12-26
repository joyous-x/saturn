package main

import (
	"testing"
	"runtime"
    "path/filepath"
	"github.com/joyous-x/saturn/common/xlog"
)

func TestConfigParser(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	curDir := filepath.Dir(filename)
	xlog.Debug("-------test start ")
	xlog.Debug("-------test end : %v ", curDir)
}
