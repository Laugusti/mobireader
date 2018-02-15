package mobireader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

// getUint read a Uint from the raw bytes
func getUint(b []byte, offset, nBytes int) (uint64, error) {
	if nBytes > 8 {
		return 0, errors.New("no support for integers with size > 8 bytes")
	}
	r := bytes.NewReader(append(make([]byte, 8-nBytes), b[offset:offset+nBytes]...))

	var x uint64
	err := binary.Read(r, binary.BigEndian, &x)
	return x, err
}

// getTime returns a time.Time from the raw bytes
func getTime(b []byte, offset int) (time.Time, error) {
	seconds, err := getUint(b, offset, 4)
	if err != nil {
		return time.Time{}, err
	}
	date := mobiStartDate.Add(time.Duration(seconds) * time.Second)
	return date, nil
}
