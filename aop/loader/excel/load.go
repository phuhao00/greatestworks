package excel

import (
	"errors"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

func LoadExcel(configManager interface{}, dir string) error {
	infos := GetFileInfos()
	if len(infos) == 0 {
		return nil
	}
	cpus := runtime.NumCPU()
	wg := sync.WaitGroup{}
	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func LoadFile(configManager interface{}, dir string, info FileInfo) error {
	path := filepath.Join(dir, info.Name)
	openFile, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}
	for i, sheet := range info.Sheets {
		if s, ok := openFile.Sheet[sheet.Name]; ok {
			_ = i
			_ = s

		}
	}
	return nil
}

func LoadSheet(info FileInfo, sheetInfo SheetInfo,
	sheet xlsx.Sheet, beginRow, beginClo int) ([]interface{}, error) {

	if len(sheet.Rows) < beginRow || len(sheet.Cols) < beginClo {
		return nil, nil
	}
	typeOf := reflect.TypeOf(sheetInfo.Object)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("")
	}

	return nil, nil
}
