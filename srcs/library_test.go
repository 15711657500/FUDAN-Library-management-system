package main

import (
	"testing"
)

// create tables
func TestCreateTables(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
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
		{"Architecture of a Database System", "Hellerstein", "0"},
		{"Database systems : the complete book", "Garcia-Molina", "1"},
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
		{"30", "", "20", 1},
		{"22", "", "22", 1},
		{"23", "", "22", 1},
		{"24", "", "23", 1},
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

func TestQueryBooks(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("Unable to connect")
	}
	// query books by ISBN
	ISBNs := []string{
		"20",
		"21",
		"22",
		"23",
	}
	for _, value := range ISBNs {
		_, err = querybookbyISBN(value, &lib)
		if err != nil {
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
	for _, value := range ISBNs {
		_, err := querysinglebookbyISBN(value, &lib)
		if err != nil {
			t.Errorf("Fail to query singlebooks!")
		}
	}
}

func Test(t *testing.T) {

}
