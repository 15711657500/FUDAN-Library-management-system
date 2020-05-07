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
type Bookwithdate struct {
	Bookid  string
	Title   string
	ISBN    string
	DueDate string
}
type Bookforappoint struct {
	Bookid string
	Title  string
	ISBN   string
}
type Bookwithvisit struct {
	Title  string
	Author string
	ISBN   string
	rank   int
}

var (
	notreturned = fmt.Errorf("removed of not returned")
)

// drop table booklist and singlebook
func resetbook(lib *Library) error {

	_, err := lib.db.Exec(`drop table if exists singlebook`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`drop table if exists booklist`)
	return err
}

// create table booklist and singlebook
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

// add a book to booklist
func addbook(book *Book, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into booklist(title, author, ISBN) values ('%s','%s','%s')", book.Title, book.Author, book.ISBN)
	_, err := lib.db.Exec(exec)
	return err
}

// add a singlebook to singlebook
func addsinglebook(book *SingleBook, lib *Library) error {
	exec := fmt.Sprintf("insert ignore into singlebook(ISBN, bookid) values ('%s','%s')", book.ISBN, book.Bookid)
	_, err := lib.db.Exec(exec)
	return err
}

// add a singlebook to removelist and set available = 0 in table singlebook
func removesinglebook(bookid string, detail string, lib *Library) error {
	query := fmt.Sprintf("select count(*) from rent where bookid = '%s' and returndate = 'not returned yet'", bookid)
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
		if i != 0 {
			return notreturned
		}
	}
	exec := fmt.Sprintf("insert ignore into removelist(bookid, detail) values ('%s','%s')", bookid, detail)
	_, err = lib.db.Exec(exec)
	if err != nil {
		return err
	}
	exec2 := fmt.Sprintf("update singlebook set available = 0 where bookid = '%s'", bookid)
	_, err = lib.db.Exec(exec2)
	return err
}

// add books to booklist, using batch insert
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

// add singlebooks to singlebook, using batch insert
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

func bookid2Book(bookid string, lib *Library) (SingleBook, error) {
	var b, t, i string
	var a int
	query := fmt.Sprintf("select bookid, title, singlebook.ISBN, available from singlebook, booklist where bookid = '%s' and singlebook.ISBN = booklist.ISBN", bookid)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return SingleBook{"", "", "", 0}, err
	}
	for rows.Next() {
		err = rows.Scan(&b, &t, &i, &a)
		if err != nil {
			return SingleBook{"", "", "", 0}, err
		}
	}
	return SingleBook{b, t, i, a}, nil
}
