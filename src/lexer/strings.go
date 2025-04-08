package lexer

import (
	"fmt"
	"go/token"
)

var symbols = map[byte]byte{
	't':  '\t',
	'n':  '\n',
	'"':  '"',
	'\\': '\\',
}

func getSymbol(ch byte) (byte, error) {
	sym, ok := symbols[ch]
	if !ok {
		return byte(token.ILLEGAL),
			fmt.Errorf("char=\\%s is not a string symbol", string(ch))
	}

	return sym, nil
}
