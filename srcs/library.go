package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "xxx"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	err := fmt.Errorf("0")
	//err = resetbook(lib)
	//if err != nil {
	//	return err
	//}
	//err = resetusers(lib)
	//if err != nil {
	//	return err
	//}
	//err = resetrent(lib)
	//if err != nil {
	//	return err
	//}
	err = createbook(lib)
	if err != nil {
		return err
	}
	err = createusers(lib)
	if err != nil {
		return err
	}
	err = createrent(lib)
	if err != nil {
		return err
	}
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, auther, ISBN string) error {

	return nil
}

// etc...
func (lib *Library) CreateUser(username string, password string) error {
	user1 := user{username, password}
	err := createuser(&user1, lib)
	return err
}
func main() {
	fmt.Println("Welcome to the Library Management System!")
}
