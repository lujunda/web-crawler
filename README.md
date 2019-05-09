# 爬虫
利用go协程（默认40个协程）实现的web爬虫。

# 使用方法
1. 自定义定义一个包并实现接口`Template`(`main.go`)，可以参考`webs/www.ishish.com/web.go`。
```
/**
 * 实现一个页面爬虫需要实现两个方法:
 * root()方法返回爬虫的起始路径,如["www.abc.com/main/page/1", "www.abc.com/main/page/2"].
 * analysis()方法执行对每个页面的具体分析逻辑,返回值为需要进一步分析的地址的切片.
 */
type Template interface {
    Root() []string
    Analysis(doc *goquery.Document) []string
}
```
2. 在`main.go`的`import`列表中引入，并设置别名为`web`，例如：
```
import (
    web "web-crawler/webs/www.ishsh.com"
    ...
    ...
)
```
3. `go run main.go`
