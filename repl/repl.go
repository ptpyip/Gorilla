package repl

import (
	"bufio"
	"fmt"
	"gorilla/lexer"
	"gorilla/token"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		ok := scanner.Scan()
		if !ok {
			return
		}

		newLine := scanner.Text()
		lx := lexer.NewLexer(newLine)
		for tok := lx.NextToken(); tok.Type != token.EOF; tok = lx.NextToken() {
			fmt.Printf("%+v\n", tok)
		}

	}
}
