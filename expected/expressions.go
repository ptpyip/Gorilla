package expected

import (
	"gorilla/ast"
	"gorilla/token"
	"strconv"
	"testing"
)

type BoolLiteral struct {
	Value bool
}

func NewBoolLiteral(value bool) *BoolLiteral {
	return &BoolLiteral{Value: value}
}

func (expected *BoolLiteral) getTokenType() token.TokenType {
	if expected.Value {
		return token.TRUE
	}
	return token.FALSE
}

func (expected *BoolLiteral) getTokenLiteral() string {
	if expected.Value {
		return "true"
	}
	return "false"
}

func (expected *BoolLiteral) Test(t *testing.T, node ast.Node) bool {
	boolLit, ok := node.(*ast.BoolLiteral)
	if !ok {
		t.Errorf("Expected BoolLiteral. got %T expression", node.GetTokenType())
		return false
	}

	// intValue, _ := strconv.ParseInt(expected., 0, 64)
	if boolLit.GetValue() != expected.Value {
		t.Errorf("boolLit.Value not %s. got=%s",
			expected.getTokenLiteral(), boolLit.GetTokenLiteral(),
		)
		return false
	}
	return ok
}

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
	identifier, ok := node.(*ast.IdentifierExpression)
	if !ok {
		t.Errorf("Expected IdentifierExpression. got %T expression", node.GetTokenType())
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
	Operand      ExpressionNode
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
		println(node.ToString())
		t.Errorf("Expected Prefix. got %T expression", node.GetTokenType())
		return false
	}

	if prefixOp.GetOperatorType() != expected.OperatorType {
		t.Errorf("Expected prefixOp.Operator.Type = %s. got = %s",
			expected.getTokenLiteral(), prefixOp.GetOperatorType(),
		)
		return false
	}

	// if prefixOp.Operand.ToString() != expected.Operand {
	// 	t.Errorf("Expected prefixOp.Operand.TokenLiteral = %s. got = %s",
	// 		expected.Operand, prefixOp.Operand.ToString(),
	// 	)
	// 	return false
	// }
	return expected.Operand.Test(t, prefixOp.Operand)
}

type Infix struct {
	OperatorType token.TokenType
	Left         ExpressionNode
	Right        ExpressionNode
}

func (expected *Infix) getTokenType() token.TokenType {
	return expected.OperatorType
}

func (expected *Infix) getTokenLiteral() string {
	return string(expected.OperatorType)
}

func (expected *Infix) Test(t *testing.T, node ast.Node) bool {
	pass := true
	if node == nil {
		t.Errorf("Expected Infix. got nil")
		return !pass
	}

	inFix, ok := node.(*ast.Infix)
	if !ok {
		t.Logf("Expected Infix. got %s expression", node.ToString())
		t.Logf("%s", node.GetTokenLiteral())
		return !pass
	}

	if inFix.GetOperatorType() != expected.OperatorType {
		t.Errorf("Expected inFix.Operator.Type = %s. got = %s",
			expected.getTokenLiteral(), inFix.GetOperatorType(),
		)
		return !pass
	}

	if inFix.Left == nil {
		t.Errorf("Invalid Infix satement: inFix.Left is nil")
		return !pass
	} else if inFix.Right == nil {
		t.Errorf("Invalid Infix satement: inFix.Right is nil")
		return !pass
	}

	return expected.Left.Test(t, inFix.Left) && expected.Right.Test(t, inFix.Right)
}
