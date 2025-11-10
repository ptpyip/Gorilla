package parser

import (
	"gorilla/ast"
	"gorilla/parser/precedences"
	"gorilla/token"
)

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	prog := &ast.Program{}

	for p.currentToken.Type != token.EOF {
		// println("Parsing statement:", p.currentToken.Literal)
		statement := p.parseStatement()
		if statement == nil {
			// for _, msg := range p.errors {
			// 	println(msg)
			// }
			p.raiseParseProgramError()
			return prog, false
		}
		// if statement == nil {
		// 	p.raiseParseStatementError(statement)
		// 	return prog, false
		// }

		prog.Statements = append(prog.Statements, statement)
		// if p.currentToken.Type == token.SEMICOLON {
		// 	p.loadNextToken()
		// 	// p.raiseNextTokenError(token.SEMICOLON)
		// 	// return prog, !ok
		// }

		// p.loadNextToken()
	}
	return prog, true
}

func (p *Parser) parseStatement() ast.StatementNode {
	switch p.currentToken.Type {

	// === Ends with ';' === //
	case token.LET:
		if p.nextToken.Type != token.IDENT {
			p.raiseNextTokenError(token.IDENT)
			return nil
		}
		p.loadNextToken()

		identifier := &ast.IdentifierExpression{p.currentToken}
		p.loadNextToken()

		if p.currentToken.Type != token.ASSIGN {
			p.raiseNextTokenError(token.ASSIGN)
			return nil
		}
		p.loadNextToken()

		expression, ok := p.parseExpression(precedences.LOWEST)
		if !ok {
			p.raiseExpressionError()
			p.raiseParseStatementError(token.LET, nil)
			return nil
		}
		p.loadNextToken()

		stmt := &ast.LetStatement{identifier, expression}
		if p.currentToken.Type != token.SEMICOLON {
			p.raiseTokenError(token.SEMICOLON)
			p.raiseParseStatementError(token.LET, stmt)
			return nil
		}
		p.loadNextToken()
		// p.skipToSemicolon()

		return stmt

	case token.RETURN:
		if p.nextToken.Type == token.SEMICOLON {
			// empty return
			return &ast.ReturnStatement{}
		}
		p.loadNextToken()

		returnValue, ok := p.parseExpression(precedences.LOWEST)
		if !ok {
			p.raiseParseStatementError(token.RETURN, nil)
			return nil
		}
		p.loadNextToken()

		stmt := &ast.ReturnStatement{returnValue}
		if p.currentToken.Type != token.SEMICOLON {
			p.raiseTokenError(token.SEMICOLON)
			p.raiseParseStatementError(token.RETURN, stmt)
			return nil
		}
		p.loadNextToken()

		return stmt

	// === Ends with '}' === //
	case token.IF:
		stmt, ok := p.parseIfElseStatement()
		if !ok {
			p.raiseError("Could not parse if statement")
			return nil
		}
		return stmt

	case token.LBRACE:
		block, ok := p.parseBlockStatement()
		if !ok {
			p.raiseError("Could not parse block statement")
			return nil
		}
		return block

	// === Edge cases === //
	case token.RBRACE:
		p.raiseError("Could not parse statement with token: " + string(p.currentToken.Type))
		return nil

	case token.SEMICOLON:
		p.raiseError("Unexpected token: " + string(p.currentToken.Type))
		return nil

	default:
		p.raiseError("Unexpected token: " + string(p.currentToken.Type))
		p.loadNextToken()
		return nil
	}

}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, bool) {
	if p.currentToken.Type != token.LBRACE {
		p.raiseNextTokenError(token.LBRACE)
		return nil, false
	}
	p.loadNextToken()

	block := &ast.BlockStatement{}
	for p.currentToken.Type != token.RBRACE {
		// println("Parsing block statement: ", p.currentToken.Literal)
		statement := p.parseStatement()
		if statement == nil {
			p.raiseBlockStatementError(block.Statements)
			// if len(block.Statements) > 0 {
			// 	return block, false
			// }
			return nil, false
		}

		block.AppendStatement(statement)
		if p.currentToken.Type == token.SEMICOLON {
			p.loadNextToken()
		}
	}
	p.loadNextToken() // load token after '}'
	// println("Finsh parsing block statement: with ", len(block.Statements), " statements")

	return block, true
}

func (p *Parser) parseIfElseStatement() (ast.StatementNode, bool) {
	if p.nextToken.Type != token.LPAREN {
		p.raiseNextTokenError(token.LPAREN)
		return nil, false
	}
	p.loadNextToken()

	condition, ok := p.parseExpression(precedences.LOWEST)
	if !ok {
		p.raiseExpressionError()
		return nil, false
	}
	p.loadNextToken()

	block, ok := p.parseBlockStatement()
	if !ok {
		p.raiseError("Could not parse block statement")
		return nil, false
	}

	if p.currentToken.Type != token.ELSE {
		return ast.NewIfStatement(condition, block), true
	}
	p.loadNextToken()

	var elseBlock ast.StatementNode
	if p.currentToken.Type == token.IF {
		elseBlock, ok = p.parseIfElseStatement()
		if !ok {
			p.raiseError("Could not parse else if statement")
			// if elseBlock != nil {
			// 	return elseBlock, !ok
			// }
			return nil, false
		}

	} else {
		elseBlock, ok = p.parseBlockStatement()
		if !ok {
			p.raiseError("Could not parse else statement")
			return nil, false
		}
	}
	return ast.NewIfElseStatement(condition, block, elseBlock), true
}

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
