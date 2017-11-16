package output

import (
	"bytes"
	"log"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"os"
	"debbie/common"
)


type Lead struct {
	Magic    []byte // 4 bytes, the magic. hex values ED AB EE DB
	Major    byte   // number, the 'major' version of this rpm file
	Minor    byte   // number, the 'minor' version of this rpm file
	Type     uint16 // short, the 'type' of this rpm
	Arch     uint16 // short, architecture number.
	Name     []byte // string, the package name. Any unused bytes are nulls
	OS       uint16 // short, "os number"
	SigType  uint16 // short, signature type
	Extra    []byte // short, signature type
}


type Header struct {
	Magic  uint64 // begin with the 8-byte header magic value: 8D AD E8 01 00 00 00 00
	Count  uint32 // 4 byte 'tag count' 
	Length uint32 // 4 byte 'tag count'
}


func CreateRpm(metadata common.PackageMetaData, dataFiles []common.TarFile) string {
	var err error
	rpmFileName := fmt.Sprintf("%s_%s_all.rpm", metadata.Name, metadata.Version)
	rpmFilePath := filepath.Join(metadata.OutputDir, rpmFileName)
	rpmFile, _ := os.Create(rpmFilePath)
	defer rpmFile.Close()

	rpmBuffer := new(bytes.Buffer)

	lead := []byte{
		0xED, 0xAB, 0xEE, 0xDB, // Magic
		0x03, 0x00,             // Major.Minor
		0x00, 0x00,             // Type
		0x01,                   // Arch (i386)
		0x72, 0x70, 0x6D, 0x2D, // Name (rpm)
	}
	err = binary.Write(rpmBuffer, binary.BigEndian, lead)
	if err != nil {
		log.Fatalf("binary.Write failed: %v", err)
	}

	// header

	header := Header{
		Magic:  binary.BigEndian.Uint64([]byte{0x8D, 0xAD, 0xE8, 0x01, 0, 0, 0, 0}),
		Count:  binary.BigEndian.Uint32([]byte{0, 0, 0, 0}),
		Length: binary.BigEndian.Uint32([]byte{0, 0, 0, 0}),
	}

	err = binary.Write(rpmBuffer, binary.BigEndian, header)
	if err != nil {
		log.Fatalf("binary.Write failed: %v", err)
	}

	// write the rpmBuffer into the rpmFile
	rpmFile.Write(rpmBuffer.Bytes())
	
	return rpmFilePath
}
