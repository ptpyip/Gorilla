package repl

import (
	"gorilla/ast"
	"gorilla/lexer"
	"gorilla/object"
	"gorilla/parser"
	"testing"
)

func TestEvalProgram(t *testing.T) {
	testEvalProgram(t, `
		return 5;
	`, []object.Object{
		&object.Int{Value: 5},
	})
}

func testEvalProgram(t *testing.T,
	testInput string, expectedObjects []object.Object,
) {
	lx := lexer.NewLexer(testInput)
	p := parser.NewParser(lx)

	prog, ok := p.ParseProgram()
	if !ok {
		// raiseParserErrors(t, prog.Statements, p.Errors)
		t.Error("ParseProgram() returned nil")
		return
	}

	// test program statements
	// num_nodes := len(expectedNodes)
	// if len(prog.Statements) != num_nodes {
	// 	raiseParserErrors(t, prog.Statements, p.Errors)
	// 	raiseExpectedNotEqParsedError(t, prog, expectedNodes)
	// 	return
	// }

	for i, stmt := range prog.Statements {
		testEvalStatement(t, stmt, expectedObjects[i])
	}

}

func testEvalStatement(t *testing.T,
	stmt ast.StatementNode,
	expectedObj object.Object,
) {
	switch stmt := stmt.(type) {
	case *ast.ReturnStatement:
		if stmt.ReturnValue == nil {
			t.Errorf("Invalid Return statement: ReturnValue is nil")
			return
		}
		obj := evalExpression(stmt.ReturnValue)
		if !objIsEqual(obj, expectedObj) {
			t.Errorf("Return statement error: got=%s, expected=%s",
				obj.Inspect(), expectedObj.Inspect(),
			)
		}

	default:
		return
	}

}
