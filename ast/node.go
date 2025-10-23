package ast

import "gorilla/token"

type Node interface {
	GetTokenLiteral() string
	GetTokenType() token.TokenType
}

type StatementNode interface {
	Node
	statementNode()
}

type ExpressionNode interface {
	Node
	expressionNode()
}

type IdentifierNode struct {
	Token token.Token
	// Name  string
}

func (in *IdentifierNode) expressionNode() {}

func (in *IdentifierNode) GetTokenLiteral() string {
	return in.Token.Literal
}

func (in *IdentifierNode) GetTokenType() token.TokenType {
	return in.Token.Type
}

func (in *IdentifierNode) GetName() string {
	return in.GetTokenLiteral()
}
