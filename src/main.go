package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	b, err := os.ReadFile("./examples/00.lang")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
