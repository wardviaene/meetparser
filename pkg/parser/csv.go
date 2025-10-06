package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
)

func MarshalCSV(data interface{}) ([]byte, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("MarshalCSV expects a slice of structs or pointers")
	}

	if v.Len() == 0 {
		return nil, fmt.Errorf("empty slice")
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Get the element type
	elem := v.Index(0)
	elemType := elem.Type()
	if elemType.Kind() == reflect.Ptr { // <-- handle pointer
		elemType = elemType.Elem()
		elem = elem.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("slice elements must be structs or pointers to structs")
	}

	// Write header
	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		headers[i] = elemType.Field(i).Name
	}
	w.Write(headers)

	// Write rows
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		row := make([]string, elem.NumField())
		for j := 0; j < elem.NumField(); j++ {
			row[j] = fmt.Sprint(elem.Field(j).Interface())
		}
		w.Write(row)
	}

	w.Flush()
	return buf.Bytes(), w.Error()
}
