package main

import (
	"fmt"
	"time"
)

var (
	du         time.Duration = 30
	maxrent    int           = 30
	maxoverdue int           = 3
)

const (
	dateformat = "2006-01-02 15:04:05"
)

func resetrent(lib *Library) error {
	_, err := lib.db.Exec(`
	drop table if exists rent;
`)
	return err
}
func createrent(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists rent
(
    rentdate nvarchar(200),
    duedate nvarchar(200),
    returndate nvarchar(200) default "not returned yet",
    fine float default 0,
    rentid int primary key auto_increment,
    username nvarchar(200),
    bookid nvarchar(200),
    extend int default 0,
	foreign key (username) references users(username),
	foreign key (bookid) references singlebook(bookid)
);
`)
	return err
}
func rent(book *Book, user *User, lib *Library) error {
	// choose an available book
	notfound := fmt.Errorf("111")
	bookid := 0
	title, author, ISBN := book.title, book.author, book.ISBN
	found := false
	query1 := fmt.Sprintf("select bookid from book where title = '%s' and author = '%s' and ISBN = '%s' and available = 1 order by bookid desc", title, author, ISBN)
	rows, err := lib.db.Queryx(query1)
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan(&bookid)
		found = true
	}
	if !found {
		return notfound
	}

	// set available = 0, insert information into rent

	exec1 := fmt.Sprintf("update book set available = 0 where bookid = %d", bookid)
	_, err = lib.db.Exec(exec1)
	if err != nil {
		return err
	}
	rentdate := time.Now().Format(dateformat)
	duedate := time.Unix(time.Now().Unix(), int64(du*time.Hour*24)).Format(dateformat)
	exec2 := fmt.Sprintf("insert into rent(rentdate, duedate, username, bookid) values ('%s','%s','%s','%d')", rentdate, duedate, user.username, bookid)
	_, err = lib.db.Exec(exec2)
	if err != nil {
		return err
	}
	return nil
}
func querybookbyISBN(ISBN string, lib *Library) ([]Book, error) {
	var books []Book
	query := fmt.Sprintf("select title, author, ISBN from booklist where ISBN = '%s'", ISBN)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, b, c string
		err = rows.Scan(&a, &b, &c)
		if err != nil {
			return nil, err
		}
		books = append(books, Book{a, b, c})
	}
	return books, nil
}
func querybookbyauthor(author string, lib *Library) ([]Book, error) {
	var books []Book
	query := fmt.Sprintf("select title, author, ISBN from booklist where author = '%s'", author)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, b, c string
		err = rows.Scan(&a, &b, &c)
		if err != nil {
			return nil, err
		}
		books = append(books, Book{a, b, c})
	}
	return books, nil
}
func querybookbytitle(title string, lib *Library) ([]Book, error) {
	var books []Book
	query := fmt.Sprintf("select title, author, ISBN from booklist where title like '%%%s%%'", title)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, b, c string
		err = rows.Scan(&a, &b, &c)
		if err != nil {
			return nil, err
		}
		books = append(books, Book{a, b, c})
	}
	return books, nil
}
func querysinglebookbyISBN(ISBN string, lib *Library) ([]SingleBook, error) {
	var books []SingleBook
	query := fmt.Sprintf("select ISBN, bookid from singlebook where ISBN = '%s'", ISBN)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, b string
		err = rows.Scan(&a, &b)
		if err != nil {
			return nil, err
		}
		books = append(books, SingleBook{a, b})
	}
	return books, nil
}
func checkrent(username string, lib *Library) (bool, error) {
	// exceed rent limit
	query1 := fmt.Sprintf("select count(*) from rent where username = '%s' and returndate = 'not returned yet'", username)
	rows1, err1 := lib.db.Queryx(query1)
	if err1 != nil {
		return false, err1
	}
	var i1 int
	for rows1.Next() {
		err1 = rows1.Scan(&i1)
		if err1 != nil {
			return false, err1
		}
	}
	if i1 > maxrent {
		return false, nil
	}
	// exceed overdue limit
	now := time.Now().Format(dateformat)
	query2 := fmt.Sprintf("select count(*) from rent where username = '%s' and returndate = 'not returned yet' and duedate < '%s'", username, now)
	rows2, err2 := lib.db.Queryx(query2)
	if err2 != nil {
		return false, err2
	}
	var i2 int
	for rows2.Next() {
		err2 = rows2.Scan(&i2)
		if err2 != nil {
			return false, err2
		}
	}
	if i2 > maxoverdue {
		return false, nil
	}
	return true, nil
}
