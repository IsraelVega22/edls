package main

import (
	"flag"
	"fmt"
	"github.com/AJRDRGZ/fileinfo"
	"github.com/fatih/color"
	"golang.org/x/exp/constraints"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

const Windows = "windows"

type file struct {
	name             string
	fileType         int
	isDir            bool
	isHidden         bool
	userName         string
	groupName        string
	size             int64
	modificationTime time.Time
	mode             string
}

const (
	fileRegular int = iota + 1
	fileDirectory
	fileExecutable
	fileCompress
	fileImage
	fileLink
)

//file extenstion
const (
	exe = ".exe"
	deb = ".deb"
	zip = ".zip"
	gz  = ".gz"
	tar = ".tar"
	rar = ".rar"
	png = ".png"
	jpp = ".jpg"
	gif = ".gif"
)

type stylefile struct {
	icon   string
	color  color.Attribute
	symbol string
}

var mapStyleByfileType = map[int]stylefile{
	fileRegular:    {icon: "documento", color: color.FgMagenta},
	fileDirectory:  {icon: "carpeta", color: color.FgYellow, symbol: "/"},
	fileExecutable: {icon: "Exe", color: color.FgRed, symbol: "*"},
	fileCompress:   {icon: "Compress", color: color.FgMagenta},
	fileImage:      {icon: "Image", color: color.FgCyan},
	fileLink:       {icon: "link", color: color.FgGreen},
}

var (
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
)

func main() {

	//filter patter
	flagPatter := flag.String("p", "", "filter by pattern")
	flaAll := flag.Bool("a", false, "all files including hide files")
	flagNumberRecords := flag.Int("n", 0, "number of records")

	//order flags
	hasOrderByTime := flag.Bool("t", false, "sort by time, oldest first")
	hasOrderBySize := flag.Bool("s", false, "sort by file size,smallest first")
	hasOrderReverse := flag.Bool("r", false, "reverse order while sorting ")

	flag.Parse()

	path := flag.Arg(0)
	fmt.Println(path)
	if path == "" {
		path = "."
	}
	fmt.Println("PAth", path)
	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fs := []file{}

	for _, dir := range dirs {

		isHidden := isHidden(dir.Name(), path)
		//fmt.Println("nombre", dir.Name(), isHidden)

		if isHidden && !*flaAll {
			continue
		}

		if *flagPatter != "" {
			isMatched, err := regexp.MatchString("(?i)"+*flagPatter, dir.Name())
			if err != nil {
				panic(err)
			}

			if !isMatched {
				continue
			}
		}
		f, err := getFile(dir, isHidden)
		if err != nil {
			panic(err)
		}

		//add file sliceFile
		fs = append(fs, f)
	}

	if !*hasOrderBySize || !*hasOrderByTime {
		orderByName(fs, *hasOrderReverse)
	}

	if *hasOrderBySize && !*hasOrderByTime {
		orderBySize(fs, *hasOrderReverse)
	}

	if *hasOrderByTime {
		orderByTime(fs, *hasOrderReverse)
	}

	if *flagNumberRecords == 0 || *flagNumberRecords > len(fs) {
		*flagNumberRecords = len(fs)
	}
	printList(fs, *flagNumberRecords)

}

func orderByTime(files []file, isReverse bool) {
	sort.SliceStable(files, func(i int, j int) bool {
		return mySort(
			files[i].modificationTime.Unix(),
			files[j].modificationTime.Unix(),
			isReverse)
	})
}

func orderBySize(files []file, isReverse bool) {
	sort.SliceStable(files, func(i int, j int) bool {
		return mySort(
			files[i].size,
			files[j].size,
			isReverse)
	})
}

func orderByName(files []file, isReverse bool) {
	sort.SliceStable(files, func(i int, j int) bool {
		return mySort(
			strings.ToLower(files[i].name),
			strings.ToLower(files[j].name),
			isReverse)
	})
}

func mySort[T constraints.Ordered](i, j T, isReverse bool) bool {
	if isReverse {
		return i > j
	}

	return i < j
}

func printList(fs []file, nRecord int) {
	for _, fl := range fs[:nRecord] {
		style, _ := mapStyleByfileType[fl.fileType]
		fmt.Printf("%s %s %s %10d %s %s %s %s %s\n", fl.mode, fl.userName, fl.groupName, fl.size,
			fl.modificationTime.Format("2006-01-02 15:04:05"), style.icon, setColor(fl.name, style.color), fl.name, style.symbol)
	}
}

func getFile(dir fs.DirEntry, isHidden bool) (file, error) {
	var (
		userName, groupName string
	)

	info, err := dir.Info()
	if err != nil {
		return file{}, err
	}

	if userName, groupName = fileinfo.GetUserAndGroup(info.Sys()); userName == "" || groupName == "" {
		userName, groupName = "Icenteno", "TDA"
	}

	f := file{
		name:             dir.Name(),
		isDir:            dir.IsDir(),
		isHidden:         isHidden,
		size:             info.Size(),
		mode:             info.Mode().String(),
		modificationTime: info.ModTime(),
		userName:         userName,
		groupName:        groupName,
	}
	//mt.Println(f)
	setFile(&f)

	return f, nil
}

func setFile(f *file) {
	switch {
	case isLink(*f):
		f.fileType = fileLink
	case isExec(*f):
		f.fileType = fileDirectory
	case isCompress(*f):
		f.fileType = fileCompress
	case isImage(*f):
		f.fileType = fileImage
	default:
		f.fileType = fileRegular
	}

}

func setColor(nameFile string, styleColor color.Attribute) string {
	switch styleColor {
	case color.FgYellow:
		return yellow(nameFile)
	case color.FgGreen:
		return green(nameFile)
	case color.FgRed:
		return red(nameFile)
	case color.FgMagenta:
		return magenta(nameFile)
	case color.FgCyan:
		return cyan(nameFile)
	}
	return nameFile

}

func isLink(f file) bool {
	return strings.HasPrefix(strings.ToUpper(f.name), "L")
}

func isExec(f file) bool {
	if runtime.GOOS == Windows {
		return strings.HasSuffix(f.name, exe)
	}

	return strings.Contains(f.mode, "x")
}

func isCompress(f file) bool {
	return strings.HasSuffix(f.name, zip) ||
		strings.HasSuffix(f.name, gz) ||
		strings.HasSuffix(f.name, tar) ||
		strings.HasSuffix(f.name, rar) ||
		strings.HasSuffix(f.name, deb)
}

func isImage(f file) bool {
	return strings.HasSuffix(f.name, png) ||
		strings.HasSuffix(f.name, jpp) ||
		strings.HasSuffix(f.name, gif)
}

func isHidden(fileName string, basePath string) bool {
	filePath := fileName

	if runtime.GOOS == Windows {
		filePath = path.Join(basePath, fileName)
	}

	return fileinfo.IsHidden(filePath)
}
