package parser

import (
	"gorilla/ast"
	"gorilla/token"
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
		p.loadNextToken()
		return nil
	}
}

func (p *Parser) parseExpression() (ast.ExpressionNode, bool) {
	left, ok := p.parsePrefix()
	if !ok {
		return nil, !ok
	}

	return left, ok

	// for p.currentToken.Type != token.SEMICOLON {

	// switch p.currentToken.Type {
	// case token.INT:
	// 	expr, err := ast.NewIntegerLiteral(p.currentToken)
	// 	if err != nil {
	// 		p.raiseError(err.Error())
	// 		return nil, !ok
	// 	}

	// 	return expr, ok
	// 	// value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	// 	// if err != nil {
	// 	// 	p.raiseError("Could not parse integer literal: " + p.currentToken.Literal)
	// 	// 	return nil, !ok
	// 	// }

	// 	// return &ast.IntegerLiteral{
	// 	// 	Token: p.currentToken,
	// 	// 	Value: value,
	// 	// }, ok

	// case token.IDENT:
	// 	return &ast.IdentifierNode{p.currentToken}, ok
	// default:
	// 	p.loadNextToken()
	// }

	// }
	// return nil, ok
}

// func (p *Parser) parseIdentifier() (ast.ExpressionNode, bool) {
// 	ok := true

// 	identifier := &ast.IdentifierNode{p.currentToken}
// 	p.loadNextToken()

// 	switch p.currentToken.Type {
// 	case token.PLUS, token.MINUS:

// 		if p.nextToken.Type != p.currentToken.Type {
// 			p.raiseNextTokenError(p.currentToken.Type)
// 			return nil, !ok
// 		}
// 		p.loadNextToken()

// 		if p.nextToken.Type != token.SEMICOLON {
// 			p.raiseNextTokenError(token.SEMICOLON)
// 			return nil, !ok
// 		}
// 		p.loadNextToken()

// 		return &ast.PrefixOperation{
// 			UnaryOperation: ast.UnaryOperation{
// 				Operator: p.currentToken,
// 				Operand:  identifier,
// 			},
// 		}, ok

// 	}
// 	return identifier, ok
// }

func (p *Parser) parsePrefix() (ast.ExpressionNode, bool) {
	ok := true
	switch p.currentToken.Type {

	case token.IDENT:
		return &ast.IdentifierNode{p.currentToken}, ok

	case token.INT:
		if p.nextToken.Type != token.SEMICOLON {
			p.raiseNextTokenError(token.SEMICOLON)
			return nil, !ok
		}

		expr, err := ast.NewIntegerLiteral(p.currentToken)
		if err != nil {
			p.raiseError(err.Error())
			return nil, !ok
		}
		// p.loadNextToken()

		return expr, ok

	case token.BANG, token.MINUS:
		operator := p.currentToken
		p.loadNextToken()

		expr, ok := p.parseExpression()
		if !ok {
			p.raiseError(
				"Could not parse expression after prefix operator " + operator.Literal,
			)
			return nil, !ok
		}

		// optimize negative value
		if operator.Type == token.MINUS && expr.GetTokenType() == token.INT {
			intLit, err := ast.NewIntegerLiteral(
				token.Token{
					Type:    token.INT,
					Literal: "-" + expr.GetTokenLiteral(),
				},
			)
			if err != nil {
				p.raiseError(err.Error())
				return nil, !ok
			}
			return intLit, ok
		}

		return &ast.Prefix{operator, expr}, ok

	default:
		p.raiseError(
			"Could not parse expression given token: " + p.currentToken.Literal,
		)
		return nil, !ok
	}
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
		// println("Skip:", p.currentToken.Literal)
	}
}

// func parseNegativeInteger(intLit *ast.IntegerLiteral) (ast.IntegerLiteral, bool) {
// 	value, err := strconv.ParseInt(token.Literal, 0, 64)
// 	if err != nil {
// 		return nil, false
// 	}
// 	return &ast.IntegerLiteral{Token: token, Value: -value}, true
// }
