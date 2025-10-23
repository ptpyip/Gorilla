package parser

import (
	"fmt"
	"gorilla/ast"
	"gorilla/lexer"
	"gorilla/token"
)

type Parser struct {
	lx *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string
}

func NewParser(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx, errors: []string{}}
	// read two tokens, so currentToken and nextToken are both set
	p.loadNextToken()
	p.loadNextToken()

	return p
}

func (p *Parser) loadNextToken() {
	// if p.nextToken != (token.Token{}) {
	// 	println("Loading statement:", p.nextToken.Literal)
	// }
	p.currentToken = p.nextToken
	p.nextToken = p.lx.GetNextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	ok := true

	prog := &ast.Program{}
	prog.Statements = []ast.StatementNode{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement == nil {
			return prog, !ok
		}

		prog.Statements = append(prog.Statements, statement)
		p.loadNextToken()
	}
	return prog, ok
}

func (p *Parser) parseStatement() ast.StatementNode {
	// defer p.loadNextToken()

	switch p.currentToken.Type {
	case token.LET:
		if p.nextToken.Type != token.IDENT {
			p.raiseNextTokenError(token.IDENT)
			return nil
		}
		p.loadNextToken()

		if p.nextToken.Type != token.ASSIGN {
			p.raiseNextTokenError(token.ASSIGN)
			return nil
		}
		stmt := &ast.LetStatement{
			Identifier: &ast.IdentifierNode{
				Token: p.currentToken,
			},
			Expression: p.parseExpression(),
		}
		p.loadNextToken()

		// handle expressiion

		p.skipToSemicolon()
		return stmt
	default:
		return nil
	}
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
		// println("Skip:", p.currentToken.Literal)
	}
}

func (p *Parser) parseExpression() *ast.ExpressionNode {
	return nil
}

func (p *Parser) reset() {
	p.lx = p.lx.Copy() // Reset lexer to initial state
	p.errors = []string{}

	p.loadNextToken()
	p.loadNextToken()
}

func (p *Parser) GetTokens() []token.Token {
	p.reset()

	tokens := []token.Token{}
	for p.currentToken.Type != token.EOF {
		tokens = append(tokens, p.currentToken)
		p.loadNextToken()
	}
	return tokens
}

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
