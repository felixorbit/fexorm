package main

import (
	"fmt"

	"github.com/felixorbit/fexorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//testSQLite()
	testMySQL()
}

func testSQLite() {
	engine, _ := fexorm.NewEngine("sqlite3", "fex.db")
	defer engine.Close()

	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected rows\n", count)
}

type User struct {
	Id   int    `fexorm:"COLUMN:id"`
	Name string `fexorm:"COLUMN:name"`
	Age  int    `fexorm:"COLUMN:age"`
}

func (User) TableName() string {
	return "user"
}

func testMySQL() {
	engine, err := fexorm.NewEngine("mysql", "root@tcp(127.0.0.1:3306)/fexorm?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer engine.Close()

	exist := engine.NewSession().Model(&User{}).HasTable()
	if !exist {
		fmt.Println("table not exist")
		return
	}
	affected, err := engine.NewSession().Insert(&User{Name: "Tom", Age: 18})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Insert rows :", affected)

	var rows []User
	err = engine.NewSession().Where("name=?", "Tom").Find(&rows)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Find rows:", rows)
}
