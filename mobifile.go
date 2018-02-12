package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Laugusti/mobireader/palmdoc"
)

type MOBIFile struct {
	PDBFormat  *PalmDatabaseFormat
	PDHeader   *PalmDocHeader
	MobiHeader *MOBIHeader
	ExthHeader *EXTHHeader
	Records    []*DataRecord
}

type DataRecord struct {
	Data []byte
}

// Create creats a MOBIFile from the Reader interface
func Create(r io.Reader) (*MOBIFile, error) {
	mobi := &MOBIFile{}
	reader := countReader(r)

	// read MobiPocket file format
	err := readFormat(mobi, reader)
	if err != nil {
		return nil, err
	}

	// read the first record (contains headers)
	err = readFirstRecord(mobi, reader)
	if err != nil {
		return nil, err
	}

	// read the rest of Record 0
	endRecord0 := mobi.PDBFormat.RecordInfoEntries[1].Offset
	fmt.Printf("read: %d, eor: %d\n", *reader.count, endRecord0)
	if *reader.count < endRecord0 {
		_, err := io.ReadFull(reader, make([]byte, endRecord0-*reader.count))
		if err != nil {
			return nil, fmt.Errorf("failed the remaining of record 0: %v", err)
		}
	}

	// read the other records
	err = readRemainingRecords(mobi, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read all records: %v", err)
	}
	return mobi, nil
}

// readFormat reads the Palm Database Format structure into the MOBIFile
func readFormat(mobi *MOBIFile, reader io.Reader) error {
	pdbFormat, err := readPalmDatabaseFormat(reader)
	if err != nil {
		return fmt.Errorf("failed to read palm database format: %v", err)
	}
	mobi.PDBFormat = pdbFormat
	return nil
}

// readFirstRecord reads the Headers from Record 0 into the MOBIFile
func readFirstRecord(mobi *MOBIFile, reader io.Reader) error {
	palmDocHeader, err := readPalmDocHeader(reader)
	if err != nil {
		return fmt.Errorf("failed to read palm doc header: %v", err)
	}
	mobi.PDHeader = palmDocHeader

	mobiHeader, err := readMobiHeader(reader)
	if err != nil {
		return fmt.Errorf("failed to read MOBI header: %v", err)
	}
	mobi.MobiHeader = mobiHeader

	exthHeader, err := readExthHeader(reader)
	if err != nil {
		return fmt.Errorf("failed to read EXTH header: %v", err)
	}
	mobi.ExthHeader = exthHeader

	return nil
}

// readRemainingRecords reads all ContentRecords (1 - NumRecords) into the MOBIFile
func readRemainingRecords(mobi *MOBIFile, reader io.Reader) error {
	numRecords := int(mobi.PDBFormat.NumRecords)
	records := make([]*DataRecord, numRecords-1)
	for i := 1; i < 2; i++ {
		var data []byte // raw bytes from record
		if i != (numRecords - 1) {
			startOffset := mobi.PDBFormat.RecordInfoEntries[i].Offset
			endOffset := mobi.PDBFormat.RecordInfoEntries[i+1].Offset
			data = make([]byte, endOffset-startOffset+1)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return err
			}
		} else {
			b, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			data = b
		}
		isContent := (i >= int(mobi.MobiHeader.FirstContentRecordNumber)) &&
			(i <= int(mobi.MobiHeader.LastContentRecordNumber))
		compressionType := int(mobi.PDHeader.Compression)
		encryptionType := int(mobi.PDHeader.EncryptionType)
		dataRecord, err := getDataRecord(data, compressionType, encryptionType, isContent)
		if err != nil {
			return err
		}
		records[i-1] = dataRecord
	}
	mobi.Records = records
	return nil
}

// returns the data record (decompress/unencrypt if necessary)
func getDataRecord(data []byte, compressionType, encryptionType int, isContent bool) (*DataRecord, error) {
	if compressionType != 1 && compressionType != 2 {
		return nil, fmt.Errorf("compression type(%d) is not supported", compressionType)
	}
	if encryptionType != 0 {
		return nil, fmt.Errorf("encryption type(%d) is not supported", encryptionType)
	}
	if !isContent || compressionType == 1 {
		return &DataRecord{data}, nil
	}
	b, err := palmdoc.Decompress(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress record: %v", err)
	}
	return &DataRecord{b}, nil
}
