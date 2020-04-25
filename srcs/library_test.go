package main

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}

}
func TestLibrary_CreateUser(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateUser("abc", "abc")
	if err != nil && err.Error() != "already exists" {
		t.Errorf(err.Error())
	}
}
func TestLibrary_Login(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.Login("abc", "abc")
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestLibrary_AddBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.AddBook("a", "b", "c")
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestLibrary_Rent(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.Rent(&Book{"a", "b", "c"}, &User{"abc", "abc"})
	if err != nil {
		t.Errorf(err.Error())
	}
}
