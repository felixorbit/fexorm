package clause

import "strings"

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// 多次调用 Set() 设置子句，调用 Build() 构造完整 SQL

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	gSql, gVars := generators[name](vars...)
	c.sql[name] = gSql
	c.sqlVars[name] = gVars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	sqls := make([]string, 0)
	vars := make([]interface{}, 0)
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
