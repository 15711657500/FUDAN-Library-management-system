package main

import (
	"fmt"
	"time"
)

const (
	du = 30
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
    rentdate varchar(50),
    duedate varchar(50),
    returndate varchar(50) default "not returned yet",
    fine float default 0,
    rentid int primary key auto_increment,
    username varchar(50) references users(username),
    bookid int references singlebook(bookid),
    extend int default 0
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
	dateformat := "2006-01-02 15:04:05"
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
