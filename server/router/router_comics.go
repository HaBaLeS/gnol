package router

import (
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func (ah *AppHandler) comicsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := getUserSession(r.Context())
		if us.ComicList == nil {
			us.ComicList = ah.bs.Comic.GetComiList()
		}
		ah.renderTemplate("index.gohtml", w, r, nil) //TODO move template selection out!
	}
}

func (ah *AppHandler) comicsLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comicID := chi.URLParam(r, "comicId")
		meta, nfe := ah.bs.Comic.GetMetadata(comicID)

		if nfe != nil {
			renderError(nfe, w)
			return
		}
		ah.renderTemplate("jqviewer.gohtml", w, r, meta)
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

		//get file from cache
		var err error
		file, hit := ah.cache.GetFileFromCache(comicID, num)
		if !hit {
			file, err = ah.bs.Comic.GetPageImage(comicID, num)
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


