package conversion

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (j *JobRunner) CreatePFCConversionJob(pdfFile string){
	bgjob := &BGJob{
		JobType: PdfToCbz,
		InputFile: pdfFile,
		DisplayName: "Create CBZ from PDF",
		JobStatus: NotStarted,
	}
	bgjob.save()

}

func convertToPDF(job *BGJob) {
	fmt.Printf("Running conversion\n")
	job.JobStatus = Error //If we do not set  finished its an error
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

	job.JobStatus =Done

}