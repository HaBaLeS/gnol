package conversion

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/gen2brain/go-fitz"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

//CreatePFCConversionJob creates a new job for processing PDF files and crate a CBZ out of it
func (j *JobRunner) CreatePFCConversionJob(pdfFile string) {
	bgjob := &BGJob{
		JobType:     PdfToCbz,
		InputFile:   pdfFile,
		DisplayName: "Create CBZ from PDF",
		JobStatus:   NotStarted,
		BaseEntity:  dao.CreateBaseEntity(),
	}
	j.save(bgjob)

}

func convertToPDF(job *BGJob) int {
	fmt.Printf("Running conversion\n")

	doc, err := fitz.New(job.InputFile)
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	tmpDir, err := ioutil.TempDir(os.TempDir(), "fitz")
	if err != nil {
		panic(err)
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		//FIXME add shrinking to max size here

		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("test%03d.jpg", n)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		f.Close()
	}

	//TODO create ZIP

	//todo cleanup unpacked, and tmp

	return Done
}
