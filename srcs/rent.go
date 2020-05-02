package main

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	du         time.Duration = 30
	maxrent    int           = 30
	maxoverdue int           = 3
)

type Rent struct {
	Rentdate   string
	Duedate    string
	Returndate string
	Fine       float32
	Bookid     string
	ISBN       string
	Title      string
	Author     string
}

const (
	dateformat = "2006-01-02 15:04:05"
)

// errors
var (
	returnnotfound = fmt.Errorf("Unable to return!")
)

// drop table rent
func resetrent(lib *Library) error {
	_, err := lib.db.Exec(`
	drop table if exists rent;
`)
	return err
}

// create table rent
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

// borrow a book
func rentsinglebook(bookid string, username string, lib *Library) error {
	restricted := fmt.Errorf("Not able to borrow!")
	// check the user
	able, err := checkrent(username, lib)
	if err != nil {
		return err
	}
	if !able {
		return restricted
	}
	query1 := fmt.Sprintf("select count(*), ISBN from singlebook where bookid = '%s' and available = 1 group by bookid", bookid)
	rows, err := lib.db.Queryx(query1)
	if err != nil {
		return err
	}
	var i int
	var ISBN string
	for rows.Next() {
		err = rows.Scan(&i, &ISBN)
		if err != nil {
			return err
		}

	}
	if i != 1 {
		return restricted
	}
	rentdate := time.Now().Format(dateformat)
	duedate := time.Unix(time.Now().Unix(), int64(du*time.Hour*24)).Format(dateformat)
	exec1 := fmt.Sprintf("insert into rent(bookid, username, rentdate, duedate) values ('%s','%s', '%s', '%s')", bookid, username, rentdate, duedate)
	_, err = lib.db.Exec(exec1)
	if err != nil {
		return err
	}
	exec2 := fmt.Sprintf("update singlebook set available = 0 where bookid = '%s'", bookid)
	_, err = lib.db.Exec(exec2)
	if err != nil {
		return err
	}
	exec3 := fmt.Sprintf("update booklist set visits = visits + 1 where ISBN = '%s'", ISBN)
	_, err = lib.db.Exec(exec3)
	return err
}

// query the booklist by ISBN
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

// query the booklist by author
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

// query the booklist by title
func querybookbytitle(title []string, lib *Library) ([]Book, error) {
	var books []Book
	// query := fmt.Sprintf("select title, author, ISBN from booklist where title like '%%%s%%'", title)
	query := "select title, author, ISBN from booklist "
	for index, value := range title {
		if index == 0 {
			query = query + "where "
		} else {
			query = query + "and "
		}
		query = query + fmt.Sprintf("title like '%%%s%%' ", value)
	}
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

// query table singlebook by ISBN
func querysinglebookbyISBN(ISBN string, lib *Library) ([]SingleBook, error) {
	var books []SingleBook
	query := fmt.Sprintf("select bookid, title, singlebook.ISBN, available from singlebook, booklist where singlebook.ISBN = '%s' and singlebook.ISBN = booklist.ISBN", ISBN)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, b, c string
		var d int
		err = rows.Scan(&a, &b, &c, &d)
		if err != nil {
			return nil, err
		}
		books = append(books, SingleBook{
			Bookid:    a,
			Title:     b,
			ISBN:      c,
			Available: d,
		})
	}
	return books, nil
}

// check whether the user has exceeded the borrow limit or overdue limit
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

// return a book
func returnsinglebook(bookid string, username string, lib *Library) error {
	query := fmt.Sprintf("select count(*), rentid from rent where bookid = '%s' and username = '%s' and returndate = 'not returned yet' group by rentid", bookid, username)
	rows, err := lib.db.Queryx(query)
	var rentid, i int
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&i, &rentid)
		if err != nil {
			return err
		}
	}
	if i != 1 {
		return returnnotfound
	}
	returndate := time.Now().Format(dateformat)
	exec1 := fmt.Sprintf("update rent set returndate = '%s' where rentid = %d", returndate, rentid)
	_, err = lib.db.Exec(exec1)
	if err != nil {
		return err
	}
	exec2 := fmt.Sprintf("update singlebook set available = 1 where bookid = '%s'", bookid)
	_, err = lib.db.Exec(exec2)
	return err
}

// query the borrow record of a user
func queryrentrecord(username string, lib *Library) ([]Rent, error) {
	var rentdate, duedate, returndate, bookid, ISBN, title, author string
	var fine float32
	var rent []Rent
	query := fmt.Sprintf(`select rentdate, duedate, returndate, fine, rent.bookid, booklist.ISBN, title, author 
								from rent, booklist, singlebook 
								where username = '%s' and rent.bookid = singlebook.bookid 
								and singlebook.ISBN = booklist.ISBN
								order by rentdate`, username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&rentdate, &duedate, &returndate, &fine, &bookid, &ISBN, &title, &author)
		if err != nil {
			return nil, err
		}
		rent = append(rent, Rent{Rentdate: rentdate, Duedate: duedate, Returndate: returndate, Fine: fine, Bookid: bookid, ISBN: ISBN, Title: title, Author: author})

	}
	return rent, nil
}

// query the duedate of a book
func queryduedate(bookid string, lib *Library) (string, error) {
	notfound := fmt.Errorf("You have not borrowed this book!")
	query1 := fmt.Sprintf("select count(*), duedate from rent where returndate = 'not returned yet' and bookid = '%s'", bookid)
	rows1, err := lib.db.Queryx(query1)
	if err != nil {
		return "", err
	}
	var duedate string
	var i int
	for rows1.Next() {
		err = rows1.Scan(&i, &duedate)
		if err != nil {
			return "", err
		}

	}
	if i != 1 {
		return "", notfound
	}
	return duedate, nil
}

// extend a book
func extend(bookid string, username string, lib *Library) error {
	fail := fmt.Errorf("Unable to extend!")
	query1 := fmt.Sprintf("select duedate, rentid from rent where returndate = 'not returned yet' and bookid = '%s' and username = '%s'", bookid, username)
	rows1, err := lib.db.Queryx(query1)
	if err == sql.ErrNoRows {
		return fail
	}
	if err != nil {
		return err
	}

	var duedate string
	var rentid int
	for rows1.Next() {
		err = rows1.Scan(&duedate, &rentid)
		if err != nil {
			return err
		}
	}
	if duedate == "" {
		return fail
	}
	newduedate, err := time.Parse(dateformat, duedate)
	if err != nil {
		return err
	}
	newduedate1 := time.Unix(newduedate.Unix(), int64(du*time.Hour*24)).Format(dateformat)
	exec1 := fmt.Sprintf("update rent set duedate = '%s' where rentid = %d", newduedate1, rentid)
	_, err = lib.db.Exec(exec1)
	return err
}
