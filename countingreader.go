package main

import "io"

type countingReader struct {
	inner io.Reader
	count *uint64
}

func (cr countingReader) Read(p []byte) (n int, err error) {
	n, err = cr.inner.Read(p)
	*cr.count += uint64(n)
	return
}

func countReader(r io.Reader) countingReader {
	var x uint64
	return countingReader{r, &x}
}
