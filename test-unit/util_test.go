package main

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParser(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	curDir := filepath.Dir(filename)
	xlog.Debug("-------test start ")
	xlog.Debug("-------test end : %v ", curDir)
	var a = -1.2
	var b = -2.11
	xlog.Debug("-------test end : %v %v", int(a), int(b))
}
