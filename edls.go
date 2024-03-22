package main

import "time"

//file types
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

type stylefile struct {
	icon   string
	color  string
	symbol string
}

var mapStyleByfileType = map[int]stylefile{
	fileRegular:    {icon: "documento"},
	fileDirectory:  {icon: "carpeta", color: "blue", symbol: "/"},
	fileExecutable: {icon: "Exe", color: "GREEN", symbol: "*"},
	fileCompress:   {icon: "Compres", color: "RED"},
	fileImage:      {icon: "Image", color: "MAGENTA"},
	fileLink:       {icon: "link", color: "CYAN"},
}
