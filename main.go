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
	"strings"
)

func ExampleScrape() {
	// 首先获取总页数
	// 遍历每页的总文章
	// 输入到 excel

	// Request the HTML page.
	res, err := http.Get("https://www.aquanliang.com/blog/page/1")
	if err != nil {
		log.Fatal(err)
	}
	//b, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("%s", b)
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

	// Find the review items
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
	})
}

func trimImg(img string) string {
	img = strings.TrimLeft(img, "<img src=\"")
	img = strings.TrimRight(img, "\" decoding=\"async\" style=\"position:absolute;top:0;left:0;bottom:0;right:0;box-sizing:border-box;padding:0;border:none;margin:auto;display:block;width:0;height:0;min-width:100%;max-width:100%;min-height:100%;max-height:100%;object-fit:cover\"/>")
	return img
}

func importExcel() {

	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Insert a picture.
	if err := f.AddPicture("Sheet1", "A2", "image.png", ""); err != nil {
		fmt.Println(err)
	}
	// Insert a picture to worksheet with scaling.
	if err := f.AddPicture("Sheet1", "D2", "image.jpg",
		`{"x_scale": 0.5, "y_scale": 0.5}`); err != nil {
		fmt.Println(err)
	}
	// Insert a picture offset in the cell with printing support.
	if err := f.AddPicture("Sheet1", "H2", "image.gif", `{
        "x_offset": 15,
        "y_offset": 10,
        "print_obj": true,
        "lock_aspect_ratio": false,
        "locked": false
    }`); err != nil {
		fmt.Println(err)
	}
	// Save the spreadsheet with the origin path.
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	ExampleScrape()
}
