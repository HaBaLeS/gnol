package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (ah *AppHandler) comicsInSeriesList(ctx *gin.Context) {
	us := getUserSession(ctx)
	sID := ctx.Param("seriesID")
	us.ComicList = ah.dao.ComicsForUserInSeries(us.UserID, sID)

	ah.renderTemplate("index.gohtml", ctx, nil)
}

func (ah *AppHandler) seriesList(ctx *gin.Context) {
	us := getUserSession(ctx)
	us.SeriesList = ah.dao.SeriesForUser(us.UserID)
	ah.renderTemplate("series.gohtml", ctx, nil)

}

func (ah *AppHandler) createSeries(ctx *gin.Context) {
	//us := getUserSession(ctx)
	name, _ := ctx.GetPostForm("name")
	imgB64, _ := ctx.GetPostForm("previewImage")
	//FIXME this is a hacke ... i cant render the image if this is prefixed, as the template engine rejects it!
	//fix would be to store full string and make the template render correctly
	imgB64 = strings.ReplaceAll(imgB64, "data:image/png;base64,", "")
	imgB64 = strings.ReplaceAll(imgB64, "data:image/jpeg;base64,", "")
	ah.dao.DB.MustExec("insert into series (name, cover_image_base64) values ($1,$2)", name, imgB64)
	ctx.Redirect(http.StatusCreated, "/comics")
}
