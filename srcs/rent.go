package main

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	du         time.Duration = -1
	maxrent    int           = 30
	maxoverdue int           = 3
	fineperday float32       = 0.1
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
	returnnotfound   = fmt.Errorf("Unable to return!")
	booknotfound     = fmt.Errorf("Book not found!")
	appointavailable = fmt.Errorf("This book is already available!")
	appointborrowed  = fmt.Errorf("You've already borrowed this book!")
	appointed        = fmt.Errorf("You've already appointed this book!")
	restricted       = fmt.Errorf("Not able to borrow!")
	extenderror      = fmt.Errorf("Unable to extend!")
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
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(`
	create table if not exists appointment
(
    id int primary key auto_increment,
    username nvarchar(200),
    bookid nvarchar(200),
    borrowed nvarchar(200) default "No",
    foreign key (username) references users(username),
	foreign key (bookid) references singlebook(bookid)
)
`)
	return err
}

// borrow a book
func rentsinglebook(bookid string, username string, lib *Library) error {

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

	query2 := fmt.Sprintf("select username from appointment where bookid = '%s' and borrowed = 'No' order by id asc limit 1 ", bookid)
	rows2, err := lib.db.Queryx(query2)
	if err != nil {
		return err
	}
	var earliest string
	for rows2.Next() {
		err = rows2.Scan(&earliest)
		if err != nil {
			return err
		}
	}
	if i != 1 && earliest != username {
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
	if err != nil {
		return err
	}
	exec4 := fmt.Sprintf("update appointment set borrowed = 'Yes' where bookid = '%s' and username = '%s' and borrowed = 'No'", bookid, username)
	_, err = lib.db.Exec(exec4)
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
	query := fmt.Sprintf("select title, author, ISBN from booklist where author like '%%%s%%'", author)
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
	if i1 >= maxrent {
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
	query := fmt.Sprintf("select count(*), rentid, duedate from rent where bookid = '%s' and username = '%s' and returndate = 'not returned yet' group by rentid", bookid, username)
	rows, err := lib.db.Queryx(query)
	var rentid, i int
	var duedate string
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&i, &rentid, &duedate)
		if err != nil {
			return err
		}
	}
	if i != 1 {
		return returnnotfound
	}
	returndate := time.Now().Format(dateformat)
	duedate1, err := time.Parse(dateformat, duedate)
	if err != nil {
		return err
	}
	now, err := time.Parse(dateformat, time.Now().Format(dateformat))
	if err != nil {
		return err
	}
	var fine float32 = 0.0
	if duedate1.Unix() < now.Unix() {
		fine = float32(((-duedate1.Unix()+now.Unix()-1)/86400)+1) * fineperday
	}
	exec1 := fmt.Sprintf("update rent set returndate = '%s', fine = %f where rentid = %d", returndate, fine, rentid)
	_, err = lib.db.Exec(exec1)
	if err != nil {
		return err
	}

	i2, err := queryappoint(bookid, lib)
	if err != nil {
		return err
	}
	if i2 != 0 {
		return nil
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

// query the books a student has borrowed and not returned yet
func querynotreturned(username string, lib *Library) ([]Rent, error) {
	var rentdate, duedate, returndate, bookid, ISBN, title, author string
	var fine float32
	var rent []Rent
	query := fmt.Sprintf(`select rentdate, duedate, returndate, fine, rent.bookid, booklist.ISBN, title, author 
								from rent, booklist, singlebook 
								where username = '%s' and rent.bookid = singlebook.bookid 
								and singlebook.ISBN = booklist.ISBN and returndate = 'not returned yet'
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

// query the overdue books of a user
func queryoverdue(username string, lib *Library) ([]Rent, error) {
	var rentdate, duedate, returndate, bookid, ISBN, title, author string
	var fine float32
	var rent []Rent
	query := fmt.Sprintf(`select rentdate, duedate, returndate, fine, rent.bookid, booklist.ISBN, title, author 
								from rent, booklist, singlebook 
								where username = '%s' and returndate = 'not returned yet' and rent.bookid = singlebook.bookid 
								and singlebook.ISBN = booklist.ISBN
								order by rentid`, username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	now, _ := time.Parse(dateformat, time.Now().Format(dateformat))
	for rows.Next() {
		err = rows.Scan(&rentdate, &duedate, &returndate, &fine, &bookid, &ISBN, &title, &author)
		if err != nil {
			return nil, err
		}
		duedate1, _ := time.Parse(dateformat, duedate)
		if duedate1.Unix() < now.Unix() {
			rent = append(rent, Rent{Rentdate: rentdate, Duedate: duedate, Returndate: returndate, Fine: fine, Bookid: bookid, ISBN: ISBN, Title: title, Author: author})
		}
	}
	return rent, nil
}

// query the duedate of a book
func queryduedate(bookid string, lib *Library) (string, error) {
	query1 := fmt.Sprintf("select duedate from rent where returndate = 'not returned yet' and bookid = '%s'", bookid)
	rows1, err := lib.db.Queryx(query1)
	if err != nil {
		return "", err
	}
	var duedate string
	for rows1.Next() {
		err = rows1.Scan(&duedate)
		if err != nil {
			return "", err
		}
	}
	if duedate == "" {
		return "", booknotfound
	}
	return duedate, nil
}

// extend a book
func extend(bookid string, username string, lib *Library) error {

	query1 := fmt.Sprintf("select duedate, rentid from rent where returndate = 'not returned yet' and bookid = '%s' and username = '%s' and extend < 3", bookid, username)
	rows1, err := lib.db.Queryx(query1)
	if err == sql.ErrNoRows {
		return extenderror
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
		return extenderror
	}
	newduedate, err := time.Parse(dateformat, duedate)
	if err != nil {
		return err
	}
	newduedate1 := time.Unix(newduedate.Unix(), int64(du*time.Hour*24)).Format(dateformat)
	exec1 := fmt.Sprintf("update rent set duedate = '%s', extend = extend + 1 where rentid = %d", newduedate1, rentid)
	_, err = lib.db.Exec(exec1)
	return err
}

// appoint
func appoint(bookid string, username string, lib *Library) (int, error) {
	// wrong bookid?
	query1 := fmt.Sprintf("select count(*) from singlebook where bookid = '%s'", bookid)
	rows1, err := lib.db.Queryx(query1)
	if err != nil {
		return 0, err
	}
	i1 := 3
	for rows1.Next() {
		err = rows1.Scan(&i1)
		if err != nil {
			return 0, err
		}
	}
	if i1 == 0 {
		return 0, booknotfound
	}
	// this book is removed?
	query2 := fmt.Sprintf("select count(*) from removelist where bookid = '%s'", bookid)
	rows2, err := lib.db.Queryx(query2)
	if err != nil {
		return 0, err
	}
	i2 := 3
	for rows2.Next() {
		err = rows2.Scan(&i2)
		if err != nil {
			return 0, err
		}
	}
	if i2 == 1 {
		return 0, booknotfound
	}
	// this book is already awailable?
	query3 := fmt.Sprintf("select count(*) from singlebook where available = 1 and bookid = '%s'", bookid)
	rows3, err := lib.db.Queryx(query3)
	if err != nil {
		return 0, err
	}
	i3 := 3
	for rows3.Next() {
		err = rows3.Scan(&i3)
		if err != nil {
			return 0, err
		}
	}
	if i3 == 1 {
		return 0, appointavailable
	}
	// have borrowed this book?
	query4 := fmt.Sprintf("select count(*) from rent where returndate = 'not returned yet' and bookid = '%s' and username = '%s'", bookid, username)
	rows4, err := lib.db.Queryx(query4)
	if err != nil {
		return 0, err
	}
	i4 := 3
	for rows4.Next() {
		err = rows4.Scan(&i4)
		if err != nil {
			return 0, err
		}
	}
	if i4 == 1 {
		return 0, appointborrowed
	}
	// have appointed this book?
	i5, err := checkappoint(bookid, username, lib)
	if i5 {
		return 0, appointed
	}
	// appoint
	exec := fmt.Sprintf("insert into appointment(bookid, username) values('%s','%s')", bookid, username)
	_, err = lib.db.Exec(exec)
	if err != nil {
		return 0, err
	}
	i6, err := queryappointbehinduser(bookid, username, lib)
	if err != nil {
		return 0, err
	}
	return i6, nil
}

// query how many users have appointed but not borrowed the book
func queryappoint(bookid string, lib *Library) (int, error) {
	query := fmt.Sprintf("select count(*) from appointment where bookid = '%s' and borrowed = 'No'", bookid)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return 0, err
	}
	i := 3
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return 0, err
		}
	}
	return i, nil
}

// query how many users have appointed but not borrowed the book before the user
func queryappointbehinduser(bookid string, username string, lib *Library) (int, error) {
	done, err := checkappoint(bookid, username, lib)
	if err != nil {
		return 0, err
	}
	if !done {
		return 0, err
	}
	query := fmt.Sprintf("select count(*) from appointment A where A.bookid = '%s' and A.borrowed = 'No' and id < all(select B.id from appointment B where B.bookid = '%s' and B.username = '%s' and B.borrowed = 'No')", bookid, bookid, username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return 0, err
	}
	var i int
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return 0, err
		}
	}
	return i, nil
}

// check where the user has appointed but not borrowed the book
func checkappoint(bookid string, username string, lib *Library) (bool, error) {
	query := fmt.Sprintf("select count(*) from appointment where bookid = '%s' and username = '%s' and borrowed = 'No'", bookid, username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return false, err
	}
	i := 0
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return false, err
		}
	}
	if i == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// query the books near duedate when logining
func loginduedate(username string, lib *Library) ([]Bookwithdate, error) {
	var books []Bookwithdate
	query := fmt.Sprintf("select bookid, duedate from rent where returndate = 'not returned yet' and username = '%s'", username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b, d string
		err = rows.Scan(&b, &d)
		if err != nil {
			return nil, err
		}
		duedate, err := time.Parse(dateformat, d)
		if err != nil {
			return nil, err
		}
		if duedate.Unix()-time.Now().Unix() > 3600*24*7 {
			continue
		}
		book, err := bookid2Book(b, lib)
		if err != nil {
			return nil, err
		}
		books = append(books, Bookwithdate{
			Bookid:  book.Bookid,
			Title:   book.Title,
			ISBN:    book.ISBN,
			DueDate: d,
		})
	}
	return books, nil
}

// query the books the user has appointed and is available now when logining
func loginappoint(username string, lib *Library) ([]Bookforappoint, error) {
	var books []Bookforappoint
	query := fmt.Sprintf(`
	select A.bookid from appointment A
	where A.username = '%s' 
	and A.borrowed = 'No'
	and A.id <= all(select B.id from appointment B where B.bookid = A.bookid and B.borrowed = 'No')
	and not exists (select * from rent where bookid = A.bookid and returndate = 'not returned yet')
`, username)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b string
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}
		book, err := bookid2Book(b, lib)
		if err != nil {
			return nil, err
		}
		books = append(books, Bookforappoint{
			Bookid: book.Bookid,
			Title:  book.Title,
			ISBN:   book.ISBN,
		})
	}
	return books, nil
}

// query the topten popular books
func topten(lib *Library) ([]Bookwithvisit, error) {
	var books []Bookwithvisit
	rows, err := lib.db.Queryx("select title, author, ISBN from booklist order by visits desc limit 10")
	if err != nil {
		return nil, err
	}
	j := 1
	for rows.Next() {
		var t, a, i string
		err = rows.Scan(&t, &a, &i)
		if err != nil {
			return nil, err
		}
		books = append(books, Bookwithvisit{t, a, i, j})
		j = j + 1
	}
	return books, nil
}
