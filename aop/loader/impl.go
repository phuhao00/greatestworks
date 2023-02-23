package loader

import "reflect"

func SliceLoader(fieldName string) LoadReader {
	return func(confManager interface{}, objects []interface{}) error {
		fieldValue := reflect.ValueOf(confManager).Elem().FieldByName(fieldName)
		_ = fieldValue
		return nil
	}
}
