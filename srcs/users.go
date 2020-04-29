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

var (
	loginerror = fmt.Errorf("No such user, or wrong password!")
)

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
	var exec string
	if user.root != 0 {
		exec = fmt.Sprintf("insert ignore into users(username, password, admin) values ('%s', '%s', 1)", username1, password1)
	} else {
		exec = fmt.Sprintf("insert ignore into users(username, password) values ('%s', '%s')", username1, password1)
	}
	_, err = lib.db.Exec(exec)
	if err != nil {
		return err
	}

	return nil
}
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
		rows.Scan(&j)
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
	}
	if password1 != password2 {
		return loginerror
	}
	root = r
	return nil
}
func createuser_batch(user *[]User, lib *Library) error {
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
