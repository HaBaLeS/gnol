package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func (ah *AppHandler) comicsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := getUserSession(r.Context())
		if us.IsLoggedIn()  {
			us.ComicList = ah.dao.ComicsForUser(us.UserID)
			ah.renderTemplate("index.gohtml", w, r, nil)
		} else {
			//FIXME Render different Template if user is not logged in
			us.ComicList = new([]storage.Comic)
			ah.renderTemplate("index.gohtml", w, r, nil)
		}
	}
}

func (ah *AppHandler) seriesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := getUserSession(r.Context())
		if us.IsLoggedIn()  {
			us.SeriesList = ah.dao.SeriesForUser(us.UserID)
			ah.renderTemplate("index.gohtml", w, r, nil)
		} else {
			//FIXME Render different Template if user is not logged in
			us.ComicList = new([]storage.Comic)
			ah.renderTemplate("index.gohtml", w, r, nil)
		}
	}
}


func (ah *AppHandler) comicsLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comicID := chi.URLParam(r, "comicId")
		comicID,_ = url.QueryUnescape(comicID)
		comic  := ah.dao.ComicById(comicID)

		if comic == nil {
			renderError(fmt.Errorf("comic with id %s not found", comicID), w)
			return
		}
		ah.renderTemplate("jqviewer.gohtml", w, r, comic)
	}
}

func (ah *AppHandler) comicsPageImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comicID := chi.URLParam(r, "comicId")
		image := chi.URLParam(r, "imageId")
		num, ce := strconv.Atoi(image)
		if ce != nil {
			renderError(ce, w)
			return
		}

		comic := ah.dao.ComicById(comicID) //FIXME change to getfilename for Comic

		//get file from cache
		var err error
		file, hit := ah.cache.GetFileFromCache(comic.FilePath, num)
		if !hit {
			file, err = ah.bs.Comic.GetPageImage(comic.FilePath,comicID, num)
			if err != nil {
				renderError(err, w)
				return
			}
			ah.cache.AddFileToCache(file)
		}

		//as a image-provider module not the cache directly
		img, oerr := os.Open(file)
		if oerr != nil {
			renderError(oerr, w)
			return
		}

		data, rerr := ioutil.ReadAll(img)
		if rerr != nil {
			renderError(rerr, w)
			return
		}
		_, re := w.Write(data)
		if re != nil {
			panic(re)
		}
	}
}


