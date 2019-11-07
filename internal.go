// Package graphdb makes heavy use of `reflect` to construct values out of specified types (an empty, initialized value of a type).
// This makes it easy to read arbitrary values from neo4j client.
// Usage:
// 		- pass empty, initialized type(s) as the last argument(s) of `ReadSingleRow` and `ReadSingleRow`
// 		- get back values and cast them back into the required types
package graphdb

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// convertRecordToTypesFunc returns a transformer function that converts neo4j.Record into the specified types
func convertRecordToTypesFunc(blankTypes []interface{}) func(record neo4j.Record) interface{} {
	return func(record neo4j.Record) interface{} {
		var (
			vals    []interface{}
			rowVals []interface{}
		)
		if record == nil {
			return nil
		}
		vals = record.Values()
		if len(vals) != len(blankTypes) {
			return vals
		}
		for i, val := range vals {
			propWithVal, ok := val.(interface{ Props() map[string]interface{} })
			if ok {
				structVal, err := buildStructWithMap(propWithVal.Props(), blankTypes[i])
				if err != nil {
					return err
				}
				rowVals = append(rowVals, structVal)
				continue
			}
			switch val.(type) {
			case []interface{}:
				newList, err := buildSliceFromList(val.([]interface{}), blankTypes[i])
				if err != nil {
					return err
				}
				rowVals = append(rowVals, newList)
			default:
				rowVals = append(rowVals, val)
			}
		}
		return rowVals
	}
}

// buildSliceFromList reads a list of interface{} and construct a new list based on the specified []<type>
func buildSliceFromList(list []interface{}, blankSlice interface{}) (interface{}, error) {
	if reflect.TypeOf(blankSlice).Kind() != reflect.Slice {
		return nil, errors.New("Received specification slice is not a slice type")
	}
	sliceType := reflect.TypeOf(blankSlice).Elem()
	newSlice := reflect.MakeSlice(reflect.SliceOf(sliceType), len(list), len(list))
	for i, val := range list {
		newSlice.Index(i).Set(reflect.ValueOf(val))
	}
	return newSlice.Interface(), nil
}

// buildStructWithMap reads a map and construct a new struct based on the specified type
func buildStructWithMap(props map[string]interface{}, blankStruct interface{}) (interface{}, error) {
	// create the list container
	itemT := reflect.TypeOf(blankStruct)

	itemPtrVal := reflect.New(itemT)
	if err := scanMapToStruct(props, itemPtrVal.Interface()); err != nil {
		return nil, err
	}
	return itemPtrVal.Elem().Interface(), nil
}

// scanMapToStruct scan a KV value into a struct of specified Type
func scanMapToStruct(props map[string]interface{}, blankStruct interface{}) error {
	reflection := reflect.ValueOf(blankStruct)
	if reflection.Kind() != reflect.Ptr {
		return errors.New("scan requires parameter to be a pointer")
	}

	el := reflection.Elem()
	if el.Kind() != reflect.Struct {
		return errors.New("scan requires parameter to be a struct")
	}

	for i := 0; i < el.NumField(); i++ {
		val := el.FieldByIndex([]int{i})
		valType := val.Kind().String()

		fieldName := el.Type().Field(i).Name
		prop := props[fieldName]

		if val.CanSet() && prop != nil {
			switch valType {
			case "int", "int64":
				val.SetInt(int64(prop.(int)))
			case "bool":
				val.SetBool(prop.(bool))
			case "string":
				val.SetString(prop.(string))
			default:
				return fmt.Errorf("unsupported field type: %s", valType)
			}
		}
	}

	return nil
}
