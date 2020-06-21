package server

import (
	"fmt"
	"github.com/HaBaLeS/gnol/util"
	"github.com/mholt/archiver/v3"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type ImageCache struct {
	dao        *DAOHandler
	cacheTable map[string]*CacheEntry
	cacheDir   string
}

type CacheEntry struct {
	id       string
	filename string
	//add things like last access etc...
}

func NewImageCache(session *Session) *ImageCache {
	return &ImageCache{
		cacheTable: make(map[string]*CacheEntry, 100),
		dao:        session.dao,
		cacheDir:   session.config.TempDirectory,
	}
}

func (i *ImageCache) Balance() {

}

func (i *ImageCache) GetImage(comicID string, pageNum int) (ImageLoader, error) {
	cacheId := fmt.Sprintf("%s-%d", comicID, pageNum)
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

	//FIXME --- learn how to use channels to make a image know to the cache as soon as it is created in order to not have
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
	}, nil
}

type ImageLoader func() ([]byte, error)
