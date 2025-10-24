package ast

import (
	"bytes"
	"gorilla/token"
)

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
	Expression ExpressionNode
}

func (letStmt *LetStatement) statementNode() {}

func (letStmt *LetStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("let ")
	out.WriteString(letStmt.Identifier.GetTokenLiteral())
	out.WriteString(" = ")
	out.WriteString(letStmt.Expression.ToString())
	return out.String()
}

func (letStmt *LetStatement) GetTokenType() token.TokenType {
	return token.LET
}

func (letStmt *LetStatement) GetTokenLiteral() string {
	return "let"
}

type ReturnStatement struct {
	ReturnValue ExpressionNode
}

func (returnStmt *ReturnStatement) statementNode() {}

func (returnStmt *ReturnStatement) GetTokenType() token.TokenType {
	return token.RETURN
}

func (returnStmt *ReturnStatement) GetTokenLiteral() string {
	return "return"
}

func (returnStmt *ReturnStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("return ")
	out.WriteString(returnStmt.ReturnValue.ToString())
	return out.String()
}
