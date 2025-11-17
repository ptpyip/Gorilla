package repl

import (
	"bufio"
	"fmt"
	"gorilla/ast"
	"gorilla/lexer"
	"gorilla/object"
	"gorilla/parser"
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

func evalExpression(expr ast.ExpressionNode) object.Object {
	switch node := expr.(type) {
	case *ast.BoolLiteral:
		return &object.Bool{node.GetValue()}

	case *ast.IntegerLiteral:
		return &object.Int{node.GetValue()}

	case *ast.Infix:
		left := evalExpression(node.Left)
		right := evalExpression(node.Right)
		switch node.GetOperatorType() {
		case token.EQ:
			return &object.Bool{objIsEqual(left, right)}
		default:
			panic("")
		}

	default:
		return &object.None{}
	}
}

func objIsEqual(left object.Object, right object.Object) bool {
	if left == nil || right == nil {
		panic("Cannot compare nil objects")
	}

	if left.GetType() != right.GetType() {
		return false
	}

	switch left.GetType() {
	case object.BOOL:
		leftBool := left.(*object.Bool)
		rightBool := right.(*object.Bool)
		return leftBool.Value == rightBool.Value

	case object.INT:
		leftInt := left.(*object.Int)
		rightInt := right.(*object.Int)
		return leftInt.Value == rightInt.Value

	case object.NONE:
		return true

	default:
		return false
	}
}
