package ast

import (
	"bytes"
	"gorilla/token"
)

type Operation interface {
	ExpressionNode
	GetOperatorType() token.TokenType
	GetOperands() []ExpressionNode
	// GetPrecedence() int
}

type Prefix struct {
	Operator token.Token
	Operand  ExpressionNode
}

func (prefixOp *Prefix) expressionNode() {}

func (prefixOp *Prefix) GetTokenType() token.TokenType {
	return prefixOp.Operator.Type
}

func (prefixOp *Prefix) GetTokenLiteral() string {
	return prefixOp.Operator.Literal
}

func (prefixOp *Prefix) GetOperatorType() token.TokenType {
	return prefixOp.GetTokenType()
}

func (prefixOp *Prefix) GetOperands() []ExpressionNode {
	return []ExpressionNode{prefixOp.Operand}
}

func (prefixOp *Prefix) ToString() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefixOp.Operator.Literal)
	out.WriteString(" ")
	out.WriteString(prefixOp.Operand.ToString())
	out.WriteString(")")
	return out.String()
}
