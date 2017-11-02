package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"text/template"
	"path/filepath"
	"path"
	"os"
	"time"
	"io/ioutil"
	//
	"github.com/blakesmith/ar"
)

type PackageMetaData struct {
	Name            string
	DebFile         string
	SourcePath      string
	SourcePathBase  string
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

const controlTemplateText = `
Package: {{.Name}}
Version: {{.Version}}-1
Section: devel
Priority: optional
Architecture: all
Install-Size: {.InstallSize}}
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


func createDeb(metadata PackageMetaData) {
	debFileName := fmt.Sprintf("/tmp/%s_%s_all.deb", metadata.Name, metadata.Version)
	debFile, _ := os.Create(debFileName)
	defer debFile.Close()
	log.Printf("GOAL: %s", debFileName)

	// The debFile is an AR file.
	// arBuffer := new(bytes.Buffer)
	arWriter := ar.NewWriter(debFile)
	err := arWriter.WriteGlobalHeader();
	if err != nil {
		log.Fatalf("Failed to write %s header: %v", debFileName, err)
	}

	// The DataFileTarGz is a TAR file
	tarBuffer := new(bytes.Buffer)
	gzBuffer := gzip.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(gzBuffer)

	// find all data files
	var totalSize int64;
	var ignoreDirs = []string{".bzr", ".hg", ".git"}
	var DataFiles = []TarFiles{}
	var md5sums = new(bytes.Buffer)

	filepath.Walk(metadata.SourcePath, populateDataFiles(ignoreDirs, &DataFiles))

	for _, file := range DataFiles {

		if uint32(file.Type) != uint32(tar.TypeReg) {
			continue
		}

		totalSize += file.Size
		hdr := tar.Header{
			Name:     file.Name,
			Size:     file.Size,
			ModTime:  time.Now(),
			Mode:     0644,
			Typeflag: file.Type,
		}
		md5sums.WriteString(fmt.Sprintf("%s  %s", file.Md5Sum, file.Name))

		err := tarWriter.WriteHeader(&hdr)
		if err != nil {
			log.Fatalf("Writing datafile.tar.gz Header error: %v", err)
		}

		_, err = tarWriter.Write(file.Content.Bytes())
		if err != nil {
			log.Fatalf("Writing datafile.tar.gz Content error: %v", err)
		}
	}

	metadata.InstallSize = totalSize

	err = tarWriter.Close()
	if err != nil {
		log.Fatalf("Failed to close control.tar.gz: %v", err)
	}

	err = gzBuffer.Close()
	if err != nil {
		log.Fatalf("Failed to close control.tar.gz: %v", err)
	}

	controlTarGz, err := createControl(metadata, md5sums.Bytes())
	if err != nil {
		log.Fatalf("CreateControl failed: %v", err)
	}

	// debian-binary
	hdr := ar.Header{
		Name: "debian-binary",
		Size: 4,
		Mode: 0644,
	}

	err = arWriter.WriteHeader(&hdr)
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

	_, err = arWriter.Write([]byte("2.0\n"))
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

	// control.tar.gz
	hdr = ar.Header{
		Name: "control.tar.gz",
		Size: int64(len(controlTarGz)),
		Mode: 0644,
	}

	err = arWriter.WriteHeader(&hdr)
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

	_, err = arWriter.Write(controlTarGz)
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

	// data.tar.gz
	hdr = ar.Header{
		Name: "control.tar.gz",
		Size: int64(len(tarBuffer.Bytes())),
		Mode: 0644,
	}

	err = arWriter.WriteHeader(&hdr)
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

	_, err = arWriter.Write(tarBuffer.Bytes())
	if err != nil {
		log.Fatalf("cannot write file header: %v", err)
	}

}


func createControl(metadata PackageMetaData, md5sums []byte) (controllTarGz []byte, err error) {
	tarBuffer := new(bytes.Buffer)
	gzBuffer := gzip.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(gzBuffer)

	// control file
	controlBuffer := new(bytes.Buffer)
	controlTemplate := template.New("control")
	controlTemplate, err = controlTemplate.Parse(controlTemplateText)
	controlTemplate.Execute(controlBuffer, metadata)

	hdr := tar.Header{
		Name:     "control",
		Size:     int64(controlBuffer.Len()),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	err = tarWriter.WriteHeader(&hdr)
	if err != nil {
		return nil, fmt.Errorf("Failed to write control to control.tar.gz: %v", err)
	}
	_, err = tarWriter.Write(controlBuffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Failed to write control of control.tar.gz: %v", err)
	}

	// md5sums
	hdr = tar.Header{
		Name:     "md5sums",
		Size:     int64(len(md5sums)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	err = tarWriter.WriteHeader(&hdr)
	if err != nil {
		return nil, fmt.Errorf("Failed to write md5sums to control.tar.gz: %v", err)
	}
	_, err = tarWriter.Write(md5sums)
	if err != nil {
		return nil, fmt.Errorf("Failed to write md5sums of control.tar.gz: %v", err)
	}

	err = tarWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to close control.tar.gz: %v", err)
	}

	err = gzBuffer.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to close control.tar.gz: %v", err)
	}

	return tarBuffer.Bytes(), nil
}

func main() {
	flag.Parse()

	sourcePathAbs, _ := filepath.Abs(*strPath)
	sourcePathBase := path.Base(sourcePathAbs)

	metadata := PackageMetaData{
		Name: *strPackageName,
		SourcePath: sourcePathAbs,
		SourcePathBase: sourcePathBase,
		Version: *strVersion,
		Maintainer: *strMaintainer,
		MaintainerEmail: *strMaintainerEmail,
		Time: time.Now()}

	if *strPackageName == "" {
		metadata.Name = sourcePathBase
	}

	createDeb(metadata)




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

		*dataFiles = append(*dataFiles, TarFiles{path, mode, size, byte(fileType), *fileBuffer, md5Sum})
		return nil
    }
}
