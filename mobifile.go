// Package mobireader reads a MOBI File into a struct from a
// io.Reader
package mobireader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Laugusti/mobireader/palmdoc"
)

// MOBIFile contains the format, headers, and data records of a MOBI file
type MOBIFile struct {
	PDBFormat  *PalmDatabaseFormat
	PDHeader   *PalmDOCHeader
	MobiHeader *MOBIHeader
	ExthHeader *EXTHHeader
	Records    []*DataRecord
}

// DataRecord represents the bytes of a Record in the MOBI file
type DataRecord struct {
	Data []byte
}

// Create creats a MOBIFile from the Reader interface
func Create(reader io.Reader) (*MOBIFile, error) {
	mobi := &MOBIFile{}

	// read MobiPocket file format
	err := readFormat(mobi, reader)
	if err != nil {
		return nil, err
	}

	// read all records
	err = readRawRecords(mobi, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read all records: %v", err)
	}

	// set Headers on the MOBI file using the first record
	err = addHeaders(mobi)
	if err != nil {
		return nil, err
	}

	// decompress/unencrpt text Records
	err = decompressTextRecords(mobi)
	if err != nil {
		return nil, err
	}

	return mobi, nil
}

// Text returns all TextRecords as a single string
func (m *MOBIFile) Text() (string, error) {
	var buf bytes.Buffer
	for i := m.MobiHeader.FirstContentRecordNumber; i < m.MobiHeader.FirstImageIndex; i++ {
		_, err := buf.Write(m.Records[i].Data)
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
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
func addHeaders(mobi *MOBIFile) error {
	// PalmDOC Header is the first 16 bytes
	palmDocHeader, err := readPalmDocHeader(mobi.Records[0].Data[:16])
	if err != nil {
		return fmt.Errorf("failed to read palm doc header: %v", err)
	}
	mobi.PDHeader = palmDocHeader

	// MOBI Header follows PalmDOC header
	mobiHeader, err := readMobiHeader(mobi.Records[0].Data[16:])
	if err != nil {
		return fmt.Errorf("failed to read MOBI header: %v", err)
	}
	mobi.MobiHeader = mobiHeader

	// EXTH Header follows the MOBI header
	if mobiHeader.hasEXTH() {
		exthHeader, err := readExthHeader(mobi.Records[0].Data[16+mobiHeader.Length:])
		if err != nil {
			return fmt.Errorf("failed to read EXTH header: %v", err)
		}
		mobi.ExthHeader = exthHeader
	}

	return nil
}

// readRawRecords reads all records into the MOBIFile,
func readRawRecords(mobi *MOBIFile, reader io.Reader) error {
	// read raw bytes from reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	numRecords := int(mobi.PDBFormat.NumRecords)
	// format was already read from the reader, there Record Offsets
	// needs to be offset
	offset := 0
	if numRecords > 0 {
		offset = int(mobi.PDBFormat.RecordInfoEntries[0].Offset)
	}
	// create DataRecords from raw bytes
	mobi.Records = make([]*DataRecord, numRecords)
	for i := 0; i < numRecords; i++ {
		var start, end int
		start = int(mobi.PDBFormat.RecordInfoEntries[i].Offset) - offset
		if i != (numRecords - 1) {
			end = int(mobi.PDBFormat.RecordInfoEntries[i+1].Offset) - offset
		} else {
			end = len(data)
		}
		mobi.Records[i] = &DataRecord{data[start:end]}
	}
	return nil
}

// decompressTextRecords decompress all the text records in the MOBI file
func decompressTextRecords(mobi *MOBIFile) error {
	compressionType := int(mobi.PDHeader.Compression)
	encryptionType := int(mobi.PDHeader.EncryptionType)

	// loop through text records and decompress
	for i := mobi.MobiHeader.FirstContentRecordNumber; i < mobi.MobiHeader.FirstImageIndex; i++ {

		err := decompressRecord(mobi.Records[i], compressionType, encryptionType)
		if err != nil {
			return err
		}
	}
	return nil
}

// deCompressRecord decompresses/unencrypts the DataRecord
func decompressRecord(dataRecord *DataRecord, compressionType, encryptionType int) error {
	// compression type must be either non or palmdoc
	if compressionType != 1 && compressionType != 2 {
		return fmt.Errorf("compression type(%d) is not supported", compressionType)
	}

	// no encryption is supported
	if encryptionType != 0 {
		return fmt.Errorf("encryption type(%d) is not supported", encryptionType)
	}

	// decompress if compression type is palmdoc
	if compressionType == 2 {
		b, err := palmdoc.Decompress(dataRecord.Data)
		if err != nil {
			return fmt.Errorf("failed to decompress record: %v", err)
		}
		dataRecord.Data = b
	}
	return nil
}
