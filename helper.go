package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

func getUint(buf []byte, offset, nBytes int) (uint64, error) {
	if nBytes > 8 {
		return 0, errors.New("no support for integers with size > 8 bytes")
	}
	b := bytes.NewReader(append(make([]byte, 8-nBytes), buf[offset:offset+nBytes]...))

	var x uint64
	err := binary.Read(b, binary.BigEndian, &x)
	return x, err
}

func getTime(buf []byte, offset int) (time.Time, error) {
	seconds, err := getUint(buf, offset, 4)
	if err != nil {
		return time.Time{}, err
	}
	date := mobiStartDate.Add(time.Duration(seconds) * time.Second)
	return date, nil
}
