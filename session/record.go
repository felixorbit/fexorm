package session

import (
	"errors"
	"reflect"

	"github.com/felixorbit/fexorm/clause"
)

// 实现记录的增删改查功能

// 多次调用 clause.Set() 构造每一个子句
// 调用一次 clause.Build() 按照传入的顺序构造最终语句
// 调用 Raw().Exec() 执行

// Insert like s.Insert(u1, u2, ...)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.CallMethod(BeforeInsert, value)
		s.clause.Set(clause.INSERT, table.Table, table.DBFieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// Find like s.Find(&users);
func (s *Session) Find(values interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(values))
	elementType := sliceValue.Type().Elem()
	table := s.Model(reflect.New(elementType).Elem().Interface()).RefTable()
	s.CallMethod(BeforeQuery, nil)

	s.clause.Set(clause.SELECT, table.Table, table.DBFieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(elementType).Elem()
		var fields []interface{}
		for _, fieldName := range table.FieldNames {
			fields = append(fields, dest.FieldByName(fieldName).Addr().Interface())
		}
		if err = rows.Scan(fields...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		sliceValue.Set(reflect.Append(sliceValue, dest))
	}
	return rows.Close()
}

// Update accepts map[string]interface{} or field-value pairs
func (s *Session) Update(values ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := values[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			m[values[i].(string)] = values[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Table, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	err := s.Limit(1).Find(destSlice.Addr().Interface())
	if err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("record not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Table)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Table)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, values ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), values...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}
