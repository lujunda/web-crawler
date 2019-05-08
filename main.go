package main

import (
	chttp "./lib/http"
	"./lib/safemap"
	"./lib/safequeue"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "io/ioutil"
	"net/http"
	//"regexp"
	"strings"
	"strconv"
	"time"
)

//待处理队列
var queue = safequeue.Create()

//管道限制协程并发数量
var running = make(chan int, 40)

//避免重复抓取
var visited = safemap.Create()

//对应具体网站(www.ishsh.com)的爬虫逻辑,待独立封装todo...
func analysis(doc *goquery.Document) (ret []string) {
	
	doc.Find(".post-thumbnail").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if ok {
			ret = append(ret, href)
		}
	})
	doc.Find(".page-numbers").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			ret = append(ret, href)
		}
	})
	doc.Find(".image_cx_cont img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok {
			title := doc.Find("h1").Text()
			title = strings.Replace(title, "(", "_", -1)
			title = strings.Replace(title, ")", "_", -1)
			title = strings.Replace(title, "/", "_", -1)
			title = strings.Replace(title, " ", "_", -1)
			chttp.GetPic(src, "/data/file/" + title + ".jpg")
			fmt.Println("/data/file/" + title + ".jpg")
		}
	})
	return
}

func main() {

	//起始路径,待独立封装todo...
	for i:=1; i <= 30; i++ {
		queue.Push("/mingzhan/page/" + strconv.Itoa(i))
	}

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
				defer resp.Body.Close()

				doc, _ := goquery.NewDocumentFromReader(resp.Body)

				//解析
				nexts := analysis(doc)
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
