package main

import (
	"flag"
	"path/filepath"
	"fmt"
	"log"
	"os"
	"path"
	"time"
	//
	"debbie/common"
	"debbie/output"
	"debbie/input"
)

var strPackageName = flag.String("name", "", "name of package")
var strPath = flag.String("path", "", "path to sources files")
var strInstallPath = flag.String("install-path", "", "installPath (eg. /usr/local/<name>)")
var strOutputDir = flag.String("output-dir", "/tmp/", "directory where the .deb file will be written")
var strVersion = flag.String("version", "0.0.1", "version of page")
var strMaintainer = flag.String("maintainer", "Dainel Lawrence", "maintainer")
var strMaintainerEmail = flag.String("maintainer-email", "dannyla@linux.com", "maintainer email")
var strPackageType = flag.String("package-type", "deb", "type package to create (only 'deb' for now)")
var strInputType = flag.String("input-type", "dir", "Input type (only 'dir' for now)")


func main() {
	var ignoreDirs = []string{".bzr", ".hg", ".git"}
	var dataFiles = []common.TarFile{}
	var outputFile string

	flag.Parse()

	if *strPackageName == "" || *strPath == "" {
		fmt.Printf("missing a required info, -name\n")
		fmt.Printf("missing a required info, -path\n")
		os.Exit(1)
	}

	sourcePathAbs, _ := filepath.Abs(*strPath)
	sourcePathBase := path.Base(sourcePathAbs)

	metadata := common.PackageMetaData{
		Name:            *strPackageName,
		SourcePath:      sourcePathAbs,
		SourcePathBase:  sourcePathBase,
		InstallPath:     *strInstallPath,
		OutputDir:       *strOutputDir,
		Version:         *strVersion,
		Maintainer:      *strMaintainer,
		MaintainerEmail: *strMaintainerEmail,
		PackageType:     *strPackageType,
		Time:            time.Now()}

	if *strInputType == "dir" {
		filepath.Walk(metadata.SourcePath, input.PopulateDataFiles(ignoreDirs, &dataFiles, metadata))
	} else {
		fmt.Printf("package-type '%s' is invalid, must be 'dir'\n", *strInputType)
		os.Exit(2)

	}

	if *strPackageType == "deb" {
		outputFile = output.CreateDeb(metadata, dataFiles)
	} else {
		fmt.Printf("package-type '%s' is invalid, must be 'deb'\n", *strPackageType)
		os.Exit(2)
	}

	log.Printf("Created file: %s\n", outputFile)
}

