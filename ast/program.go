package ast

import "gorilla/token"

type Program struct {
	Statements []StatementNode
}

func (prog *Program) GetTokenLiteral(idx int) string {
	if len(prog.Statements) > 0 {
		return prog.Statements[idx].GetTokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// tok        token.Token
	Identifier *IdentifierNode
	Expression *ExpressionNode
}

func (letStmt *LetStatement) statementNode() {}

func (letStmt *LetStatement) GetTokenType() token.TokenType {
	return token.LET
}

func (letStmt *LetStatement) GetTokenLiteral() string {
	return "let"
}
