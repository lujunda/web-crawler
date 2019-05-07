package main

import (
	"./lib/safemap"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Web struct {
	Url    string
	Html   string
	Target []string
	Next   []string
}

var ch = make(chan string, 40)
var wg = sync.WaitGroup{}
var mp = safemap.Create()

func analysis(html string) []string {
	reg_pic := regexp.MustCompile(`<a class="image_cx_cont" href="(.+?)" title="(.+)" ><img src="(.+?)"  alt="(.+?)"`)
	pics := reg_pic.FindStringSubmatch(html)

	if len(pics) > 0 {

		reg_page := regexp.MustCompile(`(.*?)-第([\d]+)页`)
		page := reg_page.FindStringSubmatch(pics[4])
		dir_name := ""
		file_name := ""
		if len(page) > 0 {
			dir_name = page[1]
			file_name = pics[4]
		} else {
			dir_name = pics[4]
			file_name = pics[4] + "-第1页"
		}

		reg_fmt := regexp.MustCompile(`[\.]{1}([a-zA-Z]+?)$`)
		file_fmt := reg_fmt.FindStringSubmatch(pics[3])
		if len(file_fmt) > 0 {
			file_name = file_name + file_fmt[0]
		}

		file_name = strings.Replace(file_name, " ", "_", -1)
		dir_name = strings.Replace(dir_name, " ", "_", -1)
		fmt.Println(file_name)

		full_cmd := "mkdir -p /data/file/" + dir_name + " && wget " + pics[3] + " -O " + "/data/file/" + dir_name + "/" + file_name + " -P " + "/data/file/" + dir_name

		exec.Command("/bin/bash", "-c", full_cmd).Output()

	}

	return regexp.MustCompile(`/([\d]+)([\_]*)([\d]*).html`).FindAllString(html, -1)
}

func dfs(url string) {
	ch <- url

	if mp.Get(url) {
		<-ch
		wg.Done()
		return
	}

	mp.Set(url, true)

	resp, _ := http.Get("https://www.ishsh.com" + url)

	if resp == nil {
		fmt.Println(url)
		mp.Set(url, false)
		<-ch

		defer func(url string) { go dfs(url) }(url)

		return
	}

	fmt.Println(url)
	body, _ := ioutil.ReadAll(resp.Body)
	html := string(body)
	resp.Body.Close()

	nexts := analysis(html)

	for _, v := range nexts {
		if !mp.Get(v) {
			wg.Add(1)
			defer func(v string) { go dfs(v) }(v)
		}
	}

	<-ch
	wg.Done()
}

func main() {
	wg.Add(1)
	go dfs("/")
	wg.Wait()
}
