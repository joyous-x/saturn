package utils

import (
	"os"
	"testing"
)

func Test_MakeParentDir(t *testing.T) {
	datas := []string{
		"./testdata/a/foo1/foo/",
		"./testdata/a/foo2/foo",
		"./testdata/a/foo3/foo.js",
	}
	checks := []string{
		"./testdata/a/foo1/foo",
		"./testdata/a/foo2",
		"./testdata/a/foo3",
	}
	addons := []string{
		"",
		"./testdata/a/foo2/foo",
		"./testdata/a/foo3/foo.js",
	}

	for i := range datas {
		_, erra := MakeParentDir(datas[i])
		if nil != erra {
			panic(erra)
		}
		if ok, errb := Exists(checks[i]); !ok {
			panic(errb)
		}
		if len(addons[i]) > 0 {
			if ok, errc := Exists(addons[i]); ok {
				panic(errc)
			}
		}
		err := os.Remove(checks[i])
		if err != nil {
			panic(err)
		}
	}
}
