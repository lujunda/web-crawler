package main

import (
	chttp "./lib/http"
	"./lib/safemap"
	"./lib/safequeue"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//待处理队列
var queue = safequeue.Create()

//管道限制协程并发数量
var running = make(chan int, 40)

//避免重复抓取
var visited = safemap.Create()

//对应具体网站(www.ishsh.com)的爬虫逻辑,待独立封装todo...
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

		fmt.Println(pics[3])
		chttp.GetPic(pics[3], "/data/file/"+file_name)
	}

	return regexp.MustCompile(`/([\d]+)([\_]*)([\d]*).html`).FindAllString(html, -1)
}

func main() {

	//起始路径,待独立封装todo...
	queue.Push("/")

	for true {
		if queue.Len() > 0 {

			//fmt.Println(queue.Len())
			//fmt.Println(len(running))

			path, _ := queue.Pop()
			running <- 1

			go func(url string) {
				defer func() { <-running }()

				//避免重复抓取页面
				if visited.Get(url) {
					return
				}
				visited.Set(url, true)

				//http,todo...
				resp, _ := http.Get("https://www.ishsh.com" + url)
				if resp == nil {
					fmt.Println(url)

					visited.Set(url, false)
					queue.Push(url)

					return
				}
				body, _ := ioutil.ReadAll(resp.Body)
				html := string(body)
				resp.Body.Close()

				//解析
				nexts := analysis(html)
				for _, v := range nexts {
					if !visited.Get(v) {
						queue.Push(v)
					}
				}
			}(path)

			continue
		}
		if len(running) > 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}
