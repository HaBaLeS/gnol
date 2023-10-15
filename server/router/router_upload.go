package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
	"strconv"
)

func (ah *AppHandler) uploadArchive(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		panic(err)
	}
	fh := form.File["arc"][0]                            //fixme this will crash for sure when api is misused!
	order, _ := strconv.Atoi(form.Value["ordernum"][0])  //fixme this will crash for sure when api is misused!
	series, _ := strconv.Atoi(form.Value["seriesID"][0]) //fixme this will crash for sure when api is misused!
	outName := path.Join(ah.config.DataDirectory, fh.Filename)
	out, err := os.Create(outName)
	if err != nil {
		panic(fmt.Sprintf("Error creating: %s\n %v", outName, err))
	}
	in, _ := fh.Open()
	_, cpe := io.Copy(out, in)
	if cpe != nil {
		panic(fmt.Sprintf("Error copying to: %s\n %v", outName, cpe))
	}

	jd := &jobs.JobMeta{
		Filename: outName,
		OrderNum: order,
		SeriesId: series,
	}

	us := getUserSession(ctx)
	ah.bgJobs.CreateNewArchiveJob(jd, us.UserID)

	ah.renderTemplate("upload.gohtml", ctx, nil)
}

func (ah *AppHandler) uploadUrl(ctx *gin.Context) {

	url := ctx.PostForm("comicurl")
	us := getUserSession(ctx)

	ah.bgJobs.CreateNewURLJob(url, us.UserID)
	ah.renderTemplate("upload.gohtml", ctx, nil)
}

func (ah *AppHandler) uploadPdf(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		panic(err)
	}
	fh := form.File["pdffile"][0]

	outName := path.Join(os.TempDir(), fh.Filename)
	out, _ := os.Create(outName)
	in, _ := fh.Open()
	_, cpe := io.Copy(out, in)
	if cpe != nil {
		panic(cpe)
	}

	us := getUserSession(ctx)
	ah.bgJobs.CreatePFCConversionJob(outName, us.UserID)

	ah.renderTemplate("upload.gohtml", ctx, nil)
}
