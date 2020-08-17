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
			if us.MetadataList == nil || len(us.MetadataList.Comics) == 0 {
				user := ah.bs.User.UserByID([]byte(us.UserID))
				us.MetadataList = ah.bs.Comic.MetadataForList(user.MetadataList)
			}
			ah.renderTemplate("index.gohtml", w, r, nil)
		} else {
			//FIXME Render different Template if user is not logges in
			us.MetadataList = &storage.MetadataList{
				Comics: make([]*storage.Metadata,0),
			}
			ah.renderTemplate("index.gohtml", w, r, nil)
		}
	}
}

func (ah *AppHandler) comicsLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comicID := chi.URLParam(r, "comicId")
		comicID,_ = url.QueryUnescape(comicID)
		meta  := ah.bs.Comic.GetMetadata([]byte(comicID))

		if meta == nil {
			renderError(fmt.Errorf("comic with id %s not found", comicID), w)
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


