package excel

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func LoadExcel(configManager interface{}, dir string) error {
	infos := GetFileInfos()
	if len(infos) == 0 {
		return nil
	}
	var errCh = make(chan error)
	for _, info := range infos {
		info := info
		go func() {
			err := LoadFile(configManager, dir, &info)
			errCh <- err
		}()
	}
	for err := range errCh {
		fmt.Println(err)
	}
	return nil
}

func LoadFile(configManager interface{}, dir string, info *FileInfo) error {
	path := filepath.Join(dir, info.Name)
	openFile, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}
	for _, sheet := range info.Sheets {
		if s, ok := openFile.Sheet[sheet.Name]; ok {
			loadSheet, err := LoadSheet(info, &sheet, s, 4, 1)
			if err != nil {
				return err
			}
			err = sheet.Reader(configManager, loadSheet)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadSheet(info *FileInfo, sheetInfo *SheetInfo,
	sheet *xlsx.Sheet, beginRow, beginClo int) ([]interface{}, error) {

	if len(sheet.Rows) < beginRow || len(sheet.Cols) < beginClo {
		return nil, nil
	}
	objectT := reflect.TypeOf(sheetInfo.Object)
	if objectT.Kind() == reflect.Ptr {
		objectT = objectT.Elem()
	}
	if objectT.Kind() != reflect.Struct {
		return nil, errors.New("object not a struct type")
	}
	cellInfos, err := GetCellInfos(objectT, sheet.Row(0), beginClo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get cell info err:%v", err.Error()))
	}
	var (
		res   = make([]interface{}, 0)
		errFn = func(cloName string, cloIndex, rowIndex int, err error) error {
			return nil
		}
		getCellString = func(cell *xlsx.Cell) (string, error) {

			return cell.String(), nil
		}
	)
	for i, row := range sheet.Rows[beginRow:] {
		if row == nil || len(row.Cells) == 0 {
			break
		}
		elem := reflect.New(objectT).Elem()
		for i2, cellInfo := range cellInfos {
			if cellInfo.CloIndex >= len(row.Cells) {
				continue
			}
			cell := row.Cells[cellInfo.CloIndex]
			cellString, err := getCellString(cell)
			if err != nil {
				return nil, errFn(cellInfo.ColName, cellInfo.CloIndex, i, err)
			}
			if len(cellString) == 0 {
				if i2 == 0 {
					break
				}
				continue
			}
			field := elem.Field(cellInfo.FieldIndex)
			if !field.CanSet() {
				return nil, errFn(cellInfo.ColName, cellInfo.CloIndex, i, errors.New("can not set value "))
			}
			err = ReflectSetVal(field, cellString)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("set value fail:%v", err.Error()))
			}
		}
		res = append(res, elem.Addr().Interface())
	}

	return res, nil
}

type CellInfo struct {
	CloIndex   int
	FieldIndex int
	Field      *reflect.StructField
	Group      string
	ColName    string
}

func GetCellInfos(p reflect.Type, row *xlsx.Row, startClo int) ([]*CellInfo, error) {
	startIndex := startClo - 1
	getCloIndex := func(fieldName string) int {
		for i, cell := range row.Cells {
			if i < startIndex {
				continue
			}
			if cell.Value == fieldName {
				return i
			}
		}
		return -1
	}
	colNames := make([]*CellInfo, 0, len(row.Cells))
	for i := 0; i < p.NumField(); i++ {
		field := p.Field(i)
		fieldName := strings.TrimSpace(field.Tag.Get("col"))
		if fieldName == "" {
			continue
		}
		cloIndex := getCloIndex(fieldName)
		if cloIndex < 0 {
			return nil, errors.New("")
		}
		cellInfo := &CellInfo{
			CloIndex:   cloIndex,
			FieldIndex: i,
			Field:      &field,
			Group:      field.Tag.Get("group"),
			ColName:    row.Cells[cloIndex].Value,
		}
		colNames = append(colNames, cellInfo)
	}
	if len(colNames) == 0 {
		return nil, errors.New("")
	}
	return colNames, nil
}

func ReflectSetVal(value reflect.Value, val string) error {
	trimSpace := strings.TrimSpace(val)
	kind := value.Kind()
	if kind == reflect.Ptr {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}
	switch kind {
	case reflect.Bool:
		parseBool, err := strconv.ParseBool(trimSpace)
		if err != nil {
			return err
		}
		value.SetBool(parseBool)
	case reflect.Int, reflect.Int8, reflect.Uint16:
		parseFloat, err := strconv.ParseFloat(trimSpace, 10)
		if err != nil {
			return err
		}
		if parseFloat < 0 {
			parseFloat -= 0.5
		} else {
			parseFloat += 0.5
		}
		value.SetInt(int64(parseFloat))
	}
	return nil
}
