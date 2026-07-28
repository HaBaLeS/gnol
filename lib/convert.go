package lib

import (
	"fmt"
	"io"
	"math"
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
	OutDir           string //Output directory Default is os.TempDir()
	MaxSize          int    //Size in MB
	Parts            int
	pagesPerPart     int
	CurrentPart      int
	ImageNamePattern string //Pattern to use
	ImgWidth         int    // max Width for Images
	ImgHeight        int    // max Height for Images

	openDoc     *responses.OpenDocument
	CurrentPage int
	NumPages    int
}

func NewPdfToImages(inFile, outdir string) (*PdfToImages, error) {
	retVal := &PdfToImages{
		infile:           inFile,
		MaxSize:          1500, //Megabyte
		ImgHeight:        2560,
		ImgWidth:         1440,
		CurrentPage:      0,
		CurrentPart:      1,
		OutDir:           outdir,
		ImageNamePattern: "page%03d.jpg",
	}

	if !isDirectory(outdir) {
		return nil, fmt.Errorf("%s is not a directory", outdir)
	}

	retVal.init()
	info, err := os.Stat(retVal.infile)
	if err != nil {
		return nil, err
	}

	pdfSize := int(info.Size() / 1024 / 1024)
	retVal.Parts = pdfSize/retVal.MaxSize + 1 //min is 1

	retVal.NumPages, err = api.PageCountFile(retVal.infile)
	if err != nil {
		return nil, err
	}
	retVal.pagesPerPart = retVal.NumPages / retVal.Parts

	return retVal, nil
}

func isDirectory(outdir string) bool {
	info, err := os.Stat(outdir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (p *PdfToImages) init() error {
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
	return err
}

func (p *PdfToImages) open() error {
	//Load the PDF file into a byte array.
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
		//File:           &pdfBytes,
		FileReader:     pdfBytes,
		FileReaderSize: stat.Size(),
		//FilePath: &p.infile,
	})
	if err != nil {
		return err
	}
	p.openDoc = doc

	return nil
}

func (p *PdfToImages) closeDoc() error {
	if p.openDoc != nil {
		// Always close the document, this will release its resources.
		_, err := p.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
			Document: p.openDoc.Document,
		})
		p.openDoc = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PdfToImages) CloseInstance() error {
	err := p.instance.Close()
	if err != nil {
		return err
	}
	p.instance = nil
	return err
}

func (p *PdfToImages) NextImage() (string, error) {
	if p.openDoc == nil {
		err := p.open()
		if err != nil {
			return "", err
		}
	}

	if p.CurrentPage == p.NumPages {
		return "", io.EOF
	}

	outFile := path.Join(p.OutDir, (fmt.Sprintf(p.ImageNamePattern, p.CurrentPage)))

	renderedImg, err := p.instance.RenderToFile(&requests.RenderToFile{
		RenderPageInPixels: &requests.RenderPageInPixels{
			Height:     p.ImgHeight,
			Width:      p.ImgWidth,
			RenderForm: false,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: p.openDoc.Document,
					Index:    p.CurrentPage,
				},
			},
		},
		Progressive:    false,
		OutputTarget:   requests.RenderToFileOutputTargetFile,
		OutputFormat:   requests.RenderToFileOutputFormatJPG,
		TargetFilePath: outFile,
		OutputQuality:  80,
	})

	if err != nil {
		return "", err
	}

	if math.IsInf(renderedImg.PointToPixelRatio, 0) {
		return "", fmt.Errorf("Problem rendering image, reported Pixel Ratio is %f. Page %d", renderedImg.PointToPixelRatio, p.CurrentPage)
	}

	//Close Document"
	if p.CurrentPage == p.CurrentPart*p.pagesPerPart {
		p.CurrentPart++
		err := p.closeDoc()
		if err != nil {
			return "", err
		}
	}

	p.CurrentPage++
	return outFile, nil
}
