package main

import (
	"fmt"

	// mysql connector
	"crypto/sha256"
	"encoding/hex"
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

func getSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
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
	//err := resetbook(lib)
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
	err := createbook(lib)
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

func main() {
	fmt.Println("Welcome to the Library Management System!")
}
