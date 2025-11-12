package parser

import (
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

func (p *Parser) GetTokens() []token.Token {
	p.reset()

	tokens := []token.Token{}
	for p.currentToken.Type != token.EOF {
		tokens = append(tokens, p.currentToken)
		p.loadNextToken()
	}
	return tokens
}

func (p *Parser) reset() {
	p.lx = p.lx.Copy() // Reset lexer to initial state
	// p.errors = []string{}

	p.loadNextToken()
	p.loadNextToken()
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
	}
}
