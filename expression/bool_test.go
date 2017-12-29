package builtin

import (
	"reflect"
	"testing"

	_ "qlang.io/lib/builtin"
)

func Test_execbool(t *testing.T) {
	type args struct {
		scriptCode []byte
		resultKey  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{[]byte("a=1==1"), "a"}, true, false},
		{"2", args{[]byte("a=1+1"), "a"}, true, false},
		{"3", args{[]byte("a=1-1"), "a"}, false, false},
		{"4+", args{[]byte("a=(1==2)"), "a"}, false, false},
		{"5", args{[]byte(`a="b"`), "a"}, true, false},
		{"6", args{[]byte(`a="true"`), "a"}, true, false},
		{"7", args{[]byte(`a="on"`), "a"}, true, false},
		{"8", args{[]byte(`a="off"`), "a"}, false, false},
		{"9", args{[]byte(`a="false"`), "a"}, false, false},
		{"10", args{[]byte(`a=[1,2]`), "a"}, true, false},
		{"11", args{[]byte(`a=[]`), "a"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := execbool(tt.args.scriptCode, tt.args.resultKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("execbool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execbool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ResultBool(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"1", args{"1==1"}, true},
		{"2", args{"1+1"}, true},
		{"3", args{"1-1"}, false},
		{"4", args{"1==2"}, false},
		{"5", args{`"b"`}, true},
		{"6", args{`"true"`}, true},
		{"7", args{`"on"`}, true},
		{"8", args{`"off"`}, false},
		{"9", args{`"false"`}, false},
		{"10", args{`[1,2]`}, true},
		{"11", args{`[]`}, false},
		{"12", args{`nil`}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResultBool(tt.args.scriptCode)

			if got != tt.want {
				t.Errorf("ResultBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultString(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{`"false"`}, "false"},
		{"2", args{`"true"`}, "true"},
		{"3", args{`"a"+"b"`}, "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResultString(tt.args.scriptCode)

			if got != tt.want {
				t.Errorf("ResultString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultStringSlice(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{"1", args{`["a","b"]`}, []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResultStringSlice(tt.args.scriptCode)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResultStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
