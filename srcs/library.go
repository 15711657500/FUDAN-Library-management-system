package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	USER     = "root"
	Password = "xxx"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", USER, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	err := fmt.Errorf("0")
	err = resetrent(lib)
	if err != nil {
		return err
	}
	err = resetusers(lib)
	if err != nil {
		return err
	}
	err = resetbook(lib)
	if err != nil {
		return err
	}
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
func (lib *Library) AddBook(title string, auther string, ISBN string) error {
	book1 := Book{title, auther, ISBN}
	err := addbook(&book1, lib)
	return err
}
func (lib *Library) AddSingleBook(ISBN string, bookid string) error {
	book1 := SingleBook{ISBN, bookid}
	err := addsinglebook(&book1, lib)
	return err
}
func (lib *Library) CreateUser(username string, password string, admin bool) error {
	user1 := User{username, password}
	err := createuser(&user1, lib, admin)
	return err
}

func (lib *Library) Login(username string, password string) error {
	user1 := User{username, password}
	err := login(&user1, lib)
	return err
}
func (lib *Library) Rent(book *Book, user *User) error {
	err := rent(book, user, lib)
	return err
}
func (lib *Library) Query(input string, mode string) error {
	var books []Book
	err := fmt.Errorf("0")
	switch mode {
	case "ISBN":
		books, err = querybookbyISBN(input, lib)
	case "author":
		books, err = querybookbyauthor(input, lib)
	case "title":
		books, err = querybookbytitle(input, lib)
	default:
		err = fmt.Errorf("Wrong Mode!")
	}
	if err != nil {
		return err
	}
	if books == nil {
		println("Not found!")
		return nil
	}
	for _, value := range books {
		println(value.title, "\t", value.author, "\t", value.ISBN)
	}
	return nil
}
func main() {
	fmt.Println("Welcome to the Library Management System!")
}
