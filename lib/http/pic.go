package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//从地址src下载图片并存储为dest
func GetPic(src string, dest string) {

	ret, err := http.Get(src)
	if err != nil {
		return
	}
	defer ret.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer file.Close()

	pic, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return
	}

	io.Copy(file, bytes.NewReader(pic))
}
