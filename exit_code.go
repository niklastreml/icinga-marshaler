package marshaler

import (
	"errors"
	"log"
	"reflect"
	"strconv"
)

func ExitCode(v any) (int, error) {
	return exitCode(v, OK)
}

func exitCode(v any, code int) (int, error) {
	value := reflect.Indirect(reflect.ValueOf(v))
	log.Printf("%+v\n", reflect.DeepEqual(value, reflect.Value{}))
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		currentField := value.Field(i)
		if fieldIsExported(fieldType) {
			if currentField.Kind() == reflect.Struct {
				return exitCode(currentField.Interface(), code)
			} else {
				warn, err := strconv.ParseFloat(fieldType.Tag.Get("warn"), 64)
				log.Println(warn, err)
				if err != nil {
					return UNKNOWN, err
				}
				crit, err := strconv.ParseFloat(fieldType.Tag.Get("crit"), 64)
				if err != nil {
					return UNKNOWN, err
				}
				min, err := strconv.ParseFloat(fieldType.Tag.Get("min"), 64)
				if err != nil {
					return UNKNOWN, err
				}
				max, err := strconv.ParseFloat(fieldType.Tag.Get("max"), 64)
				if err != nil {
					return UNKNOWN, err
				}
				f, ok := currentField.Interface().(float64)

				if !ok {
					return 0, errors.New("field is not float64")
				}
				if f >= min && f <= max {
					if f < warn && f < crit {
						// everything ok
						if code <= OK {
							code = OK
						}
					} else if f >= warn && f < crit {
						// set to warning
						if code <= WARNING {
							code = WARNING
						}
					} else if f >= crit {
						// set to critical
						if code <= CRITICAL {
							code = CRITICAL
						}
					}
				} else {
					return CRITICAL, nil
				}
			}
		}
	}
	return code, nil

}

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)
