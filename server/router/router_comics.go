package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
)

func (ah *AppHandler) comicsList(ctx *gin.Context) {
	us := getUserSession(ctx)
	us.ComicList = ah.dao.ComicsForUser(us.UserID)
	ah.renderTemplate("index.gohtml", ctx, nil)
}

func (ah *AppHandler) seriesList(ctx *gin.Context) {
	us := getUserSession(ctx)
	if us.IsLoggedIn() {
		us.SeriesList = ah.dao.SeriesForUser(us.UserID)
		ah.renderTemplate("index.gohtml", ctx, nil)
	} else {
		//FIXME Render different Template if user is not logged in
		us.ComicList = new([]storage.Comic)
		ah.renderTemplate("index.gohtml", ctx, nil)
	}
}

func (ah *AppHandler) createSeries(ctx *gin.Context) {
	us := getUserSession(ctx)
	if us.IsLoggedIn() {
		us.SeriesList = ah.dao.SeriesForUser(us.UserID)
		ah.renderTemplate("index.gohtml", ctx, nil)
	} else {
		//FIXME Render different Template if user is not logged in
		us.ComicList = new([]storage.Comic)
		ah.renderTemplate("index.gohtml", ctx, nil)
	}
}

func (ah *AppHandler) comicsLoad(ctx *gin.Context) {
	comicID := ctx.Param("comicId")
	comicID, _ = url.QueryUnescape(comicID)
	comic := ah.dao.ComicById(comicID)

	if comic == nil {
		renderError(fmt.Errorf("comic with id %s not found", comicID), ctx.Writer)
		return
	}
	ah.renderTemplate("jqviewer.gohtml", ctx, comic)
}

func (ah *AppHandler) comicsPageImage(ctx *gin.Context) {
	comicID := ctx.Param("comicId")
	image := ctx.Param("imageId")
	num, ce := strconv.Atoi(image)
	if ce != nil {
		renderError(ce, ctx.Writer)
		return
	}

	comic := ah.dao.ComicById(comicID) //FIXME change to getfilename for Comic

	//get file from cache
	var err error
	file, hit := ah.cache.GetFileFromCache(comic.FilePath, num)
	if !hit {
		file, err = storage.GetPageImage(ah.config, comic.FilePath, comicID, num)
		if err != nil {
			renderError(err, ctx.Writer)
			return
		}
		ah.cache.AddFileToCache(file)
	}

	//as a image-provider module not the cache directly
	img, oerr := os.Open(file)
	if oerr != nil {
		renderError(oerr, ctx.Writer)
		return
	}

	data, rerr := ioutil.ReadAll(img)
	if rerr != nil {
		renderError(rerr, ctx.Writer)
		return
	}
	_, re := ctx.Writer.Write(data)
	if re != nil {
		panic(re)
	}
}
