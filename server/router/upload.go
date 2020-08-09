
package router

import (
	"io"
	"net/http"
	"os"
	"path"
)

//SetupUploads define routes for File upload like PDF and Comic Archives
func (ah *AppHandler) SetupUploads() {
	ah.Router.Post("/uploadArc", func(w http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 * 1024)
		if err != nil {
			panic(err)
		}
		fh := request.MultipartForm.File["arc"][0]

		outName := path.Join(ah.config.DataDirectory, fh.Filename)
		out, _ := os.Create(outName)
		in, _ := fh.Open()
		_, cpe := io.Copy(out, in)
		if cpe != nil {
			panic(cpe)
		}

		s2 := request.FormValue("public")
		us := getUserSession(request.Context())
		ah.bgJobs.CreateNewArchiveJob(outName, us.UserName, s2)

		ah.renderTemplate("upload.gohtml",w,request,nil)
	})

	ah.Router.Get("/uploadArchive", func(w http.ResponseWriter, req *http.Request) {
		ah.renderTemplate("upload_archive.gohtml", w, req, nil)
	})

	ah.Router.Get("/uploadPdf", func(w http.ResponseWriter, req *http.Request) {
		ah.renderTemplate("upload_pdf.gohtml", w, req, nil)
	})

}
