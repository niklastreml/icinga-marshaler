// Holds function for marshaling go structs into icinga readable strings
package marshaler

import (
	"fmt"
	"reflect"
	"strings"
)

// Pass any struct into Marshal to create an icinga readable string.
func Marshal(v any) []byte {
	return []byte(strings.TrimSpace(string(marshal(v, ""))))
}

// Recursively walks through a struct and generates a string from it.
func marshal(v any, parent string) []byte {
	result := []byte{}
	if v == nil {
		return result
	}
	s := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < s.NumField(); i++ {
		fieldType := s.Type().Field(i)
		currentField := s.Field(i)
		if fieldIsExported(fieldType) { // Exported-check must be evaluated first to avoid panic.
			if currentField.Kind() == reflect.Ptr { // case when it's a pointer or struct pointer
				if currentField.IsNil() {
					continue
				}
				result = append(result, marshal(currentField.Elem().Interface(), parent+fieldType.Name+".")...)
			} else if currentField.Kind() == reflect.Struct {
				result = append(result, marshal(currentField.Interface(), parent+fieldType.Name+".")...)
			} else {
				uom := fieldType.Tag.Get("uom")

				keys := []string{"warn", "crit", "min", "max"}

				perfNumbers := ""
				for _, key := range keys {
					value := fieldType.Tag.Get(key)
					perfNumbers = fmt.Sprintf("%v;%v", perfNumbers, value)
				}
				if perfNumbers == ";;;;" {
					perfNumbers = ""
				}
				fieldName := fmt.Sprintf("%v%v", parent, fieldType.Name)
				customName := fieldType.Tag.Get("icinga")
				if customName != "" {
					fieldName = customName
				}

				result = append(result, []byte(fmt.Sprintf("'%v'=%v%v%v ", fieldName, currentField, uom, perfNumbers))...)
			}
		}
	}
	return result
}

func fieldIsExported(field reflect.StructField) bool {
	return field.Name[0] >= 65 && field.Name[0] <= 90
}
