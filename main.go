package main

import (
	"fmt"
	_ "io"
	"io/ioutil"
	"net/http"
	_ "os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Web struct {
	Url    string
	Html   string
	Target []string
	Next   []string
}

var ch chan string
var wg sync.WaitGroup
var lock sync.RWMutex

func dfs(url string, m map[string]bool) {
	lock.RLock()
	_, ok := m["https://www.ishsh.com"+url]
	lock.RUnlock()
	if ok {
		<-ch
		wg.Done()
		return
	} else {
		lock.Lock()
		m["https://www.ishsh.com"+url] = true
		lock.Unlock()
	}

	resp, _ := http.Get("https://www.ishsh.com" + url)
	body, _ := ioutil.ReadAll(resp.Body)
	html := string(body)
	resp.Body.Close()

	reg_pic := regexp.MustCompile(`<a class="image_cx_cont" href="(.+?)" title="(.+)" ><img src="(.+?)"  alt="(.+?)"`)
	pics := reg_pic.FindStringSubmatch(html)
	if len(pics) > 0 {
		//fmt.Println(pics[4])
		//fmt.Println(pics[3])
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

		exec.Command("/bin/bash", "-c", full_cmd)

		/*
		   pic, _ := http.Get(pics[3])
		   file, err := os.Create(full_name)
		   fmt.Println(err)
		   io.Copy(file, pic.Body)
		*/

	}

	<-ch

	reg := regexp.MustCompile(`/([\d]+)([\_]*)([\d]*).html`)
	nexts := reg.FindAllString(html, -1)
	for _, v := range nexts {
		ch <- v
		wg.Add(1)
		go dfs(v, m)
	}
	wg.Done()
}

func main() {
	ch = make(chan string, 40)
	m := make(map[string]bool)
	lock = sync.RWMutex{}

	ch <- "root"
	wg.Add(1)
	go dfs("/", m)

	time.Sleep(time.Second * 5)
	wg.Wait()
}
