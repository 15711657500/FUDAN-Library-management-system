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

const (
	help = `Type "help" can get this help.
quit-quit
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

	default:
		print(help)
	}
	return
}

func outputbook() {

}
