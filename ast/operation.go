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

type Infix struct {
	Operator token.Token
	Left     ExpressionNode
	Right    ExpressionNode
}

func (inFix *Infix) expressionNode() {}

func (inFix *Infix) GetTokenType() token.TokenType {
	return inFix.Operator.Type
}

func (inFix *Infix) GetTokenLiteral() string {
	return inFix.Operator.Literal
}

func (inFix *Infix) GetOperatorType() token.TokenType {
	return inFix.GetTokenType()
}

func (inFix *Infix) GetOperands() []ExpressionNode {
	return []ExpressionNode{inFix.Left, inFix.Right}
}

func (inFix *Infix) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(inFix.Left.ToString())
	out.WriteString(" ")
	out.WriteString(inFix.Operator.Literal)
	out.WriteString(" ")
	out.WriteString(inFix.Right.ToString())
	out.WriteString(")")

	return out.String()
}
