package lib

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullWorkflow(t *testing.T) {

	infile := "/home/falko/Downloads/comics/Humble Comic Bundle: Critical Role Comics Collection by Dark Horse/Critical Role: The Mighty Nein Origins-Jester Lavorre.pdf"
	cv, err := NewPdfToImages(infile, os.TempDir())
	assert.NoError(t, err, "NewPdfToImages() error")
	if err != nil {
		return
	}

	defer cv.Close()

	cv.Open()

	for true {
		imagePath, nextError := cv.NextImage()
		if errors.Is(nextError, io.EOF) {
			break
		}
		if nextError != nil {
			assert.NoError(t, nextError, "NextImage() error")
			break
		}
		assert.NotEmpty(t, imagePath, "Path should not be empty")
	}
}

func TestIsDir(t *testing.T) {

	tempDir := isDirectory(os.TempDir())
	assert.True(t, tempDir, "Should Exist, there should be a Temdir")

	nonExist := isDirectory("-1")
	assert.False(t, nonExist, "Dir with name -1 should not exist")

}
