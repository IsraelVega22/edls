package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"runtime"
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

func main() {

	//filter patter
	//flagPatter := flag.String("p", "", "filter by pattern")
	//flaAll := flag.Bool("a", false, "all files including hide files")
	//flagNumberRecords := flag.Int("n", 0, "number of records")

	//order flags
	//hasOrderByTime := flag.Bool("t", false, "sort by time, oldest first")
	//hasOrderBySize := flag.Bool("s", false, "sort by file size,smallest first")
	//hasOrderReverse := flag.Bool("r", false, "reverse order while sorting ")

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
		f, err := getFile(dir, false)
		if err != nil {
			fmt.Println("Error final", err)
		}
		//
		////add file sliceFile
		fs = append(fs, f)
	}

	printList(fs)

	//fmt.Println("Error final", fs)

}

func printList(fs []file) {
	for _, fl := range fs {
		fmt.Printf("%s %s %s %10d %s\n", fl.mode, fl.userName, fl.groupName, fl.size, fl.modificationTime.Format("2006-01-02 15:04:05"))
	}
}

func getFile(dir fs.DirEntry, isHidden bool) (file, error) {
	info, err := dir.Info()
	if err != nil {
		return file{}, err
	}

	f := file{
		name:             dir.Name(),
		isDir:            dir.IsDir(),
		isHidden:         isHidden,
		size:             info.Size(),
		mode:             info.Mode().String(),
		modificationTime: info.ModTime(),
		userName:         "Icenteno",
		groupName:        "TDA",
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
