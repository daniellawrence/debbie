package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"log"
	"text/template"
	"path/filepath"
	"path"
	"os"
	//
	"github.com/blakesmith/ar"
)

type PackageMetaData struct {
	Name string
	SourcePath string
	Version string
	Maintainer string
	MaintainerEmail string
}

type TarFiles struct {
	Name string
	Mode uint32
	Content bytes.Buffer
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

	// control file
	controlBuffer := new(bytes.Buffer)
	controlTemplate := template.New("control")
	controlTemplate, err := controlTemplate.Parse(controlTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = controlTemplate.Execute(controlBuffer, metadata)

	if err != nil {
		log.Fatal(err)
	}

	// changelog
	changelogBuffer := new(bytes.Buffer)
	changelogTemplate := template.New("changelog")
	changelogTemplate, err = changelogTemplate.Parse(changelogTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = changelogTemplate.Execute(changelogBuffer, metadata)

	if err != nil {
		log.Fatal(err)
	}

	//compat
	compatBuffer := new(bytes.Buffer)
	compatTemplate := template.New("compat")
	compatTemplate, err = compatTemplate.Parse(compatTemplateText)

	if err != nil {
		log.Fatal(err)
	}

	err = compatTemplate.Execute(compatBuffer, metadata)

	if err != nil {
		log.Fatal(err)
	}

	// write control tarball
	controlPathAbs := filepath.Join(*strPath, "control.tar.gz")
	controlFile, err := os.Create(controlPathAbs)
	
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)
	controlFiles := [...]TarFiles{
		TarFiles{"debian/control", 0600, *controlBuffer},
		TarFiles{"debian/changelog", 0600, *changelogBuffer},
		TarFiles{"debian/compat", 0600, *compatBuffer},
	}

	for _, file := range controlFiles {
		tarHeader := &tar.Header{
			Name: file.Name,
			Mode: int64(file.Mode),
			Size: int64(file.Content.Len()),
		}
		if err = tarWriter.WriteHeader(tarHeader); err != nil {
			log.Fatal(err)
		}

		ContentBytes := make([]byte, file.Content.Len())
		file.Content.Read(ContentBytes)
		if _, err := tarWriter.Write(ContentBytes); err != nil {
			log.Fatal(err)
		}
	}

	if err = tarWriter.Close(); err != nil {
		log.Fatal(err)
	}

	tarBuffer.WriteTo(controlFile)

	// find all data files
	var ignoreDirs = []string{".bzr", ".hg", ".git"}
	var DataFiles = []TarFiles{}

	arBuffer := new(bytes.Buffer)
	arWriter := ar.NewWriter(arBuffer)
	
	
	filepath.Walk(sourcePathAbs, populateDataFiles(ignoreDirs, &DataFiles))
	
	// write data arball

	// write debian-binary (file)

	// TOOD: ar
	// add 'debian-binary' + control.tar.gz + data.tar.gz via ar as a deb
	// $ ar r package-version.deb debian-binary control.tar.gz ata.tar.gz
	// I really dont want to shell out to running ar...
	// but I dont to rewrite ar in go either.
	// https://github.com/blakesmith/ar might be the way to go here.
}

func populateDataFiles(ignoreDirs []string, dataFiles *[]TarFiles) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if info.IsDir() {
			dir := filepath.Base(path)
			for _, d := range ignoreDirs {
				if d == dir {
					return filepath.SkipDir
				}
			}
		}
		mode := uint32(info.Mode())
		// ContentBytes := make([]byte, file.Content.Len())
		// file.Content.Read(ContentBytes)
		// fileContent = 
		*dataFiles = append(*dataFiles, TarFiles{path, mode, *new(bytes.Buffer)})
		fmt.Println(path)
		return nil
    }
}

