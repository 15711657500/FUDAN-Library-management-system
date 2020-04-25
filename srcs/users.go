package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type User struct {
	username string
	password string
}

func resetusers(lib *Library) error {
	_, err := lib.db.Exec(`
drop table if exists users;
`)
	return err
}
func createusers(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists users
(
    username varchar(50) primary key,
    password varchar(150),
    permission varchar(10) default "default"
);
`)
	return err
}
func getSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}
func createuser(user *User, lib *Library) error {
	username1, password1 := user.username, user.password
	query := fmt.Sprintf("select count(*) from users where username = '%s'", username1)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return err
	}
	j := 0
	rows.Next()
	err = rows.Scan(&j)

	if err != nil {
		return err
	}
	if j != 0 {
		err = fmt.Errorf("already exists")
		return err
	}
	password1 = getSHA256(password1)
	exec := fmt.Sprintf("insert into users(username, password) values ('%s', '%s')", username1, password1)
	_, err = lib.db.Exec(exec)
	if err != nil {
		return err
	}
	return nil
}
func login(user *User, lib *Library) error {
	wrong := fmt.Errorf("No such user, or wrong password!")
	username1, password1 := user.username, user.password
	password1 = getSHA256(password1)
	query := fmt.Sprintf("select count(*) from users where username = '%s'", username1)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		j := 0
		rows.Scan(&j)
		if j == 0 {
			return wrong
		}
	}
	query2 := fmt.Sprintf("select password from users where username = '%s'", username1)
	rows2, err := lib.db.Queryx(query2)
	if err != nil {
		return err
	}
	var password2 string
	for rows2.Next() {
		err = rows2.Scan(&password2)
	}
	if password1 != password2 {
		return wrong
	}
	return nil
}
