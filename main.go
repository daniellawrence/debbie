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

const controlTemplateText = `Package: {{.Name}}
Version: {{.Version}}-1
Section: base
Priority: optional
Architecture: all
Depends: bash (>= 2.05a-11)
Maintainer: {{.Maintainer}} <{{.MaintainerEmail}}>
Description: Built with debbie
`

const changelogTemplateText = `{{.name}} ({{.Version}}-1) unstable; urgent=medium

  * Initial release

-- {{.Maintianer}} <{{.MaintainerEmail}}> Mon, 22 Mar 2010 00:37:31 +0100
`

const compatTemplateText = `10`

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

	// debian directory
	debianPathAbs := filepath.Join(*strPath, "debian")
	os.Mkdir(debianPathAbs, 0644)
		
	// control file
	controlPathAbs := filepath.Join(debianPathAbs, "control")
	println(controlPathAbs)
	controlFile, _ := os.Create(controlPathAbs)
	
	controlTemplate := template.New("control")
	controlTemplate, _ = controlTemplate.Parse(controlTemplateText)
	controlTemplate.Execute(controlFile, metadata)
	
	// changelog
	changelogPathAbs := filepath.Join(*strPath, "debian", "changelog")
	changelogFile, _ := os.Create(changelogPathAbs)

	changelogTemplate := template.New("changelog")
	changelogTemplate, _ = changelogTemplate.Parse(changelogTemplateText)
	changelogTemplate.Execute(changelogFile, metadata)	

	//compant
	compatTemplate := template.New("compat")
	compatTemplate, _ = compatTemplate.Parse(compatTemplateText)
	compatTemplate.Execute(os.Stdout, metadata)	
	
	// write control tarball

	// find md5

	// write md5

	// generate output name

	// convert source into .deb

}
