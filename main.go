package main

import (
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
	"time"
)

var TotalPage int
var offset = 2
var falseSet []int

func crawler(page int, f *excelize.File) {

	doc := getDoc(page)
	// 测试一下 html 是否完整，否则重试
	a := doc.Find("._1rGJJd-K0-f7qJoR9CzyeL ._1sC8pER1GUhouLkB66Mb0I").Nodes
	l := len(a)
	// 失败最多重试5次
	for i := 0; i < 5 && l == 0; i++ {
		// 等待3秒再进行请求，以防qps过高对请求进行拦截
		time.Sleep(3 * time.Second)
		doc = getDoc(page)
		a = doc.Find("._1rGJJd-K0-f7qJoR9CzyeL ._1sC8pER1GUhouLkB66Mb0I").Nodes
		l = len(a)
	}

	if l == 0 {
		log.Printf("页面 %d 获取失败", page)
		falseSet = append(falseSet, page)
		return
	}

	// 设置总页数
	if TotalPage == 0 {
		b := a[l-1].Attr[1].Val
		b = strings.TrimLeft(b, "/blog/page/")
		TotalPage, _ = strconv.Atoi(b)
	}
	parse(doc, f)
}

func getDoc(page int) *goquery.Document {
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
	return doc
}

func parse(doc *goquery.Document, f *excelize.File) {

	// 查找
	doc.Find("._1ySUUwWwmubujD8B44ZDzy span ._3gcd_TVhABEQqCcXHsrIpT").Each(func(i int, s *goquery.Selection) {
		// 图片
		img := s.Find("a").Find("._1wTUfLBA77F7m-CM6YysS6").Find("._2ahG-zumH-g0nsl6xhsF0s").
			Find("noscript").Nodes[0].FirstChild.Data
		img = trimImg(img)
		//log.Println(img)

		// 标题
		s = s.Find("._3HG1uUQ3C2HBEsGwDWY-zw")
		title := s.Find("._3_JaaUmGUCjKZIdiLhqtfr").Text()
		//log.Println(title)

		// 日期
		date := s.Find("._3TzAhzBA-XQQruZs-bwWjE").Nodes[0].LastChild.Data
		//log.Println(date)

		// 访问量
		view := s.Find("._2gvAnxa4Xc7IT14d5w8MI1").Nodes[0].LastChild.Data
		//log.Println(view)

		insertExcel(offset, title, date, view, img, f)
		offset++
	})

}

func trimImg(img string) string {
	img = strings.TrimLeft(img, "<img src=\"")
	img = strings.TrimRight(img, "\" decoding=\"async\" style=\"position:absolute;top:0;left:0;bottom:0;right:0;box-sizing:border-box;padding:0;border:none;margin:auto;display:block;width:0;height:0;min-width:100%;max-width:100%;min-height:100%;max-height:100%;object-fit:cover\"/>")
	return img
}

func insertExcel(i int, title string, date string, view string, img string, f *excelize.File) {
	index := strconv.Itoa(i)
	a := "A" + index
	b := "B" + index
	c := "C" + index
	d := "D" + index
	err := f.SetCellValue("Sheet1", a, title)
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", b, date)
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", c, view)
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", d, img)
	if err != nil {
		log.Fatal(err)
	}
}

func initExcel() *excelize.File {
	f := excelize.NewFile()
	err := f.SetCellValue("Sheet1", "A1", "标题")
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", "B1", "日期")
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", "C1", "访问量")
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("Sheet1", "D1", "图片")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func main() {
	// 首先获取总页数
	// 遍历每页的总文章
	// 输入到 excel
	f := initExcel()
	crawler(1, f)
	for i := 2; i <= TotalPage; i++ {
		if TotalPage%40 == 0 {
			time.Sleep(20 * time.Second)
		}
		crawler(i, f)
	}
	for _, v := range falseSet {
		crawler(v, f)
	}
	if err := f.SaveAs("Data.xlsx"); err != nil {
		log.Println(err)
	}
}
