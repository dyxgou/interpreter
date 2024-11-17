package repl

import (
	"bufio"
	"fmt"
	"io"

	"gihub.com/dyxgou/parser/src/lexer"
	"gihub.com/dyxgou/parser/src/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		text := scanner.Text()

		l := lexer.New(text)

		for tok := l.NextToken(); tok.Kind != token.EOF; tok = l.NextToken() {
			fmt.Printf("tok: %v\n", tok)
		}
	}
}
