package schema

import (
	"testing"

	"github.com/felixorbit/fexorm/dialect"
)

type User struct {
	Name string `fexorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User schema")
	}
	if schema.GetFieldByName("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
