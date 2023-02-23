package json

import (
	"encoding/json"
	"greatestworks/aop/loader"
	"io/ioutil"
)

type JsonInfo struct {
	Name   string
	Object interface{}
	reader loader.LoadReader
}

func NewJsonInfo(Name string, object interface{}, reader loader.LoadReader) JsonInfo {
	return JsonInfo{
		Name:   Name,
		Object: object,
		reader: reader,
	}
}

var (
	jsonInfos []JsonInfo
)

func GetJsonInfos() []JsonInfo {
	return jsonInfos
}

func ParseJsonFile2Slice(path string, ignoreFileNotExist bool, out interface{}) bool {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		if !ignoreFileNotExist {
			panic(err)
		} else {
			return false
		}
	}
	err = json.Unmarshal(f, out)
	if err != nil {
		panic(err)
	}
	return true
}
