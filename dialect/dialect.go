package dialect

import "reflect"

// 兼容不同数据库的差异
// 例如数据类型与 Go 语言中类型的关系
type Dialect interface {
	DataTypeOf(reflect.Value) string
	TableExistSQL(string) (string, []interface{})
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectsMap[name]
	return dialect, ok
}
