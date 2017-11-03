package main

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"flag"
	"log"
	"path/filepath"
	"path"
	"os"
	"time"
	"io/ioutil"
	"strings"
	//
	"debbie/common"
       	"debbie/deb"
)


var strPackageName = flag.String("name", "", "name of package")
var strPath = flag.String("path", "", "path to sources files")
var strInstallPath = flag.String("install-path", "", "installPath (eg. /usr/local/<name>)")
var strOutputDir = flag.String("output-dir", "/tmp/", "directory where the .deb file will be written")
var strVersion = flag.String("version", "0.0.1", "version of page")
var strMaintainer = flag.String("maintainer", "Dainel Lawrence", "maintainer")
var strMaintainerEmail = flag.String("maintainer-email", "dannyla@linux.com", "maintainer email")


func main() {
	flag.Parse()

	sourcePathAbs, _ := filepath.Abs(*strPath)
	sourcePathBase := path.Base(sourcePathAbs)

	metadata := common.PackageMetaData{
		Name: *strPackageName,
		SourcePath: sourcePathAbs,
		SourcePathBase: sourcePathBase,
		InstallPath: *strInstallPath,
		OutputDir: *strOutputDir,
		Version: *strVersion,
		Maintainer: *strMaintainer,
		MaintainerEmail: *strMaintainerEmail,
		Time: time.Now()}

	if *strPackageName == "" {
		metadata.Name = sourcePathBase
	}

	var ignoreDirs = []string{".bzr", ".hg", ".git"}
	var dataFiles = []common.TarFiles{}

	filepath.Walk(metadata.SourcePath, populateDataFiles(ignoreDirs, &dataFiles, metadata))

	outputFile := deb.CreateDeb(metadata, dataFiles)
	log.Printf("Created file: %s\n", outputFile)
}


func populateDataFiles(ignoreDirs []string, dataFiles *[]common.TarFiles, metadata common.PackageMetaData) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// sourcePathLength := len(metadata.SourcePath)
		relPath := strings.TrimPrefix(path, metadata.SourcePath)
		targetPath := filepath.Join(metadata.InstallPath, relPath)		
		
		if err != nil {
			log.Print(err)
			return nil
		}
		var fileType = tar.TypeReg
		if info.IsDir() {
			fileType = tar.TypeDir
			dir := filepath.Base(path)
			for _, d := range ignoreDirs {
				if d == dir {
					return filepath.SkipDir
				}
			}
		}
		mode := uint32(info.Mode())
		size := info.Size()
		fileContent, err := ioutil.ReadFile(path)
		fileBuffer := bytes.NewBuffer(fileContent)
		md5Sum := md5.New().Sum(fileContent)

		*dataFiles = append(*dataFiles, common.TarFiles{targetPath, mode, size, byte(fileType), *fileBuffer, md5Sum})
		return nil
    }
}
