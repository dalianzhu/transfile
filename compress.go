package main

import (
	"bytes"
	"compress/gzip"
	"log"
)

// GUnzipData ...
func GUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		log.Printf("gUnzipData error:%v\n", err)
		return
	}

	resData = resB.Bytes()
	return
}

// GZipData ...
func GZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}
	compressedData = b.Bytes()
	return
}
