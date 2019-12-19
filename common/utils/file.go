package utils

import (
	"crypto/md5"
	"github.com/joyous-x/enceladus/common/xlog"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FileMd5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		xlog.Error("FileMd5 error: %s, %v", path, err)
		return ""
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		xlog.Error("FileMd5 error: %s, %v", path, err)
		return ""
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

func CopyFile(src string, dst string) int {
	ret := -1
	for {
		fs, err := os.Open(src)
		if err != nil {
			xlog.Error("CopyFile error: %v", err)
			break
		}
		defer fs.Close()
		MakeParentDir(dst)
		fd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			xlog.Error("CopyFile error: %v", err)
			break
		}
		defer fd.Close()
		_, err = io.Copy(fd, fs)
		if err != nil {
			xlog.Error("CopyFile error: %v", err)
			break
		}
		ret = 0
		break
	}
	return ret
}

func MoveFile(src string, dst string) int {
	MakeParentDir(dst)
	err := os.Rename(src, dst)
	if err != nil {
		xlog.Error("MoveFile error: %v", err)
		return -1
	}
	if !IsFileExist(dst) {
		return -1
	}
	return 0
}

func RemoveFile(path string) int {
	err := os.Remove(path)
	if err != nil {
		xlog.Error("RemoveFile error: %v", err)
		return -1
	}
	return 0
}

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeParentDir(file string) bool {
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, 0644)
	if nil != err {
		return false
	}
	return true
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}
