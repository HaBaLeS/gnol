package router

import (
	"io"
	"net/http"
	"os"
	"path"
)

func (r *AppHandler) SetupUploads() {
	r.Router.Post("/uploadArc", func(w http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 * 1024)
		if err != nil {
			panic(err)
		}
		fh := request.MultipartForm.File["arc"][0]

		outName := path.Join(r.config.DataDirectory, fh.Filename)
		out, _ := os.Create(outName)
		in, _ := fh.Open()
		io.Copy(out, in)

		s2 := request.FormValue("public")
		us := getUserSession(request.Context())
		r.bgJobs.CreateNewArchiveJob(outName, us.UserName, s2)

		tpl, err := r.getTemplate("upload.gohtml")
		if err != nil {
			panic(err)
		}

		err = tpl.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	})

	r.Router.Get("/uploadArchive", func(w http.ResponseWriter, req *http.Request) {
		tpl, err := r.getTemplate("upload_archive.gohtml")
		if err != nil {
			panic(err)
		}

		err = renderTemplate(tpl, w, req, nil)
		if err != nil {
			panic(err)
		}
	})

	r.Router.Get("/uploadPdf", func(w http.ResponseWriter, req *http.Request) {
		tpl, err := r.getTemplate("upload_pdf.gohtml")
		if err != nil {
			panic(err)
		}

		err = renderTemplate(tpl, w, req, nil)
		if err != nil {
			panic(err)
		}
	})

}
