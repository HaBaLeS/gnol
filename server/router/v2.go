package router

import "github.com/gin-gonic/gin"

func (ah *AppHandler) v2SeriesList(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	gctx.SeriesList = ah.dao.SeriesForUser(gctx.Session.UserId)
	ah.renderTemplate("index.htmx", ctx, gctx)
}
