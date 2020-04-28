package main

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAnonymousFuncName(f func()) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func Test_AnonymousFunName(t *testing.T) {
	fna := getAnonymousFuncName(func() {
		t.Log("this is func A")
	})
	fnb := getAnonymousFuncName(func() {
		t.Log("this is func B")
	})
	assert.Equal(t, fna, fnb, "i think they should be the same")
}
