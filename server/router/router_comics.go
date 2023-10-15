package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/command"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

/**
func (ah *AppHandler) comicsList(ctx *gin.Context) {
	us := getUserSession(ctx)
	us.ComicList = ah.dao.ComicsForUser(us.UserID)

	ah.renderTemplate("comic_list.gohtml", ctx, nil)
}*/

func (ah *AppHandler) deleteComic(ctx *gin.Context) {
	comicID := ctx.Param("comicId")
	us := getUserSession(ctx)
	comic := ah.dao.ComicById(comicID)
	ah.dao.DB.MustExec("delete from user_to_comic where comic_id = $1 and user_id = $2", comicID, us.UserID)
	ctx.JSON(200, command.NewRedirectCommand(fmt.Sprintf("/series/%d", comic.SeriesId)))
}

func (ah *AppHandler) comicsLoad(ctx *gin.Context) {
	comicID := ctx.Param("comicId")
	lastPage := ctx.Param("lastpage")
	comicID, _ = url.QueryUnescape(comicID)
	comic := ah.dao.ComicById(comicID)
	if lastPage != "" {
		lp, _ := strconv.Atoi(lastPage) //Fixme ignoring errors is bad
		comic.LastPage = lp
	}

	if comic == nil {
		renderError(fmt.Errorf("comic with id %s not found", comicID), ctx.Writer)
		return
	}
	ah.renderTemplate("jqviewer.gohtml", ctx, comic)
}

func (ah *AppHandler) comicsEdit(ctx *gin.Context) {
	type ComicData struct {
		Issue      *storage.Comic
		SeriesList []storage.Series
	}
	cd := &ComicData{}
	comicID := ctx.Param("comicId")
	cd.Issue = ah.dao.ComicById(comicID)
	cd.SeriesList = ah.dao.AllSeries()
	ah.renderTemplate("edit_comic.gohtml", ctx, cd)
}

func (ah *AppHandler) updateComic(ctx *gin.Context) {
	type ChangeReq struct {
		ComicID  int    `form:"comicID"`
		Name     string `form:"name"`
		Nsfw     string `form:"nsfw"`
		nsfwbool bool
		SeriesID int `form:"seriesID"`
		OrderNum int `form:"orderNum"`
	}
	cr := &ChangeReq{}
	berr := ctx.ShouldBind(cr)
	if berr != nil {
		panic(berr)
	}

	if cr.Nsfw == "on" {
		cr.nsfwbool = true
	}
	us := getUserSession(ctx)
	old := ah.dao.ComicById(strconv.Itoa(cr.ComicID))
	ah.dao.DB.MustExec("update comic set series_id = $1, nsfw = $2, name = $3, orderNum = $6 where id = $4 and ownerID = $5", cr.SeriesID, cr.nsfwbool, cr.Name, cr.ComicID, us.UserID, cr.OrderNum)

	//execute Updates
	ctx.JSON(http.StatusCreated, command.NewRedirectCommand(fmt.Sprintf("/series/%d/", old.SeriesId)))
}

func (ah *AppHandler) comicSetLastPage(ctx *gin.Context) {
	us := getUserSession(ctx)
	comicID := ctx.Param("comicId")
	lastpage := ctx.Param("lastpage")

	ah.dao.DB.MustExec("update user_to_comic set last_page = $1 where user_id = $2 and comic_id = $3", lastpage, us.UserID, comicID)
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
