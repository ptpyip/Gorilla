package expected

import (
	"gorilla/ast"
	"gorilla/token"
	"testing"
)

type LetStatement struct {
	Name       string
	Expression ExpressionNode
}

func (expected *LetStatement) getTokenType() token.TokenType {
	return token.LET
}

func (expected *LetStatement) getTokenLiteral() string {
	return "let"
}

func (expected *LetStatement) Test(t *testing.T, node ast.Node) bool {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		t.Errorf("Let satetment not found. Got %q token", node.GetTokenType())
		return false
	}

	if letStmt.Identifier == nil {
		t.Errorf("Invalid Let satement: Identifier is nil")
		return false
	}

	if letStmt.Identifier.GetName() != expected.Name {
		// t.Log(letStmt.Expression.GetTokenLiteral())
		t.Errorf("letStmt.Identifier.Value not %s. got=%s",
			expected.Name, letStmt.Identifier.GetName(),
		)
		return false
	}

	if letStmt.Expression == nil {
		t.Errorf("Invalid Let satement: Expression is nil")
		return false
	}

	return expected.Expression.Test(t, letStmt.Expression)
}

type ReturnStatement struct {
	Expression ExpressionNode
}

func (expected *ReturnStatement) getTokenType() token.TokenType {
	return token.RETURN
}

func (expected *ReturnStatement) getTokenLiteral() string {
	return "return"
}

func (expected *ReturnStatement) Test(t *testing.T, node ast.Node) bool {
	returnStmt, ok := node.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("let satetment not found. Got %q token", node.GetTokenType())
		return false
	}

	if returnStmt.ReturnValue == nil {
		t.Errorf("Invalid Return satement: ReturnValue is nil")
		return false
	}

	return expected.Expression.Test(t, returnStmt.ReturnValue)
}

type BlockStatement struct {
	Statements []StatementNode
}

func NewBlockStatement(statements ...StatementNode) *BlockStatement {
	return &BlockStatement{Statements: statements}
}

func (expected *BlockStatement) getTokenType() token.TokenType {
	return token.LBRACE
}

func (expected *BlockStatement) getTokenLiteral() string {
	return "{"
}

func (expected *BlockStatement) Test(t *testing.T, node ast.Node) bool {
	blockStmt, ok := node.(*ast.BlockStatement)
	if !ok {
		t.Errorf("Block satetment not found. Got %q token", node.GetTokenType())
		return !ok
	}

	if len(blockStmt.Statements) != len(expected.Statements) {
		t.Errorf("Invalid Block satement: len(blockStmt.Statements) != len(expected.Statements)")
		return !ok
	}

	for i, expectedStmt := range expected.Statements {
		pass := expectedStmt.Test(t, blockStmt.Statements[i])
		if !pass {
			t.Errorf("Block statement error")
			return !ok
		}
	}
	return ok
}
