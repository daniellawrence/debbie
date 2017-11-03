package common


import (
	"bytes"
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
}


type TarFiles struct {
	Name    string
	Mode    uint32
	Size    int64
	Type    byte
	Content bytes.Buffer
	Md5Sum  []byte
}
