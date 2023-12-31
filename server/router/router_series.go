package router

import (
	"github.com/HaBaLeS/gnol/server/command"
	"github.com/gin-gonic/gin"
	"strings"
)

func (ah *AppHandler) comicsInSeriesList(ctx *gin.Context) {
	us := getUserSession(ctx)
	sID := ctx.Param("seriesId")
	cl := ah.dao.ComicsForUserInSeries(us.UserID, sID)

	ah.renderTemplate("comic_list.gohtml", ctx, &RenderContext{ComicList: cl, USess: us})
}

func (ah *AppHandler) seriesList(ctx *gin.Context) {
	us := getUserSession(ctx)
	sl := ah.dao.SeriesForUser(us.UserID)
	ah.renderTemplate("series.gohtml", ctx, &RenderContext{SeriesList: sl, USess: us})

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
	rc := NewRenderContext(ctx)
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
	ah.dao.DB.MustExec("update series s set name=$3, orderNum=$4, nsfw=$5 where s.id = $1 and s.ownerid = $2", cr.SeriesId, rc.USess.UserID, cr.Name, cr.OrderNum, cr.nsfwbool)
}

func (ah *AppHandler) seriesEdit(ctx *gin.Context) {
	rc := NewRenderContext(ctx)
	sID := ctx.Param("seriesId")
	rc.Series = ah.dao.SeriesById(sID, rc.USess.UserID)
	ah.renderTemplate("edit_series.gohtml", ctx, rc)
}
