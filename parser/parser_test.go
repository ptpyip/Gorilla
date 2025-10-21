package parser

import (
	"gorilla/lexer"
	"strings"
	"testing"
)

func testPraseLetStatement(t *testing.T, input string) bool {
	lx := lexer.NewLexer(input)
	p := NewParser(lx)

	prog := p.ParseProgram()
	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	num_lines := len(strings.Split(input, ";"))
	if len(prog.Statements) != num_lines {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			num_lines, len(prog.Statements),
		)
		return false
	}
}
