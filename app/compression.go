package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
)

func gzipCompress(body []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(body)
	w.Close()
	return buf.Bytes()
}

func deflate(body []byte) []byte {
	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.BestCompression)
	w.Write(body)
	w.Close()
	return buf.Bytes()
}
