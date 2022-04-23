package modules

import "github.com/PuerkitoBio/goquery"

type PageLeechModule interface {
	FindNextPage(doc *goquery.Document) string
	FindCurrentImage(doc *goquery.Document) string
}

//------ //

type IROModule struct {

}


func (m *IROModule) FindCurrentImage(doc *goquery.Document) string{
	image := ""
	doc.Find("div#comic").Each(func(i int, selection *goquery.Selection) {
		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			image, _ = selection.Attr("src")
		})
	})
	return image
}

func (m *IROModule) FindNextPage(doc *goquery.Document) string{
	next := ""
	doc.Find("span.nav-next").Each(func(i int, selection *goquery.Selection) {
		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			next, _ = selection.Attr("href")
		})
	})
	return next
}

// ----- //

type OglafModule struct {

}


func (m *OglafModule) FindCurrentImage(doc *goquery.Document) string{
	image := ""
	doc.Find("img#strip").Each(func(i int, selection *goquery.Selection) {
		image, _ = selection.Attr("src")
	})
	return image
}

func (m *OglafModule) FindNextPage(doc *goquery.Document) string{
	next := ""
	doc.Find("a.next").Each(func(i int, selection *goquery.Selection) {
		next, _ = selection.Attr("href")
	})
	return next
}

// ---- //

type Chester5000Module struct {

}


func (m *Chester5000Module) FindCurrentImage(doc *goquery.Document) string{
	image := ""
	doc.Find("div#comic img").Each(func(i int, selection *goquery.Selection) {
		image, _ = selection.Attr("src")
	})
	return image
}

func (m *Chester5000Module) FindNextPage(doc *goquery.Document) string{
	next := ""
	doc.Find("a[rel~=\"next\"]").Each(func(i int, selection *goquery.Selection) {
		next, _ = selection.Attr("href")
	})
	return next
}


// ---- //

type Generic struct {
	NextSelector string
	ImageSelector string
	StopOnURl string
	stop bool
}


func (m *Generic) FindCurrentImage(doc *goquery.Document) string{
	image := ""
	doc.Find(m.ImageSelector).Each(func(i int, selection *goquery.Selection) {
		image, _ = selection.Attr("src")
	})
	return image
}

func (m *Generic) FindNextPage(doc *goquery.Document) string{
	next := ""
	if m.stop {
		return next
	}
	doc.Find(m.NextSelector).Each(func(i int, selection *goquery.Selection) {
		next, _ = selection.Attr("href")
	})
	return next
}

