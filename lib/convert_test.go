package lib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullWorkflow(t *testing.T) {

	infile := "/home/falko/Downloads/oglaf_book_3-dl.pdf"

	qual := 95
	os.Mkdir(path.Join(os.TempDir(), fmt.Sprintf("%d", qual)), os.ModePerm)
	conv, err := NewPdfToImages(infile, path.Join(os.TempDir(), fmt.Sprintf("%d", qual)), qual, false)
	assert.NoError(t, err, "NewPdfToImages() error")
	if err != nil {
		return
	}

	defer conv.Close()

	cnt := 0
	for true {
		imagePath, nextError := conv.NextImage()
		if errors.Is(nextError, io.EOF) {
			break
		}
		if nextError != nil {
			assert.NoError(t, nextError, "NextImage() error")
			break
		}
		assert.NotEmpty(t, imagePath, "Path should not be empty")
		cnt++
		if cnt > 10 {
			break
		}
	}
}

func TestIsDir(t *testing.T) {

	tempDir := isDirectory(os.TempDir())
	assert.True(t, tempDir, "Should Exist, there should be a Temdir")

	nonExist := isDirectory("-1")
	assert.False(t, nonExist, "Dir with name -1 should not exist")

}
