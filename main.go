package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var TotalPage int
var offset = 2

func Crawler(page int, f *excelize.File) {
	res, err := http.Get("https://www.aquanliang.com/blog/page/" + strconv.Itoa(page))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 总页数
	if page == 1 {
		a := doc.Find("._1rGJJd-K0-f7qJoR9CzyeL ._1sC8pER1GUhouLkB66Mb0I").Nodes[6].Attr[1].Val
		a = strings.TrimLeft(a, "/blog/page/")
		TotalPage, _ = strconv.Atoi(a)
	}
	Crawl(doc, f)
}

func Crawl(doc *goquery.Document, f *excelize.File) {

	// 查找
	doc.Find("._1ySUUwWwmubujD8B44ZDzy span ._3gcd_TVhABEQqCcXHsrIpT").Each(func(i int, s *goquery.Selection) {
		// 图片
		img := s.Find("a").Find("._1wTUfLBA77F7m-CM6YysS6").Find("._2ahG-zumH-g0nsl6xhsF0s").
			Find("noscript").Nodes[0].FirstChild.Data
		img = trimImg(img)
		fmt.Println(img)

		// 标题
		s = s.Find("._3HG1uUQ3C2HBEsGwDWY-zw")
		title := s.Find("._3_JaaUmGUCjKZIdiLhqtfr").Text()
		fmt.Println(title)

		// 日期
		date := s.Find("._3TzAhzBA-XQQruZs-bwWjE").Nodes[0].LastChild.Data
		fmt.Println(date)

		// 访问量
		view := s.Find("._2gvAnxa4Xc7IT14d5w8MI1").Nodes[0].LastChild.Data
		fmt.Println(view)

		importExcel(offset, title, date, view, img, f)
		offset++
	})

}

func trimImg(img string) string {
	img = strings.TrimLeft(img, "<img src=\"")
	img = strings.TrimRight(img, "\" decoding=\"async\" style=\"position:absolute;top:0;left:0;bottom:0;right:0;box-sizing:border-box;padding:0;border:none;margin:auto;display:block;width:0;height:0;min-width:100%;max-width:100%;min-height:100%;max-height:100%;object-fit:cover\"/>")
	return img
}

func importExcel(i int, title string, date string, view string, img string, f *excelize.File) {
	index := strconv.Itoa(i)
	a := "A" + index
	b := "B" + index
	c := "C" + index
	d := "D" + index
	f.SetCellValue("Sheet1", a, title)
	f.SetCellValue("Sheet1", b, date)
	f.SetCellValue("Sheet1", c, view)
	f.SetCellValue("Sheet1", d, img)
}

func main() {
	// 首先获取总页数
	// 遍历每页的总文章
	// 输入到 excel
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "标题")
	f.SetCellValue("Sheet1", "B1", "日期")
	f.SetCellValue("Sheet1", "C1", "访问量")
	f.SetCellValue("Sheet1", "D1", "图片")
	Crawler(1, f)
	for i := 2; i <= TotalPage; i++ {
		Crawler(i, f)
	}
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
