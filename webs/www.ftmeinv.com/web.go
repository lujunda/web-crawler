package ftmeinv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"web-crawler/lib/http"
)

type Web struct {
}

//对应具体网站(www.ishsh.com)的爬虫逻辑
func (w Web) Analysis(url string, doc *goquery.Document) (ret []string) {

	doc.Find("#list a").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("href")
		if ok {
			src = "http://www.ftmeinv.com/" + src
			fmt.Println(src)
			ret = append(ret, src)
		}
	})
	doc.Find(".page img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok {
			path := strings.Split(url, "/")
			filename := path[len(path)-1]
			filename = strings.ReplaceAll(filename, ".html", ".jpg")
			http.GetPic(src, "/Users/bytedance/data/file/ftmeinv/"+filename)
			fmt.Println("/Users/bytedance/data/file/ftmeinv/" + filename)
		}
	})
	doc.Find(".pagelist a").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("href")
		if ok {
			page_name := strings.Split(url, "/")
			url = strings.ReplaceAll(url, page_name[len(page_name)-1], src)
			ret = append(ret, url)
		}
	})
	return
}

//起始地址
func (w Web) Root() (ret []string) {
	for i := 1; i <= 20; i++ {
		ret = append(ret, "http://www.ftmeinv.com/fei/tungirlsp/aiss/"+strconv.Itoa(i)+".html")
	}
	return
}
