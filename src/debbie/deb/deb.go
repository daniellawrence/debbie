package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"text/template"
	"path/filepath"
	"os"
	"time"
	//
	"debbie/common"
	//
	"github.com/blakesmith/ar"
)


const controlTemplateText = `
Package: {{.Name}}
Version: {{.Version}}-1
Section: devel
Priority: optional
Architecture: all
Install-Size: {{.InstallSize}}
Maintainer: {{.Maintainer}} <{{.MaintainerEmail}}>
Description: Built with debbie, http://github.com/daniellawrence/debbie
`

const changelogTemplateText = `{{.Name}} ({{.Version}}-1) unstable; urgent=medium

  * Initial release

-- {{.Maintainer}} <{{.MaintainerEmail}}> Mon, 22 Mar 2010 00:37:31 +0100
`

const compatTemplateText = `10
`

func CreateDeb(metadata common.PackageMetaData, dataFiles []common.TarFiles) string {
	debFileName := fmt.Sprintf("%s_%s_all.deb", metadata.Name, metadata.Version)
	debFilePath := filepath.Join(metadata.OutputDir, debFileName)
	debFile, _ := os.Create(debFilePath)
	defer debFile.Close()

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
	var md5sums = new(bytes.Buffer)

	for _, file := range dataFiles {

		if uint32(file.Type) != uint32(tar.TypeReg) {
			continue
		}

		totalSize += file.Size
		hdr := tar.Header{
			Name:     file.Name,
			Size:     file.Size,
			ModTime:  time.Now(),
			Mode:     int64(file.Mode),
			Typeflag: file.Type,
		}
		md5sums.WriteString(fmt.Sprintf("%s  %s", file.Md5Sum, file.Name))

		err := tarWriter.WriteHeader(&hdr)
		if err != nil {
			log.Fatalf("Writing '%s' into datafile.tar.gz Header error: %v",
				file.Name, err)
		}

		_, err = tarWriter.Write(file.Content.Bytes())
		if err != nil {
			log.Fatalf("Writing '%s' into datafile.tar.gz Content error: %v",
				file.Name, err)
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
		Name: "data.tar.gz",
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

	return debFilePath

}


func createControl(metadata common.PackageMetaData, md5sums []byte) (controllTarGz []byte, err error) {
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
