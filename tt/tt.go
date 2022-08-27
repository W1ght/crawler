package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	reader := strings.NewReader("<html><img src=\"123\"><h1>111<h1/></html>")
	doc, _ := goquery.NewDocumentFromReader(reader)
	str, _ := doc.Find("img").Attr("src")
	fmt.Println(str)
}
