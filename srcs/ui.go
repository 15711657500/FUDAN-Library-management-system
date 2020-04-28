package main

import "C"
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	//"golang.org/x/crypto/ssh"
	terminal "golang.org/x/crypto/ssh/terminal"
)

var (
	books       []Book
	singlebooks []SingleBook
)

const (
	help = `Type "help" can get this help.
quit			quit
login			login
logout			logout
ISBN <ISBN> 	search by ISBN
title <title>   search by title
author <author> search by author
bookid <ISBN>   get bookid of books whose ISBN is <ISBN>"
borrow <bookid> borrow the book whose id is <bookid>
return <bookid> 
`
)

func handleinput(input string, lib *Library) {
	args := strings.Split(input, " ")
	switch args[0] {
	case "quit":
		if len(args) == 1 {
			done = true
		} else {
			print(help)
		}
	case "login":
		if len(args) == 1 {
			print("username:")
			user1, err := bufio.NewReader(os.Stdin).ReadString('\n')

			if err != nil {
				panic(err)
			}
			user1 = strings.TrimSpace(user1)
			print("password:")
			passwd1, err := terminal.ReadPassword(0)
			if err != nil {
				panic(err)
			}
			password1 := strings.TrimSpace(string(passwd1))
			err = login(&User{user1, password1}, lib)
			if err != nil && err.Error() == loginerror.Error() {
				println("")
				println(loginerror.Error())
				return
			} else if err != nil {
				panic(err)
			}
			visiter = false
			username = user1
			println("")
			fmt.Printf("Welcome %s!\n", username)
		} else {
			print(help)
		}
	case "logout":
		if len(args) == 1 {
			if visiter {
				println("Please login first!")
			} else {
				visiter = true
				username = "visitor"
			}
		} else {
			print(help)
		}
	case "ISBN":
		if len(args) == 2 {
			books, err := querybookbyISBN(args[1], lib)
			if err != nil {
				println(err.Error())
			}
			outputbook(&books)
		} else {
			print(help)
		}
	case "title":
		if len(args) == 2 {
			books, err := querybookbytitle(args[1], lib)
			if err != nil {
				println(err.Error())
			}
			outputbook(&books)
		} else {
			print(help)
		}
	case "author":
		if len(args) == 2 {
			books, err := querybookbyauthor(args[1], lib)
			if err != nil {
				println(err.Error())
			}
			outputbook(&books)
		} else {
			print(help)
		}
	case "return":
		if len(args) == 2 {
			err := returnsinglebook(args[1], username, lib)
			if err != nil {
				println(err.Error())
			}
		} else {
			print(help)
		}
	case "bookid":
		if len(args) == 2 {
			books, err := querysinglebookbyISBN(args[1], lib)
			if err != nil {
				println(err.Error())
			}
			outputsinglebook(&books)
		} else {
			println(help)
		}
	case "borrow":
		if len(args) == 2 {
			if visiter {
				println("Please login first!")
			} else {
				err := rentsinglebook(args[1], username, lib)
				if err != nil {
					println(err.Error())
				}
			}

		} else {
			print(help)
		}

	default:
		print(help)
	}
	return
}
