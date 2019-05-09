package ishish

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	chttp "web-crawler/lib/http"
)

type Web struct {
}

//对应具体网站(www.ishsh.com)的爬虫逻辑
func (w Web) Analysis(doc *goquery.Document) (ret []string) {

	doc.Find(".post-thumbnail").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if ok {
			ret = append(ret, "https://www.ishsh.com"+href)
		}
	})
	doc.Find(".page-numbers").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			ret = append(ret, "https://www.ishsh.com"+href)
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
			chttp.GetPic(src, "/data/file/"+title+".jpg")
			fmt.Println("/data/file/" + title + ".jpg")
		}
	})
	return
}

//起始地址
func (w Web) Root() (ret []string) {
	for i := 1; i <= 30; i++ {
		ret = append(ret, "https://www.ishsh.com/mingzhan/page/"+strconv.Itoa(i))
	}
	return
}
