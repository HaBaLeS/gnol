package router

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/HaBaLeS/gnol/server/command"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/gin-gonic/gin"
)

func (ah *AppHandler) deleteComic(ctx *gin.Context) {
	comicId := ctx.Param("comicId")
	uid := getGnolContext(ctx).Session.UserId
	comic := ah.dao.ComicById(comicId)
	if comic.OwnerID != uid {
		panic("you are not the owner of this comic!")
	}
	ah.dao.DB.MustExec("delete from user_to_comic where comic_id = $1", comicId)
	ah.dao.DB.MustExec("delete from comic where id = $1", comicId)
	err := os.Remove(comic.FilePath)
	if err != nil {
		panic(fmt.Errorf("failed to remove file: %s. got error %v", comic.FilePath, err)) //does the transaction roll back?
	}
	ctx.JSON(200, command.NewRedirectCommand(fmt.Sprintf("/series/%d", comic.SeriesId)))
}

func (ah *AppHandler) removeComic(ctx *gin.Context) {
	comicId := ctx.Param("comicId")
	uid := getGnolContext(ctx).Session.UserId
	comic := ah.dao.ComicById(comicId)
	ah.dao.DB.MustExec("delete from user_to_comic where comic_id = $1 and user_id = $2", comicId, uid)
	ctx.JSON(200, command.NewRedirectCommand(fmt.Sprintf("/series/%d", comic.SeriesId)))
}

func (ah *AppHandler) comicsLoad(ctx *gin.Context) {
	comicId := ctx.Param("comicId")
	lastPage := ctx.Param("lastpage")
	comicId, _ = url.QueryUnescape(comicId)

	comic := ah.dao.ComicById(comicId)
	if lastPage != "" {
		lp, _ := strconv.Atoi(lastPage) //Fixme ignoring errors is bad
		comic.LastPage = lp
	}

	if comic == nil {
		renderError(fmt.Errorf("comic with id %s not found", comicId), ctx.Writer)
		return
	}

	gctx := getGnolContext(ctx)
	gctx.Issue = comic
	ah.renderTemplate("jqviewer.gohtml", ctx, gctx)
}

func (ah *AppHandler) comicsEdit(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	comicID := ctx.Param("comicId")
	gctx.Issue = ah.dao.ComicById(comicID)
	nsfwUser := getGnolContext(ctx).UserInfo.Nsfw
	gctx.SeriesList = ah.dao.AllSeries(nsfwUser)

	gctx.UserList = ah.dao.AllUsers()
	ah.renderTemplate("edit_comic.gohtml", ctx, gctx)
}

func (ah *AppHandler) updateComic(ctx *gin.Context) {
	type ChangeReq struct {
		ComicID  int    `form:"comicID"`
		Name     string `form:"name"`
		Nsfw     string `form:"nsfw"`
		nsfwbool bool
		SeriesID int `form:"seriesID"`
		OrderNum int `form:"orderNum"`
		ArcId    int `form:"seriesArcID"`
	}
	cr := &ChangeReq{}
	berr := ctx.ShouldBind(cr)
	if berr != nil {
		panic(berr)
	}

	if cr.Nsfw == "on" {
		cr.nsfwbool = true
	}
	us := getGnolContext(ctx).Session
	old := ah.dao.ComicById(strconv.Itoa(cr.ComicID))

	// if comic is moved to another series set arcID to o (default)
	if old.SeriesId != cr.SeriesID {
		cr.ArcId = 0
	}

	if us.UserId == old.OwnerID {
		ah.dao.DB.MustExec("update comic set series_id = $1, nsfw = $2, name = $3, orderNum = $6, arcId =$7 where id = $4 and ownerID = $5", cr.SeriesID, cr.nsfwbool, cr.Name, cr.ComicID, us.UserId, cr.OrderNum, cr.ArcId)
	} else {
		//TODO Log or panic error user is not allowed to do that
	}

	//execute Updates
	ctx.JSON(http.StatusCreated, command.NewRedirectCommand(fmt.Sprintf("/series/%d/", old.SeriesId)))
}

func (ah *AppHandler) shareComic(ctx *gin.Context) {
	us := getGnolContext(ctx).Session

	comicId := ctx.Param("comicId")
	targetUser := ctx.Param("userId")

	//fixme add some validation code here ... unsure which pattern to use!
	comic := ah.dao.ComicById(comicId)
	if comic.OwnerID == us.UserId {
		ah.dao.AddComicToUser(comicId, targetUser)
		ctx.JSON(http.StatusOK, command.NewGoBackCommand())
	} else {
		ctx.JSON(http.StatusForbidden, "Only Owner of comic can share it")
	}

}

func (ah *AppHandler) shareSeries(ctx *gin.Context) {
	us := getGnolContext(ctx).Session

	seriesId := ctx.Param("seriesId")
	targetUser := ctx.Param("userId")

	_, ok := ah.dao.SeriesByIdAndUser(seriesId, us.UserId)
	if !ok {
		ctx.JSON(http.StatusForbidden, "Only Owner of comic can share it")
	}
	ids := make([]int, 0)
	err := ah.dao.DB.Select(&ids, "select c.id  from comic c join user_to_comic utc on utc.comic_id = c.id where utc.user_id  = $1 and c.series_id = $2", us.UserId, seriesId)
	if err != nil {
		panic(err)
	}
	for _, id := range ids {
		ah.dao.AddComicToUser(strconv.Itoa(id), targetUser)
	}
	ctx.JSON(http.StatusOK, command.NewGoBackCommand())
}

func (ah *AppHandler) comicSetLastPage(ctx *gin.Context) {
	us := getGnolContext(ctx).Session
	comicID := ctx.Param("comicId")
	lastpage := ctx.Param("lastpage")

	ah.dao.DB.MustExec("update user_to_comic set last_page = $1 where user_id = $2 and comic_id = $3", lastpage, us.UserId, comicID)

	num, _ := strconv.Atoi(lastpage)
	if num+1 >= ah.dao.ComicById(comicID).NumPages {
		ah.dao.SetFinished(us.UserId, comicID)
	}
}

func (ah *AppHandler) comicsPageImage(ctx *gin.Context) {
	comicId := ctx.Param("comicId")
	imageNum := ctx.Param("imageId")
	num, ce := strconv.Atoi(imageNum)
	if ce != nil {
		renderError(ce, ctx.Writer)
		return
	}

	data, err := ah.fileStorage.FetchImageData(comicId, num)
	if err != nil {
		renderError(err, ctx.Writer)
	}

	_, re := ctx.Writer.Write(data)
	if re != nil {
		panic(re)
	}
}

func (ah *AppHandler) recreateComicCover(ctx *gin.Context) {
	type reqObj struct {
		ComicId string `form:"comicId" binding:"required"`
		Page    int    `form:"page" binding:"required"`
	}
	var req reqObj
	if err := ctx.ShouldBind(&req); err != nil {
		renderError(err, ctx.Writer)
	}

	data, err := ah.fileStorage.FetchImageData(req.ComicId, req.Page)
	if err != nil {
		panic(fmt.Errorf("fetch image page %v error: %v", req.Page, err))
	}

	//fixme move to a util
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("error decoding image: %v", err))
	}

	//fixme move to a util
	m := util.Thumbnail(240, 300, img)
	buf := *new(bytes.Buffer)
	if err := jpeg.Encode(&buf, m, nil); err != nil {
		panic(err)
	}

	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	ah.dao.DB.MustExec("update comic set cover_image_base64 = $1 where id = $2", enc, req.ComicId)

	ctx.JSON(http.StatusOK, command.NewGoBackCommand())
}

func (ah *AppHandler) replaceComicCover(ctx *gin.Context) {

}
