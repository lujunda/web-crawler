package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
	"web-crawler/lib/safemap"
	"web-crawler/lib/safequeue"
	web "web-crawler/webs/www.ishsh.com"
)

/**
 * 实现一个页面爬虫需要实现两个方法:
 * root()方法返回爬虫的起始路径,如["www.abc.com/main/page/1", "www.abc.com/main/page/2"].
 * analysis()方法执行对每个页面的具体分析逻辑,返回值为需要进一步分析的地址的切片.
 */
type Template interface {
	Root() []string
	Analysis(doc *goquery.Document) []string
}

//待处理队列
var queue = safequeue.Create()

//管道限制协程并发数量
var running = make(chan int, 40)

//避免重复抓取
var visited = safemap.Create()

func main() {

	var t Template
	t = web.Web{}

	//将起始路径压入待处理队列
	for _, v := range t.Root() {
		queue.Push(v)
	}

	for true {
		if queue.Len() > 0 {

			path, _ := queue.Pop()
			fmt.Println(path)
			running <- 1

			go func(url string) {
				defer func() { <-running }()

				//避免重复抓取页面
				if visited.Get(url) {
					return
				}
				visited.Set(url, true)

				//http
				resp, _ := http.Get(url)
				if resp == nil {
					fmt.Println(url)

					visited.Set(url, false)
					queue.Push(url)

					return
				}
				defer resp.Body.Close()

				doc, _ := goquery.NewDocumentFromReader(resp.Body)

				//解析
				nexts := t.Analysis(doc)
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
