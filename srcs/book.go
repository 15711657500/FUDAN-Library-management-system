package main

import "fmt"

type Book struct {
	title  string
	author string
	ISBN   string
}
type SingleBook struct {
	ISBN   string
	bookid string
}

func resetbook(lib *Library) error {
	_, err := lib.db.Exec(`drop table if exists booklist`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`drop table if exists singlebook`)
	return err
}
func createbook(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists booklist
(
    title     varchar(50),
    author    varchar(50),
    ISBN      varchar(100) primary key,
    visits 	  int default 0
);
`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`
	create table if not exists singlebook
(
    ISBN      varchar(100) references booklist(ISBN) on delete cascade,
    bookid 	  varchar(100) primary key,
    available bool default true
)
`)
	return err
}
func addbook(book *Book, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into booklist(title, author, ISBN) values ('%s','%s','%s')", book.title, book.author, book.ISBN)
	_, err := lib.db.Exec(exec)
	return err
}
func addsinglebook(book *SingleBook, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into singlebook(ISBN, bookid) values ('%s','%s')", book.ISBN, book.bookid)
	_, err := lib.db.Exec(exec)
	return err
}
