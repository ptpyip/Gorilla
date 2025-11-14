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
			// for _, msg := range p.Errors {
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
		p.loadNextToken()
		return stmt

	case token.LBRACE:
		block, ok := p.parseBlockStatement()
		if !ok {
			p.raiseError("Could not parse block statement")
			return nil
		}
		p.loadNextToken()
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
	// p.loadNextToken() // load token after '}'
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
	p.loadNextToken()

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
