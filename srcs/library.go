package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	USER     = "root"
	Password = "xxx"
	DBName   = "ass3"
)

var (
	done     = false
	visitor  = true
	username = "visitor"
	root     = 0
	reader   = bufio.NewReader(os.Stdin)
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() error {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", USER, Password, DBName))
	if err != nil {
		return err
	}
	lib.db = db
	return nil
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	err := fmt.Errorf("0")
	//err = resetrent(lib)
	//if err != nil {
	//	return err
	//}
	//err = resetusers(lib)
	//if err != nil {
	//	return err
	//}
	//err = resetbook(lib)
	//if err != nil {
	//	return err
	//}
	err = createbook(lib)
	if err != nil {
		return err
	}
	err = createusers(lib)
	if err != nil {
		return err
	}
	err = createrent(lib)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to open. Please try again.")
		return
	}
	err = lib.CreateTables()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to open. Please try again.")
		return
	}
	fmt.Println("Welcome to the Library Management System!")
	for {
		output := fmt.Sprintf("%s@FUDAN<", username)
		print(output)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Unexpected error! Force to quit!")
			return
		}
		input = strings.TrimSpace(input)
		if input != "" {
			handleinput(input, &lib)
		}
		if done {
			fmt.Println("Bye!")
			break
		}
	}
}
