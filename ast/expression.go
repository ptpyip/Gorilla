package ast

import "gorilla/token"

type IdentifierNode struct {
	Token token.Token
	// Name  string
}

func (in *IdentifierNode) expressionNode() {}

func (in *IdentifierNode) GetTokenType() token.TokenType {
	return token.IDENT
}

func (in *IdentifierNode) GetTokenLiteral() string {
	return in.Token.Literal
}
func (in *IdentifierNode) GetName() string {
	return in.GetTokenLiteral()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (intLit *IntegerLiteral) expressionNode() {}

func (intLit *IntegerLiteral) GetTokenType() token.TokenType {
	return token.INT
}

func (intLit *IntegerLiteral) GetTokenLiteral() string {
	return intLit.Token.Literal
}

func (intLit *IntegerLiteral) GetValue() int64 {
	return intLit.Value
}
