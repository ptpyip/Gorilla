package repl

import (
	"bufio"
	"fmt"
	"gorilla/ast"
	"gorilla/lexer"
	"gorilla/parser"
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
		p := parser.NewParser(lx)

		prog, ok := p.ParseProgram()
		if !ok {
			printParserErrors(out, prog.Statements, p.Errors)
			continue
		}

		// print statements
		io.WriteString(out, "Parsed Program:\n")
		for _, stmt := range prog.Statements {
			io.WriteString(out, stmt.ToString())
		}
		io.WriteString(out, "\n[END]\n")

		// for tok := lx.GetNextToken(); tok.Type != token.EOF; {
		// 	fmt.Printf("%+v\n", tok)
		// 	tok = lx.GetNextToken()
		// }

	}
}

func printParserErrors(out io.Writer, stmts []ast.StatementNode, errors []string) {
	if out == nil {
		panic("Cannot print Parser errors: nil writer")
	}

	// print statements
	// for _, stmt := range stmts {
	// 	io.WriteString(out, "Parsed statement: "+stmt.ToString())
	// }

	if len(errors) == 0 {
		io.WriteString(out, "No Parser error.")
		return
	}

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
	// io.WriteString(out, "[END]\n")
	// t.FailNow()
}
