package router

import (
	"github.com/HaBaLeS/gnol/server/command"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (ah *AppHandler) comicsInSeriesList(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	sID := ctx.Param("seriesId")
	gctx.ComicList = ah.dao.ComicsForUserInSeries(gctx.Session.UserId, sID)

	ah.renderTemplate("comic_list.gohtml", ctx, gctx)
}

func (ah *AppHandler) seriesList(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	gctx.SeriesList = ah.dao.SeriesForUser(gctx.Session.UserId)
	ah.renderTemplate("series.gohtml", ctx, gctx)

}

func (ah *AppHandler) createSeries(ctx *gin.Context) {
	//us := getUserSession(ctx)
	name, _ := ctx.GetPostForm("name")
	if name == "" {
		ctx.JSON(500, command.NewValidationErrorCommand("missing name"))
		return
	}
	imgB64, _ := ctx.GetPostForm("previewImage")
	//FIXME this is a hack ... i cant render the image if this is prefixed, as the template engine rejects it!
	//fix would be to store full string and make the template render correctly
	imgB64 = strings.ReplaceAll(imgB64, "data:image/png;base64,", "")
	imgB64 = strings.ReplaceAll(imgB64, "data:image/jpeg;base64,", "")
	ah.dao.DB.MustExec("insert into series (name, cover_image_base64) values ($1,$2)", name, imgB64)
	ctx.JSON(200, command.NewRedirectCommand("/comics"))
}

func (ah *AppHandler) updateSeries(ctx *gin.Context) {
	rc := getGnolContext(ctx)
	type ChangeReq struct {
		SeriesId string `form:"seriesId"`
		Name     string `form:"name"`
		Nsfw     string `form:"nsfw"`
		OrderNum int    `form:"orderNum"`
		nsfwbool bool
	}
	cr := &ChangeReq{}
	if err := ctx.ShouldBind(cr); err != nil {
		panic(err)
	}

	if cr.Nsfw == "on" {
		cr.nsfwbool = true
	}
	ah.dao.DB.MustExec("update series s set name=$3, orderNum=$4, nsfw=$5 where s.id = $1 and s.ownerid = $2", cr.SeriesId, rc.Session.UserId, cr.Name, cr.OrderNum, cr.nsfwbool)

	//execute Updates
	ctx.JSON(http.StatusCreated, command.NewRedirectCommand("/series/"))
}

func (ah *AppHandler) seriesEdit(ctx *gin.Context) {
	rc := getGnolContext(ctx)
	sID := ctx.Param("seriesId")
	rc.Series = ah.dao.SeriesById(sID, rc.Session.UserId)
	ah.renderTemplate("edit_series.gohtml", ctx, rc)
}
