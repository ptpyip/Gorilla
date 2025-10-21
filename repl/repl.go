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
		fmt.Print(PROMPT)

		ok := scanner.Scan()
		if !ok {
			return
		}

		newLine := scanner.Text()
		lx := lexer.NewLexer(newLine)
		for tok := lx.NextToken(); tok.Type != token.EOF; {
			fmt.Printf("%+v\n", tok)
			tok = lx.NextToken()
		}

	}
}
