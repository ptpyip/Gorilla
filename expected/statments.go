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
	pass := true

	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		t.Errorf("Let statement not found. Got %q token", node.GetTokenType())
		return !pass
	}

	if letStmt.Identifier == nil {
		t.Errorf("Invalid Let statement: Identifier is nil")
		return !pass
	}

	if letStmt.Identifier.GetName() != expected.Name {
		// t.Log(letStmt.Expression.GetTokenLiteral())
		t.Errorf("letStmt.Identifier.Value not %s. got=%s",
			expected.Name, letStmt.Identifier.GetName(),
		)
		return !pass
	}

	if letStmt.Expression == nil {
		t.Errorf("Invalid Let statement: Expression is nil")
		return !pass
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
		t.Errorf("let statement not found. Got %q token", node.GetTokenType())
		return false
	}

	if returnStmt.ReturnValue == nil {
		t.Errorf("Invalid Return statement: ReturnValue is nil")
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
	return token.LPAREN
}

func (expected *BlockStatement) getTokenLiteral() string {
	return "{"
}

func (expected *BlockStatement) Test(t *testing.T, node ast.Node) bool {
	pass := true

	if node == nil {
		t.Errorf("Expected Block statement but got nil")
		return !pass
	}

	blockStmt, ok := node.(*ast.BlockStatement)
	if !ok {
		t.Errorf("Block statement not found. Got %q token", node.GetTokenType())
		return !pass
	}

	if len(blockStmt.Statements) != len(expected.Statements) {
		t.Errorf("Invalid Block statement: len(blockStmt.Statements) != len(expected.Statements)")
		return !pass
	}

	for i, expectedStmt := range expected.Statements {
		pass := expectedStmt.Test(t, blockStmt.Statements[i])
		if !pass {
			t.Errorf("Block statement error")
			return !pass
		}
	}
	return pass
}

type IfStatement struct {
	Condition ExpressionNode
	Statement StatementNode
	Else      *ElseStatement
}

func (expected *IfStatement) getTokenType() token.TokenType {
	return token.IF
}

func (expected *IfStatement) getTokenLiteral() string {
	return "if"
}

func (expected *IfStatement) Test(t *testing.T, node ast.Node) bool {
	pass := true

	ifStmt, ok := node.(*ast.IfStatement)
	if !ok {
		t.Errorf("If statement not found. Got %q token", node.GetTokenType())
		return !pass
	}

	if !expected.Condition.Test(t, ifStmt.Condition) {
		t.Errorf("Invalid If statement: Incorrect condition")
		return !pass
	}

	if !expected.Statement.Test(t, ifStmt.Statement) {
		t.Errorf("Invalid If statement: Incorrect statement")
		return !pass
	}

	if expected.Else != nil {
		return expected.Else.Test(t, ifStmt.Else)
	}

	return pass
}

type ElseStatement struct {
	Statement StatementNode
}

func (expected *ElseStatement) getTokenType() token.TokenType {
	return token.ELSE
}

func (expected *ElseStatement) getTokenLiteral() string {
	return "else"
}

func (expected *ElseStatement) Test(t *testing.T, node ast.Node) bool {
	pass := true

	if node == nil {
		t.Errorf("Expected Else statement but got nil")
		return !pass
	}

	elseStmt, ok := node.(*ast.ElseStatement)
	if !ok {
		t.Errorf("Else statement not found. Got %q token", node.GetTokenType())
		return !pass
	}

	if !expected.Statement.Test(t, elseStmt.Statement) {
		t.Errorf("Invalid Else statement: Incorrect statement")
		return !pass
	}

	return pass
}

func NewElseIfStatement(condition ExpressionNode, stmt StatementNode, elseNode *ElseStatement) *ElseStatement {
	return &ElseStatement{&IfStatement{condition, stmt, elseNode}}
}
