package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/dto"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	API_GNOL_TOKEN = "gnol-token"
	API_USER_ID    = "api-user-id"
	API_SERIES_ID  = "series-id"
	API_NSFW       = "nsfw"
	API_ODER_NUM   = "order-num"
)

func (ah *AppHandler) requireAPIToken(ctx *gin.Context) {
	gt := ctx.GetHeader(API_GNOL_TOKEN)
	if "" == gt {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "missing header 'gnol-token'")
	}
	err, uid := ah.dao.GetUserForApiToken(gt)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "unknown gnol token")
	}
	ctx.Set(API_USER_ID, uid)
	ctx.Next()
}

// @Summary Get list of dto.ComicEntry in gnol
// @Schemes
// @Description tbd
// @Tags Comic Management
// @Produce json
// @Success 200 {object} []dto.ComicEntry
// @Router  /list [get]
// @Security ApiKeyAuth
func (ah *AppHandler) apiListComics(ctx *gin.Context) {
	uidi, _ := ctx.Get(API_USER_ID)
	uid := uidi.(int)
	query := "select c.id, c.name, c.series_id, s.name as \"sname\", c.nsfw, c.num_pages, c.sha256sum from comic c, series s where c.series_id = s.id and c.ownerid = $1 order by c.id;"

	var resList = []dto.ComicEntry{}
	err := ah.dao.DB.Select(&resList, query, uid)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, resList)
}

// @Summary Upload a cbz file
// @Schemes
// @Description tbd
// @Tags Upload
// @Produce json
// @Param   series-id      query     integer     false  "string valid"
// @Param   nsfw      query     string     false  "string valid"
// @Param   order-num      query     string     false  "string valid"
// @Success 200 {string} string
// @Router  /upload [post]
// @Security ApiKeyAuth
func (ah *AppHandler) apiUploadComic(ctx *gin.Context) {
	uid := 0
	seriesId := 0
	nsfw := false
	orderNum := 100

	if uidi, exist := ctx.Get(API_USER_ID); !exist {
		panic("Missing patrameter api-user-id")
	} else {
		uid = uidi.(int)
	}
	if val := ctx.Query(API_SERIES_ID); val != "" {
		seriesId, _ = strconv.Atoi(val)
	}
	if val := ctx.Query(API_NSFW); val != "" {
		nsfw = true
	}
	if val := ctx.Query(API_ODER_NUM); val != "" {
		orderNum, _ = strconv.Atoi(val)
	}

	fmt.Printf("Storing file for uid: %d", uid)

	fn := path.Join(ah.config.DataDirectory, uuid.New().String())

	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	num, err := io.Copy(f, ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d bytes written\n", num)
	jd := &jobs.JobMeta{
		Filename: fn,
		SeriesId: seriesId,
		OrderNum: orderNum,
		Nsfw:     nsfw,
	}
	ah.bgJobs.CreateNewArchiveJob(jd, uid)

	ctx.JSON(http.StatusOK, "Thx for uploading")
}

// @Summary Get list of dto.Series in gnol
// @Schemes
// @Description tbd
// @Tags Series Mangement
// @Produce json
// @Success 200 {object} []storage.Series
// @Router  /series [get]
// @Security ApiKeyAuth
func (ah *AppHandler) apiSeries(ctx *gin.Context) {
	var series []storage.Series
	err := ah.dao.DB.Select(&series, "select Id, Name from series order by Id;")
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, series)
}

// @Summary Get dto.ComicEntry for hash
// @Schemes
// @Description Check if a comic with the given cbz hash
// @Param   hash      path     string     false  "string valid"
// @Tags Upload
// @Produce json
// @Success 200 {object} dto.ComicEntry
// @Router  /checkhash/:hash [get]
// @Security ApiKeyAuth
func (ah *AppHandler) apiCheckHash(ctx *gin.Context) {
	uidi, _ := ctx.Get(API_USER_ID)
	hash := ctx.Param("hash")
	var retVal dto.ComicEntry
	query := "select c.id, c.name, c.series_id, s.name as \"sname\", c.nsfw, c.num_pages, c.sha256sum from comic c, series s where c.series_id = s.id and c.ownerid = $1 and c.sha256sum = $2;"
	err := ah.dao.DB.Get(&retVal, query, uidi, hash)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "no file with that hash for user")
	} else {
		ctx.JSON(http.StatusOK, retVal)
	}
}
