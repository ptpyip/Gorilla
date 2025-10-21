package token

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(inputType TokenType, inputChar byte) Token {
	return Token{
		Type:    inputType,
		Literal: string(inputChar),
	}
}

// func New2CharToken(inputType TokenType, inputChar byte) Token {
// 	return Token{
// 		Type:    inputType,
// 		Literal: string(inputChar),
// 	}
// }
