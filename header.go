package main

import (
	"io"
)

type palmDocHeader struct {
	Compression    uint64
	Length         uint64
	RecordCount    uint64
	RecordSize     uint64
	EncryptionType uint64
}

type mobiHeader struct {
	Identifier               string
	Length                   uint64
	MobiType                 uint64
	Encoding                 uint64
	UniqueId                 uint64
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

type exthHeader struct {
	Identifier  string
	Length      uint64
	RecordCount uint64
	Records     []exthRecord
}

type exthRecord struct {
	RecordType uint64
	Length     uint64
	Data       []byte
}

func readPalmDocHeader(r io.Reader) (*palmDocHeader, error) {
	header := &palmDocHeader{}
	buf := make([]byte, 16)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	header.Compression, err = getUint(buf, 0, 2)
	if err != nil {
		return nil, err
	}

	header.Length, err = getUint(buf, 4, 4)
	if err != nil {
		return nil, err
	}

	header.RecordCount, err = getUint(buf, 8, 2)
	if err != nil {
		return nil, err
	}

	header.RecordSize, err = getUint(buf, 10, 2)
	if err != nil {
		return nil, err
	}

	header.EncryptionType, err = getUint(buf, 12, 2)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func readMobiHeader(r io.Reader) (*mobiHeader, error) {
	header := &mobiHeader{}
	// read MOBI identifier and length
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	header.Identifier = string(buf[0:4])
	header.Length, err = getUint(buf, 4, 4)
	if err != nil {
		return nil, err
	}

	// read the rest of the MOBI header
	rest := make([]byte, header.Length-8)
	_, err = io.ReadFull(r, rest)
	if err != nil {
		return nil, err
	}
	buf = append(buf, rest...)

	header.MobiType, err = getUint(buf, 8, 4)
	if err != nil {
		return nil, err
	}
	header.Encoding, err = getUint(buf, 12, 4)
	if err != nil {
		return nil, err
	}
	header.UniqueId, err = getUint(buf, 16, 4)
	if err != nil {
		return nil, err
	}
	header.FileVersion, err = getUint(buf, 20, 4)
	if err != nil {
		return nil, err
	}
	header.OrthographicIndex, err = getUint(buf, 24, 4)
	if err != nil {
		return nil, err
	}
	header.InflectionIndex, err = getUint(buf, 28, 4)
	if err != nil {
		return nil, err
	}
	header.IndexNames = string(buf[32:36])
	header.IndexKeys = string(buf[36:40])
	header.ExtraIndex0 = string(buf[40:44])
	header.ExtraIndex1 = string(buf[44:48])
	header.ExtraIndex2 = string(buf[48:52])
	header.ExtraIndex3 = string(buf[52:56])
	header.ExtraIndex4 = string(buf[56:60])
	header.ExtraIndex5 = string(buf[60:64])
	header.FirstNonBookIndex, err = getUint(buf, 64, 4)
	if err != nil {
		return nil, err
	}
	header.FullNameOffset, err = getUint(buf, 68, 4)
	if err != nil {
		return nil, err
	}
	header.FullNameLength, err = getUint(buf, 72, 4)
	if err != nil {
		return nil, err
	}
	header.Locale, err = getUint(buf, 76, 4)
	if err != nil {
		return nil, err
	}
	header.InputLanguage = string(buf[80:84])
	header.OutputLangugage = string(buf[84:88])
	header.MinVersion, err = getUint(buf, 88, 4)
	if err != nil {
		return nil, err
	}
	header.FirstImageIndex, err = getUint(buf, 92, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanRecordOffset, err = getUint(buf, 96, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanRecordCount, err = getUint(buf, 100, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanTableOffset, err = getUint(buf, 104, 4)
	if err != nil {
		return nil, err
	}
	header.HuffmanTableLength, err = getUint(buf, 108, 4)
	if err != nil {
		return nil, err
	}
	header.EXTHFlags, err = getUint(buf, 112, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown1 = buf[112:144]
	header.Unknown2 = buf[144:148]
	header.DRMOffset, err = getUint(buf, 152, 4)
	if err != nil {
		return nil, err
	}
	header.DRMCount, err = getUint(buf, 156, 4)
	if err != nil {
		return nil, err
	}
	header.DRMSize, err = getUint(buf, 160, 4)
	if err != nil {
		return nil, err
	}
	header.DRMFlags, err = getUint(buf, 164, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown3 = buf[168:176]
	header.FirstContentRecordNumber, err = getUint(buf, 176, 2)
	if err != nil {
		return nil, err
	}
	header.LastContentRecordNumber, err = getUint(buf, 178, 2)
	if err != nil {
		return nil, err
	}
	header.Unknown4 = buf[180:184]
	header.FCISRecordNumber, err = getUint(buf, 184, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown5 = buf[188:192]
	header.FLISRecordNumber, err = getUint(buf, 192, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown6 = buf[196:200]
	header.Unknown7 = buf[200:208]
	header.FirstCompDataSecCount, err = getUint(buf, 212, 4)
	if err != nil {
		return nil, err
	}
	header.NumCompDataSections, err = getUint(buf, 216, 4)
	if err != nil {
		return nil, err
	}
	header.Unknown8 = buf[220:224]
	header.ExtraRecordDataFlag, err = getUint(buf, 224, 4)
	if err != nil {
		return nil, err
	}
	if header.Length > 228 {
		header.INDXRecordOffset, err = getUint(buf, 228, 4)
		if err != nil {
			return nil, err
		}
	}
	if header.Length > 232 {
		header.Unknown9 = buf[232:236]
		header.Unknown10 = buf[236:240]
		header.Unknown11 = buf[240:244]
		header.Unknown12 = buf[244:248]
		header.Unknown13 = buf[248:252]
		header.Unknown14 = buf[252:256]
	}

	return header, nil
}

func readExthHeader(r io.Reader) (*exthHeader, error) {
	header := &exthHeader{}

	// read EXTH identifier, length, and count
	buf := make([]byte, 12)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	header.Identifier = string(buf[0:4])
	header.Length, err = getUint(buf, 4, 4)
	if err != nil {
		return nil, err
	}
	header.RecordCount, err = getUint(buf, 8, 4)
	if err != nil {
		return nil, err
	}

	// read the EXTH records
	header.Records = make([]exthRecord, header.RecordCount)
	for i := 0; i < int(header.RecordCount); i++ {
		record := &(header.Records[i])
		*record = exthRecord{}
		_, err = io.ReadFull(r, buf[:8])
		if err != nil {
			return nil, err
		}
		record.RecordType, err = getUint(buf, 0, 4)
		if err != nil {
			return nil, err
		}
		record.Length, err = getUint(buf, 4, 4)
		if err != nil {
			return nil, err
		}
		data := make([]byte, record.Length-8)
		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}
		record.Data = data
	}

	// skip EXTH padding (EXTH header length is padded to multiple of four)
	padBytes := 4 - (header.Length % 4)
	_, err = io.ReadFull(r, buf[:padBytes])
	if err != nil {
		return nil, err
	}

	return header, nil
}
