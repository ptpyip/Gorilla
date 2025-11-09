package precedences

import "gorilla/token"

const (
	_ int = iota
	LOWEST
	// TRINARY
	EQUALS  // ==
	COMPARE // > or <
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // myFunction(X)
)

var Precedence = map[token.TokenType]int{
	// token.IF:       TRINARY,
	// token.ELSE:     TRINARY,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LE:       EQUALS,
	token.GE:       EQUALS,
	token.LT:       COMPARE,
	token.GT:       COMPARE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	// token.BANG:     PREFIX,
	// token.MINUS:    PREFIX,
	token.AND: PRODUCT,
	token.OR:  PRODUCT,
}
