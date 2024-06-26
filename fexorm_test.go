package fexorm

import (
	"errors"
	"github.com/felixorbit/fexorm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "fex.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	ns := engine.NewSession()
	_ = ns.Model(&User{}).DropTable()
	_ = ns.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{"Tom", 18})
		_, err = s.Insert(&User{"Bob", 18})
		return nil, errors.New("Error")
	})
	if err == nil {
		t.Fatal("failed to rollback")
	}
	count, _ := ns.Model(&User{}).Count()
	if count != 0 {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	ns := engine.NewSession()
	_ = ns.Model(&User{}).DropTable()
	_ = ns.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{"Tom", 18})
		_, err = s.Insert(&User{"Bob", 18})
		return
	})
	count, _ := ns.Count()
	if err != nil || count != 2 {
		t.Fatal("failed to commit")
	}
}
