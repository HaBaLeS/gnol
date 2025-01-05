package session

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/fatih/color"
	"github.com/gen2brain/go-fitz"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Session struct {
	TempDir    string
	Verbose    bool
	InputFile  string
	OutputFile string
	MetaData   *CbzMetaData `json:"-"`
	HasErrors  bool
	DryRun     bool
	ListOrder  bool
	From       int
	To         int
	//IssueName     string
	DirectUpload  bool
	GnolHost      string
	ApiToken      string
	MonitorFolder string
	OrderNum      string
	SeriesId      string
}

type CbzMetaData struct {
	Name             string
	Author           string
	Series           string
	NumInSeries      int
	CoverPage        int
	CoverImageBase64 string
	NumPages         int
	Tags             []string
	Genre            []string
	Nsfw             bool
}

const (
	SERIES_ID = "seriesId"
	ORDER_NUM = "orderNum"
)

func NewSession() *Session {
	td, err := os.MkdirTemp("", "gnol_utils")
	if err != nil {
		panic(err)
	}
	s := &Session{
		TempDir: td,
		Verbose: false,
		DryRun:  false,
		MetaData: &CbzMetaData{
			Tags:      make([]string, 0),
			Genre:     make([]string, 0),
			CoverPage: 1,
		},
		From:         0,
		To:           math.MaxInt64,
		DirectUpload: false,
		SeriesId:     "0",
		OrderNum:     "100",
	}
	return s
}

func (s *Session) processOptionsAndValidate(args []string, options map[string]string) bool {
	//XDG Compatible https://farbenmeer.de/blog/the-power-of-the-xdg-base-directory-specification
	confDir := os.Getenv("XDG_CONFIG_HOME")
	if confDir == "" {
		hd, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error reading UserHomeDir: %v", err) //fixme choose one logger
		}
		confDir = path.Join(hd, ".config")
	}
	cf := path.Join(confDir, "gnol", "config.json")
	configJson, err := os.Open(cf)
	defer configJson.Close()
	if err != nil {
		fmt.Printf("ConfigFile not found: %s: %v\n", cf, err) //fixme choose one logger
		if err := os.MkdirAll(path.Join(confDir, "gnol"), os.ModePerm); err != nil {
			log.Fatalf("Could not Create config Dir: %v", err)
		}
		if configJson, err := os.Create(cf); err != nil {
			log.Fatalf("Could not Create config File: %v", err)
		} else {
			json.NewEncoder(configJson).Encode(&Session{})
			configJson.Close()
		}

	} else {
		err := json.NewDecoder(configJson).Decode(s)
		if err != nil {
			log.Panicf("Could not parse %s: %v", cf, err) //fixme proper logger
		}
	}

	if options["apitoken"] != "" {
		s.ApiToken = options["apitoken"]
	}

	if options["gnolhost"] != "" {
		s.GnolHost = options["gnolhost"]
	}

	if s.InputFile == "" {
		s.InputFile = args[0]
	}

	if options["verbose"] != "" {
		s.Verbose = true
	}
	if options["nsfw"] != "" {
		s.MetaData.Nsfw = true
		s.MetaData.Tags = append(s.MetaData.Tags, "nsfw")
	}

	if options["listOrder"] != "" {
		s.DryRun = true
		s.ListOrder = true
	}

	if options["tags"] != "" {
		for _, t := range strings.Split(options["tags"], ",") {
			s.MetaData.Tags = append(s.MetaData.Tags, strings.TrimSpace(t))
		}
	}

	if options[SERIES_ID] != "" {
		s.SeriesId = options[SERIES_ID]
	}
	if options[ORDER_NUM] != "" {
		s.OrderNum = options[ORDER_NUM]
	}

	if options["from"] != "" {
		c, _ := strconv.Atoi(options["from"])
		s.From = c
	}

	if options["to"] != "" {
		c, _ := strconv.Atoi(options["to"])
		s.To = c
	}
	if options["coverpage"] != "" {
		c, _ := strconv.Atoi(options["coverpage"])
		s.MetaData.CoverPage = c
	}

	if options["upload"] != "" {
		s.DirectUpload = true
	}

	if options["name"] != "" {
		s.MetaData.Name = options["name"]
	} else {
		if s.MetaData.Name == "" {
			s.MetaData.Name = path.Base(s.InputFile)
		}
	}

	if options["out_cbz"] != "" {
		s.OutputFile = options["out_cbz"]
	} else {
		dir := s.MetaData.Name
		s.OutputFile = strings.ReplaceAll(dir, " ", "_") + ".cbz"
	}

	if err := s.validate(); err != "" {
		s.Error("Error: %s", err)
		return false
	}

	return true
}

func (s *Session) Log(text string, v ...interface{}) {
	if s.Verbose {
		green := color.New(color.FgGreen).SprintFunc()
		msg := fmt.Sprintf(text, v...)
		log.Printf("[X] %s", green(msg))
	}
}
func (s *Session) Error(text string, v ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	msg := fmt.Sprintf(text, v...)
	log.Printf("[e] %s", red(msg)) //FIXME should be log.panic !?
	s.HasErrors = true
}

func (s *Session) Panic(text string, err error) {
	red := color.New(color.FgRed).SprintFunc()
	log.Panicf("[p] %s, %v", red(text), err) //FIXME should be log.panic !?
}

func (s *Session) Warn(text string, v ...interface{}) {
	log.Printf(text, v...)
	s.HasErrors = true
}

func (s *Session) validate() string {
	if _, err := os.Stat(s.InputFile); err != nil {
		return fmt.Sprintf("File: %s not found", s.InputFile)
	}
	return ""
}

func (s *Session) cleanup() {
	s.Log("Delete TempDir %s", s.TempDir)
	if remErr := os.RemoveAll(s.TempDir); remErr != nil {
		s.Panic("Error deleting Directory", remErr)
	}
}

// fillMetaData extracts pdf metadata and populates session metadata
func (s *Session) fillMetaData(doc *fitz.Document) {
	md := doc.Metadata()
	s.MetaData.NumPages = doc.NumPage()
	s.MetaData.Name = strings.Trim(md["title"], "\x00")
	if s.MetaData.Name == "" {
		s.MetaData.Name = path.Base(s.InputFile)
	}

	s.Log("Extract Metadata from PDF with title: '%s'", s.MetaData.Name)
	s.Log("%d Pages", s.MetaData.NumPages)

	for _, v := range strings.Split(strings.Trim(md["keywords"], "\x00"), " ") {
		s.MetaData.Tags = append(s.MetaData.Tags, v)
	}
	//TODO other metadata might be: Subject, author, creator ....
	s.Log("Extracted Keywords: %v \n\t", s.MetaData.Tags)
}

func (s *Session) SetCoverImage(img image.Image) {
	m := util.Thumbnail(240, 300, img)
	buf := *new(bytes.Buffer)
	if err := jpeg.Encode(&buf, m, nil); err != nil {
		panic(err)
	}
	s.MetaData.CoverImageBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (s *Session) LoadImage(file string) (image.Image, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(fi)
	return i, err
}

func (s *Session) storeAsJpg(idx int, img image.Image) error {
	name := fmt.Sprintf("page%03d.jpg", idx)
	of, err := os.Create(path.Join(s.TempDir, name))
	if err != nil {
		return err
	}
	defer of.Close()
	return jpeg.Encode(of, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
}

func (s *Session) zipFilesTempDir() error {
	s.Log("Create file: %s from Folder: %s", s.OutputFile, s.TempDir)
	ofz, e := os.Create(s.OutputFile)
	if e != nil {
		return e
	}
	defer ofz.Close()

	outw := zip.NewWriter(ofz)
	defer outw.Close()
	err := filepath.Walk(s.TempDir, func(p string, fi fs.FileInfo, err error) error {
		if fi.IsDir() {
			return nil //Ignore the . dir
		}
		w, e := outw.CreateHeader(&zip.FileHeader{
			Name:   fi.Name(),
			Method: zip.Store, //make sure to not compress, helps to unpack much faster!
		})

		if e != nil {
			return e
		}
		r, e := os.Open(p)
		if e != nil {
			return e
		}
		io.Copy(w, r)
		return nil
	})
	return err
}

func (s *Session) WriteMetadataJson() error {
	meta, err := os.Create(path.Join(s.TempDir, "gnol.json"))
	if err != nil {
		return err
	}
	defer meta.Close()
	enc := json.NewEncoder(meta)
	encErr := enc.Encode(s.MetaData)
	if encErr != nil {
		return err
	}
	panic("not good! to printf")
	if s.Verbose {
		out, err := json.MarshalIndent(s.MetaData, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Meta:\n%s\n", out)
	}
	return nil
}
