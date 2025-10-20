package router

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/HaBaLeS/gnol/server/command"
	"github.com/HaBaLeS/gnol/server/dto"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-gonic/gin"
)

func (ah *AppHandler) comicsInSeriesList(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	sID := ctx.Param("seriesId")
	coimicList := ah.dao.ComicsForUserInSeries(gctx.Session.UserId, sID)

	arcs := ah.dao.ListSeriesArcs(sID)
	arcs = addDefaultArcIfNecessary(arcs)

	comicByArc := make([]*dto.ArcDTO, 0)
	for _, arc := range arcs {
		arcDto := dto.ArcDTO{
			SeriesArc: arc,
			Comics:    filterComicByArcId(coimicList, arc.Id),
		}
		comicByArc = append(comicByArc, &arcDto)
	}

	gctx.ArcList = comicByArc

	ah.renderTemplate("comic_list.gohtml", ctx, gctx)
}

func addDefaultArcIfNecessary(arcs []*storage.SeriesArc) []*storage.SeriesArc {
	defaultArc := storage.SeriesArc{
		Name:        "Unsorted Arc",
		Description: sql.NullString{"[No Story Arc defined]", true},
		OrderNum:    0,
		Id:          0,
	}
	arcs = append([]*storage.SeriesArc{&defaultArc}, arcs...)
	return arcs
}

func filterComicByArcId(list []*storage.Comic, id int) []*storage.Comic {
	retVal := make([]*storage.Comic, 0)
	for _, comic := range list {
		if comic.ArcId == id {
			retVal = append(retVal, comic)
		}
	}
	return retVal
}

func (ah *AppHandler) seriesList(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	gctx.SeriesList = ah.dao.SeriesForUser(gctx.Session.UserId)
	ah.renderTemplate("series.gohtml", ctx, gctx)
}

func (ah *AppHandler) seriesListNSFW(ctx *gin.Context) {
	gctx := getGnolContext(ctx)
	gctx.SeriesList = ah.dao.NSFWSeriesForUser(gctx.Session.UserId)
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
	var ok bool
	rc.Series, ok = ah.dao.SeriesById(sID, rc.Session.UserId)
	if ok {
		rc.UserList = ah.dao.AllUsers()
		ah.renderTemplate("edit_series.gohtml", ctx, rc)
	} else {
		rc.Flash = "error_not_the_owner_of_series"
		ah.renderTemplate("error.gohtml", ctx, rc)
	}
}

func (ah *AppHandler) xSeriesArcs(ctx *gin.Context) {
	rc := getGnolContext(ctx)
	sID := ctx.Param("seriesId")
	rc.SeriesArcs = ah.dao.ListSeriesArcs(sID)
	ah.renderTemplate("x_arc_table.gohtml", ctx, rc)
}

func (ah *AppHandler) xSeriesAddArc(ctx *gin.Context) {
	type addArc struct {
		SeriesId    string
		Name        string
		Description string
		Link        string
	}

	payload := &addArc{}
	if err := ctx.ShouldBind(payload); err != nil {
		panic(err)
	}

	ah.dao.AddSeriesArc(payload.SeriesId, payload.Name)

	ah.xSeriesArcs(ctx)
}

func (ah *AppHandler) xSeriesArcOptions(ctx *gin.Context) {
	rc := getGnolContext(ctx)
	sID := ctx.Param("seriesId")
	cID := ctx.Param("comicId")
	rc.Issue = ah.dao.ComicById(cID)
	rc.SeriesArcs = ah.dao.ListSeriesArcs(sID)
	ah.renderTemplate("x_arc_options.gohtml", ctx, rc)
}
