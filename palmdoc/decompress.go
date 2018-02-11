// Package palmdoc provides a method to decompress a byte slice that
// was compressed using the PalmDoc compression algorithm
package palmdoc

import (
	"bytes"
)

func Decompress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	for i := 0; i < len(b); {
		c := b[i]
		i++

		switch {
		case c >= 0x1 && c <= 0x8:
			for j := 0; j < int(c); j++ {
				err := buf.WriteByte(b[i])
				i++
				if err != nil {
					return nil, err
				}
			}
		case c <= 0x7f: // 0x0, 0x9-0x7f
			err := buf.WriteByte(c)
			if err != nil {
				return nil, err
			}
		case c >= 0xc0:
			err := buf.WriteByte(' ')
			if err != nil {
				return nil, err
			}
			err = buf.WriteByte(c ^ 0x80)
			if err != nil {
				return nil, err
			}
		default:
			c := (int(c) << 8) + int(b[i])
			i++

			di := (c & 0x3fff) >> 3
			for n := (c & 7) + 3; n > 0; n-- {
				err := buf.WriteByte(buf.Bytes()[buf.Len()-di])
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return buf.Bytes(), nil
}
