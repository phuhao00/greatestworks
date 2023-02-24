package excel_object

import "greatestworks/aop/loader/excel"

func init() {
	excel.AppendFileInfo("drop", excel.SheetInfo{
		Name:   "drop",
		Reader: nil,
		Object: Drop{},
	})

}
