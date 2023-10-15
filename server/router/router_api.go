package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path"
)

const API_USER_ID = "api-user-id"

func (ah *AppHandler) requireAPIToken(ctx *gin.Context) {
	gt := ctx.GetHeader("gnol-token")
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

func (ah *AppHandler) apiListComics(ctx *gin.Context) {
	uidi, _ := ctx.Get(API_USER_ID)
	uid := uidi.(int)

	var comix []storage.Comic
	err := ah.dao.DB.Select(&comix, storage.ALL_COMICS_FOR_USER, uid, storage.NO_TAG_FILTER)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, comix)
}

func (ah *AppHandler) apiUploadComic(ctx *gin.Context) {
	uidi, _ := ctx.Get(API_USER_ID)
	uid := uidi.(int)

	/*sidString := ctx.DefaultQuery("series", "0")
	sid, err := strconv.Atoi(sidString)
	if err != nil {
		panic(err)
	}*/

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
	}
	ah.bgJobs.CreateNewArchiveJob(jd, uid)

	ctx.JSON(http.StatusOK, "Thx for uploading")
}

func (ah *AppHandler) apiSeries(ctx *gin.Context) {
	var series []storage.Series
	err := ah.dao.DB.Select(&series, "select * from series where name != '';")
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, series)
}
