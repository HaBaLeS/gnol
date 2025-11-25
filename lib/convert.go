package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfToImages struct {

	// Be sure to close pools/instances when you're done with them.
	instance         pdfium.Pdfium
	infile           string
	outDir           string
	maxSize          int
	parts            int
	pagesPerPart     int
	imageNamePattern string
	imgWidth         int
	imgHeight        int

	openDoc     *responses.OpenDocument
	currentPage int
	numPages    int
}

func NewPdfToImages(inFile string, maxSize int) (*PdfToImages, error) {
	retVal := &PdfToImages{
		infile:      inFile,
		maxSize:     maxSize,
		imgHeight:   2560,
		imgWidth:    1440,
		currentPage: 0,
		outDir:      "/tmp",
	}

	retVal.init()
	info, err := os.Stat(retVal.infile)
	if err != nil {
		return nil, err
	}

	pdfSize := int(info.Size() / 1024 / 1024)
	retVal.parts = pdfSize/maxSize + 1 //min is 1

	retVal.numPages, err = api.PageCountFile(retVal.infile)
	if err != nil {
		return nil, err
	}
	retVal.pagesPerPart = retVal.numPages / retVal.parts
	retVal.instance.Close()

	return retVal, nil
}

func (p *PdfToImages) init() {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 2, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		panic(err)
	}

	p.instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PdfToImages) Open() error {
	// Load the PDF file into a byte array.
	pdfBytes, err := os.Open(p.infile)
	if err != nil {
		return err
	}
	stat, err := os.Stat(p.infile)
	if err != nil {
		return err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := p.instance.OpenDocument(&requests.OpenDocument{
		//File: &pdfBytes,
		FileReader:     pdfBytes,
		FileReaderSize: stat.Size(),
	})
	if err != nil {
		return err
	}
	p.openDoc = doc

	return nil
}

func (p *PdfToImages) Close() error {
	// Always close the document, this will release its resources.
	_, err := p.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: p.openDoc.Document,
	})
	return err
}

func (p *PdfToImages) NextImage() (string, error) {
	if p.currentPage == p.numPages {
		return "", io.EOF
	}

	outFile := path.Join(p.outDir, (fmt.Sprintf(p.imageNamePattern, p.currentPage)))

	_, err := p.instance.RenderToFile(&requests.RenderToFile{
		RenderPageInPixels: &requests.RenderPageInPixels{
			Height:     p.imgWidth,
			Width:      p.imgHeight,
			RenderForm: false,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: p.openDoc.Document,
					Index:    p.currentPage,
				},
			},
		},
		Progressive:    false,
		OutputTarget:   requests.RenderToFileOutputTargetFile,
		OutputFormat:   requests.RenderToFileOutputFormatJPG,
		TargetFilePath: outFile,
	})

	if err != nil {
		return "", err
	}

	p.currentPage++
	return outFile, nil
}
