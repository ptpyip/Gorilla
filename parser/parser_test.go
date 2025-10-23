package parser

import (
	"gorilla/ast"
	"gorilla/lexer"
	"strings"
	"testing"
)

type ExpectedPrasedLetStatement struct {
	name string
	// expression string
}

func testPraseLetStatements(t *testing.T,
	input string, tests []ExpectedPrasedLetStatement,
) {
	// vlidate testcases
	num_lines := len(strings.Split(input, ";")) - 1
	expected_lines := len(tests)
	if num_lines != expected_lines {
		t.Fatalf("Invalid testcases: num_lines %d != expected_lines %d",
			num_lines, expected_lines,
		)
		return
	}

	// parse input
	lx := lexer.NewLexer(input)
	p := NewParser(lx)

	prog, ok := p.ParseProgram()
	if !ok {
		testParserErrors(t, p)
		return

	} else if prog == nil {
		t.Error("ParseProgram() returned nil")
		testParserErrors(t, p)
		return
	}

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

	for i, test := range tests {
		if !testPraseLetStatement(t, prog.Statements[i], test.name) {
			return
		}
	}
}

func testPraseLetStatement(t *testing.T,
	stmt ast.StatementNode, name string,
) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("let satetment not found. Got %q token", stmt.GetTokenType())
		return false
	}

	if letStmt.Identifier == nil {
		t.Errorf("Invalid Let satement: letStmt.Identifier is nil")
		return false
	}

	if letStmt.Identifier.GetName() != name {
		t.Errorf("letStmt.Identifier.Value not %s. got=%s",
			name, letStmt.Identifier.GetName(),
		)
		return false
	}

	// if letStmt.Expression == nil {
	// 	t.Errorf("Invalid Let satement: letStmt.Expression is nil")
	// 	return false
	// }

	return true
}

func testParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		t.Errorf("Expected parser errors, but got none")
	}

	for _, msg := range p.errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func TestParseLetStatements(t *testing.T) {

	testPraseLetStatements(t, `
		let x = 5;
		let _ = 10;
		let foobar = 838383;
	`, []ExpectedPrasedLetStatement{
		{name: "x"},
		{name: "_"},
		{name: "foobar"},
	})

	testPraseLetStatements(t, `
			let = 5;
		`, []ExpectedPrasedLetStatement{
		{name: "x"},
	})

	// assertPanic(t, func() {

	// })
	// testPraseLetStatements(t, `
	// 	let x = 5;
	// 	let y = 10;
	// 	let foobar = 838383;
	// `, []ExpectedPrasedLetStatement{
	// 	{name: "x"},
	// 	{name: "y"},
	// 	{name: "foobar"},
	// })
}

func assertPanic(t *testing.T, f func()) {
	t.Helper() // marks this function as a test helper
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panic is expected, but function did not panic")
		}
	}()
	f()
}
