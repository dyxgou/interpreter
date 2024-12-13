package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dyxgou/parser/src/evaluator"
	"github.com/dyxgou/parser/src/lexer"
	"github.com/dyxgou/parser/src/object"
	"github.com/dyxgou/parser/src/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnviroment()

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

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func Execute(text string, out io.Writer) {
	env := object.NewEnviroment()

	l := lexer.New(text)
	p := parser.New(l)

	program := p.ParseProgram()

	if p.ErrorsLen() != 0 {
		printParserErrors(out, p.Errors())
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "   ")
		io.WriteString(out, err.Error())
		io.WriteString(out, "\n")
	}
}
