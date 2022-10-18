package marshaler

import (
	"fmt"
	"reflect"
)

func Marshal(v any) []byte {
	// prevent primitives
	return marshal(v, "", []byte{})
}

func marshal(v any, parent string, result []byte) []byte {
	// get struct keys
	// if value add to result as key=value
	// else call marshal on value and pass path

	parentValue := reflect.ValueOf(v)
	parentType := reflect.TypeOf(v)
	/* 	fmt.Println("parent:", parent)
	   	fmt.Println("result:", result)
	   	fmt.Println("v:", v)
	   	fmt.Println("------") */
	for i := 0; i < parentValue.NumField(); i++ {
		currentFieldKind := parentValue.Field(i).Kind()
		fmt.Println("PARENTKIND", currentFieldKind)
		if parentType.Kind().String() == "struct" {

			switch currentFieldKind.String() {
			case reflect.Struct.String():
				{ // if is a struct, marshal the struct and append the result
					child := parentValue.Field(i).Interface()
					if parentValue.Field(i).Type().Kind() == reflect.Ptr {
						fmt.Println()
						child = parentValue.Field(i).Type().Elem()
					}

					var x string
					if parent == "" {
						x = parentType.Field(i).Name
					} else {
						x = fmt.Sprintf("%v.%v", parent, parentType.Field(i).Name)
					}

					marshaledStruct := marshal(child, x, []byte{})

					result = []byte(fmt.Sprintf("%v%v ", string(result), string(marshaledStruct)))
				}

			default:
				{
					var res string
					if parent == "" {
						res = fmt.Sprintf("%v=%v ", parentType.Field(i).Name, parentValue.Field(i))
					} else {
						res = fmt.Sprintf("%v.%v=%v ", parent, parentType.Field(i).Name, parentValue.Field(i))
					}

					result = append(result, []byte(res)...)
				}
			}
		}
	}

	// removes last space
	return result[:len(result)-1]
}
