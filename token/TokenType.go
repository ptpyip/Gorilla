package token

type TokenType string

const (
	// Special Tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT TokenType = "IDENT" // add, foobar, x, y, ...
	INT   TokenType = "INT"   // 134345

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	LT TokenType = "<"
	GT TokenType = ">"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LE     TokenType = "<="
	GE     TokenType = ">="

	AND TokenType = "&&"
	OR  TokenType = "||"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"

	// Keywords
	FUNCTION TokenType = "FN"
	LET      TokenType = "LET"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	// ELIF     TokenType = "ELIF"
	RETURN TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"let":   LET,
	"True":  TRUE,
	"False": FALSE,
	"if":    IF,
	"else":  ELSE,
	// "elif":   ELIF,
	"return": RETURN,
}

func GetTokenType(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENT
}

var PrefixOperatior = map[string]TokenType{
	"!": BANG,
	"-": MINUS,
}

// func isPrefixOperator(operator string) bool {
// 	_, ok := prefixOperatior[operator]
// 	return ok
// }
