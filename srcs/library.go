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
	_, err := lib.db.Exec(`
drop table if exists book;
create table if not exists book
(
    title     varchar(50),
    author    varchar(50),
    ISBN      varchar(100) primary key,
    total     int,
    constraint tcs check (total >= 0),
    remaining int,
    constraint rcs check (remaining <= total and remaining >= 0)
)
;
	`)
	return err
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, auther, ISBN string) error {
	return nil
}

// etc...

func main() {
	fmt.Println("Welcome to the Library Management System!")
}
