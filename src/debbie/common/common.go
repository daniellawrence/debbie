package common


import (
	"archive/tar"
	"bytes"
	"os"
	"time"
)


type PackageMetaData struct {
	Name            string
	DebFile         string
	SourcePath      string
	SourcePathBase  string
	InstallPath     string
	OutputDir       string
	Version         string
	Maintainer      string
	MaintainerEmail string
	InstallSize     int64
	Time            time.Time
	PackageType     string
}


type TarFile struct {
	Name      string
	BasePath  string
	Path      string
	Content   bytes.Buffer
	Md5Sum    []byte
	Info      os.FileInfo
	TarHeader *tar.Header
}


func (tf TarFile) TarType() byte {
	switch mode := tf.Info.Mode(); {                   
	case mode.IsRegular():
		return tar.TypeReg
	case mode.IsDir():
		return tar.TypeDir
	case mode&os.ModeSymlink != 0:
		return tar.TypeSymlink
	default:
		return 1
	}

}
