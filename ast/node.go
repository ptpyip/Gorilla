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
