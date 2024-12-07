package repl

import (
	"bufio"
	"fmt"
	"io"

	"gihub.com/dyxgou/parser/src/evaluator"
	"gihub.com/dyxgou/parser/src/lexer"
	"gihub.com/dyxgou/parser/src/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()

		if p.ErrorsLen() != 0 {
			printParserErrors(out, p.Errors())
		}

		evaluated := evaluator.Eval(program)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "   ")
		io.WriteString(out, err.Error())
		io.WriteString(out, "\n")
	}
}
