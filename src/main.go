package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/dyxgou/parser/src/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello %s! This is the Monkey Parser\n", user.Username)
	fmt.Println("Feel free to type in the commands")

	repl.Start(os.Stdin, os.Stdout)
}
