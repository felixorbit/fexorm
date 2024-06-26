package session

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/felixorbit/fexorm/log"
	"github.com/felixorbit/fexorm/schema"
)

// 实现 Table 相关的 SQL 语句
// 借助 refTable(Schema) 获取表结构信息，拼接 SQL 语句

// Model 解析表结构。为了减少耗时，结构体无变化时不会重复解析
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

// CreateTable 根据 Model 创建表。必须先执行 Model()
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.DBName, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Table, desc)).Exec()
	return err
}

// DropTable 根据 Model 删除表。必须先执行 Model()
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.RefTable().Table)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Table)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	row.Scan(&tmp)
	return tmp == s.RefTable().Table
}
