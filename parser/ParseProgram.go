package parser

import (
	"gorilla/ast"
	"gorilla/token"
	"strconv"
)

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	ok := true

	prog := &ast.Program{}
	prog.Statements = []ast.StatementNode{}

	for p.currentToken.Type != token.EOF {
		// println("Parsing statement:", p.currentToken.Literal)
		statement := p.parseStatement()
		if statement == nil {
			// for _, msg := range p.errors {
			// 	println(msg)
			// }
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

		identifier := &ast.IdentifierNode{
			Token: p.currentToken,
		}
		p.loadNextToken()

		if p.currentToken.Type != token.ASSIGN {
			p.raiseNextTokenError(token.ASSIGN)
			return nil
		}
		p.loadNextToken()

		expression, ok := p.parseExpression()
		if !ok {
			p.raiseExpressionError()
			return nil
		}

		stmt := &ast.LetStatement{
			Identifier: identifier,
			Expression: expression,
		}
		p.loadNextToken()

		// handle expressiion

		p.skipToSemicolon()
		return stmt

	case token.RETURN:
		if p.nextToken.Type == token.SEMICOLON {
			p.raiseExpressionError()
			return nil
		}
		p.loadNextToken()

		returnValue, ok := p.parseExpression()
		if !ok {
			p.raiseExpressionError()
			return nil
		}
		p.loadNextToken()

		return &ast.ReturnStatement{ReturnValue: returnValue}

	case token.SEMICOLON:
		p.raiseError("Unexpected semicolon")
		return nil
	default:
		return nil
	}
}

func (p *Parser) parseExpression() (ast.ExpressionNode, bool) {
	ok := true

	for p.currentToken.Type != token.SEMICOLON {
		if p.currentToken.Type == token.EOF {
			return nil, !ok
		}

		switch p.currentToken.Type {
		case token.INT:
			value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
			if err != nil {
				p.raiseError("Could not parse integer literal: " + p.currentToken.Literal)
				return nil, !ok
			}

			return &ast.IntegerLiteral{
				Token: p.currentToken,
				Value: value,
			}, ok

		case token.IDENT:
			return &ast.IdentifierNode{
				Token: p.currentToken,
			}, ok
		default:
			p.loadNextToken()
		}

	}
	return nil, ok
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
		// println("Skip:", p.currentToken.Literal)
	}
}
