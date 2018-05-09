package expression_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/antlinker/flow/expression"
)

func TestMain(m *testing.M) {
	expression.Import("test", map[string]interface{}{
		"testAdd": func(a, b int) int {
			return a + b
		},
	})
	os.Exit(m.Run())
}

func Test_ResultBool(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{"1==1"}, true, false},
		{"2", args{"1+1"}, true, false},
		{"3", args{"1-1"}, false, false},
		{"4+", args{"(1==2)"}, false, false},
		{"5", args{`"b"`}, true, false},
		{"6", args{`"true"`}, true, false},
		{"7", args{`"on"`}, true, false},
		{"8", args{`"off"`}, false, false},
		{"9", args{`"false"`}, false, false},
		{"10", args{`[1,2]`}, true, false},
		{"11", args{`[]`}, false, false},
		{"12", args{`]`}, false, true},
		{"float1", args{`0.1`}, true, false},
		{"float2", args{`1.1`}, true, false},
		{"byte1", args{`byte(1)`}, true, false},
		{"byte2", args{`byte(0)`}, false, false},
		{"nil", args{`nil`}, false, false},
		{"un", args{`a`}, false, true},
		{"un", args{`1==a`}, false, false},
		{"un", args{`a==b`}, true, false},
		{"global1", args{`global.test_1==1`}, true, false},
		{"global2", args{`global.test_a=="a"`}, true, false},
		{"fun1_1", args{`fun1(global.test_1)`}, true, false},
		{"fun1_2", args{`fun1(2)`}, false, false},
	}
	exp := createTestExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execbool(tt.args.scriptCode)
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

func Test_ResultInt(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, 2, false},
		{"int 2", args{"1*1"}, 1, false},
		{"int 3", args{"4/2"}, 2, false},
		{"int 4", args{"5/2"}, 2, false},
		{"int 5", args{"5"}, 5, false},
		{"var1", args{"a"}, 0, true},
		{"str1", args{`"a"`}, 0, true},
		{"str2", args{`"10"`}, 10, false},
		{"slice1", args{`[1,2]`}, 2, false},
		{"slice2", args{`[]`}, 0, false},
		{"map1", args{`{}`}, 0, false},
		{"map2", args{`{"a":1}`}, 1, false},
		{"map3", args{`{"a":1,"b":"a"}`}, 2, false},
		{"map4", args{`{"a":1,b:"a"}`}, 0, true},
		{"bool1", args{`true`}, 1, false},
		{"bool2", args{`1==2`}, 0, false},
		{"bool3", args{`2>1`}, 1, false},
		{"bool4", args{`2<1`}, 0, false},
		{"float1", args{`1.1`}, 1, false},
		{"float2", args{`0.9`}, 0, false},
		{"float3", args{`-0.9`}, 0, false},
		{"float4", args{`-1.9`}, -1, false},
		{"byte1", args{`byte(1)`}, 1, false},
		{"byte2", args{`byte(2)`}, 2, false},
		{"byte3", args{`byte(0)`}, 0, false},
		{"testAdd", args{`test.testAdd(1,2)`}, 3, false},
		{"testAdd", args{`test.testAdd(3,5)`}, 8, false},
		{"testAdd", args{`test.testAdd(10,20)`}, 30, false},
		{"ctx_10", args{`test.testAdd(ctx_10,20)`}, 30, false},

		{"nil", args{`nil`}, 0, false},
	}
	exp := createTestExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execint(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Resultfloat(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, 2, false},
		{"int 2", args{"1*1"}, 1, false},
		{"int 3", args{"4/2"}, 2, false},
		{"int 4", args{"5/2"}, 2, false},
		{"int 5", args{"5"}, 5, false},
		{"var1", args{"a"}, 0, true},
		{"str1", args{`"a"`}, 0, true},
		{"str2", args{`"10"`}, 10, false},
		{"slice1", args{`[1,2]`}, 0, true},
		{"slice2", args{`[]`}, 0, true},
		{"map1", args{`{}`}, 0, true},
		{"map2", args{`{"a":1}`}, 0, true},
		{"map3", args{`{"a":1,"b":"a"}`}, 0, true},
		{"map4", args{`{"a":1,b:"a"}`}, 0, true},
		{"bool1", args{`true`}, 1, false},
		{"bool2", args{`1==2`}, 0, false},
		{"bool3", args{`2>1`}, 1, false},
		{"bool4", args{`2<1`}, 0, false},
		{"float1", args{`1.1`}, 1.1, false},
		{"float2", args{`0.9`}, 0.9, false},
		{"float3", args{`-0.9`}, -0.9, false},
		{"float4", args{`5.0/2.0`}, 2.5, false},
		{"float5", args{`-1.9`}, -1.9, false},
		{"float6", args{`float(5)/float(2)`}, 2.5, false},
		{"byte1", args{`byte(1)`}, 1, false},
		{"byte2", args{`byte(2)`}, 2, false},
		{"byte3", args{`byte(0)`}, 0, false},

		{"nil", args{`nil`}, 0, false},
	}
	exp := createTestExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execfloat(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execfloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execfloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Resultstring(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, "2", false},
		{"int 2", args{"1*1"}, "1", false},
		{"int 3", args{"4/2"}, "2", false},
		{"int 4", args{"5/2"}, "2", false},
		{"int 5", args{"5"}, "5", false},
		{"var1", args{"a"}, "", true},
		{"str1", args{`"a"`}, "a", false},
		{"str2", args{`"10"`}, "10", false},
		{"slice1", args{`[1,2]`}, "[1 2]", false},
		{"slice2", args{`[]`}, "[]", false},
		{"map1", args{`{}`}, "map[]", false},
		{"map2", args{`{"a":1}`}, `map[a:1]`, false},
		//{"map3", args{`{"a":1,"b":"a"}`}, `map[a:1 b:a]`, false},
		{"map4", args{`{"a":1,b:"a"}`}, "", true},
		{"bool1", args{`true`}, "true", false},
		{"bool2", args{`1==2`}, "false", false},
		{"bool3", args{`2>1`}, "true", false},
		{"bool4", args{`2<1`}, "false", false},
		{"float1", args{`1.1`}, "1.1", false},
		{"float2", args{`0.9`}, "0.9", false},
		{"float3", args{`-0.9`}, "-0.9", false},
		{"float4", args{`5.0/2.0`}, "2.5", false},
		{"float5", args{`-1.9`}, "-1.9", false},
		{"float6", args{`float(5)/float(2)`}, "2.5", false},
		{"byte1", args{`byte(1)`}, "1", false},
		{"byte2", args{`byte(2)`}, "2", false},
		{"byte3", args{`byte(0)`}, "0", false},

		{"nil", args{`nil`}, "", false},
		{"ctx_a", args{"ctx_a"}, "a", false},
	}
	exp := createTestExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execstr(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execstr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execstr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Resultslicestring(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{`["1","2","3"]`}, []string{"1", "2", "3"}, false},
		{"2", args{`["1"]`}, []string{"1"}, false},
		{"3", args{`[""]`}, []string{""}, false},
		{"4", args{`[]`}, nil, true},
		{"5", args{`nil`}, nil, false},
		{"6", args{`a`}, nil, true},
	}
	exp := createTestExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execslicestr(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execslicestr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execslicestr() = %v, want %v", got, tt.want)
			}

		})
	}
}

func createTestExpression() *testExpression {
	exp := expression.CreateExecer("")
	exp.PredefinedJson("global", map[string]interface{}{
		"test_1": 1,
		"test_a": "a",
	})
	exp.PredefinedVar("fun1", `fn(a ) {
		return 1==a
	}`)
	exp.PredefinedVar("fun2", `fn(a ) {
		return 1==a
	}`)
	// exp.ImportAlias("foo/aa", "bb")
	// exp.ImportAlias("foo/cc", "")
	return &testExpression{
		exe: exp,
	}
}

type testExpression struct {
	exe expression.Execer
}

func (t *testExpression) exec(exp string) (*expression.OutData, error) {
	ectx := expression.CreateExpContext(context.Background())

	ectx.AddVar("ctx_10", 10)
	ectx.AddVar("ctx_a", "a")

	return t.exe.Exec(ectx, exp)

}

func (t *testExpression) execbool(exp string) (bool, error) {

	out, err := t.exec(exp)
	if err != nil {
		return false, err
	}
	return out.Bool()
}

func (t *testExpression) execint(exp string) (int, error) {
	out, err := t.exec(exp)
	if err != nil {
		return 0, err
	}

	return out.Int()
}
func (t *testExpression) execfloat(exp string) (float64, error) {
	out, err := t.exec(exp)
	if err != nil {
		return 0, err
	}

	return out.Float()
}
func (t *testExpression) execstr(exp string) (string, error) {
	out, err := t.exec(exp)
	if err != nil {
		return "", err
	}

	return out.String()
}
func (t *testExpression) execslicestr(exp string) ([]string, error) {
	out, err := t.exec(exp)
	if err != nil {
		return nil, err
	}
	return out.SliceStr()
}
