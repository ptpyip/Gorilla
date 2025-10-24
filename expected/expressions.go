package expected

import (
	"gorilla/ast"
	"gorilla/token"
	"strconv"
	"testing"
)

type IntegerLiteral struct {
	Value int
}

func NewIntegerLiteral(value int) *IntegerLiteral {
	return &IntegerLiteral{Value: value}
}

func (expected *IntegerLiteral) getTokenType() token.TokenType {
	return token.INT
}

func (expected *IntegerLiteral) getTokenLiteral() string {
	return strconv.Itoa(expected.Value)
}

func (expected *IntegerLiteral) Test(t *testing.T, node ast.Node) bool {
	intLit, ok := node.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected IntegerLiteral. got %T expression", node.GetTokenType())
		return false
	}

	if intLit.GetTokenLiteral() != expected.getTokenLiteral() {
		t.Errorf("Expected intLit.TokenLiteral = %s. got = %s",
			expected.getTokenLiteral(), intLit.GetTokenLiteral(),
		)
		return false
	}

	// intValue, _ := strconv.ParseInt(expected., 0, 64)
	if intLit.GetValue() != int64(expected.Value) {
		t.Errorf("intLit.Value not %s. got=%d",
			expected.getTokenLiteral(), intLit.GetValue(),
		)
		return false
	}
	return true
}

type Identifier struct {
	Name string
}

func (expected *Identifier) getTokenType() token.TokenType {
	return token.IDENT
}

func (expected *Identifier) getTokenLiteral() string {
	return expected.Name
}

func (expected *Identifier) Test(t *testing.T, node ast.Node) bool {
	identifier, ok := node.(*ast.IdentifierNode)
	if !ok {
		t.Errorf("Expected IdentifierNode. got %T expression", node.GetTokenType())
		return false
	}

	if identifier.GetTokenLiteral() != expected.getTokenLiteral() {
		t.Errorf("Expected identifier.TokenLiteral = %s. got = %s",
			expected.getTokenLiteral(), identifier.GetTokenLiteral(),
		)
		return false
	}
	return true
}

type Prefix struct {
	OperatorType token.TokenType
	Operand      string
}

func (expected *Prefix) getTokenType() token.TokenType {
	return expected.OperatorType
}

func (expected *Prefix) getTokenLiteral() string {
	return string(expected.OperatorType)
}

func (expected *Prefix) Test(t *testing.T, node ast.Node) bool {

	prefixOp, ok := node.(*ast.Prefix)
	if !ok {
		t.Errorf("Expected Prefix. got %T expression", node.GetTokenType())
		return false
	}

	if prefixOp.GetOperatorType() != expected.OperatorType {
		t.Errorf("Expected prefixOp.Operator.Type = %s. got = %s",
			expected.getTokenLiteral(), prefixOp.GetOperatorType(),
		)
		return false
	}

	if prefixOp.Operand.ToString() != expected.Operand {
		t.Errorf("Expected prefixOp.Operand.TokenLiteral = %s. got = %s",
			expected.Operand, prefixOp.Operand.ToString(),
		)
		return false
	}
	return true
}
