package main

import (
	"flag"
	"text/template"
	"os"
	"path/filepath"
)

type PackageMetaData struct {
	Name string
	SourcePath string
	Version string
	Maintainer string
	MaintainerEmail string
}

const controlTemplate = `Package: {{.Name}}
Version: {{.Version}}-1
Section: base
Priority: optional
Architecture: all
Depends: bash (>= 2.05a-11)
Maintainer: {{.Maintainer}} <{{.MaintainerEmail}}>
Description: Built with debbie
`

var strPackageName = flag.String("name", "", "name of package")
var strPath = flag.String("path", "", "path to sources files")
var strVersion = flag.String("version", "0.0.1", "version of page")
var strMaintainer = flag.String("maintainer", "Dainel Lawrence", "maintainer")
var strMaintainerEmail = flag.String("maintainer-email", "dannyla@linux.com", "maintainer email")


func main() {
	flag.Parse()

	sourcePathAbs, _ := filepath.Abs(*strPath)
	
	
	metadata := PackageMetaData{Name: *strPackageName,
		SourcePath: sourcePathAbs,
		Version: *strVersion,
		Maintainer: *strMaintainer,
		MaintainerEmail: *strMaintainerEmail}

	// control file
	// controlPathAbs := filepath.Join(*strPath, "control")
	// perm0644 := os.FileMode(0644)
	t := template.New("control")
	t, _ = t.Parse(controlTemplate)
	t.Execute(os.Stdout, metadata)

	// changelog

	// write control tarball

	// find md5

	// write md5

	// generate output name

	// convert source into .deb

}
