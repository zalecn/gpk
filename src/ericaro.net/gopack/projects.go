package gopack

import (
	"encoding/json"
	"os"
)

const (
	GpkFile = ".gpk"
)

//type RemoteFile struct {
//	Name, Url string
//}

func JsonReadFile(path string, v interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func JsonWriteFile(path string, v interface{}) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}