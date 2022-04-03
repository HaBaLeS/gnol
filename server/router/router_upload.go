package router

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func (ah *AppHandler) uploadArchive() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 * 1024)
		if err != nil {
			panic(err)
		}
		fh := request.MultipartForm.File["arc"][0]

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

		us := getUserSession(request.Context())
		ah.bgJobs.CreateNewArchiveJob(outName, us.UserID)

		ah.renderTemplate("upload.gohtml", w, request, nil)
	}
}

func (ah *AppHandler) uploadUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.FormValue("comicurl")
		us := getUserSession(r.Context())

		ah.bgJobs.CreateNewURLJob(url, us.UserID)
		ah.renderTemplate("upload.gohtml", w, r, nil)
	}
}

func (ah *AppHandler) uploadPdf() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 * 1024)
		if err != nil {
			panic(err)
		}
		fh := request.MultipartForm.File["pdffile"][0]

		outName := path.Join(os.TempDir(), fh.Filename)
		out, _ := os.Create(outName)
		in, _ := fh.Open()
		_, cpe := io.Copy(out, in)
		if cpe != nil {
			panic(cpe)
		}

		us := getUserSession(request.Context())
		ah.bgJobs.CreatePFCConversionJob(outName, us.UserID)

		ah.renderTemplate("upload.gohtml", w, request, nil)
	}
}
