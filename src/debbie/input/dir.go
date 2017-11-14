package input


import (
	"archive/tar"
	"path/filepath"
	"crypto/md5"
	"strings"
	"log"
	"io/ioutil"
	"bytes"
	"os"

	"debbie/common"

)

func PopulateDataFiles(ignoreDirs []string, dataFiles *[]common.TarFile, metadata common.PackageMetaData) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// sourcePathLength := len(metadata.SourcePath)
		relPath := strings.TrimPrefix(path, metadata.SourcePath)
		targetPath := filepath.Join(metadata.InstallPath, relPath)
		
		if relPath == "" {
			return nil
		}

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

		fileContent, err := ioutil.ReadFile(path)
		fileBuffer := bytes.NewBuffer(fileContent)
		md5Sum := md5.New().Sum(fileContent)

		hdr, err := tar.FileInfoHeader(info, path)
		if err != nil {
			log.Fatalf("failed to FileInfoHeader %s: %v", path, err)
			return nil
		}

		// isSymlink
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				log.Fatalf("failed to read symlink %s: %v", path, err)
				return nil
			}

			hdr, err = tar.FileInfoHeader(info, target)
			if err != nil {
				log.Fatalf("failed to FileInfoHeader %s: %v", path, err)
				return nil
			}

		}
		
		hdr.Name = relPath

		tarFile := common.TarFile{
			Name:      targetPath,
			BasePath:  filepath.Dir(path),
			Path:      path,
			Info:      info,
			Content:   *fileBuffer,
			Md5Sum:    md5Sum,
			TarHeader: hdr,
		}

		*dataFiles = append(*dataFiles, tarFile)
		return nil
	}
}
