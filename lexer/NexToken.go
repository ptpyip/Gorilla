package lexer

import "gorilla/token"

func (lx *Lexer) NextToken() token.Token {
	var nextTokenType token.TokenType
	defer lx.readChar()

	lx.skip()
	switch lx.currentChar {
	case 0:
		return token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	case '=':
		if lx.getNextChar() == '=' {
			lx.readChar()
			return token.Token{
				Type:    token.EQ,
				Literal: "==",
			}
		} else {
			nextTokenType = token.ASSIGN
		}
	case '!':
		if lx.getNextChar() == '=' {
			lx.readChar()
			return token.Token{
				Type:    token.NOT_EQ,
				Literal: "!=",
			}
		} else {
			nextTokenType = token.BANG
		}
	case '<':
		if lx.getNextChar() == '=' {
			lx.readChar()
			return token.Token{
				Type:    token.LE,
				Literal: "<=",
			}
		} else {
			nextTokenType = token.LT
		}
	case '>':
		if lx.getNextChar() == '=' {
			lx.readChar()
			return token.Token{
				Type:    token.GE,
				Literal: ">=",
			}
		} else {
			nextTokenType = token.GT
		}
	case '+':
		nextTokenType = token.PLUS
	case '-':
		nextTokenType = token.MINUS
	case '*':
		nextTokenType = token.ASTERISK
	case '/':
		nextTokenType = token.SLASH
	case '(':
		nextTokenType = token.LPAREN
	case ')':
		nextTokenType = token.RPAREN
	case '{':
		nextTokenType = token.LBRACE
	case '}':
		nextTokenType = token.RBRACE
	case ',':
		nextTokenType = token.COMMA
	case ';':
		nextTokenType = token.SEMICOLON
	case ':':
		nextTokenType = token.COLON
	default:
		if isValidLetter(lx.currentChar) {
			ident := lx.readIdentifier()
			return token.Token{
				Type:    token.GetTokenType(ident),
				Literal: ident,
			}
		} else if isNumber(lx.currentChar) {
			return token.Token{
				Type:    token.INT,
				Literal: lx.readInterger(),
			}
		} else {
			nextTokenType = token.ILLEGAL
		}
	}

	// get new token from char
	return token.NewToken(nextTokenType, lx.currentChar)

}
