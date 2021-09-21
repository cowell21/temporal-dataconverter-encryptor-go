package dataconverter

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func compressGZV1(uncompressedData []byte) ([]byte, error) {
	var compressedData bytes.Buffer
	zipWriter := gzip.NewWriter(&compressedData)
	if _, err := zipWriter.Write(uncompressedData); err != nil {
		return nil, err
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}
	return compressedData.Bytes(), nil
}

func decompressGZV1(compressedData []byte) ([]byte, error) {
	if compressedData == nil || len(compressedData) == 0 {
		return []byte(""), nil
	}
	var uncompressedData bytes.Buffer
	reader, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return nil, err
	}
	compressedData, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if _, err = uncompressedData.Write(compressedData); err != nil {
		return nil, err
	}
	return uncompressedData.Bytes(), nil
}

