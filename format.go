package main

import (
	"io"
	"strings"
	"time"
)

type PalmDatabaseFormat struct {
	Name               string
	Attributes         uint64
	Version            uint64
	CreationDate       time.Time
	ModificationDate   time.Time
	LastBackupDate     time.Time
	ModificationNumber uint64
	AppInfoId          uint64
	SortInfoId         uint64
	FormatType         string
	Creator            string
	UniqueIdSeed       uint64
	NextRecordListId   uint64
	NumRecords         uint64
	RecordInfoEntries  []*PDBRecordInfo
	Unknown1           []byte
}

type PDBRecordInfo struct {
	Offset     uint64
	Attributes uint64
	Id         uint64
}

var mobiStartDate = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)

func readPalmDatabaseFormat(r io.Reader) (*PalmDatabaseFormat, error) {
	format := &PalmDatabaseFormat{}

	// read 78 bytes from Reader
	buf := make([]byte, 78)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	// populate fields using byte slice
	format.Name = strings.TrimRight(string(buf[0:32]), string(0))
	format.Attributes, err = getUint(buf, 32, 2)
	if err != nil {
		return nil, err
	}
	format.Version, err = getUint(buf, 34, 2)
	if err != nil {
		return nil, err
	}
	format.CreationDate, err = getTime(buf, 36)
	if err != nil {
		return nil, err
	}
	format.ModificationDate, err = getTime(buf, 40)
	if err != nil {
		return nil, err
	}
	format.LastBackupDate, err = getTime(buf, 44)
	if err != nil {
		return nil, err
	}
	format.ModificationNumber, err = getUint(buf, 48, 4)
	if err != nil {
		return nil, err
	}
	format.AppInfoId, err = getUint(buf, 52, 4)
	if err != nil {
		return nil, err
	}
	format.SortInfoId, err = getUint(buf, 56, 4)
	if err != nil {
		return nil, err
	}
	format.FormatType = string(buf[60:64])
	format.Creator = string(buf[64:68])
	format.UniqueIdSeed, err = getUint(buf, 68, 4)
	if err != nil {
		return nil, err
	}
	format.NextRecordListId, err = getUint(buf, 72, 4)
	if err != nil {
		return nil, err
	}
	format.NumRecords, err = getUint(buf, 76, 2)
	if err != nil {
		return nil, err
	}

	// read record info entries (NumRecord entries of 8 bytes each)
	format.RecordInfoEntries = make([]*PDBRecordInfo, format.NumRecords)
	for i := 0; i < int(format.NumRecords); i++ {
		record := &PDBRecordInfo{}
		buf := make([]byte, 8)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}

		record.Offset, err = getUint(buf, 0, 4)
		if err != nil {
			return nil, err
		}
		record.Attributes, err = getUint(buf, 4, 1)
		if err != nil {
			return nil, err
		}
		record.Id, err = getUint(buf, 5, 3)
		if err != nil {
			return nil, err
		}
		format.RecordInfoEntries[i] = record
	}

	// skip gap to data
	_, err = io.ReadFull(r, buf[0:2])
	if err != nil {
		return nil, err
	}
	format.Unknown1 = buf[0:2]

	return format, nil
}
