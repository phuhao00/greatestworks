package excel

import "greatestworks/aop/loader"

type SheetInfo struct {
	Name   string
	Reader loader.LoadReader
	Object interface{}
}

func NewSheetInfo(name string, reader loader.LoadReader, object interface{}) SheetInfo {
	return SheetInfo{
		Name:   name,
		Reader: reader,
		Object: object,
	}
}

type FileInfo struct {
	Name   string
	Sheets []SheetInfo
}

var (
	fileInfos []FileInfo
)

func GetFileInfos() []FileInfo {
	return fileInfos
}

func AppendFileInfo(name string, sheets ...SheetInfo) {
	f := FileInfo{Name: name}
	for _, sheet := range sheets {
		f.Sheets = append(f.Sheets, sheet)
	}
	fileInfos = append(fileInfos, f)
}
