package rpm

// name: name-version-relase.architecture.rpm
// eg:   deb-0.0.1-1.i386.rpm
// ref: http://ftp.rpm.org/max-rpm/ch-rpm-file-format.html

// lead
// headers
// payload
// ref: http://www.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html

type rpmHeader struct {
	Magic  uint64
	Count  uint32
	Length uint32
}


type rpmLead struct {
}


func CreateRPMLead()  rpmLead {
	var magicMajorMinorTypeArch = []byte{0xED, 0xAB, 0xEE, 0xDB, 3, 0, 0x00, 0x00}
	var packageName := make([]byte, 66)
	lead := rpmLead{}
	return lead
}


func CreateRPMHeader() rpmHeader {
	var data = []byte{0x8D, 0xAD, 0xE8, 0x01, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 4, 132}
	buf := bytes.NewBuffer(data)
	header := rpmHeader{}
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		log.Fatalf("binary.Read failed: %v", err)
	}
	return header
}
