package parser

import (
	"gorilla/ast"
	"gorilla/expected"
	"gorilla/lexer"
	"gorilla/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	testParseProgram(t, `
		let x = 5;
		let _ = 10;
		let foobar = 838383;
		let y = x;
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
		&expected.LetStatement{"_", expected.NewIntegerLiteral(10)},
		&expected.LetStatement{"foobar", expected.NewIntegerLiteral(838383)},
		&expected.LetStatement{"y", &expected.Identifier{Name: "x"}},
	})
}

func TestReturnStatements(t *testing.T) {
	testParseProgram(t, `
		let x = 5;
		return 0;
		return x;
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
		&expected.ReturnStatement{expected.NewIntegerLiteral(0)},
		&expected.ReturnStatement{&expected.Identifier{Name: "x"}},
	})
}

func TestPrefixExpressions(t *testing.T) {
	testParseProgram(t, `
		let x = -5;
		let y = True;
		let z = !y;
		return -x;
		return !y;
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(-5)},
		&expected.LetStatement{"y", expected.NewBoolLiteral(true)},
		&expected.LetStatement{"z",
			&expected.Prefix{
				token.BANG,
				&expected.Identifier{Name: "y"},
			},
		},
		&expected.ReturnStatement{
			&expected.Prefix{
				token.MINUS,
				&expected.Identifier{Name: "x"},
			},
		},
		&expected.ReturnStatement{
			&expected.Prefix{token.BANG,
				&expected.Identifier{Name: "y"},
			},
		},
	})
}

func TestInfixExpressions(t *testing.T) {

	testParseProgram(t, `
		return 3 < 5 == True;
		let a = 5 < 4 != 3 > 4;
		return 3 + 4 * 5 == 3 * 1 + 4 * 5;
	`, []expected.Node{
		&expected.ReturnStatement{
			&expected.Infix{
				token.EQ,
				&expected.Infix{
					token.LT,
					&expected.IntegerLiteral{3},
					&expected.IntegerLiteral{5},
				},
				&expected.BoolLiteral{true},
			},
		},
		&expected.LetStatement{"a",
			&expected.Infix{
				token.NOT_EQ,
				&expected.Infix{
					token.LT,
					&expected.IntegerLiteral{Value: 5},
					&expected.IntegerLiteral{Value: 4},
				},
				&expected.Infix{
					token.GT,
					&expected.IntegerLiteral{Value: 3},
					&expected.IntegerLiteral{Value: 4},
				},
			},
		},
		&expected.ReturnStatement{
			&expected.Infix{
				token.EQ,
				&expected.Infix{
					token.PLUS,
					&expected.IntegerLiteral{Value: 3},
					&expected.Infix{
						token.ASTERISK,
						&expected.IntegerLiteral{Value: 4},
						&expected.IntegerLiteral{Value: 5},
					},
				},
				&expected.Infix{
					token.PLUS,
					&expected.Infix{
						token.ASTERISK,
						&expected.IntegerLiteral{Value: 3},
						&expected.IntegerLiteral{Value: 1},
					},
					&expected.Infix{
						token.ASTERISK,
						&expected.IntegerLiteral{Value: 4},
						&expected.IntegerLiteral{Value: 5},
					},
				},
			},
		},
	})
}

func TestGroupedExpressions(t *testing.T) {
	testParseProgram(t, `
		let x = (5 + 5) * 2;
		return ((x / 2) * 2) == 20;
	`, []expected.Node{
		&expected.LetStatement{"x",
			&expected.Infix{
				token.ASTERISK,
				&expected.Infix{
					token.PLUS,
					&expected.IntegerLiteral{5},
					&expected.IntegerLiteral{5},
				},
				&expected.IntegerLiteral{2},
			},
		},
		&expected.ReturnStatement{
			&expected.Infix{
				token.EQ,
				&expected.Infix{
					token.ASTERISK,
					&expected.Infix{
						token.SLASH,
						&expected.Identifier{Name: "x"},
						&expected.IntegerLiteral{2},
					},
					&expected.IntegerLiteral{2},
				},
				&expected.IntegerLiteral{20},
			},
		},
	})
}

func TestBlockStatements(t *testing.T) {
	testParseProgram(t, `

		let x = 5;
		{
			{
				let a = 5 < 4 != 3 > 4;
				return !x; 
			}
			return (x > 1) && !a;
		}
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
		expected.NewBlockStatement(
			expected.NewBlockStatement(
				&expected.LetStatement{"a",
					&expected.Infix{
						token.NOT_EQ,
						&expected.Infix{
							token.LT,
							&expected.IntegerLiteral{Value: 5},
							&expected.IntegerLiteral{Value: 4},
						},
						&expected.Infix{
							token.GT,
							&expected.IntegerLiteral{Value: 3},
							&expected.IntegerLiteral{Value: 4},
						},
					},
				},
				&expected.ReturnStatement{
					&expected.Prefix{
						token.BANG,
						&expected.Identifier{Name: "x"},
					},
				},
			),
			&expected.ReturnStatement{
				&expected.Infix{
					token.AND,
					&expected.Infix{
						token.GT,
						&expected.Identifier{Name: "x"},
						&expected.IntegerLiteral{Value: 1},
					},
					&expected.Prefix{
						token.BANG,
						&expected.Identifier{Name: "a"},
					},
				},
			},
		),
	})
}

func TestIfElseStatements(t *testing.T) {
	testParseProgram(t, `
	if (x == y) { 
		let x = 5; 
	} else if (x > y) {
		let x = 6;
	} else {
	 	{ let y = 10; }
		let x = 10;
	}
	`, []expected.Node{
		&expected.IfStatement{
			&expected.Infix{
				token.EQ,
				&expected.Identifier{Name: "x"},
				&expected.Identifier{Name: "y"},
			},
			expected.NewBlockStatement(
				&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
			),
			expected.NewElseIfStatement(
				&expected.Infix{
					token.GT,
					&expected.Identifier{Name: "x"},
					&expected.Identifier{Name: "y"},
				},
				expected.NewBlockStatement(
					&expected.LetStatement{"x", expected.NewIntegerLiteral(6)},
				),
				&expected.ElseStatement{
					Statement: expected.NewBlockStatement(
						expected.NewBlockStatement(
							&expected.LetStatement{"y", expected.NewIntegerLiteral(10)},
						),
						&expected.LetStatement{"x", expected.NewIntegerLiteral(10)},
					),
				},
			),
		},
	})
}

func TestIfElseExpressions(t *testing.T) {
	testParseProgram(t, `
	let x = True;
	return 1 if x else 2;
	if (x == y) { 
		return 1 if x else 2;
	} else {
	 	{ let y = 1 if x else 2; }
		return 1 if y else 2;
	}
		
	`, []expected.Node{
		&expected.LetStatement{"x",
			expected.NewBoolLiteral(true),
		},
		&expected.ReturnStatement{
			&expected.Trinary{
				expected.NewIntegerLiteral(1),
				&expected.Identifier{Name: "x"},
				expected.NewIntegerLiteral(2),
			},
		},
		&expected.IfStatement{
			&expected.Infix{
				token.EQ,
				&expected.Identifier{Name: "x"},
				&expected.Identifier{Name: "y"},
			},
			expected.NewBlockStatement(
				&expected.ReturnStatement{
					&expected.Trinary{
						expected.NewIntegerLiteral(1),
						&expected.Identifier{Name: "x"},
						expected.NewIntegerLiteral(2),
					},
				},
			),
			&expected.ElseStatement{
				expected.NewBlockStatement(
					expected.NewBlockStatement(
						&expected.LetStatement{"y",
							&expected.Trinary{
								expected.NewIntegerLiteral(1),
								&expected.Identifier{Name: "x"},
								expected.NewIntegerLiteral(2),
							},
						},
					),
					&expected.ReturnStatement{
						&expected.Trinary{
							expected.NewIntegerLiteral(1),
							&expected.Identifier{Name: "y"},
							expected.NewIntegerLiteral(2),
						},
					},
				),
			},
		},
	})

}

func TestFunctionDeclarations(t *testing.T) {
	testParseProgram(t, `
		let add = fn(x, y) {
			return x + y;
		};
	`, []expected.Node{
		&expected.LetStatement{"add",
			&expected.SkipNode{},
			// &expected.FunctionDeclaration{
			// 	Name: "add",
			// 	Parameters: []expected.Node{
			// 		&expected.Identifier{Name: "x"},
			// 		&expected.Identifier{Name: "y"},
			// 	},
			// 	Body: &expected.BlockStatement{
			// 		Statements: []expected.Node{
			// 			&expected.ReturnStatement{
			// 				&expected.Infix{
			// 					token.PLUS,
			// 					&expected.Identifier{Name: "x"},
			// 					&expected.Identifier{Name: "y"},
			// 				},
			// 			},
			// 		}},
			// },
		},
	})
}

func TestParser(t *testing.T) {

	testParseProgram(t, `
	let x = 5;
	`, []expected.Node{
		&expected.LetStatement{"x", expected.NewIntegerLiteral(5)},
	})

}

func testParseProgram(t *testing.T,
	testInput string, expectedNodes []expected.Node,
) {
	// parse input
	lx := lexer.NewLexer(testInput)
	p := NewParser(lx)

	prog, ok := p.ParseProgram()
	if !ok {
		raiseParserErrors(t, prog.Statements, p.errors)
		return
	}

	// test program statements
	num_nodes := len(expectedNodes)
	if len(prog.Statements) != num_nodes {
		raiseParserErrors(t, prog.Statements, p.errors)
		raiseExpectedNotEqParsedError(t, prog, expectedNodes)
		return
	}

	for i, expectedNode := range expectedNodes {
		pass := expectedNode.Test(t, prog.Statements[i])
		if !pass {
			raiseParserErrors(t, prog.Statements[:i+1], p.errors)
			// t.FailNow()
			return
		}
	}

}

func raiseParserErrors(t *testing.T, stmts []ast.StatementNode, errors []string) {
	// print statements
	for _, stmt := range stmts {
		t.Logf("Parsed statement: %q", stmt.ToString())
	}

	if len(errors) == 0 {
		t.Errorf("No Parser error.")
		return
	}

	for _, msg := range errors {
		t.Error(msg)
	}
	// t.FailNow()
}

func raiseExpectedNotEqParsedError(t *testing.T,
	prog *ast.Program, expectedNodes []expected.Node,
) {
	t.Errorf("Test: Expected %d statement(s), but got %d parsed.",
		len(expectedNodes), len(prog.Statements),
	)
}

// func testParseProgram(t *testing.T,
// 	testInput string, expectedNodes []expected.Node,
// ) {
// 	// parse input
// 	lx := lexer.NewLexer(testInput)
// 	p := NewParser(lx)

// 	prog, ok := p.ParseProgram()
// 	if !ok {
// 		// tokens := p.GetTokens()
// 		// for i, tok := range tokens {
// 		// 	t.Logf("token %d: %q", i, tok.Literal)

// 		raiseParserErrors(t, prog.Statements, p.errors)
// 		return

// 	}

// 	// if prog == nil {
// 	// 	t.Error("ParseProgram() returned nil")
// 	// 	raiseParserErrors(t, prog.Statements, p.errors)
// 	// 	return
// 	// }

// 	// test program statements
// 	num_nodes := len(expectedNodes)
// 	if len(prog.Statements) != num_nodes {
// 		raiseParserErrors(t, prog.Statements, p.errors)
// 		raiseExpectedNotEqParsedError(t, prog, expectedNodes)
// 		return
// 	}

// 	for i, expectedNode := range expectedNodes {
// 		pass := expectedNode.Test(t, prog.Statements[i])
// 		if !pass {
// 			raiseParserErrors(t, prog.Statements[:i+1], p.errors)
// 			t.FailNow()
// 			return
// 		}
// 	}

// 	// if len(prog.Statements) < num_nodes {
// 	// 	t.Log("len(prog.Statements) < num_nodes")
// 	// 	t.Logf("len(prog.Statements) = %d", len(prog.Statements))
// 	// 	t.Logf("num_nodes = %d", num_nodes)
// 	// 	raiseParserErrors(t, p.errors)
// 	// } else if len(prog.Statements) > num_nodes {
// 	// 	for _, stmt := range prog.Statements {
// 	// 		letStmt, _ := stmt.(*ast.LetStatement)
// 	// 		t.Logf("stmt %d: %q", stmt, letStmt.Identifier.GetName())

// 	// 	}
// 	// 	tokens := p.GetTokens()
// 	// 	for i, tok := range tokens {
// 	// 		t.Logf("token %d: %q", i, tok.Literal)
// 	// 	}
// 	// 	t.Fatalf("program.Statements does not contain %d statements. got=%d",
// 	// 		num_nodes, len(prog.Statements),
// 	// 	)
// 	// 	return
// 	// }

// 	// for i, expectedNode := range expectedNodes {
// 	// 	pass := expectedNode.Test(t, prog.Statements[i])
// 	// 	if !pass {
// 	// 		raiseParserErrors(t, p.errors)
// 	// 		t.FailNow()
// 	// 		return
// 	// 	}
// 	// }

// }

// func raiseParserErrors(t *testing.T, stmts []ast.StatementNode, errors []string) {
// 	// print statements
// 	for _, stmt := range stmts {
// 		t.Logf("Parsed statement: %q", stmt.ToString())
// 	}

// 	if len(errors) == 0 {
// 		t.Errorf("No Parser error.")
// 		return
// 	}

// 	for _, msg := range errors {
// 		t.Error(msg)
// 	}
// 	// t.FailNow()
// }

// func raiseExpectedNotEqParsedError(t *testing.T,
// 	prog *ast.Program, expectedNodes []expected.Node,
// ) {
// 	t.Errorf("Test: Expected %d statement(s), but got %d parsed.",
// 		len(expectedNodes), len(prog.Statements),
// 	)
// }

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
// 	num_nodes := len(strings.Split(input, ";")) - 1
// 	expected_lines := len(tests)
// 	if num_nodes != expected_lines {
// 		t.Fatalf("Invalid testcases: num_nodes %d != expected_lines %d",
// 			num_nodes, expected_lines,
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

// 	if len(prog.Statements) != num_nodes {
// 		for _, stmt := range prog.Statements {
// 			letStmt, _ := stmt.(*ast.LetStatement)
// 			t.Logf("stmt %d: %q", stmt, letStmt.Identifier.GetName())

// 		}
// 		tokens := p.GetTokens()
// 		for i, tok := range tokens {
// 			t.Logf("token %d: %q", i, tok.Literal)
// 		}
// 		t.Fatalf("program.Statements does not contain %d statements. got=%d",
// 			num_nodes, len(prog.Statements),
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
