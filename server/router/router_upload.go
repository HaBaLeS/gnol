
package router

import (
	"io"
	"net/http"
	"os"
	"path"
)

func (ah *AppHandler) uploadArchive() http.HandlerFunc{
	return func(w http.ResponseWriter, request *http.Request) {
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
	}
}

