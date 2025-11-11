package ast

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
