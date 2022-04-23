//Package cache offerst cached access to image files.
//Files in cache are resized and converted to save traffic
package cache

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//ImageCache holds the data for the cache and references to other modules needed
type ImageCache struct {
	cacheTable map[string]CacheEntry
	cacheDir   string
}

//CacheEntry is a mapping between files on Disk and the in memory cache map
type CacheEntry struct {
	id       string
	filename string
	//add things like last access etc...
}

//AddFileToCache makes file Known to the Cache. The filename will be the ID
func (i *ImageCache) AddFileToCache(imagepath string) CacheEntry {
	_, filename := path.Split(imagepath)
	cacheID := strings.TrimSuffix(filename, filepath.Ext(filename))
	pce := CacheEntry{
		id:       cacheID,
		filename: imagepath,
	}
	i.cacheTable[pce.id] = pce
	return pce
}

//GetFileFromCache checks if a CacheEntry for "comicID-pageNum" exists. If it exists it returns the filename of the file.
//If it does not exist second return parameter hit will be false
func (i *ImageCache) GetFileFromCache(comicID string, pageNum int) (string, bool) {
	id := fmt.Sprintf("%s-%d", comicID, pageNum)
	ce, hit := i.cacheTable[id]
	if !hit {
		return "", hit
	}
	return ce.filename, hit
}

//NewImageCache is the constructor for the Cache. Takes a gnolsession to access the Configuration
func NewImageCache(config *util.ToolConfig) *ImageCache {
	return &ImageCache{
		cacheTable: make(map[string]CacheEntry, 100),
		cacheDir:   config.TempDirectory,
	}
}

//RecoverCacheDir will check the configured cache directory for files image files and add them to the cache
//Should be invoked as go routine as this can happen async
func (i *ImageCache) RecoverCacheDir() {
	startPath := i.cacheDir
	filepath.Walk(startPath, func(fp string, info os.FileInfo, err error) error {
		if strings.ToLower(path.Ext(fp)) == ".gnol" {
			i.AddFileToCache(fp)
		}
		return nil
	})
}
