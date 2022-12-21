package main

import (
	"flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

var (
	excelPath string
	jsonPath  string
)

func init() {
	flag.StringVar(&excelPath, "excelPath", "../../excel", "please")
	flag.StringVar(&jsonPath, "jsonPath", "../../../aop/json", "please")
}

func main() {
	flag.Parse()
	filelist := getFileList(excelPath)
	fmt.Println(filelist)
	for _, file := range filelist {
		parseFile(file)
	}
}

func getFileList(path string) []string {
	var all_file []string
	finfo, _ := os.ReadDir(path)
	for _, info := range finfo {
		if filepath.Ext(info.Name()) == ".xlsx" {
			real_path := path + "/" + info.Name()
			if info.IsDir() {
				//all_file = append(all_file, getFileList(real_path)...)
			} else {
				all_file = append(all_file, real_path)
			}
		}
	}

	return all_file
}

type meta struct {
	Key string
	Idx int
	Typ string
}

type rowdata []interface{}

func parseFile(file string) {

	fmt.Println("\n\n\n\n", file)

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		panic(err.Error())
	}
	//[line][colidx][data]

	sheets := xlsx.GetSheetList()
	for _, s := range sheets {
		rows, err := xlsx.GetRows(s)
		if err != nil {
			return
		}
		if len(rows) < 5 {
			return
		}

		colNum := len(rows[1])
		fmt.Println("col num:", colNum)
		metaList := make([]*meta, 0, colNum)
		dataList := make([]rowdata, 0, len(rows)-4)
		sheetName := "template"
		for line, row := range rows {
			switch line {
			case 0: // sheet 名
				sheetName = row[0]
			case 1: // col name

				for idx, colname := range row {
					fmt.Println(idx, colname, len(metaList))

					metaList = append(metaList, &meta{Key: colname, Idx: idx})
				}
			case 2: // data type

				fmt.Println("meta cot:%d, rol cot:%d", len(metaList), len(row))
				for idx, typ := range row {
					metaList[idx].Typ = typ
				}
			case 3: // desc

			default: //>= 4 row data
				data := make(rowdata, colNum)

				for k := 0; k < colNum; k++ {
					if k < len(row) {
						data[k] = row[k]
					}
				}

				dataList = append(dataList, data)
			}
		}

		//sheetName := xlsx.GetSheetName(idx)
		// to json, save
		//filename := filepath.Base(file)
		//suf := filepath.Ext(filename)
		jsonFile := fmt.Sprintf("%s.json", sheetName)
		err = output(jsonFile, toJson(dataList, metaList))
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(toJson(dataList, metaList))

	}

}

func toJson(datarows []rowdata, metalist []*meta) string {
	ret := "["

	for _, row := range datarows {
		ret += "\n\t{"
		for idx, meta := range metalist {
			ret += fmt.Sprintf("\n\t\t\"%s\":", meta.Key)
			if meta.Typ == "string" {
				if row[idx] == nil {
					ret += "\"\""
				} else {
					ret += fmt.Sprintf("\"%s\"", row[idx])
				}
			} else {
				if row[idx] == nil || row[idx] == "" {
					ret += "0"
				} else {
					ret += fmt.Sprintf("%s", row[idx])
				}
			}
			ret += ","
		}
		ret = ret[:len(ret)-1]

		ret += "\n\t},"
	}
	ret = ret[:len(ret)-1]

	ret += "\n]"
	return ret
}

func output(filename string, str string) error {

	f, err := os.OpenFile(jsonPath+"\\"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		return err
	}

	return nil
}
