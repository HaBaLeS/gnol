package main

import (
	"fmt"
	session "github.com/HaBaLeS/gnol/cmd/gnol-tools/tool-session"
	"github.com/teris-io/cli"
	_ "image/png"
	"os"
)

var VersionNum = "undefined"
var BuildDate = "undefined"

const (
	SERIES_ID = "seriesId"
	ORDER_NUM = "orderNum"
)

func main() {
	s := session.NewSession()

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
	seriesId := cli.NewOption(SERIES_ID, "Specify which series the file belongs to").WithType(cli.TypeNumber).WithChar('s')
	orderNum := cli.NewOption(ORDER_NUM, "Specify a ordering number for sorting within a Series").WithType(cli.TypeNumber).WithChar('o')

	//FIXME move keys to constants

	pdf2cbz := cli.NewCommand("pdf2cbz", "PDF to CBZ/CBR converter with support for GNOL Metadata").
		WithArg(inPdfArg).
		WithOption(outFile).
		WithOption(tags).
		WithOption(nsfw).
		WithOption(coverImage).
		WithOption(name).
		WithOption(verbose).
		WithAction(s.Convert)

	folder2cbz := cli.NewCommand("folder2cbz", "Pack folder of images to CBZ with support for GNOL Metadata. Files will be converted to JPEG and Downsized").
		WithArg(inDirArg).
		WithOption(outFile).
		WithOption(tags).
		WithOption(nsfw).
		WithOption(coverImage).
		WithOption(listOrder).
		WithOption(upload).
		WithOption(name).
		WithOption(verbose).
		WithAction(s.Packfolder)

	uploadcmd := cli.NewCommand("upload", "Upload CBZ/CBR to a Gnol instance").
		WithArg(inFile).
		WithOption(nsfw).
		WithOption(seriesId).
		WithOption(orderNum).
		WithAction(s.Upload)

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
		WithAction(s.Repack)

	series := cli.NewCommand("series", "Gnol Series management Command. See subcommands for details").
		WithCommand(cli.NewCommand("list", "list existing series").WithAction(listSeries)).
		WithCommand(cli.NewCommand("create", "create a new series").WithArg(cli.NewArg("name", "Name for Series")).WithAction(createSeries))

	version := cli.NewCommand("version", "Print Version number").WithAction(func(args []string, options map[string]string) int {
		fmt.Printf("gnol-tools %s from %s\n", VersionNum, BuildDate)
		return 0
	})

	monitor := cli.NewCommand("monitor", "Monitor folder and auto-process files in it").WithArg(inDirArg).WithAction(s.Monitor)

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

	app.Run(os.Args, os.Stdout)

}

func listSeries(args []string, options map[string]string) int {
	fmt.Printf("Series: XXXX")
	return 0
}

func createSeries(args []string, options map[string]string) int {
	fmt.Printf("Comics: XXXX")
	return 0
}
