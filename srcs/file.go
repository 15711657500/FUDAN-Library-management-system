package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// read users from csv file
func readuser(filename string) ([]User, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	reader := csv.NewReader(fs)
	var users []User
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		u, p := row[0], row[1]
		r, err := strconv.Atoi(row[2])
		if err != nil {
			return nil, err
		}
		users = append(users, User{u, p, r})
	}
	return users, nil
}

// read books from csv file
func readbook(filename string) ([]Book, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	reader := csv.NewReader(fs)
	var books []Book
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		t, a, i := row[0], row[1], row[2]
		books = append(books, Book{t, a, i})
	}
	return books, nil
}

// read singlebooks from csv file
func readsinglebook(filename string) ([]SingleBook, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	reader := csv.NewReader(fs)
	var books []SingleBook
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		b, i := row[0], row[1]
		books = append(books, SingleBook{b, "", i, 1})
	}
	return books, nil
}
