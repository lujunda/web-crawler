package filesystem

import (
	"os"
)

//判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建路径
func MakePath(path string) {
	ok, _ := PathExists(path)
	if !ok {
		os.Mkdir(path, os.ModePerm)
	}
}
