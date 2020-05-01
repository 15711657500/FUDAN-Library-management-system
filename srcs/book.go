package main

import (
	"fmt"
)

type Book struct {
	Title  string
	Author string
	ISBN   string
}
type SingleBook struct {
	Bookid    string
	Title     string
	ISBN      string
	Available int
}

func resetbook(lib *Library) error {

	_, err := lib.db.Exec(`drop table if exists singlebook`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`drop table if exists booklist`)
	return err
}
func createbook(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists booklist
(
    title     nvarchar(200),
    author    nvarchar(200),
    ISBN      nvarchar(200) primary key,
    visits 	  int default 0
);
`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`
	create table if not exists singlebook
(
    ISBN      nvarchar(200),
    bookid 	  nvarchar(200) primary key,
    available bool default true,
    foreign key (ISBN) references booklist(ISBN) on delete cascade
)
`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`
	create table if not exists removelist
(
    bookid nvarchar(200) primary key,
    detail nvarchar(200) default "The book is lost.",
    foreign key (bookid) references singlebook(bookid) on delete cascade
)
`)
	return err
}
func addbook(book *Book, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into booklist(title, author, ISBN) values ('%s','%s','%s')", book.Title, book.Author, book.ISBN)
	_, err := lib.db.Exec(exec)
	return err
}
func addsinglebook(book *SingleBook, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into singlebook(ISBN, bookid) values ('%s','%s')", book.ISBN, book.Bookid)
	_, err := lib.db.Exec(exec)
	return err
}
func removesinglebook(bookid string, detail string, lib *Library) error {
	query := fmt.Sprintf("select count(*) from singlebook where bookid = '%s' and available = 1", bookid)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return err
		}
		if i != 1 {
			return fmt.Errorf("removed of not returned")
		}
	}
	exec := fmt.Sprintf("insert ignore into removelist(bookid, detail) values ('%s','%s')", bookid, detail)
	_, err = lib.db.Exec(exec)
	return err
}
func addbook_batch(books *[]Book, lib *Library) error {
	exec := "insert ignore into booklist(title, author, ISBN) values "
	if len(*books) < 1 {
		return nil
	}
	for index, value := range *books {
		t, a, i := value.Title, value.Author, value.ISBN
		exec = exec + fmt.Sprintf("('%s', '%s', '%s')", t, a, i)
		if index < len(*books)-1 {
			exec = exec + ","
		}
	}
	_, err := lib.db.Exec(exec)
	return err
}
func addsinglebook_batch(books *[]SingleBook, lib *Library) error {
	exec := "insert ignore into singlebook(bookid, ISBN) values "
	if len(*books) < 1 {
		return nil
	}
	for index, value := range *books {
		b, i := value.Bookid, value.ISBN
		exec = exec + fmt.Sprintf("('%s','%s')", b, i)
		if index < len(*books)-1 {
			exec = exec + ","
		}
	}
	_, err := lib.db.Exec(exec)
	return err
}
