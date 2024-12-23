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

type ComicSite int

const (
	ILIKECOMIX ComicSite = iota
	HDOTNAME
)

func Leech(site ComicSite, url string) string {

	var sitehtml []byte
	switch site {
	case ILIKECOMIX:
		sitehtml = getUrlData(url)
	case HDOTNAME:
		baseUrl = fmt.Sprintf("https://%s", url)
		sitehtml = getUrlData(url)
	}

	//sitehtml := getUrlData(url)

	//baseUrl = fmt.Sprintf("https://%s/g/%s/", site, comicid)
	//baseUrl = fmt.Sprintf("https://%s", url)

	otitle := extractName(sitehtml, "dnf-comic-name")

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
	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)

	//req.Header.Add(":authority", "ilikecomix.com")
	//req.Header.Add(":method", "GET")
	//req.Header.Add(":path", "/western-eroticism/komi-san-cant-fornicate-nudiedoodles/")
	//req.Header.Add(":scheme", "https")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Accept-Language", "de")
	//req.Header.Add("Cache-Control", "max-age=0")
	//req.Header.Add("Dnt", "1")
	//req.Header.Add("Sec-Ch-Ua", "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	//req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	//req.Header.Add("Sec-Ch-Ua-Platform", "\"Linux\"")
	//req.Header.Add("Sec-Fetch-Dest", "document")
	//req.Header.Add("Sec-Fetch-Mode", "navigate")
	//req.Header.Add("Sec-Fetch-Site", "none")
	//req.Header.Add("Sec-Fetch-User", "?1")
	//req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; U; Linux armv7l like Android; en-us) AppleWebKit/531.2+ (KHTML, like Gecko) Version/5.0 Safari/533.2+ Kindle/3.0+")

	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
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

func extractName(html []byte, comicid string) string {
	//data := getUrlData(baseUrl)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
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
