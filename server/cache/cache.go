//Very simple cache implementation with a hashmap as lookup for images that would otherwise be extracted from an cbz or cbr
//There is not yet any cleanup implemented
//There is not yet any statistics implemented
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
	cacheTable map[string]*CacheEntry
	cacheDir   string
}

//CacheEntry is a mapping between files on Disk and the in memory cache map
type CacheEntry struct {
	id       string
	filename string
	//add things like last access etc...
}

//AddFileToCache makes file Known to the Cache. The filename will be the ID
func (i *ImageCache) AddFileToCache(imagepath string) *CacheEntry {
	_, filename := path.Split(imagepath)
	cacheID := strings.TrimSuffix(filename, filepath.Ext(filename))
	pce := &CacheEntry{
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

//NewImageCache is the constructor for the Cache. Takes a session to access the Configuration
func NewImageCache(config *util.ToolConfig) *ImageCache {
	return &ImageCache{
		cacheTable: make(map[string]*CacheEntry, 100),
		cacheDir:   config.TempDirectory,
	}
}

//RecoverCacheDir will check the configured cache directory for files image files and add them to the cache
//Should be invoked as go routine as this can happen async
func (i *ImageCache) RecoverCacheDir() {
	startPath := i.cacheDir
	filepath.Walk(startPath, func(fp string, info os.FileInfo, err error) error {
		if strings.ToLower(path.Ext(fp)) == ".jpg" {
			i.AddFileToCache(fp)
		}
		return nil
	})
}

func getImageDepricated(comicID string, pageNum int) (ImageLoader, error) {
	/*cacheId := fmt.Sprintf("%s-%d", comicID, pageNum)
	ce, hit := i.cacheTable[cacheId]
	if hit {
		//Check if cache Contains the "comicID-pageNum"
		//Log Cache hit
		//Return loader func or similar to access the data
		fmt.Println("Cache Hit")
		return func() ([]byte, error) {
			return ioutil.ReadFile(ce.filename)
		}, nil
	}

	fmt.Println("Cache Miss")
	me, notfound := i.dao.getMetadata(comicID)
	if notfound != nil {
		return nil, fmt.Errorf("Unknown ComicID: %s", comicID)
	}

	comicDir := path.Join(i.cacheDir, comicID)
	if _, err := os.Stat(comicDir); os.IsNotExist(err) {
		os.Mkdir(comicDir, os.ModePerm)
	}

	//to wait until the whole file is processed
	//alternatively if an be reverted to only process images on deman but enable caching
	//probably for a online solution the better choice

	//TODO add a config parameter to enforce jpeg instead of preserving the original

	cnt := 0
	extractError := me.arc.Walk(me.FilePath, func(f archiver.File) error {
		if !isImageFile(f.Name()) {
			return nil
		}

		ne := &CacheEntry{
			id: fmt.Sprintf("%s-%d", comicID, cnt),
		}
		ne.filename = path.Join(comicDir, ne.id)

		//num, _ := strconv.Atoi(pageNum)
		if cnt == pageNum {
			ce = ne
		}
		cnt++
		out, cerr := os.Create(ne.filename)
		if cerr != nil {
			panic(cerr)
		}

		ext := path.Ext(f.Name())
		newImg, convErr := util.LimitSize(f.ReadCloser, ext, 2560, 1440)
		if convErr != nil {
			return convErr
		}
		io.Copy(out, newImg)
		i.cacheTable[ne.id] = ne

		return nil
	})

	if extractError != nil {
		return nil, extractError
	}

	//If cache does not know the Image
	//Create a lock for comicID ... if there is a lock already return a loader func which is locked waiting for the lock to open
	//
	//If there is no lock, set lock, start the process to create unpack the file and check if thumbnails need to be created
	//Once that's done unlock the lock so any other loaderFuncs can run
	//Return loader func for the requested file

	return func() ([]byte, error) {
		return ioutil.ReadFile(ce.filename)
	}, nil*/
	return nil, nil
}

//ImageLoader is a loader func that allows fast returning of reference and delayed loading of the image
type ImageLoader func() ([]byte, error)
