package parser

import (
	"gorilla/ast"
	"gorilla/parser/precedences"
	"gorilla/token"
)

func (p *Parser) parseExpression(parentPrecedence int) (ast.ExpressionNode, bool) {
	var expr ast.ExpressionNode
	switch p.currentToken.Type {
	case token.IDENT:
		expr = &ast.IdentifierExpression{p.currentToken}

	case token.TRUE, token.FALSE:
		expr = &ast.BoolLiteral{p.currentToken}

	case token.INT:
		intLit, err := ast.NewIntegerLiteral(p.currentToken)
		if err != nil {
			p.raiseError(err.Error())
			return nil, false
		}
		expr = intLit

	case token.FUNCTION:
		// fn_definition
		if p.nextToken.Type != token.LPAREN {
			p.raiseNextTokenError(token.LPAREN)
			return nil, false
		}
		p.loadNextToken()

		signiture := []ast.IdentifierExpression{}
		if p.nextToken.Type == token.RPAREN {
			p.loadNextToken()
		} // skip for fn(), ie empty signature

		for p.currentToken.Type != token.RPAREN {
			if p.nextToken.Type != token.IDENT {
				p.raiseNextTokenError(token.IDENT)
				return nil, false
			}
			p.loadNextToken()

			signiture = append(signiture, ast.IdentifierExpression{p.currentToken})
			p.loadNextToken()

			// if p.currentToken.Type == token.COMMA {
			// 	p.loadNextToken()
			// }
		}

		body, ok := p.parseBlockStatement()
		if !ok {
			p.raiseError("Could not parse block statement")
			return nil, false
		} else if body == nil {
			p.raiseError("Invalid function definition: body is nil")
			return nil, false
		}

		// print("After parsing body: ", p.currentToken.Literal) // epxected to be after '}
		expr = &ast.FunctionDefinition{signiture, body}

	case token.LPAREN, token.BANG, token.MINUS:
		prefix, ok := p.parsePrefix()
		if !ok {
			return nil, false
		}
		expr = prefix

	// case token.RBRACE:
	// 	p.raiseError()

	default:
		p.raiseErrorAndPanic(
			"Unexpected token for parseExpression: " + string(p.currentToken.Type),
		)
	}

	if p.nextToken.Type == token.SEMICOLON {
		return expr, true
	}

	// Next token is infix operator
	for p.getNextPrecedence() > parentPrecedence {
		expr = p.parseInfix(expr)
		if expr == nil {
			p.raiseError("Could not parse infix expression")
			return nil, false
		}

		if p.nextToken.Type == token.SEMICOLON {
			break
		}
	}

	if p.nextToken.Type == token.IF {
		trinary, ok := p.parseIfElseExpression(expr)
		if !ok {
			p.raiseError("Could not parse if-else expression")
			return nil, false
		}
		return trinary, true
	}

	return expr, true
}

func (p *Parser) parsePrefix() (ast.ExpressionNode, bool) {
	// println("parsePrefix on token: " + p.currentToken.Literal)
	precedence := precedences.PREFIX
	switch p.currentToken.Type {
	case token.BANG, token.MINUS:
		operator := p.currentToken
		p.loadNextToken()

		operand, ok := p.parseExpression(precedence)
		if !ok {
			p.raiseError(
				"Could not parse expression after prefix operator " + operator.Literal,
			)
			return nil, !ok
		}

		// optimize negative value
		if operator.Type == token.MINUS && operand.GetTokenType() == token.INT {
			intLit, err := ast.NewIntegerLiteral(
				token.Token{
					Type:    token.INT,
					Literal: "-" + operand.GetTokenLiteral(),
				},
			)
			if err != nil {
				p.raiseError(err.Error())
				return nil, !ok
			}
			return intLit, ok
		}

		return &ast.Prefix{operator, operand}, ok

	case token.LPAREN:
		p.loadNextToken()

		inner_expression, ok := p.parseExpression(precedences.LOWEST)
		if !ok {
			p.raiseError("Could not parse LPAREN expression")
			return nil, !ok
		}

		if p.nextToken.Type != token.RPAREN {
			p.raiseNextTokenError(token.RPAREN)
			return nil, !ok
		}
		p.loadNextToken()

		return inner_expression, ok

	// case token.IDENT:
	// 	return &ast.IdentifierNode{p.currentToken}, ok

	// case token.INT:
	// 	if p.nextToken.Type != token.SEMICOLON {
	// 		p.raiseNextTokenError(token.SEMICOLON)
	// 		return nil, !ok
	// 	}

	// 	expr, err := ast.NewIntegerLiteral(p.currentToken)
	// 	if err != nil {
	// 		p.raiseError(err.Error())
	// 		return nil, !ok
	// 	}
	// 	// p.loadNextToken()

	// 	return expr, ok
	default:
		p.raiseError("Unexpected Prefix operator " + string(p.currentToken.Type))
		return nil, false
	}
}

func (p *Parser) parseInfix(left ast.ExpressionNode) ast.ExpressionNode {
	// println("parsePrefix on token: " + p.currentToken.Literal)
	if left == nil {
		p.raiseError("Left operand is nil")
		return nil
	}
	p.loadNextToken()

	operator := p.currentToken
	precedence := p.getCurrentPrecedence()
	p.loadNextToken()

	// println("Infix operator: " + operator.Literal)
	// println("\twith operator precedence = ", precedence)

	ok := true
	var right ast.ExpressionNode
	switch operator.Type {
	// algrithmic operators
	case token.PLUS, token.MINUS:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" algorithmic error " + operator.Literal)
			return nil
		}

	case token.ASTERISK, token.SLASH:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" algorithmic error " + operator.Literal)
			return nil

		}

	// comparison operators
	case token.EQ, token.NOT_EQ, token.LE, token.GE:
		right, ok = p.parseExpression(precedence)
		if !ok {
			return nil
		}

	case token.LT, token.GT:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" comparison error " + operator.Literal)
			return nil
		}

	// logical operators
	case token.AND, token.OR:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError("Could not parse logical expression")
			return nil
		}

	// case token.IF:
	// 	right, ok = p.parseExpression(precedence)

	default:
		p.raiseError("Unexpected Infix operator " + string(operator.Type))
		return nil
	}

	if right == nil {
		p.raiseError("Right operand is nil")
		return nil
	}

	return &ast.Infix{operator, left, right}
}

func (p *Parser) parseIfElseExpression(left ast.ExpressionNode) (*ast.Trinary, bool) {
	if left == nil {
		p.raiseError("Left operand is nil")
		return nil, false
	}
	p.loadNextToken()

	if p.currentToken.Type != token.IF {
		p.raiseTokenError(token.IF)
		return nil, false
	}
	p.loadNextToken()

	condition, ok := p.parseExpression(precedences.LOWEST)
	if !ok {
		p.raiseError("Could not parse if expression")
		return nil, false
	}
	p.loadNextToken()

	if p.currentToken.Type != token.ELSE {
		p.raiseTokenError(token.ELSE)
	}
	p.loadNextToken()

	right, ok := p.parseExpression(precedences.LOWEST)
	if !ok {
		p.raiseError("Could not parse else expression")
		return nil, false
	}

	return &ast.Trinary{left, condition, right}, ok

}

func (p *Parser) getCurrentPrecedence() int {
	if precedence, ok := precedences.Precedence[p.currentToken.Type]; ok {
		return precedence
	}
	return precedences.LOWEST
}

func (p *Parser) getNextPrecedence() int {
	if p.nextToken.Type == token.SEMICOLON {
		p.raiseError("Unexpected semicolon;")
		return precedences.LOWEST
	}

	nextIsPrefix := (p.nextToken.Type == token.MINUS || p.nextToken.Type == token.BANG)
	if nextIsPrefix && p.getCurrentPrecedence() > precedences.LOWEST {
		return precedences.PREFIX
	}

	if precedence, ok := precedences.Precedence[p.nextToken.Type]; ok {
		return precedence
	}

	return precedences.LOWEST
}

// func (p *Parser) parseExpression(parentPrecedence int) (ast.ExpressionNode, bool) {
// 	ok := true

// 	var expr ast.ExpressionNode
// 	switch p.currentToken.Type {
// 	case token.IDENT:
// 		expr = &ast.IdentifierNode{p.currentToken}

// 	case token.INT:
// 		var err error
// 		expr, err = ast.NewIntegerLiteral(p.currentToken)
// 		if err != nil {
// 			p.raiseError(err.Error())
// 			return nil, !ok
// 		}
// 		// p.loadNextToken()
// 	case token.BANG, token.MINUS:
// 		expr, ok = p.parsePrefix()
// 		if !ok {
// 			return nil, !ok
// 		}
// 		// return nil, ok
// 	}

// 	if p.nextToken.Type == token.SEMICOLON {
// 		return expr, ok
// 	}

// 	// var left ast.ExpressionNode = expr
// 	for p.getNextPrecedence() > parentPrecedence {
// 		expr, ok = p.parseInfix(expr)
// 		if !ok {
// 			p.raiseError("Could not parse infix expression")
// 			break
// 		}

// 		if p.nextToken.Type == token.SEMICOLON {
// 			break
// 		}

// 	}

// 	return expr, ok
// }

// func parseNegativeInteger(intLit *ast.IntegerLiteral) (ast.IntegerLiteral, bool) {
// 	value, err := strconv.ParseInt(token.Literal, 0, 64)
// 	if err != nil {
// 		return nil, false
// 	}
// 	return &ast.IntegerLiteral{Token: token, Value: -value}, true
// }
