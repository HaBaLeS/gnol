package engine

import (
	"bytes"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var baseUrl = ""

func Leech(site string) string {


	//baseUrl = fmt.Sprintf("https://%s/g/%s/", site, comicid)
	baseUrl = fmt.Sprintf("https://%s", site)

	otitle := extractName("dnf-comic-name")
	ctitle := cleanTitle(otitle)

	pages := extractPageCount()

	thefolder := path.Join("leech-data", ctitle)
	err := os.Mkdir(thefolder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= pages; i++ {

		link := extractImageUrl(i)
		fmt.Printf("DL: %s\n", link)
		ext := path.Ext(link)

		var webpName string
		if ext == ".webp" {
			webpName = path.Join("leech-data", ctitle, fmt.Sprintf("%06d%s", i, ext))
		}
		pngName := path.Join("leech-data", ctitle, fmt.Sprintf("%06d%s", i, ".png"))
		if webpName != "" {
			fh, err := os.Create(webpName)
			defer fh.Close()
			if err != nil {
				panic(err)
			}
			imgData := getUrlData(link)
			io.Copy(fh, bytes.NewBuffer(imgData)) //fixme you so not need to persist to disk!!
			util.Webp2Png(webpName, pngName)
			os.Remove(webpName)
		}

	}
	fmt.Printf("Comic:\n%s Pages: %d\n", otitle, pages)
	return thefolder
}

func getUrlData(target string) []byte {
	res, err := http.Get(target)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		panic(fmt.Errorf("wrong status: %s", res.Status))
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func cleanTitle(title string) string {
	ct := strings.ReplaceAll(title, " ", "_")
	ct = strings.ReplaceAll(ct, "(", "")
	ct = strings.ReplaceAll(ct, ")", "")
	ct = strings.ReplaceAll(ct, "[", "")
	ct = strings.ReplaceAll(ct, "]", "")
	ct = strings.ReplaceAll(ct, "|", "-")
	ct = strings.ReplaceAll(ct, "/", "-")
	ct = strings.ReplaceAll(ct, "\\", "-")

	return ct
}

func extractName(comicid string) string {
	data := getUrlData(baseUrl)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	name := fmt.Sprint("undefined_%s", comicid)
	infoBlock := doc.Find("#info > h1")
	infoBlock.Each(func(i int, selection *goquery.Selection) {
		for _, n := range selection.Nodes {
			name = n.FirstChild.Data
		}
	})
	return name
}

func extractPageCount() int {
	//#pagination-page-top > button > span.num-pages
	url := fmt.Sprintf("%s1/", baseUrl)
	data := getUrlData(url)

	numPages := "0"
	doc2, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	infoBlock2 := doc2.Find("span.num-pages")
	infoBlock2.Each(func(i int, selection *goquery.Selection) {
		for _, n := range selection.Nodes {
			numPages = n.FirstChild.Data
		}
	})
	pages, err := strconv.Atoi(numPages)
	if err != nil {
		panic(err)
	}
	return pages
}

func extractImageUrl(page int) string {
	url3 := fmt.Sprintf("%s%d/", baseUrl, page)
	data3 := getUrlData(url3)
	imgDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(data3))
	if err != nil {
		panic(err)
	}
	img := imgDoc.Find("#image-container > a > img")
	link, _ := img.Attr("src")

	return link
}
