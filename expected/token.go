package expected

import "gorilla/token"

type Token struct {
	ExpectedType    token.TokenType
	ExpectedLiteral string
}
