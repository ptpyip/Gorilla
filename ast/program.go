package ast

import "gorilla/token"

type Program struct {
	Statements []Statement
}

func (prog *Program) GetTokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].GetTokenLiteral()
	} else {
		return ""
	}
}

type LetStaement struct {
	Identifier *IdentifierNode
	Value      Expression
}

func (ls *LetStaement) statementNode() {}

func (ls *LetStaement) GetTokenLiteral() string {
	return ls.Identifier.Token.Literal
}

type IdentifierNode struct {
	Token token.Token
	Value string
}

func (in *IdentifierNode) expressionNode() {}

func (in *IdentifierNode) GetTokenLiteral() string {
	return in.Token.Literal
}
