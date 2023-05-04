# icinga-marshaler ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/niklastreml/icinga-marshaler) [![GoDoc](https://godoc.org/github.com/NiklasTreml/icinga-marshaler?status.svg)](https://godoc.org/github.com/NiklasTreml/icinga-marshaler) [![codecov](https://codecov.io/gh/NiklasTreml/icinga-marshaler/branch/main/graph/badge.svg?token=BX811094VU)](https://codecov.io/gh/NiklasTreml/icinga-marshaler)

Small library that provides a utility function for easily marshalling any go struct into a format that icinga2 can understand. The library follows the specification mentioned [here](https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#performance-data-metrics).

# Example

```golang
package main

func main() {
    type Check struct {
		BasicValue          string
		FieldWithThresholds int64   `warn:"800" crit:"1024" min:"64" max:"2048"`
		FieldWithCustomName float64 `icinga:"MyCustomField"`
		EverythingTogether  int32   `icinga:"Complex" warn:"800" crit:"1024" min:"64" max:"2048"`
	}

	status := Check{
		BasicValue:          "WARN",
		FieldWithThresholds: 1024,
		FieldWithCustomName: 63.5,
		EverythingTogether:  100,
	}

	bytes := Marshal(status)
	fmt.Println(string(bytes))
}
```

Creates this output:

`'BasicValue'=WARN 'FieldWithThresholds'=1024;800;1024;64;2048 'MyCustomField'=63.5 'Complex'=100;800;1024;64;2048`

## Automatic exit codes
You can get an exit code based on performance data and the Min,Max,Warn and Crit tags. This is useful if you want to exit your program with the correct exit code. You can use the ExitCode() function for this. It will compare the performance data to the thresholds and return the correct exit code. If the performance data is not a number, it will return 3 (UNKNOWN).

```go
package main

import (
	"os"
)

func main() {
    type Check struct {
		MyField int64   `warn:"800" crit:"1024" min:"64" max:"2048"`
	}

	status := Check{
		MyField: 1200,
	}

	code := ExitCode(status)
	os.Exit(code) // exits with code 2, i.e. CRITICAL
}
```