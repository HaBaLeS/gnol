package jobs

import (
//	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	//"github.com/gen2brain/go-fitz"
/*	"github.com/mholt/archiver/v3"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"*/
)

//CreatePFCConversionJob creates a new job for processing PDF files and crate a CBZ out of it
func (j *JobRunner) CreatePFCConversionJob(pdfFile string,uid int) {
	bgjob := &BGJob{
		JobType:     PdfToCbz,
		InputFile:   pdfFile,
		DisplayName: "Create CBZ from PDF",
		JobStatus:   NotStarted,
		BaseEntity:  storage.CreateBaseEntity(bucketJobOpen),
		UserID: uid,
	}
	j.save(bgjob)

}

func (j *JobRunner) convertToPDF(job *BGJob) error {
	/*fmt.Printf("Running conversion\n")

	tmpDir, err := ioutil.TempDir(os.TempDir(), "fitz")
	if err != nil {
		panic(err)
	}

	doc, err := fitz.New(job.InputFile)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	//Create ZIP
	outZipPath := path.Join(j.cfg.DataDirectory, path.Base(job.InputFile) + ".cbz")
	outZip, cer :=  os.Create(outZipPath)
	if cer != nil {
		panic(cer)
	}
	zip := archiver.NewZip()
	zer := zip.Create(outZip)
	if zer != nil {
		panic(zer)
	}
	defer zip.Close()

	// Extract pages as images
	j.log.InfoF("Processing %d pages from PDF", doc.NumPage() )
	for n := 0; n < doc.NumPage(); n++ {
		img, e1 := doc.Image(n)
		if e1 != nil {
			panic(e1)
		}

		pagename := fmt.Sprintf("page%03d.jpg", n)
		tp := filepath.Join(tmpDir, pagename)
		f, e2 := os.Create(tp)
		if e2 != nil {
			panic(e2)
		}
		defer f.Close()

		m := resize.Thumbnail(2560, 1440, img, resize.Bicubic)
		e3 := jpeg.Encode(f, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if e3 != nil {
			panic(e3)
		}


		fz, e4 := os.Open(tp)
		defer fz.Close()
		if e4!=nil {
			panic(e4)
		}
		info, _ := os.Stat(tp)
		zip.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: pagename,
			},
			ReadCloser: fz,
		})

	}

	j.CreateNewArchiveJob(outZipPath,job.UserID,"")
	//FIXME cleanup unpacked, and tmp
*/
	return nil
}
