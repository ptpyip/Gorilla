package parser

// import (
// 	"gorilla/ast"
// 	"gorilla/token"
// 	"strconv"
// 	"testing"
// )

// // type ExpectedNode interface {
// // 	test(t *testing.T, node ast.Node) bool
// // 	getTokenLiteral() string
// // 	getTokenType() token.TokenType
// // }

// type ExpectedNode interface {
// 	getTokenType() token.TokenType
// 	getTokenLiteral() string
// 	test(t *testing.T, node ast.Node) bool
// }

// type ExpectedStatement interface {
// 	ExpectedNode
// }
// type ExpectedExpression interface {
// 	ExpectedNode
// }

// type ExpectedLetStatement struct {
// 	name       string
// 	expression ExpectedExpression
// }

// func (els *ExpectedLetStatement) getTokenLiteral() string {
// 	return "let"
// }

// func (els *ExpectedLetStatement) getTokenType() token.TokenType {
// 	return token.LET
// }

// func (expected *ExpectedLetStatement) test(t *testing.T, node ast.Node) bool {
// 	letStmt, ok := node.(*ast.LetStatement)
// 	if !ok {
// 		t.Errorf("let satetment not found. Got %q token", node.GetTokenType())
// 		return false
// 	}

// 	if letStmt.Identifier == nil {
// 		t.Errorf("Invalid Let satement: letStmt.Identifier is nil")
// 		return false
// 	}

// 	if letStmt.Identifier.GetName() != expected.name {
// 		t.Log(letStmt.Expression.GetTokenLiteral())
// 		t.Errorf("letStmt.Identifier.Value not %s. got=%s",
// 			expected.name, letStmt.Identifier.GetName(),
// 		)
// 		return false
// 	}

// 	if letStmt.Expression == nil {
// 		t.Errorf("Invalid Let satement: letStmt.Expression is nil")
// 		return false
// 	} else {
// 		expected.expression.test(t, letStmt.Expression)
// 	}

// 	return true
// }

// type ExpectedReturnStatment struct {
// 	expression ExpectedExpression
// }

// func (expected *ExpectedReturnStatment) getTokenLiteral() string {
// 	return "return"
// }
// func (expected *ExpectedReturnStatment) getTokenType() token.TokenType {
// 	return token.RETURN
// }

// func (expected *ExpectedReturnStatment) test(t *testing.T, node ast.Node) bool {
// 	returnStmt, ok := node.(*ast.ReturnStatement)
// 	if !ok {
// 		t.Errorf("let satetment not found. Got %q token", node.GetTokenType())
// 		return false
// 	}

// 	if returnStmt.ReturnValue == nil {
// 		t.Errorf("Invalid Return satement: ReturnValue is nil")
// 		return false
// 	}

// 	expected.expression.test(t, returnStmt.ReturnValue)

// 	return true
// }

// type ExpectedIntegerLiteral struct {
// 	value int
// }

// func NewExpectedInt(value int) *ExpectedIntegerLiteral {
// 	return &ExpectedIntegerLiteral{value: value}
// }

// func (expected *ExpectedIntegerLiteral) getTokenLiteral() string {
// 	return strconv.Itoa(expected.value)
// }

// func (expected *ExpectedIntegerLiteral) getTokenType() token.TokenType {
// 	return token.INT
// }

// func (expected *ExpectedIntegerLiteral) test(t *testing.T, node ast.Node) bool {
// 	intLit, ok := node.(*ast.IntegerLiteral)
// 	if !ok {
// 		t.Errorf("Expected IntegerLiteral. got %T expression", node.GetTokenType())
// 		return false
// 	}

// 	if intLit.GetTokenLiteral() != expected.getTokenLiteral() {
// 		t.Errorf("Expected intLit.TokenLiteral = %s. got = %s",
// 			expected.getTokenLiteral(), intLit.GetTokenLiteral(),
// 		)
// 		return false
// 	}

// 	// intValue, _ := strconv.ParseInt(expected., 0, 64)
// 	if intLit.GetValue() != int64(expected.value) {
// 		t.Errorf("intLit.Value not %s. got=%d",
// 			expected.getTokenLiteral(), intLit.GetValue(),
// 		)
// 		return false
// 	}
// 	return true
// }

// type ExpectedIdentifier struct {
// 	name string
// }

// func (expected *ExpectedIdentifier) getTokenLiteral() string {
// 	return expected.name
// }

// func (expected *ExpectedIdentifier) getTokenType() token.TokenType {
// 	return token.IDENT
// }

// func (expected *ExpectedIdentifier) test(t *testing.T, node ast.Node) bool {
// 	identifier, ok := node.(*ast.IdentifierNode)
// 	if !ok {
// 		t.Errorf("Expected IdentifierNode. got %T expression", node.GetTokenType())
// 		return false
// 	}

// 	if identifier.GetTokenLiteral() != expected.getTokenLiteral() {
// 		t.Errorf("Expected identifier.TokenLiteral = %s. got = %s",
// 			expected.getTokenLiteral(), identifier.GetTokenLiteral(),
// 		)
// 		return false
// 	}
// 	return true
// }

// // func (expected *ExpectedExpression) getTokenLiteral() string {
// // 	return expected.value
// // }

// // func (expected *ExpectedExpression) getTokenType() token.TokenType {
// // 	return nil
// // }

// // func (expected *ExpectedExpression) test(t *testing.T, node ast.Node) bool {
// // 	letStmt, ok := node.(*ast.ExpressionNode)

// // 	// if expr == nil {
// // 	// 	t.Errorf("Expression is nil")
// // 	// 	return false
// // 	// }
// // 	switch expr.GetTokenType() {
// // 	case token.INT:
// // 		intLit, ok := expr.(*ast.IntegerLiteral)
// // 		if !ok {
// // 			t.Errorf("Expected IntegerLiteral. got %T expression", expr.GetTokenType())
// // 			return false
// // 		}

// // 		if intLit.GetTokenLiteral() != value {
// // 			t.Errorf("Expected intLit.TokenLiteral = %s. got = %s",
// // 				value, intLit.GetTokenLiteral(),
// // 			)
// // 			return false
// // 		}
// // 		intValue, _ := strconv.ParseInt(value, 0, 64)
// // 		if intLit.GetValue() != intValue {
// // 			t.Errorf("intLit.Value not %s. got=%d",
// // 				value, intLit.GetValue(),
// // 			)
// // 			return false
// // 		}
// // 	default:
// // 		t.Errorf("Unexpected expression type %s", expr.GetTokenType())
// // 		return false
// // 	}
// // 	return true
// // }

// // f
