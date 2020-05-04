package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type User struct {
	username string
	password string
	root     int
}

// errors
var (
	loginerror = fmt.Errorf("No such user, or wrong password!")
)

// drop table users
func resetusers(lib *Library) error {
	_, err := lib.db.Exec(`
drop table if exists users;
`)
	return err
}

// create table users, and add default root account
func createusers(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists users
(
    username nvarchar(200) primary key,
    password nvarchar(200),
    root bool default 0
);
`)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec(fmt.Sprintf(`
	insert ignore into users(username, password, root) 
	values ("root", "%s", 1)
`, getSHA256("root")))
	return err
}

// get SHA256 code for passwords
func getSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

// add an account to table users
func adduser(user *User, lib *Library) error {
	username1, password1, root1 := user.username, user.password, user.root
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

	exec := fmt.Sprintf("insert ignore into users(username, password, root) values ('%s', '%s', '%d')", username1, password1, root1)

	_, err = lib.db.Exec(exec)
	if err != nil {
		return err
	}

	return nil
}

// login
func login(user *User, lib *Library) error {

	username1, password1 := user.username, user.password
	password1 = getSHA256(password1)
	query := fmt.Sprintf("select count(*) from users where username = '%s'", username1)
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		j := 0
		err = rows.Scan(&j)
		if err != nil {
			return err
		}
		if j == 0 {
			return loginerror
		}
	}
	query2 := fmt.Sprintf("select password, root from users where username = '%s'", username1)
	rows2, err := lib.db.Queryx(query2)
	if err != nil {
		return err
	}
	var password2 string
	var r int
	for rows2.Next() {
		err = rows2.Scan(&password2, &r)
		if err != nil {
			return err
		}
	}
	if password1 != password2 {
		return loginerror
	}
	root = r
	visitor = false
	username = username1
	return nil
}

// add accounts to table users, using batch insert
func adduser_batch(user *[]User, lib *Library) error {
	exec := "insert ignore into users(username, password, root) values "
	if len(*user) < 1 {
		return nil
	}
	for index, value := range *user {
		u, p, r := value.username, getSHA256(value.password), value.root
		exec = exec + fmt.Sprintf("('%s','%s',%d)", u, p, r)
		if index < len(*user)-1 {
			exec = exec + ","
		}
	}
	_, err := lib.db.Exec(exec)
	return err
}

// change password of one's own account
func changepassword(username string, password string, lib *Library) error {
	exec1 := fmt.Sprintf("update users set password = '%s' where username = '%s'", getSHA256(password), username)
	_, err := lib.db.Exec(exec1)
	return err
}

// check if the password of the user is correct
func checkpassword(username string, password string, lib *Library) (bool, error) {
	query := fmt.Sprintf("select count(*) from users where username = '%s' and password = '%s'", username, getSHA256(password))
	rows, err := lib.db.Queryx(query)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return false, err
		}
		if i == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}
