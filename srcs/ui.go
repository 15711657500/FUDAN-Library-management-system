package main

import "C"
import (
	"fmt"
	"github.com/modood/table"
	"os"
	"strconv"
	"strings"
	"syscall"
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
		{"title <title> [title, ...]", "search by title in mutiple keywords", "title math analysis"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
		{"duedate <bookid>", "query the duedate of a borrowed book", "duedate 1"},
		{"topten", "query the top10 bestsellers", "topten"},
	}
	helpforuser = []MyHelp{
		{"quit", "quit", "quit"},
		{"login", "login", "login"},
		{"logout", "logout", "logout"},
		{"ISBN <ISBN>", "search by ISBN", "ISBN abc"},
		{"title <title> [title, ...]", "search by title in mutiple keywords", "title math analysis"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
		{"borrow <bookid>", "borrow the book whose id is <bookid>", "borrow 2"},
		{"return <bookid>", "return the book whose id is <bookid>", "return 2"},
		{"extend <bookid>", "extend the duedate of book whose id is <bookid>", "extend 2"},
		{"changepassword", "change your password", "changepassword"},
		{"list", "query your borrow record", "list"},
		{"duedate <bookid>", "query the duedate of a borrowed book", "duedate 1"},
		{"overdue", "query overdue books of your account", "overdue"},
		{"topten", "query the top10 bestsellers", "topten"},
	}
	helpforroot = []MyHelp{
		{"quit", "quit", "quit"},
		{"login", "login", "login"},
		{"logout", "logout", "logout"},
		{"ISBN <ISBN>", "search by ISBN", "ISBN abc"},
		{"title <title> [title, ...]", "search by title in mutiple keywords", "title math analysis"},
		{"author <author>", "search by author", "author a"},
		{"bookid <ISBN>", "get bookid of books whose ISBN is <ISBN>", "bookid abc"},
		{"duedate <bookid>", "query the duedate of a borrowed book", "duedate 1"},
		{"borrow <bookid>", "borrow the book whose id is <bookid>", "borrow 2"},
		{"return <bookid>", "return the book whose id is <bookid>", "return 2"},
		{"extend <bookid>", "extend the duedate of book whose id is <bookid>", "extend 2"},
		{"changepassword", "change your password", "changepassword"},
		{"add user <username> <password> <root>", "add user", "add user root1 root1 1"},
		{"add users [filepath]", "add user from csv file, default filepath'../data/users.csv'", "add users"},
		{"add book <title> <author> <ISBN>", "add book to booklist", "add book a b c"},
		{"add books [filepath]", "add book to booklist from csv file, default filepath'../data/books.csv'", "add books"},
		{"add sbook <bookid> <ISBN>", "add singlebook", "add a b"},
		{"add sbooks [filepath]", "add singlebook from csv file, default filepath'../data/sbooks.csv'", "add sbooks"},
		{"list [username]", "query borrow record of [username], default yours", "list 18307130001"},
		{"overdue [username]", "query overdue books of the account, default yours", "overdue 18307130001"},
		{"topten", "query the top10 bestsellers", "topten"},
	}
)

const (
	unexpectederror = "Unexpected error!"
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
			user1, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("")
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
				return
			}
			user1 = strings.TrimSpace(user1)
			fmt.Print("password:")
			//passwd1, err := terminal.ReadPassword(0)
			//if err != nil {
			//	panic(err)
			//}
			//password1 := strings.TrimSpace(string(passwd1))
			disableecho(true)
			password1, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			disableecho(false)
			password1 = strings.TrimSpace(password1)
			err = login(&User{user1, password1, 0}, lib)
			if err != nil && err.Error() == loginerror.Error() {
				fmt.Println("")
				fmt.Println(loginerror.Error())
				return
			} else if err != nil {
				fmt.Println("")
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
				return
			}

			fmt.Println("")
			fmt.Printf("Welcome %s!\n", username)
			books1, err := loginduedate(username, lib)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if books1 != nil {
				fmt.Println("You should return the following books as soon as possible!")
				table.OutputA(books1)
			}
			books2, err := loginappoint(username, lib)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if books2 != nil {
				fmt.Println("You can borrow the following books that you have appointed!")
				table.OutputA(books2)
			}
		} else {
			Help()
		}
	case "logout":
		if len(args) == 1 {
			if visitor {
				fmt.Println("Please login first!")
			} else {
				visitor = true
				username = "visitor"
				root = 0
				fmt.Println("Successfully logout!")
			}
		} else {
			Help()
		}
	case "ISBN":
		if len(args) == 2 {
			books, err := querybookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
				return
			}
			outputbook(&books)
		} else {
			Help()
		}
	case "title":
		if len(args) >= 2 {
			books, err := querybookbytitle(args[1:], lib)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
				return
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
				fmt.Println(unexpectederror)
				return
			}
			outputbook(&books)
		} else {
			Help()
		}
	case "return":
		if len(args) == 2 {
			err := returnsinglebook(args[1], username, lib)
			if err != nil && err.Error() == returnnotfound.Error() {
				fmt.Println(returnnotfound.Error())
			} else if err != nil {
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
			} else {
				fmt.Println("Successfully returned!")
			}
		} else {
			Help()
		}
	case "bookid":
		if len(args) == 2 {
			books, err := querysinglebookbyISBN(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
				return
			}
			outputsinglebook(&books)
		} else {
			Help()
		}
	case "borrow":
		if len(args) == 2 {
			if visitor {
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
						err = adduser(&User{u, p, r}, lib)
						if err != nil {
							fmt.Println(err.Error())
							fmt.Println(unexpectederror)
						}
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
						err = adduser_batch(&users, lib)
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
	case "changepassword":
		if visitor {
			Help()
		} else {
			fmt.Println("Current password:")
			disableecho(true)
			curp1, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			disableecho(false)
			curp1 = strings.TrimSpace(curp1)
			check, err := checkpassword(username, curp1, lib)
			fmt.Println("")
			if err != nil {
				//TODO:
				return
			}
			if !check {
				fmt.Println("Wrong password!")
				return
			}
			fmt.Println("Input new password:")
			disableecho(true)
			newp1, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			disableecho(false)
			newp1 = strings.TrimSpace(string(newp1))
			fmt.Println("")
			fmt.Println("Repeat new password:")
			disableecho(true)
			rep1, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			disableecho(false)
			rep1 = strings.TrimSpace(string(rep1))
			fmt.Println("")
			if newp1 != rep1 {
				fmt.Println("Password dismatched!")
			} else {
				err = changepassword(username, newp1, lib)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(unexpectederror)
				} else {
					fmt.Println("Successfully changed password!")
				}

			}
		}
	case "list":
		if visitor {
			Help()
		} else if root == 0 {
			if len(args) > 1 {
				Help()
			} else {
				rents, err := queryrentrecord(username, lib)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(unexpectederror)
				} else {
					outputrent(&rents)
				}
			}
		} else {
			switch len(args) {
			case 1:
				rents, err := queryrentrecord(username, lib)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(unexpectederror)
				} else {
					outputrent(&rents)
				}
			case 2:
				rents, err := queryrentrecord(args[1], lib)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(unexpectederror)
				} else {
					outputrent(&rents)
				}
			default:
				Help()
			}
		}
	case "appoint":
		if len(args) != 2 {
			Help()
		} else if visitor {
			Help()
		} else {
			i, err := appoint(args[1], username, lib)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Successfully appointed!")
				fmt.Printf("There are %d people behind!\n", i)
			}
		}
	case "duedate":
		if len(args) != 2 {
			Help()
		} else {
			duedate, err := queryduedate(args[1], lib)
			if err != nil && err.Error() == booknotfound.Error() {
				fmt.Println(booknotfound.Error())
			} else if err != nil {
				fmt.Println(err.Error())
				fmt.Println(unexpectederror)
			} else {
				fmt.Println(duedate)
			}
		}
	case "overdue":
		if visitor {
			Help()
		} else if len(args) == 1 {
			rents, err := queryoverdue(username, lib)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				outputrent(&rents)
			}
		} else if len(args) == 2 && root > 0 {
			rents, err := queryoverdue(args[1], lib)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				outputrent(&rents)
			}
		} else {
			Help()
		}
	case "topten":
		if len(args) == 1 {
			books, err := topten(lib)
			if err != nil {
				fmt.Println(err)
			} else {
				table.OutputA(books)
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
func outputrent(rents *[]Rent) {
	if len(*rents) > 0 {
		table.OutputA(*rents)
	} else {
		fmt.Println("Borrow record not found!")
	}
}
func Help() {
	fmt.Println(help)
	if root == 1 {
		table.OutputA(helpforroot)
	} else if visitor == true {
		table.OutputA(helpforvisitor)
	} else {
		table.OutputA(helpforuser)
	}
}
func disableecho(b bool) {
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil}
	var ws syscall.WaitStatus
	if b {
		pid, err := syscall.ForkExec(
			"/bin/stty",
			[]string{"stty", "-echo"},
			&attrs)
		if err != nil {
			panic(err)
		}
		_, err = syscall.Wait4(pid, &ws, 0, nil)
		if err != nil {
			panic(err)
		}
	} else {
		attrs := syscall.ProcAttr{
			Dir:   "",
			Env:   []string{},
			Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
			Sys:   nil}
		var ws syscall.WaitStatus
		pid, err := syscall.ForkExec(
			"/bin/stty",
			[]string{"stty", "echo"},
			&attrs)
		if err != nil {
			panic(err)
		}
		_, err = syscall.Wait4(pid, &ws, 0, nil)
		if err != nil {
			panic(err)
		}
	}

}
