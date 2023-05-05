package marshaler

import (
	"errors"
	"reflect"
	"strconv"
)

// Calculates the exit code for a given struct
// The struct must have the following tags:
// - warn: warning threshold
// - crit: critical threshold
// - min: minimum value
// - max: maximum value
func ExitCode(v any) (int, error) {
	return exitCode(v, OK)
}

// Recursively walks through the struct and calculates the exit code
// The exit code is the highest exit code of all fields
// If a field is a struct, it will be walked through recursively
// If a field is not a struct, it will be compared to the thresholds
func exitCode(v any, code int) (int, error) {
	value := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		currentField := value.Field(i)
		if fieldIsExported(fieldType) {
			if currentField.Kind() == reflect.Struct {
				c, err := exitCode(currentField.Interface(), code)
				if err != nil {
					return UNKNOWN, err
				}
				if c > code {
					code = c
				}
				continue
			} else {
				warnTag := fieldType.Tag.Get("warn")
				critTag := fieldType.Tag.Get("crit")
				minTag := fieldType.Tag.Get("min")
				maxTag := fieldType.Tag.Get("max")

				f, ok := currentField.Interface().(float64)

				if !ok {
					return 0, errors.New("field is not float64")
				}
				critWarnCode, err := calculateCritWarnCode(f, critTag, warnTag)
				if err != nil {
					return UNKNOWN, err
				}
				minMaxCode, err := calculateMinMaxCode(f, minTag, maxTag)
				if err != nil {
					return UNKNOWN, err
				}

				if critWarnCode > code {
					code = critWarnCode
				} else if minMaxCode > code {
					code = minMaxCode
				}

			}
		}
	}
	return code, nil

}

func calculateCritWarnCode(value float64, critTag, warnTag string) (int, error) {
	if critTag == "" && warnTag == "" {
		return OK, nil
	}
	crit, err := strconv.ParseFloat(critTag, 64)
	if err != nil {
		return UNKNOWN, err
	}
	warn, err := strconv.ParseFloat(warnTag, 64)
	if err != nil {
		return UNKNOWN, err
	}
	if value >= warn && value < crit {
		return WARNING, nil
	} else if value >= crit {
		return CRITICAL, nil
	}
	return OK, nil
}
func calculateMinMaxCode(value float64, minTag, maxTag string) (int, error) {
	if minTag == "" && maxTag == "" {
		return OK, nil
	}
	min, err := strconv.ParseFloat(minTag, 64)
	if err != nil {
		return UNKNOWN, err
	}
	max, err := strconv.ParseFloat(maxTag, 64)
	if err != nil {
		return UNKNOWN, err
	}
	if value >= min && value <= max {
		return OK, nil
	}
	return CRITICAL, nil
}

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)
