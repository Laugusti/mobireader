package mobireader

// PalmDOCHeader represents the PalmDOC header of a MOBI file
type PalmDOCHeader struct {
	Compression    uint64
	Length         uint64
	RecordCount    uint64
	RecordSize     uint64
	EncryptionType uint64
	Unknown1       []byte
}

// MOBIHeader represents the MOBI header of a MOBI file
type MOBIHeader struct {
	Identifier               string
	Length                   uint64
	MobiType                 uint64
	Encoding                 uint64
	UniqueID                 uint64
	FileVersion              uint64
	OrthographicIndex        uint64
	InflectionIndex          uint64
	IndexNames               string
	IndexKeys                string
	ExtraIndex0              string
	ExtraIndex1              string
	ExtraIndex2              string
	ExtraIndex3              string
	ExtraIndex4              string
	ExtraIndex5              string
	FirstNonBookIndex        uint64
	FullNameOffset           uint64
	FullNameLength           uint64
	Locale                   uint64
	InputLanguage            string
	OutputLangugage          string
	MinVersion               uint64
	FirstImageIndex          uint64
	HuffmanRecordOffset      uint64
	HuffmanRecordCount       uint64
	HuffmanTableOffset       uint64
	HuffmanTableLength       uint64
	EXTHFlags                uint64
	Unknown1                 []byte
	Unknown2                 []byte
	DRMOffset                uint64
	DRMCount                 uint64
	DRMSize                  uint64
	DRMFlags                 uint64
	Unknown3                 []byte
	FirstContentRecordNumber uint64
	LastContentRecordNumber  uint64
	Unknown4                 []byte
	FCISRecordNumber         uint64
	Unknown5                 []byte
	FLISRecordNumber         uint64
	Unknown6                 []byte
	Unknown7                 []byte
	FirstCompDataSecCount    uint64
	NumCompDataSections      uint64
	Unknown8                 []byte
	ExtraRecordDataFlag      uint64
	INDXRecordOffset         uint64
	Unknown9                 []byte
	Unknown10                []byte
	Unknown11                []byte
	Unknown12                []byte
	Unknown13                []byte
	Unknown14                []byte
}

// EXTHHeader represents the EXTH header of a MOBI file
type EXTHHeader struct {
	Identifier  string
	Length      uint64
	RecordCount uint64
	Records     []*EXTHRecord
}

// EXTHRecord represents an EXTH record
type EXTHRecord struct {
	RecordType uint64
	Length     uint64
	Data       []byte
}

// readPalmDocHeader creates a PalmDoc header from the byte slice
func readPalmDocHeader(b []byte) (*PalmDOCHeader, error) {
	header := &PalmDOCHeader{}

	var err error

	// creates header from the 16 byte slice
	header.Compression, err = getUint(b, 0, 2)
	if err != nil {
		return nil, err
	}

	header.Length, err = getUint(b, 4, 4)
	if err != nil {
		return nil, err
	}

	header.RecordCount, err = getUint(b, 8, 2)
	if err != nil {
		return nil, err
	}

	header.RecordSize, err = getUint(b, 10, 2)
	if err != nil {
		return nil, err
	}

	header.EncryptionType, err = getUint(b, 12, 2)
	if err != nil {
		return nil, err
	}

	header.Unknown1 = b[14:16]

	return header, nil
}

// readMobiHeader creates a MOBI header from the byte slice
func readMobiHeader(b []byte) (*MOBIHeader, error) {
	header := &MOBIHeader{}

	var err error
	header.Identifier = string(b[0:4])
	header.Length, err = getUint(b, 4, 4)
	header.MobiType, err = getUint(b, 8, 4)
	if err != nil {
		return nil, err
	}
	header.Encoding, err = getUint(b, 12, 4)
	if err != nil {
		return nil, err
	}
	header.UniqueID, err = getUint(b, 16, 4)
	if err != nil {
		return nil, err
	}
	header.FileVersion, err = getUint(b, 20, 4)
	if err != nil {
		return nil, err
	}
	header.OrthographicIndex, err = getUint(b, 24, 4)
	if err != nil {
		return nil, err
	}
	header.InflectionIndex, err = getUint(b, 28, 4)
	if err != nil {
		return nil, err
	}
	header.IndexNames = string(b[32:36])
	header.IndexKeys = string(b[36:40])
	header.ExtraIndex0 = string(b[40:44])
	header.ExtraIndex1 = string(b[44:48])
	header.ExtraIndex2 = string(b[48:52])
	header.ExtraIndex3 = string(b[52:56])
	header.ExtraIndex4 = string(b[56:60])
	header.ExtraIndex5 = string(b[60:64])
	header.FirstNonBookIndex, err = getUint(b, 64, 4)
	if err != nil {
		return nil, err
	}
	header.FullNameOffset, err = getUint(b, 68, 4)
	if err != nil {
		return nil, err
	}
	header.FullNameLength, err = getUint(b, 72, 4)
	if err != nil {
		return nil, err
	}
	header.Locale, err = getUint(b, 76, 4)
	if err != nil {
		return nil, err
	}
	header.InputLanguage = string(b[80:84])
	header.OutputLangugage = string(b[84:88])
	header.MinVersion, err = getUint(b, 88, 4)
	if err != nil {
		return nil, err
	}
	header.FirstImageIndex, err = getUint(b, 92, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanRecordOffset, err = getUint(b, 96, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanRecordCount, err = getUint(b, 100, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanTableOffset, err = getUint(b, 104, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanTableLength, err = getUint(b, 108, 4)
	if err != nil {
		return nil, err
	}
	header.EXTHFlags, err = getUint(b, 112, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown1 = b[112:144]
	header.Unknown2 = b[144:148]
	header.DRMOffset, err = getUint(b, 152, 4)
	if err != nil {
		return nil, err
	}
	header.DRMCount, err = getUint(b, 156, 4)
	if err != nil {
		return nil, err
	}
	header.DRMSize, err = getUint(b, 160, 4)
	if err != nil {
		return nil, err
	}
	header.DRMFlags, err = getUint(b, 164, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown3 = b[168:176]
	header.FirstContentRecordNumber, err = getUint(b, 176, 2)
	if err != nil {
		return nil, err
	}
	header.LastContentRecordNumber, err = getUint(b, 178, 2)
	if err != nil {
		return nil, err
	}
	header.Unknown4 = b[180:184]
	header.FCISRecordNumber, err = getUint(b, 184, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown5 = b[188:192]
	header.FLISRecordNumber, err = getUint(b, 192, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown6 = b[196:200]
	header.Unknown7 = b[200:208]
	header.FirstCompDataSecCount, err = getUint(b, 212, 4)
	if err != nil {
		return nil, err
	}
	header.NumCompDataSections, err = getUint(b, 216, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown8 = b[220:224]
	header.ExtraRecordDataFlag, err = getUint(b, 224, 4)
	if err != nil {
		return nil, err
	}
	if header.Length > 228 {
		header.INDXRecordOffset, err = getUint(b, 228, 4)
		if err != nil {
			return nil, err
		}
	}
	if header.Length > 232 {
		header.Unknown9 = b[232:236]
		header.Unknown10 = b[236:240]
		header.Unknown11 = b[240:244]
		header.Unknown12 = b[244:248]
		header.Unknown13 = b[248:252]
		header.Unknown14 = b[252:256]
	}

	return header, nil
}

// hasEXTH returns true if the MOBIHeader indicates there is a EXTH header
func (mh *MOBIHeader) hasEXTH() bool {
	return (mh.EXTHFlags & 0x40) != 0
}

// readExthHeader creates a EXTH header from the byte slice
func readExthHeader(b []byte) (*EXTHHeader, error) {
	header := &EXTHHeader{}

	var err error
	header.Identifier = string(b[0:4])
	header.Length, err = getUint(b, 4, 4)
	if err != nil {
		return nil, err
	}
	header.RecordCount, err = getUint(b, 8, 4)
	if err != nil {
		return nil, err
	}

	// read the EXTH records
	pos := 12
	header.Records = make([]*EXTHRecord, header.RecordCount)
	for i := 0; i < int(header.RecordCount); i++ {
		record := &EXTHRecord{}
		record.RecordType, err = getUint(b, pos, 4)
		if err != nil {
			return nil, err
		}
		record.Length, err = getUint(b, pos+4, 4)
		if err != nil {
			return nil, err
		}

		record.Data = b[pos+8 : pos+int(record.Length)]
		header.Records[i] = record
		pos += int(record.Length)
	}

	return header, nil
}
