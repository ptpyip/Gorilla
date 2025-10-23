package lexer

type Lexer struct {
	input       string
	pos         int
	nextPos     int
	currentChar byte
}

func NewLexer(input string) *Lexer {
	lx := &Lexer{input: input}
	lx.readChar()
	return lx
}

func (lx *Lexer) readChar() {
	lx.pos = lx.nextPos
	if lx.nextPos < len(lx.input) {
		lx.currentChar = lx.input[lx.pos]
	} else {
		lx.currentChar = 0x0
	}

	lx.nextPos += 1
}

func (lx *Lexer) getNextChar() byte {
	if lx.nextPos >= len(lx.input) {
		return 0x0
	}

	return lx.input[lx.nextPos]
}

func (lx *Lexer) readIdentifier() string {
	startPos := lx.pos
	for isValidLetter(lx.getNextChar()) {
		lx.readChar()
	}
	return lx.input[startPos : lx.pos+1]
}

func isValidLetter(inputChar byte) bool {
	return ('a' <= inputChar && inputChar <= 'z' ||
		'A' <= inputChar && inputChar <= 'Z' ||
		inputChar == '_')
}

func isNumber(inputChar byte) bool {
	return '0' <= inputChar && inputChar <= '9'
}

func (lx *Lexer) readInterger() string {
	startPos := lx.pos
	for isNumber(lx.getNextChar()) {
		lx.readChar()
	}
	return lx.input[startPos : lx.pos+1]
}

func (lx *Lexer) skip() {
	for lx.currentChar == ' ' ||
		lx.currentChar == '\t' ||
		lx.currentChar == '\n' ||
		lx.currentChar == '\r' {
		lx.readChar()
	}
}

func (lx *Lexer) Copy() *Lexer {
	return NewLexer(lx.input)
}
