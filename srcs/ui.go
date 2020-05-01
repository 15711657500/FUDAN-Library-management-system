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

type MyHelp struct {
	Command     string
	Description string
	Example     string
}

var (
	help           = `Type "help" can get this help.`
	helpforvisitor = []MyHelp{
		{"quit", "quit", "quit"},
		{"login", "login", "login"},
		{"logout", "logout", "logout"},
		{"ISBN <ISBN>", "search by ISBN", "ISBN abc"},
		{"title <title>", "search by title", "title abc"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
	}
	helpforuser = []MyHelp{
		{"quit", "quit", "quit"},
		{"login", "login", "login"},
		{"logout", "logout", "logout"},
		{"ISBN <ISBN>", "search by ISBN", "ISBN abc"},
		{"title <title>", "search by title", "title abc"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
		{"borrow <bookid>", "borrow the book whose id is <bookid>", "borrow 2"},
		{"return <bookid>", "return the book whose id is <bookid>", "return 2"},
		{"extend <bookid>", "extend the duedate of book whose id is <bookid>", "extend 2"},
	}
	helpforroot = []MyHelp{
		{"quit", "quit", "quit"},
		{"login", "login", "login"},
		{"logout", "logout", "logout"},
		{"ISBN <ISBN>", "search by ISBN", "ISBN abc"},
		{"title <title>", "search by title", "title abc"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
		{"borrow <bookid>", "borrow the book whose id is <bookid>", "borrow 2"},
		{"return <bookid>", "return the book whose id is <bookid>", "return 2"},
		{"extend <bookid>", "extend the duedate of book whose id is <bookid>", "extend 2"},
		{"add user <username> <password> <root>", "add user", "add user root1 root1 1"},
		{"add users [filepath]", "add user from csv file, default filepath'../data/users.csv'", "add users"},
		{"add book <title> <author> <ISBN>", "add book to booklist", "add book a b c"},
		{"add books [filepath]", "add book to booklist from csv file, default filepath'../data/books.csv'", "add books"},
		{"add sbook <bookid> <ISBN>", "", ""},
	}
)

func handleinput(input string, lib *Library) {
	args := strings.Split(input, " ")
	switch args[0] {
	case "quit":
		if len(args) == 1 {
			done = true
		} else {
			Help()
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
			Help()
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
			Help()
		}
	case "ISBN":
		if len(args) == 2 {
			books, err := querybookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			Help()
		}
	case "title":
		if len(args) == 2 {
			books, err := querybookbytitle(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			Help()
		}
	case "author":
		if len(args) == 2 {
			books, err := querybookbyauthor(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputbook(&books)
		} else {
			Help()
		}
	case "return":
		if len(args) == 2 {
			err := returnsinglebook(args[1], username, lib)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			Help()
		}
	case "bookid":
		if len(args) == 2 {
			books, err := querysinglebookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			}
			outputsinglebook(&books)
		} else {
			Help()
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
			Help()
		}
	case "add":
		if root == 0 {
			Help()
		} else {
			if len(args) >= 2 {
				switch args[1] {
				case "user":
					if len(args) != 5 {
						Help()
					} else {
						u, p := args[2], args[3]
						r, err := strconv.Atoi(args[4])
						if err != nil {
							Help()
							return
						}
						createuser(&User{u, p, r}, lib)
					}
				case "users":
					if len(args) > 3 {
						Help()
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/users.csv"
						} else {
							filepath = args[2]
						}
						users, err := readuser(filepath)
						if err != nil {
							//Help()
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
						Help()
					} else {
						t, a, i := args[2], args[3], args[4]
						err := addbook(&Book{t, a, i}, lib)
						if err != nil {
							panic(err)
						}
					}
				case "books":
					if len(args) > 3 {
						Help()
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/books.csv"
						} else {
							filepath = args[2]
						}
						books, err := readbook(filepath)
						if err != nil {
							//Help()
							panic(err)
							return
						}
						err = addbook_batch(&books, lib)
						if err != nil {
							panic(err)
						}
					}
				case "sbook":
					if len(args) != 4 {
						Help()
					} else {
						b, i := args[2], args[3]
						err := addsinglebook(&SingleBook{b, "", i, 1}, lib)
						if err != nil {
							panic(err)
						}
					}
				case "sbooks":
					if len(args) > 3 {
						Help()
					} else {
						var filepath string
						if len(args) == 2 {
							filepath = "../data/sbooks.csv"
						} else {
							filepath = args[2]
						}
						sbooks, err := readsinglebook(filepath)
						if err != nil {
							//Help()
							panic(err)
							return
						}
						err = addsinglebook_batch(&sbooks, lib)
						if err != nil {
							panic(err)
						}
					}
				default:
					Help()
				}
			} else {
				Help()
			}
		}
	case "extend":
		if len(args) != 2 {
			Help()
		} else {
			bookid := args[1]
			err := extend(bookid, username, lib)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Successfully extend!")
			}
		}
	default:
		Help()
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
func Help() {
	fmt.Println(help)
	if root == 1 {
		table.OutputA(helpforroot)
	} else if visiter == true {
		table.OutputA(helpforvisitor)
	} else {
		table.OutputA(helpforuser)
	}
}
