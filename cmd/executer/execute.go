package main

import (
	"log"
	"os"

	"github.com/dyxgou/parser/src/repl"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatalf("args expected 1 argument. got=%d", len(args))
	}

	path := args[0]

	file, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	repl.Execute(string(file), os.Stdout)
}
