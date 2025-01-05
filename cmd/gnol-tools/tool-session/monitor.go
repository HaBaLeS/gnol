package session

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var lockMap map[string]int64 = make(map[string]int64)
var procMap map[string]struct{} = make(map[string]struct{})
var doneMap map[string]struct{} = make(map[string]struct{})
var LOCK sync.Mutex

type void struct {
}

func (s *Session) Monitor(args []string, options map[string]string) int {
	s.MonitorFolder = args[0]

	createFolderIfNotExist(s.MonitorFolder)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					lockMap[event.Name] = time.Now().UnixNano()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("file error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path.Join(s.MonitorFolder))
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				scanImportForlder(s.MonitorFolder)
			}
		}
	}()

	ticker2 := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker2.C:
				processList(s)
			}
		}
	}()

	// Block main goroutine forever.
	<-make(chan struct{})

	return 0
}

func createFolderIfNotExist(folder string) {
	processed := path.Join(folder, "processed")
	if !exists(processed) {
		os.Mkdir(processed, os.ModePerm)
	}
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func processList(s *Session) {
	if !LOCK.TryLock() {
		log.Print("Busy still locked")
		return
	}
	if len(procMap) == 0 {
		LOCK.Unlock()
		return
	}
	theTing := ""
	for k, _ := range procMap {
		delete(procMap, k)
		doneMap[k] = void{}
		theTing = k
		break
	}
	td, e := ioutil.TempDir(os.TempDir(), "gnol_utils")
	if e != nil {
		panic(e)
	}
	s.TempDir = td
	s.Verbose = true
	s.MetaData.Name = createName(theTing)
	s.MetaData.Tags = make([]string, 0)
	s.MetaData.NumPages = 0

	if path.Ext(theTing) == ".pdf" {
		log.Printf("PDF to process! %s", theTing)
		s.Convert([]string{theTing}, map[string]string{})
	} else if path.Ext(theTing) == ".cbz" {
		log.Printf("CBZ to process! %s", theTing)
		s.DirectUpload = false
		s.Repack([]string{theTing}, map[string]string{})
	} else {
		log.Printf("NOT processing: %s", theTing)
		LOCK.Unlock()
		return
	}

	//upload cbz
	fmt.Printf("Directly Uploading %s\n", s.OutputFile)
	s.InputFile = s.OutputFile
	s.uploadInternal()
	fmt.Printf("Deleting CBZ %s\n", s.OutputFile)
	os.Remove(s.OutputFile)

	//remove file
	log.Print("Moving file")
	os.Rename(theTing, path.Join(s.MonitorFolder, "processed", path.Base(theTing)))

	s.InputFile = ""
	s.OutputFile = ""

	LOCK.Unlock()
}

func createName(ting string) string {
	name := path.Base(ting)
	name = strings.ReplaceAll(name, path.Ext(name), "")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.Title(strings.ToLower(name))

	return name
}

func scanImportForlder(dirpath string) {
	ls, err := os.ReadDir(dirpath)
	if err != nil {
		panic(err)
	}
	for _, file := range ls {
		zeFile := path.Join(dirpath, file.Name())
		//log.Printf("-> %s <-", zeFile)
		if _, exists := procMap[zeFile]; exists {
			continue
		}
		if _, exists := doneMap[zeFile]; exists {
			continue
		}
		lastWrite := lockMap[zeFile]
		if lastWrite+time.Second.Nanoseconds() < time.Now().UnixNano() {
			procMap[zeFile] = void{}
		}
	}
}
