package cache

import (
	"github.com/HaBaLeS/gnol/server/util"
	"path/filepath"
	"testing"
)

var (
	mockConfig = &util.ToolConfig{}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestNewImageCache(t *testing.T) {
	NewImageCache(mockConfig)
}

func TestImageCache_AddFileToCache(t *testing.T) {
	c := NewImageCache(mockConfig)
	ce := c.AddFileToCache("/lot/of/path/comicid/comicid-35.jpg")

	got := ce.id
	want := "comicid-35"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestImageCache_GetFileFromCache(t *testing.T) {
	c := NewImageCache(mockConfig)
	c.AddFileToCache("/lot/of/path/comicid/comicid-34.jpg")

	got, hit := c.GetFileFromCache("asdf", 234)
	if hit {
		t.Errorf("got '%s' where nothing was expected", got)
	}
	got, hit = c.GetFileFromCache("comicid", 34)
	want := "/lot/of/path/comicid/comicid-35"
	if !hit {
		t.Errorf("got '%s' where '%s' was expected", got, want)
	}

}

func TestImageCache_RecoverCacheDir(t *testing.T) {
	absPath, _ := filepath.Abs("testdata/")
	mockConfig.TempDirectory = absPath
	c := NewImageCache(mockConfig)

	c.RecoverCacheDir()

	_, hit := c.GetFileFromCache("comicID", 666)
	if !hit {
		t.Errorf("Did not find recovered file!")
	}
}
