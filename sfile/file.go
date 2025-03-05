package sfile

import (
	"os"
)

// Exist 判断文件是否存在
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MakeDir 创建目录，如果目录已存在则不做任何操作
func MakeDir(path string) error {
	exist, err := Exist(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
