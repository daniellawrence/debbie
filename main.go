package main

import (
	"flag"
	"log"
	"text/template"
	"os"
	"path/filepath"
	"path"
)

type PackageMetaData struct {
	Name string
	SourcePath string
	Version string
	Maintainer string
	MaintainerEmail string
}

const controlTemplateText = `Package: {{.Name}}
Version: {{.Version}}-1
Section: base
Priority: optional
Architecture: all
Depends: bash (>= 2.05a-11)
Maintainer: {{.Maintainer}} <{{.MaintainerEmail}}>
Description: Built with debbie
`

const changelogTemplateText = `{{.Name}} ({{.Version}}-1) unstable; urgent=medium

  * Initial release

-- {{.Maintainer}} <{{.MaintainerEmail}}> Mon, 22 Mar 2010 00:37:31 +0100
`

const compatTemplateText = `10
`

var strPackageName = flag.String("name", "", "name of package")
var strPath = flag.String("path", "", "path to sources files")
var strVersion = flag.String("version", "0.0.1", "version of page")
var strMaintainer = flag.String("maintainer", "Dainel Lawrence", "maintainer")
var strMaintainerEmail = flag.String("maintainer-email", "dannyla@linux.com", "maintainer email")


func main() {
	flag.Parse()

	sourcePathAbs, _ := filepath.Abs(*strPath)
	sourcePathBase := path.Base(sourcePathAbs)

	metadata := PackageMetaData{
		Name: *strPackageName,
		SourcePath: sourcePathAbs,
		Version: *strVersion,
		Maintainer: *strMaintainer,
		MaintainerEmail: *strMaintainerEmail}

	
	if *strPackageName == "" {
		metadata.Name = sourcePathBase
	}

	// debian directory
	debianPathAbs := filepath.Join(*strPath, "debian")
	os.Mkdir(debianPathAbs, 0755)

	// control file
	controlPathAbs := filepath.Join(debianPathAbs, "control")
	controlFile, err := os.Create(controlPathAbs)

	if err != nil {
		log.Fatal(err)
	}

	controlTemplate := template.New("control")
	controlTemplate, err = controlTemplate.Parse(controlTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = controlTemplate.Execute(controlFile, metadata)

	if err != nil {
		log.Fatal(err)
	}

	// changelog
	changelogPathAbs := filepath.Join(debianPathAbs, "changelog")
	changelogFile, err := os.Create(changelogPathAbs)

	if err != nil {
		log.Fatal(err)
	}

	changelogTemplate := template.New("changelog")
	changelogTemplate, err = changelogTemplate.Parse(changelogTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = changelogTemplate.Execute(changelogFile, metadata)

	if err != nil {
		log.Fatal(err)
	}

	//compat
	compatPathAbs := filepath.Join(debianPathAbs, "compat")
	compatFile, err := os.Create(compatPathAbs)

	if err != nil {
		log.Fatal(err)
	}

	compatTemplate := template.New("compat")
	compatTemplate, err = compatTemplate.Parse(compatTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = compatTemplate.Execute(compatFile, metadata)

	if err != nil {
		log.Fatal(err)
	}

	// write control tarball

	// find md5

	// write md5

	// write data tarball

	// generate output name

	// convert source into .deb

}
