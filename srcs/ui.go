package main

import "C"
import (
	"bufio"
	"fmt"
	"github.com/modood/table"
	"os"
	"strconv"
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
	helpforroot = `Type "help" can get this help.
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
			fmt.Print(Help())
		}
	case "login":
		if len(args) == 1 {
			fmt.Print("username:")
			user1, err := bufio.NewReader(os.Stdin).ReadString('\n')

			if err != nil {
				panic(err)
			}
			user1 = strings.TrimSpace(user1)
			fmt.Print("password:")
			passwd1, err := terminal.ReadPassword(0)
			if err != nil {
				panic(err)
			}
			password1 := strings.TrimSpace(string(passwd1))
			err = login(&User{user1, password1, 0}, lib)
			if err != nil && err.Error() == loginerror.Error() {
				fmt.Println("")
				fmt.Println(loginerror.Error())
				return
			} else if err != nil {
				panic(err)
			}
			visiter = false
			username = user1
			fmt.Println("")
			fmt.Printf("Welcome %s!\n", username)
		} else {
			fmt.Print(Help())
		}
	case "logout":
		if len(args) == 1 {
			if visiter {
				fmt.Println("Please login first!")
			} else {
				visiter = true
				username = "visitor"
				root = 0
			}
		} else {
			fmt.Print(Help())
		}
	case "ISBN":
		if len(args) == 2 {
			books, err := querybookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			fmt.Print(Help())
		}
	case "title":
		if len(args) == 2 {
			books, err := querybookbytitle(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			fmt.Print(Help())
		}
	case "author":
		if len(args) == 2 {
			books, err := querybookbyauthor(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			fmt.Print(Help())
		}
	case "return":
		if len(args) == 2 {
			err := returnsinglebook(args[1], username, lib)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Print(Help())
		}
	case "bookid":
		if len(args) == 2 {
			books, err := querysinglebookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputsinglebook(&books)
		} else {
			fmt.Print(Help())
		}
	case "borrow":
		if len(args) == 2 {
			if visiter {
				fmt.Println("Please login first!")
			} else {
				err := rentsinglebook(args[1], username, lib)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Successfully borrowed!")
				}

			}

		} else {
			fmt.Print(Help())
		}
	case "add":
		if root == 0 {
			fmt.Print(Help())
		} else {
			if len(args) >= 2 {
				switch args[1] {
				case "user":
					if len(args) != 5 {
						fmt.Print(Help())
					} else {
						u, p := args[2], args[3]
						r, err := strconv.Atoi(args[4])
						if err != nil {
							fmt.Print(Help())
							return
						}
						createuser(&User{u, p, r}, lib)
					}
				case "users":
					if len(args) > 3 {
						fmt.Print(Help())
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/users.csv"
						} else {
							filepath = args[2]
						}
						users, err := readuser(filepath)
						if err != nil {
							//fmt.Print(Help())
							panic(err)
							return
						}
						err = createuser_batch(&users, lib)
						if err != nil {
							panic(err)
						}
					}
				case "book":
					if len(args) != 5 {
						fmt.Print(Help())
					} else {
						t, a, i := args[2], args[3], args[4]
						err := addbook(&Book{t, a, i}, lib)
						if err != nil {
							panic(err)
						}
					}
				case "books":
					if len(args) > 3 {
						fmt.Print(Help())
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/books.csv"
						} else {
							filepath = args[2]
						}
						books, err := readbook(filepath)
						if err != nil {
							//fmt.Print(Help())
							panic(err)
							return
						}
						err = addbook_batch(&books, lib)
						if err != nil {
							panic(err)
						}
					}
				case "sbook":
					if len(args) != 5 {
						fmt.Print(Help())
					} else {
						b, t, i := args[2], args[3], args[4]
						err := addsinglebook(&SingleBook{b, t, i, 1}, lib)
						if err != nil {
							panic(err)
						}
					}
				case "sbooks":
					if len(args) > 3 {
						fmt.Print(Help())
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/sbooks.csv"
						} else {
							filepath = args[2]
						}
						sbooks, err := readsinglebook(filepath)
						if err != nil {
							//fmt.Print(Help())
							panic(err)
							return
						}
						err = addsinglebook_batch(&sbooks, lib)
						if err != nil {
							panic(err)
						}
					}
				default:
					fmt.Print(Help())
				}
			} else {
				fmt.Print(Help())
			}
		}

	default:
		fmt.Print(Help())
	}
	return
}
func outputbook(books *[]Book) {
	if len(*books) > 0 {
		table.OutputA(*books)
	} else {
		fmt.Println("Book not found!")
	}
	return
}
func outputsinglebook(books *[]SingleBook) {
	if len(*books) > 0 {
		table.OutputA(*books)
	} else {
		fmt.Println("Book not found!")
	}
	return
}
func Help() string {
	if root == 1 {
		return helpforroot
	} else {
		return help
	}
}
