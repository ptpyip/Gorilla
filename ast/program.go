package ast

import (
	"bytes"
	"gorilla/token"
)

type Program struct {
	Statements []StatementNode
}

func (prog *Program) GetTokenLiteral(idx int) string {
	if len(prog.Statements) > 0 {
		return prog.Statements[idx].GetTokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// tok        token.Token
	Identifier *IdentifierExpression
	Expression ExpressionNode
}

func (letStmt *LetStatement) statementNode() {}

func (letStmt *LetStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("let ")
	out.WriteString(letStmt.Identifier.GetTokenLiteral())
	out.WriteString(" = ")
	out.WriteString(letStmt.Expression.ToString())
	out.WriteString(";")
	return out.String()
}

func (letStmt *LetStatement) GetTokenType() token.TokenType {
	return token.LET
}

func (letStmt *LetStatement) GetTokenLiteral() string {
	return "let"
}

type ReturnStatement struct {
	ReturnValue ExpressionNode
}

func (returnStmt *ReturnStatement) statementNode() {}

func (returnStmt *ReturnStatement) GetTokenType() token.TokenType {
	return token.RETURN
}

func (returnStmt *ReturnStatement) GetTokenLiteral() string {
	return "return"
}

func (returnStmt *ReturnStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("return ")
	out.WriteString(returnStmt.ReturnValue.ToString())
	out.WriteString(";")
	return out.String()
}

type BlockStatement struct {
	Statements []StatementNode
}

func (blockStmt *BlockStatement) statementNode() {}

func (blockStmt *BlockStatement) GetTokenType() token.TokenType {
	return token.LBRACE
}

func (blockStmt *BlockStatement) GetTokenLiteral() string {
	return "{"
}

func (blockStmt *BlockStatement) AppendStatement(statement StatementNode) {
	blockStmt.Statements = append(blockStmt.Statements, statement)
}

func (blockStmt *BlockStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for i, stmt := range blockStmt.Statements {
		out.WriteString("\t")
		out.WriteString(stmt.ToString())
		if i+1 != len(blockStmt.Statements) {
			out.WriteString(";\n")
		}
	}
	out.WriteString("\n}")
	return out.String()
}

func NewIfStatement(condition ExpressionNode, block *BlockStatement) *IfStatement {
	return &IfStatement{condition, block, nil}
}

func NewIfElseStatement(condition ExpressionNode, block *BlockStatement, elseBlock StatementNode) *IfStatement {
	return &IfStatement{condition, block, &ElseStatement{elseBlock}}
}

type IfStatement struct {
	// IfToken token.Token
	// Conditions []ExpressionNode
	// Statements []StatementNode
	Condition ExpressionNode
	Statement *BlockStatement
	Else      *ElseStatement // nullable
}

func (ifStmt *IfStatement) statementNode() {}

func (ifStmt *IfStatement) GetTokenType() token.TokenType {
	return token.IF
}

func (ifStmt *IfStatement) GetTokenLiteral() string {
	return "if"
}

func (ifStmt *IfStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString("if " + ifStmt.Condition.ToString() + " ")
	out.WriteString(ifStmt.Statement.ToString())
	// out.WriteString("\n} ")
	if ifStmt.Else != nil {
		out.WriteString(ifStmt.Else.ToString())
		// out.WriteString(" else {\n\t")
		// out.WriteString(ifStmt.Else.Statement.ToString())
		// out.WriteString("\n}")
	}
	return out.String()
}

type ElseStatement struct {
	// ElseToken token.Token
	Statement StatementNode // IfStatement or BlockStatement
}

func (elseStmt *ElseStatement) statementNode() {}

func (elseStmt *ElseStatement) GetTokenType() token.TokenType {
	return token.ELSE
}

func (elseStmt *ElseStatement) GetTokenLiteral() string {
	return "else"
}

func (elseStmt *ElseStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString(" else ")
	// if elseStmt.Statement.GetTokenType() != token.IF {
	// 	// out.WriteString("{\n\t")
	// 	out.WriteString(elseStmt.Statement.ToString())
	// 	// out.WriteString("\n}")
	// } else {
	// 	out.WriteString(elseStmt.Statement.ToString())
	// }

	out.WriteString(elseStmt.Statement.ToString())

	return out.String()
}
