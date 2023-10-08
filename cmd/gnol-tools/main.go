package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	"github.com/teris-io/cli"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

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

type Session struct {
	TempDir    string
	Verbose    bool
	Logger     *log.Logger `json:"-"`
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
}

var VersionNum = "undefined"
var BuildDate = "undefined"

func NewSession() *Session {
	td, err := ioutil.TempDir(os.TempDir(), "gnol_utils")
	if err != nil {
		panic(err)
	}
	return &Session{
		Logger:  log.Default(),
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
	}
}

func main() {
	s := NewSession()

	inFile := cli.NewArg("infile", "CBZ/CBR to process")
	inPdfArg := cli.NewArg("inpdf", "Input PDF")
	inDirArg := cli.NewArg("indir", "Input Folder")
	verbose := cli.NewOption("verbose", "Verbose Logging").WithType(cli.TypeBool).WithChar('v')
	gnolHost := cli.NewOption("gnolhost", "GnolHost and Path").WithChar('h')
	apiToken := cli.NewOption("apitoken", "API-Token to use").WithChar('k')
	tags := cli.NewOption("tags", "Comma separated list of Tags for Metadata").WithType(cli.TypeString).WithChar('t')
	nsfw := cli.NewOption("nsfw", "Mark Graphic Novel as NSFW").WithType(cli.TypeBool).WithChar('x')
	coverImage := cli.NewOption("coverpage", "Select page to use a cover. Starting from 1").WithType(cli.TypeInt).WithChar('c')
	outFile := cli.NewOption("out_cbz", "Output file").WithType(cli.TypeString).WithChar('o')
	listOrder := cli.NewOption("listOrder", "preview order of file.(e.g. or cover selection) CBZ will not be created").WithChar('l').WithType(cli.TypeBool)

	from := cli.NewOption("from", "StartPage Default 0 ").WithType(cli.TypeInt)
	to := cli.NewOption("to", "LastPage Default 0").WithType(cli.TypeInt)
	name := cli.NewOption("name", "Name of Issue/Novel").WithType(cli.TypeString).WithChar('n')
	upload := cli.NewOption("upload", "Directly upload CBZ after creation").WithType(cli.TypeBool).WithChar('u')

	pdf2cbz := cli.NewCommand("pdf2cbz", "PDF to CBZ/CBR converter with support for GNOL Metadata").
		WithArg(inPdfArg).
		WithOption(outFile).
		WithOption(tags).
		WithOption(nsfw).
		WithOption(coverImage).
		WithOption(name).
		WithAction(s.convert)

	folder2cbz := cli.NewCommand("folder2cbz", "Pack folder of images to CBZ with support for GNOL Metadata. Files will be converted to JPEG and Downsized").
		WithArg(inDirArg).
		WithOption(outFile).
		WithOption(tags).
		WithOption(nsfw).
		WithOption(coverImage).
		WithOption(listOrder).
		WithOption(upload).
		WithOption(name).
		WithAction(s.packfolder)

	uploadcmd := cli.NewCommand("upload", "Upload CBZ/CBR to a Gnol instance").
		WithArg(inFile).
		WithAction(s.upload)

	repack := cli.NewCommand("repack", "Repackage a CBZ/CBR. Remove compression, Images Downsized if neccesary and add/update of GNOL Metadata").
		WithArg(inFile).
		WithOption(tags).
		WithOption(nsfw).
		WithOption(coverImage).
		WithOption(from).
		WithOption(to).
		WithOption(listOrder).
		WithOption(name).
		WithOption(upload).
		WithAction(s.repack)

	series := cli.NewCommand("series", "Gnol Series management Command. See subcommands for details").
		WithCommand(cli.NewCommand("list", "list existing series").WithAction(listSeries)).
		WithCommand(cli.NewCommand("create", "create a new series").WithArg(cli.NewArg("name", "Name for Series")).WithAction(createSeries))

	version := cli.NewCommand("version", "Print Version number").WithAction(func(args []string, options map[string]string) int {
		fmt.Printf("gnol-tools %s from %s\n", VersionNum, BuildDate)
		return 0
	})

	monitor := cli.NewCommand("monitor", "Monitor folder and auto-process files in it").WithArg(inDirArg).WithAction(s.monitor)

	app := cli.New("CLI utils for GNOL").
		WithCommand(pdf2cbz).
		WithCommand(folder2cbz).
		WithCommand(uploadcmd).
		WithCommand(repack).
		WithCommand(series).
		WithCommand(version).
		WithCommand(monitor).
		WithOption(verbose).
		WithOption(gnolHost).
		WithOption(apiToken)

	os.Exit(app.Run(os.Args, os.Stdout))
}

func listSeries(args []string, options map[string]string) int {
	fmt.Printf("Series: XXXX")
	return 0
}

func createSeries(args []string, options map[string]string) int {
	fmt.Printf("Comics: XXXX")
	return 0
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

	s.InputFile = args[0]

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
		s.Logger.Printf(text, v...)
	}
}
func (s *Session) Error(text string, v ...interface{}) {
	s.Logger.Printf(text, v...)
	s.HasErrors = true
}

func (s *Session) validate() string {
	if _, err := os.Stat(s.InputFile); err != nil {
		return fmt.Sprintf("File: %s not found", s.InputFile)
	}
	return ""
}

func (s *Session) cleanup() {
	if remErr := os.RemoveAll(s.TempDir); remErr != nil {
		panic(remErr)
	}
}

func (s *Session) fillMetaData(doc *fitz.Document) {
	md := doc.Metadata()
	s.MetaData.NumPages = doc.NumPage()
	s.MetaData.Name = strings.Trim(md["title"], "\x00")
	for _, v := range strings.Split(strings.Trim(md["keywords"], "\x00"), " ") {
		s.MetaData.Tags = append(s.MetaData.Tags, v)
	}
}

func (s *Session) SetCoverImage(img image.Image) {
	m := resize.Thumbnail(240, 300, img, resize.MitchellNetravali)
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

func (s *Session) StoreAsJpg(idx int, img image.Image) error {
	name := fmt.Sprintf("page%03d.jpg", idx)
	of, err := os.Create(path.Join(s.TempDir, name))
	if err != nil {
		return err
	}
	return jpeg.Encode(of, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
}

func (s *Session) ZipFilesInWorkFolder() error {
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
		s.Log("Writing ZipEntry: %s", fi.Name())
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
	if s.Verbose {
		out, err := json.MarshalIndent(s.MetaData, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Meta:\n%s\n", out)
	}
	return nil
}
