package parser

import (
	"gorilla/ast"
	"gorilla/expected"
	"gorilla/lexer"
	"gorilla/token"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	testParseProgram(t, `
		let x = 5;
		let _ = 10;
		let foobar = 838383;
		return 0;
		return x;
		return !x;
		let x = -5;
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
		&expected.LetStatement{"_", expected.NewIntegerLiteral(10)},
		&expected.LetStatement{"foobar", expected.NewIntegerLiteral(838383)},
		&expected.ReturnStatement{expected.NewIntegerLiteral(0)},
		&expected.ReturnStatement{&expected.Identifier{Name: "x"}},
		&expected.ReturnStatement{
			&expected.Prefix{token.BANG, &expected.Identifier{Name: "x"}},
		},
		&expected.LetStatement{"x",
			// &expected.Prefix{token.MINUS, &expected.IntegerLiteral{Value: 5}},
			&expected.IntegerLiteral{-5},
		},
		// &expected.LetStatement{"y",
		// 	&expected.Infix{
		// 		token.PLUS,
		// 		&expected.Identifier{Name: "x"},
		// 		&expected.IntegerLiteral{Value: 5},
		// 	},
		// },
	})

	testParseProgram(t, `
		let x = -5;
		let y = x + 5;
		return -y * -1 + 1;
		let z = x && y;
		return !z || x;
	`, []expected.Node{
		&expected.LetStatement{"x",
			&expected.IntegerLiteral{-5},
		},
		&expected.LetStatement{"y",
			&expected.Infix{
				token.PLUS,
				&expected.Identifier{Name: "x"},
				&expected.IntegerLiteral{Value: 5},
			},
		},
		&expected.ReturnStatement{
			&expected.Infix{
				token.PLUS,
				&expected.Infix{
					token.ASTERISK,
					&expected.Prefix{
						token.MINUS,
						&expected.Identifier{Name: "y"},
					},
					&expected.IntegerLiteral{Value: -1},
				},
				&expected.IntegerLiteral{Value: 1},
			},
		},
		&expected.LetStatement{"z",
			&expected.Infix{
				token.AND,
				&expected.Identifier{Name: "x"},
				&expected.Identifier{Name: "y"},
			},
		},
		&expected.ReturnStatement{
			&expected.Infix{
				token.OR,
				&expected.Prefix{
					token.BANG,
					&expected.Identifier{Name: "z"},
				},
				&expected.Identifier{Name: "x"},
			},
		},
	})

	// testParseProgram(t, `
	// let x = 5;
	// return x;
	// let y = -5;
	// return -y;
	// `, []expected.Node{
	// 	&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
	// 	&expected.ReturnStatement{&expected.Identifier{Name: "x"}},
	// 	// &expected.LetStatement{"y", expected.NewIntegerLiteral(-5)},
	// 	&expected.ReturnStatement{
	// 		&expected.Prefix{token.MINUS, &expected.Identifier{Name: "y"}}},
	// })

}

func testParseProgram(t *testing.T,
	tests string, expectedNodes []expected.Node,
) {
	// vlidate testcases
	num_lines := len(strings.Split(tests, ";")) - 1
	expected_lines := len(expectedNodes)
	if num_lines != expected_lines {
		t.Fatalf("Invalid testcases: num_lines %d != expected_lines %d",
			num_lines, expected_lines,
		)
		return
	}

	// parse input
	lx := lexer.NewLexer(tests)
	p := NewParser(lx)

	prog, ok := p.ParseProgram()
	if !ok {
		errors := p.errors

		tokens := p.GetTokens()
		for i, tok := range tokens {
			t.Logf("token %d: %q", i, tok.Literal)
		}

		// print statements
		for _, stmt := range prog.Statements {
			t.Logf("stmt %d: %q", stmt, stmt.GetTokenLiteral())

		}

		testParserErrors(t, errors)

		return

	} else if prog == nil {
		t.Error("ParseProgram() returned nil")
		testParserErrors(t, p.errors)
		return
	}

	// test program statements
	if len(prog.Statements) != num_lines {
		for _, stmt := range prog.Statements {
			letStmt, _ := stmt.(*ast.LetStatement)
			t.Logf("stmt %d: %q", stmt, letStmt.Identifier.GetName())

		}
		tokens := p.GetTokens()
		for i, tok := range tokens {
			t.Logf("token %d: %q", i, tok.Literal)
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			num_lines, len(prog.Statements),
		)
		return
	}

	for i, expectedNode := range expectedNodes {
		pass := expectedNode.Test(t, prog.Statements[i])
		if !pass {
			testParserErrors(t, p.errors)
			t.FailNow()
			return
		}
	}

}

func testParserErrors(t *testing.T, errors []string) {
	if len(errors) == 0 {
		t.Errorf("Expected parser errors, but got none")
	}

	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

// func testParseStatement(t *testing.T, prog *ast.Program, expectedNodes []ExpectedNode) {
// }

// // type ExpectedNode interface {
// // 	test(t *testing.T, stmt ast.StatementNode) bool
// // }

// type ExpectedPrasedLetStatement struct {
// 	name  string
// 	value string
// }

// type ExpectedPrasedReturnStatment struct {
// 	expression string
// }

// func testPraseLetStatements(t *testing.T,
// 	input string, tests []ExpectedPrasedLetStatement,
// ) {
// 	// vlidate testcases
// 	num_lines := len(strings.Split(input, ";")) - 1
// 	expected_lines := len(tests)
// 	if num_lines != expected_lines {
// 		t.Fatalf("Invalid testcases: num_lines %d != expected_lines %d",
// 			num_lines, expected_lines,
// 		)
// 		return
// 	}

// 	// parse input
// 	lx := lexer.NewLexer(input)
// 	p := NewParser(lx)

// 	prog, ok := p.ParseProgram()
// 	if !ok {
// 		testParserErrors(t, p)
// 		return

// 	} else if prog == nil {
// 		t.Error("ParseProgram() returned nil")
// 		testParserErrors(t, p)
// 		return
// 	}

// 	if len(prog.Statements) != num_lines {
// 		for _, stmt := range prog.Statements {
// 			letStmt, _ := stmt.(*ast.LetStatement)
// 			t.Logf("stmt %d: %q", stmt, letStmt.Identifier.GetName())

// 		}
// 		tokens := p.GetTokens()
// 		for i, tok := range tokens {
// 			t.Logf("token %d: %q", i, tok.Literal)
// 		}
// 		t.Fatalf("program.Statements does not contain %d statements. got=%d",
// 			num_lines, len(prog.Statements),
// 		)
// 		return
// 	}

// 	for i, test := range tests {
// 		if !testPraseLetStatement(t, prog.Statements[i], test) {
// 			return
// 		}
// 	}
// }

// func testPraseLetStatement(t *testing.T,
// 	stmt ast.StatementNode, expected ExpectedPrasedLetStatement,
// ) bool {
// 	letStmt, ok := stmt.(*ast.LetStatement)
// 	if !ok {
// 		t.Errorf("let satetment not found. Got %q token", stmt.GetTokenType())
// 		return false
// 	}

// 	if letStmt.Identifier == nil {
// 		t.Errorf("Invalid Let satement: letStmt.Identifier is nil")
// 		return false
// 	}

// 	if letStmt.Identifier.GetName() != expected.name {
// 		t.Log(letStmt.Expression.GetTokenLiteral())
// 		t.Errorf("letStmt.Identifier.Value not %s. got=%s",
// 			expected.name, letStmt.Identifier.GetName(),
// 		)
// 		return false
// 	}

// 	if letStmt.Expression == nil {
// 		t.Errorf("Invalid Let satement: letStmt.Expression is nil")
// 		return false
// 	} else {
// 		testPraseExpression(t, letStmt.Expression, expected.value)
// 	}

// 	return true
// }

// func TestParseLetStatements(t *testing.T) {
// 	TestParseLetStatements()

// 	testPraseLetStatements(t, `
// 		let x = 5;
// 		let _ = 10;
// 		let foobar = 838383;
// 	`, []ExpectedPrasedLetStatement{
// 		{name: "x", value: "5"},
// 		{name: "_", value: "10"},
// 		{name: "foobar", value: "838383"},
// 	})

// 	testPraseReStatements(t, `
// 		let x = 5;
// 		return 5;
// 	`, []ExpectedPrasedLetStatement{
// 		{name: "x", value: "5"},
// 		{name: "foobar", value: "83"},
// 	})
// 	// testPraseLetStatements(t, `
// 	// 		let = 5;
// 	// 	`, []ExpectedPrasedLetStatement{
// 	// 	{name: "x", value: "5"},
// 	// })

// 	// assertPanic(t, func() {

// 	// })
// 	// testPraseLetStatements(t, `
// 	// 	let x = 5;
// 	// 	let y = 10;
// 	// 	let foobar = 838383;
// 	// `, []ExpectedPrasedLetStatement{
// 	// 	{name: "x"},
// 	// 	{name: "y"},
// 	// 	{name: "foobar"},
// 	// })
// }

// func assertPanic(t *testing.T, f func()) {
// 	t.Helper() // marks this function as a test helper
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Panic is expected, but function did not panic")
// 		}
// 	}()
// 	f()
// }

// func testPraseReturnStatement(t *testing.T,
// 	stmt ast.StatementNode, expected ExpectedPrasedReturnStatment,
// ) bool {
// 	returnStmt, ok := stmt.(*ast.ReturnStatement)
// 	if !ok {
// 		t.Errorf("let satetment not found. Got %q token", stmt.GetTokenType())
// 		return false
// 	}

// 	if returnStmt.ReturnValue == nil {
// 		t.Errorf("Invalid Let satement: letStmt.Identifier is nil")
// 		return false
// 	}

// 	return true
// }

// func testPraseExpression(t *testing.T,
// 	expr ast.ExpressionNode, value string,
// ) bool {
// 	if expr == nil {
// 		t.Errorf("Expression is nil")
// 		return false
// 	}
// 	switch expr.GetTokenType() {
// 	case token.INT:
// 		intLit, ok := expr.(*ast.IntegerLiteral)
// 		if !ok {
// 			t.Errorf("Expected IntegerLiteral. got %T expression", expr.GetTokenType())
// 			return false
// 		}

// 		if intLit.GetTokenLiteral() != value {
// 			t.Errorf("Expected intLit.TokenLiteral = %s. got = %s",
// 				value, intLit.GetTokenLiteral(),
// 			)
// 			return false
// 		}
// 		intValue, _ := strconv.ParseInt(value, 0, 64)
// 		if intLit.GetValue() != intValue {
// 			t.Errorf("intLit.Value not %s. got=%d",
// 				value, intLit.GetValue(),
// 			)
// 			return false
// 		}
// 	default:
// 		t.Errorf("Unexpected expression type %s", expr.GetTokenType())
// 		return false
// 	}
// 	return true
// }
