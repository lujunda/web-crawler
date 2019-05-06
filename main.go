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
)

type Web struct {
	Url    string
	Html   string
	Target []string
	Next   []string
}

//map加锁
type UrlMap struct {
    m   map[string]bool
    w   sync.RWMutex
}

//设置时加写独占锁
func (u *UrlMap)set(key string, val bool) {
    u.w.Lock()
    defer u.w.Unlock()

    u.m[key] = val
}

//获取时加读共享锁
func (u *UrlMap)get(key string) bool {
    u.w.RLock()
    defer u.w.RUnlock()

	ret, ok := u.m[key]

    if ok {
        return ret
    }

    return false
}

var ch chan string
var wg sync.WaitGroup
var mp UrlMap

func dfs(url string) {

    if mp.get("https://www.ishsh.com"+url) {
		<-ch
		wg.Done()
		return
    }

    mp.set("https://www.ishsh.com"+url, true)

	resp, _ := http.Get("https://www.ishsh.com" + url)
	if resp == nil {
		fmt.Println(url)
		<-ch
		wg.Done()
		return
	}
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
		//fmt.Println(file_name)

		full_cmd := "mkdir -p /data/file/" + dir_name + " && wget " + pics[3] + " -O " + "/data/file/" + dir_name + "/" + file_name + " -P " + "/data/file/" + dir_name

		exec.Command("/bin/bash", "-c", full_cmd).Output()

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
		go dfs(v)
	}
	wg.Done()
}

func main() {
	ch = make(chan string, 40)
	mp = UrlMap{w: sync.RWMutex{}, m: make(map[string]bool)}

	ch <- "root"
	wg.Add(1)
	go dfs("/")

	wg.Wait()
}
