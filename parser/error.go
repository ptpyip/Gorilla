package parser

import (
	"fmt"
	"gorilla/token"
)

func (p *Parser) raiseErrorAndPanic(msg string) {
	p.raiseError(msg)
	panic(p.errors)

}

func (p *Parser) raiseError(msg string) {
	p.errors = append(p.errors, "Parser error: "+msg)
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
			p.nextToken.Type,
		),
	)
}
