package parser

import (
	"fmt"
	"gorilla/ast"
	"gorilla/token"
	"strconv"
)

func (p *Parser) raiseErrorAndPanic(msg string) {
	p.raiseError(msg)
	panic(p.errors)
}

func (p *Parser) raiseError(msg string) {
	p.errors = append(p.errors, "Parser error: "+msg)
}

func (p *Parser) raiseParseProgramError() {
	p.raiseError("Failed to parse program. \n[END]")
	// if stmt == nil { // stmt is nil => failed to parse statement

	// 	return
	// }
}

func (p *Parser) raiseParseStatementError(tokenType token.TokenType, stmt ast.StatementNode) {
	if stmt == nil { // stmt is nil => failed to parse statement
		p.raiseError(
			fmt.Sprintf("Failed to parse %s statement.", tokenType),
		)
		return
	}

	// missing semicolon or something
	stmt_str := stmt.ToString()
	if len(stmt_str) > 0 {
		stmt_str = stmt_str[:len(stmt_str)-1]
	}

	p.raiseError(
		fmt.Sprintf("Failed to parse %s statement: %s", tokenType, stmt_str),
	)
}

func (p *Parser) raiseTokenError(expectedTokenType token.TokenType) {
	// println("Current token:", p.currentToken.Literal)

	p.raiseError(
		fmt.Sprintf("Expected %s", expectedTokenType),
	)
}

func (p *Parser) raiseNextTokenError(expectedTokenType token.TokenType) {
	// println("Current token:", p.currentToken.Literal)
	p.raiseError(
		fmt.Sprintf("Expected %s token, got %s token instead",
			expectedTokenType, p.nextToken.Type,
		),
	)
}

func (p *Parser) raiseExpressionError() {
	// println("Current token:", p.currentToken.Literal)
	p.raiseError(
		fmt.Sprintf("Expected expression, got %s token instead",
			p.currentToken.Type,
		),
	)
}

func (p *Parser) raiseBlockStatementError(statements []ast.StatementNode) {
	// print statements
	for i, stmt := range statements {
		p.raiseError("Parsed block statement [" + strconv.Itoa(i) + "]: " + stmt.ToString())

	}

	p.raiseError(
		"current token: " + p.currentToken.Literal +
			" Next token: " + p.nextToken.Literal)

}
