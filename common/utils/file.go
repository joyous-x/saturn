package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joyous-x/saturn/common/xlog"
)

// FileMd5 make md5 for the file
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
	if !IsPathExist(dst) {
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

// IsPathExist check whether the path is exist or not
func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// MakeParentDir 生成指定路径的上级目录
//               如: /foo/bar/baz.js   会确保 /foo/bar 存在
//                   /foo/bar/baz      会确保 /foo/bar 存在
//                   /foo/bar/baz/     会确保 /foo/bar/baz 存在
func MakeParentDir(path string) (bool, error) {
	dir := filepath.Dir(path)
	if ok := IsPathExist(dir); ok {
		return true, nil
	}

	err := os.MkdirAll(dir, 0755)
	if nil != err {
		return false, err
	}
	return true, nil
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

func PathRelative2Bin(relative string) (string, error) {
	binpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	rst := filepath.Join(binpath, relative)
	return rst, nil
}

func GetExecDirPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return filepath.Dir(absPath), nil
}
