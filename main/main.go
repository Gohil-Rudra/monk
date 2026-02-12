package main

import (
	"fmt"
	"monk/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %q This is Monk Programming Language ! \n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
