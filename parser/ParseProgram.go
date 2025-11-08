package parser

import (
	"gorilla/ast"
	"gorilla/parser/precedences"
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
			p.raiseError("Could not parse statement with token: " + p.currentToken.Literal)
			return prog, !ok
		}

		prog.Statements = append(prog.Statements, statement)
		if p.currentToken.Type == token.SEMICOLON {
			p.loadNextToken()
			// p.raiseNextTokenError(token.SEMICOLON)
			// return prog, !ok
		}

		// p.loadNextToken()
	}
	return prog, ok
}

func (p *Parser) parseStatement() ast.StatementNode {
	switch p.currentToken.Type {
	case token.LET:
		if p.nextToken.Type != token.IDENT {
			p.raiseNextTokenError(token.IDENT)
			return nil
		}
		p.loadNextToken()

		identifier := &ast.IdentifierNode{p.currentToken}
		p.loadNextToken()

		if p.currentToken.Type != token.ASSIGN {
			p.raiseNextTokenError(token.ASSIGN)
			return nil
		}
		p.loadNextToken()

		expression, ok := p.parseExpression(precedences.LOWEST)
		if !ok {
			p.raiseExpressionError()
			return nil
		}
		p.loadNextToken()

		stmt := &ast.LetStatement{identifier, expression}
		p.skipToSemicolon()

		return stmt

	case token.RETURN:
		p.loadNextToken()

		if p.currentToken.Type == token.SEMICOLON {
			// empty return
			return &ast.ReturnStatement{}
		}

		returnValue, ok := p.parseExpression(precedences.LOWEST)
		if !ok {
			p.raiseExpressionError()
			return nil
		}
		p.skipToSemicolon()

		return &ast.ReturnStatement{returnValue}

	case token.IF:
		stmt, ok := p.parseIfElseStatement()
		if !ok {
			p.raiseError("Could not parse if statement")
			return nil
		}
		return stmt

	case token.LBRACE:
		block, ok := p.parseBlockStament()
		if !ok {
			p.raiseError("Could not parse block statement")
			return nil
		}
		return block
	case token.RBRACE:
		return nil

	case token.SEMICOLON:
		p.raiseError("Unexpected token: " + p.currentToken.Literal)
		return nil

	default:
		p.raiseError("Unexpected token: " + p.currentToken.Literal)
		p.loadNextToken()
		return nil
	}
}

func (p *Parser) parseBlockStament() (*ast.BlockStatement, bool) {
	ok := true

	if p.currentToken.Type != token.LBRACE {
		p.raiseNextTokenError(token.LBRACE)
		return nil, !ok
	}
	p.loadNextToken()

	block := &ast.BlockStatement{}
	for p.currentToken.Type != token.RBRACE {
		// println("Parsing block statement: ", p.currentToken.Literal)
		statement := p.parseStatement()
		if statement == nil {
			p.raiseBloackStatementError(block.Statements)
			return nil, !ok
		}

		block.AppendStatement(statement)
		if p.currentToken.Type == token.SEMICOLON {
			p.loadNextToken()
		}
	}
	p.loadNextToken() // load token after '}'
	// println("Finsh parsing block statement: with ", len(block.Statements), " statements")

	return block, ok
}

func (p *Parser) parseIfElseStatement() (ast.StatementNode, bool) {
	ok := true

	if p.nextToken.Type != token.LPAREN {
		p.raiseNextTokenError(token.LPAREN)
		return nil, !ok
	}
	p.loadNextToken()

	condition, ok := p.parseExpression(precedences.LOWEST)
	if !ok {
		p.raiseExpressionError()
		return nil, !ok
	}
	p.loadNextToken()

	block, ok := p.parseBlockStament()
	if !ok {
		p.raiseError("Could not parse block statement")
		return nil, !ok
	}

	if p.currentToken.Type != token.ELSE {
		return ast.NewIfStatement(condition, block), ok
	}
	p.loadNextToken()

	var elseBlock ast.StatementNode
	if p.currentToken.Type == token.IF {
		elseBlock = p.parseStatement()
		if elseBlock == nil {
			p.raiseError("Could not parse else if statement")
			return nil, !ok
		}

	} else {
		elseBlock, ok = p.parseBlockStament()
		if !ok {
			p.raiseError("Could not parse else statement")
			return nil, !ok
		}
	}
	return ast.NewIfElseStatement(condition, block, elseBlock), ok
}

func (p *Parser) parseExpression(parentPrecedence int) (ast.ExpressionNode, bool) {
	// println("Parsing :" + p.currentToken.Literal)
	ok := true

	var expr ast.ExpressionNode
	switch p.currentToken.Type {
	case token.IDENT:
		expr = &ast.IdentifierNode{p.currentToken}

	case token.TRUE, token.FALSE:
		expr = &ast.BoolLiteral{p.currentToken}

	case token.INT:
		var err error
		expr, err = ast.NewIntegerLiteral(p.currentToken)
		if err != nil {
			p.raiseError(err.Error())
			return nil, !ok
		}

	case token.LPAREN, token.BANG, token.MINUS:
		expr, ok = p.parsePrefix()
		if !ok {
			return nil, !ok
		}

	// case token.RBRACE:
	// 	p.raiseError()
	default:
		panic("Unexpected token for parseExpression: " + p.currentToken.Literal)
	}

	if p.nextToken.Type == token.SEMICOLON {
		return expr, ok
	}

	// Next token is infix operator
	for p.getNextPrecedence() > parentPrecedence {
		expr, ok = p.parseInfix(expr)
		if !ok {
			p.raiseError("Could not parse infix expression")
			break
		}

		if p.nextToken.Type == token.SEMICOLON {
			break
		}

	}

	return expr, ok
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
		panic("Invalid token type: " + p.currentToken.Literal)
	}
}

func (p *Parser) parseInfix(left ast.ExpressionNode) (ast.ExpressionNode, bool) {
	// println("parsePrefix on token: " + p.currentToken.Literal)
	ok := true

	if left == nil {
		p.raiseError("Left operand is nil")
		return nil, !ok
	}
	p.loadNextToken()

	operator := p.currentToken
	precedence := p.getCurrentPrecedence()
	p.loadNextToken()

	// println("Infix operator: " + operator.Literal)
	// println("\twith operator precedence = ", precedence)

	var right ast.ExpressionNode
	switch operator.Type {
	// algrithmic operators
	case token.PLUS, token.MINUS:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" algrithmic error " + operator.Literal)
			return nil, !ok
		}

	case token.ASTERISK, token.SLASH:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" algrithmic error " + operator.Literal)
			return nil, !ok

		}

	// comparison operators
	case token.EQ, token.NOT_EQ, token.LE, token.GE:
		right, ok = p.parseExpression(precedence)
		if !ok {
			return nil, !ok
		}

	case token.LT, token.GT:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError(" comparison error " + operator.Literal)
			return nil, !ok
		}

	// logical operators
	case token.AND, token.OR:
		right, ok = p.parseExpression(precedence)
		if !ok {
			p.raiseError("Could not parse logical expression")
			return nil, !ok
		}

	default:
		p.raiseError(" Unexpected operator " + operator.Literal)
		return nil, !ok
	}

	if right == nil {
		p.raiseError("Right operand is nil")
		return nil, !ok
	}

	return &ast.Infix{operator, left, right}, ok
}

func (p *Parser) skipToSemicolon() {
	for p.currentToken.Type != token.SEMICOLON {
		p.loadNextToken()
		// println("Skip:", p.currentToken.Literal)
	}
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

// func (p *Parser) parse() (ast.ExpressionNode, bool) {
// 	ok := true
// 	switch p.currentToken.Type {

// 	case token.IDENT:
// 		return &ast.IdentifierNode{p.currentToken}, ok

// 	case token.INT:
// 		if p.nextToken.Type != token.SEMICOLON {
// 			p.raiseNextTokenError(token.SEMICOLON)
// 			return nil, !ok
// 		}

// 		expr, err := ast.NewIntegerLiteral(p.currentToken)
// 		if err != nil {
// 			p.raiseError(err.Error())
// 			return nil, !ok
// 		}
// 		// p.loadNextToken()

// 		return expr, ok

// 	case token.BANG, token.MINUS:
// 		operator := p.currentToken
// 		p.loadNextToken()

// 		expr, ok := p.parseExpression()
// 		if !ok {
// 			p.raiseError(
// 				"Could not parse expression after prefix operator " + operator.Literal,
// 			)
// 			return nil, !ok
// 		}

// 		// optimize negative value
// 		if operator.Type == token.MINUS && expr.GetTokenType() == token.INT {
// 			intLit, err := ast.NewIntegerLiteral(
// 				token.Token{
// 					Type:    token.INT,
// 					Literal: "-" + expr.GetTokenLiteral(),
// 				},
// 			)
// 			if err != nil {
// 				p.raiseError(err.Error())
// 				return nil, !ok
// 			}
// 			return intLit, ok
// 		}

// 		return &ast.Prefix{operator, expr}, ok

// 	default:
// 		return nil, ok
// 	}
// }
