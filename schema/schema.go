package schema

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/felixorbit/fexorm/dialect"
)

const (
	FexORMTagName       = "fexorm"
	FexORMColumnNameTag = "COLUMN"
)

type Field struct {
	Name, Type, Tag, DBName string
}

// Schema 利用反射完成 结构体和数据库表结构的映射
type Schema struct {
	Model        interface{}
	Name         string   // 结构体名称
	FieldNames   []string // 结构体字段名
	Fields       []*Field
	fieldMap     map[string]*Field
	Table        string   // 数据库表名，取自 TableName() 方法，兜底结构体名
	DBFieldNames []string // 数据库字段名，取自 COLUMN 标签，兜底字段名
}

type NamedTable interface {
	TableName() string
}

func (schema *Schema) GetFieldByName(name string) *Field {
	return schema.fieldMap[name]
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	if v, ok := dest.(NamedTable); ok {
		schema.Table = v.TableName()
	} else {
		schema.Table = schema.Name
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if p.Anonymous || !ast.IsExported(p.Name) { // 匿名字段或者非导出字段不解析
			continue
		}
		field := &Field{
			Name: p.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
		}
		tagSetting := make(map[string]string)
		if v, ok := p.Tag.Lookup(FexORMTagName); ok {
			field.Tag = v
			tagSetting = ParseTagSetting(v)
		}
		if len(tagSetting[FexORMColumnNameTag]) > 0 {
			field.DBName = tagSetting[FexORMColumnNameTag]
		} else {
			field.DBName = field.Name
		}
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, p.Name)
		schema.DBFieldNames = append(schema.DBFieldNames, field.DBName)
		schema.fieldMap[p.Name] = field
	}
	return schema
}

func ParseTagSetting(tag string) map[string]string {
	setting := make(map[string]string)
	tags := strings.Split(tag, ",")
	for _, v := range tags {
		a := strings.Split(v, ":")
		if len(a) != 2 {
			continue
		}
		setting[strings.TrimSpace(strings.ToUpper(a[0]))] = a[1]
	}
	return setting
}
