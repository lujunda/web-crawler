package youtube

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

type Web struct {
}

//对应具体网站(www.youtube.com)的爬虫逻辑
func (w Web) Analysis(doc *goquery.Document) (ret []string) {

	doc.Find(".yt-lockup-content").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".yt-uix-tile-link").Text()
		link, _ := s.Find(".yt-uix-tile-link").Attr("href")
		Write(link, title)
	})
	return
}

//起始地址
func (w Web) Root() (ret []string) {
	for i := 1; i <= 1; i++ {
		ret = append(ret, "https://www.youtube.com/results?search_query=xbox+games&page="+strconv.Itoa(i))
	}
	return
}
