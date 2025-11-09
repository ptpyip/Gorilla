package ast

import (
	"gorilla/token"
	"strconv"
)

type Literal interface {
	Node
}

type IdentifierExpression struct {
	Token token.Token
	// Name  string
}

func (in *IdentifierExpression) expressionNode() {}

func (in *IdentifierExpression) GetTokenType() token.TokenType {
	return token.IDENT
}

func (in *IdentifierExpression) GetTokenLiteral() string {
	return in.Token.Literal
}
func (in *IdentifierExpression) GetName() string {
	return in.GetTokenLiteral()
}

func (in *IdentifierExpression) ToString() string {
	return in.GetTokenLiteral()
}

// type BoolToken

type BoolLiteral struct {
	Token token.Token
	// Value bool
}

func (boolLit *BoolLiteral) expressionNode() {}

func (boolLit *BoolLiteral) GetTokenType() token.TokenType {
	return boolLit.Token.Type
}

func (boolLit *BoolLiteral) GetTokenLiteral() string {
	return boolLit.Token.Literal
}

func (boolLit *BoolLiteral) GetValue() bool {
	return boolLit.Token.Type == token.TRUE
}

func (boolLit *BoolLiteral) ToString() string {
	return boolLit.GetTokenLiteral()
}

type IntegerLiteral struct {
	token token.Token
	value int64
}

func NewIntegerLiteral(token token.Token) (*IntegerLiteral, error) {
	value, err := strconv.ParseInt(token.Literal, 0, 64)
	if err != nil {
		return nil, err
		// return nil, fmt.Errorf("Could not parse integer literal: " + token.Literal)
	}
	return &IntegerLiteral{token, value}, nil
}

func (intLit *IntegerLiteral) expressionNode() {}

func (intLit *IntegerLiteral) GetTokenType() token.TokenType {
	return token.INT
}

func (intLit *IntegerLiteral) GetTokenLiteral() string {
	return intLit.token.Literal
}

func (intLit *IntegerLiteral) GetValue() int64 {
	return intLit.value
}

func (intLit *IntegerLiteral) ToString() string {
	return intLit.GetTokenLiteral()
}
