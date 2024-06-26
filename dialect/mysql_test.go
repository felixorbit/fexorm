package dialect

import (
	"reflect"
	"testing"
)

func TestDataTypeOfMySQL(t *testing.T) {
	dial := &mysql{}
	cases := []struct {
		Value interface{}
		Type  string
	}{
		{"Tom", "varchar(255)"},
		{int(123), "int"},
		{float32(1.2), "float"},
		{[]int{1, 2, 3}, "blob"},
	}

	for _, c := range cases {
		if typ := dial.DataTypeOf(reflect.ValueOf(c.Value)); typ != c.Type {
			t.Fatalf("expect %s, but got %s", c.Type, typ)
		}
	}
}
