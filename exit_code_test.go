package marshaler

import "testing"

func TestExitCodeSimple(t *testing.T) {
	type args struct {
		v Simple
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "OK", args: args{v: Simple{Warning: 5}}, want: OK, wantErr: false},
		{name: "WARNING", args: args{v: Simple{Warning: 15}}, want: WARNING, wantErr: false},
		{name: "CRITICAL", args: args{v: Simple{Warning: 25}}, want: CRITICAL, wantErr: false},
		{name: "Over Max", args: args{v: Simple{Warning: 150}}, want: CRITICAL, wantErr: false},
		{name: "Under Min", args: args{v: Simple{Warning: -50}}, want: CRITICAL, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExitCode(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Simple struct {
	Warning float64 `warn:"10" crit:"20" min:"0" max:"100"`
}

func TestExitCodeComplex(t *testing.T) {
	type args struct {
		v Complex
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "OK", args: args{v: Complex{Warning: 5, Sub: Simple{Warning: 5}}}, want: OK, wantErr: false},
		{name: "WARNING Outer", args: args{v: Complex{Warning: 15, Sub: Simple{Warning: 5}}}, want: WARNING, wantErr: false},
		{name: "CRITICAL Outer", args: args{v: Complex{Warning: 25, Sub: Simple{Warning: 5}}}, want: CRITICAL, wantErr: false},
		{name: "WARNING Inner", args: args{v: Complex{Warning: 5, Sub: Simple{Warning: 15}}}, want: WARNING, wantErr: false},
		{name: "CRITICAL Inner", args: args{v: Complex{Warning: 5, Sub: Simple{Warning: 25}}}, want: CRITICAL, wantErr: false},
		{name: "WARNING Both", args: args{v: Complex{Warning: 15, Sub: Simple{Warning: 15}}}, want: WARNING, wantErr: false},
		{name: "CRITICAL Both", args: args{v: Complex{Warning: 25, Sub: Simple{Warning: 25}}}, want: CRITICAL, wantErr: false},
		{name: "WARNING Outer CRITICAL Inner", args: args{v: Complex{Warning: 25, Sub: Simple{Warning: 15}}}, want: CRITICAL, wantErr: false},
		{name: "CRITICAL Outer WARNING Inner", args: args{v: Complex{Warning: 15, Sub: Simple{Warning: 25}}}, want: CRITICAL, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExitCode(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Complex struct {
	Warning float64 `warn:"10" crit:"20" min:"0" max:"100"`
	Sub     Simple
}

func TestExitCodeFail(t *testing.T) {
	warn := FailWarn{Warning: 5}
	crit := FailCrit{Warning: 5}
	min := FailMin{Warning: 5}
	max := FailMax{Warning: 5}

	t.Run("FailWarn", func(t *testing.T) {
		if _, err := ExitCode(warn); err == nil {
			t.Errorf("Expected err got nil")
		}
	})

	t.Run("FailCrit", func(t *testing.T) {
		if _, err := ExitCode(crit); err == nil {
			t.Errorf("Expected err got nil")
		}
	})

	t.Run("FailMin", func(t *testing.T) {
		if _, err := ExitCode(min); err == nil {
			t.Errorf("Expected err got nil")
		}
	})

	t.Run("FailMax", func(t *testing.T) {
		if _, err := ExitCode(max); err == nil {
			t.Errorf("Expected err got nil")
		}
	})

}

type FailWarn struct {
	Warning float64 `warn:"af" crit:"5" min:"5" max:"5"`
}

type FailCrit struct {
	Warning float64 `warn:"5" crit:"af" min:"5" max:"15"`
}

type FailMin struct {
	Warning float64 `warn:"5" crit:"5" min:"af" max:"5"`
}

type FailMax struct {
	Warning float64 `warn:"5" crit:"5" min:"5" max:"af"`
}

func TestExitCodePointer(t *testing.T) {
	v := 12.12
	pointer := Pointer{Warning: &v}
	nilpointer := Pointer{Warning: nil}
	t.Run("OK", func(t *testing.T) {
		if _, err := ExitCode(pointer); err == nil {
			t.Errorf("Expected err got nil")
		}
	})

	t.Run("WARNING", func(t *testing.T) {
		if _, err := ExitCode(nilpointer); err == nil {
			t.Errorf("Expected err got nil")
		}
	})
}

type Pointer struct {
	Warning *float64 `warn:"10" crit:"20" min:"0" max:"100"`
}
