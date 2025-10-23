package expected

import (
	"gorilla/ast"
	"gorilla/token"
	"testing"
)

// type ExpectedNode interface {
// 	test(t *testing.T, node ast.Node) bool
// 	getTokenLiteral() string
// 	getTokenType() token.TokenType
// }

type Node interface {
	getTokenType() token.TokenType
	getTokenLiteral() string

	Test(t *testing.T, node ast.Node) bool
}

type StatementNode interface {
	Node
}
type ExpressionNode interface {
	Node
}

// func (els *Le) getTokenLiteral() string {
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
