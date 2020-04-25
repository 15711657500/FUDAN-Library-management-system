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
	if err != nil {
		t.Errorf(err.Error())
	}
}
