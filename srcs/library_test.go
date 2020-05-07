package main

import (
	"fmt"
	"testing"
)

/* If the tables are initially not empty,
* the primary key of existing items may contradict with item to be added
* and causes error in tests
* I won't drop the tables in this test, but you should be cautious to this
 */

// create tables
func TestCreateTables(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	err = lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
}

// add users, books, singlebooks

func TestAddinformation(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	// add users
	users := []User{
		{"18307130001", "1111", 0},
		{"18307130002", "1111", 0},
		{"18307130003", "1111", 0},
		{"root1", "root1", 1},
		{"root2", "root2", 1},
	}
	for _, value := range users {
		err = adduser(&value, &lib)
		if err != nil {
			t.Errorf("Fail to add a user!")
		}
	}
	// add books
	books := []Book{
		{"book20", "1", "20"},
		{"book21", "2", "21"},
		{"book21", "3", "22"},
		{"book23", "3", "23"},
		{"Architecture of a Database System", "Hellerstein", "799"},
		{"Database systems : the complete book", "Garcia-Molina", "800"},
	}
	for _, value := range books {
		err = addbook(&value, &lib)
		if err != nil {
			t.Errorf("Fail to add books!")
		}
	}
	// add books using batch insert
	err = addbook_batch(&books, &lib)
	if err != nil {
		t.Errorf("Fail to add books using batch insert!")
	}
	// add singlebooks
	sbooks := []SingleBook{
		{"51", "", "20", 1},
		{"52", "", "22", 1},
		{"53", "", "22", 1},
		{"54", "", "23", 1},
		{"1", "", "23", 1},
		{"2", "", "23", 1},
		{"3", "", "23", 1},
		{"4", "", "23", 1},
		{"5", "", "23", 1},
		{"6", "", "23", 1},
	}
	for _, value := range sbooks {
		err = addsinglebook(&value, &lib)
		if err != nil {
			t.Errorf("Fail to add singlebooks!")
		}
	}
	// add singlebooks using batch insert
	err = addsinglebook_batch(&sbooks, &lib)
	if err != nil {
		t.Errorf("Fail to add singlebooks using batch insert!")
	}
}

// query books by ISBN, title, author; query singlebooks by ISBN
func TestQueryBooks(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	// query books by ISBN
	ISBNs := []struct {
		Title string
		ISBN  string
	}{
		{"book20", "20"},
		{"book21", "21"},
		{"book21", "22"},
		{"book23", "23"},
	}
	for _, value := range ISBNs {
		books, err := querybookbyISBN(value.ISBN, &lib) // only return 1 book, I check the title of the books
		if err != nil || len(books) != 1 || books[0].Title != value.Title {
			t.Errorf("Fail to query books by ISBN!")
		}
	}
	// query books by author
	authors := []struct {
		author string // single keywords
		least  int    // each query may receive mutiple books, at least this number of books
	}{
		{"1", 1},
		{"3", 2},
		{"i", 2}, // Hellerstein and Garcia-Molina
	}
	for _, value := range authors {
		books, err := querybookbyauthor(value.author, &lib)
		if err != nil {
			t.Errorf("Fail to query books by author!")
		}
		if len(books) < value.least {
			t.Errorf("Insert error or query error in querying by author!")
		}
	}
	//query books by title
	titles := []struct {
		title []string // mutiple keywords
		least int      // each query may receive mutiple books, at least this number of books
	}{
		{[]string{"1"}, 2},                     // book21 and book21
		{[]string{"book", "2"}, 4},             // book20, book21, book21, book23
		{[]string{"Data", "base", "ystem"}, 2}, //Architecture of a Database System, Database systems : the complete book
	}
	for _, value := range titles {
		books, err := querybookbytitle(value.title, &lib)
		if err != nil {
			t.Errorf("Fail to query books by author!")
		}
		if len(books) < value.least {
			t.Errorf("Insert error or query error in querying by author!")
		}
	}
	// query singlebooks by ISBN
	ISBNs2 := []struct {
		ISBN  string
		least int
	}{
		{"20", 1},
		{"22", 2},
		{"23", 7},
	}
	for _, value := range ISBNs2 {
		book, err := querysinglebookbyISBN(value.ISBN, &lib)
		if err != nil || len(book) < value.least {
			t.Errorf("Fail to query singlebooks!")
		}
	}
}

// login
func TestLogin(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	users := []struct {
		user User
		err  error
	}{
		{User{"18307130001", "1111", 1}, nil}, // here the root is useless
		{User{"root1", "root1", 1}, nil},
		{User{"18307130001", "1234", 0}, loginerror}, // wrong password
		{User{"1234", "1234", 0}, loginerror},        // invalid username
	}
	for _, value := range users {
		err = login(&value.user, &lib)
		if (err == nil && value.err == nil) || (err != nil && value.err != nil && err.Error() == value.err.Error()) {
			logout()
		} else {
			t.Errorf("Error in login!")
		}
	}
}

// test borrow, return, exceed borrow limit, extend, appoint
func TestBorrow(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	err = login(&User{"18307130001", "1111", 0}, &lib)
	if err != nil || username != "18307130001" {
		t.Errorf("Failed to login!")
		return
	}
	// borrow books
	borrow1 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
		{"999", restricted},
	}
	for _, value := range borrow1 {
		err = rentsinglebook(value.bookid, username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Borrow error")
		}
	}
	// return books
	return1 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
		{"999", returnnotfound}, // 18307130001 did not borrow book 999
	}
	for _, value := range return1 {
		err = returnsinglebook(value.bookid, username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Return error")
		}
	}
	// exceed maxrent limit
	maxrent = 2 // default 30
	borrow2 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
		{"3", restricted},
	}
	for _, value := range borrow2 {
		err = rentsinglebook(value.bookid, username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Check maxrent limit error")
		}
	}
	return2 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
	}
	for _, value := range return2 {
		err = returnsinglebook(value.bookid, username, &lib)
		if err != nil {
			t.Errorf("Return error")
		}
	}
	maxrent = 30

	// exceed overdue limit, default 3
	du = -1 // du = -1 means: Once I borrow a book, it's overdue. du is default 30
	borrow3 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
		{"3", nil},
		{"4", nil},
		{"5", restricted},
	}
	for _, value := range borrow3 {
		err = rentsinglebook(value.bookid, username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Check maxoverdue limit error")
		}
	}
	return3 := []struct {
		bookid string
		err    error
	}{
		{"1", nil},
		{"2", nil},
		{"3", nil},
		{"4", nil},
	}
	for _, value := range return3 {
		err = returnsinglebook(value.bookid, username, &lib)
		if err != nil {
			t.Errorf("Return error")
		}
	}
	du = 30

	//extend
	err = rentsinglebook("1", username, &lib)
	if err != nil {
		t.Errorf("Borrow error")
	}
	errs := []error{
		nil,
		nil,
		nil,
		extenderror, // Users cannot extend more than three times in one rent
	}
	for _, value := range errs {
		err = extend("1", username, &lib)
		if !(err == nil && value == nil) && !(err != nil && value != nil && err.Error() == value.Error()) {
			t.Errorf("Extend error")
		}
	}
	err = returnsinglebook("1", username, &lib)
	if err != nil {
		t.Errorf("Return error")
	}

	// appoint
	// 18307130001 has borrowed book 1, now 18307130002 can appoint book 1
	err = rentsinglebook("1", "18307130001", &lib)
	if err != nil {
		t.Errorf("Borrow error")
	}
	appoints := []struct {
		bookid   string
		username string
		err      error
	}{
		{"1", "18307130001", appointborrowed},                 // Users cannot appoint books they have borrowed
		{"1", "18307130002", nil},                             // Other users can appoint books which have been borrowed by others
		{"2", "18307130001", appointavailable},                //Users cannot appoint books which is available
		{"book not in booklist", "18307130001", booknotfound}, //invalid bookid
		{"1", "18307130002", appointed},                       // Users cannot appoint a book twice before they borrow it
		{"1", "18307130003", nil},
	}
	for _, value := range appoints {
		_, err = appoint(value.bookid, value.username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Appoint error")
		}
	}
	err = returnsinglebook("1", "18307130001", &lib)
	if err != nil {
		t.Errorf("Return error")
	}
	// Now 18307130002 and 18307130003 have appointed book 1
	// 183007130002 first appoints, he can borrow it
	borrow4 := []struct {
		bookid   string
		username string
		err      error
	}{
		{"1", "18307130001", restricted}, // 18307130001 has not appointed
		{"1", "18307130003", restricted}, // 18307130003 has to wait
		{"1", "18307130002", nil},
	}
	for _, value := range borrow4 {
		err = rentsinglebook(value.bookid, value.username, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Appoint selection error")
		}
	}
	err = returnsinglebook("1", "18307130002", &lib)
	if err != nil {
		t.Errorf("Return error")
	}
	err = rentsinglebook("1", "18307130003", &lib)
	if err != nil {
		t.Errorf("Borrow error")
	}
	err = returnsinglebook("1", "18307130003", &lib)
	if err != nil {
		t.Errorf("Return error")
	}
	// logout
	logout()
}

// query the borrow history, books not returned, duedate of a book and whether a user has overdue books
func TestQueryRecord(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	// query rent record
	record1 := []struct {
		username string
		least    int
	}{
		{"18307130001", 10},
		{"18307130002", 1},
		{"invalid username", 0},
	}
	for _, value := range record1 {
		r, err := queryrentrecord(value.username, &lib)
		if err != nil || len(r) < value.least {
			t.Errorf("Query rent record error")
		}
	}

	// query overdue books
	du = -1
	err = rentsinglebook("4", "18307130002", &lib)
	if err != nil {
		t.Errorf("Borrow error")
	}
	rent1, err := queryoverdue("18307130002", &lib)
	if err != nil || rent1 == nil || len(rent1) != 1 || rent1[0].Bookid != "4" {
		t.Errorf("Query overdue error!")
	}
	du = 30
	// query the duedate of an unreturned book
	duedate, err := queryduedate("4", &lib)
	if err != nil || duedate == "" {
		t.Errorf("Query duedate error!")
	}

	// query books a user has not returned
	rent2, err := querynotreturned("18307130002", &lib)
	if err != nil || rent2 == nil || len(rent2) != 1 || rent2[0].Bookid != "4" {
		t.Errorf("Query notreturned book error!")
	}
}

// remove books
func TestRemovebooks(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	removelist := []struct {
		bookid string
		datail string
		err    error
	}{
		{"4", "lost", notreturned}, // book 4 was borrowed by 18307130002 and not returned
		{"5", "lost", nil},
		{"this book does not exist", "lost", booknotfound},
	}
	for _, value := range removelist {
		err = removesinglebook(value.bookid, value.datail, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("Remove book error!")
			fmt.Println(err)
		}
	}
}

func TestTopten(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	// data is limited, I only test the most popular book
	books, err := topten(&lib)
	if err != nil || len(books) == 0 || books[0].ISBN != "23" {
		t.Errorf("Topten book error!")
	}
}

// read from csv
func TestReadfromfile(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	ans := []Book{
		{"Architecture of a Database System", "Hellerstein", "0"},
		{"Database systems : the complete book", "Garcia-Molina", "1"},
		{"Introduction to algorithms", "Cormen", "2"},
		{"book3", "author1", "3"},
		{"book3", "author2", "4"},
		{"book4", "author2", "5"},
		{"book5", "author3", "6"},
	}
	books, err := readbook("../data/books.csv")
	if err != nil || len(books) != len(ans) {
		t.Errorf("read csv file error!")
	}
	if len(books) == len(ans) {
		for index, value := range books {
			if value != ans[index] {
				t.Errorf("read csv file error!")
				break
			}
		}
	}
}

func TestChangePassword(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
		return
	}
	err = changepassword("18307130002", "1234", &lib)
	if err != nil {
		t.Errorf("change password error")
	}
	trials := []struct {
		password string
		err      error
	}{
		{"1111", loginerror}, //initial password
		{"1234", nil},        //new password
	}
	for _, value := range trials {
		err = login(&User{"18307130002", value.password, 1}, &lib)
		if !(err == nil && value.err == nil) && !(err != nil && value.err != nil && err.Error() == value.err.Error()) {
			t.Errorf("change password error")
		}
	}
}
