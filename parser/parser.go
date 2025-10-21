package parser

import (
	"gorilla/ast"
	"gorilla/lexer"
	"gorilla/token"
)

type Parser struct {
	lx *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token
}

func NewParser(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx}
	// read two tokens, so currentToken and nextToken are both set
	p.loadNextToken()
	p.loadNextToken()

	return p
}

func (p *Parser) loadNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lx.GetNextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			prog.Statements = append(prog.Statements, statement)
		}
		p.loadNextToken()
	}
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	defer p.loadNextToken()

	switch p.currentToken.Type {
	case token.LET:
		if p.nextToken.Type != token.IDENT {
			return nil
		}
		p.loadNextToken()

		if p.nextToken.Type != token.ASSIGN {
			return nil
		}
		statment := &ast.LetStaement{
			Identifier: &ast.IdentifierNode{
				Token: p.currentToken,
				Value: p.currentToken.Literal,
			},
			Value: nil, // set expression
		}
		p.loadNextToken()

		p.skipToSemicolon()
		return statment
	default:
		return nil
	}
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
	}
}
